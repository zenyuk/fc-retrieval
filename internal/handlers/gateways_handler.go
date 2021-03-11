package handlers

import (
	"context"

	"encoding/json"

	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
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

	err := rdb.HSet(ctx, registerType, register.NodeID, register).Err()
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
	for _, register := range registers {
		registerData := models.GatewayRegister{}
		json.Unmarshal([]byte(register), &registerData)
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

	register, err := rdb.HGet(ctx, registerType, registerID).Result()
	if err != nil {
		msg := "Register not found"
		log.Error(msg)
		return op.NewGetGatewayRegistersByIDDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	registerData := models.GatewayRegister{}
	json.Unmarshal([]byte(register), &registerData)

	payload := registerData
	return op.NewGetGatewayRegistersByIDOK().WithPayload(&payload)
}

// DeleteGatewayRegisters deletes all Gateways
func DeleteGatewayRegisters(params op.DeleteGatewayRegisterParams) middleware.Responder {
	registerType := "gateway"

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	ctx := context.Background()

	registers, err := rdb.HGetAll(ctx, registerType).Result()
	if err != nil {
		log.Error("Unable to get Redis value")
		panic(err)
	}

	for index, _ := range registers {
		log.Info("DELETE %v", index)
		err := rdb.HDel(ctx, registerType, index).Err()
		if err != nil {
			log.Error("Unable to set Redis value")
			panic(err)
		}
	}

	payload := models.Ack{Status: "success", Message: "All Gateways have been deleted"}
	return op.NewDeleteGatewayRegisterOK().WithPayload(&payload)
}