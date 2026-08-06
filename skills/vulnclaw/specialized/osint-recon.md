# stage: recon
# category: specialized

> OSINT  openSourceIntelligence Gatheringknowledge base — four维Information Gatheringmodule type（Server→networkstand→Domain Name→人member），dimensionfour（人memberinformation）conditiontrigger

# OSINT  openSourceIntelligence Gatheringknowledge base

针 for/toInformation Gathering/Reconnaissance/社will/can工程scenario's/ofsolid战knowledge base，provide**four维Information Gatheringmodule type**（Serverinformation → networkstandinformation → Domain Nameinformation → 人memberinformation），as well as具body's/oftoolusemethodanddataextracttip/trick。

**and/with `recon` Skill 's/ofdifferencepart**：
- `recon` → techniquelayer面Reconnaissance（PortScanning、DNS、DirectoryEnumeration）— foundation/basis版
- `osint-recon` → all/fulldimensionReconnaissance（Server + networkstand + Domain Name + 人member/社will/can工程）— deepdegree/measure版

## coreprinciple

1. **fourdimensionall/full覆stamp** — Server/networkstand/Domain Namethree (counter)dimension start endExecute，人memberdimensionpress/according toconditiontrigger
2. ** frompageextractone切canextract's/ofinformation** — not (classifier)look/see HTTP headers，stillneed tolook/see HTML content、JS File、comment
3. ** firstPassive after/backActive** —  firstlook/seeresponsehead/top、DNS、WHOIS（Passive）， again doPortScanning/DirectoryEnumeration（Active）
4. **dimension complete become/successdegree/measure自检** — everyroundInspect/Check哪些dimensioncompleted ✅，哪些incomplete ❌，all complete become/success after/backjustallows [DONE]
5. **Externallinki.e.线索** — pageascend's/ofeachExternallinkallcancanisinformationcomeSource
6. **structure-izeinputexit** — placehas/havediscover汇totalfor/is Markdown Report

## four维Information Gatheringmodule type

### dimensionone：Serverinformation
| Inspect/Checkitem | tool/method | explanation |
|--------|----------|------|
|  openrelease/putPort & Serviceversion | MCP nmap / `python_execute` + socket | all/fullPortScanningorcommonPort（21/22/80/443/3306/6379/8080/8443） |
| truesolid IP detect/probe | DNS historicalLog/Record / all/fullgame Ping / 邮 (classifier)head/topextract | CDN  after/back's/ofOrigin Server IP — SecurityTrails/DNSHistory/all/fullgamePing |
| Operating SystemFingerprint | TTL inference + nmap OS detection | Linux TTL≈64, Windows TTL≈128, Unix TTL≈255 |
| middle (classifier)version | responsehead/top Server + error/mistake页 + special征File | Apache/Nginx/IIS/Tomcat versionidentify |
| Databaseidentify | Portdetect/probe + error/mistakeinformation + special征rowfor/is | MySQL(3306)/Redis(6379)/MongoDB(27017)/MSSQL(1433) |

### dimensiontwo：networkstandinformation
| Inspect/Checkitem | tool/method | explanation |
|--------|----------|------|
| networkstand架construct | responsehead/top + pagespecial征 + JS Library | OS + middle (classifier) + Database +  language speech/language + Framework →  completewhole/integertechniqueStack |
| Web Fingerprint | `fetch` + responsespecial征Match | CMS type、FrontendFramework、JS Library、templatelead/guide擎 |
| WAF detection | wafw00f logic + responsespecial征 | Interceptpage/specialresponsehead/top/Exceptionstatecode |
| SensitiveDirectory & SensitiveFile | `python_execute` + commonPathDictionary | /admin /backup /config /api /robots.txt /sitemap.xml |
| SourcecodeLeak/Disclosure | Inspect/CheckcommonLeak/DisclosurePath | .git/.svn/.DS_Store/.env/web.config/backup file(.bak/.swp/.old) |
| 旁standquery | same/together IP negative/reverse查Domain Name | standgrowtool/micro步online/crt.sh same/together IP query |
| C  paragraph/segmentquery | same/togetherNetwork SegmentexistactivehostScanning | nmap -sn Scanning /24 Network Segment |

### dimensionthree：Domain Nameinformation
| Inspect/Checkitem | tool/method | explanation |
|--------|----------|------|
| WHOIS registerinformation | `python_execute` + whois API/command | register人/register商/NS Server/registerdate/ to期date |
| ICP backup案information | 工message partbackup案query API | onlymiddle/center国large陆Domain Name需查，境outDomain Nameno/withoutbackup案 |
| 子Domain Namediscover | crt.sh + brute force + Searchlead/guide擎 + DNS Zone Transfer | multi/multiplemethod交叉Validate，Ensure覆stampall/full面 |
| DNS Log/Recordall/fullquantity/measure | `python_execute` + dnspython/socket | A/CNAME/MX/TXT/NS/SPF/SOA all/fullquantity/measurequery |
| CertificatetransparencyLog | crt.sh / Censys / certspotter | discoverhistoricalCertificate、子Domain Name、associate/relatedDomain Name |

### dimensionfour：人memberinformation ⚡ conditiontrigger
**⚠️ thisdimensiononlyat/in with/bydescendcondition之onefull足timejustExecute：**
- usercommandmiddle/centerbrightcertainmentions"社will/can工程/社工/人memberinformation/ as/do者trace/人物画像"etc.
- goal/targetnetworkstandhas/havebrightcertain as/do者information（meta author、about page、联 system/relationshipway/manner）

**notshould do社工's/of情况**：regular/normal企业官networkno/without (counter)人 as/do者 / user (classifier)need to求"Scanninggoal/target" / goal/targetis IP/intranet/internal networkaddress

| tracedirection | method | explanation |
|----------|------|------|
|  as/do者identifierextract | page meta author、about page | user name、昵 call、邮箱 |
| GitHub trace | `fetch` + GitHub API | 仓Library、 language speech/language偏good、贡献Log/Record、邮箱 |
| Social Media |  frompageextractlink → Access | Bstand、micro博、know乎、Twitter、LinkedIn |
| 跨platformassociate/related |  useuser name/邮箱Searchotherplatform | same/identical ID 跨platformSearch |
| historicalCommit | GitHub commits → Commit邮箱 | associate/relatedotheritemeye/lookandIdentity |
| Leak/Disclosuredetection | GitHub historicalcodeSearch | .env、config、KeyLeak/Disclosure |

## First-Pass Workflow

1. **Accessgoal/target** → `fetch` Getfirst/head页，extract HTTP headers + HTML content
2. **dimensionone：Serverinformation** → PortScanning、truesolid IP、OS Fingerprint、middle (classifier)/Databaseidentify
3. **dimensiontwo：networkstandinformation** → Web Fingerprint、WAF detection、SensitiveDirectory/SourcecodeLeak/Disclosure、旁stand/C paragraph/segment
4. **dimensionthree：Domain Nameinformation** → WHOIS、ICP backup案、子Domain Name、DNS Log/Record、Certificatetransparency
5. **dimensionfour（conditiontrigger）** → extract as/do者information、跨platformtrace、information汇total
6. **dimension complete become/successdegree/measure自检** → Acknowledgmenteachdimensionarrivedecrease do past/excessiveoneroundInspect/Check
7. **汇totalReport** → generate Markdown formatReconnaissanceReport

## scenarioRoute

