package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"compiler"
	"types"
)

var PORT = 8090

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func parseNodeMapSchema(w http.ResponseWriter, req *http.Request) {

	type Body struct {
		Csv    string              `json:"csv"`
		Schema types.NodeMapSchema `json:"schema"`
	}
	type Response struct {
		Csv    string               `json:"csv"`
		Errors []types.NodeMapError `json:"errors"`
	}
	var body Body

	fmt.Println("===")

	sache, csvDecodeErr := ioutil.ReadAll(req.Body)

	if csvDecodeErr != nil {
		panic(csvDecodeErr)
	}
	fmt.Println(string(sache))

	fmt.Println("===")

	jsonDecodeErr := json.Unmarshal(sache, &body)

	if jsonDecodeErr != nil {
		fmt.Println("error decoding the json body")
		http.Error(w, jsonDecodeErr.Error(), http.StatusBadRequest)
		return
	}
	strReader := strings.NewReader(body.Csv)
	csvReader := csv.NewReader(strReader)
	data, csvDecodeErr := csvReader.ReadAll()

	if csvDecodeErr != nil {
		fmt.Println("error decoding the csv body")
		http.Error(w, csvDecodeErr.Error(), http.StatusBadRequest)
		return
	}
	// headerRow := data[0]
	dataRows := data[1:]

	var results [][]types.Primitive

	results = append(results, toPrimitiveSlice(body.Schema.Out.Columns))

	for _, row := range dataRows {

		results = append(results, compiler.ParseSache(types.DefaultNodeTypes, body.Schema, toPrimitiveSlice(row)))
	}

	csvbuffer := new(bytes.Buffer)
	csvwriter := csv.NewWriter(csvbuffer)
	csvEncodeErr := csvwriter.WriteAll(toStringSliceSlice(results))

	if csvEncodeErr != nil {
		log.Fatal(csvEncodeErr)
	}

	jsonBytes, jsonEncodeErr := json.Marshal(Response{Csv: csvbuffer.String(), Errors: []types.NodeMapError{}})

	if jsonEncodeErr != nil {
		panic(jsonEncodeErr)
	}
	fmt.Fprintf(w, string(jsonBytes))
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/calc", handleCors(parseNodeMapSchema))

	fmt.Printf("Listening on http://localhost:%v/\n", PORT)

	http.ListenAndServe(":"+strconv.Itoa(8090), nil)
}

func GetTestCsv() string {

	content, err := ioutil.ReadFile("data.csv")

	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func toPrimitiveSlice(data []string) []types.Primitive {

	var slice []types.Primitive

	for _, item := range data {

		slice = append(slice, types.Primitive(item))
	}
	return slice
}
func toStringSlice(data []types.Primitive) []string {
	var slice []string

	for _, item := range data {

		str, strOk := item.(string)
		float, floatOk := item.(float64)

		if strOk {
			slice = append(slice, str)

		} else if floatOk {
			slice = append(slice, fmt.Sprint(float))
			// slice = append(slice, strconv.FormatFloat(float, 'E', -1, 64))

		} else {

		}
	}
	return slice
}
func toStringSliceSlice(data [][]types.Primitive) [][]string {
	var slice [][]string

	for _, item := range data {

		slice = append(slice, toStringSlice(item))
	}
	return slice
}

func handleCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		sache := w.Header()

		sache.Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		sache.Set("Access-Control-Allow-Headers", "*")
		sache.Set("Content-Type", "application/json")

		if req.Method == "OPTIONS" {

			fmt.Println("preflight")

		} else {
			fmt.Println("request einkommend")
			h.ServeHTTP(w, req)
		}
	}
}
