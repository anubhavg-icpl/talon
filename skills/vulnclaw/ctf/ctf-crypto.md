# stage: exploit
# category: ctf

> CTF cryptography attack knowledge base — RSA attacks (small exponent / common modulus / Wiener / Coppersmith), AES attacks (Padding Oracle / ECB byte flipping / GCM nonce reuse), ECC attacks, LFSR/LCG/PRNG attacks, classical ciphers, LWE lattice attacks

# CTF Cryptography Attack Knowledge Base

A practical attack knowledge base for CTF Crypto challenges, providing **concrete attack parameters, mathematical formulas, and Python code snippets**.

**Difference from `crypto-toolkit`**:
- `crypto-toolkit` → encoding/decoding utility (base64 decoding, MD5 hashing, AES encrypt/decrypt)
- `ctf-crypto` → cryptographic attack knowledge (how to perform an RSA small-exponent attack, how to exploit Padding Oracle)

## Core Principles

1. **Identify the cryptosystem first** — look at key length, cipher mode, and known quantities to determine the attack direction
2. **Verify with tools** — use `python_execute` to run attack code, use `crypto_decode` for auxiliary encoding/decoding
3. **Parameter sensitivity** — cryptographic attacks are extremely sensitive to parameters; compute precisely

## Scenario Routing

| Scenario | Reference doc | Core attacks |
|------|---------|---------|
| RSA attacks | `rsa-attacks-cheatsheet.md` | small e / common modulus / Wiener / Pollard / Fermat / Coppersmith |
| AES/block cipher attacks | `aes-and-block-cipher-attacks.md` | ECB flipping / Padding Oracle / GCM nonce reuse |
| ECC attacks | `ecc-attacks-cheatsheet.md` | small subgroup / invalid curve / Smart / Pohlig-Hellman |
| PRNG/stream cipher attacks | `prng-and-stream-cipher-attacks.md` | MT19937 / LCG / LFSR / RC4 |
| Classical ciphers | `classic-cipher-attacks.md` | Vigenere / XOR frequency analysis / OTP reuse |
| Lattice attacks | `lattice-and-lwe-attacks.md` | LLL / BKZ / HNP / LWE embedding |

## Quick Challenge Identification Guide

| Challenge feature | Likely attack | Recommended reference |
|---------|---------|---------|
| Given n, e, c | RSA | rsa-attacks-cheatsheet.md |
| e=3 or very small e | RSA small-exponent attack | rsa-attacks-cheatsheet.md |
| Multiple (n, e, c) sets with the same n | RSA common-modulus attack | rsa-attacks-cheatsheet.md |
| n is large but e is large | Wiener attack | rsa-attacks-cheatsheet.md |
| AES-CBC + decryption oracle | Padding Oracle | aes-and-block-cipher-attacks.md |
| AES-ECB + controllable plaintext | ECB byte flipping | aes-and-block-cipher-attacks.md |
| Elliptic curve parameters | ECC attacks | ecc-attacks-cheatsheet.md |
| Given a random number sequence | PRNG prediction | prng-and-stream-cipher-attacks.md |
| Given ciphertext and partial plaintext | XOR/stream cipher | classic-cipher-attacks.md |
| Matrix/vector operations | Lattice attacks | lattice-and-lwe-attacks.md |

## References — aes-and-block-cipher-attacks

# AES and Block Cipher Attacks

## Cipher Mode Cheat Sheet

| Mode | Characteristic | Exploitable weakness |
|------|------|-----------|
| ECB | Same plaintext → same ciphertext | Pattern recognition, block rearrangement attacks |
| CBC | Previous ciphertext block feeds into current encryption | IV flipping, Padding Oracle |
| CTR | Stream-style encryption | nonce reuse → XOR leakage |
| CFB | Similar to stream cipher | IV flipping |
| OFB | Similar to stream cipher | nonce reuse |
| GCM | Authenticated encryption | nonce reuse → keystream recovery |

## ECB Byte Flipping

```python
from Crypto.Cipher import AES

# In ECB mode, identical plaintext blocks produce identical ciphertext blocks
# Attack: identify repeated ciphertext blocks → infer plaintext structure
# Ciphertext blocks can be rearranged to alter plaintext structure

def ecb_detect(ciphertext, block_size=16):
    """Detect ECB mode (look for repeated blocks)"""
    blocks = [ciphertext[i:i+block_size] for i in range(0, len(ciphertext), block_size)]
    return len(blocks) != len(set(blocks))
```