| scenario | reference document | corecontent |
|------|---------|---------|
| ServerInformation Gathering | `server-recon.md` | PortScanning、truesolid IP、OS Fingerprint、middle (classifier)/Databaseidentify |
| networkstandInformation Gathering | `website-recon.md` | 架construct/Fingerprint/WAF/SensitiveDirectory/SourcecodeLeak/Disclosure/旁stand/C paragraph/segment |
| Web Fingerprintidentify | `web-fingerprinting.md` | Frameworkdetection、versionidentify、techniqueStackinference |
|  as/do者tracemethod | `author-tracking.md` |  frompageextract as/do者 → 跨platformtrace → information汇total |
| OSINT tooluse | `osint-toolkit.md` | crt.sh、GitHub API、Searchlead/guide擎 dork、旁stand/C paragraph/segment/ICP |
| 社will/can工程information汇total | `social-engineering-intel.md` | 人物画像、 close/shut system/relationshipnetwork、information交叉Validate |
| ReconnaissanceReporttemplate | `recon-report-template.md` | standard Markdown Reportformat（fourdimension） |

## ⭐ often useextractcode (classifier) paragraph/segment

###  from HTML extractplacehas/haveExternallink
```python
import re
html = "..."  # fetch Get's/of HTML
links = re.findall(r'href=["\'](https?://[^"\']+)["\']', html)
for link in set(links):
    print(link)
```

###  from HTML extract as/do者information
```python
import re
# meta author
author = re.findall(r'<meta\s+name=["\']author["\']\s+content=["\']([^"\']+)["\']', html)
# about pagelink
about_links = re.findall(r'href=["\']([^"\']*(?:about|me|contact)[^"\']*)["\']', html, re.I)
```

### query crt.sh 子Domain Name
```python
import requests
domain = "example.com"
r = requests.get(f"https://crt.sh/?q=%.{domain}&output=json")
if r.status_code == 200:
    for entry in r.json():
        print(entry['name_value'])
```

### GitHub userinformation
```python
import requests
username = "target_user"
r = requests.get(f"https://api.github.com/users/{username}")
if r.status_code == 200:
    data = r.json()
    print(f"Name: {data.get('name')}")
    print(f"Bio: {data.get('bio')}")
    print(f"Email: {data.get('email')}")
    print(f"Blog: {data.get('blog')}")
    print(f"Location: {data.get('location')}")
    print(f"Company: {data.get('company')}")
```

### WAF detection（responsespecial征 method/law）
```python
import requests
url = "https://target.com"
# normalrequest
r1 = requests.get(url)
# trigger WAF 's/ofrequest（bring/carryattackspecial征）
r2 = requests.get(url + "/?id=1' OR 1=1--")
# comparisonresponse
if r1.status_code != r2.status_code or len(r1.text) != len(r2.text):
    print("[!] cancanexistat/in WAF")
    print(f"normalstatecode: {r1.status_code}, attackstatecode: {r2.status_code}")
```

### 旁standquery（same/together IP negative/reverse查Domain Name）
```python
import requests
ip = "1.2.3.4"
# use chinaz API orothernegative/reverse查interface
# alsocanvia/through crt.sh querysame/together IP 's/ofCertificate
r = requests.get(f"https://crt.sh/?q={ip}&output=json")
```

## References — author-tracking

#  as/do者tracemethod

## coreprocess

```
pageextract as/do者identifier → determines唯oneidentifiersymbol/character(user name/邮箱) → 跨platformSearch → information汇total
```

## Step 1:  frompageextract as/do者identifier

### HTML Meta tag/label
```python
import re

def extract_author_from_meta(html):
    """ from HTML meta tag/labelextract as/do者information"""
    authors = []
    
    # <meta name="author" content="XXX">
    m = re.findall(r'<meta\s+name=["\']author["\']\s+content=["\']([^"\']+)["\']', html)
    authors.extend(m)
    
    # <meta name="copyright" content="XXX">
    m = re.findall(r'<meta\s+name=["\']copyright["\']\s+content=["\']([^"\']+)["\']', html)
    authors.extend(m)
    
    # OG tag/label
    m = re.findall(r'<meta\s+property=["\']article:author["\']\s+content=["\']([^"\']+)["\']', html)
    authors.extend(m)
    
    return list(set(authors))
```

### pagelinkextract
```python
def extract_social_links(html):
    """ frompageextractSocial Medialink"""
    links = re.findall(r'href=["\'](https?://[^"\']+)["\']', html)
    
    social = {}
    for link in links:
        if 'github.com' in link:
            social['github'] = link
        elif 'bilibili.com' in link:
            social['bilibili'] = link
        elif 'weibo.com' in link or 'weibo.cn' in link:
            social['weibo'] = link
        elif 'zhihu.com' in link:
            social['zhihu'] = link
        elif 'twitter.com' in link or 'x.com' in link:
            social['twitter'] = link
        elif 'linkedin.com' in link:
            social['linkedin'] = link
        elif 'youtube.com' in link:
            social['youtube'] = link
        elif 'facebook.com' in link:
            social['facebook'] = link
    
    return social
```

## Step 2: GitHub trace

### userinformation API
```python
import requests

def get_github_profile(username):
    """Get GitHub userPublicinformation"""
    r = requests.get(f"https://api.github.com/users/{username}")
    if r.status_code != 200:
        return None
    
    data = r.json()
    return {
        'name': data.get('name'),
        'bio': data.get('bio'),
        'email': data.get('email'),
        'blog': data.get('blog'),
        'location': data.get('location'),
        'company': data.get('company'),
        'public_repos': data.get('public_repos'),
        'followers': data.get('followers'),
        'following': data.get('following'),
        'created_at': data.get('created_at'),
        'avatar_url': data.get('avatar_url'),
    }

def get_github_repos(username):
    """GetuserPublic仓Library（inferencetechniqueStack）"""
    r = requests.get(f"https://api.github.com/users/{username}/repos?per_page=100")
    if r.status_code != 200:
        return []
    
    repos = r.json()
    languages = {}
    for repo in repos:
        lang = repo.get('language')
        if lang:
            languages[lang] = languages.get(lang, 0) + 1
    
    return {
        'top_languages': sorted(languages.items(), key=lambda x: -x[1])[:5],
        'repo_count': len(repos),
        'starred_total': sum(r.get('stargazers_count', 0) for r in repos),
    }
```

###  from GitHub CommitLog/Recordextract邮箱
```python
def get_github_commit_email(username, repo):
    """ from GitHub CommitLog/Recordextract as/do者邮箱"""
    r = requests.get(f"https://api.github.com/repos/{username}/{repo}/commits?per_page=10")
    if r.status_code != 200:
        return []
    
    emails = set()
    for commit in r.json():
        author = commit.get('commit', {}).get('author', {})
        if author.get('email'):
            emails.add(author['email'])
    
    return list(emails)
```

## Step 3: 跨platformassociate/related

###  useuser nameSearchotherplatform
```python
# commonplatformdetection
PLATFORMS = {
    'GitHub': 'https://github.com/{username}',
    'Bstand': 'https://space.bilibili.com/search?keyword={username}',
    'know乎': 'https://www.zhihu.com/search?type=content&q={username}',
    'CSDN': 'https://blog.csdn.net/{username}',
    '掘金': 'https://juejin.cn/user/{username}',
    'Twitter': 'https://twitter.com/{username}',
    'LinkedIn': 'https://www.linkedin.com/in/{username}',
}

async def cross_platform_search(username, fetch_tool):
    """ useuser nameat/inmulti/multiple (counter)platformSearch"""
    results = {}
    for platform, url_template in PLATFORMS.items():
        url = url_template.format(username=username)
        try:
            resp = await fetch_tool(url=url)
            if resp.get('status') == 200:
                results[platform] = f"✅ 找 to ({url})"
            else:
                results[platform] = f"❌ un-找 to"
        except:
            results[platform] = f"⚠️ detectionfailure"
    return results
```

