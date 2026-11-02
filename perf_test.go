package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestPerformanceOverBudget(t *testing.T) {
	cases := []struct {
		elapsed time.Duration
		matched int
		want    bool
	}{
		{2*time.Second + time.Nanosecond, 200, true},
		{2 * time.Second, 200, false},
		{3 * time.Second, 201, false},
		{3 * time.Second, 0, false},
		{1 * time.Second, 1, false},
	}
	for _, c := range cases {
		if got := performanceOverBudget(c.elapsed, c.matched); got != c.want {
			t.Fatalf("performanceOverBudget(%s, %d) = %v, want %v", c.elapsed, c.matched, got, c.want)
		}
	}
}

func fakeElapsed(d time.Duration) func() {
	old := elapsedSince
	elapsedSince = func(time.Time) time.Duration { return d }
	return func() { elapsedSince = old }
}

func TestPerformanceBudgetDegradesToReportOnly(t *testing.T) {
	restore := fakeElapsed(3 * time.Second)
	t.Cleanup(restore)

	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 0 {
		t.Fatalf("exit = %d, want 0 after performance degrade; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "performance budget exceeded") {
		t.Fatalf("stderr = %q, want performance warning", errs)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 {
		t.Fatalf("findings = %+v, want the violation still reported", r.Findings)
	}
}

func TestPerformanceBudgetKeepsPhase2Blocking(t *testing.T) {
	restore := fakeElapsed(3 * time.Second)
	t.Cleanup(restore)

	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "2")
	if code != 1 {
		t.Fatalf("exit = %d, want 1 in phase 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "performance budget exceeded") {
		t.Fatalf("stderr = %q, want performance warning", errs)
	}
}

func TestPerformanceBudgetWithinLimit(t *testing.T) {
	restore := fakeElapsed(1500 * time.Millisecond)
	t.Cleanup(restore)

	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stderr=%s", code, errs)
	}
	if strings.Contains(errs, "performance budget exceeded") {
		t.Fatalf("stderr = %q, unexpected performance warning", errs)
	}
}

func perfCorpus(t *testing.T, specs int) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	for i := 0; i < specs; i++ {
		writeFile(t, root, fmt.Sprintf("docs/specs/%04d.md", i), doc("Approved", "docs/architecture/003.md"))
	}
	return root
}

func TestNFR001ThousandDocsBudget(t *testing.T) {
	if os.Getenv("MDSTANDARDS_PERF") != "1" {
		t.Skip("set MDSTANDARDS_PERF=1 to run the wall-clock budget check")
	}
	root := perfCorpus(t, 999)

	start := time.Now()
	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "0")
	elapsed := time.Since(start)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	t.Logf("NFR-001: 1000 documents in %s", elapsed.Round(time.Millisecond))
	if elapsed >= 10*time.Second {
		t.Fatalf("1000 documents took %s, want < 10s", elapsed)
	}
}

func TestNFR002P95Budget(t *testing.T) {
	if os.Getenv("MDSTANDARDS_PERF") != "1" {
		t.Skip("set MDSTANDARDS_PERF=1 to run the wall-clock budget check")
	}
	root := perfCorpus(t, 199)

	const runs = 20
	durations := make([]time.Duration, 0, runs)
	for i := 0; i < runs; i++ {
		start := time.Now()
		code, _, errs := runCLI(t, root, "--format", "json", "--phase", "0")
		durations = append(durations, time.Since(start))
		if code != 0 {
			t.Fatalf("run %d: exit = %d; stderr=%s", i, code, errs)
		}
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p95 := durations[int(math.Ceil(0.95*runs))-1]
	t.Logf("NFR-002: 200 documents p95 = %s over %d runs", p95.Round(time.Millisecond), runs)
	if p95 >= 2*time.Second {
		t.Fatalf("p95 = %s, want < 2s", p95)
	}
}
