/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"github.com/mattagohni/kbuddy/internal"
	v1 "k8s.io/api/apps/v1"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// analizeCmd represents the analize command
var analizeCmd = &cobra.Command{
	Use:   "analize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var messages []goopenai.Message
		apiKey := os.Getenv("OPEN_AI_API_KEY")
		organization := os.Getenv("OPEN_AI_API_ORG")

		client := goopenai.NewClient(apiKey, organization)

		var pathToFile, err = cmd.Flags().GetString("file")
		check(err)

		/*if len(args) >= 1 && args[0] != "" {
			pathToFile = args[0]
		}*/
		fmt.Printf("given file %s\n", pathToFile)
		fmt.Printf("using API-Key %s\n", apiKey)
		fmt.Printf("for Org %s\n", organization)

		data, err := os.ReadFile(pathToFile)
		if err != nil {
			log.Fatal(err)
		}

		for _, resourceYAML := range strings.Split(string(data), "---") {

			// skip empty documents, `Decode` will fail on them
			if len(resourceYAML) == 0 {
				continue
			}

			// - obj is the API object (e.g., Deployment)
			// - groupVersionKind is a generic object that allows
			//   detecting the API type we are dealing with, for
			//   accurate type casting later.
			obj, groupVersionKind, err := internal.Parse(resourceYAML)
			check(err)

			// Figure out from `Kind` the resource type, and attempt
			// to cast appropriately.
			if groupVersionKind.Group == "apps" &&
				groupVersionKind.Version == "v1" &&
				groupVersionKind.Kind == "Deployment" {
				deployment := obj.(*v1.Deployment)
				message := goopenai.Message{
					Role:    "user",
					Content: fmt.Sprintf("analize violated best practices in the following Spec of a kubernetes resource and provide tips: \n%s", deployment.Spec.String()),
				}
				messages = append(messages, message)
				log.Print(deployment.ObjectMeta.Name)
			}
		}
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
	rootCmd.AddCommand(analizeCmd)

	analizeCmd.Flags().String("file", "", "path to a file for explanation")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// analizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// analizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
