package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AppPort                  string `yaml:"appPort"`
	LLM                      string `yaml:"llm"`
	Scheme                   string `yaml:"scheme"`
	WvHost                   string `yaml:"wvHost"`
	WvPort                   string `yaml:"wvPort"`
	OllamaUrl                string `yaml:"ollamaUrl"`
	OllamaPort               string `yaml:"ollamaPort"`
	IndexName                string `yaml:"indexName"`
	DocumentsRetrievalNumber int    `yaml:"documentsRetrievalNumber"`
	MQEndpoint               string `yaml:"mqEndpoint"`
	AnswerExchange           string `yaml:"answerExchange"`
	FrontEndUrl              string `yaml:"frontEndUrl"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
