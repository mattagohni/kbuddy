/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/franciscoescher/goopenai"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Will explain given resources using ChatGPT",
	Long:  `given a resource in yaml format this command will return an explanation what is happening`,
	Run: func(cmd *cobra.Command, args []string) {
		var messages []goopenai.Message
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

		// new code
		// Load the file into a buffer

		data, err := os.ReadFile(givenResource)
		if err != nil {
			log.Fatal(err)
		}

		// Create a runtime.Decoder from the Codecs field within
		// k8s.io/client-go that's pre-loaded with the schemas for all
		// the standard Kubernetes resource types.
		decoder := scheme.Codecs.UniversalDeserializer()

		for _, resourceYAML := range strings.Split(string(data), "---") {

			// skip empty documents, `Decode` will fail on them
			if len(resourceYAML) == 0 {
				continue
			}

			// - obj is the API object (e.g., Deployment)
			// - groupVersionKind is a generic object that allows
			//   detecting the API type we are dealing with, for
			//   accurate type casting later.
			obj, groupVersionKind, err := decoder.Decode(
				[]byte(resourceYAML),
				nil,
				nil)
			if err != nil {
				log.Print(err)
				continue
			}

			// Figure out from `Kind` the resource type, and attempt
			// to cast appropriately.
			if groupVersionKind.Group == "apps" &&
				groupVersionKind.Version == "v1" &&
				groupVersionKind.Kind == "Deployment" {
				deployment := obj.(*v1.Deployment)
				message := goopenai.Message{
					Role:    "user",
					Content: fmt.Sprintf("explain %s", deployment.Spec.String()),
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
