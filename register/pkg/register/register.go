package register

// GatewayRegister data model
type GatewayRegister struct {
	NodeID              string
	Address             string
	NetworkGatewayInfo  string
	NetworkProviderInfo string
	NetworkProviderInfo string
	NetworkClientInfo   string
	RegionCode          string
	RootSigningKey      string
	SigingKey           string
}

// ProviderRegister data model
type ProviderRegister struct {
	NodeID         string
	Address        string
	NetworkInfo    string
	RegionCode     string
	RootSigningKey string
	SigingKey      string
}