## CBC IV Flipping Attack

```python
"""
Principle: in CBC, P[i] = Decrypt(C[i]) XOR C[i-1]
Modifying a byte of C[i-1] → the corresponding byte of P[i] is also flipped

Use: modifying the IV changes the first plaintext block; modifying C[i-1] changes the i-th plaintext block
Cost: the plaintext P[i-1] corresponding to C[i-1] gets corrupted
"""

def cbc_iv_flip(ciphertext, known_plain, target_plain, block_size=16):
    """Flip the first CBC plaintext block (by modifying the IV)"""
    iv = bytearray(ciphertext[:block_size])
    for i in range(block_size):
        iv[i] = iv[i] ^ known_plain[i] ^ target_plain[i]
    return bytes(iv) + ciphertext[block_size:]
```

## Padding Oracle Attack

```python
"""
Principle: during CBC decryption, if the padding is invalid the server returns a different error
By brute-forcing byte by byte, the error/valid difference reveals the plaintext

Conditions:
1. CBC mode is used
2. The server returns different responses for padding errors vs. other ciphertext errors
3. Modified ciphertexts can be submitted repeatedly
"""

def padding_oracle_attack(oracle, ciphertext, block_size=16):
    """Recover plaintext with a Padding Oracle attack
    
    oracle: function that accepts a ciphertext and returns True (padding valid) / False (padding invalid)
    """
    blocks = [ciphertext[i:i+block_size] for i in range(0, len(ciphertext), block_size)]
    plaintext = b''
    
    for block_idx in range(1, len(blocks)):
        prev_block = bytearray(blocks[block_idx - 1])
        curr_block = blocks[block_idx]
        intermediate = bytearray(block_size)
        
        for byte_pos in range(block_size - 1, -1, -1):
            padding_val = block_size - byte_pos
            
            # Construct test ciphertext
            test_block = bytearray(block_size)
            for k in range(byte_pos + 1, block_size):
                test_block[k] = intermediate[k] ^ padding_val
            
            found = False
            for guess in range(256):
                test_block[byte_pos] = guess
                test_cipher = bytes(test_block) + curr_block
                
                if oracle(test_cipher):
                    intermediate[byte_pos] = guess ^ padding_val
                    found = True
                    break
            
            if not found:
                raise Exception(f"Padding oracle attack failed at byte {byte_pos}")
        
        # Recover plaintext
        for i in range(block_size):
            plaintext += bytes([intermediate[i] ^ prev_block[i]])
    
    return plaintext
```

## GCM Nonce Reuse Attack

```python
"""
When the same nonce is used for two encryptions:
- Both encryptions use the same keystream
- C1 = P1 XOR keystream
- C2 = P2 XOR keystream
- C1 XOR C2 = P1 XOR P2

If P1 is known, P2 can be recovered
"""

def gcm_nonce_reuse(c1, c2, p1):
    """Recover plaintext by exploiting GCM nonce reuse"""
    return bytes(a ^ b ^ c for a, b, c in zip(c1, c2, p1))
```

## CTR Nonce Reuse

```python
"""
In CTR mode, nonce reuse is equivalent to stream-cipher key reuse
C1 = P1 XOR keystream
C2 = P2 XOR keystream
C1 XOR C2 = P1 XOR P2
"""

def ctr_nonce_reuse(c1, c2, known_p1):
    """Recover plaintext by exploiting CTR nonce reuse"""
    return bytes(a ^ b ^ c for a, b, c in zip(c1, c2, known_p1))
```

## References — classic-cipher-attacks

# Classical Cipher Attacks

## Caesar Cipher

```python
def caesar_break(ciphertext):
    """Iterate over all shifts"""
    for shift in range(26):
        result = ""
        for c in ciphertext:
            if c.isalpha():
                base = ord('A') if c.isupper() else ord('a')
                result += chr((ord(c) - base + shift) % 26 + base)
            else:
                result += c
        print(f"Shift {shift}: {result}")
```

## Vigenère Cipher

