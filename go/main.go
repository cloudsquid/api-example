package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var file = flag.String("f", "", "file path to upload")

type CloudsquidClient struct {
	client   *http.Client
	apikey   string
	endpoint string
	sourceId string
}

type UploadRequest struct {
	Mimetype string `json:"mimetype,omitempty"`
	Filename string `json:"filename,omitempty"`
	Filetype string `json:"file_type,omitempty"`
	File     string `json:"file,omitempty"`
}

type UploadResponse struct {
	FileId string `json:"file_id"`
}

type RunRequest struct {
	FileId   string `json:"file_id"`
	Pipeline string `json:"pipeline"`
}
type RunResponse struct {
	RunId string `json:"run_id"`
}

type GetStatusRequest struct {
	RunId string `json:"run_id"`
}

type StatusResponse struct {
	Status string `json:"status"`
	Result any    `json:"result"`
}

func main() {
	config, err := Load("./config.env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	flag.Parse()

	if *file == "" {
		flag.Usage()
		return
	}

	client := CloudsquidClient{
		client:   &http.Client{},
		apikey:   config.CsKey,
		endpoint: config.CsEndpoint,
		sourceId: config.CsSourceID,
	}

	f, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	fileName := path.Base(*file)
	// file content must be base64 encoded
	fileContent := base64.StdEncoding.EncodeToString(f)

	uploadPayload := UploadRequest{
		Mimetype: "application/pdf",
		Filename: fileName,
		Filetype: "binary",
		File:     fileContent,
	}

	uploadResponse := client.UploadFile(uploadPayload)
	var uploadResponseBody UploadResponse
	err = json.NewDecoder(uploadResponse.Body).Decode(&uploadResponseBody)
	if err != nil {
		panic(err)
	}

	printResponse(uploadResponse, uploadResponseBody)

	runPayload := RunRequest{
		FileId:   uploadResponseBody.FileId,
		Pipeline: "cloudsquid-flash",
	}

	runResponse := client.RunFile(runPayload)
	var runResponseBody RunResponse
	err = json.NewDecoder(runResponse.Body).Decode(&runResponseBody)
	if err != nil {
		log.Printf("response body: %s", runResponse)
		panic(err)
	}
	printResponse(runResponse, runResponseBody)

	var extraction []byte
	// Polling for status
	for {

		statusResponse := client.GetStatus(runResponseBody.RunId)

		var statusResponseBody StatusResponse
		err = json.NewDecoder(statusResponse.Body).Decode(&statusResponseBody)
		if err != nil {
			panic(err)
		}
		printResponse(statusResponse, statusResponseBody)

		// Check if the status is "done", break the loop if true
		if statusResponseBody.Status == "done" {
			bs, err := json.MarshalIndent(statusResponseBody.Result, "", "  ")
			if err != nil {
				panic(err)
			}
			extraction = bs
			break
		}
		if statusResponseBody.Status == "error" {
			log.Fatalf("Error in processing: %s", statusResponseBody.Result)
		}

		// Optional: Delay to avoid too many requests in a short time
		time.Sleep(2 * time.Second) // Adjust the time as needed
	}

	fmt.Printf("Final result: \n%s\n", string(extraction))
}

func (client *CloudsquidClient) UploadFile(req UploadRequest) *http.Response {
	log.Print("Uploading file")

	baseURL, err := url.Parse(client.endpoint)
	if err != nil {
		panic(err)
	}
	baseURL.Path = path.Join(baseURL.Path, "datasources", client.sourceId, "documents")

	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	uploadRequest, err := http.NewRequest(http.MethodPost, baseURL.String(), bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	uploadResponse, err := client.doRequest(uploadRequest)
	if err != nil {
		panic(err)
	}

	log.Print("Successfully sent out upload request")

	return uploadResponse
}

func (c *CloudsquidClient) RunFile(
	req RunRequest,
) *http.Response {
	log.Print("running file")
	defer log.Print("successfully sent out request to run file")

	runEndpoint, err := url.Parse(c.endpoint)
	if err != nil {
		panic(err)
	}
	runEndpoint.Path = path.Join(runEndpoint.Path, "datasources", c.sourceId, "run")

	bs, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	log.Printf("running file: %s", runEndpoint.String())
	runRequest, err := http.NewRequest(http.MethodPost, runEndpoint.String(), bytes.NewReader(bs))
	if err != nil {
		panic(err)
	}

	runResponse, err := c.doRequest(runRequest)
	if err != nil {
		panic(err)
	}

	return runResponse
}

func (c *CloudsquidClient) GetStatus(runID string) *http.Response {
	log.Print("Getting status")

	statusURL, err := url.Parse(c.endpoint)
	if err != nil {
		panic(err)
	}
	statusURL.Path = path.Join(statusURL.Path, "datasources", c.sourceId, "run", runID)

	statusGetRequest, err := http.NewRequest(http.MethodGet, statusURL.String(), nil)
	if err != nil {
		panic(err)
	}

	statusResponse, err := c.doRequest(statusGetRequest)
	if err != nil {
		panic(err)
	}

	log.Print("successfully got status of file")

	return statusResponse
}

func printResponse(response *http.Response, body any) {
	fmt.Printf("Response: %#v\n", body)
	fmt.Printf("ResponseCode: %d\n", response.StatusCode)
}

func (c *CloudsquidClient) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apikey)

	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request %#v: %w", req, err)
	}

	return response, nil
}
