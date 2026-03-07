package scrubber

import (
	"testing"
)

func TestScrub(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string literal",
			input:    "SELECT * FROM users WHERE email = 'test@example.com'",
			expected: "SELECT * FROM users WHERE email = ?",
		},
		{
			name:     "numeric literal",
			input:    "SELECT * FROM orders WHERE id = 42",
			expected: "SELECT * FROM orders WHERE id = ?",
		},
		{
			name:     "multiple literals",
			input:    "SELECT * FROM users WHERE age > 25 AND name = 'Alice'",
			expected: "SELECT * FROM users WHERE age > ? AND name = ?",
		},
		{
			name:     "IN clause",
			input:    "SELECT * FROM users WHERE id IN (1, 2, 3, 4, 5)",
			expected: "SELECT * FROM users WHERE id IN (?, ?, ?, ?, ?)",
		},
		{
			name:     "already parameterized",
			input:    "SELECT * FROM users WHERE id = $1 AND name = $2",
			expected: "SELECT * FROM users WHERE id = $? AND name = $?",
		},
		{
			name:     "decimal numbers",
			input:    "SELECT * FROM products WHERE price > 19.99",
			expected: "SELECT * FROM products WHERE price > ?",
		},
		{
			name:     "whitespace normalization",
			input:    "SELECT *   FROM   users\n\tWHERE id = 1",
			expected: "SELECT * FROM users WHERE id = ?",
		},
		{
			name:     "empty string literal",
			input:    "SELECT * FROM users WHERE name = ''",
			expected: "SELECT * FROM users WHERE name = ?",
		},
		{
			name:     "timestamp literal",
			input:    "SELECT * FROM logs WHERE created_at > '2024-01-01 00:00:00'",
			expected: "SELECT * FROM logs WHERE created_at > ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Scrub(tt.input)
			if got != tt.expected {
				t.Errorf("Scrub(%q)\n  got:  %q\n  want: %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFingerprint(t *testing.T) {
	q1 := "SELECT * FROM users WHERE id = ?"
	q2 := "select * from users where id = ?"

	f1 := Fingerprint(q1)
	f2 := Fingerprint(q2)

	if f1 != f2 {
		t.Errorf("Fingerprint should be case-insensitive: %q != %q", f1, f2)
	}

	if len(f1) != 64 {
		t.Errorf("Fingerprint should be 64 hex chars, got %d", len(f1))
	}

	f3 := Fingerprint("SELECT * FROM orders WHERE id = ?")
	if f1 == f3 {
		t.Error("Different queries should have different fingerprints")
	}
}
