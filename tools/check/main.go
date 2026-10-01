// Command check runs the quality gates used by the git hooks and CI.
//
//	go run ./tools/check coverage [-min 80]   go test -race, fails below min % statement coverage
//	go run ./tools/check deadcode             fails if deadcode -test reports unreachable code
//	go run ./tools/check mutation             gremlins, fails if any mutant survives
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/tools/cover"
)

const (
	coverPkg     = "./internal/..."
	coverProfile = "coverage.out"
	mutantReport = "gremlins.json"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: check coverage [-min N] | deadcode | mutation")
		os.Exit(2)
	}

	var err error
	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "coverage":
		err = coverage(args)
	case "deadcode":
		err = deadcode()
	case "mutation":
		err = mutation()
	default:
		err = fmt.Errorf("unknown gate %q", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "check:", err)
		os.Exit(1)
	}
}

func coverage(args []string) error {
	fs := flag.NewFlagSet("coverage", flag.ExitOnError)
	minPct := fs.Float64("min", 80, "minimum statement coverage (%) of "+coverPkg)
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := run("go", "test", "-race", "-covermode=atomic", "-coverpkg="+coverPkg,
		"-coverprofile="+coverProfile, "./..."); err != nil {
		return err
	}

	profiles, err := cover.ParseProfiles(coverProfile)
	if err != nil {
		return err
	}
	var total, covered int
	for _, p := range profiles {
		for _, b := range p.Blocks {
			total += b.NumStmt
			if b.Count > 0 {
				covered += b.NumStmt
			}
		}
	}
	pct := 100.0
	if total > 0 {
		pct = 100 * float64(covered) / float64(total)
	}
	fmt.Printf("statement coverage %.2f%% (%d/%d, min %g%%)\n", pct, covered, total, *minPct)
	if pct < *minPct {
		return fmt.Errorf("coverage below %g%%", *minPct)
	}
	return nil
}

func deadcode() error {
	var out bytes.Buffer
	cmd := exec.Command("go", "tool", "deadcode", "-test", "./...")
	cmd.Stdout, cmd.Stderr = &out, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if out.Len() > 0 {
		return fmt.Errorf("unreachable code found:\n%s", out.Bytes())
	}
	return nil
}

// mutation gates on the JSON report: gremlins v0.6.0 ignores --threshold-efficacy (always exits 0).
func mutation() error {
	if err := run("go", "tool", "gremlins", "unleash", "--exclude-files", "^tools/", "--output", mutantReport); err != nil {
		return err
	}
	data, err := os.ReadFile(mutantReport)
	if err != nil {
		return err
	}
	var report struct {
		Lived int `json:"mutants_lived"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		return err
	}
	if report.Lived > 0 {
		return fmt.Errorf("%d surviving mutants (LIVED above)", report.Lived)
	}
	return nil
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