## Step 4: information汇totaltemplate

```markdown
## 人物画像：{昵 call}

### foundation/basisinformation
- **昵 call**：xxx
- **truesolid姓 name**：xxx（like/such ashas/have）
- **邮箱**：xxx
- **location**：xxx
- **职业/公司**：xxx

### technique画像
- **main力 language speech/language**：Python / JavaScript / ...
- **techniqueStack偏good**：...
- ** openSource贡献**：N  (counter)仓Library，M 颗星
- **感prosper趣leaddomain**：...

### Social Media
- GitHub: xxx
- Bstand: xxx
- know乎: xxx
- ...

### associate/relatedinformation
- 跨platformsame/identical ID：xxx
- Knownitemeye/look：xxx
- historicalLeak/Disclosure：xxx
```

## References — osint-toolkit

# OSINT toolusemanual

## 1. crt.sh — Certificatetransparency子Domain Namequery

###  use method/law
```python
import requests

def query_crtsh(domain):
    """via/through crt.sh query子Domain Name"""
    url = f"https://crt.sh/?q=%25.{domain}&output=json"
    try:
        r = requests.get(url, timeout=30)
        if r.status_code == 200:
            data = r.json()
            subdomains = set()
            for entry in data:
                name = entry.get('name_value', '')
                for n in name.split('\n'):
                    n = n.strip().lower()
                    if n and '*' not in n:
                        subdomains.add(n)
            return sorted(subdomains)
    except Exception as e:
        return [f"queryfailure: {e}"]
    return []
```

### Note
- crt.sh cancan较slow，setting 30s Timeout
- result/outcomeincludes/containscommonmatchsymbol/characterCertificate（`*.example.com`），需Filter
- go/leave re-/heavy after/backreturns

## 2. GitHub API — codeand/withuserSearch

### Searchcode（detectionLeak/Disclosure）
```python
def search_github_code(query, max_results=10):
    """Search GitHub code（detectionKey/configurationLeak/Disclosure）"""
    url = "https://api.github.com/search/code"
    params = {'q': query, 'per_page': max_results}
    headers = {'Accept': 'application/vnd.github.v3+json'}
    
    r = requests.get(url, params=params, headers=headers)
    if r.status_code == 200:
        items = r.json().get('items', [])
        return [{
            'repo': item['repository']['full_name'],
            'path': item['path'],
            'url': item['html_url'],
        } for item in items]
    return []
```

### often useSearch dork
```
"domain.com" password
"domain.com" api_key
"domain.com" secret
"domain.com" .env
filename:.env domain.com
filename:config domain.com
org:company-name password
```

## 3. DNS query

### Python inner/insideplace DNS query
```python
import socket

def dns_lookup(domain):
    """foundation/basis DNS query"""
    results = {}
    try:
        # A Log/Record
        results['A'] = socket.gethostbyname_ex(domain)[2]
    except:
        results['A'] = 'parsingfailure'
    
    return results
```

###  completewhole/integer DNS query（need dnspython）
```python
# like/such as resultenvironmenthas/have dnspython
try:
    import dns.resolver
    
    def full_dns_lookup(domain):
        record_types = ['A', 'AAAA', 'CNAME', 'MX', 'TXT', 'NS']
        results = {}
        for rtype in record_types:
            try:
                answers = dns.resolver.resolve(domain, rtype)
                results[rtype] = [str(r) for r in answers]
            except:
                pass
        return results
except ImportError:
    pass
```

## 4. WHOIS query

### online WHOIS API
```python
def whois_lookup(domain):
    """via/throughonline API query WHOIS"""
    # use whoisjson.com 免费 API
    url = f"https://whoisjson.com/api/v1/whois?domain={domain}"
    try:
        r = requests.get(url, timeout=10)
        if r.status_code == 200:
            data = r.json()
            return {
                'registrar': data.get('registrar'),
                'creation_date': data.get('creation_date'),
                'expiration_date': data.get('expiration_date'),
                'name_servers': data.get('name_servers'),
                'registrant': data.get('registrant'),
            }
    except:
        pass
    return {}
```

## 5. Google Dorking

### often useSearch language method/law
|  language method/law |  use途 | example |
|------|------|------|
| `site:` | limitdefineDomain Name | `site:github.com "unclec"` |
| `intitle:` | title close/shutkey word | `intitle:"index of" site:example.com` |
| `inurl:` | URL  close/shutkey word | `inurl:admin site:example.com` |
| `filetype:` | Filetype | `filetype:pdf site:example.com` |
| `"exact phrase"` | exactMatch | `"UncleCheng" security` |
| `related:` | 相 close/shutnetworkstand | `related:github.com` |

### Information Gatheringoften use dork
```
site:github.com "goal/targetuser name"
site:bilibili.com "goal/targetuser name"
site:zhihu.com "goal/targetuser name"
"邮箱@domain.com"
"手machinenumber"
```

## 6. Shodan/Censys（需 API Key）

### Shodan Search
```python
def shodan_search(api_key, query):
    import shodan
    api = shodan.Shodan(api_key)
    try:
        results = api.search(query)
        return [{
            'ip': result['ip_str'],
            'port': result['port'],
            'org': result.get('org', ''),
            'data': result['data'][:200],
        } for result in results['matches'][:10]]
    except Exception as e:
        return [f"Shodan queryfailure: {e}"]
```

## 7. Wayback Machine

### queryhistoricalSnapshot
```python
def wayback_query(domain):
    """query Wayback Machine historicalSnapshot"""
    url = f"http://archive.org/wayback/available?url={domain}"
    try:
        r = requests.get(url, timeout=10)
        if r.status_code == 200:
            data = r.json()
            snapshots = data.get('archived_snapshots', {})
            if snapshots.get('closest'):
                return snapshots['closest']['url']
    except:
        pass
    return None
```

## 8. 旁standquery（same/together IP negative/reverse查Domain Name）

### onlinetool
| tool | URL | explanation |
|------|-----|------|
| standgrowtool | https://stool.chinaz.com/same | 国inner/insidemostoften use |
| micro步online | https://x.threatbook.cn | Threat Intelligence+旁stand |
| crt.sh | https://crt.sh |  use IP 查Certificateassociate/relatedDomain Name |
| Censys | https://search.censys.io | all/full球AssetSearch |
| Fofa | https://fofa.info | emptybetweenSearchlead/guide擎 |

### python_execute 旁standquery
```python
import requests

def reverse_ip_lookup(ip):
    """via/through crt.sh negative/reverse查same/together IP Domain Name"""
    domains = set()
    try:
        r = requests.get(f"https://crt.sh/?q={ip}&output=json", timeout=30)
        if r.status_code == 200:
            for entry in r.json():
                for name in entry.get('name_value', '').split('\n'):
                    name = name.strip()
                    if name and '*' not in name:
                        domains.add(name)
    except Exception as e:
        print(f"crt.sh queryfailure: {e}")
    return sorted(domains)

# use
ip = "1.2.3.4"
result = reverse_ip_lookup(ip)
print(f"[+] same/together IP Domain Name ({len(result)}):")
for d in result:
    print(f"  - {d}")
```

## 9. C  paragraph/segmentquery（same/togetherNetwork Segmentexistactivehost）

### onlinetool
| tool | URL | explanation |
|------|-----|------|
| Fofa | https://fofa.info | `ip="1.2.3.0/24"` |
| Shodan | https://www.shodan.io | `net:1.2.3.0/24` |
| Censys | https://search.censys.io | `ip:/1.2.3.0-1.2.3.255/` |

