package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/hpcloud/tail"
)

const apiKey = "cl5d7071f3082be"

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage adif2logcloud <adiflog> <logcloudurl>")
	}
	file := os.Args[1]
	t, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Error tailing %s: %v", file, err)
	}
	for line := range t.Lines {
		sendLine(line.Text)
	}
}

type cloudlogRequest struct {
	APIKey string `json:"key"`
	Type   string `json:"type"`
	Line   string `json:"string"`
	Status string `json:"status,omitempty"`
}

func sendLine(line string) {
	log.Printf("Sending '%s", line)
	req := cloudlogRequest{APIKey: apiKey, Type: "adif", Line: line}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(req)
	resp, err := http.Post(os.Args[2], "application/json", buf)
	if err != nil {
		log.Fatalf("Error posting to %s: %v", os.Args[2], err)
	}
	defer resp.Body.Close()
	result := new(cloudlogRequest)
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	if result.Status != "created" {
		log.Printf("Failed to create record.\nResult status: %s", result.Status)
	}
}
