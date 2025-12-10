package mdl

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

// FromEnv defines an llm model from environment variable MODEL
func FromEnv() model.LLM {
	if os.Getenv("GOOGLE_API_KEY") == "" {
		log.Printf("model: %q, project: %q", os.Getenv("MODEL"), os.Getenv("GOOGLE_CLOUD_PROJECT"))
		model, err := gemini.NewModel(context.Background(),
			os.Getenv("MODEL"),
			&genai.ClientConfig{}, // values set from environment
		)
		if err != nil {
			log.Fatalf("could not create model on Vertex AI: %v", err)
		}

		return model
	}

	model, err := gemini.NewModel(context.Background(),
		os.Getenv("MODEL"),
		&genai.ClientConfig{APIKey: os.Getenv("GOOGLE_API_KEY")})
	if err != nil {
		log.Fatalf("could not create model: %v", err)
	}

	return model

}