### python_execute C  paragraph/segmentScanning
```python
import socket
from concurrent.futures import ThreadPoolExecutor, as_completed

def scan_c_segment(ip, timeout=1, max_workers=100):
    """Scanning C  paragraph/segmentexistactivehost"""
    prefix = ".".join(ip.split(".")[:3])
    alive = []

    def check(host_ip):
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.settimeout(timeout)
            result = s.connect_ex((host_ip, 80))
            s.close()
            if result == 0:
                return host_ip
        except:
            pass
        return None

    targets = [f"{prefix}.{i}" for i in range(1, 255)]
    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = {executor.submit(check, t): t for t in targets}
        for future in as_completed(futures):
            result = future.result()
            if result:
                alive.append(result)

    return sorted(alive, key=lambda x: int(x.split(".")[-1]))

# use
ip = "1.2.3.4"
hosts = scan_c_segment(ip)
print(f"[+] C  paragraph/segmentexistactivehost ({len(hosts)}):")
for h in hosts:
    print(f"  - {h}")
```

## 10. ICP backup案query

### onlinetool
| tool | URL | explanation |
|------|-----|------|
| 工message partbackup案query | https://beian.miit.gov.cn | 官directionright威 |
| standgrowtoolbackup案query | https://icp.chinaz.com | then捷query |
| 天眼查 | https://www.tianyancha.com | 企业+backup案associate/related |
| 爱standbackup案query | https://www.aizhan.com/cha/ | Batchquery |

### python_execute ICP backup案query
```python
import requests

def icp_lookup(domain):
    """query ICP backup案information（usePublic API）"""
    # method1: use chinaz API（need API key）
    # method2: usePublicqueryinterface
    try:
        # use whois querymiddle/center国Domain Nameinformation
        url = f"https://whois.chinaz.com/{domain}"
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
        }
        r = requests.get(url, headers=headers, timeout=10)
        # parsingbackup案information
        import re
        icp_match = re.search(r'backup案number[：:]\s*([^<\s]+)', r.text)
        if icp_match:
            return icp_match.group(1)
    except:
        pass

    # like/such as resultis境outDomain Name，usuallyno/without ICP backup案
    return "un-查 tobackup案（cancanfor/is境outDomain Name）"
```

## 11. 子Domain Namediscover（multi/multiplemethod）

### methodcombinationstrategy
1. **crt.sh** — Certificatetransparency（mostfast）
2. **Searchlead/guide擎 dork** — Google/Bing site: Search
3. **DNS brute force** — common before/front缀Dictionary
4. **DNS Zone Transfer** — attempt axfr
5. **JS FileAnalysis** —  frompage JS middle/centerextract子Domain Name

### python_execute 子Domain Namebrute force
```python
import socket
from concurrent.futures import ThreadPoolExecutor

def subdomain_brute(domain, wordlist=None, max_workers=20):
    """子Domain Namebrute force"""
    if wordlist is None:
        wordlist = [
            'www', 'mail', 'ftp', 'admin', 'blog', 'dev', 'staging',
            'api', 'test', 'portal', 'cdn', 'ns1', 'ns2', 'mx',
            'app', 'web', 'git', 'ci', 'jenkins', 'jira',
            'vpn', 'remote', 'shop', 'store', 'news',
        ]

    found = []
    def check(sub):
        fqdn = f"{sub}.{domain}"
        try:
            ip = socket.gethostbyname(fqdn)
            return (fqdn, ip)
        except:
            return None

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        results = executor.map(check, wordlist)
        found = [r for r in results if r]

    return sorted(found, key=lambda x: x[0])

# use
domain = "example.com"
subs = subdomain_brute(domain)
print(f"[+] discover子Domain Name ({len(subs)}):")
for sub, ip in subs:
    print(f"  - {sub} → {ip}")
```

### DNS Zone Transferattempt
```python
import socket

def try_zone_transfer(domain):
    """attempt DNS Zone Transfer"""
    # Get NS Log/Record
    try:
        ns_servers = socket.getaddrinfo(domain, None)
    except:
        return []

    # attempt for/toeach NS ServeradvancerowZone Transfer
    # Note：presentgeneration/proxy DNS ServerusuallyalreadyDisablethismeritcan
    import subprocess
    results = []
    try:
        result = subprocess.run(
            ['dig', 'axfr', domain, '@' + domain],
            capture_output=True, text=True, timeout=10
        )
        if 'XFR size' in result.stdout:
            results.append(result.stdout)
    except:
        pass

    return results
```

## References — recon-report-template

# ReconnaissanceReporttemplate

## useexplanation

at/inInformation GatheringTask complete become/successtime，use `python_execute` toolwill/shall with/bydescendtemplatePaddingfor/is completewhole/integerReport，
save touser指definePathor桌面。

## Markdown Reporttemplate

```markdown
# 🦞 {goal/target} ReconnaissanceReport

> generatetime：{datetime}
> tool：VulnClaw v0.3.7

---

## 1. goal/target概view

| itemeye/look | content |
|------|------|
| goal/target URL | {url} |
| IP address | {ip} |
| Server | {server} |
| Framework/CMS | {framework} |
| CDN | {cdn} |
| SSL Certificate | {ssl_info} |

---

## 2. techniqueReconnaissance

### 2.1 HTTP responsehead/top
| responsehead/top | value | securityTip |
|--------|---|---------|
| Server | {value} | {Notepoint} |
| X-Powered-By | {value} | Leak/DisclosuretechniqueStack |
| ... | ... | ... |

### 2.2 DNS Log/Record
| type | value |
|------|---|
| A | {ip} |
| CNAME | {cname} |
| MX | {mx} |
| TXT | {txt} |

### 2.3 子Domain Name
| 子Domain Name | IP | explanation |
|--------|---|------|
| {sub} | {ip} | {note} |

### 2.4  openrelease/putPort
| Port | Service | version |
|------|------|------|
| 80 | HTTP | nginx/1.18 |
| 443 | HTTPS | nginx/1.18 |

### 2.5 Directoryand/withFile
| Path | statecode | explanation |
|------|--------|------|
| /robots.txt | 200 | {contentabstract} |
| /sitemap.xml | 200 | {contentabstract} |
| /.git/HEAD | 403/200 | {isno/notLeak/Disclosure} |

---

## 3. contentReconnaissance

### 3.1 pageMetadata
- **Title**：{title}
- **Description**：{desc}
- **Keywords**：{keywords}
- **Author**：{author}

### 3.2 Externallink
| link | type | explanation |
|------|------|------|
| {url} | GitHub |  (counter)人homepage |
| {url} | Bstand | videoemptybetween |
| {url} | CDN | resourceSourceLoad |

### 3.3 JavaScript File
| File |  close/shutkeydiscover |
|------|---------|
| {path} | {api_endpoint/config/key} |

### 3.4 hide/concealinformation
- HTML comment：{comments}
- hide/concealword paragraph/segment：{hidden_fields}
- 邮箱/联 system/relationshipway/manner：{contacts}

---

## 4. 人物trace

### 4.1  as/do者information
| itemeye/look | content | comeSource | placemessagedegree/measure |
|------|------|------|--------|
| 昵 call | {name} | {source} | 🟢/🟡/🔴 |
| GitHub | {url} | {source} | 🟢 |
| Bstand | {url} | {source} | 🟢 |
| 邮箱 | {email} | {source} | 🟡 |
| location | {location} | {source} | 🟡 |

### 4.2 technique画像
- **main力 language speech/language**：{languages}
- **techniqueStack**：{stack}
- ** openSourceitemeye/look**：{repos}
- ** close/shutnoteleaddomain**：{interests}

### 4.3 跨platformassociate/related
| platform | user name/ID | Matchdegree/measure | explanation |
|------|----------|--------|------|
| {platform} | {id} | high/middle/center/low | {note} |

---

## 5.  close/shutkeydiscover

| # | discover | Risk Level | explanation |
|---|------|---------|------|
| 1 | {finding} | 🔴high/🟡middle/center/🟢low | {detail} |

---

## 6. Recommendation

1. {suggestion_1}
2. {suggestion_2}

---

*thisReport by/from VulnClaw Automaticgenerate，placehas/haveinformationcomeSource at/inPublicchannel。*
```

