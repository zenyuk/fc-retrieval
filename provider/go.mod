module github.com/ConsenSys/fc-retrieval/provider

go 1.16

replace github.com/ConsenSys/fc-retrieval/common => ../common

require (
	github.com/ConsenSys/fc-retrieval/common v0.0.0-00010101000000-000000000000
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
)