```python
def vigenere_break(ciphertext, max_keylen=20):
    """Break Vigenère with Kasiski + frequency analysis"""
    from collections import Counter

    # 1. Kasiski: find repeated sequences, estimate key length
    def kasiski(text):
        distances = []
        for length in range(3, 6):
            seqs = {}
            for i in range(len(text) - length):
                seq = text[i:i+length]
                if seq in seqs:
                    distances.append(i - seqs[seq])
                seqs[seq] = i
        return distances

    # 2. Index of Coincidence (IC) to estimate key length
    def ic(text):
        freq = Counter(text.upper())
        n = len(text)
        return sum(f * (f - 1) for f in freq.values()) / (n * (n - 1))

    # 3. Frequency analysis to solve a single letter
    def solve_char(text, key_char):
        ENGLISH_FREQ = 'ETAOINSHRDLCUMWFGYPBVKJXQZ'
        key_base = ord(key_char.upper()) - ord('A')
        best_score = 0
        best_char = 'E'
        for shift in range(26):
            freq = Counter()
            for c in text:
                if c.isalpha():
                    shifted = chr((ord(c.upper()) - ord('A') - shift) % 26 + ord('A'))
                    freq[shifted] += 1
            score = sum(ENGLISH_FREQ.index(k) * freq[k] for k in freq if k in ENGLISH_FREQ)
            if score > best_score:
                best_score = score
                best_char = chr(ord('A') + shift)
        return best_char
```

## Multi-Byte XOR Encryption

```python
def multi_byte_xor_break(ciphertext, max_keylen=16):
    """Multi-byte XOR attack: Hamming distance + frequency analysis"""
    from collections import Counter

    def hamming_distance(b1, b2):
        return sum(bin(a ^ b).count('1') for a, b in zip(b1, b2))

    # Use Hamming distance to estimate key length
    best_keylen = 1
    best_score = float('inf')
    for keylen in range(2, max_keylen + 1):
        chunks = [ciphertext[i:i+keylen] for i in range(0, len(ciphertext), keylen)]
        avg_dist = sum(hamming_distance(c1, c2) for c1, c2 in zip(chunks[:4], chunks[1:5])) / 4
        normalized = avg_dist / keylen
        if normalized < best_score:
            best_score = normalized
            best_keylen = keylen

    # Group by key length; run single-byte XOR on each group
    key = b''
    for i in range(best_keylen):
        block = bytes(ciphertext[j] for j in range(i, len(ciphertext), best_keylen))
        # Frequency analysis to find the best single-byte key
        best = 0
        best_score = 0
        for k in range(256):
            decrypted = bytes(b ^ k for b in block)
            score = sum(1 for b in decrypted if chr(b).isalpha() or chr(b).isspace())
            if score > best_score:
                best_score = score
                best = k
        key += bytes([best])

    return key
```

## One-Time Pad (OTP) Reuse Attack

```python
"""
If the same OTP key is used to encrypt two messages:
C1 = P1 XOR key
C2 = P2 XOR key
C1 XOR C2 = P1 XOR P2

Exploit language redundancy (English word frequencies) to break it
"""
from collections import Counter

def otp_reuse_attack(c1, c2):
    """OTP key reuse attack"""
    xor_result = bytes(a ^ b for a, b in zip(c1, c2))
    # Recover plaintext via frequency analysis
```

## Rail Fence Cipher

```python
def railfence_break(ciphertext, max_rails=10):
    """Iterate over rail counts to decrypt"""
    for rails in range(2, max_rails + 1):
        # Rebuild the fence structure
        fence = [[] for _ in range(rails)]
        rail = 0
        direction = 1
        for c in ciphertext:
            fence[rail].append(c)
            rail += direction
            if rail == 0 or rail == rails - 1:
                direction = -direction
        # Read row by row
        result = ''.join(''.join(row) for row in fence)
        print(f"Rails {rails}: {result}")
```

## References — ecc-attacks-cheatsheet

# ECC Attack Cheat Sheet

## Elliptic Curve Basics

```python
# Elliptic curve: y² = x³ + ax + b (mod p)
# Point operations: P + Q, k*P
# ECDLP: given P, Q=k*P, find k
```

## Attack Selection

| Condition | Attack method | Applicable scenario |
|------|---------|---------|
| Order n is smooth | Pohlig-Hellman | All factors of n are small |
| Anomalous curve (p=n) | Smart attack | Anomalous curves |
| Small subgroup order | Small subgroup attack | Order has a large prime factor |
| Suspicious curve parameters | Invalid Curve attack | Non-standard curves |
| ECDSA nonce reuse | Deterministic attack | Same k used to sign twice |
| Very small order | Brute force / Baby-step Giant-step | n < 2^40 |

