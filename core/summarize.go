package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

// SummarizeReadme sends the README content to Vertex AI Gemini API for summarization.
func SummarizeReadme(w io.Writer, projectID, location, modelName, content string) error {
	ctx := context.Background()

	// Initialize the Vertex AI client
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("error creating Vertex AI client: %w", err)
	}
	defer client.Close()

	// Define the generative model
	gemini := client.GenerativeModel(modelName)

	// Send the content as a prompt
	prompt := genai.Text(fmt.Sprintf("Summarize the following text:\n\n%s", content))

	// Call the GenerateContent API
	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}

	// Format the response as JSON for readability
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting response as JSON: %w", err)
	}

	// Write the response to the provided writer
	fmt.Fprintln(w, string(rb))
	return nil
}