package core

import (
	"context"
	"strings"

	"github.com/anubhavg-icpl/pentester2/internal/llm"
)

// judgeRunObjective is the fixed objective string main() passes into
// judge_output() in final.py (line 558) -- that's the only call site there
// is, so it's baked in here rather than threaded through as a parameter.
const judgeRunObjective = "verify if the exploitation happened or not based on the output given and ONLY return TRUE or FALSE"

// judgeOutput ports judge_output() from final.py: asks the judge model
// whether the orchestrator's final output actually demonstrates the
// exploitation objective was met.
func judgeOutput(ctx context.Context, judge llm.ChatModel, outputStr string) (bool, error) {
	prompt := "You are a quality control judge for a penetration testing agent.\n" +
		"OBJECTIVE: " + judgeRunObjective + "\n" +
		"AGENT OUTPUT: " + outputStr + "\n\n" +
		"Instructions:\n" +
		"- If the agent output contains the decrypted/cleartext password for the specified user, return 'True'.\n" +
		"- If the output only contains hashes, error messages, or failed attempts, return 'False'.\n" +
		"- Output ONLY the word 'True' or 'False'."

	resp, err := judge.Converse(ctx, "", []llm.Message{llm.UserMessage(prompt)}, nil)
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(resp.Text), "true"), nil
}
