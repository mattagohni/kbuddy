package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"github.com/mattagohni/kbuddy/cmd/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestExplainCmd(t *testing.T) {
	// arrange
	var mockClient = new(mocks.OpenAiClient)

	var givenSearchTerm = "Deployment"
	var requestFormat, err = os.ReadFile("./request/explain_request.json")
	check(err)

	expectedRequest := createExpectedRequest(requestFormat, givenSearchTerm)

	initTestServer(err, mockClient, expectedRequest)

	buf := bytes.NewBufferString("")
	cmd := NewExplainCommand(mockClient.CreateCompletions)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"Deployment"})

	// act
	cmd.Execute()

	var expectedOutput = `Deployment

A Deployment provides declarative updates for Pods and ReplicaSets.

This information may not be accurate.



Pod
An overview of Pods.
https://kubernetes.io/docs/concepts/workloads/pods/

ReplicaSet
An overview of ReplicaSets.
https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/
`
	// assert
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output: %s\nGot: %s", expectedOutput, buf.String())
	}
}

func initTestServer(err error, mockClient *mocks.OpenAiClient, expectedRequest goopenai.CreateCompletionsRequest) {
	// Mock response
	mockResponse := mockOpenAiResponse(err, mockClient, expectedRequest)

	// Create a new server to mock the API response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer testServer.Close()
	os.Setenv("OPEN_AI_API_KEY", "test-api-key")
	os.Setenv("OPEN_AI_API_ORG", "test-org")
	os.Setenv("OPEN_AI_API_URL", testServer.URL)

}

func createExpectedRequest(responseFormat []byte, givenSearchTerm string) goopenai.CreateCompletionsRequest {
	expectedRequest := goopenai.CreateCompletionsRequest{
		Model: "gpt-3.5-turbo",
		Messages: []goopenai.Message{
			{
				Role: "user",
				Content: fmt.Sprintf(
					"explain %s in context of kubernetes. in your response make a new line every 80"+
						" charecters. also structure your response in a json with the following format "+string(responseFormat), givenSearchTerm),
			},
		},
		Temperature: 0.2,
	}
	return expectedRequest
}

func mockOpenAiResponse(err error, mockClient *mocks.OpenAiClient, expectedRequest goopenai.CreateCompletionsRequest) []byte {
	mockResponse, err := os.ReadFile("../test/explain_response.json")
	if err != nil {
		panic(err)
	}

	response := goopenai.CreateCompletionsResponse{
		Model: "gpt-3.5-turbo",
		Choices: []struct {
			Message struct {
				Role    string `json:"role,omitempty"`
				Content string `json:"content,omitempty"`
			} `json:"message"`
			Text         string      `json:"text,omitempty"`
			Index        int         `json:"index,omitempty"`
			Logprobs     interface{} `json:"logprobs,omitempty"`
			FinishReason string      `json:"finish_reason,omitempty"`
		}{
			{
				Message: goopenai.Message{
					Content: string(mockResponse),
				},
			},
		},
	}
	mockClient.On("CreateCompletions", context.Background(), expectedRequest).Return(
		response,
		nil,
	)
	return mockResponse
}
