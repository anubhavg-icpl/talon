package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// cryptoToolSpec returns the agent-facing crypto_decode tool spec.
// Ported from pentest agent skills/crypto_tools.py (~29 operations).
func cryptoToolSpec() llm.ToolSpec {
	ops := strings.Join(cryptoOpNames(), ", ")
	return llm.ToolSpec{
		Name:        "crypto_decode",
		Description: "Encode/decode/hash/encrypt/decrypt operations for penetration testing. " +
			"Use this instead of guessing: supports " + ops + ". " +
			"Call with operation + input (+ optional key/shift for ciphers).",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"operation": map[string]any{
					"type":        "string",
					"description": "The crypto operation to perform",
				},
				"input": map[string]any{
					"type":        "string",
					"description": "The input string to process",
				},
				"key": map[string]any{
					"type":        "string",
					"description": "Key for AES/DES operations",
				},
				"iv": map[string]any{
					"type":        "string",
					"description": "IV for AES/DES CBC mode",
				},
				"shift": map[string]any{
					"type":        "integer",
					"description": "Shift for Caesar cipher (default 3)",
				},
			},
			"required": []any{"operation", "input"},
		},
	}
}

func cryptoOpNames() []string {
	return []string{
		"base64_encode", "base64_decode", "base32_encode", "base32_decode",
		"base58_encode", "base58_decode", "hex_encode", "hex_decode",
		"url_encode", "url_decode", "html_encode", "html_decode",
		"unicode_encode", "unicode_decode", "rot13", "caesar_encode", "caesar_decode",
		"morse_encode", "morse_decode",
		"md5_hash", "sha1_hash", "sha256_hash", "sha512_hash",
		"aes_encrypt", "aes_decrypt", "des_encrypt", "des_decrypt",
		"jwt_decode", "auto_decode",
	}
}

// handleCryptoDecode dispatches to the appropriate crypto operation.
func handleCryptoDecode(args map[string]any, tr *tracker) (string, bool) {
	op, _ := args["operation"].(string)
	input, _ := args["input"].(string)
	key, _ := args["key"].(string)
	iv, _ := args["iv"].(string)
	shift := 3
	if s, ok := args["shift"].(float64); ok {
		shift = int(s)
	}

	result, err := cryptoExecute(op, input, key, iv, shift)
	if err != nil {
		if tr != nil {
			tr.record("crypto_decode", args, fmt.Sprintf("op=%s error=%s", op, err))
		}
		return fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error()), true
	}

	payload := map[string]any{
		"success":   true,
		"operation": op,
		"result":    result,
	}
	if tr != nil {
		tr.record("crypto_decode", args, fmt.Sprintf("op=%s ok=true", op))
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	return string(raw), false
}

// cryptoExecute performs the actual crypto operation.
func cryptoExecute(op, input, key, iv string, shift int) (string, error) {
	switch op {
	// Base encodings
	case "base64_encode":
		return base64.StdEncoding.EncodeToString([]byte(input)), nil
	case "base64_decode":
		data, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return base64DecodeLenient(input)
		}
		return string(data), nil
	case "base32_encode":
		return base32.StdEncoding.EncodeToString([]byte(input)), nil
	case "base32_decode":
		data, err := base32.StdEncoding.DecodeString(strings.ToUpper(input))
		if err != nil {
			return "", fmt.Errorf("base32 decode failed: %v", err)
		}
		return string(data), nil
	case "base58_encode":
		return base58Encode([]byte(input)), nil
	case "base58_decode":
		data, err := base58Decode(input)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "hex_encode":
		return hex.EncodeToString([]byte(input)), nil
	case "hex_decode":
		data, err := hex.DecodeString(input)
		if err != nil {
			return "", fmt.Errorf("hex decode failed: %v", err)
		}
		return string(data), nil

	// Web encodings
	case "url_encode":
		return url.QueryEscape(input), nil
	case "url_decode":
		decoded, err := url.QueryUnescape(input)
		if err != nil {
			return "", fmt.Errorf("url decode failed: %v", err)
		}
		return decoded, nil
	case "html_encode":
		return html.EscapeString(input), nil
	case "html_decode":
		return html.UnescapeString(input), nil

	// Unicode
	case "unicode_encode":
		var sb strings.Builder
		for _, r := range input {
			if r > 127 {
				sb.WriteString(fmt.Sprintf("\\u%04x", r))
			} else {
				sb.WriteRune(r)
			}
		}
		return sb.String(), nil
	case "unicode_decode":
		return unicodeDecode(input), nil

	// Classical ciphers
	case "rot13":
		return caesarShift(input, 13), nil
	case "caesar_encode":
		return caesarShift(input, shift), nil
	case "caesar_decode":
		return caesarShift(input, -shift), nil
	case "morse_encode":
		return morseEncode(input), nil
	case "morse_decode":
		return morseDecode(input), nil

	// Hashes
	case "md5_hash":
		h := md5.Sum([]byte(input))
		return hex.EncodeToString(h[:]), nil
	case "sha1_hash":
		h := sha1.Sum([]byte(input))
		return hex.EncodeToString(h[:]), nil
	case "sha256_hash":
		h := sha256.Sum256([]byte(input))
		return hex.EncodeToString(h[:]), nil
	case "sha512_hash":
		h := sha512.Sum512([]byte(input))
		return hex.EncodeToString(h[:]), nil

	// Symmetric encryption
	case "aes_encrypt":
		return aesEncrypt(input, key, iv)
	case "aes_decrypt":
		return aesDecrypt(input, key, iv)
	case "des_encrypt":
		return desEncrypt(input, key)
	case "des_decrypt":
		return desDecrypt(input, key)

	// JWT
	case "jwt_decode":
		return jwtDecode(input)

	// Auto-detect
	case "auto_decode":
		return autoDecode(input), nil

	default:
		return "", fmt.Errorf("unknown operation: %s", op)
	}
}