## Python savecode

```python
import os
from datetime import datetime

def save_recon_report(target, report_content, output_path=None):
    """saveReconnaissanceReport toFile"""
    if not output_path:
        # defaultsave to桌面
        desktop = os.path.join(os.path.expanduser('~'), 'Desktop')
        safe_name = re.sub(r'[^\w]', '_', target)[:30]
        date_str = datetime.now().strftime('%Y%m%d_%H%M')
        output_path = os.path.join(desktop, f'{safe_name}_ReconnaissanceReport_{date_str}.md')
    
    os.makedirs(os.path.dirname(output_path) or '.', exist_ok=True)
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(report_content)
    
    return output_path
```

## References — server-recon

# ServerInformation Gatheringreference

## 1.  openrelease/putPort & Serviceversionidentify

### nmap often usecommand
```bash
# all/fullPortScanning（slow但all/full面）
nmap -p- -sV <target>

# commonPortfastspeed/fastScanning
nmap -sV -top-ports 1000 <target>

# UDP PortScanning
nmap -sU --top-ports 100 <target>

# Serviceversionidentify + OS detection
nmap -sV -O <target>
```

### python_execute way/manner（no/without nmap time）
```python
import socket

def scan_port(host, port, timeout=2):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(timeout)
        result = s.connect_ex((host, port))
        s.close()
        return result == 0
    except:
        return False

host = "target.com"
common_ports = [21,22,23,25,53,80,110,143,443,445,993,995,1433,1521,3306,3389,5432,6379,8080,8443,9200,27017]
open_ports = [p for p in common_ports if scan_port(host, p)]
print(f" openrelease/putPort: {open_ports}")
```

### Serviceversionidentify（Banner Grabbing）
```python
import socket

def grab_banner(host, port, timeout=3):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(timeout)
        s.connect((host, port))
        # HTTP ServiceSendrequestGet banner
        if port in [80, 443, 8080, 8443]:
            s.send(b"HEAD / HTTP/1.1\r\nHost: " + host.encode() + b"\r\n\r\n")
        else:
            s.send(b"\r\n")
        banner = s.recv(1024).decode('utf-8', errors='ignore')
        s.close()
        return banner[:200]
    except:
        return None
```

## 2. truesolid IP detect/probe（CDN  after/back's/ofOrigin Server IP）