## Pohlig-Hellman Attack

```python
# Sage implementation
# When all factors of the group order n are small

P = EllipticCurve(GF(p), [a, b])
G = P(P_x, P_y)  # base point
Q = P(Q_x, Q_y)  # target point

n = P.order()  # group order
factors = factor(n)

# Pohlig-Hellman
k = discrete_log(Q, G, operation='+')
# or specify the method
k = Q.discrete_log(G)
```

## Smart Attack (Anomalous Curves)

```python
# When the curve's order equals the characteristic p (anomalous curve)
# E.lift_x() may fail but p-adic lifting can be exploited

# Sage implementation
def smart_attack(P, Q, p, a, b):
    """Smart attack, applicable to anomalous curves with #E = p"""
    E = EllipticCurve(Qp(p), [a, b])
    P_lift = E.lift_x(ZZ(P.xy()[0]))
    Q_lift = E.lift_x(ZZ(Q.xy()[0]))
    
    pP = p * P_lift
    pQ = p * Q_lift
    
    x1 = pP.xy()[0] / pP.xy()[1]
    x2 = pQ.xy()[0] / pQ.xy()[1]
    
    k = ZZ(x2) / ZZ(x1) % p
    return k
```

## Invalid Curve Attack

```python
# When the server does not verify that a point is on the curve
# You can send a point not on the curve; it may lie on a different curve
# If that curve's order is smooth, Pohlig-Hellman can be used

# Construction: choose a' such that y² = x³ + a'*x + b has smooth order
```

## ECDSA Nonce Reuse Attack

```python
"""
If the same nonce k is used for two ECDSA signatures:
s1 = k^(-1) * (h1 + r*d) mod n
s2 = k^(-1) * (h2 + r*d) mod n

s1 - s2 = k^(-1) * (h1 - h2) mod n
k = (h1 - h2) * (s1 - s2)^(-1) mod n
d = (s1 * k - h1) * r^(-1) mod n  (private key)
"""

def ecdsa_nonce_reuse(r1, s1, h1, r2, s2, h2, n):
    """Recover private key from ECDSA nonce reuse"""
    from gmpy2 import invert
    # Confirm r is identical
    assert r1 == r2
    k = ((h1 - h2) * invert(s1 - s2, n)) % n
    d = ((s1 * k - h1) * invert(r1, n)) % n
    return int(d)
```

## Common ECC CTF Challenge Types

| Challenge type | Feature | Attack |
|------|------|------|
| Standard curve + small order | n < 2^40 | Brute force |
| Standard curve + smooth order | n has small factors | Pohlig-Hellman |
| Anomalous curve | #E = p | Smart attack |
| Custom curve | suspicious a, b | Invalid Curve / factor the order |
| ECDSA signatures | Multiple signatures | Nonce reuse |
| Twisted Edwards | x² + a*y² = 1 + d*x²*y² | Convert to Weierstrass |

## References — lattice-and-lwe-attacks

# Lattice Attacks and LWE

## Basic Concepts

```
Lattice: a discrete additive subgroup of Z^n
Basis: a linearly independent set of vectors generating the lattice
LLL algorithm: finds an approximately shortest vector of a lattice basis (SVP approximation)
CVP (Closest Vector Problem): find the closest vector
SVP (Shortest Vector Problem): find the shortest vector
```

## LLL Algorithm

```python
# SageMath implementation
"""
A = matrix(ZZ, [[...], [...], ...])  # lattice basis matrix
B = A.LLL()  # LLL-reduced basis
# The column vectors of B are close to the shortest lattice vectors
```

## Hidden Number Problem (HNP)

```python
"""
Known: partial bits of (d_i, (t_i * a + k_i * d_i) mod p)
Recover: a (private key)
Use Coppersmith to solve for k_i
"""
# SageMath
def hnp_attack(d, t, bits, p):
    F.<x> = PolynomialRing(Zmod(p))
    # construct the polynomial...
```

## Coppersmith Related

```python
"""
Coppersmith finds small roots of polynomials:
f(x) = 0 mod n, |x| < n^(1/d)
where d is the polynomial degree
"""

# SageMath
def coppersmith_small_root(f, n, d, m):
    """f(x) = 0 mod n, find small root x, |x| < n^(1/(d*omega))"""
    # construct a lattice and run LLL
```

