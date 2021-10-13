package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/ghodss/yaml"
	"github.com/thehivecorporation/log"
	"github.com/tidwall/pretty"
)

var (
	elasticsearchUrl  = kingpin.Flag("elasticsearch", "URL of Elasticsearch. Scheme and port are mandatory (but HTTPS is not tested xD)").Default("http://localhost:9200").Short('e').String()
	pipelineName      = kingpin.Flag("pipeline-name", "Name of the pipeline to create on elasticsearch").Default("testing").Short('p').String()
	ugly              = kingpin.Flag("ugly", "Deactivate pretty printing").Short('u').Bool()
	bulk              = kingpin.Flag("bulk", "Use bulk if you want to make a single request to ES and get a single response with many documents").Short('b').Default("false").Bool()
	skip              = kingpin.Flag("skip", "Skip the N first documents").Default("0").Short('s').Int()
	total             = kingpin.Flag("total", "Process a total of N documents. You can use it with skip").Short('n').Default("0").Int()
	useStdin          = kingpin.Flag("stdin", "Instead of using an input log file, read from stdin").Short('i').Bool()
	inputPipelineFile = kingpin.Arg("pipeline", "File path of the pipeline to use").Required().String()
	inputLogFile      = kingpin.Arg("logs", "The log filepath to parse. Use stdin flag if you want to work just on the cli").String()
)

func main() {
	kingpin.Parse()

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	insertPipelineToElasticsearch(client)

	var reader io.Reader = os.Stdin
	if !*useStdin {
		file, err := os.Open(*inputLogFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = file
	}

	if *bulk {
		simulateBulk(client, reader)
	} else {
		simulate(client, reader)
	}
}

// simulate will make a different HTTP request for each line in the incoming log file. It uses the Simulate API from
// Elasticsearch https://www.elastic.co/guide/en/elasticsearch/reference/master/simulate-pipeline-api.html
func simulate(client http.Client, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	var i, t int
	for scanner.Scan() {
		if *skip != 0 && i < *skip {
			i++
			continue
		}

		if *total != 0 && t >= *total {
			break
		}
		t++

		text := scanner.Text()

		//This "uglyness" handles JSON input and makes kittens cry
		var message interface{} = text
		temp := map[string]interface{}{}
		if err := json.Unmarshal([]byte(text), &temp); err == nil {
			//This is a JSON
			message = temp
		}

		docs := make(map[string]interface{})
		docs["docs"] = []SimulateDoc{{
			Source: map[string]interface{}{"message": message},
		}}

		doRequest(*elasticsearchUrl, *pipelineName, docs, &client, *ugly)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// simulateBulk will build a single request to elasticsearch with all the incoming log lines in a single request. It
// uses the Simulate API https://www.elastic.co/guide/en/elasticsearch/reference/master/simulate-pipeline-api.html
func simulateBulk(client http.Client, reader io.Reader) {
	simulateDocs := make([]SimulateDoc, 0)

	scanner := bufio.NewScanner(reader)
	var i, t int
	for scanner.Scan() {
		if *skip != 0 && i < *skip {
			i++
			continue
		}

		if *total != 0 && t >= *total {
			break
		}
		t++

		simulateDocs = append(simulateDocs, SimulateDoc{
			Source: map[string]interface{}{"message": scanner.Text()},
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	docs := make(map[string]interface{})
	docs["docs"] = simulateDocs

	doRequest(*elasticsearchUrl, *pipelineName, docs, &client, *ugly)
}

// doRequest takes the elasticsearch ready payload and does a single request to Elasticsearch. Panics on error.
func doRequest(elasticsearchURL, pipelineName string, outputData interface{}, client *http.Client, uglyPrint bool) {
	outputBodyBytes, err := json.Marshal(outputData)
	if err != nil {
		log.Fatal(err)
	}

	// doRequest is shared between `simulate` and `simulateBulk` so this block might be used repeatedly to always
	// get the same result. Whatever...
	elasticsearchURL = strings.Trim(elasticsearchURL, "/")
	urlS := elasticsearchURL + "/_ingest/pipeline/" + pipelineName + "/_simulate"
	docCreationURL, err := url.Parse(urlS)
	if err != nil {
		log.WithError(err).WithField("url", urlS).Fatal("error parsing URL")
	}

	req := buildRequest(docCreationURL, ioutil.NopCloser(bytes.NewReader(outputBodyBytes)))
	req.Method = http.MethodPost
	res, err := client.Do(req)
	if err != nil {
		log.WithError(err).Fatal("Error doing request")
	}
	defer res.Body.Close() //Ignore error

	if res.StatusCode/100 != 2 {
		log.WithField("status_code", res.StatusCode).Fatal("Error doing request")
	}

	inputBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Fatal("error reading body of response")
	}

	// Do not pretty print output. Useful to pipe data in bash
	if uglyPrint {
		fmt.Printf(string(inputBody))
	} else {
		log.Infof("Response:\n%s\n", string(pretty.Color(pretty.Pretty(pretty.PrettyOptions(inputBody, &pretty.Options{
			Width:    120,
			Prefix:   "",
			Indent:   "  ",
			SortKeys: true,
		})), nil)))
	}
}

// insertPipelineToElasticsearch reads the incoming pipeline file and inserts it into Elasticsearch as a mandatory
// step before doing the parsing step. Existing pipelines with the same name will be overridden
func insertPipelineToElasticsearch(client http.Client) {
	pipelineFile, err := os.Open(*inputPipelineFile)
	if err != nil {
		log.Fatal(err)
	}
	defer pipelineFile.Close()

	*elasticsearchUrl = strings.Trim(*elasticsearchUrl, "/")
	pipelineCreationUrl, err := url.Parse(*elasticsearchUrl + "/_ingest/pipeline/" + *pipelineName)
	if err != nil {
		log.Fatal(err)
	}

	byt, err := ioutil.ReadAll(pipelineFile)
	if err != nil {
		log.Fatal(err)
	}

	jsonPipelineByt, err := yaml.YAMLToJSON(byt)
	if err != nil {
		log.Fatal(err)
	}

	req := buildRequest(pipelineCreationUrl, ioutil.NopCloser(bytes.NewReader(jsonPipelineByt)))
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if *ugly {
		fmt.Println(string(body))
	} else {
		log.Infof("Pipeline insertion output:\n%s", string(pretty.Color(pretty.Pretty(body), nil)))
	}

	if res.StatusCode != 200 {
		log.Fatalf("Error writing pipeline")
	}
}

func buildRequest(u *url.URL, r io.ReadCloser) *http.Request {
	return &http.Request{
		Method: http.MethodPut,
		URL:    u,
		Body:   r,
		Header: http.Header{
			"Content-Type":    []string{"application/json"},
			"Accept":          []string{"application/json; charset=UTF-8"},
			"Accept-encoding": []string{"gzip"},
		},
	}
}
