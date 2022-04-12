package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"./compiler"
	"./types"
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
		Csv    string
		Schema types.NodeMapSchema
	}
	var body Body

	fmt.Println("===")

	sache, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}
	fmt.Println(string(sache))

	fmt.Println("===")

	ferr := json.Unmarshal(sache, &body)

	if ferr != nil {
		fmt.Println("error decoding the json body")
		http.Error(w, ferr.Error(), http.StatusBadRequest)

		fmt.Println(ferr.Error())
		return
	}
	strReader := strings.NewReader(body.Csv)
	csvReader := csv.NewReader(strReader)
	data, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println("error decoding the csv body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// headerRow := data[0]
	dataRows := data[1:]

	parsedData := compiler.ParseSache(types.DefaultNodeTypes, body.Schema, toPrimitiveSlice(dataRows[0]))

	fmt.Println(parsedData)

	fmt.Fprintf(w, "lol")
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/calc", parseNodeMapSchema)

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

	var sache []types.Primitive

	for i := range data {

		sache = append(sache, types.Primitive(i))

	}
	return sache
}
