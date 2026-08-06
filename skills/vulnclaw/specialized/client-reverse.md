# stage: exploit
# category: specialized

> ClientReverseand/withBurpReplay —  repeatmixedClientSignaturerecovery、Encryptionrestoration、requestchaintrace、稳defineReplay，适used for/forauthorizedsafe卓ApppenetrationTest、BrowserJSSignature、桌面ClientReverse

# ClientReverseand/with Burp Replay Skill

whenrequest by/fromClient（safe卓App、BrowserJS、桌面Client）construct，且existat/inSignature、Encryption、tokenstate、set upbackupBindornegative/reverseAutomatic-izelogicleads to Burp cannotdirectreceive/connectReplaytime，usethis Skill。

## coreprinciple

**Packet-First**： firstCapture并Analysistruesolid's/of HTTP/HTTPS requestor WebSocket Streamquantity/measure，Acknowledgmentcan useproperty/nature， againpress/according to需Reverse阻塞point。Reverseis阻塞solvestep，notisdefaultenter口。

## scenarioRoute

### authorizedsafe卓 App penetrationTest

**do not first use jadx、ida_pro_mcp Analysis APK**，press/according to with/bydescendsequentialoperation：

1. Acknowledgmentgoal/target App alreadyInstallationat/inConnectionset upbackupascend
2. accuratebackupgood Burp or Charles grab/capturePackage
3.  use scrcpy_vision 打 open App，Drivertruesolid业务process
4. each close/shutkeyaction after/backInspect/Check Burp/Charles isno/notexitpresent HTTP/HTTPS or WebSocket dataPackage
5. like/such as resultPackagecansee/meet且canReplay → immediatelyenter `web-security-advanced`  do Web/API securityTest
6.  re-/heavy repeat"boundary面action → grab/capturePackage → Web securityAnalysis"循环
7.  (classifier)has/havegrab/capturenot toPackage/Packageby (passive)Encryption/cannotReplaytime → Upgrade to jadx → frida_mcp → ida_pro_mcp

**MCP toolchain**：scrcpy_vision → burp/charles → adb_mcp → jadx → frida_mcp → ida_pro_mcp

### Browser JS Signature、negative/reversecrawl、WebSocket grasp手

1. chrome_devtools viewpagestateandrequestchain
2. js_reverse definebit token/sign generatelogic
3. burp ValidateReplay并determinescan变word paragraph/segment

**phase/stagemodule type**：locate → recover → runtime → validation → replay

**MCP toolchain**：chrome_devtools → js_reverse → burp

### 桌面Client / Local signer

1. everything_search definebit相 close/shutFile
2. ida_pro_mcp staticAnalysisSignaturefunction
3. frida_mcp GetRuntimeparameter
4. burp Validate稳defineReplay

**MCP toolchain**：everything_search → ida_pro_mcp → frida_mcp → burp

## Replaythen绪Inspect/Checkclearsingle

at/inenter Payload Test before/front，mustcanreturnanswer：

- Request Bodylike/such as何construct？
- Signature/Encryptioninputenterfromwhere？
- 哪些 cookie、header、token、set upbackupvalue、time戳、nonce ismust's/of？
- requestisno/notdepend onsequentialorSessionstate？
- 哪些word paragraph/segment改move after/backwill not破badReplay？

## 证据keepstay/keep

- builder/signer/crypto codelocation
-  close/shutkey hook pointandRuntimeobservevalue
- can use's/of replay request样this
-  before/frontplacecondition、failurepatternandnegative/reverseAutomatic-izerowfor/isexplanation

## reference document

- `references/02-client-api-reverse-and-burp.md` — ClientReverse to Burp ReplaytotalWorkflow
- `references/android-authorized-app-pentest-sop.md` — safe卓 App penetration SOP
- `references/browser-js-signing-workflow.md` — Browser JS SignatureWorkflow
- `references/android-signing-and-crypto-workflow.md` — safe卓Signatureand/withEncryptionWorkflow
- `references/android-ui-driven-observation-and-packet-loop.md` — safe卓 UI Driverobserve循环
- `references/android-external-url-runtime-first-workflow.md` — safe卓External URL Test
- `references/android-network-layer-testing-quick-reference.md` — safe卓networklayerTestquick reference
- `references/MCP.md` — MCP can力totaldocument
- `references/tool-selection-map.md` — tool选 type (adverbial)Graph

## References — 02-client-api-reverse-and-burp

# 02 Client API Reverse And Burp Workflow

This integrated file merges the client-side reverse workflow, MCP execution order, tool-selection rules, and evidence expectations needed to move from opaque client traffic to reproducible Burp replay.

## Use This File When

- Burp cannot directly replay the request
- the client computes `sign`, `token`, `nonce`, `timestamp`, encrypted body, or device-bound fields
- the request sequence is stateful or tied to runtime values
- you need a reliable workflow from runtime packet capture to blocker-driven reverse recovery to replay

## Core Objective

Recover the full request-production recipe:

- where the request is assembled
- where crypto or signing is applied
- which runtime values are mandatory
- which state transitions must happen before replay

The first priority is not reverse for its own sake. The first priority is to capture the real HTTP/HTTPS request or WebSocket message that the client emits, confirm whether it is already usable, and only then reverse the missing blocker.

## Front Rule For Authorized Android App Pentest

When the task is to pentest an authorized Android app, do not start with `jadx`, `ida_pro_mcp`, or APK-first reverse analysis.
Start in this order instead:

1. confirm the target app is actually installed on the connected Android device
2. get `burp` or `charles` ready before driving the feature
3. use `scrcpy_vision` to open the app and drive real business features
4. after each important action, inspect `burp` or `charles` for HTTP/HTTPS requests or WebSocket messages
5. if packets are already visible and usable, move directly into `web-playbook-index.md` and test the server-side surface
6. repeat the UI action -> packet capture -> Web security analysis loop for the next business feature
7. escalate into `jadx`, `frida_mcp`, or `ida_pro_mcp` only when packets are absent, encrypted, opaque, still not replayable, or when runtime anomalies clearly point to a client-side blocker

For this Android pentest path, reverse engineering is a blocker-resolution step, not the default entrypoint.

## Recommended Read Path

1. Read `Goal`, `Stages`, and the platform-specific section for Android, desktop, or browser JS.
2. Read `Priority` and `Primary Chains` to choose the smallest MCP chain.
3. If the target is browser JS, continue into `browser-js-signing-workflow.md`.
4. If the target is Android external URL testing, continue into `android-external-url-runtime-first-workflow.md`.
5. If the target is Android reverse or crypto recovery, continue into `android-signing-and-crypto-workflow.md`.
6. If Android runtime progress depends on app UI state, continue into `android-ui-driven-observation-and-packet-loop.md`.
7. Read `Rule` and `reporting-and-evidence.md` content before switching to Burp.
8. After replay is stable, move into `web-playbook-index.md` or `04-ai-and-mcp-security-integrated.md`.

## Replay Readiness Checklist

- you can name the builder, signer, or serializer location
- you know which cookies, headers, tokens, timestamps, or device values are required
- you know whether request ordering matters
- you have one working replay outside the client
- you know which fields are safe to mutate during later testing

## Platform Branch Rules

### Browser JS

- decide the stage from engineering state, not from clue words alone
- stay in `locate` until the request, sink, and upstream dependency chain are real
- only enter `recover` after the boundary is proven
- only enter `runtime` when the boundary is clear but browser and local execution diverge
- only enter `validation` when the remaining work is checkpoint proof

Detailed branch file: `references/browser-js-signing-workflow.md`
Stage references: `references/browser-locate-and-request-chain.md`, `references/browser-recover-and-shell-reduction.md`, `references/browser-runtime-fit-and-risk.md`, `references/browser-validation-and-handoff.md`
Record template: `references/browser-request-chain-template.md`

### Android

- for external URL testing, start with live app interaction and packet visibility, not reverse engineering
- first confirm the target app is installed on a connected device and can actually be launched
- use `scrcpy_vision` to navigate, inspect screenshots, and decide the next action
- check `burp` or `charles` for HTTP/HTTPS requests or WebSocket messages after each important action
- use `adb_mcp` to review logs after important actions
- once packets are visible and replayable, move directly into `web-playbook-index.md` and keep the UI-action to packet to Web-analysis loop going for the next business feature
- reverse Java only when packets are absent, encrypted, still opaque, or otherwise blocked
- escalate into JNI or `.so` work only when Java stops exposing the required inputs or outputs
- use `frida_mcp` when hook-based plaintext recovery is faster than reimplementation

Detailed branch files: `references/android-external-url-runtime-first-workflow.md`, `references/android-signing-and-crypto-workflow.md`
Phase references: `references/android-static-triage-and-callflow.md`, `references/android-dynamic-hooking-and-replay.md`, `references/android-ui-driven-observation-and-packet-loop.md`, `references/android-native-signature-analysis.md`
Record template: `references/android-signature-reverse-template.md`

## Included Sources

- references/client-reverse-workflow.md
- references/mcp-first-methodology.md
- references/tool-selection-map.md
- references/reporting-and-evidence.md

---

## Source: client-reverse-workflow.md

Path: references/client-reverse-workflow.md

# Complex Client Reverse Workflow

## Goal

Recover the real request-production chain so the interface can be reproduced outside the client.

## Stages

1. classify the client
2. choose the smallest platform branch that can prove the request chain
3. for Android app pentests, confirm app presence on the device and try runtime packet capture before any reverse step
4. dynamically confirm signer, serializer, and state values only when runtime packet proof is no longer enough
5. statically recover the missing blocker only after runtime visibility, plaintext, or replay stalls
6. rebuild the request recipe
7. replay in Burp
8. move into Web or AI attack testing only after replay is stable

## Android

- start by confirming the target app exists on the connected device, then use `scrcpy_vision`, logs, and proxy visibility for external URL testing
- move to `jadx` only when packets are missing, encrypted, or blocked
- reverse Java before native
- use `frida_mcp` when runtime hook proof or plaintext recovery is faster than deeper reverse
- dump and analyze `.so` only after Java has stopped answering the blocker
- move to `burp`, then into Web security analysis once replay is stable

## Native desktop

- locate files with `everything_search`
- reverse code with `ida_pro_mcp`
- capture runtime values with `frida_mcp`
- move to `burp`

## Browser JS

- inspect live requests with `chrome_devtools`
- choose the current stage from `locate`, `recover`, `runtime`, or `validation`
- trace initiators and signer functions with `js_reverse`
- replay with `burp`

## Android sign and crypto

- enter this branch only after runtime-first packet checks prove reverse is required, or when the task is already an explicit Android sign or crypto reverse problem
- decompile and triage in `jadx`
- trace request flow from manifest and entry components
- locate request builder, interceptor, signer, encryptor, and JNI handoff
- confirm final on-wire values with `frida_mcp` or `charles` only after static triage narrows the target
- replay with `burp`

## Android external URL runtime-first

- drive the app with `scrcpy_vision`
- inspect screenshots for visible anomalies and state changes
- review logs with `adb_mcp`
- verify whether `burp` or `charles` sees traffic
- only then decide whether Java reverse, Frida hooks, or dumped `.so` analysis is necessary

## Detailed Branches

- browser JS staged flow: `browser-js-signing-workflow.md`
- Android external URL runtime-first flow: `android-external-url-runtime-first-workflow.md`
- Android sign and crypto flow: `android-signing-and-crypto-workflow.md`

For staged browser work, continue into `references/browser-js-signing-workflow.md`.
For Android external URL testing, continue into `references/android-external-url-runtime-first-workflow.md`.
For Android blocker recovery or explicit sign and crypto reverse work, continue into `references/android-signing-and-crypto-workflow.md`.


---

## Source: mcp-first-methodology.md

Path: references/mcp-first-methodology.md

# MCP-First Methodology

This file is a navigation aid. The full methodology lives in `references/methodology/MCP.md`.

## Priority

1. Read the raw `MCP.md`
2. Select the minimal MCP chain for the target
3. Capture the real HTTP/HTTPS request or WebSocket message before deeper reverse
4. Restore the request lifecycle before Burp replay

## Primary Chains

### Android

- `scrcpy_vision`
- `burp`
- `charles`
- `adb_mcp`
- `jadx` only when packets are blocked
- `frida_mcp`
- `ida_pro_mcp`

### Native or packed desktop

- `everything_search`
- `ida_pro_mcp`
- `frida_mcp`
- `burp`

### Browser JS reverse

- `chrome_devtools`
- `js_reverse`
- `burp`


---

## Source: tool-selection-map.md

Path: references/tool-selection-map.md

# Tool Selection Map

## Reverse Layer

- `jadx`
- `ida_pro_mcp`
- `frida_mcp`
- `scrcpy_vision`
- `adb_mcp`
- `charles`
- `js_reverse`
- `chrome_devtools`
- `burp`

## Support Layer

- `everything_search`
- `context7`
- `fetch`
- `memory`
- `sequential_thinking`

## Platform Sequences

### Browser JS sign or anti-bot

- boundary and request proof: `chrome_devtools` -> `js_reverse`
- browser/local divergence: `js_reverse`
- replay confirmation: `burp`

### Android sign or encrypt

- runtime-first app-presence check and packet check: `scrcpy_vision` -> `adb_mcp` -> `charles` / `burp`
- Java recovery when blocked: `jadx`
- UI-state steering and screenshot-guided next actions: `scrcpy_vision`
- device state and runtime context: `adb_mcp`
- narrow Java or JNI hooks: `frida_mcp`
- dumped `.so` analysis when required: `ida_pro_mcp`
- wire validation or Charles-assisted observation: `charles`
- replay confirmation: `burp`

## Rule

Do not start in reverse when the relevant HTTP/HTTPS request or WebSocket message has not even been checked in Burp or Charles.
For Android app pentests, first confirm the target app is installed on the connected device before deeper workflow branching.
For Android external URL testing, do not reverse first when screenshot, logs, and packet visibility can answer the question.
Do not choose browser references by clue words before the current stage is known.


---

## Source: reporting-and-evidence.md

Path: references/reporting-and-evidence.md

# Reporting And Evidence

Minimum output:

- scope and client type
- chosen MCP chain
- static findings
- runtime proof
- recovered request recipe
- Burp-ready baseline request
- security finding and mitigation

## References — MCP

# MCP can力totaldocument

## 1. documenteye/look's/of

thisdocumentwhole/integer principle/logic(past tense)when before/frontSessioninIcandirectreceive/connectcall/invoke's/of MCP can力，goal/targetnotis (classifier) doone份“toolclearsingle”，而isprovideone份suitable for after/back续editwrite `skills` 's/ofreferencebottom稿。  
 re-/heavypoint覆stamp with/bydescendcontent：

- each MCP Server/命 nameemptybetween's/ofdefinebit
- eachmethod's/ofcall/invokeway/manner
- mainneed toparameter's/ofcontain/include义
- returnsresult/outcomelarge致will/canincludes/containswhat
- 典 typeusescenario
- and/withother MCP combinationtime's/ofcommonWorkflow

this文default面 to/towards Codex / Agent  category/classtoolorchestration，notisgeneral/universal SDK document。thereforewill/can更strong调“whatwhen/time useit”“write skill timehowdescriptioncall/invokestrategy”。

---

## 2. general/universalcall/invokeabout/approximatelydefine

### 2.1 tool命 nameformat

when before/frontenvironmentin's/of MCP tool namelargemulti/multiple遵循descend面format：

```text
mcp__<server_name>__<tool_name>
```

for example：

- `mcp__adb_mcp__list_devices`
- `mcp__chrome_devtools__navigate_page`
- `mcp__ida_pro_mcp__decompile`

fewand/with MCP resourceSourceAccess相 close/shut's/offunctionnotbring/carry `mcp__`  before/front缀，但this质ascendalsois MCP generate/live态can力：

- `list_mcp_resources`
- `list_mcp_resource_templates`
- `read_mcp_resource`

### 2.2 call/invokeparameterformat

placehas/have MCP toolalluse JSON 风format/gridparameter for/to象。典 typeformat：

```json
{
  "device_id": "emulator-5554",
  "lines": 200
}
```

Notepoint：

-  (classifier)transmitneed's/ofword paragraph/segment，do notno/withoutmeaning/intent义 (adverbial)塞emptyArrayor `null`
- `optional` parametergenerallycan省strategy
- certain/sometoolneed to求绝 for/toPath，尤its/theiris截Graph、saveSourcecode、PullFile、record屏inputexitPathetc.
- certain/sometooluse part/point页parameter，like/such as `offset`、`count`、`pageIdx`、`pageSize`

### 2.3 write skill timeRecommendationdescription's/ofneed topoint

like/such as resultyouneed to (object marker)thesecan力write become/success skill，Recommendationeach skill brightcertainwriteexit：

1. triggercondition  
2. advantage firstuse's/of MCP  
3. toolbetween's/of first after/backsequential  
4. 哪些parameterismust补all/full's/of  
5. what情况descendswitch toother MCP  
6. like/such as resultinputexitfor/isempty/failure，descendone步shouldhow补救

### 2.4 MCP 选 typequick reference

| Tasktype | advantage first MCP |
| --- | --- |
| Android set upbackupmanage、Installation APK、point击滑move、拉File | `adb_mcp` |
| Android canlook-izecontrol、UI Treedefinebit、no/without线 ADB、Real-time画面 | `scrcpy_vision` |
| Android grab/capture HTTP/HTTPS Streamquantity/measure、Charles SessionAnalysis | `charles` |
| Burp historical、Repeater、Collaborator、Intruder | `burp` |
| network页Automatic-ize、截Graph、tablesingle、networkrequest、Console | `chrome_devtools` |
| JS break/judgepoint、SourcecodeSearch、XHR send/issuestartchain、functiontrace | `js_reverse` |
| 官directiondocument检索、codeexamplequery | `context7` |
| general/universalnetwork页grab/capturetake/get/Pullnetwork页content | `fetch` |
| LocalFileextremespeed/fastSearch | `everything_search` |
| Android dynamicInject、Frida attach/spawn | `frida_mcp` |
| BinarystaticAnalysis、IDA Batch re-/heavy命 name/negative/reverseCompile/typerepair/fix | `ida_pro_mcp` |
| APK negative/reverseCompile、Manifest、 category/class/method/xref query | `jadx` |
| remember忆Graph谱、grow期structure-izeremember忆 | `memory` |
|  repeatmixedissue/problem part/point步thinktest | `sequential_thinking` |

### 2.5 commoncombinationWorkflow

#### Android App Analysis

- static：`jadx`
- dynamic：`frida_mcp`
- grab/capturePackage：`charles`
- set upbackupcontrol：`adb_mcp`
- canlook-ize/UI Automatic-ize：`scrcpy_vision`

#### Web FrontendReverse

- pageoperation：`chrome_devtools`
- JS break/judgepointand/withSourcecodeSearch：`js_reverse`
- HTTP Replayand/withsecurityTest：`burp`

#### Native / APK So Reverse