## LWE (Learning With Errors)

```python
"""
LWE problem:
Known: (A, b = As + e) mod q
Recover: s (private key)
where e is a small error vector

Common attacks:
1. Enumerate small errors (when e is very small)
2. BKW algorithm
3. Reduction to SVP/CVP
"""
```

## HNP Attack Template

```python
# SageMath: recover an RSA private key from partial key material
"""
DCP (Diffie-Hellman Claw Problem) variant
Solved using lattice reduction
"""

# Basic template
"""
F = GF(p)
P.<x> = PolynomialRing(F)

# Construct the lattice basis matrix
# Apply LLL
# Extract the private key from the reduced basis
"""
```

## Generic Lattice Attack Template

```python
# Consider a lattice attack when you encounter:
# 1. Multiple equations with unknowns and "small errors"
# 2. Partial private key / partial plaintext recovery
# 3. Reduction to a closest vector problem on a lattice

# Steps:
# 1. Model the problem as CVP/SVP on a lattice
# 2. Construct the lattice basis matrix
# 3. Reduce with LLL/BKZ
# 4. Extract the solution from the reduced basis
```

## References — prng-and-stream-cipher-attacks

# PRNG and Stream Cipher Attacks

## MT19937 ( Mersenne Twister ) Attack

```python
# MT19937 state recovery (given 624 outputs)
from ctypes import *

def untemper(y):
    y ^= y >> 18
    y ^= (y << 15) & 0xefc60000
    y ^= (y << 7) & 0x9d2c5680
    y ^= (y << 14) & 0x9d2c5680
    y ^= (y << 13) & 0x9d2c5680
    y ^= (y << 11) & 0x9d2c5680
    y ^= y >> 18
    return y

def recover_mt(outputs):
    """Recover internal state from 624 consecutive MT19937 outputs"""
    state = [untemper(y) for y in outputs[:624]]
    MT = c_ulong * 624
    mt = MT(*state)
    index = 624
    def twist():
        global index, mt
        for i in range(227):
            y = (mt[i] & 0x80000000) + (mt[(i+1)%624] & 0x7fffffff)
            mt[i] = mt[(i+397) % 624] ^ (y >> 1)
            if y & 1:
                mt[i] ^= 0x9908b0df
        index = 0
    return mt, twist, index
```

## LCG (Linear Congruential Generator) Attack

```python
"""
LCG: s_{n+1} = a * s_n + c (mod m)
With known parameters: iterate directly
With unknown parameters: 3 known (s, s_next) pairs yield a, c, m
"""

def lcg_attack(states):
    """Recover LCG parameters (a, c, m) from 3 consecutive states"""
    s0, s1, s2 = states[0], states[1], states[2]
    # s1 = a*s0 + c (mod m)
    # s2 = a*s1 + c (mod m)
    # s2 - s1 = a*(s1 - s0) (mod m)
    # Extended Euclidean algorithm to solve for a, m
```

## LFSR (Linear Feedback Shift Register) Attack

```python
"""
Berlekamp-Massey algorithm: recover the LFSR feedback polynomial from the output sequence
"""

def berlekamp_massey(s):
    """Recover the shortest LFSR feedback polynomial from a binary sequence"""
    # Sage implementation
    # F.<x> = GF(2)[]
    # s_seq = sequence(s)
    # return list(lfsr_sequence(f, [1]+[0]*15, len(s)))
```

## Known-Plaintext Attack (XOR Stream Cipher)

```python
"""
Stream cipher: C = P XOR keystream
If part of the plaintext P is known, recover keystream = C XOR P
The keystream can be used to decrypt other ciphertexts
"""

def xor_attack(ciphertext, known_plaintext):
    """Known-plaintext attack against an XOR stream cipher"""
    key = bytes(a ^ b for a, b in zip(ciphertext, known_plaintext))
    return key

def xor_decrypt(key, ciphertext):
    """Decrypt with the recovered keystream"""
    return bytes(a ^ b for a, b in zip(key, ciphertext))
```

## RC4 Attacks

```python
"""
Known RC4 weaknesses:
1. RC4 Drop (after discarding the first N bytes, the keystream is close to random)
2. Some key initializations are biased
"""

def rc4_drop(ciphertext, drop=3072):
    """Decrypt after RC4 Drop of N bytes"""
```

