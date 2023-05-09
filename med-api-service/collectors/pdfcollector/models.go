package pdfcollector

import "encoding/xml"

/*
*	The XML models that the pubmed pdf api sends back on request
*/

type OA struct {
	XMLName    xml.Name   `xml:"OA"`
	RecordList RecordList `xml:"records"`
	Error      string     `xml:"error"`
}

type RecordList struct {
	XMLName xml.Name `xml:"records"`
	Records []Record `xml:"record"`
}

type Record struct {
	XMLName xml.Name   `xml:"record"`
	Link    RecordLink `xml:"link"`
}

type RecordLink struct {
	XMLName xml.Name `xml:"link"`
	Value   string   `xml:"href,attr"`
	Format  string   `xml:"format,attr"`
}
