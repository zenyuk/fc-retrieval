package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

// SendMessage request Send JSON
func SendMessage(url string, message *fcrmessages.FCRMessage) (*fcrmessages.FCRMessage, error) {
	var data fcrmessages.FCRMessage
	jsonData, _ := json.Marshal(message)
	log.Info("Sending JSON to url: %v", url)
	contentReader := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", url, contentReader)
	req.Header.Set("Content-Type", "application/json")

	r, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &data, err
	}

	if r.StatusCode != 200 {
		err := errors.New("SendMessage receive error code: " + r.Status)
		return nil, err
	}

	json.Unmarshal(bytes, &data)
	return &data, nil
}
