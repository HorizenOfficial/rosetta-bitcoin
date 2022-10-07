// Copyright 2022 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"net/http"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
)

// LoggerAPIServicer defines the api actions for the LoggerAPI service
type LoggerAPIServicer interface {
	LogLevel(w http.ResponseWriter, r *http.Request)
}

// A LoggerApiController binds http requests to an api service and writes the service results to
// the http response
type LoggerApiController struct {
	service  LoggerAPIServicer
	asserter *asserter.Asserter
}

// NewLoggerApiController creates a default api controller
func NewLoggerApiController(
	s LoggerAPIServicer,
	asserter *asserter.Asserter,
) server.Router {
	return &LoggerApiController{
		service:  s,
		asserter: asserter,
	}
}

// Routes returns all the api routes for the LoggerApiController
func (c *LoggerApiController) Routes() server.Routes {
	return server.Routes{
		{
			"LogLevel",
			strings.ToUpper("Get"),
			"/log/level",
			c.LogLevel,
		},
		{
			"UpdateLogLevel",
			strings.ToUpper("Put"),
			"/log/level",
			c.LogLevel,
		},
	}
}

// LogLevel Get the log level or update it
func (c *LoggerApiController) LogLevel(w http.ResponseWriter, r *http.Request) {
	c.service.LogLevel(w, r)
}