- IDA staticAnalysis：`ida_pro_mcp`
- Runtime hook：`frida_mcp`
- set upbackupend(side)协助：`adb_mcp` / `scrcpy_vision`

---

## 3. MCP resourceSource category/classgeneral/universalinterface

这three category/classfunctionnotis具body业务Server，而is“Access MCP ServerExposeresourceSource”'s/ofgeneral/universalcan力。

### 3.1 `list_mcp_resources`

- effect/function：listsome/certain (counter) MCP Serverorplacehas/haveServerPublic's/ofresourceSource
- 典 type use途：找candirectreceive/connectRead's/ofFile、context、Database schema、configuration (classifier) paragraph/segment
- parameter：
  - `server`：can选，指defineServer name
  - `cursor`：can选， part/point页swim标
- suitable for skill 's/ofdescription： firstEnumerationresourceSource， againdecidesisno/notcall/invoke `read_mcp_resource`

example：

```json
{
  "server": "some_server"
}
```

### 3.2 `list_mcp_resource_templates`

- effect/function：listparameter-izeresourceSourcetemplate
- 典 type use途：discover“bring/carryparameterRead”'s/ofresourceSource，for examplepress/according totable name、press/according toprimary key、press/according toPathquery's/ofresourceSource
- parameter：
  - `server`
  - `cursor`
- suitable for skill 's/ofdescription：whenresourceSourcenotisfixed URI，而is“template URI”time first查this

### 3.3 `read_mcp_resource`

- effect/function：Read具bodyresourceSourcecontent
- parameter：
  - `server`：Server name
  - `uri`：resourceSource URI
- suitable forscenario：
  - 读configuration
  - 读 schema
  - 读Servicecontext
  - 读together/shareenjoystate

example：

```json
{
  "server": "some_server",
  "uri": "resource://example/path"
}
```

---

## 4. `adb_mcp`：Android set upbackupcontroland/withFileinteractive

### 4.1 definebit

`adb_mcp` ismostfoundation/basis's/of Android set upbackupinteractivelayer，suitable for do：

- set upbackupcolumntableand/withstateAcknowledgment
- Installation/Unmount APK
- 截Graph、record屏
- inputenter文this、point击、滑move、send/issuepress/according tokey
- Pull/PushFile
- Read logcat、电池、memory、storeinformation

like/such as resultyou's/of skill need“controlset upbackupthis身”，advantage firsttestconsiderit。

### 4.2 commonWorkflow

1. `list_devices` Acknowledgmentset upbackup  
2. `get_device_info` / `get_battery_info` judgebreak/judgeenvironment  
3. `install_app` or `list_packages`  
4. `send_tap` / `send_swipe` / `send_text` Driverinteractive  
5. `take_screenshot` / `record_screen` stay/keep证据  
6. `get_logcat` linewrong  

### 4.3 methodclearsingle

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__adb_mcp__list_devices` | no/without | listConnection's/of Android set upbackup | Taskenter口， firstAcknowledgmentset upbackupisno/notonline |
| `mcp__adb_mcp__get_device_info` | `device_id?` | Readset upbackup详fineinformation | look/see typenumber、systemversion、Sequence Number |
| `mcp__adb_mcp__get_battery_info` | `device_id?` | Read电池state | growtimeTest before/frontAcknowledgment电quantity/measure |
| `mcp__adb_mcp__get_memory_info` | `device_id?` | Readmemoryinformation | property/naturecan/稳qualitativeline查 |
| `mcp__adb_mcp__get_storage_info` | `device_id?` | Readstoreinformation | look/seeemptybetweenisno/notenoughInstallation/record屏 |
| `mcp__adb_mcp__clear_logcat` | `device_id?` | empty logcat |  doone next/timedrycleangrab/captureLog |
| `mcp__adb_mcp__get_logcat` | `device_id?`, `filter_tag?`, `lines?` | ReadLog | 崩溃、network、SSL、Debuglinewrong |
| `mcp__adb_mcp__install_app` | `apk_path`, `device_id?` | Installation APK | deploymentTestPackage |
| `mcp__adb_mcp__uninstall_app` | `package_name`, `device_id?` | Unmountshould use | Cleanupenvironment |
| `mcp__adb_mcp__list_packages` | `device_id?`, `system_apps?` | listInstallationPackage name | 找goal/targetPackage name |
| `mcp__adb_mcp__list_files` | `remote_path`, `device_id?` | viewset upbackupDirectory | 找cache、configuration、ExportFile |
| `mcp__adb_mcp__pull_file` | `remote_path`, `local_path`, `device_id?` |  fromset upbackup拉File toLocal | ExportDatabase、Log、cache |
| `mcp__adb_mcp__push_file` | `local_path`, `remote_path`, `device_id?` | 推File toset upbackup | 推Certificate、footthis、Patch |
| `mcp__adb_mcp__send_keyevent` | `keycode`, `device_id?` | Sendpress/according tokeyevent | returnskey、Home、菜singlekey |
| `mcp__adb_mcp__send_tap` | `x`, `y`, `device_id?` | point击sit标 | Automatic-izeoperation |
| `mcp__adb_mcp__send_swipe` | `x1`,`y1`,`x2`,`y2`,`duration?`,`device_id?` | 滑move | 滚movecolumntable、unlock、切页 |
| `mcp__adb_mcp__send_text` | `text`, `device_id?` | inputenter文this | Search、login、tablesingleinputenter |
| `mcp__adb_mcp__take_screenshot` | `save_path`, `device_id?` | 截Graph toLocal | 证据keepstay/keep、UI stateAcknowledgment |
| `mcp__adb_mcp__record_screen` | `duration?`, `save_path?`, `device_id?` | record屏 |  repeatpresentprocessstay/keep证 |

### 4.4 典 typecall/invokeexample

columnset upbackup：

```json
{}
```

截Graph：

```json
{
  "device_id": "emulator-5554",
  "save_path": "C:\\Users\\28484\\Desktop\\screen.png"
}
```

Readmostnear 200 rowLog：

```json
{
  "device_id": "emulator-5554",
  "lines": 200
}
```

### 4.5 write skill time's/ofNotepoint

- any Android Task几乎allshould firstrunone next/time `list_devices`
- `take_screenshot` brightcertainneed to求Local绝 for/toPath
- `get_logcat` at/in repeatmixedscenarioinRecommendation first `clear_logcat`
- `send_tap` / `send_swipe`  completeall/fulldepend onsit标，suitable forfixedboundary面，notsuitable forstrongdynamic布game
- `push_file` and/with `pull_file` is doCertificateInstallation、LogExport、datastay/keep证's/ofhigh频tool

---

## 5. `charles`：Charles grab/capturePackageand/withSessionAnalysis

### 5.1 definebit

`charles` defeat责ReadandAnalysis Charles Proxy alreadyCapture's/ofStreamquantity/measure， re-/heavypointnotis“directreceive/connectcontrol Android Proxy”，而is：

- Inspect/Check Charles isno/notonline、isno/notalreadyhas/haveactivitygrab/capturePackageSession
- Startorreceive/connect管 live capture，take to `capture_id`
- structure-izescreen live traffic oralreadysave recording
- descend钻single (classifier)request，viewhead/top、statecode、Request Body/responsebody预view
-  for/toStreamquantity/measurepress/according to host、path、status、resource class Block/GroupAnalysis
- tie/knotbindgrab/capturePackage并PersistenceSnapshot，directionthen after/back续 repeatround

### 5.2 suitable for's/of skill type

- Android API Reverse
- HTTPS grab/capturePackage
- App interfacerowfor/isAnalysis
- parameterSignature before/front after/backcomparison
- Find token、session、Encryptionword paragraph/segment
- Sessionrecordmake/control、screenand/with证据stay/keepexist

### 5.3 methodclearsingle

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__charles__charles_status` | no/without | Inspect/Check Charles 连commonproperty/natureand/with live capture state | Acknowledgmentenvironmentisno/notthen绪 |
| `mcp__charles__reset_environment` | no/without | Reset Charles environment并recoverysave's/ofconfiguration |  dodrycleansolid验 |
| `mcp__charles__start_live_capture` | `adopt_existing?`,`include_existing?`,`reset_session?` | Startorreceive/connect管 live capture | Get after/back续Analysisneed to use's/of `capture_id` |
| `mcp__charles__query_live_capture_entries` | `capture_id`,`cursor?`,`preset?`,`host_contains?`,`path_contains?`,`method_in?`,`status_in?`,`request_body_contains?`,`response_body_contains?`,`max_items?` | structure-izescreen live Streamquantity/measure | Recommendation's/ofReal-time检索enter口 |
| `mcp__charles__peek_live_capture` | `capture_id`,`cursor?`,`limit?` | 预viewwhen before/front live capture in's/ofnew (classifier)eye/look | lightquantity/measureviewmostnearrequest |
| `mcp__charles__read_live_capture` | `capture_id`,`cursor?`,`limit?` | increasequantity/measureRead并推advance live cursor | needStream style/modeReadnewStreamquantity/measuretimeuse |
| `mcp__charles__get_traffic_entry_detail` | `source`,`entry_id`,`capture_id?`,`recording_path?`,`include_full_body?`,`max_body_chars?` | descend钻single (classifier)Streamquantity/measure详情 | look/seehead/top、body 预view、requestresponsefine section |
| `mcp__charles__group_capture_analysis` | `source`,`capture_id?`,`recording_path?`,`group_by`,`preset?`,`host_contains?`,`path_contains?`,`status_in?` | press/according to host/path/status/resource class Block/Group | fastspeed/fast找hotpointinterface |
| `mcp__charles__get_capture_analysis_stats` | `source`,`capture_id?`,`recording_path?`,`preset?` | returnscoarse粒degree/measurestatistics | look/seegrab/capturePackageall/fullgame part/point布 |
| `mcp__charles__stop_live_capture` | `capture_id`,`persist?` | Stop live capture 并canPersistence | tie/knotbindsolid验并saveSnapshot |
| `mcp__charles__list_recordings` | no/without | listalreadysaverecordmake/controlFile | selecthistoricalStreamquantity/measurePackage |
| `mcp__charles__list_sessions` | no/without | compatibilityway/mannerlisthistorical session | compatibilityold命 name |
| `mcp__charles__get_recording_snapshot` | `path?` | Readalreadysaverecordmake/control's/ofSnapshot元information | offlineInspect/Check recording |
| `mcp__charles__analyze_recorded_traffic` | `recording_path?`,`preset?`,`host_contains?`,`path_contains?`,`method_in?`,`status_in?`,`request_body_contains?`,`response_body_contains?`,`max_items?` | Analysishistoricalrecordmake/control | offlinereturnlook/seeand/with repeatround |
| `mcp__charles__query_recorded_traffic` | `host_contains?`,`http_method?`,`keyword_regex?`,`keep_request?`,`keep_response?` | querylatestsave's/of recording | fastspeed/fastFilterhistoricalStreamquantity/measure |
| `mcp__charles__proxy_by_time` | `record_seconds` | press/according tofixedtimegrowgrab/capturetake/getorReadlatesthistoricalPackage | fastspeed/fasttime窗Analysis |
| `mcp__charles__filter_func` | `capture_seconds`,`host_contains?`,`http_method?`,`keyword_regex?`,`keep_request?`,`keep_response?` | press/according totime窗andconditionFilterStreamquantity/measure | fastspeed/fastshrinksmall范围 |
| `mcp__charles__throttling` | `preset` | setting Charles weaknetwork/limitspeed/fastpreset | weaknetwork repeatpresentand/withrowfor/isValidate |

### 5.4 RecommendationWorkflow

1. `charles_status`  
2. Acknowledgment Charles alreadyEnable/OnListen，Android Proxyalreadypoints tograb/capturePackagemachine，HTTPS needtimealreadyInstallation Charles Certificate  
3. `reset_environment`（can选， dodrycleansolid验）  
4. `start_live_capture`  
5. operation App  
6. `query_live_capture_entries`  
7. `get_traffic_entry_detail`  
8. `group_capture_analysis` / `get_capture_analysis_stats`  
9. `stop_live_capture`，必need totimesetting `persist: true`  
10. `analyze_recorded_traffic` / `query_recorded_traffic`

### 5.5 call/invokeexample

StartReal-timegrab/capturePackage：

```json
{
  "reset_session": true,
  "include_existing": false
}
```

screenReal-timeinterfaceStreamquantity/measure：

```json
{
  "capture_id": "capture-id-from-start",
  "preset": "api_focus",
  "host_contains": "api.example.com",
  "max_items": 10
}
```

### 5.6 Notepoint

- `charles` MCP will not替youconfiguration Android systemProxy；need to first complete become/success Charles Listen、set upbackupProxyandCertificateaccuratebackup
- Real-time检索advantage first use `query_live_capture_entries`，do notdefault usewill/can推advanceswim标's/of `read_live_capture`
- `get_traffic_entry_detail` default (classifier)look/see预view更省context， (classifier)has/haveindeedneedOriginal Texttime again open `include_full_body`
- like/such as resultthink repeatroundgrab/capturePackageresult/outcome，tie/knotbind live capture timeRecommendation `persist: true`
- like/such as result Charles alreadyat/inRunandyounotthinkemptywhen before/frontSession， use `adopt_existing: true`

---

## 6. `burp`：Burp Suite 协same/togetheroperation

### 6.1 definebit

`burp` MCP is面 to/towards Burp Suite 's/ofcontroland/withdataAccesslayer，suitable for：

- ReadProxyhistorical
-  (object marker)requestsend off to Repeater / Intruder
- send/issue HTTP/1.1、HTTP/2 request
- generate Collaborator payload
- look/seeScanningdeviceissue/problem
- read-writewhen before/frontediteditdevicecontent
- adjustmentProxyIntercept、TaskExecutestate
- read-write Burp configuration

### 6.2 methodclearsingle

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__burp__base64_encode` | `content` | Base64 Encoding | construct payload |
| `mcp__burp__base64_decode` | `content` | Base64 Decoding | look/seeEncodingdata |
| `mcp__burp__url_encode` | `content` | URL Encoding | constructparameter |
| `mcp__burp__url_decode` | `content` | URL Decoding | restorationparameter |
| `mcp__burp__generate_random_string` | `length`,`characterSet` | generatefollowmachine串 | token、boundary/perimetervalue、detect/probe串 |
| `mcp__burp__get_active_editor_contents` | no/without | Getwhen before/frontediteditdevicecontent | Read手工editeditrequest |
| `mcp__burp__set_active_editor_contents` | `text` | settingwhen before/frontediteditdevicecontent | Automatic填enterrequesttemplate |
| `mcp__burp__create_repeater_tab` | `content`,`targetHostname`,`targetPort`,`usesHttps`,`tabName?` | newbuild Repeater tag/label页 | send offrequest to Repeater |
| `mcp__burp__send_to_intruder` | `content`,`targetHostname`,`targetPort`,`usesHttps`,`tabName?` | send off to Intruder | brute force/BatchTest |
| `mcp__burp__send_http1_request` | `content`,`targetHostname`,`targetPort`,`usesHttps` | send/issue HTTP/1.1 request | exactReplay |
| `mcp__burp__send_http2_request` | `pseudoHeaders`,`headers`,`requestBody`,`targetHostname`,`targetPort`,`usesHttps` | send/issue HTTP/2 request | H2 specificscenario |
| `mcp__burp__generate_collaborator_payload` | `customData?` | generate OOB Domain Name | SSRF / RCE / Blind XXE Test |
| `mcp__burp__get_collaborator_interactions` | `payloadId?` | round询 OOB interactive | look/seeisno/notoutbound |
| `mcp__burp__get_proxy_http_history` | `count`,`offset` | ReadProxy HTTP historical | returnlook/seerequest |
| `mcp__burp__get_proxy_http_history_regex` | `count`,`offset`,`regex` | press/according tocorrect/positive rule/principleFilter HTTP historical | exactscreen |
| `mcp__burp__get_proxy_websocket_history` | `count`,`offset` | Read WS historical | Analysis WebSocket |
| `mcp__burp__get_proxy_websocket_history_regex` | `count`,`offset`,`regex` | correct/positive rule/principleFilter WS historical | 查 token、commandword paragraph/segment |
| `mcp__burp__get_scanner_issues` | `count`,`offset` | listScanningdevicediscover | vulnerability巡检 |
| `mcp__burp__output_project_options` | no/without | Exportitemeye/looklevel/gradeconfiguration | viewconfiguration schema |
| `mcp__burp__output_user_options` | no/without | Exportuserlevel/gradeconfiguration | viewconfiguration schema |
| `mcp__burp__set_project_options` | `json` | settingitemeye/looklevel/gradeconfiguration | Automatic-ize调advantage |
| `mcp__burp__set_user_options` | `json` | settinguserlevel/gradeconfiguration | userall/fullgameconfiguration |
| `mcp__burp__set_proxy_intercept_state` | `intercepting` |  open close/shutProxyIntercept |  open/ close/shut Intercept |
| `mcp__burp__set_task_execution_engine_state` | `running` |  open close/shutTaskExecutelead/guide擎 | Pause/recoveryScanningTask |

### 6.3 典 typecall/invokeexample

Create Repeater：

```json
{
  "content": "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n",
  "targetHostname": "example.com",
  "targetPort": 443,
  "usesHttps": true,
  "tabName": "home"
}
```

generate Collaborator：

```json
{
  "customData": "ssrf-test"
}
```

### 6.4 Notepoint

- `send_http2_request` 's/ofRequest Bodyandhead/topissplit open's/of，do not (object marker)head/topwriteadvance body
- 改configuration before/frontRecommendation first `output_project_options` / `output_user_options`
- OOB detectiongenerallyis：`generate_collaborator_payload` -> Inject业务point -> `get_collaborator_interactions`
- `get_proxy_http_history_regex` verysuitable forwrite skill time do“Automaticscreen相 close/shuthistoricalrequest”

---

## 7. `chrome_devtools`：BrowserAutomatic-ize、page诊break/judgeand/withproperty/naturecanAnalysis

### 7.1 definebit

`chrome_devtools` defeat责Browserpage's/ofAutomatic-izecontroland/with DevTools level/gradeobserve。corecan力including：

- 打 open/Disable/Off/selectpage
- guide航、Refresh、simulateset upbackup
- DOM Snapshot、截Graph
- point击、inputenter、UploadFile
- columntable-izenetworkrequestandConsoleinformation
- Executepagefootthis
- Lighthouse Audit
- property/naturecan trace
- HeapSnapshot

like/such as resultyouneed to“像人at/inBrowserinoperationpage”，itispreferred。

### 7.2 pageand/withcontextcontrol

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__chrome_devtools__list_pages` | no/without | listwhen before/front打 open's/ofpage |
| `mcp__chrome_devtools__new_page` | `url`,`background?`,`isolatedContext?`,`timeout?` | newbuildtag/label页并Access URL |
| `mcp__chrome_devtools__select_page` | `pageId`,`bringToFront?` | switchwhen before/frontoperationpage |
| `mcp__chrome_devtools__close_page` | `pageId` | Disable/Offpage |
| `mcp__chrome_devtools__navigate_page` | `type`,`url?`,`timeout?`,`ignoreCache?`,`handleBeforeUnload?`,`initScript?` | URL guide航、 before/frontadvance、 after/backretreat、Refresh |
| `mcp__chrome_devtools__resize_page` | `width`,`height` | adjustmentBrowser尺寸 |
| `mcp__chrome_devtools__emulate` | `viewport?`,`colorScheme?`,`geolocation?`,`networkConditions?`,`userAgent?`,`cpuThrottlingRate?` | set upbackup/network/UA simulate |

