package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const ideaConfig = `{"type":"idea"}`

func ideaContent(id, title, date, status string) string {
	return `---
id: ` + id + `
title: "` + title + `"
date: ` + date + `
tags: [产品]
status: ` + status + `
---

## 一句话描述

一句能独立讲清的话。

## 想法

具体机制。

## 为什么有意思

真实问题。

## 下一步

拿真实数据跑一次。
`
}

func hasRule(rep Report, rule string) bool {
	for _, f := range rep.Findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

func TestIdeaValidRepo(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".mdstandards.json", ideaConfig)
	writeFile(t, root, "2026/2026-11-01/ok.md", ideaContent("AABBCCDD", "测试想法", "2026-11-01", "萌芽"))

	code, out, errs := runCLI(t, root, "-phase", "1", "-format", "json")
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, out, errs)
	}
	rep := decodeReport(t, out)
	if len(rep.Findings) != 0 {
		t.Fatalf("unexpected findings: %+v", rep.Findings)
	}
}

func TestIdeaFixtureRepo(t *testing.T) {
	code, out, errs := runCLI(t, "testdata/idea-repo", "-phase", "1", "-format", "json")
	if code != 0 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, out, errs)
	}
	rep := decodeReport(t, out)
	if len(rep.Findings) != 0 {
		t.Fatalf("unexpected findings: %+v", rep.Findings)
	}
}

func TestIdeaOptionalVoiceSection(t *testing.T) {
	std, err := loadStandard(policyDir(t))
	if err != nil {
		t.Fatal(err)
	}
	var placeholder string
	for _, m := range std.DocumentTypes["idea"].PlaceholderMarkers {
		if strings.Contains(m, "outside-voice") {
			placeholder = m
			break
		}
	}
	if placeholder == "" {
		t.Fatal("standard.json has no placeholder marker for outside-voice")
	}
	base := ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽")
	cases := []struct {
		name     string
		body     string
		wantCode int
		wantRule string
	}{
		{"absent", base, 0, ""},
		{"filled", strings.Replace(base, "## 下一步", "## 外部的声音\n\n有人试过类似方案。\n\n## 下一步", 1), 0, ""},
		{"placeholder remains", strings.Replace(base, "## 下一步", "## 外部的声音\n\n"+placeholder+"\n\n## 下一步", 1), 1, "idea.placeholder_remaining"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".mdstandards.json", ideaConfig)
			writeFile(t, root, "2026/2026-11-01/a.md", tc.body)
			code, out, errs := runCLI(t, root, "-phase", "1", "-format", "json")
			if code != tc.wantCode {
				t.Fatalf("code=%d want=%d stdout=%s stderr=%s", code, tc.wantCode, out, errs)
			}
			rep := decodeReport(t, out)
			if tc.wantRule == "" {
				if len(rep.Findings) != 0 {
					t.Fatalf("unexpected findings: %+v", rep.Findings)
				}
				return
			}
			if !hasRule(rep, tc.wantRule) {
				t.Fatalf("want rule %q, got %+v", tc.wantRule, rep.Findings)
			}
		})
	}
}

