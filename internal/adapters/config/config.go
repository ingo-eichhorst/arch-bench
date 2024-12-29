package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	EvalAPIKey   string
	EvalModel    string
	EvalProvider string
}

func LoadConfig() (*Config, error) {
	pwd, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(pwd, "../../.env"))
	if err != nil {
		return nil, err
	}

	return &Config{
		EvalAPIKey:   os.Getenv("EVAL_API_KEY"),
		EvalModel:    os.Getenv("EVAL_MODEL"),
		EvalProvider: os.Getenv("EVAL_PROVIDER"),
	}, nil
}
