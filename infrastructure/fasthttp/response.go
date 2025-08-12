package fasthttp

import (
	"encoding/json"
	"github.com/google/uuid"
	fastHttpLib "github.com/valyala/fasthttp"
	"log"
	"time"
)

func WriteResponse(ctx *fastHttpLib.RequestCtx, t interface{}) error {
	response, err := json.Marshal(t)
	if err != nil {
		return err
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fastHttpLib.StatusOK)
	ctx.SetBody(response)

	return nil
}

func WriteErrorResponse(ctx *fastHttpLib.RequestCtx, message string, statusCode int, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}

	errorResponse := map[string]interface{}{
		"message":      message,
		"path":         string(ctx.Path()),
		"stackTraceId": uuid.New().String(),
		"statusCode":   statusCode,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")

	if err := json.NewEncoder(ctx).Encode(errorResponse); err != nil {
		ctx.Error("Internal Server Error", fastHttpLib.StatusInternalServerError)
	}
}
