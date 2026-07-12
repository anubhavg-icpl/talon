// Command talon is the production operator CLI for the Talon platform.
// It is a thin client over talon-core's HTTP control plane — not an
// orchestrator itself.
//
//	talon status
//	talon run start --ip ... --cve ...
//	talon run watch <run_id>
//	talon run approve <run_id>
package main

import (
	"os"

	"github.com/anubhavg-icpl/talon/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
