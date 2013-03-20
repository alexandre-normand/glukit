package main

import (
	"parser"
	"net/http"
	"fmt"
	"io"
	"encoding/json"
)

type MeterRead struct {
    LocalTime string
    Value int
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("New Request\n")
	writer.WriteHeader(200)
	value := writer.Header()
	value.Add("Content-type", "application/json")


	fmt.Fprintf(writer, " { \"data\": [")
	parser.Parse("data.xml", writer, writeReadAsJson)
	fmt.Fprintf(writer, "{}] }")
}

func writeReadAsJson(meter *parser.Meter, writer io.Writer) {
	enc := json.NewEncoder(writer)
	meterRead := MeterRead{meter.DisplayTime, meter.Value}
	enc.Encode(meterRead)
	fmt.Fprintf(writer, ",")
}

