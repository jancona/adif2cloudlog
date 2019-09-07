package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hpcloud/tail"
)

var apiKey string

func main() {
	var found bool
	apiKey, found = os.LookupEnv("CLOUDLOG_API_KEY")
	if !found {
		log.Fatal("CLOUDLOG_API_KEY must be set")
	}
	if len(os.Args) != 3 {
		log.Fatal("Usage adif2logcloud <ADIF log> <cloudlog url>\n  Example adif2logcloud ~/.local/share/WSJT-X/wsjtx_log.adi http://localhost/cloudlog")
	}
	file := os.Args[1]
	t, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Error tailing %s: %v", file, err)
	}
	url := os.Args[2]
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "index.php/api/qso"
	for line := range t.Lines {
		sendLine(line.Text, url)
	}
}

type cloudlogRequest struct {
	APIKey string `json:"key"`
	Type   string `json:"type"`
	Line   string `json:"string"`
	Status string `json:"status,omitempty"`
}

func sendLine(line string, url string) {
	log.Printf("Sending '%s", line)
	req := cloudlogRequest{APIKey: apiKey, Type: "adif", Line: line}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(req)
	resp, err := http.Post(url, "application/json", buf)
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
