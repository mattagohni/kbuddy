/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"github.com/spf13/cobra"
	"os"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Will explain given topic related to kubernetes using ChatGPT",
	Long:  `given keyword e.g. (Deployment) will be explained using ChatGPT`,
	Run: func(cmd *cobra.Command, args []string) {
		var messages []goopenai.Message
		apiKey := os.Getenv("OPEN_AI_API_KEY")
		organization := os.Getenv("OPEN_AI_API_ORG")

		client := goopenai.NewClient(apiKey, organization)

		var givenSearchTerm = ""
		if len(args) >= 1 && args[0] != "" {
			givenSearchTerm = args[0]
		}
		fmt.Printf("given search term for explanation %s\n", givenSearchTerm)
		fmt.Printf("using API-Key %s\n", apiKey)
		fmt.Printf("for Org %s\n", organization)

		// new code
		// Load the file into a buffer

		// Create a runtime.Decoder from the Codecs field within
		// k8s.io/client-go that's pre-loaded with the schemas for all
		// the standard Kubernetes resource types.

		message := goopenai.Message{
			Role:    "user",
			Content: fmt.Sprintf("explain %s in context of a kubernetes resource", givenSearchTerm),
		}
		messages = append(messages, message)
		r := goopenai.CreateCompletionsRequest{
			Model:       "gpt-3.5-turbo",
			Messages:    messages,
			Temperature: 0.2,
		}

		completions, err := client.CreateCompletions(context.Background(), r)
		if err != nil {
			panic(err)
		}

		fmt.Println(completions)
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
		panic(e)
	}
}
