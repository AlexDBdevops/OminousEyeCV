// OmnimousEyeCV: static generator for a bilingual CV site.
// It reads content/content.json, renders templates/index.html.tmpl and writes dist/,
// ready to be published on Cloudflare Pages.
package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//go:embed templates/*.tmpl static/* content/*.json content/themes/*.json content/theme.txt content/variant.txt
var fsys embed.FS

const out = "dist"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("dist/ generated")
}

func run() error {
	c, err := loadContent(envOr("CONTENT", "content.json"))
	if err != nil {
		return err
	}
	layoutSectors(c.Sectors)
	variant, err := setting("VARIANT", "content/variant.txt")
	if err != nil {
		return err
	}
	theme, err := setting("THEME", "content/theme.txt")
	if err != nil {
		return err
	}
	v := View{C: c, Disc: buildDisc(len(c.Sectors)), B: buildInfo(variant)}

	var page bytes.Buffer
	if err := renderIndex(&page, v); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(out, "fonts"), 0o755); err != nil {
		return err
	}
	if err := write("index.html", page.Bytes()); err != nil {
		return err
	}

	raw, err := fsys.ReadFile("content/themes/" + theme + ".json")
	if err != nil {
		return err
	}
	css, err := themeCSS(raw)
	if err != nil {
		return err
	}
	if err := write("theme.css", []byte(css)); err != nil {
		return err
	}

	tmpl, err := fsys.ReadFile("static/headers.txt")
	if err != nil {
		return err
	}
	if err := write("_headers", headersFile(tmpl, page.Bytes())); err != nil {
		return err
	}

	// Self-hosted fonts and the remaining static files.
	err = fs.WalkDir(fsys, "static/fonts", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		return copyStatic(p, filepath.Join("fonts", d.Name()))
	})
	if err != nil {
		return err
	}
	for _, n := range []string{"style.css", "variant-" + variant + ".css", "app.js", "favicon.svg", "favicon-32.png", "apple-touch-icon.png", "og.png"} {
		if err := copyStatic("static/"+n, n); err != nil {
			return err
		}
	}
	return nil
}

// buildInfo returns the short commit (GITHUB_SHA in CI, git locally) and today's date.
func buildInfo(variant string) BuildInfo {
	sha := os.Getenv("GITHUB_SHA")
	if sha == "" {
		if b, err := exec.Command("git", "rev-parse", "HEAD").Output(); err == nil {
			sha = strings.TrimSpace(string(b))
		} else {
			sha = "local"
		}
	}
	if len(sha) > 7 {
		sha = sha[:7]
	}
	return BuildInfo{SHA: sha, Date: time.Now().UTC().Format("2006-01-02"), Variant: variant}
}

// setting reads an environment override or, if unset, a one-word file under content/.
func setting(env, file string) (string, error) {
	if v := strings.TrimSpace(os.Getenv(env)); v != "" {
		return v, nil
	}
	b, err := fsys.ReadFile(file)
	return strings.TrimSpace(string(b)), err
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func write(name string, b []byte) error {
	return os.WriteFile(filepath.Join(out, name), b, 0o644)
}

func copyStatic(src, dst string) error {
	b, err := fsys.ReadFile(src)
	if err != nil {
		return err
	}
	return write(dst, b)
}
