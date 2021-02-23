package handlers

import (
	"context"

	"encoding/json"

	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	middleware "github.com/go-openapi/runtime/middleware"
	redis "github.com/go-redis/redis/v8"

	"github.com/ConsenSys/fc-retrieval-register/models"
	op "github.com/ConsenSys/fc-retrieval-register/restapi/operations/gateway"
)

// AddGatewayRegister to create a gateway register
func AddGatewayRegister(params op.AddGatewayRegisterParams) middleware.Responder {
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
	return op.NewAddGatewayRegisterOK().WithPayload(register)
}

// GetGatewayRegisters retrieve Gateway register list
func GetGatewayRegisters(params op.GetGatewayRegistersParams) middleware.Responder {
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
	for registerJSON, _ := range registers {
		registerData := models.GatewayRegister{}
		json.Unmarshal([]byte(registerJSON), &registerData)
		payload = append(payload, &registerData)
	}

	return op.NewGetGatewayRegistersOK().WithPayload(payload)
}

// GetGatewayRegisterByID retrieve Gateway register by ID
func GetGatewayRegisterByID(params op.GetGatewayRegistersByIDParams) middleware.Responder {
	registerType := "gateway"
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

	payload := models.GatewayRegister{}
	for registerJSON, _ := range registers {
		registerData := models.GatewayRegister{}
		json.Unmarshal([]byte(registerJSON), &registerData)
		if (registerData.NodeID == registerID) {
			log.Info("Register found")
			payload = registerData
		} else {
			log.Info("Register not found")
		}
	}

	return op.NewGetGatewayRegistersByIDOK().WithPayload(&payload)
}

