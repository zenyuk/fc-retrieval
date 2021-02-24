package handlers

import (
	"strings"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/ConsenSys/fc-retrieval-register/models"
	op "github.com/ConsenSys/fc-retrieval-register/restapi/operations/homepage"
)

// Homepage handler
func HomepageHandler() middleware.Responder {
	serviceName := apiconfig.GetString("SERVICE_NAME")

	// Response
	payload := models.Ack{Status: "success", Message: strings.Join([]string{"Welcome to", serviceName, "api, please go to /docs."}, " ")}
	return op.NewHomepageOK().WithPayload(&payload)
}
