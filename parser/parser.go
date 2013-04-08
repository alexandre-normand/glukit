package parser

import (
	"encoding/xml"
	"io"
	"os"
	"bufio"
	"container/list"
	"models"
	"utils"
)

type Meter struct {
	InternalTime string `xml:"InternalTime,attr"`
	DisplayTime  string `xml:"DisplayTime,attr"`
	Value        int    `xml:"Value,attr"`
}

func Parse(filepath string) (reads []models.MeterRead) {
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

func ParseContent(reader io.Reader) (reads []models.MeterRead) {
	decoder := xml.NewDecoder(reader)
	readsList := list.New()
	readsList.Init()
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
				if (read.Value > 0) {
					readsList.PushBack(read)
				}
			}
		}
	}

	return ConvertAsReadsArray(readsList)
}

func ConvertAsReadsArray(meterReads *list.List) (reads []models.MeterRead) {
	reads = make([]models.MeterRead, meterReads.Len())
	for e, i := meterReads.Front(), 0; e != nil; e, i = e.Next(), i + 1 {
		meter := e.Value.(Meter)
		reads[i] = models.MeterRead{meter.DisplayTime, utils.GetTimeInSeconds(meter.DisplayTime), meter.Value}
	}

	return reads
}
