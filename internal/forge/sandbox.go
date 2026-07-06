// Package forge is a Docker-sandboxed Python code executor used as the
// "codegen" subagent's one tool.
package forge

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	containerName = "python_lab"
	image         = "python:3.13-slim"
)

// ExecResult is the stdout/stderr/exit-code outcome of a sandboxed command.
type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// EnsureContainer ensures the persistent "python_lab" sandbox container
// exists and is running.
func EnsureContainer(ctx context.Context) error {
	psOut, err := exec.CommandContext(ctx, "docker", "ps", "-a",
		"--filter", fmt.Sprintf("name=^%s$", containerName),
		"--format", "{{.Names}}",
	).Output()
	if err != nil {
		return fmt.Errorf("docker ps: %w", err)
	}

	exists := false
	for _, line := range strings.Split(strings.TrimSpace(string(psOut)), "\n") {
		if line == containerName {
			exists = true
			break
		}
	}

	if !exists {
		cmd := exec.CommandContext(ctx, "docker", "run", "-dit",
			"--name", containerName,
			"--network", "host",
			"--memory", "512m",
			"--cpus", "1",
			image, "bash",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("docker run: %w: %s", err, out)
		}
		return nil
	}

	inspectOut, err := exec.CommandContext(ctx, "docker", "inspect",
		"-f", "{{.State.Running}}", containerName,
	).Output()
	if err != nil {
		return fmt.Errorf("docker inspect: %w", err)
	}
	running := strings.TrimSpace(string(inspectOut)) == "true"

	if !running {
		if out, err := exec.CommandContext(ctx, "docker", "start", containerName).CombinedOutput(); err != nil {
			return fmt.Errorf("docker start: %w: %s", err, out)
		}
	}
	return nil
}

// ExecPython executes code inside the persistent container with a strict
// timeout. timeout <= 0 defaults to 200s.
func ExecPython(ctx context.Context, code string, timeout time.Duration) (stdout, stderr string, exitCode int, err error) {
	if timeout <= 0 {
		timeout = 200 * time.Second
	}
	timeoutSecs := int(timeout.Seconds())
	if timeoutSecs < 1 {
		timeoutSecs = 1
	}

	if err = EnsureContainer(ctx); err != nil {
		return "", "", 0, err
	}

	filename := fmt.Sprintf("/root/script_%s.py", uuid.New().String()[:8])

	tmp, err := os.CreateTemp("", "codegen-*.py")
	if err != nil {
		return "", "", 0, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	defer func() {
		// best-effort cleanup of the file inside the container
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = exec.CommandContext(cleanupCtx, "docker", "exec", containerName, "rm", "-f", filename).Run()
	}()

	if _, err = tmp.WriteString(code); err != nil {
		tmp.Close()
		return "", "", 0, fmt.Errorf("write temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return "", "", 0, fmt.Errorf("close temp file: %w", err)
	}

	if out, cpErr := exec.CommandContext(ctx, "docker", "cp", tmpPath, containerName+":"+filename).CombinedOutput(); cpErr != nil {
		return "", "", 0, fmt.Errorf("docker cp: %w: %s", cpErr, out)
	}

	// Give the subprocess a 5s buffer after the container-side "timeout" fires.
	execCtx, cancel := context.WithTimeout(ctx, timeout+5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(execCtx, "docker", "exec", containerName,
		"timeout", fmt.Sprintf("%d", timeoutSecs), "python", filename,
	)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	runErr := cmd.Run()

	if execCtx.Err() == context.DeadlineExceeded {
		// Fallback in case the local docker command itself hangs.
		killCtx, killCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer killCancel()
		_ = exec.CommandContext(killCtx, "docker", "exec", containerName, "pkill", "-f", filename).Run()
		return strings.TrimSpace(outBuf.String()),
			fmt.Sprintf("Error: Docker execution timed out after %d seconds.", timeoutSecs),
			124, nil
	}

	if runErr != nil {
		if _, ok := runErr.(*exec.ExitError); !ok {
			return "", "", 0, fmt.Errorf("docker exec: %w", runErr)
		}
	}

	exitCode = cmd.ProcessState.ExitCode()
	if exitCode == 124 {
		// Exit code 124 means the Linux 'timeout' command inside the
		// container killed the script.
		return strings.TrimSpace(outBuf.String()),
			fmt.Sprintf("Error: Script execution timed out after %d seconds. Check for infinite loops or missing network timeouts.", timeoutSecs),
			124, nil
	}

	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), exitCode, nil
}

// PipInstall installs a package inside the persistent sandbox container.
func PipInstall(ctx context.Context, pkg string) (ExecResult, error) {
	if err := EnsureContainer(ctx); err != nil {
		return ExecResult{}, err
	}

	cmd := exec.CommandContext(ctx, "docker", "exec", containerName, "pip", "install", pkg)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	runErr := cmd.Run()

	exitCode := 0
	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return ExecResult{}, fmt.Errorf("docker exec pip install: %w", runErr)
		}
	}

	return ExecResult{
		Stdout:   strings.TrimSpace(outBuf.String()),
		Stderr:   strings.TrimSpace(errBuf.String()),
		ExitCode: exitCode,
	}, nil
}
