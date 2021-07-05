package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-redis/redis/v8"

	log "github.com/ConsenSys/fc-retrieval/common/pkg/logging"

	"github.com/ConsenSys/fc-retrieval/register/models"
	op "github.com/ConsenSys/fc-retrieval/register/restapi/operations/provider"
)

// AddProviderRegister to create a provider register
func AddProviderRegister(params op.AddProviderRegisterParams) middleware.Responder {
	register := params.Register
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	err := rdb.HSet(ctx, "provider", register.NodeID, register).Err()
	if err != nil {
		msg := "Unable to set Redis value"
		log.Error("%s, error in AddProviderRegister: %s", msg, err.Error())
		return op.NewAddProviderRegisterDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	log.Info("register created a provider record with ID: %s", params.Register.NodeID)

	// Response
	return op.NewAddProviderRegisterOK().WithPayload(register)
}

// GetProviderRegisters retrieve Provider register list
func GetProviderRegisters(_ op.GetProviderRegistersParams) middleware.Responder {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	providerRegisters, err := rdb.HGetAll(ctx, "provider").Result()
	if err != nil {
		msg := "Unable to get Redis value"
		log.Error("%s, error in GetProviderRegisters: %s", msg, err.Error())
		return op.NewGetProviderRegistersDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	var payload []*models.ProviderRegister
	var debugOutputSb strings.Builder
	for _, register := range providerRegisters {
		registerData := models.ProviderRegister{}
		if unmarshallErr := json.Unmarshal([]byte(register), &registerData); unmarshallErr != nil {
			log.Error("inside GetProviderRegisters - can't unmarshall JSON: %s", unmarshallErr.Error())
		}
		payload = append(payload, &registerData)
		debugOutputSb.WriteString(fmt.Sprintf("%s, ", registerData.NodeID))
	}

	return op.NewGetProviderRegistersOK().WithPayload(payload)
}

// GetProviderRegisterByID retrieve Provider register by ID
func GetProviderRegisterByID(params op.GetProviderRegistersByIDParams) middleware.Responder {
	registerID := params.ID
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	register, err := rdb.HGet(ctx, "provider", registerID).Result()
	if err != nil {
		msg := "Register not found"
		log.Error(msg)
		return op.NewGetProviderRegistersByIDDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	registerData := models.ProviderRegister{}
	if unmarshallErr := json.Unmarshal([]byte(register), &registerData); unmarshallErr != nil {
		msg := "can't unmarshall JSON"
		log.Error("inside GetProviderRegisterByID - %s: %s", msg, unmarshallErr.Error())
		return op.NewGetProviderRegistersByIDDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	payload := registerData
	return op.NewGetProviderRegistersByIDOK().WithPayload(&payload)
}

// DeleteProviderRegisters deletes all Providers
func DeleteProviderRegisters(_ op.DeleteProviderRegisterParams) middleware.Responder {
	const registerTypeProvider = "provider"

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	ctx := context.Background()

	registers, err := rdb.HGetAll(ctx, registerTypeProvider).Result()
	if err != nil {
		msg := "Unable to get Redis value"
		log.Error("%s, error: %s", msg, err.Error())
		return op.NewDeleteProviderRegisterDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	for index := range registers {
		log.Info("register deleted a provider record with ID: %s", index)
		err := rdb.HDel(ctx, registerTypeProvider, index).Err()
		if err != nil {
			log.Error("Unable to set Redis value")
			panic(err)
		}
	}

	payload := models.Ack{Status: "success", Message: "All Providers have been deleted"}
	return op.NewDeleteProviderRegisterOK().WithPayload(&payload)
}
