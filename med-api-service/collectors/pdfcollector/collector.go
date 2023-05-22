package pdfcollector

import (
	"archive/tar"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetPDFByPMCID gets the pdf link from the pubmed pdf api and if successful,
// returns the pdf in string format or an error if the pdf link retrieval was unsuccessful.
func GetPDFByPMCID(pmcid string) (string, error) {
	const baseURL = "https://www.ncbi.nlm.nih.gov/pmc/utils/oa/oa.fcgi?id="

	finalURL := baseURL + pmcid

	response, err := http.Get(finalURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data := new(OA)

	err = xml.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return "", err
	}

	if data.Error != "" {
		return "", errors.New("free pdf for provided pmcid could not be retrieved with error: " + data.Error)
	}

	link, fromGzip := getLinkFromRecords(data.RecordList.Records)

	return getPDF(link, fromGzip)
}

// getLinkFromRecords gets the link that has the pdf data from the provided records and
// returns the link and whether the link directs to a gzip download
func getLinkFromRecords(records []Record) (string, bool) {
	link := ""
	fromGzip := true

	for _, record := range records {
		link = record.Link.Value
		// direct pdf link is preferred, otherwise it's going to be a gzip file
		if record.Link.Format == "pdf" {
			fromGzip = false
			break
		}
	}

	return link, fromGzip
}

// getPDF gets the pdf from the provided link.
// fromGzip needs to be provided to specify whether the
// link is a gzip link (otherwise a pdf link is assumed),
// in order to retrieve the pdf appropriatly
func getPDF(link string, fromGzip bool) (string, error) {
	httpsLink := strings.Replace(link, "ftp", "https", 1)

	response, err := http.Get(httpsLink)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if fromGzip {
		return getPdfFromGzip(response.Body)
	}

	return convertPDFToText(response.Body)
}

// getPdfFromGzip retrieves the pdf inside of a io.Reader that is
// a compressed .tar.gz file. It returns the pdf as a string or an error
func getPdfFromGzip(r io.Reader) (string, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	tarHeader := tar.NewReader(gzipReader)

	for {
		header, err := tarHeader.Next()
		if err != nil {
			// Reached the end of the archive
			if err == io.EOF {
				break
			}

			return "", err
		}

		if filepath.Ext(header.Name) == ".pdf" {
			return convertPDFToText(tarHeader)
		}
	}

	return "", errors.New("no pdf found in tar.gz file")
}

// convertPDFToText converts the pdf file provided as an io.Reader to string
// utilizing the Linux 'pdftotext' commandline utility. The function returns the
// pdf as string and potentially an error.
func convertPDFToText(r io.Reader) (string, error) {
	f, err := os.CreateTemp(os.TempDir(), "med_api_service*")
	if err != nil {
		return "", err
	}

	cleanupTempFile := func() {
		f.Close()
		if err = os.Remove(f.Name()); err != nil {
			log.Fatal("CRITICAL: Cannot remove temporary file from pdftotext conversion with error:", err)
		}
	}
	defer cleanupTempFile()

	_, err = io.Copy(f, r)
	if err != nil {
		return "", err
	}

	data, err := exec.Command("pdftotext", "-q", "-nopgbrk", "-enc", "UTF-8", "-eol", "unix", f.Name(), "-").Output()

	return string(data), err
}
