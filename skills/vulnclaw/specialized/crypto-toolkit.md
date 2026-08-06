# stage: exploit
# category: specialized

> Encoding/decoding and encryption/decryption toolkit — base64/URL/Hex/HTML entity encode/decode, MD5/SHA hashing, AES/DES/RSA encrypt/decrypt, JWT parsing, Caesar/ROT13 cipher, Rail Fence/Vigenere cipher, Unicode escape, Morse code, etc.

# Encoding/Decoding & Encryption/Decryption Skill

Provides comprehensive encoding/decoding and encryption/decryption capabilities for common encoding, encryption, and obfuscation scenarios in penetration testing.
**Important**: When encountering any encoded/encrypted string, prioritize using the `crypto_decode` tool to decode it, rather than guessing by intuition.

## Core Principles

1. **Tool-first** — When encountering base64, hex, URL-encoded strings, call the `crypto_decode` tool to decode; do not guess
2. **Multi-format attempts** — If one decoding method yields unreasonable results, try other encoding formats
3. **Chain decoding** — Multi-layer encoding is common in CTF (e.g., base64→hex→ROT13); after decoding, check whether the result needs further decoding
4. **Validate results** — After decoding, validate the result's reasonableness (is it readable text, does it look like a path/URL/flag, etc.)

## 1. Encoding Identification & Decoding

### Common Encoding Feature Identification

| Encoding Type | Features | Example |
|---------|------|------|
| Base64 | `A-Za-z0-9+/=` often has `=` padding at end | `TnNTY1RmLnBocA==` |
| Base32 | `A-Z2-7=` | `OBZHK5DFN2A====` |
| Hex | `0-9a-f` even length | `4e73536354662e706870` |
| URL encoding | `%XX` format | `%2F%61%64%6D%69%6E` |
| HTML entity | `&#xNN;` or `&#NNN;` | `&#x3C;script&#x3E;` |
| Unicode escape | `\uXXXX` or `\UXXXXXXXX` | `\u003c\u0073\u0063` |
| JWT | Three base64 segments separated by `.` | `eyJhbG...` |

### Decoding Strategy

1. Identify encoding type → call `crypto_decode` tool with the corresponding operation
2. Check whether the decoded result is readable/reasonable
3. If unreasonable, try other encoding formats
4. If the result still looks encoded, repeat steps 1-3

## 2. Hashes & Digests

### Common Hash Types

| Type | Output Length | Features |
|------|---------|------|
| MD5 | 32 hex | `e10adc3949ba59abbe56e057f20f883e` |
| SHA1 | 40 hex | `aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d` |
| SHA256 | 64 hex | `2c26b46b68ffc68ff99b453c1d30413413422d7064...` |
| SHA512 | 128 hex | Longer hex string |
| NTLM | 32 hex | Windows hash |
| MySQL5 | 41 chars | `*E6CC90B878B948C35E92B003C792C46758BF4` |

### Hash Processing Strategy

- Identify hash type (by length and character set)
- Try online rainbow table queries (via fetch tool accessing crackstation, etc.)
- For hashes with known salt, try salted brute force

## 3. Symmetric Encryption

### AES/DES/3DES

- Requires key and mode (ECB/CBC/CTR, etc.)
- CBC mode requires IV
- Common padding: PKCS7/ZeroPadding
- Hardcoded keys are often encountered in pentesting; prioritize extracting from source code

## 4. Asymmetric Encryption

### RSA

- Extract parameters from public/private key files
- RSA with small modulus can be factored
- Known private key allows direct decryption

## 5. Classical Ciphers

| Type | Features | Crack Method |
|------|------|---------|
| Caesar/ROT13 | Letter shift | Brute-force 25 shifts |
| Vigenere | Poly-alphabetic substitution | Kasiski/frequency analysis |
| Rail Fence | Character group rearrangement | Try common rail counts |
| Bacon cipher | AB quintuplets | Table lookup |
| Morse | `.-` dots and dashes | Table lookup |

## 6. JWT Processing

- Decode Header + Payload (base64url)
- Check algorithm: `none` algorithm bypass, RS256→HS256 algorithm confusion
- Try weak key signature forgery
- Check exp/nbf and other time claims

## Tool Usage

### `crypto_decode` Tool

