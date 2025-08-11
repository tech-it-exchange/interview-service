package management

import httpSwagger "github.com/swaggo/http-swagger/v2"

type SwaggerConfig struct {
	SwaggerUrl  string
	SwaggerPath string
}

func (h *HandlerBuilder) AddSwaggerUIHandler(config SwaggerConfig) *HandlerBuilder {
	// Используем встроенный doc.json по умолчанию, чтобы UI и JSON отдавались по одному и тому же пути
	h.swaggerUI = httpSwagger.Handler()
	h.swaggerPath = config.SwaggerPath
	return h
}
