package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

var apiURL = "http://localhost:8080"
var output io.Writer = os.Stdout

func handleList(listType, action, cidr string) {
	if cidr == "" {
		fmt.Fprintln(output, "CIDR is required")
		return
	}
	url := fmt.Sprintf("%s/%s/%s?cidr=%s", apiURL, listType, action, cidr)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		fmt.Fprintln(output, "Error creating request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(output, "Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Fprintln(output, "Status:", resp.Status)
}

func handleBucket(action string) {
	if action != "reset" {
		fmt.Fprintln(output, "only reset supported")
		return
	}
	url := fmt.Sprintf("%s/bucket/reset", apiURL)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Fprintln(output, "Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Fprintln(output, "Status:", resp.Status)
}

func usage() {
	fmt.Fprintln(output, `
Usage:
  antibruteforce-cli whitelist add <cidr>
  antibruteforce-cli whitelist remove <cidr>
  antibruteforce-cli blacklist add <cidr>
  antibruteforce-cli blacklist remove <cidr>
  antibruteforce-cli bucket reset`)
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		return
	}

	switch args[0] {
	case "whitelist":
		if len(args) < 3 {
			usage()
			return
		}
		handleList("whitelist", args[1], args[2])
	case "blacklist":
		if len(args) < 3 {
			usage()
			return
		}
		handleList("blacklist", args[1], args[2])
	case "bucket":
		if len(args) < 2 {
			usage()
			return
		}
		handleBucket(args[1])
	default:
		usage()
	}
}
