package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// ResourceMetrics represents the resource utilization metrics
type ResourceMetrics struct {
	CPU       ResourceUsage `json:"cpu"`
	IO        ResourceUsage `json:"io"`
	Memory    ResourceUsage `json:"memory"`
	Timestamp int64         `json:"timestamp"`
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
	ID    uuid.UUID `json:"id"`
	Delay int       `json:"delay"`
}

// Server represents the HTTP server
type Server struct {
	mux              *http.ServeMux
	monitorEventChan <-chan *QueuedQuery
	getQueuedChan    chan<- GetQueuedRequest
	getResourcesChan chan<- GetResourcesRequest
	delayChan        chan<- DelayRequest
}

type GetQueuedRequest struct {
	ResponseChan chan<- []*QueuedQuery
}

type GetResourcesRequest struct {
	ResponseChan chan<- *ResourceMetrics
}

// NewServer creates a new HTTP server
func NewServer(monitorEventChan <-chan *QueuedQuery, getQueuedChan chan<- GetQueuedRequest, getResourcesChan chan<- GetResourcesRequest, delayChan chan<- DelayRequest) *Server {
	server := &Server{
		mux:              http.NewServeMux(),
		monitorEventChan: monitorEventChan,
		getQueuedChan:    getQueuedChan,
		getResourcesChan: getResourcesChan,
		delayChan:        delayChan,
	}

	// Set up routes
	server.mux.HandleFunc("/queued", server.handleGetQueued)
	server.mux.HandleFunc("/resources", server.handleGetResources)
	server.mux.HandleFunc("/delay", server.handlePostDelay)

	return server
}

func (s *Server) Start(addr string) {
	log.Printf("Starting server on %s", addr)
	// Start a goroutine to consume events from monitorEventChan
	go func() {
		for event := range s.monitorEventChan {
			_ = event
		}
	}()
	http.ListenAndServe(addr, s)
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleGetQueued handles GET /queued requests
func (s *Server) handleGetQueued(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a response channel
	responseChan := make(chan []*QueuedQuery, 1)

	// Send request through the channel
	s.getQueuedChan <- GetQueuedRequest{
		ResponseChan: responseChan,
	}

	// Get the response
	queued := <-responseChan

	if len(queued) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Create response array
	responses := make([]QueuedQuery, 0, len(queued))

	// Process each queued query
	for _, query := range queued {
		responses = append(responses, *query)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// handleGetResources handles GET /resources requests
func (s *Server) handleGetResources(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a response channel
	responseChan := make(chan *ResourceMetrics, 1)

	// Send request through the channel
	s.getResourcesChan <- GetResourcesRequest{
		ResponseChan: responseChan,
	}

	// Get the response
	metrics := <-responseChan

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// handlePostDelay handles POST /delay requests
func (s *Server) handlePostDelay(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request DelayRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if request.ID == uuid.Nil {
		http.Error(w, "Missing execution ID", http.StatusBadRequest)
		return
	}

	// Send request through the channel
	s.delayChan <- request

	w.WriteHeader(http.StatusOK)
}
