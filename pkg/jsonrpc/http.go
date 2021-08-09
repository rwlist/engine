package jsonrpc

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var jsonContentType = "application/json"

// TransportHTTP tries to implement JSON-RPC 2.0 Transport: HTTP spec.
// http://www.simple-is-better.org/json-rpc/transport_http.html
type TransportHTTP struct {
	handler Handler
}

func NewHTTP(handler Handler) *TransportHTTP {
	return &TransportHTTP{
		handler: handler,
	}
}

func (h *TransportHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h.handleRequest(r)

	if err := resp.validate(); err != nil {
		log.WithError(err).Warn("invalid response")
	}

	w.Header().Set("Content-Type", jsonContentType)
	// TODO: set Content-Length
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.WithError(err).Error("failed to marshal response")
	}
}

func (h *TransportHTTP) handleRequest(r *http.Request) Response {
	resp := Response{
		Version: Version,
	}

	if !h.validateRequest(r) {
		resp.Error = &InvalidRequest
		return resp
	}

	var req Request
	// TODO: limit request length
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.WithError(err).Warn("failed to parse jsonrpc request")
		resp.Error = &ParseError
		return resp
	}

	// set header from request
	req.Authorization = r.Header.Get("Authorization")

	// set result
	resp.ID = req.ID
	resp.Result, resp.Error = h.handler(&req)

	return resp
}

func (h *TransportHTTP) validateRequest(r *http.Request) bool {
	if r.Method != "POST" {
		return false
	}

	header := r.Header

	if header.Get("Content-Type") != jsonContentType {
		return false
	}

	// TODO: verify Content-Length

	if header.Get("Accept") != jsonContentType {
		return false
	}

	return true
}
