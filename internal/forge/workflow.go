package forge

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/anubhavg-icpl/pentester2/internal/llm"
)

// generatorSystemPrompt instructs the model on how to write executable
// exploit code for the sandbox: hardcode inputs, always set network
// timeouts, print clear success/failure output.
const generatorSystemPrompt = "You are an expert Python exploit developer and pentester.\n" +
	"Generate ONLY executable Python code that accomplishes the security testing task.\n\n" +
	"REQUIREMENTS:\n" +
	"- Output Python code in a python code block\n" +
	"- Hardcode all values from the request (IPs, ports, CVE IDs, etc.)\n" +
	"- Do NOT use sys.argv, argparse, input(), or environment variables\n" +
	"- Include error handling and informative output\n" +
	"- Print clear success/failure messages\n" +
	"- CRITICAL: You MUST include timeouts for ALL network operations (e.g., requests.get(..., timeout=5) or socket.setdefaulttimeout(5)).\n" +
	"- Unless specified otherwise, assume standard ports (21, 23, 6667, etc.) use plaintext sockets, not SSL.\n" +
	"COMMON TASKS:\n"

// feedbackSystemPrompt constrains the completion-check response to a
// strict boolean.
const feedbackSystemPrompt = "Respond ONLY true or false."

const (
	codeExecutorAttemptsMax  = 2
	codeGeneratorAttemptsMax = 2
)

var (
	// codeBlockRe matches a python-tagged (or untagged) fenced code block.
	codeBlockRe = regexp.MustCompile("(?s)```(?:python)?\\s*(.*?)```")

	// moduleNotFoundRe extracts the missing package name from a
	// ModuleNotFoundError so it can be auto-installed and retried.
	moduleNotFoundRe = regexp.MustCompile(`ModuleNotFoundError: No module named '([^']+)'`)

	errorKeywords = []string{"error", "traceback", "exception", "failed", "failure"}
)

// state tracks one generate/execute/feedback retry cycle.
type state struct {
	pseudocode string
	code       string
	output     string
	completed  bool

	executorAttempts  int
	generatorAttempts int
}

// Coder runs the generator/executor/feedback retry loop: generate code,
// execute it in the sandbox, auto-fix a missing-package error and retry,
// then ask the model whether the task succeeded -- up to 2 generation
// attempts x 2 execution-fix attempts.
func Coder(ctx context.Context, model llm.ChatModel, pseudocode string) (string, error) {
	s := &state{pseudocode: pseudocode}

	for {
		if err := generate(ctx, model, s); err != nil {
			return s.output, err
		}

		if err := execute(ctx, s); err != nil {
			return s.output, err
		}

		for hasError(s.output) && s.executorAttempts < codeExecutorAttemptsMax {
			handleError(ctx, s)
			if err := execute(ctx, s); err != nil {
				return s.output, err
			}
		}

		completed, err := checkFeedback(ctx, model, s)
		if err != nil {
			return s.output, err
		}
		s.completed = completed

		if s.completed {
			return s.output, nil
		}
		if s.generatorAttempts >= codeGeneratorAttemptsMax {
			return s.output, nil
		}
		// retry: fall through to another generate() with accumulated context
	}
}

// generate asks the model for code given the pseudocode plus any previous
// attempt's code/output, then extracts the fenced code block.
func generate(ctx context.Context, model llm.ChatModel, s *state) error {
	messages := []llm.Message{
		llm.UserMessage(s.pseudocode),
		llm.UserMessage(fmt.Sprintf("Previous code:\n%s\n\nOutput:\n%s", s.code, s.output)),
	}

	resp, err := model.Converse(ctx, generatorSystemPrompt, messages, nil)
	if err != nil {
		return fmt.Errorf("codegen generate: %w", err)
	}

	s.code = extractCode(resp.Text)
	s.output = ""
	s.executorAttempts = 0
	s.generatorAttempts++
	return nil
}

// extractCode pulls the fenced code block out of a model response, falling
// back to the whole response if no fenced block is found.
func extractCode(text string) string {
	if m := codeBlockRe.FindStringSubmatch(text); m != nil {
		return strings.TrimSpace(m[1])
	}
	return strings.TrimSpace(text)
}

// execute runs the current generated code in the sandbox and records its
// combined stdout/stderr.
func execute(ctx context.Context, s *state) error {
	stdout, stderr, _, err := ExecPython(ctx, s.code, 0)
	if err != nil {
		return fmt.Errorf("codegen execute: %w", err)
	}

	combined := strings.TrimSpace(stdout + stderr)
	if combined == "" {
		combined = "[No Output]"
	}
	s.output = combined
	s.executorAttempts++
	return nil
}

// hasError reports whether output contains a known error indicator.
func hasError(output string) bool {
	lower := strings.ToLower(output)
	for _, kw := range errorKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// handleError best-effort pip-installs a missing module found via
// ModuleNotFoundError, then falls through to retry.
func handleError(ctx context.Context, s *state) {
	m := moduleNotFoundRe.FindStringSubmatch(s.output)
	if m == nil {
		return
	}
	// ponytail: ignore PipInstall errors -- retry happens regardless of
	// whether the install succeeded, upgrade if silent install failures
	// need to short-circuit the retry loop.
	_, _ = PipInstall(ctx, m[1])
}

// checkFeedback asks the model whether the task was accomplished given the
// observed output.
func checkFeedback(ctx context.Context, model llm.ChatModel, s *state) (bool, error) {
	messages := []llm.Message{
		llm.UserMessage(fmt.Sprintf("Task: %s\nOutput:\n%s", s.pseudocode, s.output)),
	}

	resp, err := model.Converse(ctx, feedbackSystemPrompt, messages, nil)
	if err != nil {
		return false, fmt.Errorf("codegen feedback: %w", err)
	}

	return strings.Contains(strings.ToLower(resp.Text), "true"), nil
}
