/*
Copyright © 2023 Matthias Alt <mattagohni@gmail.com>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/franciscoescher/goopenai"
	. "github.com/mattagohni/kbuddy/internal/response"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const DefaultLanguage = "english"

func init() {
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	organization := os.Getenv("OPEN_AI_API_ORG")

	client := NewMyOpenAiClient(apiKey, organization)
	callOpenAi := client.CreateCompletions

	explainCmd := NewExplainCommand(callOpenAi)
	rootCmd.AddCommand(explainCmd)
}

func NewExplainCommand(sendExplainRequest func(ctx context.Context, req goopenai.CreateCompletionsRequest) (goopenai.CreateCompletionsResponse, error)) *cobra.Command {
	var language string
	var explainCmd = &cobra.Command{
		Use:   "explain <keyword>",
		Short: "Will explain given topic related to kubernetes using ChatGPT",
		Long:  `given keyword e.g. (Deployment) will be explained using ChatGPT`,
		Run: func(cmd *cobra.Command, args []string) {
			color.Set(color.FgHiBlue)
			var messages []goopenai.Message

			var givenSearchTerm = ""
			if len(args) >= 1 && args[0] != "" {
				givenSearchTerm = args[0]
			}

			if len(language) == 0 {
				language = DefaultLanguage
			}

			_, filename, _, _ := runtime.Caller(0)
			pathToCurrentDir := filepath.Dir(filename)

			var responseFormat, err = os.ReadFile(pathToCurrentDir + "/request/response_format.json")
			check(err)

			message := goopenai.Message{
				Role: "user",
				Content: fmt.Sprintf(
					"Explain %s in context of kubernetes. The user will need your response in the language: %s!"+
						"Add a disclaimer for your statement with reference to the docs."+
						"In your response make a new line every 80 charecters. Also structure your response in a json with the following format "+
						string(responseFormat),
					givenSearchTerm, language),
			}
			messages = append(messages, message)
			createCompletionsRequest := goopenai.CreateCompletionsRequest{
				Model:       "gpt-3.5-turbo",
				Messages:    messages,
				Temperature: 0.2,
			}

			completions, err := sendExplainRequest(context.Background(), createCompletionsRequest)
			check(err)

			var explainResponse ExplainResponse
			if err := json.Unmarshal([]byte(completions.Choices[0].Message.Content), &explainResponse); err != nil {
				panic(err)
			}

			keyword := color.HiMagentaString(explainResponse.Keyword)
			explanation := color.HiWhiteString(explainResponse.Explanation)
			var furtherReadings []string
			disclaimer := color.HiYellowString(explainResponse.Disclaimer)

			for _, reading := range explainResponse.FurtherReading {
				furtherReadings = append(furtherReadings, color.WhiteString(reading.Keyword+"\n"+reading.Description+"\n"+reading.Link))
			}

			yamlExample := explainResponse.ExampleYaml
			jsonExample := explainResponse.ExampleJson

			output := createOutput(keyword, explanation, disclaimer, yamlExample, jsonExample, furtherReadings)
			printOutput(cmd.OutOrStdout(), output)
		},
	}

	explainCmd.Flags().StringVarP(&language, "lang", "l", "english", "kbuddy explain [keyword] --lang=german")
	return explainCmd
}

func printOutput(writer io.Writer, output []string) {
	for _, outputPart := range output {
		if len(outputPart) > 0 {
			_, err := fmt.Fprintln(writer, outputPart+"\n")
			check(err)
		}
	}
}

func createOutput(keyword string, explanation string, disclaimer string, yamlExample string, jsonExample string, furtherReadings []string) []string {
	output := []string{keyword, explanation, disclaimer, yamlExample, jsonExample}

	for _, reading := range furtherReadings {
		output = append(output, reading)
	}
	return output
}

func check(e error) {
	if e != nil {
		log.Print(e)
		panic(e)
	}
}
