package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"search-service/caller"
	"search-service/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ctxTimeout is the set timeout for every mongodb operation
const ctxTimeOut = 15 * time.Second

var client *mongo.Client

// NewConn gets the db connection from the main function
func NewConn(mongo *mongo.Client) {
	client = mongo
}

// InsertInto inserts a  models.DataEntry item into the appropriate collection
// and returns the hex id and potentially an error
func InsertInto(collName string, entry models.DataEntry) (string, error) {
	collection := client.Database("search").Collection(collName)

	entry.AddDefaultData()

	res, err := collection.InsertOne(context.Background(), entry)
	if err != nil {
		log.Println("Error inserting into searches:", err)
		return "", err
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", nil
	}

	return id.Hex(), nil
}

// SearchEntriesByKeyword queries the mongodb with the provided query keyword and the SitesToSearch,
// if it doesn't find a suitable entries, it calles the appropriate scraper to collect it
// and returns a slice of SearchEntries and potentially an error
func SearchEntriesByKeyword(query *models.SearchQuery) ([]*models.SearchEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	sitesLen := len(query.SitesToSearch)

	if sitesLen < 1 {
		return nil, errors.New("no sites to search given, search cancelled")
	}

	if query.Keyword == "" {
		return nil, errors.New("no keywords to perform search")
	}

	results := make([]*models.SearchEntry, sitesLen)

	wg := new(sync.WaitGroup)

	wg.Add(sitesLen)

	// Doing this concurrently saves a lot of time if the requested entries based on keyword
	// do not already exist in mongo and need to be collected, but does slow down the collection for a couple of
	// milliseconds if the entries already exist in mongo, due to added overhead of launching new threads.
	// There's probably a better way of doing this as to only launch threads when entry does not exist,
	// but so far I haven't found a good way of doing that.
	//
	// Note: Capturing i and site in closure because their values change through each iteration
	getResult := func(i int, site string) {
		defer wg.Done()

		result, err := searchForKeyword(ctx, query.Keyword, site)
		if err != nil {
			log.Printf("Failed to fetch result for site %s with error: %s\n", site, err)
			return
		}

		results[i] = result

		// Doing this in the background since it doesn't affect the final results,
		// nor is there a returned value or error to be handled
		go checkForUpdate(result, site)

	}

	for i, site := range query.SitesToSearch {
		go getResult(i, site)
	}
	
	wg.Wait()

	return results, nil
}

// SearchForPDF queries mongdb for the requested pdf based on the keyword(PMID)
// and returns a models.SearchEntry and potentially and error
func SearchForPDF(query *models.SearchQuery) (*models.PDFEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	collection := client.Database("search").Collection("pdf_logs")

	result := new(models.PDFEntry)

	err := collection.FindOne(ctx, bson.M{"pmid": query.Keyword}).Decode(result)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			err = fmt.Errorf("could not decode pdf entry result for PMCID %s with error: %s", query.Keyword, err.Error())
			log.Println(err)
			return nil, err
		}

		result, err = caller.RequestPDFEntry(query.Keyword)
		if err != nil {
			return nil, err
		}

		result.ID, err = InsertInto("pdf_logs", result)
		if err != nil {
			log.Printf("Could not insert pdf entry result with PMCID: %s and error: %s\n", query.Keyword, err.Error())
		}
	}

	return result, nil
}

// GetOneByID queries mongodb for one SearchEntry with the given id
// and returns a SearchEntry and potentially an error
func GetOneSearchEntryByID(id, site string) (*models.SearchEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	collection := client.Database("search").Collection("search_logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("could not get objectID with ID: %s due to error: %s", id, err.Error())
	}

	entry := new(models.SearchEntry)

	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(entry)
	if err != nil {
		return nil, fmt.Errorf("could not decode search entry with error: %s", err.Error())
	}

	return entry, nil
}

// DropCollection drops a given collection
func DropCollection(collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	collection := client.Database("search").Collection(collectionName)

	err := collection.Drop(ctx)

	return err
}

// UpdateSearchEntry updates a mongodb collection according to the SearchEntry.ID
// and returns the mongo.UpdateResult and potentially an error
func UpdateSearchEntry(s *models.SearchEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	collection := client.Database("search").Collection("search_logs")

	docID, err := primitive.ObjectIDFromHex(s.ID)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"data", s.Data},
				{"updated_at", time.Now()},
			}},
		},
	)
	return err
}

// DeleteByID deletes a entry in the given collcetions
// with the correspoding id and potentially returns an error
func DeleteByIDIn(collectionName, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	collection := client.Database("search").Collection(collectionName)

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("could not get objectID with ID: %s due to error: %s", id, err.Error())
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": docID})

	return err
}

// searchForKeyword performs a query in mongo for the search entry with given keyword and site,
// if it succeeds it just returns the result, if not, it checks if error was ErrNoDocuments,
// which means no documents were found, so it requests new data to insert to the collection and
// return from the appropriate scraper
func searchForKeyword(ctx context.Context, keyword, site string) (*models.SearchEntry, error) {
	collection := client.Database("search").Collection("search_logs")

	result := new(models.SearchEntry)

	err := collection.FindOne(ctx, bson.M{"keyword": keyword, "origin": site}).Decode(result)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			err = fmt.Errorf("could not decode search result for %s with error: %s", keyword, err.Error())
			log.Println(err)
			return nil, err
		}

		result, err = caller.RequestSearchEntry(keyword, site)
		if err != nil {
			return nil, err
		}
		result.ID, err = InsertInto("search_logs", result)
		if err != nil {
			log.Printf("Could not insert search result with keyword: %s and error: %s\n", keyword, err.Error())
		}

	}

	return result, nil
}

// checkForUpdate checks if the entry needs to be updated
// and updates it if nessecary
func checkForUpdate(entry *models.SearchEntry, site string) {
	const (
		maxDays    = 7
		hoursInDay = 24
	)

	dif := time.Since(entry.UpdatedAt)

	days := int(dif.Hours() / hoursInDay)

	if days < maxDays {
		return
	}

	entry, err := caller.RequestSearchEntry(entry.Keyword, entry.Origin)
	if err != nil {
		log.Println("Could not get new entry for update from scraper with error:", err.Error())
		return
	}

	err = UpdateSearchEntry(entry)
	if err != nil {
		log.Println("Entry update failed due to error:", err)
		err = DeleteByIDIn("search_logs", entry.ID)
		if err != nil {
			log.Println("Entry delete failed due to error:", err)
		}
	}
}
