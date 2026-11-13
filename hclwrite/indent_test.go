// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package hclwrite

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestFormatWithIndent(t *testing.T) {
	tests := map[string]struct{ input, want string }{
		"blocks":      {"outer {\ninner {\na=1\n}\n}\n", "outer {\n    inner {\n        a = 1\n    }\n}\n"},
		"collections": {"a = {\nb = [\n1,\n2,\n]\n}\n", "a = {\n    b = [\n        1,\n        2\n    ]\n}\n"},
		"comments":    {"outer {\n# Keep this.\na=1 # Inline.\n}\n", "outer {\n    # Keep this.\n    a = 1 # Inline.\n}\n"},
		"heredoc":     {"outer {\na = <<EOF\n  Keep { this } exactly.\nEOF\n}\n", "outer {\n    a = <<EOF\n  Keep { this } exactly.\nEOF\n}\n"},
		"template":    {"outer {\na=" + `"hello ${var.name}"` + "\n}\n", "outer {\n    a = " + `"hello ${var.name}"` + "\n}\n"},
		"empty":       {"", ""},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := FormatWithIndent([]byte(test.input), 4)
			if string(got) != test.want {
				t.Fatalf("Unexpected output.\nGot %q\nWant %q", got, test.want)
			}
			if !bytes.Equal(FormatWithIndent(got, 4), got) {
				t.Fatal("Formatting must be idempotent.")
			}
			if !bytes.Equal(FormatWithIndent([]byte(test.input), 2), Format([]byte(test.input))) {
				t.Fatal("Two-space formatting must match Format.")
			}
			_, diags := hclsyntax.ParseConfig(got, "test.hcl", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatal(diags)
			}
		})
	}
}

func TestFormatWithIndentWidths(t *testing.T) {
	for _, width := range []int{1, 2, 4, 8} {
		got := FormatWithIndent([]byte("outer {\na=1\n}\n"), width)
		want := "outer {\n" + strings.Repeat(" ", width) + "a = 1\n}\n"
		if string(got) != want {
			t.Fatalf("Unexpected output for width %d.\nGot %q\nWant %q", width, got, want)
		}
	}
}

func TestFormatWithIndentInvalidWidth(t *testing.T) {
	for _, width := range []int{-1, 0} {
		t.Run(strings.Repeat("negative", -width)+"width", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("Expected a panic for a non-positive width.")
				}
			}()
			FormatWithIndent(nil, width)
		})
	}
}
