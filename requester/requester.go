package requester

import (
	"fmt"
	"github.com/mprokocki/interview-o/getclient"
	"log/slog"
	"net/http"
	"time"
)

type Configuration struct {
	Url            string
	RequestsAmount int
	Interval       time.Duration
}

type Requester struct {
	logger *slog.Logger
	client *getclient.GetMetadataHttpClient
}

func NewRequester(logger *slog.Logger, client *getclient.GetMetadataHttpClient) *Requester {
	if logger == nil {
		logger = slog.Default()
	}

	if client == nil {
		client = &getclient.GetMetadataHttpClient{
			&http.Client{
				Timeout: 10 * time.Second,
			},
		}
	}

	return &Requester{
		logger: logger,
		client: client,
	}
}

func DefaultConfiguration() *Configuration {
	return &Configuration{
		Url:            "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/?format=json",
		Interval:       10,
		RequestsAmount: 5,
	}
}

func (r *Requester) Run(config *Configuration) {
	if config == nil || config.Url == "" {
		config = DefaultConfiguration()
	}

	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	resChan := make(chan getclient.ResponseMetadata)
	for {
		select {
		case <-ticker.C:
			go func() {
				for i := 0; i < config.RequestsAmount; i++ {
					r.client.AsyncGet(config.Url, resChan)
				}
			}()
		case res := <-resChan:
			r.logger.Info("response", "res", fmt.Sprintf("%+v", res))
		}
	}
}
