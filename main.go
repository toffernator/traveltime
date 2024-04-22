package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "traveltime",
	Short: "Tool for calculating travel times from one address to other addresses.",
	Long: `Tool for calclating travel times from one address to other addresses. Uses the
google maps routes API to provide transit time calculations`,
}

var arriveBy string
var calculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "Calculate the traveltime between one address and many other addresses.",
	Long: `Calculate the travel time from the first argument given to each of the remainder
arguments given. If no --departAt is given, then it defaults to the next
upcoming Tuesday at 10:00 UTC.`,
	Example: "calcluate london paris brussels --arriveBy 2024-04-25T09:00:00UTC+01:00",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.TODO()

		var latestArrivalTime time.Time
		var err error
		if arriveBy != "" {
			iso8601WithUtcOffsetFormat := "2006-01-02T15:04:05-07:00"
			latestArrivalTime, err = time.Parse(iso8601WithUtcOffsetFormat, arriveBy)
			if err != nil {
				log.Fatalf("Failed to parse --arriveBy: %v", err)
				return
			}
		}

		origin := Address(args[0])
		destinatations := make([]Address, len(args)-1)
		for i := 0; i < len(args)-1; i++ {
			destinatations[i] = Address(args[i+1])
		}

		var wg sync.WaitGroup
		for i := 0; i < len(destinatations); i++ {
			wg.Add(1)
			d := destinatations[i]
			go func() {
				defer wg.Done()

				var result ComputeTravelTimeResult
				if arriveBy != "" {
					// Then latestArrivalTime contains the parsed time.Time
					result, err = ComputeTravelTime(ctx, origin, d, WithArrivalTime(latestArrivalTime))
				} else {
					result, err = ComputeTravelTime(ctx, origin, d)
				}

				if err != nil {
					log.Printf("Failed to compute travel time from %s to %s: %v", origin, d, err)
				} else {
					fmt.Println(PresentComputeTravelTimeResult(result))
				}
			}()
		}
		wg.Wait()
		fmt.Println("Powered by Google, ©2024 Google")
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
	calculateCmd.Flags().StringVar(&arriveBy, "arriveBy", "", "Arrival time in ISO 8601 format (optional)")
	rootCmd.AddCommand(calculateCmd)
}

func main() {
	rootCmd.Execute()
}
