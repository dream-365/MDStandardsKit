package main

import (
	"reflect"
	"testing"
)

func TestLoadRepoConfig(t *testing.T) {
	cases := []struct {
		name      string
		content   string
		missing   bool
		wantMode  repoMode
		wantType  string
		wantTypes []string
		wantErr   bool
	}{
		{name: "missing file is legacy spec", missing: true, wantMode: legacyMode, wantType: "spec"},
		{name: "legacy spec", content: `{"type":"spec"}`, wantMode: legacyMode, wantType: "spec"},
		{name: "legacy idea", content: `{"type":"idea"}`, wantMode: legacyMode, wantType: "idea"},
		{name: "legacy architecture rejected", content: `{"type":"architecture"}`, wantErr: true},
		{name: "empty legacy type", content: `{"type":""}`, wantErr: true},
		{name: "multi", content: `{"types":["spec","idea"]}`, wantMode: multiMode, wantTypes: []string{"spec", "idea"}},
		{name: "multi trims entries", content: `{"types":[" spec "," idea "]}`, wantMode: multiMode, wantTypes: []string{"spec", "idea"}},
		{name: "type null types absent", content: `{"type":null}`, wantErr: true},
		{name: "types null type absent", content: `{"types":null}`, wantErr: true},
		{name: "empty object", content: `{}`, wantErr: true},
		{name: "both keys", content: `{"type":"spec","types":["spec"]}`, wantErr: true},
		{name: "type null with types", content: `{"type":null,"types":["spec"]}`, wantMode: multiMode, wantTypes: []string{"spec"}},
		{name: "types null with type", content: `{"type":"spec","types":null}`, wantMode: legacyMode, wantType: "spec"},
		{name: "empty types", content: `{"types":[]}`, wantErr: true},
		{name: "empty types entry", content: `{"types":["spec"," "]}`, wantErr: true},
		{name: "duplicate types", content: `{"types":["spec"," spec "]}`, wantErr: true},
		{name: "unknown key", content: `{"type":"spec","globs":["docs/**"]}`, wantErr: true},
		{name: "trailing data", content: `{"type":"spec"} {}`, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if !tc.missing {
				writeFile(t, root, ".mdstandards.json", tc.content)
			}
			got, err := loadRepoConfig(root)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("loadRepoConfig() = %+v, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("loadRepoConfig(): %v", err)
			}
			if got.Mode != tc.wantMode || got.Type != tc.wantType || !reflect.DeepEqual(got.Types, tc.wantTypes) {
				t.Fatalf("loadRepoConfig() = %+v, want mode=%v type=%q types=%v", got, tc.wantMode, tc.wantType, tc.wantTypes)
			}
		})
	}
}

func testStandard() Standard {
	return Standard{
		StandardVersion:    "2",
		GovernedGlobs:      []string{"docs/specs/**/*.md", "docs/architecture/**/*.md"},
		DependableStatuses: []string{"Approved", "Published"},
		DocumentTypes: map[string]DocumentType{
			"spec": {RuleSet: "spec", Graph: true},
			"architecture": {
				RuleSet: "spec",
				Graph:   false,
				Globs:   []string{"docs/architecture/**/*.md"},
			},
			"changelog": {
				RuleSet: "none",
				Graph:   false,
				Globs:   []string{"docs/specs/**/changelog.md"},
			},
			"idea": {
				RuleSet:        "idea",
				Graph:          false,
				Globs:          []string{"20*/**/*.md"},
				RequiredFields: []string{"id", "title", "date", "tags", "status"},
				StatusValues:   []string{"萌芽"},
				RequiredHeadings: []string{
					"一句话描述", "想法", "为什么有意思", "下一步",
				},
				IDPattern:          "^[A-Z]{8}$",
				DateField:          "date",
				PlaceholderMarkers: []string{"<!-- placeholder -->"},
			},
		},
	}
}

