package main

import (
	"alertwest-interview-q1/lib"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// DelayRequest represents a request to delay a query execution
type DelayRequest struct {
	ID    uuid.UUID `json:"id"`
	Delay int       `json:"delay"`
}

// Server represents the HTTP server
type Server struct {
	mux *http.ServeMux
	db  *lib.DB
}

// NewServer creates a new HTTP server
func NewServer(db *lib.DB) *Server {
	server := &Server{
		mux: http.NewServeMux(),
		db:  db,
	}

	// Set up routes
	server.mux.HandleFunc("/queued", server.handleGetQueued)
	server.mux.HandleFunc("/resources", server.handleGetResources)
	server.mux.HandleFunc("/delay", server.handlePostDelay)

	return server
}

func (s *Server) Start(addr string) {
	log.Info().Str("Listening on", addr).Msg("Starting server")
	http.ListenAndServe(addr, s)
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleGetQueued handles GET /queued requests - currently polled every 5s on the client side.
// This returns the current list of queued (but not yet executed) queries.
func (s *Server) handleGetQueued(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Send request through the channel
	queued := s.db.GetQueued()

	if len(queued) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Create response array
	responses := make([]lib.QueuedQuery, 0, len(queued))

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

	// Send request through the channel
	metrics := s.db.GetResources()

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
	err := s.db.Delay(request.ID, request.Delay)
	if err != nil {
		http.Error(w, "Failed to delay query", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
