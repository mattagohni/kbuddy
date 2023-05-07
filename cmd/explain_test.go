package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"github.com/mattagohni/kbuddy/cmd/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var background = context.Background()

func TestExplainCmd(t *testing.T) {
	// arrange
	var mockClient = new(mocks.OpenAiClient)
	initTestServer(mockClient, createExpectedRequest("Deployment", "english"))

	buf, cmd := prepareExplainCommand(mockClient)

	// act
	cmd.Execute()

	expectedOutput := getExpectedExplainOutput()
	// assert
	assert.Equal(t, expectedOutput, buf.String())
}

func TestLanguageFlag(t *testing.T) {
	var mockClient = new(mocks.OpenAiClient)

	tests := map[string]struct {
		languageFlag string
	}{
		"it uses default language":    {languageFlag: ""},
		"it can use german language":  {languageFlag: "german"},
		"it can use spanish language": {languageFlag: "spanish"},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			expectedRequest := createExpectedRequest("Deployment", testCase.languageFlag)
			initTestServer(mockClient, expectedRequest)
			_, cmd := prepareExplainCommand(mockClient)
			cmd.Flags().Set("lang", testCase.languageFlag)

			// act
			cmd.Execute()

			// assert
			mockClient.AssertCalled(t, "CreateCompletions", background, expectedRequest)
		})
	}
}

func prepareExplainCommand(mockClient *mocks.OpenAiClient) (*bytes.Buffer, *cobra.Command) {
	buf := bytes.NewBufferString("")
	cmd := NewExplainCommand(mockClient.CreateCompletions)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"Deployment"})
	return buf, cmd
}

func getExpectedExplainOutput() string {
	return `Deployment

A Deployment in Kubernetes is a resource object that manages a set of replicas of a Pod. It provides declarative updates for Pods and ReplicaSets. A Deployment ensures that a specified number of replica Pods are running at any given time. If there are too few replicas, it will create more. If there are too many replicas, it will delete the excess. Deployments are useful when you need to update or roll back an application, as they provide a way to manage the deployment process without downtime. Deployments can also be used to scale an application up or down based on demand.

This information may not be accurate.

Pod
An overview of Pods.
https://kubernetes.io/docs/concepts/workloads/pods/

ReplicaSet
An overview of ReplicaSets.
https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/

`
}

func initTestServer(mockClient *mocks.OpenAiClient, expectedRequest goopenai.CreateCompletionsRequest) {
	// Mock response
	mockResponse := mockOpenAiResponse(mockClient, expectedRequest)

	// Create a new server to mock the API response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer testServer.Close()
	setEnvironmentVariables(testServer)
}

func setEnvironmentVariables(testServer *httptest.Server) {
	os.Setenv("OPEN_AI_API_KEY", "test-api-key")
	os.Setenv("OPEN_AI_API_ORG", "test-org")
	os.Setenv("OPEN_AI_API_URL", testServer.URL)
}

func createExpectedRequest(keyword string, language string) goopenai.CreateCompletionsRequest {
	if len(language) == 0 {
		language = "english"
	}
	var responseFormat, _ = os.ReadFile("./request/response_format.json")

	expectedRequest := goopenai.CreateCompletionsRequest{
		Model: "gpt-3.5-turbo",
		Messages: []goopenai.Message{
			{
				Role: "user",
				Content: fmt.Sprintf(
					"Explain %s in context of kubernetes. The user will need your response in the language: %s!"+
						"Add a disclaimer for your statement with reference to the docs."+
						"In your response make a new line every 80 charecters. Also structure your response in a json with the following format "+
						string(responseFormat),
					keyword, language),
			},
		},
		Temperature: 0.2,
	}
	return expectedRequest
}

func mockOpenAiResponse(mockClient *mocks.OpenAiClient, expectedRequest goopenai.CreateCompletionsRequest) []byte {
	mockResponse, err := os.ReadFile("../test/explain_response_deployment_en.json")
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

	mockClient.On("CreateCompletions", background, expectedRequest).Return(
		response,
		nil,
	)
	return mockResponse
}
