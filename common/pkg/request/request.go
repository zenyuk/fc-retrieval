package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetJSON request Get JSON
func GetJSON(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// SendJSON request Send JSON
func SendJSON(url string, data interface{}) error {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	r, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

// SendJSON request Send JSON
func SendMessage(url string, message *fcrmessages.FCRMessage) (*fcrmessages.FCRMessage, error) {
	jsonData, _ := json.Marshal(message)
	log.Info("Sending JSON: %v to url: %v", string(jsonData), url)
	contentReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", url, contentReader)
	req.Header.Set("Content-Type", "application/json")

	r, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var data fcrmessages.FCRMessage
	json.NewDecoder(r.Body).Decode(&data)
	return &data, nil
}
