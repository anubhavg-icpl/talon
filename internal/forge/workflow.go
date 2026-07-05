package forge

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/anubhavg-icpl/pentester2/internal/llm"
)

// generatorSystemPrompt is code_gen.py lines 196-209, verbatim.
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

// feedbackSystemPrompt mirrors feedback()'s SystemMessage.
const feedbackSystemPrompt = "Respond ONLY true or false."

const (
	codeExecutorAttemptsMax  = 2
	codeGeneratorAttemptsMax = 2
)

var (
	// codeBlockRe matches a python-tagged (or untagged) fenced code block,
	// same pattern as the Python re.search(r"```(?:python)?\s*(.*?)```", ...).
	codeBlockRe = regexp.MustCompile("(?s)```(?:python)?\\s*(.*?)```")

	// moduleNotFoundRe extracts the missing package name from a
	// ModuleNotFoundError, same pattern as code_error()'s regex.
	moduleNotFoundRe = regexp.MustCompile(`ModuleNotFoundError: No module named '([^']+)'`)

	errorKeywords = []string{"error", "traceback", "exception", "failed", "failure"}
)

// state mirrors CodeExecutorState in code_gen.py.
type state struct {
	pseudocode string
	code       string
	output     string
	completed  bool

	executorAttempts  int
	generatorAttempts int
}

// Coder runs the generator/executor/feedback retry loop described by
// get_coder_workflow() in code_gen.py, using a plain Go loop instead of a
// LangGraph state machine since the state space is bounded (2x2 attempts).
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

// generate mirrors code_generator(): asks the model for code given the
// pseudocode plus previous code/output, then extracts the fenced block.
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

// extractCode mirrors the code_match regex in code_generator(), falling
// back to the whole response if no fenced block is found.
func extractCode(text string) string {
	if m := codeBlockRe.FindStringSubmatch(text); m != nil {
		return strings.TrimSpace(m[1])
	}
	return strings.TrimSpace(text)
}

// execute mirrors code_executor().
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

// hasError mirrors executor_router()'s error-keyword check.
func hasError(output string) bool {
	lower := strings.ToLower(output)
	for _, kw := range errorKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// handleError mirrors code_error(): best-effort pip install of a missing
// module found via ModuleNotFoundError, then falls through to retry.
func handleError(ctx context.Context, s *state) {
	m := moduleNotFoundRe.FindStringSubmatch(s.output)
	if m == nil {
		return
	}
	// ponytail: ignore PipInstall errors, same as the Python which never
	// inspects pip_install()'s return value either -- retry happens
	// regardless of whether the install succeeded.
	_, _ = PipInstall(ctx, m[1])
}

// checkFeedback mirrors feedback(): asks the model whether the task was
// accomplished given the output.
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
