package encode
// internal/generator/encoder.go - XOR/BASE64 ENCODING
package generator

import (
    "encoding/base64"
    "fmt"
    "math/rand"
    "strings"
)

type Encoder struct {
    key     []byte
    charset string
}

func NewEncoder() *Encoder {
    return &Encoder{
        charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
    }
}

func (e *Encoder) XOREncode(data []byte) (string, []byte) {
    // GENERATE RANDOM KEY
    keyLength := rand.Intn(16) + 8
    key := make([]byte, keyLength)
    for i := range key {
        key[i] = byte(rand.Intn(256))
    }
    
    // XOR ENCRYPT
    encrypted := make([]byte, len(data))
    for i := range data {
        encrypted[i] = data[i] ^ key[i%len(key)]
    }
    
    // RETURN BASE64 ENCODED
    return base64.StdEncoding.EncodeToString(encrypted), key
}

func (e *Encoder) XORDecode(encoded string, key []byte) ([]byte, error) {
    encrypted, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return nil, err
    }
    
    decrypted := make([]byte, len(encrypted))
    for i := range encrypted {
        decrypted[i] = encrypted[i] ^ key[i%len(key)]
    }
    
    return decrypted, nil
}

func (e *Encoder) CustomBase64(data []byte) string {
    // CUSTOM BASE64 WITH SHUFFLED CHARSET
    shuffled := e.shuffleString(e.charset + "+/")
    
    stdEncoding := base64.StdEncoding.EncodeToString(data)
    
    // REPLACE WITH CUSTOM CHARS
    var result strings.Builder
    for _, c := range stdEncoding {
        idx := strings.IndexByte(base64.StdEncoding.EncodeToString([]byte{byte(c)})[0], byte(c))
        if idx >= 0 {
            result.WriteByte(shuffled[idx])
        } else {
            result.WriteByte(byte(c))
        }
    }
    
    return result.String()
}

func (e *Encoder) shuffleString(s string) string {
    runes := []rune(s)
    for i := range runes {
        j := rand.Intn(i + 1)
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func (e *Encoder) GenerateDecoder(key []byte, customCharset string) string {
    // GENERATE PHP DECODER CODE
    decoder := `
function _custom_decode($data) {
    $key = "%s";
    $decoded = '';
    $encrypted = base64_decode($data);
    for($i=0; $i<strlen($encrypted); $i++) {
        $decoded .= $encrypted[$i] ^ $key[$i %% strlen($key)];
    }
    return $decoded;
}
`
    return fmt.Sprintf(decoder, base64.StdEncoding.EncodeToString(key))
}