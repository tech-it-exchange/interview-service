package management

import httpSwagger "github.com/swaggo/http-swagger/v2"

type SwaggerConfig struct {
	SwaggerUrl  string
	SwaggerPath string
}

func (h *HandlerBuilder) AddSwaggerUIHandler(config SwaggerConfig) *HandlerBuilder {
	h.swaggerUI = httpSwagger.Handler(
		httpSwagger.URL(config.SwaggerUrl),
	)
	h.swaggerPath = config.SwaggerPath
	return h
}