### 7.3 pagestructureand/with截Graph

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__chrome_devtools__take_snapshot` | `filePath?`,`verbose?` | Getpage a11y TreeSnapshot，returnselement `uid` |
| `mcp__chrome_devtools__take_screenshot` | `filePath?`,`format?`,`fullPage?`,`quality?`,`uid?` | pageorelement截Graph |
| `mcp__chrome_devtools__wait_for` | `text`,`timeout?` | etc.待certain/some文thisexitpresent |

explanation：

-  first `take_snapshot`， againusein面's/of `uid` go/leave do click/fill/hover，usuallymost稳
- `uid` iswhen before/frontSnapshotcontextin's/ofelementidentifier，SnapshotUpdate after/backcancanchange

### 7.4 pageinteractive

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__chrome_devtools__click` | `uid`,`dblClick?`,`includeSnapshot?` | point击element |
| `mcp__chrome_devtools__hover` | `uid`,`includeSnapshot?` | 悬stopelement |
| `mcp__chrome_devtools__drag` | `from_uid`,`to_uid`,`includeSnapshot?` | drag |
| `mcp__chrome_devtools__fill` | `uid`,`value`,`includeSnapshot?` | 填single (counter)inputenterbox |
| `mcp__chrome_devtools__fill_form` | `elements`,`includeSnapshot?` | Batch填tablesingle |
| `mcp__chrome_devtools__type_text` | `text`,`submitKey?` |  to/towardswhen before/front焦pointinputenter文this |
| `mcp__chrome_devtools__press_key` | `key`,`includeSnapshot?` | keyroundfast捷key、specialpress/according tokey |
| `mcp__chrome_devtools__upload_file` | `uid`,`filePath`,`includeSnapshot?` | UploadFile |
| `mcp__chrome_devtools__handle_dialog` | `action`,`promptText?` | process/handle alert/confirm/prompt |

### 7.5 pagefootthisand/withDebuginformation

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__chrome_devtools__evaluate_script` | `function`,`args?` | at/inpageinner/insideExecute JS |
| `mcp__chrome_devtools__list_console_messages` | `includePreservedMessages?`,`pageIdx?`,`pageSize?`,`types?` | viewConsoleLog |
| `mcp__chrome_devtools__get_console_message` | `msgid` | Getsingle (classifier)Consolemessage详情 |
| `mcp__chrome_devtools__list_network_requests` | `includePreservedRequests?`,`pageIdx?`,`pageSize?`,`resourceTypes?` | viewnetworkrequestcolumntable |
| `mcp__chrome_devtools__get_network_request` | `reqid?`,`requestFilePath?`,`responseFilePath?` | vieworExportrequest详情/body |

### 7.6 Auditand/withproperty/naturecan

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__chrome_devtools__lighthouse_audit` | `device?`,`mode?`,`outputDirPath?` | run Lighthouse（notcontain/includeproperty/naturecan part/point） |
| `mcp__chrome_devtools__performance_start_trace` | `autoStop?`,`filePath?`,`reload?` | Startproperty/naturecan trace |
| `mcp__chrome_devtools__performance_stop_trace` | `filePath?` | Stopproperty/naturecan trace |
| `mcp__chrome_devtools__performance_analyze_insight` | `insightName`,`insightSetId` | Analysissome/certain (counter)property/naturecan insight |
| `mcp__chrome_devtools__take_memory_snapshot` | `filePath` | Export JS HeapSnapshot |

### 7.7 RecommendationWorkflow

#### pageAutomatic-ize

1. `new_page`
2. `take_snapshot`
3. `click` / `fill` / `press_key`
4. `wait_for`
5. `take_screenshot`

#### grab/capturepagerequest

1. `new_page`
2. pageinteractive
3. `list_network_requests`
4. `get_network_request`

#### property/naturecanline查

1. `navigate_page`
2. `performance_start_trace`
3. pageoperationor reload
4. `performance_stop_trace`
5. `performance_analyze_insight`

### 7.8 Notepoint

-  do DOM interactive before/frontadvantage first `take_snapshot`
- pageRefresh after/backold `uid` notonedefinestillcan use
- GetRequest Body/responsebodytime，必need totime use `requestFilePath` / `responseFilePath` fall (adverbial) toFile
- 若you close/shutnote“JS call/invokechainandbreak/judgepoint”，`js_reverse`  to/towards to/towards比这in更suitable for

---

## 8. `context7`：Real-timedocumentand/withexample检索

### 8.1 definebit

`context7` suitable forqueryNo.threedirectionLibrary、Framework、官directiondocumentandcodeexample，尤its/theirsuitable for技caneditwritein“need tocitationlatest官direction use method/law”'s/ofscenario。

### 8.2 method

#### `mcp__context7__resolve_library_id`

- effect/function： first (object marker)“Library name”parsing become/success Context7 canidentify's/ofdocument ID
- parameter：
  - `libraryName`
  - `query`
- returns re-/heavypoint：
  - `libraryId`
  - Library name
  - description
  - snippets numberquantity/measure
  - source reputation
  - benchmark score

#### `mcp__context7__query_docs`

- effect/function：based onalreadyparsingexit's/of `libraryId` 检索documentandexample
- parameter：
  - `libraryId`
  - `query`

### 8.3 RecommendationWorkflow

1. `resolve_library_id`
2. 选mostcombine适's/of `libraryId`
3. `query_docs`

### 8.4 example

 firstparsing：

```json
{
  "libraryName": "Next.js",
  "query": "App Router middleware authentication examples"
}
```

 againquery：

```json
{
  "libraryId": "/vercel/next.js",
  "query": "How to protect routes in App Router middleware?"
}
```

### 8.5 write skill 's/ofNotepoint

- like/such as resultuser to/for's/ofisfuzzy/blurLibrary name， first `resolve_library_id`
- 这is“documentaskanswer MCP”，notis联networkfollowthen搜network页
-  for/totechniqueissue/problem，advantage first (object marker)itwhen as/do“官directiondocument检索device”

---

## 9. `everything_search`：LocalFileextremespeed/fastSearch

### 9.1 definebit

这is Windows LocalFileSearch MCP，suitable forlargeDirectory、all/fullround、fuzzy/blurconditiondescendfastspeed/fast找File。

### 9.2 method

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__everything_search__search` | `query`,`maxResults?`,`parentPath?`,`filesOnly?`,`foldersOnly?`,`matchPath?`,`regex?`,`caseSensitive?`,`wholeWord?`,`sortBy?`,`sortDescending?`,`showSize?`,`showDateModified?` | SearchFileorDirectory |
| `mcp__everything_search__get_file_info` | `filename` | Getsome/certain (counter)File详fineinformation |

### 9.3 example

Search指defineDirectorydescend's/ofplacehas/have `.apk`：

```json
{
  "query": "*.apk",
  "parentPath": "C:\\Users\\28484",
  "filesOnly": true,
  "maxResults": 50
}
```

### 9.4 适 usescenario

- 找 APK / SO / Log / ExportFile
-  to/forReverse category/class skill 找goal/targetFile
- at/inlargeDirectoryin找configuration、footthis、Database、Certificate

---

## 10. `fetch`：general/universalnetwork页grab/capturetake/get

### 10.1 definebit

`fetch` is“grab/capturetake/getnetwork页/URL content”'s/ofgeneral/universaltool，suitable for：

- 拉network页content
- grab/capturedocument页
- Read HTML
-  dosimplesinglenetwork页contentextract

### 10.2 method

#### `mcp__fetch__fetch`

- parameter：
  - `url`
  - `max_length?`
  - `raw?`
  - `start_index?`
- effect/function：
  - Getnetwork页content
  - canreturnssimplifies after/back's/of markdown  style/modecontent
  - can指defineoffsetcontinue读growpage

### 10.3 example

```json
{
  "url": "https://example.com",
  "max_length": 6000
}
```

### 10.4 Notepoint

- 更suitable for“Known URL 's/ofcontentgrab/capturetake/get”，notisSearchlead/guide擎
- like/such as resultpage太grow，canvia/through `start_index` FragmentationRead
- techniquedocumentscenarioin，like/such ashas/have `context7`，usuallyadvantage first `context7`

---

## 11. `frida_mcp`：Android dynamicInjectand/withRuntime Hook

### 11.1 definebit

`frida_mcp` is Android dynamicAnalysislayer，core use途：

- Inspect/Check/Start/Stop `frida-server`
- Enumerationshould use
- Getwhen before/front before/front (classifier for machines)should use
- `spawn` or `attach`  togoal/targetProcess
- Inject Frida JS footthis
- GetfootthisinputexitLog

suitable for's/ofscenario：

- SSL Pinning bypass
- methodparameter/returnsvalue打print
- dynamicgrab/captureSignature、token、header
- native/Java layerRuntimeobserve

### 11.2 methodclearsingle

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__frida_mcp__check_frida_status` | no/without | view frida-server isno/notRun |  before/frontplaceInspect/Check |
| `mcp__frida_mcp__start_frida_server` | no/without | Start frida-server | dynamicAnalysisaccuratebackup |
| `mcp__frida_mcp__stop_frida_server` | no/without | Stop frida-server | Cleanupenvironment |
| `mcp__frida_mcp__list_applications` | no/without | listset upbackupshould use | 找Package name、look/seeisno/notRunmiddle/center |
| `mcp__frida_mcp__get_frontmost_application` | no/without | Getwhen before/front before/front (classifier for machines)should use | Acknowledgmentwhen before/frontboundary面placebelongPackage name |
| `mcp__frida_mcp__spawn` | `package_name`,`initial_script?`,`script_file_path?`,`output_file?` | 挂startStart并attachaddgoal/targetshould use | early期timemachine hook |
| `mcp__frida_mcp__attach` | `target`,`initial_script?`,`script_file_path?`,`output_file?` | attachadd to PID orPackage name |  for/toalreadyRunshould useInject |
| `mcp__frida_mcp__get_messages` | `max_messages?` | Get hook/log inputexitbuffer | look/seefootthis打printresult/outcome |

### 11.3 `attach` and/with `spawn` 's/ofdifferencepart

- `attach`
  - used for/forgoal/targetalreadyat/inRun
  - canpress/according to PID orPackage nameattachadd
  - suitable fortemporaryobserve、late期 hook

- `spawn`
  - used for/forat/inshould userecovery before/frontInjectfootthis
  - suitable forearly期 category/classLoad、Startprocess、Signatureinitial-ize、SSL pinning early期bypass

### 11.4 example

Inspect/Checkstate：

```json
{}
```

press/according toPackage nameStart并InjectfootthisFile：

```json
{
  "package_name": "com.example.app",
  "script_file_path": "C:\\Users\\28484\\Desktop\\hook.js",
  "output_file": "C:\\Users\\28484\\Desktop\\frida.log"
}
```

attachaddalreadyRunshould use并directreceive/connectwriteinner/inside联footthis：

```json
{
  "target": "com.example.app",
  "initial_script": "Java.perform(function(){ console.log('hook loaded'); });"
}
```

### 11.5 RecommendationWorkflow

1. `check_frida_status`
2. 若un-Run rule/principle `start_frida_server`
3. `list_applications` or `get_frontmost_application`
4. `spawn` or `attach`
5. `get_messages`

### 11.6 Notepoint

- needset upbackupenvironmentcorrect/positivecertaindeployment `frida-server`
- `script_file_path` Priorityhigh at/in `initial_script`
- mostSignature/EncryptiondefinebitTaskusuallyis：`jadx` staticdefinebit -> `frida_mcp` dynamicValidate

---

## 12. `ida_pro_mcp`：IDA Pro staticAnalysisand/withbatchprocess/handle re-/heavyconstruct

### 12.1 definebit

`ida_pro_mcp` iswhen before/frontcan力inmost re-/heavy's/ofstaticAnalysis MCP。itnotis“ (classifier)look/seenegative/reverseCompile”，而is覆stamp：

- 打 open/switch IDA instance
- fastspeed/fast survey Binary
- columnfunction、all/fullgame、Import、type
- 查 xref / callgraph / basic block
- negative/reverseCompile、negative/reverseAssemble、Exportfunctioninformation
- Modifycomment、 re-/heavy命 name、declaretype、CreateStackvariable
- 读memory、Patchbyte、PatchAssemble
-  use Python at/in IDA contextExecutefootthis

like/such as result skill is面 to/towards native Reverse、maliciouscodeAnalysis、Patch、Batch re-/heavy命 name，it几乎iscore。

### 12.2 strong烈Recommendation's/ofenter口tool

#### `mcp__ida_pro_mcp__survey_binary`

这ismostsuitable for doNo.one步 triage 's/oftool。itcanone next/timeproperty/nature to/forexit：

- File元information
-  paragraph/segment布game
- entry point
- statisticsinformation
- high频string
- highpricevaluefunction
- imports classification
- call/invokeGraph概况

write skill timecanbrightcertainregulation：  
“ open startAnalysis IDB  after/back， firstcall/invoke `survey_binary`，do notdirectreceive/connect盲eye/look `list_funcs`。”

### 12.3 instanceand/withSession Management

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__list_instances` | no/without | listwhen before/frontcanConnection's/of IDA instance |
| `mcp__ida_pro_mcp__select_instance` | `port`,`host?` | switchwhen before/front MCP points to's/of IDA instance |
| `mcp__ida_pro_mcp__open_file` | `file_path`,`autonomous?`,`new_database?`,`switch?`,`timeout?` | 打 openFile tonew's/of IDA instance |
| `mcp__ida_pro_mcp__server_health` | no/without | look/seewhen before/front IDB/Service健康state |
| `mcp__ida_pro_mcp__server_warmup` | `build_caches?`,`init_hexrays?`,`wait_auto_analysis?` | 预hotAnalysisenvironment |
| `mcp__ida_pro_mcp__idb_save` | `path?` | savewhen before/front IDB |

### 12.4 Binarytotalviewand/withdiscover

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__survey_binary` | `detail_level?` | Binarytotalview |
| `mcp__ida_pro_mcp__entity_query` |  repeatmixedquery for/to象 | 查 functions/globals/imports/strings/names |
| `mcp__ida_pro_mcp__find_regex` | `pattern`,`limit?`,`offset?` | at/instringmiddle/center usecorrect/positive rule/principle查 |
| `mcp__ida_pro_mcp__find` | `targets`,`type`,`limit?`,`offset?` | 查string、immediatelynumber、data/codecitation |
| `mcp__ida_pro_mcp__find_bytes` | `patterns`,`limit?`,`offset?` | bytepatternSearch |

### 12.5 functionand/withGraphAnalysis

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__list_funcs` | `queries` | columnfunction |
| `mcp__ida_pro_mcp__func_query` | FilterconditionSet | press/according tolargesmall/ nameword/isno/nothas/havetypeFilterfunction |
| `mcp__ida_pro_mcp__func_profile` | querySet |  to/forfunction do概view画像 |
| `mcp__ida_pro_mcp__lookup_funcs` | `queries` | press/according toaddressornamequeryfunction |
| `mcp__ida_pro_mcp__callees` | `addrs`,`limit?` | 查by (passive)call/invokefunction |
| `mcp__ida_pro_mcp__callgraph` | `roots`,`max_depth?`,`max_nodes?`,`max_edges?`,`max_edges_per_func?` | buildcall/invokeGraph |
| `mcp__ida_pro_mcp__basic_blocks` | `addrs`,`offset?`,`max_blocks?` | Get CFG 基thisBlock |
| `mcp__ida_pro_mcp__analyze_function` | `addr`,`include_asm?` | 紧凑singlefunctionAnalysis |
| `mcp__ida_pro_mcp__analyze_batch` | `queries` | Batchmulti/multiplefunction综combineAnalysis |
| `mcp__ida_pro_mcp__analyze_component` | `addrs` |  for/toonegroup/set相 close/shutfunction doComponentAnalysis |

### 12.6 negative/reverseCompile、negative/reverseAssembleand/withExport

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__decompile` | `addr` | negative/reverseCompilefunction |
| `mcp__ida_pro_mcp__disasm` | `addr`,`offset?`,`max_instructions?`,`include_total?` | negative/reverseAssemblefunction |
| `mcp__ida_pro_mcp__export_funcs` | `addrs`,`format?` | Exportfunctionfor/is JSON / C head/top / original type |

### 12.7 交叉citationand/withdataStream

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__xrefs_to` | `addrs`,`limit?` | Get xrefs to |
| `mcp__ida_pro_mcp__xref_query` | querySet | press/according todirection/typeBatchquery xref |
| `mcp__ida_pro_mcp__trace_data_flow` | `addr`,`direction?`,`max_depth?` | tracemulti/multiplejumpdataStream |
| `mcp__ida_pro_mcp__xrefs_to_field` | `queries` | 查structurebodyword paragraph/segmentcitation |

### 12.8 typesystemand/withstructurerecovery

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__type_query` | querySet | 查Localtype |
| `mcp__ida_pro_mcp__type_inspect` | `queries` | viewtypedeclareand/with become/successmember |
| `mcp__ida_pro_mcp__declare_type` | `decls` | Inject C typedeclare |
| `mcp__ida_pro_mcp__set_type` | `edits` | settingfunction/variable/game partvariabletype |
| `mcp__ida_pro_mcp__type_apply_batch` | `batch` | Batchshould usetype |
| `mcp__ida_pro_mcp__infer_types` | `addrs` | inferencetype |
| `mcp__ida_pro_mcp__enum_upsert` | `queries` | Create/supplementEnumeration |
| `mcp__ida_pro_mcp__search_structs` | `filter` | 搜structurebody/联combinebody |
| `mcp__ida_pro_mcp__read_struct` | `queries` | Readsome/certainaddressplacestructurebodyword paragraph/segmentvalue |

### 12.9 Stack帧and/withgame partvariable

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__stack_frame` | `addrs` | GetfunctionStack帧 |
| `mcp__ida_pro_mcp__declare_stack` | `items` | declareStackvariable |
| `mcp__ida_pro_mcp__delete_stack` | `items` | DeleteStackvariable |

### 12.10  re-/heavy命 name、commentand/withdifferenceValidate

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__rename` | `batch` | Batch re-/heavy命 namefunction/data/game part/Stackvariable |
| `mcp__ida_pro_mcp__set_comments` | `items` | settingcomment |
| `mcp__ida_pro_mcp__append_comments` | `items` | appendcomment |
| `mcp__ida_pro_mcp__diff_before_after` | `addr`,`action`,`action_args` | should use rename/type/comment  after/backcomparison before/front after/backnegative/reverseCompile |

