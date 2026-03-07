package scrubber

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	singleQuotedRe = regexp.MustCompile(`'[^']*'`)
	numericRe      = regexp.MustCompile(`\b\d+(\.\d+)?\b`)
	hexRe          = regexp.MustCompile(`(?i)(E?)'\\x[0-9a-f]+'`)
	whitespaceRe   = regexp.MustCompile(`\s+`)
)

// Scrub replaces all literal values in a SQL query with placeholder ?.
// Defensive second pass — pg_stat_statements already normalizes most queries.
func Scrub(query string) string {
	q := hexRe.ReplaceAllString(query, "?")
	q = singleQuotedRe.ReplaceAllString(q, "?")
	q = numericRe.ReplaceAllString(q, "?")
	q = whitespaceRe.ReplaceAllString(strings.TrimSpace(q), " ")
	return q
}

// Fingerprint returns a SHA-256 hex digest of the scrubbed query.
func Fingerprint(scrubbedQuery string) string {
	h := sha256.Sum256([]byte(strings.ToLower(scrubbedQuery)))
	return hex.EncodeToString(h[:])
}
