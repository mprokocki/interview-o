package main

import (
	"flag"
	"fmt"
	"github.com/mprokocki/interview-o/rates"
	"github.com/mprokocki/interview-o/requester"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "rates":
		notInRangeRates, err := rates.NotInRange(4.5, 4.7)
		if err != nil {
			panic(err)
		}

		for date, rate := range notInRangeRates {
			fmt.Printf("Date: %s, rate: %.02f\n", date.Format("2006-01-02"), rate)
		}

	default:
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		logger, closeFN, err := CreateLogger()
		if err != nil {
			panic(fmt.Errorf("unable to create logger: %w", err))
		}
		defer closeFN()

		rqster := requester.NewRequester(logger, nil)

		go rqster.Run(CreateConfig())
		fmt.Println("Script is running\n")
		<-sigCh
		fmt.Println("Killing the script\n")
	}
}

func CreateLogger() (*slog.Logger, func(), error) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, func() {}, err
	}

	closeFN := func() {
		defer file.Close()
	}

	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdin, file), nil))
	return logger, closeFN, nil
}

func CreateConfig() *requester.Configuration {
	url := flag.String("url", "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/?format=json", "Url where Get request will be made")
	req := flag.Int("req", 1, "Request amount in interval")
	interval := flag.Int("interval", 1, "Requests interval in seconds")
	flag.Parse()

	return &requester.Configuration{
		Url:            *url,
		RequestsAmount: *req,
		Interval:       time.Duration(*interval) * time.Second,
	}
}
