package handlers

import (
	"context"

	"encoding/json"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	middleware "github.com/go-openapi/runtime/middleware"
	redis "github.com/go-redis/redis/v8"

	"github.com/ConsenSys/fc-retrieval-register/models"
	opG "github.com/ConsenSys/fc-retrieval-register/restapi/operations/gateway"
	opP "github.com/ConsenSys/fc-retrieval-register/restapi/operations/provider"
)

// AddGatewayRegister to create a gateway register
func AddGatewayRegister(params opG.AddGatewayRegisterParams) middleware.Responder {
	registerType := "gateway"
	register := params.Register
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	err := rdb.HSet(ctx, registerType, register, 0).Err()
	if err != nil {
		log.Error("Unable to set Redis value")
		panic(err)
	}

	log.Info("Register created %v", registerType)

	// Response
	return opG.NewAddGatewayRegisterOK().WithPayload(register)
}

// AddProviderRegister to create a provider register
func AddProviderRegister(params opP.AddProviderRegisterParams) middleware.Responder {
	registerType := "provider"
	register := params.Register
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	err := rdb.HSet(ctx, registerType, register, 0).Err()
	if err != nil {
		log.Error("Unable to set Redis value")
		panic(err)
	}

	log.Info("Register created %v", registerType)

	// Response
	return opP.NewAddProviderRegisterOK().WithPayload(register)
}

// GetGatewayRegisters retrieve Gateway register list
func GetGatewayRegisters(params opG.GetGatewayRegistersParams) middleware.Responder {
	registerType := "gateway"
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	registers, err := rdb.HGetAll(ctx, registerType).Result()
	if err != nil {
		log.Error("Unable to get Redis value")
		panic(err)
	}

	payload := []*models.GatewayRegister{}
	for registerJson, _ := range registers {
		registerData := models.GatewayRegister{}
		json.Unmarshal([]byte(registerJson), &registerData)
		payload = append(payload, &registerData)
	}

	return opG.NewGetGatewayRegistersOK().WithPayload(payload)
}

// GetProviderRegisters retrieve Provider register list
func GetProviderRegisters(params opP.GetProviderRegistersParams) middleware.Responder {
	registerType := "provider"
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	registers, err := rdb.HGetAll(ctx, registerType).Result()
	if err != nil {
		log.Error("Unable to get Redis value")
		panic(err)
	}

	payload := []*models.ProviderRegister{}
	for registerJson, _ := range registers {
		registerData := models.ProviderRegister{}
		json.Unmarshal([]byte(registerJson), &registerData)
		payload = append(payload, &registerData)
	}

	return opP.NewGetProviderRegistersOK().WithPayload(payload)
}
