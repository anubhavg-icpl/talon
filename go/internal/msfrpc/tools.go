package msfrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register wires the 12 MetasploitMCP.py @mcp.tool() functions onto srv,
// backed by c.
func Register(srv *server.MCPServer, c *Client) {
	srv.AddTool(
		mcp.NewTool("list_exploits",
			mcp.WithDescription("List available Metasploit exploits, optionally filtered by search term."),
			mcp.WithString("search_term", mcp.DefaultString("")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			term := strOr(a["search_term"], "")
			exploits, err := c.ListExploits(ctx)
			if err != nil {
				return textResult([]string{fmt.Sprintf("Error: Metasploit RPC error: %v", err)}), nil
			}
			if term != "" {
				lower := strings.ToLower(term)
				filtered := make([]string, 0, len(exploits))
				for _, e := range exploits {
					if strings.Contains(strings.ToLower(e), lower) {
						filtered = append(filtered, e)
					}
				}
				return textResult(limitStrings(filtered, 200)), nil
			}
			return textResult(limitStrings(exploits, 100)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("list_payloads",
			mcp.WithDescription("List available Metasploit payloads, optionally filtered by platform and/or architecture."),
			mcp.WithString("platform", mcp.DefaultString("")),
			mcp.WithString("arch", mcp.DefaultString("")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			platform := strOr(a["platform"], "")
			arch := strOr(a["arch"], "")
			payloads, err := c.ListPayloads(ctx)
			if err != nil {
				return textResult([]string{fmt.Sprintf("Error: Metasploit RPC error: %v", err)}), nil
			}
			filtered := payloads
			if platform != "" {
				plat := strings.ToLower(platform)
				out := make([]string, 0, len(filtered))
				for _, p := range filtered {
					lp := strings.ToLower(p)
					if strings.HasPrefix(lp, plat+"/") || strings.Contains(lp, "/"+plat+"/") {
						out = append(out, p)
					}
				}
				filtered = out
			}
			if arch != "" {
				archLower := strings.ToLower(arch)
				out := make([]string, 0, len(filtered))
				for _, p := range filtered {
					lp := strings.ToLower(p)
					if strings.Contains(lp, "/"+archLower+"/") || containsSegment(lp, archLower) {
						out = append(out, p)
					}
				}
				filtered = out
			}
			return textResult(limitStrings(filtered, 100)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("generate_payload",
			mcp.WithDescription("Generate a Metasploit payload via module.execute (msfvenom-style). Saves the generated payload to a file on the server if successful."),
			mcp.WithString("payload_type", mcp.Required()),
			mcp.WithString("format_type", mcp.Required()),
			mcp.WithString("options", mcp.Required(), mcp.Description("Dict of required payload options (e.g. LHOST/LPORT) as a JSON object string, or 'LHOST=1.2.3.4,LPORT=4444'.")),
			mcp.WithString("encoder", mcp.DefaultString("")),
			mcp.WithNumber("iterations", mcp.DefaultNumber(0)),
			mcp.WithString("bad_chars", mcp.DefaultString("")),
			mcp.WithNumber("nop_sled_size", mcp.DefaultNumber(0)),
			mcp.WithString("template_path", mcp.DefaultString("")),
			mcp.WithBoolean("keep_template", mcp.DefaultBool(false)),
			mcp.WithBoolean("force_encode", mcp.DefaultBool(false)),
			mcp.WithString("output_filename", mcp.DefaultString("")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			payloadType := strOr(a["payload_type"], "")
			formatType := strOr(a["format_type"], "")

			parsedOptions, err := ParseOptionsGracefully(a["options"])
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Invalid options format: %v", err)}), nil
			}
			if len(parsedOptions) == 0 {
				return textResult(map[string]any{"status": "error", "message": "Payload 'options' dictionary (e.g., LHOST, LPORT) is required."}), nil
			}

			genOptions := map[string]any{}
			for k, v := range parsedOptions {
				genOptions[k] = v
			}
			genOptions["Format"] = formatType
			if enc := strOr(a["encoder"], ""); enc != "" {
				genOptions["Encoder"] = enc
			}
			if it := intOr(a["iterations"], 0); it != 0 {
				genOptions["Iterations"] = it
			}
			genOptions["BadChars"] = strOr(a["bad_chars"], "")
			if n := intOr(a["nop_sled_size"], 0); n != 0 {
				genOptions["NopSledSize"] = n
			}
			if tp := strOr(a["template_path"], ""); tp != "" {
				genOptions["Template"] = tp
			}
			if boolOr(a["keep_template"], false) {
				genOptions["KeepTemplateWorking"] = true
			}
			if boolOr(a["force_encode"], false) {
				genOptions["ForceEncode"] = true
			}

			result, err := c.Execute(ctx, "payload", payloadType, genOptions)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Payload generation failed: %v", err)}), nil
			}

			rawPayload, ok := payloadBytes(result["payload"])
			if !ok {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Payload generation failed: unexpected response shape: %v", result)}), nil
			}
			payloadSize := len(rawPayload)

			saveDir := c.PayloadSaveDir
			if err := os.MkdirAll(saveDir, 0o755); err != nil {
				return textResult(map[string]any{
					"status":       "error",
					"message":      fmt.Sprintf("Payload generated (%d bytes) but could not create save directory: %v", payloadSize, err),
					"payload_size": payloadSize,
					"format":       formatType,
				}), nil
			}

			filename := sanitizeFilename(strOr(a["output_filename"], ""))
			if filename == "" {
				filename = defaultPayloadFilename(payloadType, formatType)
			}
			savePath := filepath.Join(saveDir, filename)

			if err := os.WriteFile(savePath, rawPayload, 0o644); err != nil {
				return textResult(map[string]any{
					"status":       "error",
					"message":      fmt.Sprintf("Payload generated but failed to save to file: %v", err),
					"payload_size": payloadSize,
					"format":       formatType,
				}), nil
			}

			return textResult(map[string]any{
				"status":           "success",
				"message":          fmt.Sprintf("Payload '%s' generated successfully and saved.", payloadType),
				"payload_size":     payloadSize,
				"format":           formatType,
				"server_save_path": savePath,
			}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("run_exploit",
			mcp.WithDescription("Run a Metasploit exploit module with specified options as a background job, with session polling."),
			mcp.WithString("module_name", mcp.Required()),
			mcp.WithObject("options", mcp.Required()),
			mcp.WithString("payload_name", mcp.DefaultString("")),
			mcp.WithString("payload_options", mcp.DefaultString(""), mcp.Description("Dict of payload options as a JSON object string, or 'LHOST=1.2.3.4,LPORT=4444'.")),
			mcp.WithBoolean("run_as_job", mcp.DefaultBool(false)),
			mcp.WithBoolean("check_vulnerability", mcp.DefaultBool(false)),
			mcp.WithNumber("timeout_seconds", mcp.DefaultNumber(defaultSessionCommandTimeout)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			moduleName := strOr(a["module_name"], "")

			options, err := ParseOptionsGracefully(a["options"])
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Invalid options format: %v", err)}), nil
			}
			payloadOptions, err := ParseOptionsGracefully(a["payload_options"])
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Invalid payload_options format: %v", err)}), nil
			}

			var payload *payloadSpec
			if name := strOr(a["payload_name"], ""); name != "" {
				payload = &payloadSpec{Name: name, Options: payloadOptions}
			}

			// ponytail: check_vulnerability and run_as_job mirrored the Python
			// tool's console-based 'check' pre-flight and its console-vs-job
			// branch -- both needed the console fallback that was dropped (see
			// executeModuleJob). Every exploit now always runs as an RPC
			// background job via module.execute. Upgrade when console support
			// is ported and the check-then-run flow can be restored.
			result := executeModuleJob(ctx, c, "exploit", moduleName, options, payload)
			return textResult(result), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("run_post_module",
			mcp.WithDescription("Run a Metasploit post-exploitation module against a session as a background job."),
			mcp.WithString("module_name", mcp.Required()),
			mcp.WithNumber("session_id", mcp.Required()),
			mcp.WithObject("options"),
			mcp.WithBoolean("run_as_job", mcp.DefaultBool(false)),
			mcp.WithNumber("timeout_seconds", mcp.DefaultNumber(defaultSessionCommandTimeout)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			moduleName := strOr(a["module_name"], "")
			sessionID := intOr(a["session_id"], 0)

			options, err := ParseOptionsGracefully(a["options"])
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Invalid options format: %v", err)}), nil
			}
			options["SESSION"] = sessionID

			sessions, err := c.ListSessions(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error validating session %d: %v", sessionID, err), "module": moduleName}), nil
			}
			if _, ok := sessions[strconv.Itoa(sessionID)]; !ok {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Session %d not found.", sessionID), "module": moduleName}), nil
			}

			result := executeModuleJob(ctx, c, "post", moduleName, options, nil)
			return textResult(result), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("run_auxiliary_module",
			mcp.WithDescription("Run a Metasploit auxiliary module as a background job."),
			mcp.WithString("module_name", mcp.Required()),
			mcp.WithObject("options", mcp.Required()),
			mcp.WithBoolean("run_as_job", mcp.DefaultBool(false)),
			mcp.WithBoolean("check_target", mcp.DefaultBool(false)),
			mcp.WithNumber("timeout_seconds", mcp.DefaultNumber(defaultSessionCommandTimeout)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			moduleName := strOr(a["module_name"], "")
			options, err := ParseOptionsGracefully(a["options"])
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Invalid options format: %v", err)}), nil
			}
			// ponytail: check_target dropped along with the console fallback
			// -- see executeModuleJob's comment on run_exploit.
			result := executeModuleJob(ctx, c, "auxiliary", moduleName, options, nil)
			return textResult(result), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("list_active_sessions", mcp.WithDescription("List active Metasploit sessions with their details.")),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sessions, err := c.ListSessions(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Metasploit RPC error: %v", err)}), nil
			}
			return textResult(map[string]any{"status": "success", "sessions": sessions, "count": len(sessions)}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("send_session_command",
			mcp.WithDescription("Send a command to an active Metasploit session (Meterpreter or Shell) and get output."),
			mcp.WithNumber("session_id", mcp.Required()),
			mcp.WithString("command", mcp.Required()),
			mcp.WithNumber("timeout_seconds", mcp.DefaultNumber(defaultSessionCommandTimeout)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			sessionID := intOr(a["session_id"], 0)
			command := strOr(a["command"], "")
			timeout := time.Duration(intOr(a["timeout_seconds"], defaultSessionCommandTimeout)) * time.Second
			sessionIDStr := strconv.Itoa(sessionID)

			sessions, err := c.ListSessions(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error interacting with session %d: %v", sessionID, err)}), nil
			}
			infoRaw, ok := sessions[sessionIDStr]
			if !ok {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Session %d not found.", sessionID)}), nil
			}
			info, _ := infoRaw.(map[string]any)
			sessionType := strings.ToLower(strOr(info["type"], "unknown"))

			// ponytail: this drops pymetasploit3's session_shell_type submode
			// tracking (entering a nested `shell` channel from within
			// meterpreter via run_with_output(end_strs=["created."]), and
			// detach()-ing back out via `exit`) -- every command here is a
			// flat write-then-poll-read against the session's native RPC
			// methods. Upgrade when meterpreter's nested shell-channel UX
			// needs to be scripted through this tool.
			var status, message, output string
			switch sessionType {
			case "meterpreter":
				output, status, message = runMeterpreterCommand(ctx, c, sessionIDStr, command, timeout)
			case "shell":
				output, status, message = runShellCommand(ctx, c, sessionIDStr, command, timeout)
			default:
				status, message = "error", fmt.Sprintf("Cannot execute command: Unknown session type '%s'.", sessionType)
			}
			return textResult(map[string]any{"status": status, "message": message, "output": output}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("list_listeners", mcp.WithDescription("List all active Metasploit jobs, categorizing exploit/multi/handler jobs.")),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			jobs, err := c.ListJobs(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Metasploit RPC error: %v", err)}), nil
			}
			handlers := map[string]any{}
			other := map[string]any{}
			for jobID, raw := range jobs {
				info, _ := raw.(map[string]any)
				data := map[string]any{"job_id": jobID, "name": "Unknown", "details": raw}
				isHandler := false
				if info != nil {
					name, _ := info["name"].(string)
					if name != "" {
						data["name"] = name
					}
					data["start_time"] = info["start_time"]
					ds, _ := info["datastore"].(map[string]any)
					if ds != nil {
						data["datastore"] = ds
					}
					combined := strings.ToLower(name + strOr(info["info"], ""))
					if strings.Contains(combined, "exploit/multi/handler") {
						isHandler = true
					} else if ds != nil {
						_, hasPayload := ds["payload"]
						_, hasLhost := ds["lhost"]
						_, hasLport := ds["lport"]
						if hasPayload || (hasLhost && hasLport) {
							isHandler = true
						}
					}
				}
				if isHandler {
					handlers[jobID] = data
				} else {
					other[jobID] = data
				}
			}
			return textResult(map[string]any{
				"status":          "success",
				"handlers":        handlers,
				"other_jobs":      other,
				"handler_count":   len(handlers),
				"other_job_count": len(other),
				"total_job_count": len(jobs),
			}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("start_listener",
			mcp.WithDescription("Start a new Metasploit handler (exploit/multi/handler) as a background job."),
			mcp.WithString("payload_type", mcp.Required()),
			mcp.WithString("lhost", mcp.Required()),
			mcp.WithNumber("lport", mcp.Required()),
			mcp.WithString("additional_options", mcp.DefaultString(""), mcp.Description("Dict of additional payload options as a JSON object string, or 'LURI=/path,HandlerSSLCert=cert.pem'.")),
			mcp.WithBoolean("exit_on_session", mcp.DefaultBool(false)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			payloadType := strOr(a["payload_type"], "")
			lhost := strOr(a["lhost"], "")
			lport := intOr(a["lport"], 0)
			exitOnSession := boolOr(a["exit_on_session"], false)

			if lport < 1 || lport > 65535 {
				return textResult(map[string]any{"status": "error", "message": "Invalid LPORT. Must be between 1 and 65535."}), nil
			}

			additionalOptions, err := ParseOptionsGracefully(a["additional_options"])
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Invalid additional_options format: %v", err)}), nil
			}

			payloadOptions := map[string]any{}
			for k, v := range additionalOptions {
				payloadOptions[k] = v
			}
			payloadOptions["LHOST"] = lhost
			payloadOptions["LPORT"] = lport

			moduleOptions := map[string]any{"ExitOnSession": exitOnSession}
			result := executeModuleJob(ctx, c, "exploit", "multi/handler", moduleOptions, &payloadSpec{Name: payloadType, Options: payloadOptions})

			switch result["status"] {
			case "success":
				result["message"] = fmt.Sprintf("Listener for %s started as job %v on %s:%d.", payloadType, result["job_id"], lhost, lport)
			case "warning":
				result["message"] = fmt.Sprintf("Listener job %v started, but encountered issues: %v", result["job_id"], result["message"])
			default:
				result["message"] = fmt.Sprintf("Failed to start listener: %v", result["message"])
			}
			return textResult(result), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("stop_job",
			mcp.WithDescription("Stop a running Metasploit job (handler or other). Verifies disappearance."),
			mcp.WithNumber("job_id", mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			jobID := intOr(a["job_id"], 0)
			jobIDStr := strconv.Itoa(jobID)

			before, err := c.ListJobs(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error stopping job %d: %v", jobID, err)}), nil
			}
			raw, ok := before[jobIDStr]
			if !ok {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Job %d not found.", jobID)}), nil
			}
			jobName := "Unknown"
			if info, ok := raw.(map[string]any); ok {
				if n, ok := info["name"].(string); ok && n != "" {
					jobName = n
				}
			}

			apiResult, err := c.StopJob(ctx, jobIDStr)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error stopping job %d: %v", jobID, err)}), nil
			}

			sleepOrDone(ctx, 1*time.Second)

			after, err := c.ListJobs(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error stopping job %d: %v", jobID, err)}), nil
			}
			if _, stillThere := after[jobIDStr]; !stillThere {
				return textResult(map[string]any{
					"status": "success", "message": fmt.Sprintf("Successfully stopped job %d ('%s')", jobID, jobName),
					"job_id": jobID, "job_name": jobName, "api_result": apiResult,
				}), nil
			}
			return textResult(map[string]any{
				"status": "error", "message": fmt.Sprintf("Failed to stop job %d. Job may still be running.", jobID),
				"job_id": jobID, "job_name": jobName, "api_result": apiResult,
			}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("terminate_session",
			mcp.WithDescription("Forcefully terminate a Metasploit session using session.stop()."),
			mcp.WithNumber("session_id", mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			sessionID := intOr(a["session_id"], 0)
			sessionIDStr := strconv.Itoa(sessionID)

			sessions, err := c.ListSessions(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error terminating session %d: %v", sessionID, err)}), nil
			}
			if _, ok := sessions[sessionIDStr]; !ok {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Session %d not found.", sessionID)}), nil
			}

			if _, err := c.StopSession(ctx, sessionIDStr); err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error terminating session %d: %v", sessionID, err)}), nil
			}

			sleepOrDone(ctx, 1*time.Second)

			after, err := c.ListSessions(ctx)
			if err != nil {
				return textResult(map[string]any{"status": "error", "message": fmt.Sprintf("Error terminating session %d: %v", sessionID, err)}), nil
			}
			if _, stillThere := after[sessionIDStr]; !stillThere {
				return textResult(map[string]any{"status": "success", "message": fmt.Sprintf("Session %d terminated successfully.", sessionID)}), nil
			}
			return textResult(map[string]any{"status": "warning", "message": fmt.Sprintf("Session %d may not have been terminated properly.", sessionID)}), nil
		},
	)
}

// --- Shared execution helpers (mirrors _execute_module_rpc) ---

type payloadSpec struct {
	Name    string
	Options map[string]any
}

const (
	exploitSessionPollTimeout    = 60 * time.Second
	exploitSessionPollInterval   = 2 * time.Second
	sessionReadInactivityWindow  = 10 * time.Second
	sessionPollInterval          = 100 * time.Millisecond
	defaultSessionCommandTimeout = 60
)

var shellPromptRE = regexp.MustCompile(`([#$>]|%)\s*$`)

// executeModuleJob mirrors _execute_module_rpc (RPC/background-job path
// only). The Python source also has a slower console-based fallback with
// raw byte-stream prompt-regex matching for synchronous exploit/auxiliary/
// post runs -- that was dropped here; module.execute covers every exec path
// run_exploit/run_post_module/run_auxiliary_module/start_listener need.
func executeModuleJob(ctx context.Context, c *Client, modtype, modname string, moduleOptions map[string]any, payload *payloadSpec) map[string]any {
	fullOptions := make(map[string]any, len(moduleOptions)+2)
	for k, v := range moduleOptions {
		fullOptions[k] = v
	}
	var payloadName string
	if modtype == "exploit" && payload != nil {
		payloadName = payload.Name
		fullOptions["PAYLOAD"] = payload.Name
		for k, v := range payload.Options {
			fullOptions[k] = v
		}
	}

	fullModulePath := fmt.Sprintf("%s/%s", modtype, normalizeModuleName(modname))

	execResult, err := c.Execute(ctx, modtype, modname, fullOptions)
	if err != nil {
		msg := err.Error()
		if raw := firstNonEmptyString(mapGet(execResult, "error_message"), mapGet(execResult, "error_string")); raw != "" {
			msg = raw
		}
		if strings.Contains(strings.ToLower(msg), "could not bind") {
			return map[string]any{
				"status": "error", "module": fullModulePath,
				"message": fmt.Sprintf("Job start failed: Address/Port likely already in use. %s", msg),
			}
		}
		return map[string]any{
			"status": "error", "module": fullModulePath,
			"message": fmt.Sprintf("Failed to start job: %s", msg),
		}
	}

	jobID, hasJobID := execResult["job_id"]
	uuidVal, _ := execResult["uuid"].(string)
	if !hasJobID || jobID == nil {
		return map[string]any{
			"status":  "unknown",
			"message": fmt.Sprintf("%s executed, but no job ID returned.", capitalizeFirst(modtype)),
			"result":  execResult,
			"module":  fullModulePath,
		}
	}

	var foundSessionID any
	if modtype == "exploit" && uuidVal != "" {
		foundSessionID = pollForSession(ctx, c, uuidVal)
	}

	message := fmt.Sprintf("%s module %s started as job %v.", capitalizeFirst(modtype), fullModulePath, jobID)
	status := "success"
	if modtype == "exploit" {
		if foundSessionID != nil {
			message += fmt.Sprintf(" Session %v created.", foundSessionID)
		} else {
			message += " No session detected within timeout."
			status = "warning"
		}
	}

	return map[string]any{
		"status":       status,
		"message":      message,
		"job_id":       jobID,
		"uuid":         uuidVal,
		"session_id":   foundSessionID,
		"module":       fullModulePath,
		"options":      moduleOptions,
		"payload_name": payloadName,
	}
}

// pollForSession mirrors _execute_module_rpc's post-job session poll:
// EXPLOIT_SESSION_POLL_TIMEOUT=60s, interval EXPLOIT_SESSION_POLL_INTERVAL=2s.
func pollForSession(ctx context.Context, c *Client, targetUUID string) any {
	deadline := time.Now().Add(exploitSessionPollTimeout)
	for time.Now().Before(deadline) {
		if sessions, err := c.ListSessions(ctx); err == nil {
			for id, raw := range sessions {
				if info, ok := raw.(map[string]any); ok {
					if fmt.Sprintf("%v", info["exploit_uuid"]) == targetUUID {
						return id
					}
				}
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(exploitSessionPollInterval):
		}
	}
	return nil
}

// --- send_session_command helpers ---

func runMeterpreterCommand(ctx context.Context, c *Client, sessionID, command string, timeout time.Duration) (output, status, message string) {
	if _, err := c.WriteSession(ctx, sessionID, "meterpreter", command+"\n"); err != nil {
		return "", "error", fmt.Sprintf("Error executing Meterpreter command: %v", err)
	}
	out, ok := pollSessionRead(ctx, c, sessionID, "meterpreter", timeout)
	if !ok {
		return out, "timeout", fmt.Sprintf("Meterpreter command timed out after %s.", timeout)
	}
	return out, "success", "Meterpreter command executed successfully."
}

func runShellCommand(ctx context.Context, c *Client, sessionID, command string, timeout time.Duration) (output, status, message string) {
	if _, err := c.WriteSession(ctx, sessionID, "shell", command+"\n"); err != nil {
		return "", "error", fmt.Sprintf("Error executing Shell command: %v", err)
	}
	if strings.EqualFold(strings.TrimSpace(command), "exit") {
		return "(No output expected after exit)", "success", "Exit command sent to shell session."
	}

	deadline := time.Now().Add(timeout)
	lastData := time.Now()
	var buf strings.Builder
	for {
		if time.Now().After(deadline) {
			return buf.String(), "timeout", fmt.Sprintf("Shell command timed out after %s.", timeout)
		}
		resp, err := c.ReadSession(ctx, sessionID, "shell")
		if err == nil {
			if chunk := readChunkString(resp); chunk != "" {
				buf.WriteString(chunk)
				lastData = time.Now()
				if shellPromptRE.MatchString(buf.String()) {
					return buf.String(), "success", "Shell command executed successfully."
				}
			} else if time.Since(lastData) > sessionReadInactivityWindow {
				return buf.String(), "success", "Shell command likely completed (inactivity)."
			}
		}
		select {
		case <-ctx.Done():
			return buf.String(), "error", ctx.Err().Error()
		case <-time.After(sessionPollInterval):
		}
	}
}

// pollSessionRead is the Meterpreter equivalent of the shell read loop:
// there is no reliable end-of-output prompt to match against, so completion
// is inactivity-based, mirroring what pymetasploit3's run_with_output does
// under the hood for a plain (non-"shell"/"exit") Meterpreter command.
func pollSessionRead(ctx context.Context, c *Client, sessionID, sessionType string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)
	lastActivity := time.Now()
	var buf strings.Builder
	for {
		now := time.Now()
		if now.After(deadline) {
			return buf.String(), false
		}
		resp, err := c.ReadSession(ctx, sessionID, sessionType)
		if err == nil {
			if chunk := readChunkString(resp); chunk != "" {
				buf.WriteString(chunk)
				lastActivity = now
			}
		}
		if now.Sub(lastActivity) > sessionReadInactivityWindow {
			return buf.String(), true
		}
		select {
		case <-ctx.Done():
			return buf.String(), true
		case <-time.After(sessionPollInterval):
		}
	}
}

func readChunkString(resp map[string]any) string {
	switch v := resp["data"].(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

func sleepOrDone(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}

// --- generic small helpers ---

func textResult(v any) *mcp.CallToolResult {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("marshal result", err)
	}
	return mcp.NewToolResultText(string(b))
}

func strOr(v any, fallback string) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fallback
}

func boolOr(v any, fallback bool) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return fallback
}

func intOr(v any, fallback int) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return fallback
	}
}

func limitStrings(s []string, n int) []string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// containsSegment mirrors Python's `arch_lower in p.lower().split("/")`.
func containsSegment(path, segment string) bool {
	for _, part := range strings.Split(path, "/") {
		if part == segment {
			return true
		}
	}
	return false
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func mapGet(m map[string]any, key string) any {
	if m == nil {
		return nil
	}
	return m[key]
}

// payloadBytes normalizes module.execute's "payload" response field:
// msgpack "bin" data decodes to []byte, but some msfrpcd versions send it as
// a "str"-typed field instead.
func payloadBytes(v any) ([]byte, bool) {
	switch b := v.(type) {
	case []byte:
		return b, true
	case string:
		return []byte(b), true
	default:
		return nil, false
	}
}

var filenameSanitizeRE = regexp.MustCompile(`[^a-zA-Z0-9_.\-]`)
var tokenSanitizeRE = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func sanitizeFilename(name string) string {
	if name == "" {
		return ""
	}
	base := filepath.Base(name)
	return filenameSanitizeRE.ReplaceAllString(base, "_")
}

func defaultPayloadFilename(payloadType, formatType string) string {
	timestamp := time.Now().Format("20060102_150405")
	safeType := tokenSanitizeRE.ReplaceAllString(payloadType, "_")
	safeFormat := tokenSanitizeRE.ReplaceAllString(formatType, "_")
	return fmt.Sprintf("payload_%s_%s.%s", safeType, timestamp, safeFormat)
}
