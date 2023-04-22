/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Will explain given resources using ChatGPT",
	Long:  `given a resource in yaml format this command will return an explanation what is happening`,
	Run: func(cmd *cobra.Command, args []string) {
		var givenResource = ""
		if len(args) >= 1 && args[0] != "" {
			givenResource = args[0]
		}
		fmt.Printf("given resource %s", givenResource)
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