### 12.11 originalmemoryReadand/withPatch

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__get_bytes` | `regions` | 读byte |
| `mcp__ida_pro_mcp__get_int` | `queries` | 读whole/integernumber |
| `mcp__ida_pro_mcp__get_string` | `addrs` | 读string |
| `mcp__ida_pro_mcp__get_global_value` | `queries` | 读all/fullgamevariablevalue |
| `mcp__ida_pro_mcp__put_int` | `items` | writewhole/integernumber |
| `mcp__ida_pro_mcp__patch` | `patches` | Patchbyte |
| `mcp__ida_pro_mcp__patch_asm` | `items` | PatchAssemble |
| `mcp__ida_pro_mcp__undefine` | `items` | canceldefinefor/isoriginalbyte |
| `mcp__ida_pro_mcp__define_code` | `items` | will/shallbytedefinefor/iscode |
| `mcp__ida_pro_mcp__define_func` | `items` | definefunction |

### 12.12 Import、all/fullgame、指 makeand/withsolidbodyquery

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__imports` | `count`,`offset` | columnImport |
| `mcp__ida_pro_mcp__imports_query` | `queries` | press/according tomoduleBlock/ namewordFilterImport |
| `mcp__ida_pro_mcp__list_globals` | `queries` | columnall/fullgamevariable |
| `mcp__ida_pro_mcp__insn_query` | `queries` | query指 makepattern |
| `mcp__ida_pro_mcp__int_convert` | `inputs` | numberformatconversion |

### 12.13 Python Extension

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__ida_pro_mcp__py_eval` | `code` | at/in IDA environmentinExecute Python  (classifier) paragraph/segment |
| `mcp__ida_pro_mcp__py_exec_file` | `file_path` | Executewhole/integer (counter) Python footthisFile |

### 12.14 RecommendationWorkflow

#### initial triage

1. `server_health`
2. `server_warmup`
3. `survey_binary`
4. `find_regex` / `imports_query`
5. `analyze_function` / `decompile`

#### recovery language义

1. `decompile`
2. `stack_frame`
3. `type_query` / `type_inspect`
4. `set_type` / `declare_type`
5. `rename`
6. `diff_before_after`

#### TraceSensitivestring

1. `find_regex`
2. `xrefs_to`
3. `trace_data_flow`
4. `analyze_component`

### 12.15 skill editwriteRecommendation

- one open startthenwrite死“ first `survey_binary`”usuallyisgoodstrategy
- like/such as resultneed to doBatch re-/heavy命 name，mostgood (object marker) `diff_before_after` when become/successValidatestep
- need toAnalysis JNI / crypto / dispatch table，`trace_data_flow` veryhas/havepricevalue
- `type_apply_batch` suitable for do“Automatic修type” category/class skill
- `py_eval` / `py_exec_file` suitable for dohighlevel/gradeAutomatic-ize，但should谨慎definefootthisboundary/perimeter

---

## 13. `jadx`：APK staticnegative/reverseCompileand/with Android codeguide航

### 13.1 definebit

`jadx` MCP is Android staticAnalysisenter口，suitable for：

- 读 `AndroidManifest.xml`
- 找main Activity、Component、ExportComponent
- Search category/class/method/word paragraph/segment
- Get category/classSourcecode、methodSourcecode、smali
- 查citation close/shut system/relationship
-  re-/heavy命 name category/class/method/word paragraph/segment/variable/Package

itand `ida_pro_mcp` 's/ofdifferenceat/in at/in：

- `jadx` 更偏 Java/Kotlin layer APK
- `ida_pro_mcp` 更偏 native Binary / so / ELF / PE

### 13.2 enter口informationand/with Manifest

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__jadx__get_android_manifest` | no/without | Get Manifest all/full文 |
| `mcp__jadx__get_main_activity_class` | no/without | Getmain Activity |
| `mcp__jadx__get_main_application_classes_names` | no/without | Getmainshould usePackagedescendmainneed to category/class name |
| `mcp__jadx__get_main_application_classes_code` | `count?`,`offset?` | Getmainneed to category/classcode |
| `mcp__jadx__get_manifest_component` | `component_type`,`only_exported?` | Get activity/service/provider/receiver Componentinformation |

### 13.3  category/classand/withSourcecodeRead

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__jadx__get_all_classes` | `count?`,`offset?` | Getplacehas/have category/class name |
| `mcp__jadx__fetch_current_class` | no/without | take/get GUI when before/front选middle/center category/classSourcecode |
| `mcp__jadx__get_class_source` | `class_name` | Getsome/certain category/class Java Sourcecode |
| `mcp__jadx__get_smali_of_class` | `class_name` | Getsome/certain category/class smali |
| `mcp__jadx__get_methods_of_class` | `class_name` | columnmethod |
| `mcp__jadx__get_fields_of_class` | `class_name` | columnword paragraph/segment |
| `mcp__jadx__get_method_by_name` | `class_name`,`method_name` | take/getsome/certainmethodSourcecode |
| `mcp__jadx__get_selected_text` | no/without | Getwhen before/front选middle/center文word |

### 13.4 resourceSourceand/withstring

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__jadx__get_all_resource_file_names` | `count?`,`offset?` | columnresourceSourceFile |
| `mcp__jadx__get_resource_file` | `resource_name` | 读resourceSourceFilecontent |
| `mcp__jadx__get_strings` | `count?`,`offset?` | Get strings.xml content |

### 13.5 Searchand/withcitation

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__jadx__search_classes_by_keyword` | `search_term`,`package?`,`search_in?`,`offset?`,`count?` | 跨codeSearch category/class/method/word paragraph/segment/codecontent |
| `mcp__jadx__search_method_by_name` | `method_name` | 搜method name |
| `mcp__jadx__get_xrefs_to_class` | `class_name`,`count?`,`offset?` | 查 category/classcitation |
| `mcp__jadx__get_xrefs_to_field` | `class_name`,`field_name`,`count?`,`offset?` | 查word paragraph/segmentcitation |
| `mcp__jadx__get_xrefs_to_method` | `class_name`,`method_name`,`count?`,`offset?` | 查methodcitation |

### 13.6  re-/heavy命 name

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__jadx__rename_class` | `class_name`,`new_name` |  re-/heavy命 name category/class |
| `mcp__jadx__rename_field` | `class_name`,`field_name`,`new_name` |  re-/heavy命 nameword paragraph/segment |
| `mcp__jadx__rename_method` | `method_name`,`new_name` |  re-/heavy命 namemethod |
| `mcp__jadx__rename_variable` | `class_name`,`method_name`,`variable_name`,`new_name`,`reg?`,`ssa?` |  re-/heavy命 namevariable |
| `mcp__jadx__rename_package` | `old_package_name`,`new_package_name` |  re-/heavy命 namePackage |

### 13.7 Debug相 close/shut

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__jadx__debug_get_threads` | no/without | viewDebugThread |
| `mcp__jadx__debug_get_stack_frames` | no/without | viewwhen before/frontcall/invokeStack |
| `mcp__jadx__debug_get_variables` | no/without | viewwhen before/frontvariable |

### 13.8 RecommendationWorkflow

#### APK 初步Analysis

1. `get_android_manifest`
2. `get_main_activity_class`
3. `get_manifest_component`
4. `search_classes_by_keyword`
5. `get_class_source`

#### Signature/interfacedefinebit

1. `search_classes_by_keyword` 搜 `okhttp`, `retrofit`, `sign`, `token`, `encrypt`
2. `get_xrefs_to_method`
3. `get_method_by_name`
4. 必need totime切 to `frida_mcp` dynamicValidate

### 13.9 Notepoint

- `search_classes_by_keyword` is `jadx` innon-oftenhighpricevalue's/ofenter口tool
- `search_in` can指define `class,method,field,code,comment`
-  for/to JNI scenario，usually `jadx` 找 native registerpoint，`ida_pro_mcp` deep挖 so

---

## 14. `js_reverse`：Web Frontend JavaScript Reverseand/withbreak/judgepointDebug

### 14.1 definebit

`js_reverse` is面 to/towards Web FrontendReverse's/of专业 MCP。itand `chrome_devtools` 's/ofdifferencepart：

- `chrome_devtools` 更偏pageoperation、network、Snapshot、property/naturecan
- `js_reverse` 更偏 JS Sourcecode、break/judgepoint、call/invokechain、XHR send/issuestart者、functionTrace、Sourcecodesave

适 usescenario：

- AnalysisSignaturefunction
- trace XHR/Fetch send/issuestartchain
- definebitObfuscationfunction
- Search JS Sourcecodemiddle/center's/of close/shutkey word
- at/inExecutecontextmiddle/centertake/getvariable
- Analysis WebSocket messagepattern

### 14.2 pageand/withcontext

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__new_page` | `url`,`timeout?` | newbuildpage |
| `mcp__js_reverse__select_page` | `pageIdx?` | listorswitchpage |
| `mcp__js_reverse__navigate_page` | `type`,`url?`,`timeout?`,`ignoreCache?` | guide航/Refresh |
| `mcp__js_reverse__select_frame` | `frameIdx?` | listorswitch frame/iframe |

### 14.3 footthisEnumerationand/withSourcecodeRead

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__list_scripts` | `filter?` | listwhen before/frontpagefootthis |
| `mcp__js_reverse__search_in_sources` | `query`,`isRegex?`,`caseSensitive?`,`excludeMinified?`,`urlFilter?`,`maxResults?`,`maxLineLength?` | at/inallfootthismiddle/centerSearch |
| `mcp__js_reverse__get_script_source` | `url?`,`scriptId?`,`startLine?`,`endLine?`,`offset?`,`length?` | Readsmall (classifier) paragraph/segmentSourcecode |
| `mcp__js_reverse__save_script_source` | `filePath`,`url?`,`scriptId?` | save completewhole/integerfootthis toLocal |

explanation：

- `get_script_source` set upplan become/success“look/seegame part”，notis拉whole/integer (counter)File
- largefootthisshoulduse `save_script_source`

### 14.4 break/judgepoint、traceand/withExecutecontrol

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__set_breakpoint_on_text` | `text`,`urlFilter?`,`occurrence?`,`condition?` | press/according tocode文thisAutomaticdescendbreak/judgepoint |
| `mcp__js_reverse__list_breakpoints` | no/without | columnbreak/judgepoint |
| `mcp__js_reverse__remove_breakpoint` | `breakpointId?`,`url?` | Deletebreak/judgepointor XHR break/judgepoint |
| `mcp__js_reverse__pause_or_resume` | no/without | PauseorcontinueExecute |
| `mcp__js_reverse__step` | `direction` | single步 over/into/out |
| `mcp__js_reverse__trace_function` | `functionName`,`logArgs?`,`logThis?`,`pause?`,`traceId?`,`urlFilter?` | Tracefunctioncall/invoke |
| `mcp__js_reverse__inject_before_load` | `script?`,`identifier?` | pageLoad before/frontInjectfootthis |

### 14.5 break/judgepointHit after/back's/ofcontextAnalysis

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__get_paused_info` | `frameIndex?`,`includeScopes?`,`maxScopeDepth?` | Getbreak/judgepointHittime's/ofStackand/witheffect/functiondomainvariable |
| `mcp__js_reverse__evaluate_script` | `function`,`frameIndex?`,`mainWorld?` | at/inwhen before/frontpageorbreak/judgepoint帧middle/centerExecute JS |

### 14.6 networkand/withcall/invokechain

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__break_on_xhr` | `url` |  for/toincludes/containsgoal/target URL 's/of XHR/Fetch settingbreak/judgepoint |
| `mcp__js_reverse__list_network_requests` | `reqid?`,`pageIdx?`,`pageSize?`,`resourceTypes?`,`urlFilter?`,`includePreservedRequests?` | viewrequestcolumntableorsinglerequest详情 |
| `mcp__js_reverse__get_request_initiator` | `requestId` | viewsome/certainrequest by/from哪 paragraph/segment JS send/issuestart |
| `mcp__js_reverse__list_console_messages` | `msgid?`,`pageIdx?`,`pageSize?`,`types?`,`includePreservedMessages?` | viewConsole |

### 14.7 WebSocket Analysis

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__get_websocket_messages` | `wsid?`,`analyze?`,`groupId?`,`frameIndex?`,`direction?`,`show_content?`,`pageIdx?`,`pageSize?`,`urlFilter?`,`includePreservedConnections?` | column WS Connection、AnalysismessageBlock/Group、look/see具body帧 |

### 14.8 截Graph

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__js_reverse__take_screenshot` | `filePath?`,`format?`,`fullPage?`,`quality?` | 截Graph |

### 14.9 RecommendationWorkflow

#### definebitSignaturefunction

1. `new_page`
2. `list_scripts`
3. `search_in_sources` 搜 `sign` / `token` / Path close/shutkeyword
4. `set_breakpoint_on_text`
5. triggerrequest
6. `get_paused_info`
7. `step`
8. `evaluate_script`

#### Tracerequestis谁send/issuestart's/of

1. operationpage
2. `list_network_requests`
3. `get_request_initiator`
4. 必need totime `break_on_xhr`

#### AnalysisObfuscationfootthis

1. `search_in_sources`
2. `save_script_source`
3. `set_breakpoint_on_text`
4. `trace_function`

### 14.10 skill editwriteRecommendation

- has/haveSourcecode close/shutkey wordtime，advantage first `search_in_sources`
- has/haverequest URL time，advantage first `break_on_xhr` or `get_request_initiator`
- needat/inpagefootthiseffect/functiondomainintakeall/fullgamevariabletime，cantestconsider `mainWorld: true`
- like/such as resultpage re-/heavyload频complex，advantage firstpress/according to URL 查footthis，do not past/excessivedegree/measuredepend ontemporary `scriptId`

---

## 15. `memory`：structure-izeknowknowGraph谱remember忆

### 15.1 definebit

`memory` isgrow期structure-izeremember忆layer，notisregular/normalnotes。itmaintain's/ofis“solidbody-observe- close/shut system/relationship”'s/ofknowknowGraph谱。

suitable forused to：

- Log/Recorduser偏good
- Log/Recorditemeye/look事solid
- Log/Recordset upbackup、goal/target、Package name、interface name、vulnerabilitypointetc.structure-izeknowknow
- at/inmulti/multipleroundTaskbetweensave稳define事solid

### 15.2 core for/to象

- solidbody `entity`
  - has/have nameword `name`
  - has/havetype `entityType`
  - has/havemulti/multiple (classifier)observe `observations`

-  close/shut system/relationship `relation`
  - `from`
  - `relationType`
  - `to`

### 15.3 methodclearsingle

| tool | mainneed toparameter | effect/function |
| --- | --- | --- |
| `mcp__memory__read_graph` | no/without | Readwhole/integer (counter)Graph谱 |
| `mcp__memory__search_nodes` | `query` | 搜solidbody/type/observe |
| `mcp__memory__open_nodes` | `names` | 打 open指definesolidbody详情 |
| `mcp__memory__create_entities` | `entities` | BatchCreatesolidbody |
| `mcp__memory__delete_entities` | `entityNames` | Deletesolidbody |
| `mcp__memory__add_observations` | `observations` |  to/forsolidbodyappendobserve |
| `mcp__memory__delete_observations` | `deletions` | Deleteobserve |
| `mcp__memory__create_relations` | `relations` | Create close/shut system/relationship |
| `mcp__memory__delete_relations` | `relations` | Delete close/shut system/relationship |

### 15.4 example

Createsolidbody：

```json
{
  "entities": [
    {
      "name": "com.example.app",
      "entityType": "android_app",
      "observations": [
        "mainPackage name",
        "use OkHttp"
      ]
    }
  ]
}
```

Create close/shut system/relationship：

```json
{
  "relations": [
    {
      "from": "com.example.app",
      "relationType": "uses",
      "to": "OkHttp"
    }
  ]
}
```

### 15.5 suitable for skill 's/of use途

- at/inReverse skill middle/centerrememberlive/staygoal/targetPackage name、Encryption category/class、so  name、 close/shutkeyinterface
- at/inpenetrationTest skill middle/centerrememberlive/stayDomain Name、vulnerabilitypoint、Scanningresult/outcome
- at/inAutomatic-ize skill middle/centerrememberlive/stayaccountenvironment、deploymentway/manner、about/approximatelydefinePath

### 15.6 Notepoint

-  close/shut system/relationshipRecommendation useActive language态，for example `App uses OkHttp`
- notsuitable forexistsupergrowOriginal Text，更suitable forexist“can检索事solid”

---

## 16. `sequential_thinking`： part/point步thinktestsupplementary

### 16.1 definebit

这isone (counter)“show/display style/modemulti/multiple步thinktest”tool，used for/for repeatmixedissue/problemAnalysis、modification、Branch、Validateassumption。  
itsuitable for do：

- multi/multiplestepReverseAnalysis规plan
- notdeterminesTask's/ofsolutionexplore
- needmodification before/front面judgebreak/judge's/of repeatmixeddecidestrategy
- largeTask part/pointuntie/solve

### 16.2 method

#### `mcp__sequential_thinking__sequentialthinking`

mainneed toparameter：

- `thought`
- `thoughtNumber`
- `totalThoughts`
- `nextThoughtNeeded`
- `isRevision?`
- `revisesThought?`
- `branchFromThought?`
- `branchId?`
- `needsMoreThoughts?`

### 16.3 useway/manner principle/logicuntie/solve

thistoolnotisused to“查data”'s/of，而isused to (object marker)推 principle/logicstatestructure-ize (adverbial)Commit to/forsystem。  
youcan：

-  fromNo. 1 步 open startAnalysis
- discover before/front面wrong(past tense)then revision
-  fromsome/certainone步 part/point叉 branch
- finally形 become/successone (counter)经 past/excessiveValidate's/ofuntie/solve method/law

### 16.4 suitable for skill 's/ofscenario

- Automatic triage skill
- multi/multiplephase/stageVulnerability Exploitation路线judgebreak/judge
- Reversemiddle/center“ first Java or first native”'s/ofdecidestrategy
- multi/multiplecandidateSignaturefunctionscreen

### 16.5 example

```json
{
  "thought": " firstAcknowledgmentissue/problemisFrontendSignatureorServiceend(side)validateleads to 403。",
  "thoughtNumber": 1,
  "totalThoughts": 4,
  "nextThoughtNeeded": true
}
```

### 16.6 Notepoint

- 这isAnalysisenhancementdevice，notisExecutedevice
-  for/tosimplesingleTask没必need touse
-  for/to repeatmixed、fuzzy/blur、容easywalkwrong路's/ofissue/problem尤its/theirhas/havepricevalue

---

## 17. `scrcpy_vision`：Android canlook-izecontrol、UI definebitand/withno/without线Debug

### 17.1 definebit

`scrcpy_vision`  (object marker) ADB、scrcpy lowlatencycontrol、屏幕截Graph/串Stream、`uiautomator` UI TreeReadintegration toonegroup/settoolin，suitable for do：

