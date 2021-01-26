package handlers

import (
	"context"

	"encoding/json"

	middleware "github.com/go-openapi/runtime/middleware"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"

	"github.com/ConsenSys/fc-retrieval-register/models"
	op "github.com/ConsenSys/fc-retrieval-register/restapi/operations/registers"
)

// AddRegister to create a register
func AddRegister(params op.AddRegisterParams) middleware.Responder {
	redisHash := params.Type

	register := params.Register
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	err := rdb.HSet(ctx, redisHash, register, 0).Err()
	if err != nil {
		log.Error().Msg("Unable to set Redis value")
		panic(err)
	}

	log.Info().Str("type", redisHash).Msg("Register created")

	// Response
	return op.NewAddRegisterOK().WithPayload(register)
}

// GetRegisters retrieve register list
func GetRegisters(params op.GetRegistersParams) middleware.Responder {
	redisHash := params.Type

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	registers, err := rdb.HGetAll(ctx, redisHash).Result()
	if err != nil {
		log.Error().Msg("Unable to get Redis value")
		panic(err)
	}

	payload := []*models.Register{}
	for registerJson, _ := range registers {
		registerData := models.Register{}
		json.Unmarshal([]byte(registerJson), &registerData)
		payload = append(payload, &registerData)
	}

	return op.NewGetRegistersOK().WithPayload(payload)
}
