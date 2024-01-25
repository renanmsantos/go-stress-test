package cmd

import (
	"os"

	"github.com/renanmsantos/go-stress-test/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-stress-test",
	Short: "A go stress test tool",
	Long:  `This application is a CLI tool to stress test a given url.`,
	Run: func(cmd *cobra.Command, args []string) {
		urlFlag, _ := cmd.Flags().GetString("url")
		requestsFlag, _ := cmd.Flags().GetInt("requests")
		concurrencyFlag, _ := cmd.Flags().GetInt("concurrency")
		stressTest := internal.NewStressTest(urlFlag, requestsFlag, concurrencyFlag)
		stressTest.ExecuteTests()
		stressTest.PrintMetrics()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "Flag to set the url to be tested")
	rootCmd.Flags().IntP("requests", "r", 10, "Flag to set the number of requests to be made")
	rootCmd.Flags().IntP("concurrency", "c", 1, "Flag to set the number of concurrent requests to be made")
	rootCmd.MarkFlagRequired("url")
}
