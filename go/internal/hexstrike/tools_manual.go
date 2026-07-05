package hexstrike

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerManualTools ports the 11 hexstrike_mcp.py tools whose bodies do
// more than a straight param->body passthrough (comma-split lists, int
// clamping, or fully local playbook generation) and so don't fit the
// mechanical generatedTools table.
//
// ponytail: ai_generate_attack_suite and comprehensive_api_audit called
// `self.ai_generate_payload(...)` / `self.api_fuzzer(...)` etc. in the
// Python source with no `self` param in scope -- a real bug (NameError at
// call time). Here they call the same underlying HexStrike endpoints
// directly instead, which is what the Python code was clearly trying to do.
func registerManualTools(srv *server.MCPServer, client *Client) {
	textResult := func(v any) *mcp.CallToolResult {
		b, err := json.Marshal(v)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("marshal result", err)
		}
		return mcp.NewToolResultText(string(b))
	}
	splitCSV := func(s string) []string {
		if s == "" {
			return []string{}
		}
		parts := strings.Split(s, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if t := strings.TrimSpace(p); t != "" {
				out = append(out, t)
			}
		}
		return out
	}
	clamp := func(v, lo, hi int64) int64 {
		if v < lo {
			return lo
		}
		if v > hi {
			return hi
		}
		return v
	}

	srv.AddTool(
		mcp.NewTool("api_fuzzer",
			mcp.WithDescription("Advanced API endpoint fuzzing with intelligent parameter discovery."),
			mcp.WithString("base_url", mcp.Required()),
			mcp.WithString("endpoints", mcp.DefaultString("")),
			mcp.WithString("methods", mcp.DefaultString("GET,POST,PUT,DELETE")),
			mcp.WithString("wordlist", mcp.DefaultString("/usr/share/wordlists/api/api-endpoints.txt")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			body := map[string]any{
				"base_url":  strOr(a["base_url"], ""),
				"endpoints": splitCSV(strOr(a["endpoints"], "")),
				"methods":   strings.Split(strOr(a["methods"], "GET,POST,PUT,DELETE"), ","),
				"wordlist":  strOr(a["wordlist"], "/usr/share/wordlists/api/api-endpoints.txt"),
			}
			return textResult(client.Post("api/tools/api_fuzzer", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("discover_attack_chains",
			mcp.WithDescription("Discover multi-stage attack chains for target software with vulnerability correlation."),
			mcp.WithString("target_software", mcp.Required()),
			mcp.WithNumber("attack_depth", mcp.DefaultNumber(3)),
			mcp.WithBoolean("include_zero_days", mcp.DefaultBool(false)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			depth := clamp(int64(intOr(a["attack_depth"], 3)), 1, 5)
			body := map[string]any{
				"target_software":   strOr(a["target_software"], ""),
				"attack_depth":      depth,
				"include_zero_days": boolOr(a["include_zero_days"], false),
			}
			return textResult(client.Post("api/vuln-intel/attack-chains", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("correlate_threat_intelligence",
			mcp.WithDescription("Correlate threat intelligence across multiple sources with advanced analysis."),
			mcp.WithString("indicators", mcp.Required()),
			mcp.WithString("timeframe", mcp.DefaultString("30d")),
			mcp.WithString("sources", mcp.DefaultString("all")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			timeframe := strOr(a["timeframe"], "30d")
			switch timeframe {
			case "7d", "30d", "90d", "1y":
			default:
				timeframe = "30d"
			}
			indicators := splitCSV(strOr(a["indicators"], ""))
			if len(indicators) == 0 {
				return textResult(map[string]any{"success": false, "error": "No valid indicators provided"}), nil
			}
			body := map[string]any{
				"indicators": indicators,
				"timeframe":  timeframe,
				"sources":    strOr(a["sources"], "all"),
			}
			return textResult(client.Post("api/vuln-intel/threat-feeds", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("optimize_tool_parameters_ai",
			mcp.WithDescription("Use AI to optimize tool parameters based on target profile and context."),
			mcp.WithString("target", mcp.Required()),
			mcp.WithString("tool", mcp.Required()),
			mcp.WithString("context", mcp.DefaultString("{}")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			ctxStr := strOr(a["context"], "{}")
			var ctxDict map[string]any
			if ctxStr != "{}" {
				_ = json.Unmarshal([]byte(ctxStr), &ctxDict) // best-effort, same as the Python try/except
			}
			if ctxDict == nil {
				ctxDict = map[string]any{}
			}
			body := map[string]any{
				"target":  strOr(a["target"], ""),
				"tool":    strOr(a["tool"], ""),
				"context": ctxDict,
			}
			return textResult(client.Post("api/intelligence/optimize-parameters", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("bugbounty_reconnaissance_workflow",
			mcp.WithDescription("Create comprehensive reconnaissance workflow for bug bounty hunting."),
			mcp.WithString("domain", mcp.Required()),
			mcp.WithString("scope", mcp.DefaultString("")),
			mcp.WithString("out_of_scope", mcp.DefaultString("")),
			mcp.WithString("program_type", mcp.DefaultString("web")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			body := map[string]any{
				"domain":       strOr(a["domain"], ""),
				"scope":        splitRaw(strOr(a["scope"], "")),
				"out_of_scope": splitRaw(strOr(a["out_of_scope"], "")),
				"program_type": strOr(a["program_type"], "web"),
			}
			return textResult(client.Post("api/bugbounty/reconnaissance-workflow", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("bugbounty_vulnerability_hunting",
			mcp.WithDescription("Create vulnerability hunting workflow prioritized by impact and bounty potential."),
			mcp.WithString("domain", mcp.Required()),
			mcp.WithString("priority_vulns", mcp.DefaultString("rce,sqli,xss,idor,ssrf")),
			mcp.WithString("bounty_range", mcp.DefaultString("unknown")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			body := map[string]any{
				"domain":         strOr(a["domain"], ""),
				"priority_vulns": splitRaw(strOr(a["priority_vulns"], "rce,sqli,xss,idor,ssrf")),
				"bounty_range":   strOr(a["bounty_range"], "unknown"),
			}
			return textResult(client.Post("api/bugbounty/vulnerability-hunting-workflow", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("bugbounty_comprehensive_assessment",
			mcp.WithDescription("Create comprehensive bug bounty assessment combining all specialized workflows."),
			mcp.WithString("domain", mcp.Required()),
			mcp.WithString("scope", mcp.DefaultString("")),
			mcp.WithString("priority_vulns", mcp.DefaultString("rce,sqli,xss,idor,ssrf")),
			mcp.WithBoolean("include_osint", mcp.DefaultBool(true)),
			mcp.WithBoolean("include_business_logic", mcp.DefaultBool(true)),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			body := map[string]any{
				"domain":                 strOr(a["domain"], ""),
				"scope":                  splitRaw(strOr(a["scope"], "")),
				"priority_vulns":         splitRaw(strOr(a["priority_vulns"], "rce,sqli,xss,idor,ssrf")),
				"include_osint":          boolOr(a["include_osint"], true),
				"include_business_logic": boolOr(a["include_business_logic"], true),
			}
			return textResult(client.Post("api/bugbounty/comprehensive-assessment", body)), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("ai_generate_attack_suite",
			mcp.WithDescription("Generate comprehensive attack suite with multiple payload types."),
			mcp.WithString("target_url", mcp.Required()),
			mcp.WithString("attack_types", mcp.DefaultString("xss,sqli,lfi")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			targetURL := strOr(a["target_url"], "")
			attackList := splitCSV(strOr(a["attack_types"], "xss,sqli,lfi"))

			suites := map[string]any{}
			totalPayloads, highRisk, testCases := 0, 0, 0
			for _, attackType := range attackList {
				res := client.Post("api/ai/generate_payload", map[string]any{
					"attack_type": attackType,
					"complexity":  "advanced",
					"technology":  "",
					"url":         targetURL,
				})
				if ok, _ := res["success"].(bool); !ok {
					continue
				}
				payloadData, _ := res["ai_payload_generation"].(map[string]any)
				if payloadData == nil {
					continue
				}
				suites[attackType] = payloadData
				totalPayloads += intOr(payloadData["payload_count"], 0)
				if tc, ok := payloadData["test_cases"].([]any); ok {
					testCases += len(tc)
				}
				if payloads, ok := payloadData["payloads"].([]any); ok {
					for _, p := range payloads {
						if pm, ok := p.(map[string]any); ok && pm["risk_level"] == "HIGH" {
							highRisk++
						}
					}
				}
			}

			results := map[string]any{
				"target_url":     targetURL,
				"attack_types":   attackList,
				"payload_suites": suites,
				"summary": map[string]any{
					"total_payloads":     totalPayloads,
					"high_risk_payloads": highRisk,
					"test_cases":         testCases,
				},
			}
			return textResult(map[string]any{
				"success":      true,
				"attack_suite": results,
				"timestamp":    time.Now().Unix(),
			}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("comprehensive_api_audit",
			mcp.WithDescription("Comprehensive API security audit combining multiple testing techniques."),
			mcp.WithString("base_url", mcp.Required()),
			mcp.WithString("schema_url", mcp.DefaultString("")),
			mcp.WithString("jwt_token", mcp.DefaultString("")),
			mcp.WithString("graphql_endpoint", mcp.DefaultString("")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			baseURL := strOr(a["base_url"], "")
			schemaURL := strOr(a["schema_url"], "")
			jwtToken := strOr(a["jwt_token"], "")
			graphqlEndpoint := strOr(a["graphql_endpoint"], "")

			performed := []string{}
			totalVulns := 0
			audit := map[string]any{
				"base_url":        baseURL,
				"audit_timestamp": time.Now().Unix(),
			}

			fuzz := client.Post("api/tools/api_fuzzer", map[string]any{
				"base_url": baseURL, "endpoints": []string{}, "methods": []string{"GET", "POST", "PUT", "DELETE"},
				"wordlist": "/usr/share/wordlists/api/api-endpoints.txt",
			})
			if ok, _ := fuzz["success"].(bool); ok {
				performed = append(performed, "api_fuzzing")
				audit["api_fuzzing"] = fuzz
			}

			if schemaURL != "" {
				schema := client.Post("api/tools/api_schema_analyzer", map[string]any{"schema_url": schemaURL, "schema_type": "openapi"})
				if ok, _ := schema["success"].(bool); ok {
					performed = append(performed, "schema_analysis")
					audit["schema_analysis"] = schema
					if sd, ok := schema["schema_analysis_results"].(map[string]any); ok {
						if issues, ok := sd["security_issues"].([]any); ok {
							totalVulns += len(issues)
						}
					}
				}
			}

			if jwtToken != "" {
				jwt := client.Post("api/tools/jwt_analyzer", map[string]any{"jwt_token": jwtToken, "target_url": baseURL})
				if ok, _ := jwt["success"].(bool); ok {
					performed = append(performed, "jwt_analysis")
					audit["jwt_analysis"] = jwt
					if jd, ok := jwt["jwt_analysis_results"].(map[string]any); ok {
						if vulns, ok := jd["vulnerabilities"].([]any); ok {
							totalVulns += len(vulns)
						}
					}
				}
			}

			if graphqlEndpoint != "" {
				gql := client.Post("api/tools/graphql_scanner", map[string]any{
					"endpoint": graphqlEndpoint, "introspection": true, "query_depth": 10, "test_mutations": true,
				})
				if ok, _ := gql["success"].(bool); ok {
					performed = append(performed, "graphql_scanning")
					audit["graphql_scanning"] = gql
					if gd, ok := gql["graphql_scan_results"].(map[string]any); ok {
						if vulns, ok := gd["vulnerabilities"].([]any); ok {
							totalVulns += len(vulns)
						}
					}
				}
			}

			audit["tests_performed"] = performed
			audit["total_vulnerabilities"] = totalVulns
			audit["recommendations"] = []string{
				"Implement proper authentication and authorization",
				"Use HTTPS for all API communications",
				"Validate and sanitize all input parameters",
				"Implement rate limiting and request throttling",
				"Add comprehensive logging and monitoring",
				"Regular security testing and code reviews",
				"Keep API documentation updated and secure",
				"Implement proper error handling",
			}
			coverage := "partial"
			if len(performed) >= 3 {
				coverage = "comprehensive"
			}
			audit["summary"] = map[string]any{
				"tests_performed":       len(performed),
				"total_vulnerabilities": totalVulns,
				"audit_coverage":        coverage,
			}

			return textResult(map[string]any{"success": true, "comprehensive_audit": audit}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("threat_hunting_assistant",
			mcp.WithDescription("AI-powered threat hunting assistant with vulnerability correlation and attack simulation."),
			mcp.WithString("target_environment", mcp.Required()),
			mcp.WithString("threat_indicators", mcp.DefaultString("")),
			mcp.WithString("hunt_focus", mcp.DefaultString("general")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			targetEnv := strOr(a["target_environment"], "")
			huntFocus := strOr(a["hunt_focus"], "general")
			switch huntFocus {
			case "general", "apt", "ransomware", "insider_threat", "supply_chain":
			default:
				huntFocus = "general"
			}
			indicators := splitCSV(strOr(a["threat_indicators"], ""))

			var detectionQueries []string
			lowerEnv := strings.ToLower(targetEnv)
			switch {
			case strings.Contains(lowerEnv, "windows"):
				detectionQueries = []string{
					"Get-WinEvent | Where-Object {$_.Id -eq 4688 -and $_.Message -like '*suspicious*'}",
					"Get-Process | Where-Object {$_.ProcessName -notin @('explorer.exe', 'svchost.exe')}",
					`Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Run`,
					"Get-NetTCPConnection | Where-Object {$_.State -eq 'Established' -and $_.RemoteAddress -notlike '10.*'}",
				}
			case strings.Contains(lowerEnv, "cloud"):
				detectionQueries = []string{
					"CloudTrail logs for unusual API calls",
					"Failed authentication attempts from unknown IPs",
					"Privilege escalation events",
					"Data exfiltration indicators",
				}
			default:
				detectionQueries = []string{}
			}

			focusScenarios := map[string][]string{
				"apt": {
					"Spear phishing with weaponized documents",
					"Living-off-the-land techniques",
					"Lateral movement via stolen credentials",
					"Data staging and exfiltration",
				},
				"ransomware": {
					"Initial access via RDP/VPN",
					"Privilege escalation and persistence",
					"Shadow copy deletion",
					"Encryption and ransom note deployment",
				},
				"insider_threat": {
					"Unusual data access patterns",
					"After-hours activity",
					"Large data downloads",
					"Access to sensitive systems",
				},
			}
			threatScenarios, ok := focusScenarios[huntFocus]
			if !ok {
				threatScenarios = []string{
					"Unauthorized access attempts",
					"Suspicious process execution",
					"Network anomalies",
					"Data access violations",
				}
			}

			playbook := map[string]any{
				"target_environment":  targetEnv,
				"hunt_focus":          huntFocus,
				"indicators_analyzed": indicators,
				"detection_queries":   detectionQueries,
				"threat_scenarios":    threatScenarios,
				"investigation_steps": []string{
					"1. Validate initial indicators and expand IOC list",
					"2. Run detection queries and analyze results",
					"3. Correlate events across multiple data sources",
					"4. Identify affected systems and user accounts",
					"5. Assess scope and impact of potential compromise",
					"6. Implement containment measures if threat confirmed",
					"7. Document findings and update detection rules",
				},
				"mitigation_strategies": []string{},
			}

			if len(indicators) > 0 {
				corr := client.Post("api/vuln-intel/threat-feeds", map[string]any{
					"indicators": indicators, "timeframe": "30d", "sources": "all",
				})
				if ok, _ := corr["success"].(bool); ok {
					playbook["threat_correlation"] = corr["threat_intelligence"]
				}
			}

			return textResult(map[string]any{"success": true, "hunting_playbook": playbook}), nil
		},
	)

	srv.AddTool(
		mcp.NewTool("bugbounty_authentication_bypass_testing",
			mcp.WithDescription("Create authentication bypass testing workflow for bug bounty hunting."),
			mcp.WithString("target_url", mcp.Required()),
			mcp.WithString("auth_type", mcp.DefaultString("form")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			a := req.GetArguments()
			targetURL := strOr(a["target_url"], "")
			authType := strOr(a["auth_type"], "form")

			bypassTechniques := map[string][]map[string]any{
				"form": {
					{"technique": "SQL Injection", "payloads": []string{"admin'--", "' OR '1'='1'--"}},
					{"technique": "Default Credentials", "payloads": []string{"admin:admin", "admin:password"}},
					{"technique": "Password Reset", "description": "Test password reset token reuse and manipulation"},
					{"technique": "Session Fixation", "description": "Test session ID prediction and fixation"},
				},
				"jwt": {
					{"technique": "Algorithm Confusion", "description": "Change RS256 to HS256"},
					{"technique": "None Algorithm", "description": "Set algorithm to 'none'"},
					{"technique": "Key Confusion", "description": "Use public key as HMAC secret"},
					{"technique": "Token Manipulation", "description": "Modify claims and resign token"},
				},
				"oauth": {
					{"technique": "Redirect URI Manipulation", "description": "Test open redirect in redirect_uri"},
					{"technique": "State Parameter", "description": "Test CSRF via missing/weak state parameter"},
					{"technique": "Code Reuse", "description": "Test authorization code reuse"},
					{"technique": "Client Secret", "description": "Test for exposed client secrets"},
				},
				"saml": {
					{"technique": "XML Signature Wrapping", "description": "Manipulate SAML assertions"},
					{"technique": "XML External Entity", "description": "Test XXE in SAML requests"},
					{"technique": "Replay Attacks", "description": "Test assertion replay"},
					{"technique": "Signature Bypass", "description": "Test signature validation bypass"},
				},
			}

			workflow := map[string]any{
				"target":            targetURL,
				"auth_type":         authType,
				"bypass_techniques": bypassTechniques[authType],
				"testing_phases": []map[string]string{
					{"phase": "reconnaissance", "description": "Identify authentication mechanisms"},
					{"phase": "baseline_testing", "description": "Test normal authentication flow"},
					{"phase": "bypass_testing", "description": "Apply bypass techniques"},
					{"phase": "privilege_escalation", "description": "Test for privilege escalation"},
				},
				"estimated_time":          240,
				"manual_testing_required": true,
			}

			return textResult(map[string]any{
				"success":   true,
				"workflow":  workflow,
				"timestamp": time.Now().Format(time.RFC3339),
			}), nil
		},
	)
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

// splitRaw mirrors Python's bare `s.split(",")` (no trimming/filtering),
// used by the three bugbounty_* tools' scope/priority_vulns fields.
func splitRaw(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
