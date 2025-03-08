package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/joseph0x45/arcane/server/logger"
)

func WriteData(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal data", data, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, err = w.Write(bytes)
	if err != nil {
		logger.Error("Error when writing data to client", err.Error())
	}
}

func WriteError(w http.ResponseWriter, errStr string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(map[string]string{
		"error": errStr,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Failed to marshal data", err.Error())
		return
	}
	w.WriteHeader(statusCode)
	_, err = w.Write(bytes)
	if err != nil {
		logger.Error("Error while writing data to client", err.Error())
	}
}