func TestIdeaStructuralFindings(t *testing.T) {
	cases := []struct {
		name string
		rel  string
		body string
		want string
	}{
		{
			name: "bad id format",
			rel:  "2026/2026-11-01/a.md",
			body: ideaContent("bad-id", "标题", "2026-11-01", "萌芽"),
			want: "idea.id_format",
		},
		{
			name: "invalid status",
			rel:  "2026/2026-11-01/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026-11-01", "随便"),
			want: "idea.status_invalid",
		},
		{
			name: "empty title",
			rel:  "2026/2026-11-01/a.md",
			body: ideaContent("AABBCCDD", "   ", "2026-11-01", "萌芽"),
			want: "idea.empty_title",
		},
		{
			name: "missing heading",
			rel:  "2026/2026-11-01/a.md",
			body: strings.Replace(ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"), "## 下一步", "## 别的", 1),
			want: "idea.heading_missing",
		},
		{
			name: "template placeholder",
			rel:  "2026/2026-11-01/a.md",
			body: strings.Replace(ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"), "一句能独立讲清的话。", "<!-- 一句话说清楚这个想法是什么。写不出来就说明还没想清楚。 -->", 1),
			want: "idea.placeholder_remaining",
		},
		{
			name: "template id sentinel",
			rel:  "2026/2026-11-01/a.md",
			body: ideaContent("ABCDEFGH", "标题", "2026-11-01", "萌芽"),
			want: "idea.placeholder_remaining",
		},
		{
			name: "date directory mismatch",
			rel:  "2026/2026-11-02/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"),
			want: "idea.date_dir_mismatch",
		},
		{
			name: "year directory mismatch",
			rel:  "2025/2026-11-01/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"),
			want: "idea.year_dir_mismatch",
		},
		{
			name: "invalid calendar date",
			rel:  "2026/2026-02-31/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026-02-31", "萌芽"),
			want: "idea.date_invalid",
		},
		{
			name: "datetime instead of date",
			rel:  "2026/2026-11-01/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026-11-01T10:00:00Z", "萌芽"),
			want: "idea.date_invalid",
		},
		{
			name: "bad date format",
			rel:  "2026/2026-11-01/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026/11/01", "萌芽"),
			want: "idea.date_format",
		},
		{
			name: "two layer path",
			rel:  "2026/a.md",
			body: ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"),
			want: "idea.path_shape",
		},
		{
			name: "tags wrong type",
			rel:  "2026/2026-11-01/a.md",
			body: strings.Replace(ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"), "tags: [产品]", `tags: "产品"`, 1),
			want: "idea.tags_invalid",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".mdstandards.json", ideaConfig)
			writeFile(t, root, tc.rel, tc.body)
			code, out, errs := runCLI(t, root, "-phase", "1", "-format", "json")
			if code != 1 {
				t.Fatalf("code=%d stdout=%s stderr=%s", code, out, errs)
			}
			rep := decodeReport(t, out)
			if !hasRule(rep, tc.want) {
				t.Fatalf("want rule %q, got %+v", tc.want, rep.Findings)
			}
		})
	}
}

func TestIdeaDuplicateID(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".mdstandards.json", ideaConfig)
	writeFile(t, root, "2026/2026-11-01/a.md", ideaContent("AABBCCDD", "甲", "2026-11-01", "萌芽"))
	writeFile(t, root, "2026/2026-11-02/b.md", ideaContent("AABBCCDD", "乙", "2026-11-02", "探索中"))

	code, out, errs := runCLI(t, root, "-phase", "1", "-format", "json")
	if code != 1 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, out, errs)
	}
	rep := decodeReport(t, out)
	if !hasRule(rep, "idea.id_duplicate") {
		t.Fatalf("want idea.id_duplicate, got %+v", rep.Findings)
	}
}

func TestIdeaUnquotedDateIsNormalized(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".mdstandards.json", ideaConfig)
	writeFile(t, root, "2026/2026-11-01/a.md", ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"))

	code, out, errs := runCLI(t, root, "-dump-corpus")
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errs)
	}
	var corpus Corpus
	if err := json.Unmarshal([]byte(out), &corpus); err != nil {
		t.Fatal(err)
	}
	if len(corpus.Documents) != 1 || corpus.Documents[0].Date != "2026-11-01" {
		t.Fatalf("date did not normalize: %+v", corpus.Documents)
	}
}

func TestIdeaConfigFailsClosed(t *testing.T) {
	cases := []struct {
		name string
		cfg  string
	}{
		{"malformed json", `{`},
		{"unknown type", `{"type":"note"}`},
		{"unknown key", `{"type":"idea","globs":["20*/20*/*.md"]}`},
		{"empty type", `{}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".mdstandards.json", tc.cfg)
			code, out, errs := runCLI(t, root, "-phase", "1")
			if code != 2 {
				t.Fatalf("code=%d stdout=%s stderr=%s", code, out, errs)
			}
		})
	}
}

func TestIdeaRulesDoNotRunForSpec(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".mdstandards.json", `{"type":"spec"}`)
	writeFile(t, root, "2026/2026-11-01/a.md", ideaContent("AABBCCDD", "标题", "2026-11-01", "萌芽"))

	// Spec globs do not match a year-archived idea, so coverage fails closed
	// before any idea rule can run.
	code, out, errs := runCLI(t, root, "-phase", "1")
	if code != 2 {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, out, errs)
	}
}

func TestSpecDumpCorpusGolden(t *testing.T) {
	code, out, errs := runCLI(t, "testdata/example-repo", "-dump-corpus")
	if code != 0 {
		t.Fatalf("code=%d stderr=%s", code, errs)
	}
	golden, err := os.ReadFile("testdata/spec-corpus.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	if out != string(golden) {
		t.Fatalf("spec corpus changed:\nwant:\n%s\ngot:\n%s", golden, out)
	}
}