-  with/by `serial` for/iscore's/of Android set upbackupConnectionand/withidentify
- based onwhen before/frontpageelement文this、`resource-id`、`content-desc` 's/of UI definebit
- sit标point击、drag、growpress/according to、滑move、keyroundinputenter
- 屏幕唤wake/unlock、 before/front (classifier for machines) Activity、notification、剪paste板etc.stateAcknowledgment
- USB 转 WiFi ADB Debug
- single帧截Graphorcontinuous画面Stream，used for/forobserveboundary面changeandAutomatic-ize联move

and `adb_mcp` compared to，it更偏“canlook-izecontrol”and“UI layerdefinebit”；`adb_mcp` 更偏foundation/basisset upbackupmanage、Installation APK、logcat、record屏、Filetransmitinput。write skill time两者usuallyis互补 close/shut system/relationship，而notistwo选one。

### 17.2 suitable for's/of skill type

- Android UI Automatic-izeand/withpageregress
- App dynamicTestmiddle/center's/ofelementdefinebitand/withboundary面Driver
- no/without线Debugswitchand/withtruemachineRemotecontrol
- grab/capturePackage/Hook  before/front after/back's/ofpagestateValidate
- needvia/through UI TreeAcknowledgmentbutton、inputenterbox、弹窗location's/ofTask
- need连续viewset upbackup画面而notis (classifier)截single (classifier)Graph's/ofTask

### 17.3 methodclearsingle

#### set upbackupConnectionand/withidentify

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__scrcpy_vision__android_devices_list` | no/without | listalreadyConnectionset upbackup | Get `serial`，Acknowledgment USB/WiFi Connectionisno/notnormal |
| `mcp__scrcpy_vision__android_devices_info` | `serial` | Readset upbackupfoundation/basis `getprop` information | look/see typenumber、systemversion、ABI、set upbackupidentifier |
| `mcp__scrcpy_vision__android_adb_enableTcpip` | `serial`,`port?` | at/in USB alreadyConnectiontimeEnable/On WiFi Debug | for/isno/without线 ADB  do before/frontplaceaccuratebackup |
| `mcp__scrcpy_vision__android_adb_getDeviceIp` | `serial` | Getset upbackup WiFi IP | accuratebackup `connectWifi` |
| `mcp__scrcpy_vision__android_adb_connectWifi` | `ipAddress`,`port?` | via/through WiFi Connectionset upbackup | no/without线Debug |
| `mcp__scrcpy_vision__android_adb_disconnectWifi` | `ipAddress?` | Disconnect指defineorall WiFi ADB Connection | Cleanupno/without线DebugSession |

#### should useand/withRun态

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__scrcpy_vision__android_app_start` | `serial`,`packageName`,`activity?` | Startshould useor指define Activity | 打 opengoal/target App、directreach指definepage |
| `mcp__scrcpy_vision__android_app_stop` | `serial`,`packageName` | mandatoryStopshould use | Resetshould usestate |
| `mcp__scrcpy_vision__android_apps_list` | `serial`,`system?` | listalreadyInstallationPackage | 找Package name、Acknowledgmentshould useisno/notInstallation |
| `mcp__scrcpy_vision__android_activity_current` | `serial` | Getwhen before/front before/front (classifier for machines)Package nameand/with Activity | judgebreak/judgewhen before/frontpageisno/notswitch become/successmerit |
| `mcp__scrcpy_vision__android_notifications_get` | `serial` | Exportwhen before/frontnotification详情 | 查Validatecodenotification、Push文案、Package namecomeSource |

#### 屏幕、剪paste板and/withset upbackupstate

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__scrcpy_vision__android_screen_isOn` | `serial` | judgebreak/judge屏幕isno/notpoint亮 | Automatic-ize before/frontInspect/Checkset upbackupstate |
| `mcp__scrcpy_vision__android_screen_wake` | `serial` | point亮屏幕 | accuratebackupoperationset upbackup |
| `mcp__scrcpy_vision__android_screen_sleep` | `serial` | 熄destroy/extinguish屏幕 | collect/receivetail/endorValidatelock屏rowfor/is |
| `mcp__scrcpy_vision__android_screen_unlock` | `serial` | attempt唤wake并unlockset upbackup | no/withoutsecuritylocktimefastspeed/fastenter桌面 |
| `mcp__scrcpy_vision__android_clipboard_get` | `serial` | Read剪paste板content | take/getValidatecode、 part/pointenjoylink、copyresult/outcome |
| `mcp__scrcpy_vision__android_clipboard_set` | `serial`,`text` | attemptsetting剪paste板 |  to/towardsinputenterboxpasteaccuratebackupgood's/of文this |

#### Fileand/with Shell

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__scrcpy_vision__android_file_list` | `serial`,`path` | listset upbackupDirectorycontent | viewExportDirectory、cacheDirectory、DownloadDirectory |
| `mcp__scrcpy_vision__android_file_pull` | `serial`,`remotePath`,`localPath` |  fromset upbackup拉File toLocal | ExportLog、Graph (classifier)、DownloadFile |
| `mcp__scrcpy_vision__android_file_push` | `serial`,`localPath`,`remotePath` | PushLocalFile toset upbackup | 推configuration、TestFile、Certificate |
| `mcp__scrcpy_vision__android_shell_exec` | `serial`,`command` | Executeanymeaning/intent `adb shell` command | at/inmusttime dohighlevel/grade诊break/judge、 part/point辨率queryorset upbackupoperation |