When encountering operations requiring encoding/decoding/encryption/decryption, call this tool:

```
crypto_decode(operation="base64_decode", input="TnNTY1RmLnBocA==")
```

Supported operations:
- **Encoding** `base64_encode`, `base32_encode`, `hex_encode`, `url_encode`, `html_encode`, `unicode_encode`, `rot13_encode`, `morse_encode`, `caesar_encode`, `base58_encode`
- **Decoding** `base64_decode`, `base32_decode`, `hex_decode`, `url_decode`, `html_decode`, `unicode_decode`, `rot13_decode`, `morse_decode`, `caesar_decode`, `base58_decode`
- **Hashing** `md5_hash`, `sha1_hash`, `sha256_hash`, `sha512_hash`
- **Encryption/Decryption** `aes_encrypt`, `aes_decrypt`, `des_encrypt`, `des_decrypt`, `rsa_encrypt`, `rsa_decrypt`
- **JWT** `jwt_decode`, `jwt_encode`
- **Auto-detect**: `auto_decode` (auto-detect encoding type and decode)

## CTF Cryptographic Attack Routing

> When encountering cryptographic attack scenarios (known encryption algorithm, need to recover plaintext or key), prioritize the `ctf-crypto` skill:

| Attack Scenario | Route to ctf-crypto | Reference Doc |
|---------|-----------------|---------|
| RSA small exponent/common modulus/Wiener | `ctf-crypto` | `references/rsa-attacks-cheatsheet.md` |
| AES Padding Oracle/ECB flip | `ctf-crypto` | `references/aes-and-block-cipher-attacks.md` |
| ECC small subgroup/discrete log | `ctf-crypto` | `references/ecc-attacks-cheatsheet.md` |
| PRNG/MT19937 prediction | `ctf-crypto` | `references/prng-and-stream-cipher-attacks.md` |
| Classical ciphers (Vigenere/XOR) | `ctf-crypto` | `references/classic-cipher-attacks.md` |
| Lattice attack/LWE | `ctf-crypto` | `references/lattice-and-lwe-attacks.md` |

**This skill focuses on encoding/decoding operations and tools**; for specific cryptographic attack methods and parameters, refer to `ctf-crypto`.

## Reference Documents

- `references/encoding-cheatsheet.md` — Encoding identification quick reference
- `references/crypto-attacks.md` — Cryptographic attack techniques
- `references/crypto-attacks-roadmap.md` — Cryptographic attack classification routing (select attack method based on challenge features)

## References — crypto-attacks-roadmap

# Cryptographic Attack Classification Routing

Based on known information from the challenge, quickly determine which attack method to use.

## Decision Tree

```
Known conditions
├── Know plaintext + ciphertext?
│   ├── Same key encrypts multiple times? → XOR/stream cipher analysis
│   └── Single encryption? → Analyze encryption mode
├── Know ciphertext + key?
│   ├── Symmetric encryption → Direct decryption
│   └── Asymmetric encryption → RSA/ECC attack
├── Known n, e, c (RSA)?
│   ├── e is small → Small exponent attack
│   ├── Multiple share n → Common modulus attack
│   ├── d is small → Wiener's attack
│   ├── p-1 is smooth → Pollard p-1
│   └── Try online factoring (factordb)
├── Elliptic curve parameters?
│   ├── Smooth order → Pohlig-Hellman
│   ├── Anomalous curve → Smart's attack
│   └── ECDSA nonce reuse → Private key recovery
├── Known PRNG output sequence?
│   ├── MT19937 → State recovery
│   ├── LCG → Parameter recovery
│   └── LFSR → Berlekamp-Massey
└── Classical cipher?
    ├── Caesar/ROT13 → Brute force
    ├── Vigenere → Kasiski + frequency
    └── One-Time Pad reuse → Statistical attack
```

## RSA Attack Quick Selection

| Known | Attack |
|------|------|
| n, e, c, e=3 | Small exponent root |
| Multiple (n, c), same e, same plaintext | Håstad broadcast |
| Multiple (n, c), same n, different e | Common modulus attack |
| n, e, d very small approximation | Wiener's attack |
| n factorable, p≈q | Fermat factoring |
| n factorable, p-1 smooth | Pollard p-1 |
| Known partial plaintext | Coppersmith |
| factordb queryable | Online factoring |

