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
	"log"
	"os"
)

const (
	responseFormat = `
{
	keyword: "The keyword to explain",
	explanation: "the up to 1000 word explanation of the given keyword",
	disclaimer: "disclaimer about correctness with with a link to the kubernetes docs",
	exampleYaml: "the yaml definition of the resource of the given keyword if the keyword is a kubernetes resource, empty otherwise",
	exampleJson: "the json definition of the resource of the given keyword if the keyword is a kubernetes resource, empty otherwise",
	furtherReading: [
	    {
			keyword: the keyword
			description: short description of the topic
			link: hyperlink to the kubernetes docs if possible
		}
	]
}
`
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Will explain given topic related to kubernetes using ChatGPT",
	Long:  `given keyword e.g. (Deployment) will be explained using ChatGPT`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Set(color.FgHiBlue)
		var messages []goopenai.Message
		apiKey := os.Getenv("OPEN_AI_API_KEY")
		organization := os.Getenv("OPEN_AI_API_ORG")

		client := goopenai.NewClient(apiKey, organization)

		var givenSearchTerm = ""
		if len(args) >= 1 && args[0] != "" {
			givenSearchTerm = args[0]
		}
		message := goopenai.Message{
			Role: "user",
			Content: fmt.Sprintf(
				"explain %s in context of kubernetes. in your response make a new line every 80"+
					" charecters. also structure your response in a json with the following format "+responseFormat, givenSearchTerm),
		}
		messages = append(messages, message)
		r := goopenai.CreateCompletionsRequest{
			Model:       "gpt-3.5-turbo",
			Messages:    messages,
			Temperature: 0.2,
		}

		completions, err := client.CreateCompletions(context.Background(), r)
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

		fmt.Println(keyword + "\n")
		fmt.Println(explanation + "\n")
		fmt.Println(disclaimer + "\n")
		fmt.Println(yamlExample)
		fmt.Println(jsonExample)
		for _, reading := range furtherReadings {
			fmt.Println(reading + "\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func check(e error) {
	if e != nil {
		log.Print(e)
		panic(e)
	}
}