#### UI TreeReadand/withinputentercontrol

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__scrcpy_vision__android_ui_dump` | `serial` | Exportwhen before/frontpage's/of `uiautomator` XML | Getelement文this、 category/class name、boundary/perimeter、`resource-id` |
| `mcp__scrcpy_vision__android_ui_findElement` | `serial`,`text?`,`resourceId?`,`className?`,`contentDesc?` | press/according to UI attribute查element并returnscentersit标 | definebitbutton、inputenterbox、弹窗control (classifier) |
| `mcp__scrcpy_vision__android_input_tap` | `serial`,`x`,`y` | point击sit标 | pointbutton、columntableitem、菜single |
| `mcp__scrcpy_vision__android_input_longPress` | `serial`,`x`,`y`,`durationMs?` | growpress/according tosit标 | 呼exitcontext菜single、拖dynamicaccuratebackup |
| `mcp__scrcpy_vision__android_input_swipe` | `serial`,`x1`,`y1`,`x2`,`y2`,`durationMs?` | 滑move屏幕 | 滚movecolumntable、切页、descend拉Refresh |
| `mcp__scrcpy_vision__android_input_dragDrop` | `serial`,`startX`,`startY`,`endX`,`endY`,`durationMs?` | drag togoal/targetlocation | 拖move卡 (classifier)、Graph标、Sortitem |
| `mcp__scrcpy_vision__android_input_pinch` | `serial`,`centerX`,`centerY`,`startDistance`,`endDistance`,`durationMs?` | near似simulateshrinkrelease/put手势 |  (adverbial)Graph、Graph (classifier)shrinkrelease/putValidate |
| `mcp__scrcpy_vision__android_input_keyevent` | `serial`,`keycode` | Send Android press/according tokey | Home、Back、Enter、Delete、音quantity/measurekey |
| `mcp__scrcpy_vision__android_input_text` | `serial`,`text` | inputenter文this | login、Search、tablesingle填write |

#### lookfeelcan力

| tool | mainneed toparameter | effect/function | 典 type use途 |
| --- | --- | --- | --- |
| `mcp__scrcpy_vision__android_vision_snapshot` | `serial` | via/through `adb exec-out screencap -p` Getwhen before/front屏幕 PNG | single next/time截GraphAcknowledgmentboundary面 |
| `mcp__scrcpy_vision__android_vision_startStream` | `serial`,`frameFps?`,`maxFps?`,`maxSize?` | Start scrcpy+ffmpeg continuous画面Stream | continuousobservepagechange，with/combined withfastspeed/fastinputentercontrol |
| `mcp__scrcpy_vision__android_vision_stopStream` | `serial` | Stop画面Stream并RemoveresourceSource | collect/receivetail/end，releaseStreamresourceSource |

### 17.4 RecommendationWorkflow

#### pageAutomatic-izeand/withdefinebit

1. `android_devices_list`
2. `android_screen_isOn` / `android_screen_wake` / `android_screen_unlock`
3. like/such as result after/back续need to usesit标point击or滑move， first use `android_shell_exec` Execute `wm size` Getwhen before/front part/point辨率
4. `android_vision_snapshot` or `android_vision_startStream`
5. `android_ui_dump` or `android_ui_findElement`
6. `android_input_tap` / `android_input_text` / `android_input_swipe`
7. `android_activity_current` Acknowledgmentisno/notentergoal/targetpage
8. needcontinuousobservetimekeepstay/keep stream，tie/knotbind after/back `android_vision_stopStream`

#### WiFi ADB switch

1. USB Connectionset upbackup after/backExecute `android_adb_enableTcpip`
2. `android_adb_getDeviceIp`
3. `android_adb_connectWifi`
4. `android_devices_list` Acknowledgmentno/without线Connectionalreadyexitpresent
5. Test complete become/success after/back use `android_adb_disconnectWifi` Cleanup

### 17.5 call/invokeexample

Enable/On WiFi Debug：

```json
{
  "serial": "R58N123456A",
  "port": 5555
}
```

press/according to文this找element：

```json
{
  "serial": "R58N123456A",
  "text": "login"
}
```

Startcontinuous画面Stream：

```json
{
  "serial": "R58N123456A",
  "frameFps": 5,
  "maxSize": 1080
}
```

querywhen before/front part/point辨率：

```json
{
  "serial": "R58N123456A",
  "command": "wm size"
}
```

### 17.6 Notepoint

- divide `android_devices_list`、`android_adb_connectWifi`、`android_adb_disconnectWifi` outside，mostmethodallneed to求 firsttake toset upbackup `serial`
- like/such as result scrcpy 画面StreamalreadyStart，point击、滑move、inputenteretc.operationwill/canadvantage firstwalk更fast's/of scrcpy controlcommon道；no/not rule/principlefallback to ADB inputenter
- like/such as resultneed tosend/issuesit标point击、growpress/according to、滑move、dragor pinch， firstquerywhen before/front part/point辨率；notsame/togetherset upbackup、横竖屏、shrinkrelease/putor截Graph尺寸assumptionallcancanleads tosit标offset
- `android_ui_findElement` suitable forwhen before/frontpage's/ofstaticdefinebit，pagechange after/backRecommendation re-/heavynew `ui_dump` or re-/heavynew查element
- can use `android_ui_findElement` / `android_ui_dump` then尽quantity/measurepartdirectreceive/connectwrite死sit标； (classifier)has/haveat/inelementdefinebitnotcan靠timejustretreatreturnsit标point击
- `android_screen_unlock`  (classifier)适used for/for没has/have PIN/Password/Graph案etc.securitylock's/ofset upbackup
- `android_clipboard_set` at/in Android 10+ ascendcancanreceive tosystemlimitation，notGuaranteeplacehas/haveset upbackupallcandirectreceive/connectgenerate/live效
- `android_input_pinch` isnear似手势，notistruecorrect/positive's/ofmulti/multiplepoint触control
- `android_shell_exec`、`android_file_push` will/candirectreceive/connect改moveset upbackupenvironment，write skill timeshouldbrightcertain这ishighriskoperation
- `android_vision_startStream` produces's/ofisReal-timeresourceSource而notisfall (adverbial)File；like/such as result (classifier)issingle next/time截Graph，advantage first use `android_vision_snapshot`

---

## 18. tie/knotcombine skill editwrite's/ofRecommendationBlock/Group

for/is(past tense) after/back续write skill，更Recommendationyoupress/according to“Taskdomain”comegroup/setorganize，而notispress/according to“toolServer name”machine械splitting。

### 18.1 Android staticAnalysis skill

advantage first MCP：

- `jadx`
- `everything_search`

commonprocess：

1. 找 APK / resourceSource
2. 读 Manifest
3. 搜 close/shutkey category/class
4. 拉methodSourcecode
5. chase xref

### 18.2 Android dynamicAnalysis skill

advantage first MCP：

- `adb_mcp`
- `scrcpy_vision`
- `frida_mcp`
- `charles`

commonprocess：

1. Acknowledgmentset upbackup
2. Installationshould use
3. look情况Start scrcpy 画面StreamorRead UI Tree
4. Start Charles live capture
5. Inject hook
6. viewrequest、boundary面andLog

### 18.3 Native Reverse skill

advantage first MCP：

- `ida_pro_mcp`
- `everything_search`

commonprocess：

1. 找 so / exe
2. `survey_binary`
3. 查string/Import
4. negative/reverseCompile close/shutkeyfunction
5.  re-/heavy命 name、修type、chasedataStream

### 18.4 Web pageAutomatic-ize skill

advantage first MCP：

- `chrome_devtools`

commonprocess：

1. 打 openpage
2. GetSnapshot
3. interactivetablesingle
4. grab/capturerequest
5. 截Graphstay/keep证

### 18.5 Web JS Reverse skill

advantage first MCP：

- `js_reverse`
- `chrome_devtools`
- `burp`

commonprocess：

1. 搜Sourcecode
2.  for/torequest URL break/judgepoint
3. chasecall/invokechain
4. Exportfootthis
5. Burp Replay

### 18.6 document检索 skill

advantage first MCP：

- `context7`
- `fetch`

commonprocess：

1. `resolve_library_id`
2. `query_docs`
3. like/such as需supplementpagecontent， again use `fetch`

---

## 19. write skill timecandirectreceive/connect repeat use's/ofTip wordtemplate

descend面 to/foryouseveralsuitable fordirectreceive/connect改writeadvance skill 's/oftemplate。

### 19.1 Android Reverse skill template (classifier) paragraph/segment

```text
whenuserneed to求Analysis Android APK time：
1. 若Taskis for/toauthorized Android App  dopenetrationTest，do not firststaticAnalysis APK； firstAcknowledgmentConnectionset upbackupascendisno/notalreadyInstallationgoal/target App。
2.  firstaccuratebackup burp or charles 's/ofgrab/capturePackagecansee/meetproperty/nature， againuse scrcpy_vision 打 open App、Drivertruesolid业务point击、inputenterandguide航。
3. each close/shutkeyaction after/back， firstInspect/Check burp or charles isno/notalreadyexitpresent HTTP/HTTPS or WebSocket dataPackage，并tie/knotcombine adb_mcp viewLog、boundary面ExceptionandRuntimestate。
4. like/such as resultdataPackagealreadycansee/meet且canReplay，directreceive/connect转enter Web/API/WebSocket securityTest，press/according to“boundary面action -> dataPackage -> Web securityAnalysis”'s/of循环continue推advancenotsame/together业务meritcan。
5.  (classifier)has/haveat/ingrab/capturenot toPackage、Packageby (passive)Encryption、Plaintextnotcan (complement)、Protocol仍not透bright、cannot稳defineReplay，orExceptionclearlypoints toClientlogic阻塞time，justuse jadx Read AndroidManifest.xml、main Activity、ExportComponent，并Search okhttp/retrofit/sign/token/encrypt etc. close/shutkeyword。
6. 若 Java layer仍not够，use frida_mcp hook Java or native boundary/perimeterrecoveryPlaintext；若discover native 线索（System.loadLibrary、JNI、so File）且 Java and/with hook 仍cannotsolve， againswitch to ida_pro_mcp Analysis dump exitcome's/of so。
7. 若needcontrolset upbackup、press/according to UI elementdefinebit、observeReal-time画面or切 to WiFi Debug，use scrcpy_vision；若needInstallationshould use、record屏、logcat、foundation/basisFiletransmitinput，use adb_mcp。
```

### 19.2 Web JS Reverse skill template (classifier) paragraph/segment

```text
whenuserneed to求definebitFrontendSignature、Obfuscationfunctionorinterfacecall/invokechaintime：
1. advantage firstuse js_reverse column举footthis并 use search_in_sources Search sign/token/hash/encode/api path etc. close/shutkey word。
2. like/such as resultKnownrequest URL，advantage firstuse break_on_xhr or get_request_initiator determinessend/issuestartlocation。
3.  for/to close/shutkeyfunctionuse set_breakpoint_on_text、trace_function、get_paused_info、step and evaluate_script GetRuntimecontext。
4. 若needsave completewhole/integerfootthisused for/forofflineAnalysis，use save_script_source。
5. 若need repeatpresentorReplayrequest，with/combined with burp 's/of create_repeater_tab、send_http1_request、send_http2_request。
6. 若needpagelevel/gradeinteractiveor截Graph，with/combined with chrome_devtools。
```

### 19.3 Native BinaryAnalysis skill template (classifier) paragraph/segment

```text
whenuserneed to求AnalysisBinary、so、malicious样thisor patch pointtime：
1. 打 open IDA  after/back firstcall/invoke ida_pro_mcp.survey_binary  dototalview，do notdirectreceive/connect盲eye/look list_funcs。
2. advantage first from strings、imports、callgraph、 close/shutkeyconstant、Sensitive API enter手shrinksmall范围。
3.  for/tocan疑functionuse analyze_function / decompile / xref_query / trace_data_flow。
4. like/such as resultfunctionReadableproperty/naturedifference，use rename、set_type、declare_type、stack_frame、diff_before_after 逐步recovery language义。
5. like/such as需Modify样this，use patch / patch_asm / put_int，并at/in必need totimesave IDB。
```

---

## 20. commonNote事item汇total

### 20.1 绝 for/toPathneed to求

 with/bydescendtypetoolfrequentlyneed to求绝 for/toPath：

- `adb_mcp.take_screenshot`
- `adb_mcp.record_screen`
- `adb_mcp.pull_file` / `push_file`
- `scrcpy_vision.android_file_pull` / `android_file_push`
- `frida_mcp` 's/of `script_file_path`、`output_file`
- `js_reverse.save_script_source`
- `chrome_devtools.take_screenshot`
- `chrome_devtools.take_memory_snapshot`
- `ida_pro_mcp.open_file`

### 20.2  part/point页 category/classparameter

common part/point页/Fragmentationparameter：

- `offset`
- `count`
- `limit`
- `pageIdx`
- `pageSize`
- `start_index`
- `length`

write skill timeRecommendationshow/display style/modeexplanation：

- default firsttake/getsmallBatch样this
- 若result/outcome past/excessivemulti/multiple， againincreaselarge limit / count

### 20.3  firstdiscover， againdeepenter

verymulti/multiple MCP allhas/haveclearly's/of“discoverphase/stagetool”，do notoneascendcomethendeep挖：

- `ida_pro_mcp`: `survey_binary`
- `jadx`: `get_android_manifest` / `search_classes_by_keyword`
- `js_reverse`: `list_scripts` / `search_in_sources`
- `chrome_devtools`: `take_snapshot`
- `charles`: `query_live_capture_entries`

### 20.4 证据stay/keepexist

suitable for do证据keepstay/keep's/of MCP：

- `adb_mcp.take_screenshot`
- `adb_mcp.record_screen`
- `scrcpy_vision.android_vision_snapshot`
- `chrome_devtools.take_screenshot`
- `js_reverse.take_screenshot`
- `charles.get_traffic_entry_detail`
- `burp` historicaland/with Repeater

### 20.5 mostcommon's/ofcombination

- Android static + dynamic：`jadx` + `frida_mcp`
- Android dynamic + Streamquantity/measure：`adb_mcp` + `charles`
- Android dynamic + UI Automatic-ize：`scrcpy_vision` + `frida_mcp`
- Android grab/capturePackage + pageDriver：`scrcpy_vision` + `charles`
- Web Automatic-ize + JS Reverse：`chrome_devtools` + `js_reverse`
- Web securityReplay：`js_reverse` + `burp`
- Native static + dynamic：`ida_pro_mcp` + `frida_mcp`

---

## 21. summary

like/such as resultyou's/ofgoal/targetis“directionthen after/back续write become/success skills”，mostsolid use's/of do method/lawnotisfor/iseach MCP single独writeone (counter) skill，而ispress/according toTaskdomainsplit：

- Android staticAnalysis
- Android dynamicAnalysisand/withgrab/capturePackage
- Web Automatic-ize
- Web JS Reverse
- Native BinaryAnalysis
- document检索
- remember忆and/withTaskstatemanage

its/theirmiddle/centermostvalue (complement)advantage first围绕its/theirset upplan skill 's/of MCP is：

1. `jadx`
2. `ida_pro_mcp`
3. `js_reverse`
4. `chrome_devtools`
5. `frida_mcp`
6. `charles`
7. `adb_mcp`

like/such as result after/back面youneed to，Istillcanat/in这份documentfoundation/basisascendcontinue帮you do两 (classifier)事：

1.  againgenerateone份“suitable for skills 's/of精simple版 MCP Quick Reference Table”
2. directreceive/connect (object marker)这份documentsplit become/successmulti/multiple (counter) `SKILL.md` template骨架

## References — android-authorized-app-pentest-sop

# Android Authorized App Pentest SOP

Use this file as the default entrypoint when the task is to test an authorized Android app and you need one stable operating sequence instead of choosing between multiple Android branches manually.

## Front Rule

Do not start with APK-first reverse analysis.

For an authorized Android app pentest, the default opening order is:

1. confirm the target app is installed on the connected device
2. get `burp` or `charles` ready for traffic visibility
3. open the app with `scrcpy_vision`
4. simulate real business clicks, navigation, text input, and feature use
5. after each important action, check whether `burp` or `charles` already sees HTTP/HTTPS requests or WebSocket messages
6. if traffic is visible and replayable, move directly into `web-playbook-index.md`
7. only escalate into `jadx`, `frida_mcp`, or `ida_pro_mcp` when traffic is absent, encrypted, opaque, unreplayable, or runtime evidence points to a client-side blocker

Reverse engineering is a blocker-resolution step, not the default entrypoint.

## Operating Goal

The first goal is not "understand the APK."
The first goal is "obtain a usable request or message and prove which UI action triggered it."

The standard loop is:

`app presence -> UI action -> screenshot/log check -> packet visibility -> replay -> Web/API/WebSocket testing`

## Default Read Order

1. `01-unified-methodology.md`
2. `02-client-api-reverse-and-burp.md`
3. `android-external-url-runtime-first-workflow.md`
4. `android-ui-driven-observation-and-packet-loop.md`
5. `web-playbook-index.md`
6. `tools-reference-index.md`

Read `android-signing-and-crypto-workflow.md` only when runtime-first testing proves that Java, hook, or native recovery is actually needed.

## Default MCP Chain

1. `scrcpy_vision`
2. `burp` or `charles`
3. `adb_mcp`
4. `jadx` only after runtime packet visibility fails
5. `frida_mcp`
6. `ida_pro_mcp` only for dumped `.so` analysis or native blocker resolution

## Step 1: Confirm device and app presence

Before any testing move:

- confirm the correct Android device `serial`
- confirm the target package is installed
- confirm whether the app is already foregrounded or must be launched

Typical helpers:

- `android_devices_list`
- `android_apps_list`
- `android_activity_current`
- `android_app_start`

If the app is not present, do not start reverse work as a substitute for environment validation.

## Step 2: Prepare packet visibility first

Before driving the feature:

- ensure Burp or Charles is the active capture path
- confirm proxy and certificate assumptions are already in place
- decide whether HTTP/HTTPS, WebSocket, or both are expected

The app should be exercised only after packet visibility is ready.

## Step 3: Drive the real business flow

Use `scrcpy_vision` to:

- wake or unlock the device
- open the app
- navigate to the target feature
- input required text or credentials
- tap, swipe, back, or confirm dialogs

If coordinate-based input is required, confirm the current screen resolution first.

## Step 4: After each important action, inspect runtime evidence

Do not jump straight from UI action to reverse.

After each important action:

1. inspect the current screenshot
2. review logs with `adb_mcp`
3. check `burp` or `charles` for HTTP/HTTPS requests or WebSocket messages

Look for:

- visible error dialogs, auth blockers, certificate warnings, blank states, or redirects
- TLS, parsing, auth, WebView, okhttp, retrofit, or JNI errors in logs
- newly appeared packets, missing packets, encrypted payloads, or replayable plaintext requests

## Step 5: Branch by packet result

### Case A: Packet is visible and replayable

- do not reverse first
- move directly into `web-playbook-index.md`
- test the API, Web, or WebSocket surface from the captured request
- preserve which screen state and action produced the packet
- repeat the loop for the next business feature

### Case B: Packet is visible but encrypted or opaque

- reverse Java first with `jadx`
- locate builder, interceptor, signer, encryptor, serializer, or JNI boundary
- use `frida_mcp` when hook-based plaintext recovery is faster than deeper static work

### Case C: Packet is missing

- re-check screenshot state
- re-check logs
- verify proxy and certificate assumptions
- if runtime setup is correct and traffic is still missing, escalate into Java recovery first

## Escalation Order

When runtime-first testing is not enough:

1. `jadx`
2. `frida_mcp`
3. dump the relevant `.so`
4. `ida_pro_mcp`

Do not dump native code first when Java or hooks can answer the blocker faster.

## Exit Criteria Before Payload Testing

Do not begin mutation or exploitation work until you can explain:

- which UI action triggered the request
- whether the request is HTTP/HTTPS or WebSocket
- whether login state, cookies, tokens, nonces, timestamps, or device values matter
- whether replay works outside the app
- which fields are safe to modify

## Evidence To Keep

Keep:

- app package name and device identity
- the screen state that triggered the packet
- the triggering UI action
- relevant screenshot and log anomalies
- the captured request or message
- the reason reverse escalation was or was not necessary
- the recovery notes if `jadx`, `frida_mcp`, or `ida_pro_mcp` became necessary

## Related Files

- `02-client-api-reverse-and-burp.md`
- `android-external-url-runtime-first-workflow.md`
- `android-ui-driven-observation-and-packet-loop.md`
- `android-signing-and-crypto-workflow.md`
- `web-playbook-index.md`

## References — android-dynamic-hooking-and-replay

# Android Dynamic Hooking And Replay

Use this file only after Android static triage has narrowed the request path, or after the runtime-first Android external URL workflow has proven that screenshots, logs, and packet checks are not enough.
Do not enter this branch while Burp or Charles already has a usable replay baseline.

## Hook Order

Prefer these points in order:

1. final request object construction
2. interceptor methods
3. request execution entrypoint
4. sign or token generator
5. native boundary

## Capture For Each Hook

- class and method
- URL
- HTTP method
- headers
- body or serialized payload
- sign input tuple
- sign output or encrypted result

## Escalation Rules

- use `scrcpy_vision` when the next runtime hook or packet trigger depends on navigating app UI, entering data, or confirming the visible screen state
- use `adb_mcp` log review before deeper reverse if a runtime exception may explain the blocker
- use Frida to confirm or bridge static gaps
- keep proxy capture active throughout dynamic work so every hook result can be compared to live HTTP/HTTPS or WebSocket traffic
- treat SSL pinning bypass as a support step, not the first step

## UI-Driven Runtime Loop

When the app path is not obvious from static code alone:

1. use `scrcpy_vision` to tap, input, scroll, or navigate toward the suspected trigger
2. capture a screenshot or UI tree after each important transition
3. analyze the current state and decide the next test action before acting again
4. keep packet capture ready so the UI trigger can be tied to one or more concrete requests
5. only then place hooks or move to replay if the relevant packet and runtime values are real
6. if packets are still encrypted or absent, reverse Java first and escalate to native only when Java no longer answers the blocker

Detailed reference: `references/android-ui-driven-observation-and-packet-loop.md`

## Replay Goal

Dynamic work is complete when you can produce:

- a stable replay recipe
- the mandatory runtime inputs
- a clear answer about which fields are safe to mutate in Burp

If a stable replay recipe already exists from captured traffic alone, skip this file and move straight into network-layer testing.

## References — android-external-url-runtime-first-workflow

# Android External URL Runtime-First Workflow

Use this file when you are testing an Android app feature that reaches an external URL or remote API and you do not yet know whether reverse engineering is necessary.

This branch is packet-first and runtime-first, not reverse-first.

## Front Rule For Authorized Android App Pentest

For an authorized Android app pentest, the opening move is not APK analysis.
Before `jadx`, `frida_mcp`, or `ida_pro_mcp`, always do this first:

1. confirm the target app is installed on the connected device
2. get `burp` or `charles` ready so the next request or message is observable
3. open the app with `scrcpy_vision`
4. simulate real business clicks and navigation
5. after each important action, inspect screenshots, logs, and `burp` or `charles`
6. if packets are visible and usable, move directly into `web-playbook-index.md`
7. only if traffic is missing, encrypted, opaque, still not replayable, or runtime evidence clearly points to a client-side blocker should reverse work begin

This rule is the default for Android pentest work. Reverse is not the first step unless the task itself is already a known decryption or reverse-only problem.

## Core Rule

Do not start by reversing the interface.

First:

1. confirm the target app is installed on a connected device
2. prepare `burp` or `charles` so the next request or message can be observed
3. drive the app with `scrcpy_vision`
4. inspect the screenshot for visible anomalies
5. review logs with `adb_mcp`
6. check whether `burp` or `charles` already receives HTTP/HTTPS requests or WebSocket messages

Only if packets are encrypted, absent, still opaque, or still unusable for replay should you escalate into reverse engineering.

## When To Use This File

- the goal is to test an Android app's external URL or API behavior
- you are still in black-box or gray-box mode and want to know whether reverse work is necessary
- the request may already be visible once the right screen action is found
- you need a disciplined way to decide when to escalate from runtime observation into Java, native, or hook-based recovery

## Primary MCP Chain

1. `scrcpy_vision`
2. `burp` or `charles`
3. `adb_mcp`
4. `jadx` only when runtime visibility is insufficient
5. `frida_mcp`
6. `ida_pro_mcp` for dumped `.so` analysis

## Runtime-First Loop

### Step 1: Confirm device and app presence

Use `scrcpy_vision` to:

- list connected devices and confirm the right `serial`
- verify the target app package is installed on the device
- confirm whether the app is already foregrounded or needs to be launched

Typical helpers:

- `android_devices_list`
- `android_apps_list`
- `android_activity_current`
- `android_app_start`

Do not jump to `jadx`, `frida_mcp`, or `ida_pro_mcp` before you have confirmed that the target app is actually present and launchable on the test device.

### Step 2: Prepare packet visibility first

Before driving the target feature:

- ensure Burp or Charles is already the active capture path
- confirm proxy and certificate assumptions are in place
- decide whether the next trigger should produce HTTP/HTTPS, WebSocket, or both

Do not start business-flow driving until the next request can actually be observed.

### Step 3: Drive the app to the target feature

Use `scrcpy_vision` to:

- wake or unlock the device
- get the current screen resolution before any coordinate-based action
- start the app
- tap into the target feature
- input text, swipe, or navigate until the external URL should be triggered

Typical helpers:

- `android_devices_list`
- `android_screen_wake`
- `android_screen_unlock`
- `android_shell_exec` with `wm size`
- `android_app_start`
- `android_input_tap`
- `android_input_text`
- `android_input_swipe`

If you are about to use coordinates instead of UI element lookup, first query the current resolution with `android_shell_exec` and `wm size`.
This prevents desktop or app clicks from drifting when device resolution, orientation, scaling, or screenshot size differs from your assumption.

Before triggering the business flow, ensure Burp or Charles is already in a state where the next request can be observed.

### Step 4: Inspect the screenshot before reversing

After each important action, take a visual checkpoint:

- use `android_vision_snapshot`
- use `android_ui_dump` when UI structure matters

Check whether the screenshot already shows something abnormal:

- visible error dialog or warning
- login or permission blocker
- white screen, crash, spinner loop, or timeout
- certificate warning or network error
- redirect to an unexpected page, host, or WebView destination

Do not reverse just because the feature failed once. First determine whether the failure is already explained by the visible state.

### Step 5: Review logs for cheap evidence

Use `adb_mcp` log review after important actions:

- check for TLS failures
- serialization or parsing errors
- auth failures or token expiry
- WebView, okhttp, retrofit, or custom network stack exceptions
- crash traces or JNI load failures

If logs already explain the issue, fix the test path first instead of escalating into reverse work.

### Step 6: Check Burp and Charles

Now decide whether traffic is already visible:

- inspect `burp` history if Burp is the active proxy
- inspect `charles` if Charles is the active capture path
- confirm whether the HTTP/HTTPS request or WebSocket message exists, whether the body or frames are plaintext, and whether replay looks realistic

Three cases:

1. packet is visible and usable
2. packet is visible but encrypted or still opaque
3. packet is missing entirely

### Step 7: Branch by packet visibility

#### Case 1: Packet is visible and usable

- do not reverse first
- move directly into replay and security testing
- use `burp` as the testing baseline
- preserve the screen action that produced the packet
- continue into `web-playbook-index.md` to test the HTTP/HTTPS or WebSocket surface
- after finishing one server-side probe set, return to the app and repeat the loop for the next business action if needed

#### Case 2: Packet is visible but encrypted or opaque

- reverse Java first with `jadx`
- locate URL builders, interceptors, signers, encryptors, and serialization logic
- if Java is insufficient, use `frida_mcp` to hook the relevant Java or native boundary and recover plaintext or arguments

#### Case 3: Packet is missing entirely

- re-check screenshot state and logs
- verify proxy and certificate assumptions
- if the app path is correct but traffic is still hidden, reverse Java first
- if Java points into native code or still does not explain the missing traffic, dump the relevant `.so`

## Escalation Order

When runtime-first visibility is not enough, escalate in this order:

1. Java reverse with `jadx`
2. Java or native hook recovery with `frida_mcp`
3. dump the relevant `.so`
4. analyze the dumped `.so` with `ida_pro_mcp`

Native work is a blocker-resolution step, not the default starting point.

## Reverse Objectives

Reverse only until one of these goals is met:

- plaintext request data is recovered
- the HTTP/HTTPS request or WebSocket message becomes visible in `burp` or `charles`
- the encryption or signer boundary is understood well enough for replay
- hook-based decryption or argument capture makes the interface testable

## Handoff To Pentest

Move into pentest only after at least one of these is true:

- Burp or Charles already has a usable baseline request
- Frida hooks recover plaintext inputs or outputs reliably
- Java or native reverse has exposed the exact blocker and replay path

If Burp or Charles already has a usable baseline request, that is the preferred handoff condition.
Do not keep reversing only because static recovery also seems possible.

Then continue into:

- `web-playbook-index.md` for API and Web testing
- `04-ai-and-mcp-security-integrated.md` if the target request reaches AI, agent, or MCP surfaces
- `tools-reference-index.md` when you need the next operator tool family

## Evidence Contract

Keep:

- the screen state that triggered the external URL
- screenshot anomalies that influenced the next step
- relevant log anomalies
- whether `burp` or `charles` saw traffic
- the reason reverse escalation was or was not necessary
- Java findings, hook points, or dumped `.so` evidence when escalation happened

## Anti-Patterns

- do not open `jadx` or `ida_pro_mcp` before confirming the target app is installed on the connected device and attempting runtime packet capture
- do not reverse the app before checking screenshot, logs, and HTTP/HTTPS or WebSocket visibility
- do not dump `.so` first when Java or hooks might solve the blocker faster
- do not move into payload testing before the request is reproducible or plaintext is recoverable
- do not send coordinate clicks or swipes before confirming the current screen resolution

## References — android-native-signature-analysis

# Android Native Signature Analysis

Use this file when Android sign or crypto logic crosses from Java into JNI or `.so`.
Only enter this branch after runtime packet checks, Java triage, or hooks show that native analysis is actually needed.

## Owns

- Java-to-native boundary proof
- SO identification
- JNI style classification
- native sign-input and sign-output assessment
- decision on whether deeper native reversing is justified

## First Pass

Prove:

- which Java method declares `native`
- which `System.loadLibrary` or `System.load` call loads the target library
- whether JNI is static export or dynamic registration
- which parameters cross the boundary
- whether the return value is the final sign or an intermediate token

## Do Not Escalate Yet When

- Java still exposes the needed request values
- replay can reuse the app or hook point
- the user does not need offline execution

## Escalate Further Only When

- offline generation is required
- deeper algorithm recovery is required
- unidbg or SO-level execution is explicitly needed

## Output

- Java entrypoint
- SO name
- JNI style
- input tuple
- output role
- recommended next step

## References — android-network-layer-testing-quick-reference

# Android Network-Layer Testing Quick Reference

Use this file when the Android request is already visible or close to visible and you want one short operator card for network-layer testing instead of switching repeatedly between the Android runtime docs and the Web testing docs.

## Core Rule

For an authorized Android app pentest, network-layer testing starts as soon as one real HTTP/HTTPS request or WebSocket message is reproducible outside the app.

Do not keep reversing just because deeper recovery is possible.
If Burp or Charles already has a usable baseline request, switch to server-side testing first and only return to reverse when replay, plaintext, or state recovery stalls.

## Minimal Entry Conditions

Before payload mutation or exploitation work, confirm:

- the target app is installed and the triggering business flow is known
- the triggering UI action is known
- Burp or Charles has captured the real request or message
- the request can be replayed outside the app at least once
- required cookies, headers, tokens, timestamps, nonces, device values, or sequence prerequisites are noted
- you know which fields are safe to change first

If these conditions are not met, go back to:

- `android-authorized-app-pentest-sop.md`
- `android-external-url-runtime-first-workflow.md`
- `android-ui-driven-observation-and-packet-loop.md`
- `02-client-api-reverse-and-burp.md`

## Network-Layer Loop

Use this loop for each business feature:

1. capture one clean baseline request or WebSocket message
2. replay it unchanged in Burp to prove the baseline is stable
3. classify the surface: REST, GraphQL, WebSocket, file upload, auth flow, payment flow, admin/API gateway, or mixed
4. mutate the smallest safe field first
5. compare status code, body, timing, side effects, and server state
6. preserve evidence and note whether the change was accepted, normalized, rejected, or blocked by signer or sequencing logic
7. if the baseline breaks, stop fuzzing and restore replay before continuing

The operating sequence is:

`baseline capture -> stable replay -> small mutation -> compare response and side effect -> expand by bug class`

## What To Test First

### Auth and session

- remove or swap tokens, cookies, device identifiers, and tenant or user identifiers
- replay requests across users, roles, and sessions
- test horizontal and vertical authorization
- test whether old tokens, stale connections, or downgraded roles still work

### Business logic

- change object IDs, amounts, quantities, prices, discounts, coupon state, or workflow steps
- skip prerequisite steps
- replay requests out of order
- repeat the same request to look for race or double-spend style behavior

### Input and injection

- test the request fields that cross trust boundaries into query, render, parse, template, file, or command contexts
- prioritize fields that reach search, filter, sort, rich text, file metadata, XML, or server-side fetch behavior

### Protocol-specific behavior

- for GraphQL, test introspection, field overreach, nested object access, and resolver auth
- for WebSocket, test message auth, room or channel access, stale authorization, and message tampering
- for upload flows, test file type checks, metadata trust, parser reachability, and storage-path exposure

## Safe Mutation Order

Start from the lowest-risk mutations first:

1. duplicate the baseline unchanged
2. remove optional-looking parameters
3. modify one non-crypto business field
4. modify one identity or authorization field
5. modify one sequencing field such as nonce, timestamp, cursor, or step token
6. only then test larger payload families

Do not change many fields at once.
If the request is signed or stateful, multi-field changes hide the real blocker.

## Stop Conditions

Stop network-layer mutation and return to recovery when:

- the baseline request no longer replays consistently
- every mutation fails because a hidden signer, serializer, or state transition is missing
- the payload is still encrypted or opaque
- the response behavior suggests the app is adding unseen runtime values
- WebSocket frames or HTTP bodies are not plaintext enough to reason about safely

Escalation order:

1. Java recovery with `jadx`
2. runtime hook recovery with `frida_mcp`
3. `.so` dump and `ida_pro_mcp` only when Java and hooks still do not answer the blocker

## Best Follow-On References

- `web-playbook-index.md` for server-side bug classes and payload families
- `web-modern-protocols.md` for CORS, GraphQL, WebSocket, OAuth/OIDC, and request smuggling
- `web-logic-auth.md` for IDOR, auth bypass, reset flows, payment logic, and workflow abuse
- `web-file-infra.md` for upload, traversal, inclusion, and infrastructure issues

## Evidence To Keep

For each tested feature, keep:

- the screen state and UI action that produced the baseline request
- one clean baseline request or message
- the first successful replay outside the app
- the exact mutated field and observed difference
- whether the issue is auth, logic, injection, protocol, file, or infrastructure related
- whether reverse recovery was needed again and why

## References — android-signature-reverse-template

# Android Signature Reverse Template

Use this template for Android sign, token, encrypt, decrypt, JNI, interceptor, and replay tasks.

## Template

```markdown
# Android SignatureReverseLog/Record

