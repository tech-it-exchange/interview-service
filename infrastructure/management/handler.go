package management

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerBuilder struct {
	logger      *log.Logger
	port        string
	swaggerUI   http.Handler
	swaggerPath string
}

func NewHandlerBuilder(logger *log.Logger, port string) *HandlerBuilder {
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
			h.logger.Printf("Swagger UI available on port %s path %s", h.port, h.swaggerPath)
		}

		addr := fmt.Sprintf(":%s", h.port)

		if err := http.ListenAndServe(addr, handler); err != nil {
			panic(err)
		}
	}()
}