## AES/Block Cipher Attack Quick Selection

| Scenario | Attack |
|------|------|
| ECB mode | Pattern analysis + block rearrangement |
| CBC mode, controllable IV | IV flip attack |
| CBC mode, Padding Oracle | Padding Oracle attack |
| CTR/GCM, nonce reuse | Keystream recovery |
| Known partial plaintext | XOR keystream recovery |

## PRNG Attack Quick Selection

| Scenario | Attack |
|------|------|
| Python random(), 624 outputs | MT19937 state recovery |
| 3 consecutive LCG outputs | Parameter recovery |
| LFSR output sequence | Berlekamp-Massey |
| RC4 (after dropping first 3072 bytes) | RC4 Drop attack |

## Classical Cipher Quick Selection

| Ciphertext Feature | Attack |
|---------|------|
| Single character substitution | Frequency analysis |
| Multi-character shift | Caesar brute force |
| Poly-alphabetic substitution | Vigenere Kasiski |
| Binary XOR multi-byte | Frequency analysis + key length estimation |
| One-time pad reuse | XOR comparison attack |

## References — crypto-attacks

# Cryptographic Attack Techniques

## 1. Hash Attacks

### Rainbow Table Queries
- crackstation.net — free, supports MD5/SHA1/SHA256
- cmd5.com — broad coverage
- hashes.org — community-maintained

### Hash Length Extension Attack
- Applicable: MD5, SHA1, SHA256 and other Merkle-Damgård-based hashes
- Condition: know `H(message)` and `len(message)`, but not the message itself
- Tools: hashpump, hash_extender
- Scenario: API signature verification bypass

### Hash Collisions
- MD5: fastcoll, HashClash
- SHA1: SHAttered (theoretically feasible)
- Scenario: file integrity bypass, certificate forgery

## 2. Symmetric Encryption Attacks

### ECB Mode Attack
- Same plaintext block → same ciphertext block
- Plaintext can be rearranged by rearranging ciphertext blocks
- Can identify repeated patterns (e.g., user role fields)

### CBC Byte Flip Attack
- Modify IV or previous ciphertext block to flip the corresponding byte in the next plaintext block
- Formula: `P[i] = D(C[i]) XOR C[i-1]`
- Modify `C[i-1][j]` → `P[i][j]` flips
- Scenario: modify encrypted user ID, role fields

### Padding Oracle Attack
- Condition: server returns whether padding is correct
- Recover plaintext byte-by-byte without the key
- Tools: padbuster, padding-oracle-attacker
- Scenario: ASP.NET, Java serialized tokens

### IV Reuse Attack
- Same IV + same Key in CBC mode → information leakage
- Can infer whether plaintexts are the same

## 3. RSA Attacks

### Small Public Exponent Attack
- When e=3, if plaintext m^3 < n, recover by taking the cube root directly
- Low encryption exponent broadcast attack: same plaintext encrypted with same e, different n

### Common Modulus Attack
- Same plaintext encrypted with same n, different e
- Recover plaintext via Extended Euclidean Algorithm

### Wiener's Attack
- Can factor n when d < n^0.25
- Applicable to small private key exponent scenarios

### Fermat Factoring
- Can quickly factor n when p and q are close
- Applicable to weak key generation

### Known Key File
- Extract parameters from .pem/.der files
- openssl rsa -text -noout -in key.pem

## 4. Classical Cipher Attacks

### Caesar Brute Force
- Only 25 possibilities, iterate directly
- Combine with frequency analysis to select the most likely result

### Vigenere Analysis
- Kasiski test to determine key length
- Index of coincidence method to verify key length
- After determining length, crack each column with Caesar

### Rail Fence Cipher
- Common rail counts: 2-8
- Iterate through all possible rail counts
- Check whether the result makes sense

### Bacon Cipher
- Two fonts/styles → A/B encoding
- Decode one letter per 5 characters

## 5. JWT Attacks

### none Algorithm Bypass
```json
{"alg": "none", "typ": "JWT"}
```
- Change the algorithm to none
- Remove the signature part
- Some implementations accept unsigned tokens

### RS256 → HS256 Algorithm Confusion
- Change the algorithm from RS256 to HS256
- Sign using the public key as the HMAC key
- If the server verifies HS256 signatures with the public key → bypass