## 基thisinformation

- APK / Package name：
- goal/targetmeritcan：
- goal/targetrequest：
- goal/targetword paragraph/segment：
- when before/frontphase/stage：static / dynamic / native / replay
- when before/frontstate：🟡 advancerowmiddle/center / ✅ alreadyclose环 / ⛔ 阻塞
- goal/target：
- constraint：

## statictotalview

| itemeye/look | content |
| --- | --- |
| Manifest enter口 |  |
| Application |  |
| main Activity / goal/targetComponent |  |
| mainneed toPackagestructure |  |
| networkFramework |  |
| DI Framework |  |
| when before/frontconclusion |  |

## requestcall/invokechain

```text
Activity / Fragment / Service
-> ViewModel / Presenter / UseCase
-> Repository / DataSource
-> ApiService / RequestBuilder / Interceptor
-> Signer / Encryptor / Serializer
```

- truesolidcall/invokechain：
- request Method / Path：
- Header Writepoint：
- Body Writepoint：
- Sign inputenter汇combinepoint：
- 序column /  before/frontplacedepend on：

## Sign / Crypto definebit

| itemeye/look | content |
| --- | --- |
| Sign  category/class / method |  |
| Encrypt  category/class / method |  |
|  close/shutkeyconstant |  |
|  close/shutkey Header |  |
|  close/shutkey Token / Device value |  |
| Java-only / Java+JNI / Native-first |  |

## dynamicValidate

| Hook point | cause | Capturecontent | result/outcome |
| --- | --- | --- | --- |
| Hook1 |  |  |  |

- URL：
- Headers：
- Body：
- Sign inputenter：
- Sign inputexit：
- ProxyValidate：

## JNI / SO Analysis

| itemeye/look | content |
| --- | --- |
| Java native enter口 |  |
| SO name |  |
| JNI type | static / dynamic |
| inputenterparameter |  |
| inputexitrole | final sign / middle token / other |
| isno/notneed deeper RE |  |

## Burp Replaybaseline

- Method：
- Path：
- Query：
- Headers：
- Body：
- mustkeepstay/keepword paragraph/segment：
- can变differentword paragraph/segment：
-  before/frontplacestate：
- isno/notneedset upbackup / Hook / App 协助：

## conclusion

- when before/frontclose环程degree/measure：
- 剩extra阻塞：
- descendone步Recommendation：
```

## Minimum Required Fields

Even in a compact record, keep:

- APK or package
- target request
- real call-flow summary
- network stack
- sign or crypto location
- Java versus JNI conclusion
- one runtime hook or explicit reason why runtime is not needed
- Burp replay baseline or explicit blocker

## References — android-signing-and-crypto-workflow

# Android Signing And Crypto Workflow

Use this file when the target request is produced in an Android app and the main task is to recover sign, token, encrypt, decrypt, JNI, or request-sequencing logic so the request can be explained or replayed outside the APK.
This is not the default entrypoint for a general authorized Android app pentest.

## Core Rule

Do not jump straight into Frida or `.so` reversing.

If the task is a general authorized Android app pentest and you do not yet know whether reverse is required, do not start here.
First confirm the app is installed on the connected device, prepare Burp or Charles, use `scrcpy_vision` to drive real business features, and check after each important action whether HTTP/HTTPS requests or WebSocket messages are already visible and usable.
Start from `references/android-authorized-app-pentest-sop.md` or `references/android-external-url-runtime-first-workflow.md`, and only return here after screenshot review, logs, and packet visibility checks show that reverse is necessary.

Start with static triage in `jadx` and answer:

- which network stack is in use
- where the request is built
- where headers, body, and sign fields are written
- whether the sign or crypto path is visible in Java or handed to JNI

Use runtime work only after static evidence narrows the target.
If live traffic is already visible and replayable, prioritize network-layer testing first and use this file only to resolve the remaining signer, crypto, or sequencing blocker.

## Intake Contract

Start from this block:

```text
APK / package / target feature:
Target request / field / API path:
Trigger action:
Current symptom:
Known evidence:
Goal:
Constraints:
```

Then decide:

- is the task static triage, runtime confirmation, JNI analysis, or replay proof
- is the target request already captured or still inferred
- is the app using Java-only logic, mixed Java/JNI logic, or mostly native logic

## Static-First Workflow

### Phase 1: Entry and architecture

Read:

- `AndroidManifest.xml`
- application class
- launcher activity or target component
- package structure around `api`, `network`, `data`, `repository`, `service`, `retrofit`, `http`

Goal:

- locate entry components
- identify the app package
- identify the likely network stack and dependency injection setup

Detailed reference: `references/android-static-triage-and-callflow.md`

### Phase 2: Request-chain and call-flow proof

Trace:

```text
Activity / Fragment / Service
-> ViewModel / Presenter / UseCase
-> Repository / DataSource
-> ApiService / RequestBuilder / Interceptor
-> Signer / Encryptor / Serializer
```

Use strings, Retrofit annotations, interceptor classes, request builders, and constants as anchors.

Prove:

- request method and path
- header and body writers
- request ordering or preflight dependencies
- the exact class or method where sign inputs come together

### Phase 3: Sign and crypto locator

Search for:

- `sign`, `token`, `encrypt`, `decrypt`, `cipher`, `aes`, `rsa`, `hmac`, `md5`, `sha`
- `Interceptor`, `intercept`, `addInterceptor`
- `native`, `System.loadLibrary`, `System.load`
- hardcoded URLs, header names, key names, and device identifiers

Classify the current sign path:

- Java-only
- Java wrapper around native
- native-first
- still unknown

### Phase 4: JNI handoff triage

If Java calls native code, prove:

- which Java method declares `native`
- which library is loaded
- whether the native function is statically exported or dynamically registered
- which parameters are passed into the native boundary
- which return value comes back into the request chain

Do not start deep native reversing until the Java-side boundary is already concrete.

Detailed reference: `references/android-native-signature-analysis.md`

### Phase 5: UI-driven trigger proof

If the request depends on what screen the app is showing or which gesture submits the data, use `scrcpy_vision` after static triage has already narrowed the target.

Run this loop:

1. navigate or tap toward the suspected trigger
2. capture a screenshot or UI tree
3. analyze what screen is visible now, which controls matter, and which next action is most likely to expose the target request
4. perform the next input, tap, swipe, or back action
5. watch for the packet or state transition that proves the request path

Do not treat screenshot reasoning as a replacement for static proof. It is a runtime steering layer that helps you reach the right trigger and connect visible UI state to the request chain.

Detailed reference: `references/android-ui-driven-observation-and-packet-loop.md`

## Dynamic Escalation Rules

Escalate only when static proof is no longer enough.

### Prefer these hook points in order

1. final request object construction
2. interceptor methods
3. request execution entrypoint
4. sign or token generator
5. native boundary

For each hook, capture:

- class and method
- URL
- headers
- body or serialized payload
- sign input tuple
- sign output or encrypted result

### SSL pinning and packet capture

Treat SSL pinning bypass as a support step, not the first move.
Treat Burp or Charles as the runtime baseline that stays active so recovered signer behavior can be compared to real traffic.

Use them when:

- Java hooks still do not expose final request values
- the custom transport hides fields until after TLS setup
- you need to verify that replay matches runtime traffic

Detailed reference: `references/android-dynamic-hooking-and-replay.md`

## Native and Signature Decisions

Only escalate past Java and JNI boundary proof when the user needs:

- offline reproduction
- deeper algorithm recovery
- unidbg-based execution
- `.so` patching or native control-flow analysis

Before that, answer these questions:

- is the signature generated in Java or native code
- what exact inputs feed the signature
- which inputs are constants versus runtime values
- can replay call the app or hook the boundary instead of reimplementing the algorithm

## Android Tool Order

1. `burp` or `charles`
2. `jadx`
3. `adb_mcp`
4. `frida_mcp`
5. `ida_pro_mcp` when dumped `.so` analysis is required

The order may compress, but the logic stays the same: network visibility first, static proof second, runtime recovery third, deeper native analysis last.

## Replay Exit Criteria

Do not move into Burp mutation work until you can explain:

- where the request is built
- where sign or encryption is applied
- which runtime inputs are mandatory
- whether device identity, timestamp, nonce, token, or sequence must be preserved
- whether replay can call the app, reuse a hook point, or must reimplement the logic

If Burp or Charles already has a stable replay baseline and the remaining blocker is narrow, resolve only that blocker instead of expanding reverse scope.

## Output Contract

Deliver:

- app architecture summary
- call-flow map from entry component to request execution
- request-builder and signer location
- Java versus JNI conclusion
- runtime hook point and observed values when runtime work was needed
- Burp-ready replay recipe or the exact remaining blocker

Record template: `references/android-signature-reverse-template.md`

## Recommended Read Order Inside This Branch

1. `android-static-triage-and-callflow.md`
2. `android-dynamic-hooking-and-replay.md` only when static proof is not enough
3. `android-native-signature-analysis.md` when JNI or `.so` becomes part of the real sign path
4. `android-signature-reverse-template.md` when you need a persistent record or replay handoff

## References — android-static-triage-and-callflow

# Android Static Triage And Call Flow

Use this file first for Android request, sign, and crypto tasks after runtime-first Android pentest work has shown that network-layer testing alone is not enough.
It is not the default entrypoint for a general authorized Android app pentest.

## Owns

- manifest and entry-component reading
- package and architecture survey
- network stack identification
- call-flow tracing from UI or component to request execution
- sign-path and encrypt-path location in Java

## Static Order

1. read `AndroidManifest.xml`
2. identify application class and entry components
3. find package areas around `api`, `network`, `data`, `repository`, `service`, `retrofit`, `http`
4. identify the network framework
5. trace the request chain down to builder, interceptor, signer, encryptor, or serializer

## Common Call Flow

```text
Activity / Fragment / Service
-> ViewModel / Presenter / UseCase
-> Repository / DataSource
-> ApiService / RequestBuilder / Interceptor
-> Signer / Encryptor / Serializer
```

## Strong Anchors

- Retrofit annotations
- `Request.Builder`, `HttpUrl`, interceptor classes
- hardcoded URLs, headers, and token names
- `sign`, `token`, `encrypt`, `decrypt`, `cipher`, `sha`, `hmac`, `md5`
- `native`, `System.loadLibrary`, `System.load`

## Completion Standard

Stop static triage when you can state:

- the network stack
- the request method and path
- where headers and body are written
- where sign inputs converge
- whether the path is Java-only, mixed Java/JNI, or mostly native

## References — android-ui-driven-observation-and-packet-loop

# Android UI-Driven Observation And Packet Loop

Use this file when Android runtime progress depends on what is visible in the app, which button or field the operator should touch next, or which UI action is needed to trigger the HTTP/HTTPS request or WebSocket message that will later enter Burp and the pentest workflow.

## Core Rule

This file is a runtime steering layer.

For Android external URL testing, start here before reverse engineering: drive the app, inspect the screenshot, review logs, and check whether HTTP/HTTPS requests or WebSocket messages are already visible.
Only after those checks fail should you fall back to Java recovery, native dump, or runtime hooks.

The runtime sequence stays the same:

`app presence -> packet path ready -> UI action -> screenshot/log check -> packet capture -> replay -> pentest`

## When To Use This File

- the next request trigger is hidden behind login, navigation, wizard steps, dialogs, or feature toggles
- the target request only appears after a visible UI transition
- you need to reason from the current screenshot before deciding the next tap, text input, swipe, or back action
- you need to correlate a specific screen action with one or more packets before replay work starts
- the app is testable, but the operator still needs a disciplined observe-decide-act loop instead of blind random clicking

## Primary MCP Chain

1. `scrcpy_vision`
2. `charles` or `burp`
3. `adb_mcp`
4. `jadx` only when packets are absent, encrypted, or blocked
5. `frida_mcp`
6. `ida_pro_mcp` when dumped `.so` analysis is required

## Observe-Decide-Act Loop

### Step 1: Prepare the runtime view

- list devices and confirm the right `serial`
- confirm the target app package is installed on the device
- make sure packet capture is already ready if the next action may trigger the target request
- wake or unlock the screen if needed
- get the physical screen resolution before any coordinate-based tap or swipe
- start the app or bring it to the target feature

Typical `scrcpy_vision` helpers:

- `android_devices_list`
- `android_apps_list`
- `android_screen_wake`
- `android_screen_unlock`
- `android_shell_exec` with `wm size`
- `android_app_start`
- `android_activity_current`

If the next step will use raw coordinates, first run `android_shell_exec` with `wm size` and record the current resolution.
Do not reuse old coordinates from a different device, orientation, display mode, or screenshot scale.
Do not jump into `jadx`, `frida_mcp`, or `ida_pro_mcp` until you have confirmed the app is present and tried to trigger the target packet from the live UI.

### Step 2: Create a visual checkpoint

Capture the current state before taking the next action:

- use `android_vision_snapshot` for a single screen image
- use `android_ui_dump` when you need `resource-id`, text, class, or bounds
- use `android_ui_findElement` when you already know a likely button, text label, or content description

Do not rely on coordinates alone when the UI can still be described structurally.

### Step 3: Analyze the current screen

From the screenshot and UI tree, answer:

- what page or dialog is currently visible
- which controls are actionable now
- which field probably maps to the target request path
- what blocker is present: login, consent, captcha-like gate, empty form, step gate, cooldown, or missing prerequisite
- what next action is most likely to produce the packet you want

The output of this step should be explicit, for example:

- current screen state
- candidate next actions
- chosen next action
- why this action is the best next probe

Also decide whether the screenshot already suggests an abnormal condition:

- visible error message or warning
- auth failure or forced login
- network timeout, TLS warning, or certificate issue
- blank page, crash dialog, or repeated retry state
- redirect to an unexpected domain or WebView target

### Step 4: Execute the next UI action

Use `scrcpy_vision` to perform the chosen move:

- `android_input_tap`
- `android_input_text`
- `android_input_swipe`
- `android_input_longPress`
- `android_input_keyevent`
- `android_input_dragDrop`

Before sending any coordinate-based action such as `android_input_tap`, `android_input_swipe`, `android_input_longPress`, `android_input_dragDrop`, or `android_input_pinch`, confirm the current screen resolution first.
Prefer `android_ui_findElement` or `android_ui_dump` when possible; only fall back to raw coordinates after resolution has been confirmed.

If the action changes the screen significantly, return to Step 2 and take a new checkpoint before chaining more actions.

### Step 4.5: Review logs before reversing

After important UI transitions, check logs with `adb_mcp`:

- look for crashes, TLS errors, serialization failures, auth errors, WebView warnings, or network stack exceptions
- treat logs as a cheap discriminator before escalating into reverse work
- if logs already explain the failure or blocker, fix the test path first instead of reversing immediately

### Step 4.8: Check network evidence immediately

Before taking more UI actions or escalating into reverse:

- query `charles` or inspect `burp` for the latest HTTP/HTTPS requests or WebSocket messages
- decide whether the expected request already exists in plaintext or replayable form
- if a usable packet already exists, stop UI exploration and move to replay and pentest
- if no packet exists, return to screenshot reasoning instead of defaulting to reverse

### Step 5: Tie UI action to packet capture

After the action:

- query `charles` or inspect `burp` for the latest HTTP/HTTPS requests or WebSocket messages
- decide whether the expected packet appeared
- if it appeared, mark the triggering screen state and exact UI action that produced it
- if it did not appear, go back to screenshot analysis instead of blindly continuing

The goal is to produce a clear mapping:

`screen state -> user action -> observed packet`

### Step 6: Promote the packet into replay analysis

Once the target packet is real:

- inspect it in `charles` or `burp`
- determine host, path, method, headers, cookies, tokens, body shape, and sequencing
- if it is already usable, move directly into replay and security testing
- hand it off to `web-playbook-index.md` for API, HTTP/HTTPS, or WebSocket security analysis
- only correlate it with builder or signer logic if encryption, signatures, or replay blockers remain
- use `frida_mcp` only if runtime-only values are still missing

At this point the work leaves UI steering and enters replay and pentest mode. When one business flow has been tested, return to the app and repeat the loop for the next feature instead of defaulting to reverse engineering.
If the packet is already replayable, reverse work is optional and should not delay network-layer testing.

## Escalation Order When Packets Are Blocked

Only escalate beyond UI steering when one of these is true:

- no packet appears in `burp` or `charles`
- a packet appears but the payload is encrypted or unusable
- replay fails because mandatory plaintext values are still hidden

Escalation order:

1. reverse Java first with `jadx`
2. use `frida_mcp` to hook Java or native boundaries when hook-based plaintext recovery is faster than deeper reverse
3. if Java and hooks still do not expose enough, dump the relevant `.so`
4. analyze the dumped `.so` with `ida_pro_mcp`

The goal is not reverse for its own sake. The goal is to make HTTP/HTTPS requests or WebSocket messages visible, plaintext recoverable, or replay stable.

## Handoff To Pentest Workflow

Do not start payload mutation only because you captured a packet once.

First confirm:

- which visible UI action produced the packet
- whether the packet depends on login state, toggles, or prior screens
- whether the packet contains sign, token, nonce, timestamp, device ID, or session values that must be preserved
- whether replay is stable outside the app
- which fields are safe to change without breaking the request

Then branch:

- `web-playbook-index.md` for normal API and Web testing
- `04-ai-and-mcp-security-integrated.md` if the packet reaches AI, agent, or MCP-exposed surfaces
- `tools-reference-index.md` when you need the next operator tool family

## Evidence Contract

Keep these artifacts:

- the screen state that mattered
- the chosen next action and why it was selected
- the packet triggered by that action
- the mapping from screen action to request
- any abnormal screenshot or log evidence that justified escalation
- any runtime-only values or hook points needed for replay
- the point where the task switched from UI steering to Burp and pentest analysis

## Anti-Pattern Warnings

- do not start by randomly clicking through the app without checkpoints
- do not trust screenshot reasoning alone when logs or packet evidence can resolve uncertainty
- do not open `jadx` or `ida_pro_mcp` before confirming the target app is installed and trying to trigger the packet from the live app
- do not reverse first when screenshots, logs, and HTTP/HTTPS or WebSocket visibility can already answer the problem
- do not jump into Burp payload testing before the request is reproducible
- do not send coordinate taps or swipes before checking the current resolution, or stale coordinates may drift and hit the wrong place

## References — browser-js-signing-workflow

# Browser JS Signing Workflow

Use this file when the target request is produced in the browser and the current blocker is sign generation, token flow, cookie hops, worker or wasm indirection, anti-bot logic, or browser versus local divergence.

## Mission

Keep browser JS reverse on a staged spine:

`intake -> evidence -> locate -> recover -> runtime -> validation -> replay`

Do not pick the next step from clue words alone. Pick it from engineering state.

## Intake Contract

Start from this block:

```text
URL or target page:
Target request / field / cookie / message:
Trigger action:
Current symptom:
Known evidence:
Goal:
Constraints:
```

Then answer:

- is the target request real or still guessed
- is the write boundary proven, partial, or unknown
- is the blocker shell reduction, runtime divergence, or checkpoint proof
- what artifact must be updated next

## Evidence Rule

Do not enter stage work if the real request chain is still guessed. First capture a real sample and prove:

- the target request or message
- the trigger action
- the first dependent upstream request or response when state is involved
- whether the current sample is normal state, risk state, or still mixed

Keep a persistent request-chain record. At minimum, preserve:

- request sample
- sink or write boundary
- upstream hops
- runtime notes
- replay prerequisites

## Stage Selection

### `locate`

Enter when the request, sink, write boundary, or upstream dependency chain is still unproven.

Own these questions:

- where the target value is finally written
- which action, callback, or response triggers the write
- what upstream state feeds the write
- where normal and risk paths fork

Default boundary model:

```text
writer <- builder <- entry <- source
```

Stop when the next blocker is no longer request discovery.

Detailed reference: `references/browser-locate-and-request-chain.md`

### `recover`

Enter only after the boundary is real enough and the next blocker is shell opacity.

Typical blockers:

- webpack bootstrap
- worker bridge
- wasm loader
- dispatcher flattening
- string tables
- helper indirection
- JSVMP-style shells

Reduce only the layer that blocks progress. Stop as soon as you have a readable or callable logic contract.

Detailed reference: `references/browser-recover-and-shell-reduction.md`

### `runtime`

Enter when the boundary and shell are already clear but browser execution and local execution diverge.

Classify the first meaningful divergence before patching:

- missing object
- missing state
- anti-debugging
- unstable source
- risk branch

Use a first-divergence comparison table and keep the runtime dependency set minimal.

Detailed reference: `references/browser-runtime-fit-and-risk.md`

### `validation`

Enter when the remaining work is equivalence proof.

Compare checkpoints, not just the final output:

- request body before sign
- sign input tuple
- sign output
- encrypted payload
- header set
- cookie or storage mutation

The result must state what is proven, what is still open, and which evidence supports each claim.

Detailed reference: `references/browser-validation-and-handoff.md`

## Topic Routing Inside The Browser Branch

After the stage is selected, apply the matching topic lens:

| Current blocker | Use inside the stage |
| --- | --- |
| `sign`, `token`, dynamic headers, encrypted fields | crypto entry locating and boundary observation |
| `worker`, `wasm`, `webpack/runtime`, loader callbacks | bridge and shell reduction |
| `hasDebug`, endless `debugger`, branch flips | anti-debug and runtime diagnosis |
| `cookie` hops, WebSocket, protobuf, SSE, ack or renewal | protocol and state-chain expansion |
| `basearr`, browser/local mismatch, missing browser state | minimal environment fit |

## Browser Tool Order

1. `chrome_devtools` to capture the real request and initiator
2. `js_reverse` to trace boundary, shell, runtime, or checkpoints
3. `burp` only after one replay path is stable

## Handoff Discipline

Whenever the stage changes, output a compact handoff card:

```text
--- Stage Handoff ---
From: {previous stage}
To: {next stage}
Proven: {request, boundary, upstream chain, runtime or recovery facts}
Open: {questions the next stage must answer}
Invalidated: {stale assumptions or "none"}
```

Do not carry guesses forward as facts.

## Replay Exit Criteria

Do not move into Burp fuzzing until you can explain:

- where the target field is written
- which inputs are stable constants
- which inputs come from cookies, storage, upstream responses, or browser lifecycle
- whether request order or navigation state matters
- which fields are safe to mutate

## Output Contract

Deliver:

- current stage and why it is the correct stage
- request-chain proof
- sink or write boundary
- recovered shell or runtime conclusions when applicable
- a Burp-ready baseline request or a precise statement of the remaining blocker

Record template: `references/browser-request-chain-template.md`

## Recommended Read Order Inside This Branch

1. `browser-locate-and-request-chain.md` when the boundary is not real yet
2. `browser-recover-and-shell-reduction.md` when shell opacity is the blocker
3. `browser-runtime-fit-and-risk.md` when browser/local execution diverges
4. `browser-validation-and-handoff.md` when the remaining work is proof or stage transfer
5. `browser-request-chain-template.md` when you need a persistent record or handoff artifact

## References — browser-locate-and-request-chain

# Browser Locate And Request Chain

Use this file when the browser-side target request, sink, write boundary, or upstream state chain is still not concrete enough for shell reduction or runtime work.

## Owns

- proving the real target request from a live sample
- proving the sink or write boundary
- proving the trigger action or callback
- walking the upstream dependency chain
- separating normal-state and risk-state chains

## Boundary Model

Use this model and keep each layer distinct:

```text
writer <- builder <- entry <- source
```

- `writer`: final write into body, header, query, cookie, storage, or message envelope
- `builder`: transform, sign, encrypt, serialize, or package layer
- `entry`: UI action, callback, event, or response that starts the chain
- `source`: upstream response, storage, cookie, browser state, time, randomness, or user input

## Default Order

1. capture a real target request sample
2. observe the sink first
3. walk backward through `writer <- builder <- entry <- source`
4. expand upstream when the current source depends on prior requests or state
5. split normal-state and risk-state chains if both appear

## Strong First Observation Points

| Sink type | First point to prove |
| --- | --- |
| request body field | final serialization or submit write point |
| header field | request construction or header-set call |
| JS-written cookie | cookie setter |
| response-driven cookie dependency | response packet and first dependent request |
| WebSocket frame | final envelope before `send` |
| worker reply | `postMessage` bridge contract |

## Completion Standard

Stop locate when:

- the request sample is real
- the sink is real
- `writer`, `builder`, `entry`, and `source` are concrete enough for the next step
- the next blocker is shell opacity, runtime divergence, or checkpoint proof rather than request discovery

## Do Not Do

- broad deobfuscation before the boundary is real
- environment patching while the sink is still guessed
- relying on keyword hits as proof

## References — browser-recover-and-shell-reduction

# Browser Recover And Shell Reduction

Use this file only after the browser request boundary is already real and the next blocker is shell opacity.

## Owns

- choosing the smallest layer to open
- deciding whether the task needs semantic explanation, key-operator extraction, or a minimal rebuild
- preserving black-box reuse when deeper deobfuscation is unnecessary

## First Layer Selection

| Symptom | First layer to open |
| --- | --- |
| callable path still hidden | outer container |
| large dispatcher or VM flow | dispatcher layer |
| parameters visible but state carrier opaque | state carrier |
| logic appears after `worker` or `wasm` bridge | bridge layer |
| write-back point known but algorithm opaque | core operator |

## Recovery Levels

### Level A

Recover only the critical operator or helper needed to explain the target field.

### Level B

Recover dispatcher flow plus critical state carriers when operator meaning depends on state flow.

### Level C

Build the smallest verifiable fragment or interpreter only when levels A and B cannot support the next stage.

## Prefer Black-Box Reuse When

- input and output boundaries are already known
- the target module or bridge entry is found
- the blocker is container logic, not business logic

## Escalate Deeper When

- replay is unstable because of hidden shared state
- the bridge contract itself is opaque
- the module contains another VM or protocol shell that still blocks progress

## Completion Standard

Stop recover when the current reduction depth is already enough for runtime fit or validation.

## References — browser-request-chain-template

# Browser Request Chain Template

Use this template for browser-side sign, token, anti-bot, worker, wasm, cookie-hop, and replay tasks.

## Template

```markdown
# BrowserrequestchainLog/Record

