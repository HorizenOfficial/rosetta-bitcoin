package services

import (
	"github.com/HorizenOfficial/rosetta-zen/utils"
	"net/http"
)

// LoggerAPIService implements the LoggerAPIServicer interface.
type LoggerAPIService struct{}

// NewLoggerAPIService creates a new instance of a LoggerAPIService.
func NewLoggerAPIService() *LoggerAPIService {
	return &LoggerAPIService{}
}

// LogLevel implements the /log/level endpoint.
func (s *LoggerAPIService) LogLevel(w http.ResponseWriter, r *http.Request) {
	utils.Atom.ServeHTTP(w, r)
}
