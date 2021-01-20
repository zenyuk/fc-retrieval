package provider

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/spf13/viper"
)

// Provider configuration
type Provider struct {
	conf *viper.Viper
}

// Register data model
type Register struct {
	Address        string
	NetworkInfo    string
	RegionCode     string
	RootSigningKey string
	SigingKey      string
}

// CIDMessage data model
type CIDMessage struct {
	ProtocolVersion   string
	ProtocolSupported string
	MessageType       string
	Nonce             int
	ProviderID        string
	PricePerByte      int
	ExpiryDate        int64
	QosMetric         int
	Signature         string
	PieceCID          []string
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func postJSON(url string, data interface{}) error {
	json, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")

	r, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NewProvider creates a new provider
func NewProvider(conf *viper.Viper) *Provider {
	return &Provider{
		conf: conf,
	}
}

// Start a new provider
func Start(p *Provider) {
	scheme := p.conf.GetString("SERVICE_SCHEME")
	host := p.conf.GetString("SERVICE_HOST")
	port := p.conf.GetString("SERVICE_PORT")
	log.Info("Provider started at %s://%s:%s", scheme, host, port)
}

func generateDummyMessage() CIDMessage {
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	pieceCIDs := []string{"a", "b", "c", "d", "e"}
	dummyMessage := CIDMessage{
		ProtocolVersion:   "ProtocolVersion",
		ProtocolSupported: "ProtocolSupported",
		MessageType:       "MessageType",
		Nonce:             rand.Intn(100000),
		ProviderID:        "ProviderID",
		PricePerByte:      42,
		ExpiryDate:        expiryDate,
		QosMetric:         42,
		Signature:         "Signature",
		PieceCID:          pieceCIDs,
	}
	return dummyMessage
}

func (p *Provider) loop() {

	url := p.conf.GetString("REGISTER_API_URL") + "/registers/gateway"
	providerReg := Register{
		Address:        "f01213",
		NetworkInfo:    "127.0.0.1:80",
		RegionCode:     "FR",
		RootSigningKey: "0xABCDE123456789",
		SigingKey:      "0x987654321EDCBA",
	}

	postJSON(url, providerReg)
	for {
		gateways := []Register{}
		getJSON(url, &gateways)

		for _, gateway := range gateways {
			message := generateDummyMessage()
			fmt.Println("TODO, post to this gateway")
			fmt.Println(gateway)
			fmt.Println(message)
		}
		time.Sleep(25 * time.Second)
	}
}
