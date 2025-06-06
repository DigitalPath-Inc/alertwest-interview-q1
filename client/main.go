package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"alertwest-interview-q1/lib"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	serverURL = ""
)

// QueuedQuery represents a query in the queue
type QueuedQuery struct {
	Query struct {
		ID string `json:"id"`
	} `json:"query"`
	Execution struct {
		ID        string `json:"id"`
		Timestamp int64  `json:"timestamp"`
	} `json:"execution"`
}

// DelayRequest represents a request to delay a query execution
type DelayRequest struct {
	ID    string `json:"id"`
	Delay int    `json:"delay"`
}

// Set up the server URL from environment variable and ensure it is set
// before starting the client and log an error if it is not set
func init() {
	serverURL = os.Getenv("SERVER_URL")
	if serverURL == "" {
		log.Fatal().Msg("SERVER_URL is not set")
	}
}

// Run main function initializes the logger and starts the resource checking.
// It also starts checking the queued queries from the server.
func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting client")
	go checkResources()
	checkQueuedQueries()
}

// checkQueuedQueries periodically checks the queued queries from the server
// and logs the details of each query.
func checkQueuedQueries() {
	for {
		time.Sleep(time.Second * 1)
		resp, err := http.Get(serverURL + "/queued")
		if err != nil {
			log.Err(err).Msg("Error making request")
			continue
		}

		// Check if response is 204 No Content
		if resp.StatusCode == http.StatusNoContent {
			log.Warn().Msg("No queued queries available")
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close body in all cases

		if err != nil {
			log.Err(err).Msg("Error reading response")
			continue
		}

		var queuedQueries []QueuedQuery
		err = json.Unmarshal(body, &queuedQueries)
		if err != nil {
			log.Err(err).Msg("Error unmarshalling response")
			continue
		}

		log.Info().Int("Queue Size", len(queuedQueries)).Msg("Queries")
		for _, query := range queuedQueries {
			log.Debug().Interface("Query", query).Msg("Queued Query")
		}
	}
}

// checkResources periodically checks the server for resource metrics
// and logs the resource usage.
func checkResources() {
	for {
		time.Sleep(time.Second * 1)
		resp, err := http.Get(serverURL + "/resources")
		if err != nil {
			log.Err(err).Msg("Error making request")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close body in all cases

		if err != nil {
			log.Err(err).Msg("Error reading response")
			continue
		}

		var resourceMetrics lib.ResourceMetrics
		err = json.Unmarshal(body, &resourceMetrics)
		if err != nil {
			log.Err(err).Msg("Error unmarshalling response")
			continue
		}

		log.Info().EmbedObject(resourceMetrics).Msg("Resources")
	}
}
