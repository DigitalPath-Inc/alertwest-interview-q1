package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	serverURL = ""
)

// ResourceMetrics represents the resource utilization metrics
type ResourceMetrics struct {
	CPU struct {
		Average int `json:"average"`
		Min     int `json:"min"`
		Max     int `json:"max"`
	} `json:"cpu"`
	IO struct {
		Average int `json:"average"`
		Min     int `json:"min"`
		Max     int `json:"max"`
	} `json:"io"`
	Memory struct {
		Average int `json:"average"`
		Min     int `json:"min"`
		Max     int `json:"max"`
	} `json:"memory"`
	Timestamp int64 `json:"timestamp"`
}

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

func init() {
	serverURL = os.Getenv("SERVER_URL")
	if serverURL == "" {
		log.Fatal("SERVER_URL is not set")
	}
}

func main() {
	go checkResources()
	checkQueuedQueries()
}

func checkQueuedQueries() {
	for {
		time.Sleep(time.Second * 5)
		resp, err := http.Get(serverURL + "/queued")
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}

		// Check if response is 204 No Content
		if resp.StatusCode == http.StatusNoContent {
			fmt.Println("No queued queries available")
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close body in all cases

		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		var queuedQueries []QueuedQuery
		err = json.Unmarshal(body, &queuedQueries)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			continue
		}

		fmt.Println(len(queuedQueries), "queued queries")
	}
}

func checkResources() {
	for {
		time.Sleep(time.Second * 15)
		resp, err := http.Get(serverURL + "/resources")
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close body in all cases

		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		var resourceMetrics ResourceMetrics
		err = json.Unmarshal(body, &resourceMetrics)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			continue
		}

		fmt.Println(resourceMetrics)
	}
}
