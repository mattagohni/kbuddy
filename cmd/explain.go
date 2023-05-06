/*
Copyright © 2023 Matthias Alt <mattagohni@gmail.com>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	goopenai "github.com/franciscoescher/goopenai"
	. "github.com/mattagohni/kbuddy/internal/response"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{}

func init() {
	apiKey := os.Getenv("OPEN_AI_API_KEY")
	organization := os.Getenv("OPEN_AI_API_ORG")

	client := NewMyOpenAiClient(apiKey, organization)
	callOpenAi := client.CreateCompletions

	explainCmd := NewExplainCommand(callOpenAi)
	rootCmd.AddCommand(explainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func NewExplainCommand(sendExplainRequest func(ctx context.Context, req goopenai.CreateCompletionsRequest) (goopenai.CreateCompletionsResponse, error)) *cobra.Command {
	var explainCmd = &cobra.Command{
		Use:   "explain",
		Short: "Will explain given topic related to kubernetes using ChatGPT",
		Long:  `given keyword e.g. (Deployment) will be explained using ChatGPT`,
		Run: func(cmd *cobra.Command, args []string) {
			color.Set(color.FgHiBlue)
			var messages []goopenai.Message

			var givenSearchTerm = ""
			if len(args) >= 1 && args[0] != "" {
				givenSearchTerm = args[0]
			}

			_, filename, _, _ := runtime.Caller(0)
			pathToCurrentDir := filepath.Dir(filename)

			var responseFormat, err = os.ReadFile(pathToCurrentDir + "/request/explain_request.json")
			check(err)

			message := goopenai.Message{
				Role: "user",
				Content: fmt.Sprintf(
					"explain %s in context of kubernetes. in your response make a new line every 80"+
						" charecters. also structure your response in a json with the following format "+string(responseFormat), givenSearchTerm),
			}
			messages = append(messages, message)
			r := goopenai.CreateCompletionsRequest{
				Model:       "gpt-3.5-turbo",
				Messages:    messages,
				Temperature: 0.2,
			}

			completions, err := sendExplainRequest(context.Background(), r)
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

			output := []string{keyword, explanation, disclaimer, yamlExample, jsonExample}

			fmt.Fprintln(cmd.OutOrStdout(), keyword+"\n")
			fmt.Fprintln(cmd.OutOrStdout(), explanation+"\n")
			fmt.Fprintln(cmd.OutOrStdout(), disclaimer+"\n")
			fmt.Fprintln(cmd.OutOrStdout(), yamlExample)
			fmt.Fprintln(cmd.OutOrStdout(), jsonExample)
			for _, reading := range furtherReadings {
				output = append(output, reading)
				fmt.Fprintln(cmd.OutOrStdout(), reading+"\n")
			}
		},
	}

	// add any flags or arguments that the command needs
	// ...

	return explainCmd
}

func check(e error) {
	if e != nil {
		log.Print(e)
		panic(e)
	}
}