## Python random Module Prediction

```python
import random

# If you can access the Python random state, you can predict future random numbers
# The state is 624 * 4 = 2496 bytes
state = random.getstate()
# Advance the random numbers
random.setstate(state)
next_val = random.randint(0, 2**31)
```

## References — rsa-attacks-cheatsheet

# RSA Attack Cheat Sheet

## Attack Selection Decision Tree

```
Known n, e, c
├── e very small (e=3)?
│   ├── Same plaintext encrypted multiple times (multiple c)? → Håstad broadcast attack
│   └── Only one set? → Small-exponent root attack (low probability)
├── Multiple (n, e, c) sets?
│   ├── Same n? → Common modulus attack
│   ├── Same e? → Håstad broadcast attack
│   └── p or q share a common factor? → GCD factoring
├── e very large (>65537)?
│   └── d may be small → Wiener attack
├── n factorable?
│   ├── Fermat factoring (p≈q)
│   ├── Pollard p-1 (p-1 has small factors)
│   ├── Williams p+1 (p+1 has small factors)
│   └── Online lookup (factordb)
└── Partial information known?
    ├── Partial plaintext → Coppersmith
    ├── Partial p → Coppersmith
    └── Partial d → direct construction
```

## Small-Exponent Attack (e=3)

### Low-Exponent Broadcast Attack (Håstad)
```python
from gmpy2 import iroot
from functools import reduce

def hastard_broadcast(cs, ns, e=3):
    """When the same plaintext is encrypted under e different n values"""
    # CRT solving
    N = reduce(lambda a, b: a * b, ns)
    x = 0
    for i in range(e):
        Mi = N // ns[i]
        yi = pow(Mi, -1, ns[i])
        x += cs[i] * Mi * yi
    x %= N
    m = iroot(x, e)
    if m[1]:
        return int(m[0])
    return None
```

## Common Modulus Attack

```python
from gmpy2 import gcd

def common_modulus_attack(c1, c2, e1, e2, n):
    """Same plaintext, same n, different e"""
    g, s1, s2 = extended_gcd(e1, e2)
    if s1 < 0:
        c1 = pow(c1, -1, n)
        s1 = -s1
    if s2 < 0:
        c2 = pow(c2, -1, n)
        s2 = -s2
    m = (pow(c1, s1, n) * pow(c2, s2, n)) % n
    return m

def extended_gcd(a, b):
    if a == 0:
        return b, 0, 1
    g, x, y = extended_gcd(b % a, a)
    return g, y - (b // a) * x, x
```

## Wiener Attack (e large, d small)

```python
def wiener_attack(e, n):
    """Effective when d < n^(1/4)"""
    cf = continued_fraction(e, n)
    convergents = get_convergents(cf)
    for k, d in convergents:
        if k == 0:
            continue
        phi = (e * d - 1) // k
        # Check whether this is a valid phi
        x = n - phi + 1
        disc = x * x - 4 * n
        if disc >= 0:
            s = int(disc ** 0.5)
            if s * s == disc:
                return d
    return None
```

## Fermat Factoring (p ≈ q)

```python
from gmpy2 import is_square, iroot

def fermat_factor(n):
    """Effective when p and q are close together"""
    a = iroot(n, 2)[0] + 1
    b2 = a * a - n
    while not is_square(b2):
        a += 1
        b2 = a * a - n
    p = a + iroot(b2, 2)[0]
    q = a - iroot(b2, 2)[0]
    return int(p), int(q)
```

## Pollard p-1 Attack

```python
from math import gcd

def pollard_p1(n, B=100000):
    """Effective when all factors of p-1 are smaller than B"""
    a = 2
    for j in range(2, B):
        a = pow(a, j, n)
        d = gcd(a - 1, n)
        if 1 < d < n:
            return d, n // d
    return None
```

## Coppersmith Attack (Partial Plaintext Known)

```python
# Using SageMath
# When the high or low bits of the plaintext are known
# m = known_part + unknown_part
# unknown_part < n^(1/e)

# Sage implementation:
P.<x> = PolynomialRing(Zmod(n))
f = (known_prefix + x)^e - c
f = f.monic()
roots = f.small_roots()
if roots:
    m = known_prefix + roots[0]
```

## Online Factoring Tools

- https://factordb.com — look up already-factored n
- http://sagecell.sagemath.org — online Sage computation