func TestBuildRunPlan(t *testing.T) {
	std := testStandard()

	t.Run("legacy spec", func(t *testing.T) {
		plan, err := buildRunPlan(std, RepoConfig{Mode: legacyMode, Type: "spec"})
		if err != nil {
			t.Fatal(err)
		}
		if plan.Mode != legacyMode || !plan.Graph || len(plan.Globs) != 2 {
			t.Fatalf("plan = %+v", plan)
		}
	})

	t.Run("legacy idea", func(t *testing.T) {
		plan, err := buildRunPlan(std, RepoConfig{Mode: legacyMode, Type: "idea"})
		if err != nil {
			t.Fatal(err)
		}
		if plan.Mode != legacyMode || plan.Graph || len(plan.Globs) != 1 || plan.Globs[0] != "20*/**/*.md" {
			t.Fatalf("plan = %+v", plan)
		}
	})

	t.Run("legacy architecture rejected", func(t *testing.T) {
		if _, err := buildRunPlan(std, RepoConfig{Mode: legacyMode, Type: "architecture"}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("multi unions top-level and enabled globs", func(t *testing.T) {
		plan, err := buildRunPlan(std, RepoConfig{Mode: multiMode, Types: []string{"spec", "architecture", "changelog"}})
		if err != nil {
			t.Fatal(err)
		}
		if plan.Mode != multiMode || !plan.Graph {
			t.Fatalf("plan = %+v", plan)
		}
		wantGlobs := []string{
			"docs/specs/**/*.md",
			"docs/architecture/**/*.md",
			"docs/specs/**/changelog.md",
		}
		if !reflect.DeepEqual(plan.Globs, wantGlobs) {
			t.Fatalf("globs = %v, want %v", plan.Globs, wantGlobs)
		}
		if len(plan.Profiles) != 3 {
			t.Fatalf("profiles = %v", plan.Profiles)
		}
	})

	t.Run("top-level baseline survives without spec", func(t *testing.T) {
		plan, err := buildRunPlan(std, RepoConfig{Mode: multiMode, Types: []string{"changelog"}})
		if err != nil {
			t.Fatal(err)
		}
		if plan.Globs[0] != "docs/specs/**/*.md" || plan.Globs[1] != "docs/architecture/**/*.md" {
			t.Fatalf("globs = %v", plan.Globs)
		}
	})

	t.Run("unknown profile", func(t *testing.T) {
		if _, err := buildRunPlan(std, RepoConfig{Mode: multiMode, Types: []string{"note"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("missing rule_set", func(t *testing.T) {
		modified := testStandard()
		p := modified.DocumentTypes["changelog"]
		p.RuleSet = ""
		modified.DocumentTypes["changelog"] = p
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"changelog"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("invalid rule_set", func(t *testing.T) {
		modified := testStandard()
		p := modified.DocumentTypes["changelog"]
		p.RuleSet = "other"
		modified.DocumentTypes["changelog"] = p
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"changelog"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("graph on none", func(t *testing.T) {
		modified := testStandard()
		p := modified.DocumentTypes["changelog"]
		p.Graph = true
		modified.DocumentTypes["changelog"] = p
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"changelog"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("non-spec empty globs", func(t *testing.T) {
		modified := testStandard()
		p := modified.DocumentTypes["changelog"]
		p.Globs = nil
		modified.DocumentTypes["changelog"] = p
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"changelog"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("spec without top-level globs", func(t *testing.T) {
		modified := testStandard()
		modified.GovernedGlobs = nil
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"spec"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("empty dependable statuses", func(t *testing.T) {
		modified := testStandard()
		modified.DependableStatuses = nil
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"changelog"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("duplicate types defensive", func(t *testing.T) {
		if _, err := buildRunPlan(std, RepoConfig{Mode: multiMode, Types: []string{"changelog", "changelog"}}); err == nil {
			t.Fatal("want error")
		}
	})

	t.Run("idea profile validation", func(t *testing.T) {
		modified := testStandard()
		p := modified.DocumentTypes["idea"]
		p.RequiredHeadings = nil
		modified.DocumentTypes["idea"] = p
		if _, err := buildRunPlan(modified, RepoConfig{Mode: multiMode, Types: []string{"idea"}}); err == nil {
			t.Fatal("want error")
		}
	})
}

func TestStandardProfileValidation(t *testing.T) {
	mutate := func(fn func(*Standard)) Standard {
		std := testStandard()
		fn(&std)
		return std
	}
	legacyIdea := RepoConfig{Mode: legacyMode, Type: "idea"}
	legacySpec := RepoConfig{Mode: legacyMode, Type: "spec"}

	cases := []struct {
		name    string
		std     Standard
		cfg     RepoConfig
		wantErr bool
	}{
		{
			name: "legacy v1 idea without rule_set stays legal",
			std: mutate(func(s *Standard) {
				s.StandardVersion = "1"
				p := s.DocumentTypes["idea"]
				p.RuleSet = ""
				s.DocumentTypes["idea"] = p
			}),
			cfg: legacyIdea,
		},
		{
			name: "invalid rule_set in disabled profile",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["changelog"]
				p.RuleSet = "bogus"
				s.DocumentTypes["changelog"] = p
			}),
			cfg:     legacySpec,
			wantErr: true,
		},
		{
			name: "graph on disabled none profile",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["changelog"]
				p.Graph = true
				s.DocumentTypes["changelog"] = p
			}),
			cfg:     legacySpec,
			wantErr: true,
		},
		{
			name: "invalid rule_set on legacy idea profile",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["idea"]
				p.RuleSet = "bogus"
				s.DocumentTypes["idea"] = p
			}),
			cfg:     legacyIdea,
			wantErr: true,
		},
		{
			name: "graph on legacy idea profile",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["idea"]
				p.Graph = true
				s.DocumentTypes["idea"] = p
			}),
			cfg:     legacyIdea,
			wantErr: true,
		},
		{
			name: "legacy idea profile missing structure fields",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["idea"]
				p.RuleSet = ""
				p.RequiredFields = nil
				s.DocumentTypes["idea"] = p
			}),
			cfg:     legacyIdea,
			wantErr: true,
		},
		{
			name: "empty glob entry",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["changelog"]
				p.Globs = []string{"docs/specs/**/changelog.md", " "}
				s.DocumentTypes["changelog"] = p
			}),
			cfg:     legacySpec,
			wantErr: true,
		},
		{
			name: "empty globs on disabled profile",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["changelog"]
				p.Globs = nil
				s.DocumentTypes["changelog"] = p
			}),
			cfg:     legacySpec,
			wantErr: true,
		},
		{
			name: "malformed glob pattern",
			std: mutate(func(s *Standard) {
				p := s.DocumentTypes["changelog"]
				p.Globs = []string{"docs/["}
				s.DocumentTypes["changelog"] = p
			}),
			cfg:     legacySpec,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := buildRunPlan(tc.std, tc.cfg)
			if tc.wantErr && err == nil {
				t.Fatal("buildRunPlan() = nil error, want error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("buildRunPlan(): %v", err)
			}
		})
	}
}
