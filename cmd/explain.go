/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"os"

	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Will explain given resources using ChatGPT",
	Long:  `given a resource in yaml format this command will return an explanation what is happening`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPEN_AI_API_KEY")
		organization := os.Getenv("OPEN_AI_API_ORG")

		client := goopenai.NewClient(apiKey, organization)

		var givenResource = ""
		if len(args) >= 1 && args[0] != "" {
			givenResource = args[0]
		}
		fmt.Printf("given resource %s\n", givenResource)
		fmt.Printf("using API-Key %s\n", apiKey)
		fmt.Printf("for Org %s\n", organization)

		dat, err := os.ReadFile(givenResource)
		check(err)

		readContent := string(dat)

		r := goopenai.CreateCompletionsRequest{
			Model: "gpt-3.5-turbo",
			Messages: []goopenai.Message{
				{
					Role:    "user",
					Content: fmt.Sprintf(`explain %s`, readContent),
				},
			},
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
