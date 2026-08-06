package core

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestCrypto_Base64(t *testing.T) {
	encoded, _ := cryptoExecute("base64_encode", "Hello", "", "", 0)
	if encoded != "SGVsbG8=" {
		t.Errorf("base64_encode: expected SGVsbG8=, got %s", encoded)
	}
	decoded, _ := cryptoExecute("base64_decode", "SGVsbG8=", "", "", 0)
	if decoded != "Hello" {
		t.Errorf("base64_decode: expected Hello, got %s", decoded)
	}
}

func TestCrypto_Hex(t *testing.T) {
	encoded, _ := cryptoExecute("hex_encode", "AB", "", "", 0)
	if encoded != "4142" {
		t.Errorf("hex_encode: expected 4142, got %s", encoded)
	}
	decoded, _ := cryptoExecute("hex_decode", "4142", "", "", 0)
	if decoded != "AB" {
		t.Errorf("hex_decode: expected AB, got %s", decoded)
	}
}

func TestCrypto_URL(t *testing.T) {
	encoded, _ := cryptoExecute("url_encode", "hello world&foo=bar", "", "", 0)
	if !strings.Contains(encoded, "hello+world") {
		t.Errorf("url_encode: unexpected result %s", encoded)
	}
	decoded, _ := cryptoExecute("url_decode", "hello+world", "", "", 0)
	if decoded != "hello world" {
		t.Errorf("url_decode: expected 'hello world', got %s", decoded)
	}
}

func TestCrypto_Rot13(t *testing.T) {
	encoded, _ := cryptoExecute("rot13", "Hello", "", "", 0)
	if encoded != "Uryyb" {
		t.Errorf("rot13: expected Uryyb, got %s", encoded)
	}
	// rot13 is self-inverse
	decoded, _ := cryptoExecute("rot13", encoded, "", "", 0)
	if decoded != "Hello" {
		t.Errorf("rot13 decode: expected Hello, got %s", decoded)
	}
}

func TestCrypto_Caesar(t *testing.T) {
	encoded, _ := cryptoExecute("caesar_encode", "abc", "", "", 3)
	if encoded != "def" {
		t.Errorf("caesar_encode: expected def, got %s", encoded)
	}
	decoded, _ := cryptoExecute("caesar_decode", "def", "", "", 3)
	if decoded != "abc" {
		t.Errorf("caesar_decode: expected abc, got %s", decoded)
	}
}

func TestCrypto_Morse(t *testing.T) {
	encoded, _ := cryptoExecute("morse_encode", "SOS", "", "", 0)
	if encoded != "... --- ..." {
		t.Errorf("morse_encode: expected '... --- ...', got %s", encoded)
	}
	decoded, _ := cryptoExecute("morse_decode", "... --- ...", "", "", 0)
	if decoded != "SOS" {
		t.Errorf("morse_decode: expected SOS, got %s", decoded)
	}
}

func TestCrypto_MD5(t *testing.T) {
	result, _ := cryptoExecute("md5_hash", "test", "", "", 0)
	if result != "098f6bcd4621d373cade4e832627b4f6" {
		t.Errorf("md5_hash: unexpected result %s", result)
	}
}

func TestCrypto_SHA256(t *testing.T) {
	result, _ := cryptoExecute("sha256_hash", "test", "", "", 0)
	if len(result) != 64 {
		t.Errorf("sha256_hash: expected 64 chars, got %d", len(result))
	}
}

func TestCrypto_AES(t *testing.T) {
	key := "0123456789abcdef" // 16 bytes
	plaintext := "Hello World"

	encrypted, err := cryptoExecute("aes_encrypt", plaintext, key, "", 0)
	if err != nil {
		t.Fatalf("aes_encrypt failed: %v", err)
	}
	// Verify it's hex
	if _, err := hex.DecodeString(encrypted); err != nil {
		t.Fatalf("aes_encrypt: expected hex output, got %s", encrypted)
	}

	decrypted, err := cryptoExecute("aes_decrypt", encrypted, key, "", 0)
	if err != nil {
		t.Fatalf("aes_decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("aes round-trip: expected %s, got %s", plaintext, decrypted)
	}
}

func TestCrypto_AES_CBC(t *testing.T) {
	key := "0123456789abcdef"
	iv := "abcdef0123456789"
	plaintext := "Secret message"

	encrypted, err := cryptoExecute("aes_encrypt", plaintext, key, iv, 0)
	if err != nil {
		t.Fatalf("aes_encrypt CBC failed: %v", err)
	}
	decrypted, err := cryptoExecute("aes_decrypt", encrypted, key, iv, 0)
	if err != nil {
		t.Fatalf("aes_decrypt CBC failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("aes CBC round-trip: expected %s, got %s", plaintext, decrypted)
	}
}

func TestCrypto_JWT(t *testing.T) {
	// Build a test JWT at runtime to avoid literal token in source
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"testuser","name":"Test User","iat":1516239022}`))
	jwt := header + "." + payload + ".TESTSIG"
	result, err := cryptoExecute("jwt_decode", jwt, "", "", 0)
	if err != nil {
		t.Fatalf("jwt_decode failed: %v", err)
	}
	if !strings.Contains(result, "HS256") {
		t.Error("expected HS256 in header")
	}
	if !strings.Contains(result, "Test User") {
		t.Error("expected 'Test User' in payload")
	}
}

func TestCrypto_JWT_NoneAlg(t *testing.T) {
	// Build alg=none JWT at runtime
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"testuser","name":"Test User"}`))
	jwt := header + "." + payload + "."
	result, _ := cryptoExecute("jwt_decode", jwt, "", "", 0)
	if !strings.Contains(result, "warning") {
		t.Error("expected warning for alg=none")
	}
}

func TestCrypto_AutoDecode(t *testing.T) {
	// Base64 encoded "Hello"
	result := autoDecode("SGVsbG8=")
	if result != "Hello" {
		t.Errorf("autoDecode base64: expected Hello, got %s", result)
	}
}

func TestCrypto_UnicodeDecode(t *testing.T) {
	result := unicodeDecode("\\u0048\\u0065\\u006c\\u006c\\u006f")
	if result != "Hello" {
		t.Errorf("unicodeDecode: expected Hello, got %s", result)
	}
}

func TestCrypto_Base58(t *testing.T) {
	encoded := base58Encode([]byte("Hello"))
	decoded, err := base58Decode(encoded)
	if err != nil {
		t.Fatalf("base58 round-trip failed: %v", err)
	}
	if string(decoded) != "Hello" {
		t.Errorf("base58 round-trip: expected Hello, got %s", string(decoded))
	}
}

func TestCrypto_Des(t *testing.T) {
	key := "8bytekey" // 8 bytes
	plaintext := "TestData"

	encrypted, err := cryptoExecute("des_encrypt", plaintext, key, "", 0)
	if err != nil {
		t.Fatalf("des_encrypt failed: %v", err)
	}
	decrypted, err := cryptoExecute("des_decrypt", encrypted, key, "", 0)
	if err != nil {
		t.Fatalf("des_decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("des round-trip: expected %s, got %s", plaintext, decrypted)
	}
}
