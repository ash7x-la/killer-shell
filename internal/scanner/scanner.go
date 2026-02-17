package scanner

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	Target  string
	Stack   string
	Payload string
	Headers map[string]string
}

// Scan performs fingerprinting on the target URL
func Scan(target string) (*Result, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &Result{
		Target:  target,
		Headers: make(map[string]string),
	}

	// Analyze Headers
	stack := "unknown"
	payload := "php" // Default

	for k, v := range resp.Header {
		res.Headers[k] = strings.Join(v, ", ")
		val := strings.ToLower(res.Headers[k])

		if k == "X-Powered-By" {
			if strings.Contains(val, "php") {
				stack = "PHP"
				payload = "php"
			} else if strings.Contains(val, "express") || strings.Contains(val, "node") {
				stack = "Node.js"
				payload = "node"
			}
		}

		if k == "Server" {
			if strings.Contains(val, "gunicorn") || strings.Contains(val, "werkzeug") {
				stack = "Python (Flask/Django)"
				payload = "python"
			}
		}
	}

	res.Stack = stack
	res.Payload = payload

	return res, nil
}
