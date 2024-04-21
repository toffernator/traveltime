package main

import (
	"os"

	"example.com/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "traveltime",
	Short: "Tool for calculating travel times from one address to other addresses.",
	Long:  `TODO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var calculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "Calculate the traveltime between one address and many other addresses.",
	Long:  "Calculate the traveltime between one address and many other addresses.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.example.com.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(calculateCmd)
}

func main() {
	cmd.Execute()
}
