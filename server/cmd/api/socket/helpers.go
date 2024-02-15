package socket

import (
	"encoding/json"
	lg "gotrader/logger"
	"net/http"
)

func init() {
	logger = lg.CreateCustomLogger("server/socket")
}

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func readJson(r *http.Request, w http.ResponseWriter, v interface{}) error {
	maxBytes := 1024

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, int64(maxBytes)))

	err := decoder.Decode(v)

	if err != nil {
		logger.Printf("Error decoding json: %v", err)
		return err
	}

	return nil
}

func writeJson(w http.ResponseWriter, status int, v interface{}, headers ...http.Header) error {
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(v)

	if err != nil {
		logger.Printf("Error encoding json: %v", err)
		return err
	}

	for _, header := range headers {
		for key, value := range header {
			w.Header().Set(key, value[0])
		}
	}

	w.WriteHeader(status)
	w.Write(out)

	return nil
}

func writeError(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusInternalServerError

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse

	payload.Error = true
	payload.Message = err.Error()
	payload.Data = nil

	return writeJson(w, statusCode, payload)

}
