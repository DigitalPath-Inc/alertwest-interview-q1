package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	ID    string `json:"id"`
	Delay int    `json:"delay"`
}

// Server represents the HTTP server
type Server struct {
	queue    *Queue
	monitor  *Monitor
	mux      *http.ServeMux
	tickrate int
}

// NewServer creates a new HTTP server
func NewServer(queue *Queue, monitor *Monitor, tickrate int) *Server {
	server := &Server{
		queue:    queue,
		monitor:  monitor,
		tickrate: tickrate,
		mux:      http.NewServeMux(),
	}

	// Set up routes
	server.mux.HandleFunc("/queued", server.handleGetQueued)
	server.mux.HandleFunc("/resources", server.handleGetResources)
	server.mux.HandleFunc("/delay", server.handlePostDelay)

	return server
}

func (s *Server) Start(addr string) {
	log.Printf("Starting server on %s", addr)
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

	// Get all queued queries
	queued := s.queue.GetQueued()
	if len(queued) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Create response array
	responses := make([]QueuedQuery, 0, len(queued))

	// Process each queued query
	for _, query := range queued {
		response := QueuedQuery{}
		response.Query.ID = query.id.String()
		response.Execution.ID = uuid.New().String()
		offset := time.Duration(float64(query.delay)/float64(s.tickrate)) * time.Second
		response.Execution.Timestamp = time.Now().Add(offset).UnixMilli()
		responses = append(responses, response)
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

	lastUpdate, lastCpu, lastMemory, lastIo := s.monitor.GetResourceUsage()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResourceMetrics{
		CPU:       lastCpu,
		IO:        lastIo,
		Memory:    lastMemory,
		Timestamp: lastUpdate.UnixMilli(),
	})
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
	if request.ID == "" {
		http.Error(w, "Missing execution ID", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(request.ID)
	if err != nil {
		http.Error(w, "Invalid execution ID", http.StatusBadRequest)
		return
	}

	err = s.queue.Delay(id, request.Delay)
	if err != nil {
		http.Error(w, "Execution not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