## 基thisinformation

- goal/targetpage：
- goal/targetrequest：
- goal/targetword paragraph/segment：
- triggeraction：
- when before/frontphase/stage：locate / recover / runtime / validation
- when before/frontstate：🟡 advancerowmiddle/center / ✅ alreadyclose环 / ⛔ 阻塞
- goal/target：
- constraint：

## 样thisand/withpresent象

- normal态样this：
- 风control态样this：
- Browserpresent象：
- Localpresent象：
- when before/frontdifference：

## requestchainmaintable

| itemeye/look | content |
| --- | --- |
| writer |  |
| builder |  |
| entry |  |
| source |  |
| ascendswimdepend on |  |
| stateloadbody |  |
| 风control part/point叉point |  |
| when before/frontconclusion |  |

##  close/shutkey证据

| 证据type | location/pointbit | content | conclusion |
| --- | --- | --- | --- |
| request样this |  |  |  |
| call/invokeStack |  |  |  |
| break/judgepoint/Hook |  |  |  |
| middlevalue |  |  |  |
| Cookie/Storage |  |  |  |

## phase/stagesupplement

### Locate supplement

- Sink：
- truesolidWritepoint：
- ascendswimrequest：
- normal态 / 风control态difference part/point：

### Recover supplement

- 遮蔽layertype：
- when before/frontrecoverylevel：A / B / C
- alreadyrecovery契about/approximately：
- 仍un-recoverydisadvantage口：

### Runtime supplement

- absent for/to象：
- absentstate：
- fixedSource：
- first/head (counter) part/point歧point：
- 风control / negative/reverseDebug：

### Validation supplement

| Inspect/Checkpoint | Browserside | Local/recoveryside | result/outcome | 证据 | disadvantage口 |
| --- | --- | --- | --- | --- | --- |
| Inspect/Checkpoint1 |  |  |  |  |  |

## Burp Replaybaseline

- Method：
- Path：
- Query：
- Headers：
- Body：
- mustkeepstay/keepword paragraph/segment：
- can变differentword paragraph/segment：
-  before/frontplacestate：

## Stage Handoff

--- Stage Handoff ---
From:
To:
Proven:
Open:
Invalidated:

## descendone步

- descendoneaction：
- 预期inputexit：
- 阻塞point：
```

## Minimum Required Fields

Even in a compact record, keep:

- target page and target request
- current stage
- `writer / builder / entry / source`
- one real request sample
- one concrete evidence row
- Burp replay baseline or explicit blocker

## References — browser-runtime-fit-and-risk

# Browser Runtime Fit And Risk

Use this file when the browser boundary and shell are already clear but browser execution and local or controlled execution diverge.

## Diagnose Before Patching

Classify the first meaningful divergence as one or more of:

- missing object
- missing state
- anti-debugging
- unstable source
- risk branch

## First-Divergence Table

Always compare browser normal state and local execution using a concrete checkpoint table before adding patches.

Minimum comparison rows:

- input parameters
- cookie and storage state
- fixed time and randomness
- first stable intermediate value
- first abnormal intermediate value
- final branch or response

## Risk And Anti-Debug Refinement

When debugging changes behavior or a risk branch is suspected, answer:

- where the fork begins
- whether the issue is debug friction or a real consumer-driven risk branch
- which exact missing state or fingerprint surface triggers the split

Keep the anti-debug handling minimal. Remove only the smallest obstacle needed to keep observation going.

## Environment Fit Rules

- keep `required objects` and `required state` separate
- record why each dependency is necessary
- fix time, randomness, and seed sources before further comparison
- do not claim pure computation while upstream response, `HttpOnly` state, challenge flow, or browser lifecycle state remain open

## Completion Standard

Stop runtime when:

- the divergence class is explicit
- the first divergent checkpoint is known
- missing object and missing state are not mixed
- the next action is clearly patch, state restore, or validation

## References — browser-validation-and-handoff

# Browser Validation And Handoff

Use this file when the remaining browser-side work is checkpoint proof, equivalence proof, or stage handoff.

## Checkpoints To Compare

Do not compare only the final output. Compare:

- pre-sign payload
- sign input tuple
- sign output
- encrypted payload
- final request body
- final headers
- cookie or storage mutation when it affects later requests

## Proof Rules

Each checkpoint must state:

- fixed input sample
- browser-side value
- local or recovered-side value
- whether the checkpoint matches
- what evidence supports the comparison
- what gap remains if it does not match

## Handoff Card

When the stage changes, emit:

```text
--- Stage Handoff ---
From: {previous stage}
To: {next stage}
Proven: {request, boundary, upstream chain, recovery or runtime facts}
Open: {questions for the next stage}
Invalidated: {stale assumptions or "none"}
```

## Completion Standard

Validation is complete only when the next operator can see:

- what is equivalent
- what is not equivalent
- what evidence supports each statement
- what remains open

## References — client-reverse-workflow

# Complex Client Reverse Workflow

## Goal

Recover the real request-production chain so the interface can be reproduced outside the client.

## Stages

1. classify the client
2. choose the smallest platform branch that can prove the request chain
3. statically find request and crypto code
4. dynamically confirm signer, serializer, and state values only when static proof is no longer enough
5. rebuild the request recipe
6. replay in Burp
7. move into Web or AI attack testing only after replay is stable

## Android

- start in `jadx`
- finish manifest, package, network stack, request-builder, signer, and JNI triage first
- use `scrcpy_vision` to steer UI-dependent runtime paths when the next packet depends on what is visible on screen
- verify on-wire behavior with `adb_mcp` and `charles`
- hook signer or builder with `frida_mcp` only after the static target is narrow enough
- move to `burp`

## Native desktop

- locate files with `everything_search`
- reverse code with `ida_pro_mcp`
- capture runtime values with `frida_mcp`
- move to `burp`

## Browser JS

- inspect live requests with `chrome_devtools`
- choose the current stage from `locate`, `recover`, `runtime`, or `validation`
- trace initiators and signer functions with `js_reverse`
- replay with `burp`

## Detailed Branches

- browser JS staged flow: `browser-js-signing-workflow.md`
- Android sign and crypto flow: `android-signing-and-crypto-workflow.md`
- Android UI-driven packet trigger flow: `android-ui-driven-observation-and-packet-loop.md`

## References — reporting-and-evidence

# Reporting And Evidence

Use this file to normalize final evidence and replay handoff after browser, Android, desktop-client, Web, or AI/MCP work.

## Minimum Output

- scope and client type
- chosen MCP chain
- static findings
- runtime proof
- recovered request recipe
- Burp-ready baseline request
- security finding and mitigation

## Client-Controlled Targets

For browser or Android request-generation tasks, always include:

- target request and target field
- request-chain summary
- proven writer or sink
- upstream dependency or explicit statement that none exists
- runtime values that must be preserved
- replay-safe fields versus mutation-safe fields

## Recommended Templates

### Browser JS

- workflow: `references/browser-js-signing-workflow.md`
- persistent record: `references/browser-request-chain-template.md`

### Android

- workflow: `references/android-signing-and-crypto-workflow.md`
- persistent record: `references/android-signature-reverse-template.md`

## Final Handoff Checklist

- one real request sample is preserved
- replay prerequisites are explicit
- blockers are separated from proven facts
- the next operator can reproduce the baseline request

## References — tool-selection-map

# Tool Selection Map

## Reverse Layer

- `jadx`
- `ida_pro_mcp`
- `frida_mcp`
- `scrcpy_vision`
- `adb_mcp`
- `charles`
- `js_reverse`
- `chrome_devtools`
- `burp`

## Platform Sequences

### Browser JS sign or anti-bot

- boundary and request proof: `chrome_devtools` -> `js_reverse`
- browser/local divergence: `js_reverse`
- replay confirmation: `burp`

### Android external URL or sign/encrypt

- proxy or capture readiness first: `burp` / `charles`
- runtime-first visibility and packet check: `scrcpy_vision` -> `adb_mcp`
- Java recovery when blocked: `jadx`
- UI-state steering and screenshot-guided next actions: `scrcpy_vision`
- device state and runtime context: `adb_mcp`
- narrow Java or JNI hooks: `frida_mcp`
- dumped `.so` analysis when required: `ida_pro_mcp`
- wire validation or Charles-assisted observation: `charles`
- replay confirmation: `burp`

## Support Layer

- `everything_search`
- `context7`
- `fetch`
- `memory`
- `sequential_thinking`

## Rule

Do not start payload testing in Burp when the request is still opaque.
For Android external URL testing, do not reverse first when screenshots, logs, and packet visibility can answer the problem.
Do not choose browser references by clue words before the current stage is known.