// --- Base58 ---

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func base58Encode(input []byte) string {
	if len(input) == 0 {
		return ""
	}
	// Count leading zeros
	zeros := 0
	for zeros < len(input) && input[zeros] == 0 {
		zeros++
	}
	// Convert to big integer
	num := make([]byte, 0, len(input)*138/100+1)
	for _, b := range input {
		carry := int(b)
		for j := 0; j < len(num); j++ {
			carry += int(num[j]) << 8
			num[j] = byte(carry % 58)
			carry /= 58
		}
		for carry > 0 {
			num = append(num, byte(carry%58))
			carry /= 58
		}
	}
	// Convert to base58 string
	result := make([]byte, 0, len(num)+zeros)
	for i := 0; i < zeros; i++ {
		result = append(result, base58Alphabet[0])
	}
	for i := len(num) - 1; i >= 0; i-- {
		result = append(result, base58Alphabet[num[i]])
	}
	return string(result)
}

func base58Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return []byte{}, nil
	}
	// Map characters to values
	alphabetMap := make(map[byte]int)
	for i, c := range base58Alphabet {
		alphabetMap[byte(c)] = i
	}
	// Count leading '1's (zeros)
	zeros := 0
	for zeros < len(input) && input[zeros] == base58Alphabet[0] {
		zeros++
	}
	// Convert
	bytes := make([]byte, 0, len(input)*733/1000+1)
	for _, c := range input {
		val, ok := alphabetMap[byte(c)]
		if !ok {
			return nil, fmt.Errorf("invalid base58 character: %c", c)
		}
		carry := val
		for j := 0; j < len(bytes); j++ {
			carry += int(bytes[j]) * 58
			bytes[j] = byte(carry & 0xff)
			carry >>= 8
		}
		for carry > 0 {
			bytes = append(bytes, byte(carry&0xff))
			carry >>= 8
		}
	}
	// Add leading zeros
	result := make([]byte, 0, len(bytes)+zeros)
	for i := 0; i < zeros; i++ {
		result = append(result, 0)
	}
	for i := len(bytes) - 1; i >= 0; i-- {
		result = append(result, bytes[i])
	}
	return result, nil
}

// --- Classical ciphers ---

