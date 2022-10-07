package services

import (
	"github.com/HorizenOfficial/rosetta-zen/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLogLevel(t *testing.T) {
	logger, err := utils.NewLogger()
	assert.Nil(t, err)
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))
	assert.Equal(t, zap.DebugLevel, utils.Atom.Level())
	loggerAPIService := NewLoggerAPIService()
	req := httptest.NewRequest(http.MethodGet, "/log/level", nil)
	w := httptest.NewRecorder()
	loggerAPIService.LogLevel(w, req)
	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			assert.True(t, false)
		}
	}(res.Body)
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"level\":\"debug\"}\n", string(data))
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))
	assert.Equal(t, zap.DebugLevel, utils.Atom.Level())
}

func TestPutLogLevel(t *testing.T) {
	logger, err := utils.NewLogger()
	assert.Nil(t, err)
	assert.True(t, logger.Core().Enabled(zap.DebugLevel))
	assert.Equal(t, zap.DebugLevel, utils.Atom.Level())
	loggerAPIService := NewLoggerAPIService()
	req := httptest.NewRequest(http.MethodPut, "/log/level", strings.NewReader("{\"level\":\"info\"}"))
	w := httptest.NewRecorder()
	loggerAPIService.LogLevel(w, req)
	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			assert.True(t, false)
		}
	}(res.Body)
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "{\"level\":\"info\"}\n", string(data))
	assert.True(t, logger.Core().Enabled(zap.InfoLevel))
	assert.Equal(t, zap.InfoLevel, utils.Atom.Level())
}
