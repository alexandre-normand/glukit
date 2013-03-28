package parser

import (
	"encoding/xml"
	"io"
	"os"
	"bufio"
)

type Meter struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	Value        int    `xml:"Value,attr"`
}

func Parse(filepath string, writer io.Writer, handler func(*Meter)) {
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

	decoder := xml.NewDecoder(r)

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
			if se.Name.Local == "Meter" {
				var read Meter
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				decoder.DecodeElement(&read, &se)
				// Do some stuff with the page.
				handler(&read)
			}
		}
	}
}
