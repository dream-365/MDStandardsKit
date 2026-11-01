package main

import (
	"context"
	"testing"

	"github.com/open-policy-agent/opa/rego"
)

// T1: prove that policy/data/standard.json resolves as data.standard inside OPA.
//
// OPA merges data documents at the namespace root. A bare object would land at
// data.governed_globs, not data.standard.governed_globs, so this test exists to
// keep the wrapper from being removed by accident.
func TestStandardDataResolvesAsDataStandard(t *testing.T) {
	ctx := context.Background()

	q, err := rego.New(
		rego.Query("data.standard.dependable_statuses"),
		rego.Load([]string{"policy/data/standard.json"}, nil),
	).PrepareForEval(ctx)
	if err != nil {
		t.Fatalf("prepare for eval: %v", err)
	}

	rs, err := q.Eval(ctx)
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		t.Fatal("data.standard.dependable_statuses is undefined")
	}

	got, ok := rs[0].Expressions[0].Value.([]interface{})
	if !ok {
		t.Fatalf("unexpected value type %T", rs[0].Expressions[0].Value)
	}
	want := []string{"Approved", "Published"}
	if len(got) != len(want) {
		t.Fatalf("dependable_statuses = %v, want %v", got, want)
	}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("dependable_statuses[%d] = %v, want %q", i, got[i], w)
		}
	}
}
