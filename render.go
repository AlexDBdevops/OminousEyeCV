package main

import (
	"html/template"
	"io"
	"regexp"
	"strings"
)

var (
	iconViewBox = regexp.MustCompile(`viewBox="([^"]+)"`)
	iconInner   = regexp.MustCompile(`(?s)<svg[^>]*>(.*)</svg>`)
	iconTitle   = regexp.MustCompile(`(?s)<title>.*?</title>`)
	iconFill    = regexp.MustCompile(`\sfill="#[0-9A-Fa-f]{3,6}"`)
)

const genericIcon = `<svg class="ico ico-generic" viewBox="0 0 24 24" aria-hidden="true" focusable="false"><path d="M12 3l9 9-9 9-9-9z" fill="none" stroke="currentColor" stroke-width="2"/></svg>`

// BuildInfo is shown in the footer: the real commit and date of each deploy.
type BuildInfo struct {
	SHA, Date, Variant string
}

// View is everything the template needs.
type View struct {
	B    BuildInfo
	C    Content
	Disc Disc
}

// bilingual renders both languages; CSS shows the one matching <html lang>.
func bilingual(x L) template.HTML {
	return template.HTML(`<span data-lang="es">` + template.HTMLEscapeString(x.ES) + `</span><span data-lang="en">` + template.HTMLEscapeString(x.EN) + `</span>`)
}

// iconSVG inlines an icon from static/icons, stripping its colours so it follows currentColor.
func iconSVG(slug string) (string, error) {
	b, err := fsys.ReadFile("static/icons/" + slug + ".svg")
	if err != nil {
		return "", err
	}
	vb := "0 0 24 24"
	if m := iconViewBox.FindSubmatch(b); m != nil {
		vb = string(m[1])
	}
	inner := iconInner.FindSubmatch(b)
	if inner == nil {
		return "", &iconError{slug}
	}
	body := iconFill.ReplaceAll(iconTitle.ReplaceAll(inner[1], nil), nil)
	return `<svg class="ico" viewBox="` + vb + `" aria-hidden="true" focusable="false">` + string(body) + `</svg>`, nil
}

type iconError struct{ slug string }

func (e *iconError) Error() string { return "invalid icon: " + e.slug }

func templateFuncs(ui map[string]L) template.FuncMap {
	return template.FuncMap{
		"l":      bilingual,
		"ui":     func(k string) template.HTML { return bilingual(ui[k]) },
		"en":     func(x L) string { return x.EN },
		"es":     func(x L) string { return x.ES },
		"banner": bannerSVG,
		"slice1": func(s string) []string { return []string{s} },
		"cols":   func() []int { return []int{0, 1, 2} },
		"icons": func(slugs []string) (template.HTML, error) {
			if len(slugs) == 0 {
				return template.HTML(genericIcon), nil
			}
			var sb strings.Builder
			for _, s := range slugs {
				svg, err := iconSVG(s)
				if err != nil {
					return "", err
				}
				sb.WriteString(svg)
			}
			return template.HTML(sb.String()), nil
		},
	}
}

func renderIndex(w io.Writer, v View) error {
	t, err := template.New("index.html.tmpl").Funcs(templateFuncs(v.C.UI)).ParseFS(fsys, "templates/index.html.tmpl")
	if err != nil {
		return err
	}
	return t.Execute(w, v)
}
