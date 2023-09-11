package getclient

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type ResponseMetadata struct {
	IsJsonContentType bool
	StatusCode        int
	Time              time.Duration
	IsValidJson       bool
	Content           string
}

type GetMetadataHttpClient struct {
	Client *http.Client
}

func NewClient(client *http.Client) *GetMetadataHttpClient {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &GetMetadataHttpClient{Client: client}
}

func (r *GetMetadataHttpClient) Get(url string) (ResponseMetadata, error) {
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Technical interview golang script")
	res, err := r.Client.Do(req)
	if err != nil {
		return ResponseMetadata{}, err
	}

	defer res.Body.Close()
	duration := time.Since(start)
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return ResponseMetadata{}, err
	}

	return ResponseMetadata{
		IsJsonContentType: strings.Contains(res.Header.Get("Content-Type"), "application/json"),
		StatusCode:        res.StatusCode,
		Time:              duration,
		IsValidJson:       json.Valid(content),
		Content:           string(content),
	}, nil
}

func (r *GetMetadataHttpClient) AsyncGet(url string, resChan chan ResponseMetadata) {
	go func() {
		resMetadata, err := r.Get(url)
		if err != nil {
			panic(err)
		}

		resChan <- resMetadata
	}()
}
