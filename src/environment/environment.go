package environment

import (
	"github.com/caarlos0/env/v11"
	"log"
)

type Envs struct {
	AppPort                    string `env:"BACKEND_SERVICE_PORT" envDefault:"8080"`
	AllowedOrigins             string `env:"ALLOWED_ORIGINS"`
	OllamaHost                 string `env:"OLLAMA_SERVICE_HOST"`
	OllamaPort                 string `env:"OLLAMA_SERVICE_PORT" envDefault:"11434"`
	WvHost                     string `env:"WEAVIATE_SERVICE_HOST"`
	WvPort                     string `env:"WEAVIATE_SERVICE_PORT" envDefault:"8080"`
	WvScheme                   string `env:"WEAVIATE_SCHEME" envDefault:"http"`
	WvIndexName                string `env:"WEAVIATE_INDEX_NAME" envDefault:"Document"`
	WvDocumentsRetrievalNumber int    `env:"WV_DOCUMENTS_RETRIEVAL_NUMBER" envDefault:"3"`
	RabbitMQProtocol           string `env:"RABBITMQ_PROTOCOL" envDefault:"amqp"`
	RabbitMQUsername           string `env:"RABBITMQ_USERNAME"`
	RabbitMQPassword           string `env:"RABBITMQ_PASSWORD"`
	RabbitMQHost               string `env:"RABBITMQ_SERVICE_HOST"`
	RabbitMQPort               string `env:"RABBITMQ_SERVICE_PORT" envDefault:"5672"`
	RabbitMQVHost              string `env:"RABBITMQ_VHOST" envDefault:"/%2F"`
	LLM                        string `env:"LLM"`
	RabbitMQAnswerExchange     string `env:"RABBITMQ_ANSWER_EXCHANGE" envDefault:"answers-topic"`
}

func NewEnvs() *Envs {
	var envs Envs
	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Error reading environment variables: %v", err)
	}

	return &envs
}