### Weak Key Brute Force
- jwt-tool, jwt-cracker
- Common weak keys: secret, password, 123456, etc.

### JWK / jku Injection
- Embed public key in Header (jwk field)
- Or point to attacker-controlled jku URL
- If the server trusts the key in the Header → forgery

## 6. Encoding Chain Attack Patterns

### WAF Bypass Encoding
- Double URL encoding: `%2527` → `%27` → `'`
- Unicode normalization: `％27` → `'` (full-width to half-width)
- HTML entity: `&#39;` → `'`
- Base64-encode injection parameters

### Encoding in Deserialization
- PHP: base64-encoded serialized objects
- Java: Base64-encoded serialized byte streams
- Python: base64 pickle payloads

## 7. Tool Quick Reference

| Scenario | Tool |
|------|------|
| General encode/decode | CyberChef |
| Hash cracking | hashcat, john |
| RSA analysis | RsaCtfTool |
| JWT analysis | jwt-tool |
| Padding Oracle | padbuster |
| Hash extension | hashpump |
| Online decoding | base64decode.org, cyberchef.org |

## References — encoding-cheatsheet

# Encoding Identification Quick Reference

## Quick Identification Flow

```
Input string
  ├─ Contains %XX → URL encoding → url_decode
  ├─ Contains &# or &#x → HTML entity → html_decode
  ├─ Contains \uXXXX → Unicode escape → unicode_decode
  ├─ Contains .- with only dots, dashes, spaces → Morse → morse_decode
  ├─ Three base64 segments joined by . → JWT → jwt_decode
  ├─ Has = padding at end + A-Za-z0-9+/ → Base64 → base64_decode
  ├─ Has = padding at end + A-Z2-7 → Base32 → base32_decode
  ├─ Pure hex chars (0-9a-f) even length → Hex → hex_decode
  ├─ Pure uppercase + digits, no padding → Possibly Base58 → base58_decode
  ├─ Letter shift features (e.g., E→M, A→I) → Caesar → caesar_decode
  └─ Cannot determine → auto_decode
```

## Base64 Variants

| Variant | Character Set | Usage |
|------|--------|------|
| Standard Base64 | `A-Za-z0-9+/=` | General purpose |
| URL-safe Base64 | `A-Za-z0-9-_` | URL parameters |
| Base64url (JWT) | `A-Za-z0-9_-` no padding | JWT |

## Base58

| Variant | Excluded Chars | Usage |
|------|---------|------|
| Bitcoin | `0OIl` | Address encoding |
| Flickr | `0OIl` | Short URLs |
| Ripple | `0OIl` | Address encoding |

## Common Obfuscation Patterns

### Double Encoding
```
Original: admin
→ URL encoding: %61%64%6D%69%6E
→ Double URL encoding: %2561%2564%256D%2569%256E
```

### Base64 + Hex Chain
```
Original: NsScTf.php
→ Hex: 4e73536354662e706870
→ Base64: TnNTY1RmLnBocA==
```

### ROT13 Nesting
```
Original: password
→ ROT13: cnffjbeq
→ ROT13 again: password (ROT13 is self-inverse)
```

## Length & Encoding Reference

| Original Length | Base64 Length | Hex Length | Base32 Length |
|---------|------------|---------|------------|
| 1 byte | 4 chars | 2 chars | 8 chars |
| 4 bytes | 8 chars | 8 chars | 8 chars |
| 8 bytes | 12 chars | 16 chars | 16 chars |
| 16 bytes | 24 chars | 32 chars | 28 chars |

## Common CTF Encoding Chains

1. **Base64 → plaintext** — most common
2. **Base64 → Hex → plaintext** — double encoding
3. **Base64 → Base64 → plaintext** — nested Base64
4. **Hex → Base64 → ROT13 → plaintext** — triple encoding
5. **URL encoding → Base64 → plaintext** — common in Web scenarios
6. **Morse → Base64 → Hex → plaintext** — Crypto challenges

## Post-Decode Validation

After decoding, check whether the result:
- [ ] Is readable ASCII/UTF-8 text
- [ ] Looks like a path (/xxx/yyy.php)
- [ ] Looks like a URL (http://...)
- [ ] Contains flag format (flag{...}, NSSCTF{...})
- [ ] Is still encoded (needs further decoding)