func caesarShift(input string, shift int) string {
	shift = ((shift % 26) + 26) % 26 // Normalize to 0-25
	var sb strings.Builder
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			sb.WriteRune(rune(((int(r-'a')+shift)%26 + 26%26) % 26) + 'a')
		case r >= 'A' && r <= 'Z':
			sb.WriteRune(rune(((int(r-'A')+shift)%26 + 26%26) % 26) + 'A')
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// --- Morse code ---

var morseMap = map[rune]string{
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".",
	'F': "..-.", 'G': "--.", 'H': "....", 'I': "..", 'J': ".---",
	'K': "-.-", 'L': ".-..", 'M': "--", 'N': "-.", 'O': "---",
	'P': ".--.", 'Q': "--.-", 'R': ".-.", 'S': "...", 'T': "-",
	'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-", 'Y': "-.--",
	'Z': "--..", '0': "-----", '1': ".----", '2': "..---", '3': "...--",
	'4': "....-", '5': ".....", '6': "-....", '7': "--...", '8': "---..",
	'9': "----.", '.': ".-.-.-", ',': "--..--", '?': "..--..", '!': "-.-.--",
	'/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...", ':': "---...",
	';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-", '_': "..--.-",
	'"': ".-..-.", '@': ".--.-.",
}

var morseReverseMap map[string]rune

func init() {
	morseReverseMap = make(map[string]rune, len(morseMap))
	for k, v := range morseMap {
		morseReverseMap[v] = k
	}
}

func morseEncode(input string) string {
	var parts []string
	for _, r := range strings.ToUpper(input) {
		if r == ' ' {
			parts = append(parts, "/")
			continue
		}
		if m, ok := morseMap[r]; ok {
			parts = append(parts, m)
		}
	}
	return strings.Join(parts, " ")
}

func morseDecode(input string) string {
	var sb strings.Builder
	words := strings.Split(input, "/")
	for wi, word := range words {
		if wi > 0 {
			sb.WriteRune(' ')
		}
		for _, sym := range strings.Fields(word) {
			if r, ok := morseReverseMap[sym]; ok {
				sb.WriteRune(r)
			}
		}
	}
	return sb.String()
}

// --- Unicode decode ---

var unicodeRe = regexp.MustCompile(`\\u([0-9a-fA-F]{4})|\\U([0-9a-fA-F]{8})`)

func unicodeDecode(input string) string {
	return unicodeRe.ReplaceAllStringFunc(input, func(match string) string {
		var hexStr string
		if len(match) > 2 && match[1] == 'U' {
			hexStr = match[3:]
		} else {
			hexStr = match[2:]
		}
		val, err := strconv.ParseInt(hexStr, 16, 32)
		if err != nil {
			return match
		}
		r := rune(val)
		if !utf8.ValidRune(r) {
			return match
		}
		return string(r)
	})
}

// --- AES ---

func aesEncrypt(input, key, iv string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("AES key must be 16/24/32 bytes, got %d", len(key))
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	mode := "ECB"
	if iv != "" {
		if len(iv) != 16 {
			return "", fmt.Errorf("AES IV must be 16 bytes, got %d", len(iv))
		}
		mode = "CBC"
	}

	padded := pkcs7Pad([]byte(input), block.BlockSize())
	if mode == "CBC" {
		encrypted := make([]byte, len(padded))
		cipher.NewCBCEncrypter(block, []byte(iv)).CryptBlocks(encrypted, padded)
		return hex.EncodeToString(encrypted), nil
	}
	// ECB
	encrypted := make([]byte, len(padded))
	for i := 0; i < len(padded); i += block.BlockSize() {
		block.Encrypt(encrypted[i:i+block.BlockSize()], padded[i:i+block.BlockSize()])
	}
	return hex.EncodeToString(encrypted), nil
}

func aesDecrypt(inputHex, key, iv string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", fmt.Errorf("AES key must be 16/24/32 bytes, got %d", len(key))
	}
	ciphertext, err := hex.DecodeString(inputHex)
	if err != nil {
		return "", fmt.Errorf("invalid hex input: %v", err)
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	if len(ciphertext)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext length %d is not a multiple of block size %d",
			len(ciphertext), block.BlockSize())
	}
	decrypted := make([]byte, len(ciphertext))
	if iv != "" {
		if len(iv) != 16 {
			return "", fmt.Errorf("AES IV must be 16 bytes, got %d", len(iv))
		}
		cipher.NewCBCDecrypter(block, []byte(iv)).CryptBlocks(decrypted, ciphertext)
	} else {
		// ECB
		for i := 0; i < len(ciphertext); i += block.BlockSize() {
			block.Decrypt(decrypted[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
		}
	}
	unpadded, err := pkcs7Unpad(decrypted, block.BlockSize())
	if err != nil {
		return string(decrypted), nil // Return raw if padding is broken
	}
	return string(unpadded), nil
}

// --- DES ---

func desEncrypt(input, key string) (string, error) {
	if len(key) != 8 {
		return "", fmt.Errorf("DES key must be 8 bytes, got %d", len(key))
	}
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	padded := pkcs7Pad([]byte(input), block.BlockSize())
	encrypted := make([]byte, len(padded))
	for i := 0; i < len(padded); i += block.BlockSize() {
		block.Encrypt(encrypted[i:i+block.BlockSize()], padded[i:i+block.BlockSize()])
	}
	return hex.EncodeToString(encrypted), nil
}

func desDecrypt(inputHex, key string) (string, error) {
	if len(key) != 8 {
		return "", fmt.Errorf("DES key must be 8 bytes, got %d", len(key))
	}
	ciphertext, err := hex.DecodeString(inputHex)
	if err != nil {
		return "", fmt.Errorf("invalid hex input: %v", err)
	}
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(decrypted[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}
	unpadded, err := pkcs7Unpad(decrypted, block.BlockSize())
	if err != nil {
		return string(decrypted), nil
	}
	return string(unpadded), nil
}

// --- JWT ---

func jwtDecode(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT: expected 3 parts separated by '.', got %d", len(parts))
	}
	headerJSON, err := jwtBase64Decode(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to decode header: %v", err)
	}
	payloadJSON, err := jwtBase64Decode(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode payload: %v", err)
	}
	var header, payload map[string]any
	json.Unmarshal([]byte(headerJSON), &header)
	json.Unmarshal([]byte(payloadJSON), &payload)
	result := map[string]any{
		"header":   header,
		"payload":  payload,
		"raw_header": headerJSON,
		"raw_payload": payloadJSON,
		"signature": parts[2],
	}
	if alg, ok := header["alg"].(string); ok {
		if alg == "none" {
			result["warning"] = "alg=none detected: this JWT can be forged without a signature"
		}
		if strings.HasPrefix(alg, "HS") {
			result["note"] = "HS alg detected: vulnerable to RS256→HS256 confusion if server verifies with public key"
		}
	}
	raw, _ := json.MarshalIndent(result, "", "  ")
	return string(raw), nil
}

func jwtBase64Decode(s string) (string, error) {
	// JWT uses base64url without padding
	padded := s
	if l := len(s) % 4; l > 0 {
		padded += strings.Repeat("=", 4-l)
	}
	data, err := base64.URLEncoding.DecodeString(padded)
	if err != nil {
		data, err = base64.StdEncoding.DecodeString(padded)
		if err != nil {
			return "", err
		}
	}
	return string(data), nil
}

// --- Auto decode ---

func autoDecode(input string) string {
	result := input
	for i := 0; i < 10; i++ {
		decoded, changed := autoDecodeStep(result)
		if !changed {
			break
		}
		result = decoded
	}
	return result
}

func autoDecodeStep(input string) (string, bool) {
	// Try URL decode
	if decoded, err := url.QueryUnescape(input); err == nil && decoded != input {
		return decoded, true
	}
	// Try HTML entity decode
	if decoded := html.UnescapeString(input); decoded != input {
		return decoded, true
	}
	// Try base64
	if data, err := base64.StdEncoding.DecodeString(input); err == nil && isPrintable(data) {
		return string(data), true
	}
	// Try hex
	if isHexString(input) {
		if data, err := hex.DecodeString(input); err == nil && isPrintable(data) {
			return string(data), true
		}
	}
	return input, false
}

func isPrintable(data []byte) bool {
	printable := 0
	for _, b := range data {
		if b >= 32 && b <= 126 || b == '\n' || b == '\r' || b == '\t' {
			printable++
		}
	}
	return len(data) > 0 && printable*100/len(data) > 70
}

func isHexString(s string) bool {
	if len(s)%2 != 0 || len(s) < 2 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// --- PKCS7 padding ---

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return data, fmt.Errorf("invalid padding")
	}
	padding := int(data[len(data)-1])
	if padding < 1 || padding > blockSize {
		return data, fmt.Errorf("invalid padding length")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return data, fmt.Errorf("invalid padding bytes")
		}
	}
	return data[:len(data)-padding], nil
}

// base64DecodeLenient tries to decode with padding adjustments.
func base64DecodeLenient(input string) (string, error) {
	// Try with padding
	padded := input
	if l := len(input) % 4; l > 0 {
		padded += strings.Repeat("=", 4-l)
	}
	data, err := base64.StdEncoding.DecodeString(padded)
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(padded)
		if err != nil {
			return "", fmt.Errorf("base64 decode failed: %v", err)
		}
	}
	return string(data), nil
}

// CryptoExecutePublic is a public wrapper around cryptoExecute for use by
// the control plane (REST API crypto/decode endpoint).
func CryptoExecutePublic(op, input, key, iv string, shift int) (string, error) {
	return cryptoExecute(op, input, key, iv, shift)
}

// CryptoOpNamesPublic returns the list of supported crypto operations.
func CryptoOpNamesPublic() []string {
	return cryptoOpNames()
}
