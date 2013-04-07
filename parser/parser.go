package parser

import (
	"encoding/xml"
	"io"
	"os"
	"bufio"
	"container/list"
)

type Meter struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	Value        int    `xml:"Value,attr"`
}

func Parse(filepath string) (reads *list.List) {
	// open input file
	fi, err := os.Open(filepath)
	if err != nil { panic(err) }
	// close fi on exit and check for its returned error
	defer func() {
		if fi.Close() != nil {
			panic(err)
		}
	}()
	// make a read buffer
	r := bufio.NewReader(fi)

	return ParseContent(r)
}

func ParseContent(reader io.Reader) (reads *list.List) {
	decoder := xml.NewDecoder(reader)
	reads = list.New()
	reads.Init()
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			// ...and its name is "page"
			if se.Name.Local == "Glucose" {
				var read Meter
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				decoder.DecodeElement(&read, &se)
				reads.PushBack(read)
			}
		}
	}

	return reads
}
