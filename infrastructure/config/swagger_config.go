package config

type SwaggerConfig struct {
	SwaggerPath string `env:"SWAGGER_PATH" envDefault:"/docs/"`
	SwaggerUrl  string `env:"SWAGGER_URL" envDefault:"http://localhost:9090/docs/doc.json"`
}
