package management

import (
	"fmt"
	"log/slog"
	"net/http"
)

type HandlerBuilder struct {
	logger      *slog.Logger
	port        string
	swaggerUI   http.Handler
	swaggerPath string
}

func NewHandlerBuilder(logger *slog.Logger, port string) *HandlerBuilder {
	h := HandlerBuilder{
		port:   port,
		logger: logger,
	}
	return &h
}

func (h *HandlerBuilder) Build() {
	go func() {
		handler := http.NewServeMux()

		if h.swaggerUI != nil {
			handler.Handle(h.swaggerPath, h.swaggerUI)
			h.logger.Debug("Swagger UI available", "port", h.port, "path", h.swaggerPath)
		}

		addr := fmt.Sprintf(":%s", h.port)

		if err := http.ListenAndServe(addr, handler); err != nil {
			panic(err)
		}
	}()
}
