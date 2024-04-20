package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const ComputeRoutesUrl = "https://routes.googleapis.com/directions/v2:computeRoutes"

type Place struct {
	Address string `json:"address"`
}

type OriginDestination struct {
	Origin      Place `json:"origin"`
	Destination Place `json:"destination"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")

	originDestination := OriginDestination{
		Origin:      Place{Address: "1800 Amphitheatre Parkway, Mountain View, CA 94043"},
		Destination: Place{Address: "Sloat Blvd &, Upper Great Hwy, San Francisco, CA 94132"},
	}
	reqBody, err := json.Marshal(&originDestination)
	if err != nil {
		log.Fatalf("Failed to marshall body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ComputeRoutesUrl, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)
	req.Header.Set("X-Goog-FieldMask", "routes.duration,routes.distanceMeters")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println(string(respBody))
}

func sourceDotEnv() (string, error) {
	f, err := os.Open(".env")
	if err != nil {
		return "", err
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
