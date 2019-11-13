package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hpcloud/tail"
)

var apiKey string

func main() {
	config := args2config()
	var found bool
	apiKey, found = os.LookupEnv("CLOUDLOG_API_KEY")
	if !found {
		fmt.Println("CLOUDLOG_API_KEY must be set\n")
		flag.Usage()
		os.Exit(1)
	}
	file := flag.Args()[0]
	t, err := tail.TailFile(file, config)
	if err != nil {
		log.Fatalf("Error tailing %s: %v", file, err)
	}
	url := flag.Args()[1]
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "index.php/api/qso"
	for line := range t.Lines {
		sendLine(line.Text, url)
	}
}

func args2config() tail.Config {
	flag.Usage = func() {
		fmt.Printf(`Usage: %s [-b] <ADIF log> <cloudlog url>
  Example: adif2logcloud ~/.local/share/WSJT-X/wsjtx_log.adi http://localhost/Cloudlog

  The Cloudlog API key shoud be passed in the CLOUDLOG_API_KEY environment variable.
`, os.Args[0])
		flag.PrintDefaults()
	}
	config := tail.Config{Follow: true}
	var fromBeginning bool
	flag.BoolVar(&fromBeginning, "b", false, "If true, load entire log file from the beginning, otherwise tail the file, only posting new entries to Cloudlog.")
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if !fromBeginning {
		config.Location = &tail.SeekInfo{
			Whence: os.SEEK_END,
		}
	}
	return config
}

type cloudlogRequest struct {
	APIKey string `json:"key"`
	Type   string `json:"type"`
	Line   string `json:"string"`
	Status string `json:"status,omitempty"`
}

func sendLine(line string, url string) {
	if len(line) == 0 {
		log.Print("Skipping blank line")
		return
	}
	log.Printf("Sending '%s'", line)
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
