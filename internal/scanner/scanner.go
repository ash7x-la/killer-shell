package scanner

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

func ScanTarget() {
	var target string
	fmt.Printf("\n%s[ TARGET SCANNER ]%s\n", yellow, reset)
	fmt.Print("Target URL: ")
	fmt.Scanln(&target)
	if !strings.HasPrefix(target, "http") {
		target = "http://" + target
	}

	client := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Get(target)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	stack, rec := "Unknown", "php"
	for k, v := range resp.Header {
		val := strings.ToLower(strings.Join(v, ", "))
		if k == "X-Powered-By" {
			if strings.Contains(val, "php") {
				stack = "PHP"
				rec = "php"
			}
			if strings.Contains(val, "express") || strings.Contains(val, "node") {
				stack = "Node.js"
				rec = "node"
			}
		}
		if strings.Contains(k, "Server") && (strings.Contains(val, "gunicorn") || strings.Contains(val, "werkzeug")) {
			stack = "Python"
			rec = "python"
		}
	}
	fmt.Printf("%s[+]%s Stack: %s | Rec: %s\n", green, reset, stack, rec)
	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
