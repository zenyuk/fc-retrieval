package handlers

import (
	"context"

	"encoding/json"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	middleware "github.com/go-openapi/runtime/middleware"
	redis "github.com/go-redis/redis/v8"

	"github.com/ConsenSys/fc-retrieval-register/models"
	op "github.com/ConsenSys/fc-retrieval-register/restapi/operations/provider"
)

// AddProviderRegister to create a provider register
func AddProviderRegister(params op.AddProviderRegisterParams) middleware.Responder {
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
	return op.NewAddProviderRegisterOK().WithPayload(register)
}

// GetProviderRegisters retrieve Provider register list
func GetProviderRegisters(params op.GetProviderRegistersParams) middleware.Responder {
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
	for registerJSON, _ := range registers {
		registerData := models.ProviderRegister{}
		json.Unmarshal([]byte(registerJSON), &registerData)
		payload = append(payload, &registerData)
	}

	return op.NewGetProviderRegistersOK().WithPayload(payload)
}

// GetProviderRegisterByID retrieve Provider register by ID
func GetProviderRegisterByID(params op.GetProviderRegistersByIDParams) middleware.Responder {
	registerType := "provider"
	registerID := params.ID
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

	payload := models.ProviderRegister{}
	for registerJSON, _ := range registers {
		registerData := models.ProviderRegister{}
		json.Unmarshal([]byte(registerJSON), &registerData)
		if (registerData.NodeID == registerID) {
			log.Info("Register found")
			payload = registerData
		} else {
			log.Info("Register not found")
		}
	}

	return op.NewGetProviderRegistersByIDOK().WithPayload(&payload)
}
