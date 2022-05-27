package main

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	st "github.com/montanaflynn/stats"
)

type Data struct {
	DataArr []float64 `json:"data"`
}

func findThroughput(component string, data []float64) {
	// throughput is no of req/sec
	mean, _ := st.Mean(data)
	median, _ := st.Median(data)
	sd, _ := st.StandardDeviation(data)
	max, _ := st.Max(data)
	min, _ := st.Min(data)
	perc, _ := st.Percentile(data, float64(90))

	log.Printf("%s => mean : %.5f req/sec", component, mean)
	log.Printf("%s => median : %.5f req/sec", component, median)
	log.Printf("%s => sd : %.5f req/sec", component, sd)
	log.Printf("%s => min : %.5f req/sec", component, min)
	log.Printf("%s => max : %.5f req/sec", component, max)

	log.Printf("%s => 90th Percentile : %.5f req/sec", component, perc)
}

func findRTT(component string, data []float64) {
	mean, _ := st.Mean(data)
	median, _ := st.Median(data)
	sd, _ := st.StandardDeviation(data)
	max, _ := st.Max(data)
	min, _ := st.Min(data)
	perc, _ := st.Percentile(data, float64(90))

	log.Printf("%s => mean : %.5f secs", component, mean)
	log.Printf("%s => median : %.5f secs", component, median)
	log.Printf("%s => sd : %.5f secs", component, sd)
	log.Printf("%s => min : %.5f secs", component, min)
	log.Printf("%s => max : %.5f secs", component, max)

	log.Printf("%s => 90th Percentile : %.5f secs", component, perc)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {

	absPath, err := filepath.Abs("/Users/hysensw/Videos/lab5/lab5/lab5/gorumspaxos/cmd/paxosclient")
	if err != nil {
		log.Panicf("Error : %v\n", err)
	}

	res := findFile(absPath, findServer)
	dataServer := make([]float64, 0)
	for _, f := range res {
		data := readFile(f)
		dataServer = append(dataServer, data...)
	}

	log.Printf("Total replicas : %d", len(res))
	findThroughput("Server", dataServer)

	res = findFile(absPath, findClient)
	dataClient := make([]float64, 0)
	for _, f := range res {
		data := readFile(f)
		dataClient = append(dataClient, data...)
	}

	log.Printf("Total clients : %d", len(res))
	log.Printf("Total clients requests : %d", len(dataClient))

	findRTT("Client", dataClient)
}

// Read input from file
func readFile(filename string) []float64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicf("Error opening file : %v\n", err)
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("Error opening file : %v\n", err)
	}

	d := Data{}

	err = json.Unmarshal(fileBytes, &d)
	if err != nil {
		log.Panicf("Error: %v\n", err)
	}
	return d.DataArr
}

// find the files in the given path, it takes a function as callback to find the specific files
func findFile(dir string, fn func(string) bool) []string {
	var files []string
	filepath.WalkDir(dir, func(s string, d fs.DirEntry, e error) error {
		if fn(s) {
			files = append(files, s)
		}
		return nil
	})
	return files
}

// function to be passed as callback to select json file
func findClient(s string) bool {
	_, file := filepath.Split(s)
	return strings.HasPrefix(file, "client_")
}

// function to be passed as callback to select specific yaml files
func findServer(s string) bool {
	_, file := filepath.Split(s)
	return strings.HasPrefix(file, "server_")
}
