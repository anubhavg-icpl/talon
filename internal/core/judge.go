package core

import (
	"context"
	"errors"
	"strings"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// judgeRunObjective is the fixed objective every run is judged against --
// there's currently only one call site, so it's baked in here rather than
// threaded through as a parameter.
const judgeRunObjective = "verify whether remote code execution or session compromise was achieved; return ONLY TRUE or FALSE"

// judgeOutput asks the judge model whether the orchestrator's final output
// actually demonstrates the exploitation objective was met (shell/session
// proof — not password cracking).
func judgeOutput(ctx context.Context, judge llm.ChatModel, outputStr string) (bool, error) {
	if judge == nil {
		return false, errors.New("agent: judge model is nil")
	}
	prompt := "You are a quality control judge for a penetration testing agent.\n" +
		"OBJECTIVE: " + judgeRunObjective + "\n" +
		"AGENT OUTPUT: " + outputStr + "\n\n" +
		"Instructions:\n" +
		"- Return 'True' if the output shows proof of successful exploitation, such as: " +
		"an interactive shell/session created, Meterpreter session, remote command output " +
		"(whoami/hostname/id/sysinfo), or clear RCE evidence on the target.\n" +
		"- Return 'False' if the output only has recon results, failed exploits, errors, " +
		"hashes without a session, or no evidence of code execution on the target.\n" +
		"- Output ONLY the word 'True' or 'False'."

	resp, err := judge.Converse(ctx, "", []llm.Message{llm.UserMessage(prompt)}, nil)
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(resp.Text), "true"), nil
}
