package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
)

var inlineScript = regexp.MustCompile(`(?s)<script>(.*?)</script>`)

// scriptHashes returns the CSP sources ('sha256-…') for every inline <script> of the page,
// so the policy can drop 'unsafe-inline' and still run our own script.
func scriptHashes(page []byte) string {
	var hashes []string
	for _, m := range inlineScript.FindAllSubmatch(page, -1) {
		sum := sha256.Sum256(m[1])
		hashes = append(hashes, "'sha256-"+base64.StdEncoding.EncodeToString(sum[:])+"'")
	}
	return strings.Join(hashes, " ")
}

// headersFile fills the Cloudflare Pages _headers template with the script hashes.
func headersFile(tmpl, page []byte) []byte {
	return []byte(strings.ReplaceAll(string(tmpl), "__SCRIPT_HASHES__", scriptHashes(page)))
}

// themeCSS turns a theme JSON ({"bg":"#07060d",…}) into CSS custom properties, in a stable order.
func themeCSS(raw []byte) (string, error) {
	var tv map[string]string
	if err := json.Unmarshal(raw, &tv); err != nil {
		return "", err
	}
	keys := make([]string, 0, len(tv))
	for k := range tv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString(":root{")
	for _, k := range keys {
		sb.WriteString("--" + k + ":" + tv[k] + ";")
	}
	sb.WriteString("}\n")
	return sb.String(), nil
}
