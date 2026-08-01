package control

import (
	"errors"
	"sync"
	"testing"

	"github.com/anubhavg-icpl/talon/internal/core"
)

func TestClaimInterruptIsAtomic(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "1.1.1.1"})
	store.SetResult("r1", core.RunResult{
		Interrupted: true,
		Interrupt:   &core.PendingInterrupt{ToolName: "nmap_scan", Args: map[string]any{"target": "1.1.1.1"}},
	})

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		wins  int
		fails int
	)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, ok := store.ClaimInterrupt("r1"); ok {
				mu.Lock()
				wins++
				mu.Unlock()
			} else {
				mu.Lock()
				fails++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if wins != 1 {
		t.Fatalf("expected exactly 1 successful claim, got %d (fails=%d)", wins, fails)
	}
	if fails != 19 {
		t.Fatalf("expected 19 failed claims, got %d", fails)
	}
	sess, ok := store.Get("r1")
	if !ok {
		t.Fatal("session missing")
	}
	if sess.PendingInterrupt != nil {
		t.Fatal("pending interrupt should be cleared after claim")
	}
	if sess.Status != "running" {
		t.Fatalf("status=%q want running", sess.Status)
	}
}

func TestSetResultHistoryAndJudge(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r2", core.RunInput{SessionID: "r2", TargetIP: "2.2.2.2"})
	store.SetStatus("r2", "running")
	store.SetToolLog("r2", []core.ToolCallRecord{
		{Index: 0, ToolName: "nmap_scan", Output: "open 21"},
	})
	findings := []core.Finding{{
		ID: "FIND-001", Severity: core.SeverityCritical, Title: "RCE",
		Evidence: core.GateEvidence{Passed: true, Baseline: "none", Attack: "session", Diff: "created"},
	}}
	rep := core.BuildReport(core.RunInput{TargetIP: "2.2.2.2"}, nil, "session opened", findings, true, true)
	store.SetResult("r2", core.RunResult{
		FinalMessage: "session opened",
		ToolLog:      []core.ToolCallRecord{{Index: 0, ToolName: "nmap_scan", Output: "open 21"}},
		JudgeVerdict: true,
		JudgeSet:     true,
		Findings:     findings,
		Report:       &rep,
	})

	sess, ok := store.Get("r2")
	if !ok {
		t.Fatal("missing session")
	}
	if !sess.JudgeSet || !sess.JudgeVerdict {
		t.Fatalf("judge not stored: set=%v verdict=%v", sess.JudgeSet, sess.JudgeVerdict)
	}
	if sess.Status != "completed" {
		t.Fatalf("status=%s", sess.Status)
	}
	if len(sess.History) < 3 {
		t.Fatalf("expected history events, got %v", sess.History)
	}
	if len(sess.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(sess.Findings))
	}
	if sess.Report == nil || sess.Report.Markdown == "" {
		t.Fatal("expected structured report")
	}
}

func TestSetErrorClearsInterrupt(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r3", core.RunInput{SessionID: "r3"})
	store.SetResult("r3", core.RunResult{
		Interrupted: true,
		Interrupt:   &core.PendingInterrupt{ToolName: "nmap_scan"},
	})
	store.SetError("r3", errors.New("boom"))
	sess, ok := store.Get("r3")
	if !ok {
		t.Fatal("missing")
	}
	if sess.PendingInterrupt != nil {
		t.Fatal("interrupt should be cleared on error")
	}
	if sess.Status != "error" {
		t.Fatalf("status=%s", sess.Status)
	}
}