### methodone：DNS historicalLog/Record
- SecurityTrails (https://securitytrails.com/dns-trials)
- DNSHistory (https://dnshistory.org)
- ViewDNS (https://viewdns.info/iphistory/)
- Netcraft Site Report (https://sitereport.netcraft.com/)

### methodtwo：all/fullgame Ping
```python
import requests
# usemulti/multiple (adverbial) Ping Service
urls = [
    f"https://www.whatsmydns.net/#A/{domain}",
    f"https://ping.pe/{domain}",
    f"https://tools.keycdn.com/curl?url={domain}",
]
# like/such as resultnotsame/together (adverbial)differenceparsing tonotsame/together IP，explanationuse(past tense) CDN
# like/such as resultmulti/multiple (adverbial)parsing tosame/togetherone IP，should/this IP cancanistruesolidOrigin Server
```

### methodthree：邮 (classifier)head/topextract
- register/logingoal/targetnetworkstand，collect邮 (classifier)
- view邮 (classifier)head/topmiddle/center's/of `Received:` word paragraph/segment
- cancanExpose邮 (classifier)Server's/oftruesolid IP

### methodfour：子Domain Nameparsing
- CDN usually (classifier)for/ismainDomain NameService
- 子Domain Name（like/such as mail.ftp.dev.staging）cancandirectreceive/connectparsing toOrigin Server IP
- Inspect/Checkplacehas/have子Domain Name's/of A Log/Record，excludes CDN IP

### methodfive：SSL CertificateSearch
```python
import requests
domain = "target.com"
r = requests.get(f"https://crt.sh/?q=%.{domain}&output=json")
if r.status_code == 200:
    # Findnotsame/together子Domain Name's/ofCertificateassociate/related's/of IP
    for entry in r.json():
        print(entry.get('name_value', ''))
```

## 3. Operating SystemFingerprint

### TTL inference
| TTL value | cancan's/ofOperating System |
|--------|-------------|
| ≈ 64 | Linux / Unix / macOS |
| ≈ 128 | Windows |
| ≈ 255 | networkset upbackup / 老 style/mode Unix |

```python
import subprocess
# Ping Get TTL
result = subprocess.run(['ping', '-c', '1', host], capture_output=True, text=True)
# Windows: ping -n 1 host
#  frominputexitmiddle/centerextract TTL
import re
ttl_match = re.search(r'TTL[=:]\s*(\d+)', result.output, re.I)
if ttl_match:
    ttl = int(ttl_match.group(1))
    if ttl <= 64:
        print("speculation: Linux/Unix")
    elif ttl <= 128:
        print("speculation: Windows")
    else:
        print("speculation: networkset upbackup")
```

### nmap OS detection
```bash
nmap -O <target>
# 更激advance（need root）
sudo nmap -O --osscan-guess <target>
```

## 4. middle (classifier)versionidentify

### HTTP responsehead/topAnalysis
```
Server: Apache/2.4.49 (Ubuntu)
Server: nginx/1.18.0
Server: Microsoft-IIS/10.0
X-Powered-By: PHP/7.4.3
X-Powered-By: Express
X-AspNet-Version: 4.0.30319
```

### error/mistakepagespecial征
- Apache: default 404 pagecontain/include "Apache" word样
- Nginx: default 404 pagecontain/include "nginx" word样
- IIS: defaulterror/mistake页contain/include IIS versioninformation
- Tomcat: default 404 pagecontain/include Apache Tomcat version

### special征Filedetect/probe
```python
import requests
target = "https://target.com"
# Apache
r = requests.get(f"{target}/server-status")  # 403 = existat/in
r = requests.get(f"{target}/server-info")    # 403 = existat/in
# Nginx
r = requests.get(f"{target}/nginx_status")   # cancanExposestate
# Tomcat
r = requests.get(f"{target}/manager/html")   # Administrative Interface
# IIS
r = requests.get(f"{target}/aspnet_client/") # ASP.NET special征
```

## 5. Databaseidentify

### Portdetect/probe
| Database | defaultPort | explanation |
|--------|---------|------|
| MySQL | 3306 | mostcommon |
| PostgreSQL | 5432 | common at/in Rails/Django |
| MSSQL | 1433 | Windows environment |
| MongoDB | 27017 | NoSQL |
| Redis | 6379 | cache/Message Queue |
| Oracle | 1521 | 企业level/grade |
| Memcached | 11211 | cache |

### error/mistakeinformationspecial征
- MySQL: `You have an error in your SQL syntax`
- PostgreSQL: `ERROR: syntax error at or near`
- MSSQL: `Microsoft SQL Server`
- Oracle: `ORA-01756`

### python_execute detection
```python
import socket

def check_db(host, port, timeout=2):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(timeout)
        s.connect((host, port))
        # attemptRead banner
        s.send(b"\r\n")
        banner = s.recv(1024)
        s.close()
        return banner.hex()[:40], banner[:100]
    except:
        return None, None

db_ports = {
    3306: "MySQL", 5432: "PostgreSQL", 1433: "MSSQL",
    27017: "MongoDB", 6379: "Redis", 1521: "Oracle",
}
for port, name in db_ports.items():
    hex_banner, banner = check_db(host, port)
    if hex_banner:
        print(f"[+] {name} ({port}): {banner}")
```

## References — social-engineering-intel

# 社will/can工程information汇total

## 人物画像buildFramework

### informationdimension

| dimension | dataSource | extractmethod |
|------|--------|---------|
| Identityidentifier | page meta、GitHub | correct/positive rule/principleextract author/copyright |
| Social Network | pageExternallink | `<a href>` MatchSocial MediaDomain Name |
| technique偏good | GitHub 仓Library language speech/language part/point布 | GitHub API |
|  (adverbial) principle/logiclocation | GitHub location、博客 |  (counter)人resource料页 |
| 职业information | GitHub company、LinkedIn |  (counter)人resource料页 |
| 联 system/relationshipway/manner | GitHub email、博客联 system/relationship页 | API + pageextract |
| prosper趣leaddomain | GitHub 仓Librarytheme/topic、博客文 chapter | 仓Library topics + 文 chapterclassification |

## information交叉Validate

### principle
1. **singleonecomeSourcenot采message** —  close/shutkeyinformationneedarrivedecrease 2  (counter)independentcomeSourceAcknowledgment
2. **time效property/natureannotate** — annotateinformation's/ofGettime， past/excessivetimeinformationsingle独mark
3. **placemessagedegree/measureassesslevel/grade**：
   - 🟢 **high**：multi/multiple (counter)independentcomeSourceAcknowledgment
   - 🟡 **middle/center**：singleonecan靠comeSource
   - 🔴 **low**：inference/un-Validate

### commonassociate/relatedpattern

```
博客 GitHub link → GitHub user name → GitHub API Get邮箱
                                  → GitHub API Get仓Library → techniqueStackinference
                                  → GitHub Commit邮箱 → associate/relatedotherIdentity

博客 Bstand link → Bstand UID → Bstandhomepage →  close/shutnote/粉丝 → prosper趣tag/label
                                    → 投稿video → techniqueleaddomain

user name → 跨platformSearch → discover更multi/multiple社交account
邮箱 → haveibeenpwned → dataLeak/DisclosureLog/Record
```

## Social Mediainformationextract

### Bstand
```python
import re

def extract_bilibili_uid(url):
    """ from Bstand URL extract UID"""
    # space.bilibili.com/12345
    m = re.search(r'bilibili\.com/(\d+)', url)
    if m:
        return m.group(1)
    return None
```

### micro博
```python
def extract_weibo_uid(url):
    """ frommicro博 URL extract UID"""
    # weibo.com/u/12345 or weibo.com/username
    m = re.search(r'weibo\.com/(?:u/)?(\w+)', url)
    if m:
        return m.group(1)
    return None
```

### know乎
```python
def extract_zhihu_username(url):
    """ fromknow乎 URL extractuser name"""
    # zhihu.com/people/username
    m = re.search(r'zhihu\.com/people/([^/?]+)', url)
    if m:
        return m.group(1)
    return None
```

## information汇totalReportformat

```markdown
# goal/targetReconnaissanceReport

## 📋 基thisinformation
| itemeye/look | content | placemessagedegree/measure | comeSource |
|------|------|--------|------|
| goal/target | https://xxx | - | userinputenter |
| Framework | Hexo | 🟢 | HTTPhead/top+HTMLspecial征 |
| Server | GitHub Pages | 🟢 | Serverhead/top |
|  as/do者 | XXX | 🟢 | meta author |
| ... | ... | ... | ... |

## 👤 人物画像
- **昵 call**：XXX
- **GitHub**：https://github.com/xxx
- **Bstand**：https://space.bilibili.com/xxx
- **techniqueStack**：Python / JavaScript
- **location**：deep圳
- ...

## 🔗 associate/relateddiscover
- [discover1]
- [discover2]

## 📌  close/shutkeydiscover
1. ...
2. ...

---
*Reportgeneratetime：YYYY-MM-DD HH:MM*
*datacomeSource：goal/targetnetworkstand、GitHub API、Social MediaPublicinformation*
```

## privacyand/with伦 principle/logic

- ✅  (classifier)gather**Publicinformation**（notneedlogini.e.canAccess's/ofcontent）
- ✅ notattemptloginother人account
- ✅ notexploitgather's/ofinformationadvancerow骚扰or社will/can工程attack
- ✅ annotateinformationcomeSource，Ensurecanchase溯
- ❌ notgather私人common讯content
- ❌ notexploitinformationadvancerow钓鱼orotherSpoofrowfor/is

## References — web-fingerprinting

# Web Fingerprintidentify

## Inspect/Checkitemclearsingle

### HTTP responsehead/topFingerprint
| responsehead/top | inferenceinformation | example |
|--------|---------|------|
| `Server` | Web Server | `nginx/1.18.0`、`Apache/2.4.41`、`GitHub.com` |
| `X-Powered-By` | Backend language speech/language/Framework | `PHP/7.4.3`、`Express`、`Next.js` |
| `X-AspNet-Version` | .NET version | `4.0.30319` |
| `Set-Cookie` | Frameworkspecial征 | `PHPSESSID`→PHP、`JSESSIONID`→Java、`csrf_token`→Django |
| `X-Generator` | CMS | `Hugo`、`WordPress`、`Ghost` |
| `X-DRupal-Cache` | CMS | Drupal |
| `Via` | Proxy/CDN | `1.1 varnish`→Varnish CDN |

### HTML SourcecodeFingerprint
```python
import re

# WordPress
wp_signs = ['wp-content', 'wp-includes', 'wordpress']
# Hexo
hexo_signs = ['hexo', 'hexo-theme']
# Hugo
hugo_signs = ['hugo', 'gohugo']
# Jekyll
jekyll_signs = ['jekyll']
# Next.js
next_signs = ['__NEXT_DATA__', '_next/']
# Vue
vue_signs = ['data-v-', '__vue__']
# React
react_signs = ['data-reactroot', '__react']

def detect_framework(html):
    html_lower = html.lower()
    frameworks = []
    checks = {
        'WordPress': wp_signs,
        'Hexo': hexo_signs,
        'Hugo': hugo_signs,
        'Jekyll': jekyll_signs,
        'Next.js': next_signs,
        'Vue': vue_signs,
        'React': react_signs,
    }
    for name, signs in checks.items():
        if any(s in html_lower for s in signs):
            frameworks.append(name)
    return frameworks
```

### JavaScript FileFingerprint
- Frameworkspecialhas/have JS FilePath：`/wp-includes/js/` → WordPress
- Vue/React DevTools detection：`__VUE_DEVTOOLS_GLOBAL_HOOK__`、`__REACT_DEVTOOLS_GLOBAL_HOOK__`
- Frameworkversionusuallyat/in JS commentorvariablemiddle/center

### CSS Fingerprint
- `/wp-content/themes/` → WordPress
- Hexo theme/topicspecial征 class  name
- Bootstrap/Tailwind class special征

### special征File
| FilePath | inferenceinformation |
|---------|---------|
| `/robots.txt` | CMS information、hide/concealPath |
| `/sitemap.xml` | standpointstructure |
| `/favicon.ico` | FrameworkdefaultGraph标 |
| `/.well-known/security.txt` | security联 system/relationshipway/manner |
| `/humans.txt` |  opensend/issue者information |
| `/.git/HEAD` | Git 仓LibraryLeak/Disclosure |
| `/.env` | environmentvariableLeak/Disclosure |

## GitHub Pages special征
- responsehead/top `Server: GitHub.com`
- `X-GitHub-Request-Id` existat/in
- `X-Cache: HIT` + `X-Fastly-Request-ID` → Fastly CDN
- `Via: 1.1 varnish` → Varnish cache
- commonFramework：Jekyll、Hexo、Hugo

---

## WAF detection

### common WAF identifyspecial征
| WAF | responsehead/top/pagespecial征 | Interceptstatecode |
|-----|----------------|-----------|
| Cloudflare | `Server: cloudflare`, `CF-Ray` | 403 |
| AWS WAF | `x-amz-request-id`, `x-amz-cf-id` | 403 |
| 阿in云 WAF | Cookie contain/include `acw_tc` | 405/403 |
| 腾讯云 WAF | specific JSON Interceptpage | 403 |
| 宝塔 WAF | Interceptpagecontain/include "宝塔" | 403 |
| security狗 | Interceptpagecontain/include "safedog" | 403/404 |
| ModSecurity | specific 403 + Server head/top | 403 |
| Nginx WAF | `HTTP/1.1 444` orspecial 403 | 444/403 |

### WAF detectionmethod
1. **normalrequest vs attackrequestcomparison** — Sendbring/carryattackspecial征's/ofrequest，observeresponsedifference
2. **responsehead/topInspect/Check** — certain/some WAF will/canAddspecificresponsehead/top
3. **Cookie Inspect/Check** — partial/some WAF settingtrace Cookie
4. **statecodeException** — attackrequestreturnsExceptionstatecode（403/406/429/444）

### common WAF bypasstrigger payload
```
/?id=1' OR 1=1--
/?search=<script>alert(1)</script>
/../../../etc/passwd
/?file=php://filter/convert.base64-encode/resource=index
```

---

## SourcecodeLeak/DisclosureInspect/Check

### commonSourcecodeLeak/Disclosuretypeand/withdetection
| type | Path | detectionmethod | harmgrade/level |
|------|------|---------|---------|
| Git 仓Library | `/.git/config`, `/.git/HEAD` | 200 且contain/include git content | 🔴 Critical |
| SVN 仓Library | `/.svn/entries` | 200 且contain/include svn content | 🔴 Critical |
| .DS_Store | `/.DS_Store` | Download after/backparsingDirectorystructure | 🟡 Medium |
| .env File | `/.env` | contain/include DB_PASSWORD etc. | 🔴 Critical |
| web.config | `/web.config` | IIS configurationLeak/Disclosure | 🟡 Medium |
| backup file | `/.bak`, `/.swp`, `/.old`, `/.tar.gz` | directreceive/connectDownload | 🟡 Medium |
| Docker | `/Dockerfile`, `/docker-compose.yml` | containerconfiguration | 🟡 Medium |
| package.json | `/package.json` | Node.js depend on | 🟢 Low |
| composer.json | `/composer.json` | PHP depend on | 🟢 Low |
| webpack | `/webpack.json`, `/map Files` | SourcecodeMap | 🟡 Medium |

### Git Leak/Disclosureexploitprocess
1. Access `/.git/HEAD` → Get ref Path
2. Access `/.git/config` → GetRemote仓Libraryinformation
3. Access `/.git/objects/` → traverse/iterate Git  for/to象
4. use GitHack/scrabble toolAutomaticrecoverySourcecode

### SensitiveFileScanningPathcolumntable
```
/.git/config
/.git/HEAD
/.svn/entries
/.DS_Store
/.env
/.env.bak
/.env.local
/web.config
/config.php
/config.yml
/backup.sql
/database.sql
/db.sql
/phpinfo.php
/test/
/debug/
/console/
/admin/
/wp-config.php
/robots.txt
/sitemap.xml
/.well-known/security.txt
```

## References — website-recon

# networkstandInformation Gatheringreference

## 1. networkstand架constructidentify

### techniqueStackinferencemethod
1. **HTTP responsehead/top** — Server、X-Powered-By、Set-Cookie special征
2. **HTML Sourcecodespecial征** — meta generator、specific class/id 命 name
3. **JS FilePath** — /static/js/app.js、/wp-content/、/assets/
4. **Cookie name** — PHPSESSID(php)、JSESSIONID(Java)、_rails_session(Rails)
5. **URL Path** — ?id= (PHP)、/api/ (REST)、/wp-admin/ (WordPress)

### common架constructcombination
|  language speech/language | Framework | Database | Server | special征 |
|------|------|--------|--------|------|
| PHP | Laravel | MySQL | Apache/Nginx | Set-Cookie: laravel_session |
| PHP | WordPress | MySQL | Apache | /wp-content/, /wp-admin/ |
| Python | Django | PostgreSQL | Nginx+Gunicorn | CSRF middleware cookie |
| Python | Flask | SQLite/MySQL | Nginx+uWSGI | Set-Cookie: session= |
| Java | Spring | MySQL/Oracle | Tomcat | JSESSIONID |
| Node.js | Express | MongoDB | Nginx | X-Powered-By: Express |
| Ruby | Rails | PostgreSQL | Nginx+Puma | _rails_session |

### python_execute 架constructdetect/probe
```python
import requests

url = "https://target.com"
r = requests.get(url, timeout=10)

# 1. responsehead/topAnalysis
headers = r.headers
print(f"Server: {headers.get('Server', 'N/A')}")
print(f"X-Powered-By: {headers.get('X-Powered-By', 'N/A')}")

# 2. Cookie Analysis
cookies = r.cookies
for cookie in cookies:
    print(f"Cookie: {cookie.name} = {cookie.value[:20]}...")

# 3. HTML special征Analysis
html = r.text
# WordPress
if 'wp-content' in html or 'wp-includes' in html:
    print("[+] WordPress detection")
# Laravel
if 'laravel_session' in str(cookies):
    print("[+] Laravel detection")
# Django
if 'csrftoken' in str(cookies) or 'csrfmiddlewaretoken' in html:
    print("[+] Django detection")
# Hexo
if 'hexo' in html.lower():
    print("[+] Hexo 博客detection")
# Hugo
if 'hugo' in html.lower():
    print("[+] Hugo 博客detection")
```

## 2. Web Fingerprintidentify

### CMS Fingerprintspecial征
| CMS | special征Path | special征string |
|-----|---------|-----------|
| WordPress | /wp-login.php, /wp-content/ | wp-content, xmlrpc.php |
| Joomla | /administrator/ | /media/jui/ |
| Drupal | /misc/drupal.js | Drupal.settings |
| Discuz | /forum.php | discuz_uid |
| Typecho | /admin/login.php | typecho |
| Hexo | /archives/ | hexo |
| Ghost | /ghost/ | ghost-frontend |

### FrontendFrameworkspecial征
| Framework | special征 |
|------|------|
| React | data-reactroot, __NEXT_DATA__ |
| Vue.js | data-v-xxx, __vue__ |
| Angular | ng-version, _nghost |
| jQuery | jQuery in scripts |
| Bootstrap | bootstrap.css/js |

### python_execute Fingerprintidentify
```python
import requests, re

url = "https://target.com"
r = requests.get(url, timeout=10)
html = r.text

# CMS detection
cms_signatures = {
    "WordPress": ["wp-content", "wp-includes", "wp-admin"],
    "Joomla": ["/administrator/", "media/jui"],
    "Drupal": ["Drupal.settings", "/misc/drupal"],
    "Hexo": ["hexo", "/archives/"],
    "Hugo": ["hugo", "gohugo"],
    "Ghost": ["ghost-frontend", "/ghost/"],
}

for cms, sigs in cms_signatures.items():
    if any(sig in html for sig in sigs):
        print(f"[+] CMS: {cms}")

# FrontendFrameworkdetection
fw_signatures = {
    "React": ["data-reactroot", "__NEXT_DATA__", "react"],
    "Vue.js": ["data-v-", "__vue__", "vue"],
    "Angular": ["ng-version", "_nghost", "angular"],
    "jQuery": ["jquery", "jQuery"],
    "Bootstrap": ["bootstrap"],
}

for fw, sigs in fw_signatures.items():
    if any(sig.lower() in html.lower() for sig in sigs):
        print(f"[+] FrontendFramework: {fw}")

# JS Fileextract
js_files = re.findall(r'src=["\']([^"\']*\.js[^"\']*)["\']', html)
print(f"JS File: {js_files[:10]}")
```

## 3. WAF detection

### common WAF special征
| WAF | Interceptspecial征 |
|-----|---------|
| Cloudflare | Server: cloudflare, CF-Ray header |
| AWS WAF | Server: AmazonS3, x-amz-request-id |
| 阿in云 WAF | Set-Cookie includes/contains acw_tc |
| 腾讯云 WAF | specificInterceptpage |
| 宝塔 WAF | Interceptpagecontain/include "宝塔" |
| security狗 | Interceptpagecontain/include "safedog" |
| ModSecurity | specific 403 response |

### python_execute WAF detection
```python
import requests

url = "https://target.com"

# 1. normalrequest
r1 = requests.get(url)

# 2. trigger WAF 's/ofrequest
waf_payloads = [
    "/?id=1' OR 1=1--",
    "/?search=<script>alert(1)</script>",
    "/../../../etc/passwd",
    "/?file=php://filter/convert.base64-encode/resource=index",
]

for payload in waf_payloads:
    r2 = requests.get(url + payload, allow_redirects=False)
    # statecodechange
    if r2.status_code in [403, 406, 429, 501]:
        print(f"[!] WAF detection: {payload} → {r2.status_code}")
    # responsegrowdegree/measureshow/display著change
    if abs(len(r2.text) - len(r1.text)) > 500:
        print(f"[!] responsegrowdegree/measurechange: normal={len(r1.text)}, attack={len(r2.text)}")

# 3. Inspect/Checkspecific WAF responsehead/top
waf_headers = {
    "cloudflare": ["cf-ray", "server: cloudflare"],
    "aws": ["x-amz-request-id", "x-amz-cf-id"],
    "阿in云": ["acw_tc"],
}
for waf_name, sigs in waf_headers.items():
    for sig in sigs:
        if sig in str(r1.headers).lower():
            print(f"[+] WAF detection: {waf_name}")
```

## 4. SensitiveDirectory & SensitiveFile

### commonSensitivePathcolumntable
```
/robots.txt
/sitemap.xml
/.git/
/.svn/
/.env
/.DS_Store
/web.config
/config.php
/config.yml
/backup/
/admin/
/login/
/api/
/swagger/
/graphql
/phpinfo.php
/test/
/debug/
/console/
/actuator/
/.well-known/
```

### python_execute DirectoryScanning
```python
import requests

target = "https://target.com"
paths = [
    "/robots.txt", "/sitemap.xml", "/.git/", "/.env", "/.DS_Store",
    "/admin/", "/backup/", "/config.php", "/api/", "/phpinfo.php",
    "/.git/config", "/.git/HEAD", "/wp-config.php",
    "/swagger/", "/graphql", "/actuator/",
]

for path in paths:
    try:
        r = requests.get(target + path, timeout=5, allow_redirects=False)
        if r.status_code in [200, 301, 302, 401, 403]:
            print(f"[{r.status_code}] {path}")
    except:
        pass
```

## 5. SourcecodeLeak/DisclosureInspect/Check

### commonSourcecodeLeak/Disclosuretype
| type | Path | detectionmethod |
|------|------|---------|
| Git 仓Library | /.git/config, /.git/HEAD | 200 且contain/include git content |
| SVN 仓Library | /.svn/entries | 200 且contain/include svn content |
| .DS_Store | /.DS_Store | Download after/backparsing |
| .env File | /.env | contain/include DB_PASSWORD etc. |
| web.config | /web.config | IIS configuration |
| backup file | /.bak, /.swp, /.old, /~ | directreceive/connectDownload |
| Docker | /Dockerfile, /docker-compose.yml | containerconfiguration |
| package.json | /package.json | Node.js depend on |
| composer.json | /composer.json | PHP depend on |

### Git 仓LibraryLeak/Disclosureexploit
```python
import requests

target = "https://target.com"

# 1. Inspect/Check .git/HEAD
r = requests.get(f"{target}/.git/HEAD")
if r.status_code == 200 and "ref:" in r.text:
    print("[!] Git 仓LibraryLeak/Disclosure!")
    # 2. attemptGet ref
    ref_path = r.text.strip().split("ref: ")[1] if "ref: " in r.text else ""
    if ref_path:
        r2 = requests.get(f"{target}/.git/{ref_path}")
        if r2.status_code == 200:
            print(f"[+] Git ref: {r2.text.strip()}")

# 3. attemptGet config
r3 = requests.get(f"{target}/.git/config")
if r3.status_code == 200:
    print(f"[+] Git config:\n{r3.text}")
```

## 6. 旁standquery（same/together IP negative/reverse查Domain Name）

### querymethod
1. **standgrowtool** — https://stool.chinaz.com/same
2. **micro步online** — https://x.threatbook.cn
3. **crt.sh** —  use IP queryCertificateassociate/relatedDomain Name
4. **Censys** — https://search.censys.io

### python_execute 旁standquery
```python
import requests, json

ip = "1.2.3.4"

# method1: crt.sh querysame/together IP Certificate
r = requests.get(f"https://crt.sh/?q={ip}&output=json", timeout=15)
if r.status_code == 200:
    domains = set()
    for entry in r.json():
        for name in entry.get('name_value', '').split('\n'):
            if name.strip() and '*' not in name:
                domains.add(name.strip())
    print(f"[+] same/together IP Domain Name ({len(domains)}):")
    for d in sorted(domains):
        print(f"  - {d}")
```

## 7. C  paragraph/segmentquery（same/togetherNetwork Segmentexistactivehost）

### python_execute C  paragraph/segmentScanning
```python
import requests, socket
from concurrent.futures import ThreadPoolExecutor

#  fromDomain NameGet IP
domain = "target.com"
ip = socket.gethostbyname(domain)
# extract C  paragraph/segment
c_segment = ".".join(ip.split(".")[:3])

def check_host(ip, timeout=1):
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(timeout)
        result = s.connect_ex((ip, 80))
        s.close()
        if result == 0:
            return ip
    except:
        pass
    return None

# Scanning C  paragraph/segment（1-254）
alive_hosts = []
with ThreadPoolExecutor(max_workers=50) as executor:
    ips = [f"{c_segment}.{i}" for i in range(1, 255)]
    results = executor.map(check_host, ips)
    alive_hosts = [ip for ip in results if ip]

print(f"[+] C  paragraph/segmentexistactivehost: {alive_hosts}")
```
