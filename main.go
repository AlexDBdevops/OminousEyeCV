package main

import (
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed templates/*.tmpl static/* content/*.json content/themes/*.json content/theme.txt content/variant.txt
var fsys embed.FS

type L struct{ ES, EN string }

type Sector struct {
	ID      string
	Name    L
	Period  L
	Summary L
	Items   []L
	Roles   []Role
	Certs   []L
	Sub     *Sub
	Path    string
	Label   string
	Index   int
}
type Role struct {
	Company L
	Title   L
	Period  L
	Items   []L
}
type Sub struct {
	Name  L
	Note  L
	Items []L
}
type Skill struct {
	Name  L
	Level int
	Group string
	Icons []string
}
type Content struct {
	Name    string
	Email   string
	Title   L
	Intro   L
	UI      map[string]L
	Sectors []Sector
	Skills  []Skill
	Build   Build
	Groups  []Group
}
type Group struct {
	Key  string
	Name L
}
type Build struct {
	Text  []L
	Tools []Tool
}
type Tool struct {
	Name  string
	Icons []string
}
type Spoke struct{ X1, Y1, X2, Y2 float64 }
type BuildInfo struct {
	SHA, Date, Variant string
	Grouped            bool
}
type View struct {
	B      BuildInfo
	C      Content
	Spokes []Spoke
	Hex    string
	Ticks  []Spoke
	Bound  []Spoke
}

func pt(r, a float64) (float64, float64) {
	rad := a * math.Pi / 180
	return 300 + r*math.Sin(rad), 300 - r*math.Cos(rad)
}
func annulus(r1, r2, a0, a1 float64) string {
	x0, y0 := pt(r2, a0)
	x1, y1 := pt(r2, a1)
	x2, y2 := pt(r1, a1)
	x3, y3 := pt(r1, a0)
	return fmt.Sprintf("M%.2f %.2f A%.0f %.0f 0 0 1 %.2f %.2f L%.2f %.2f A%.0f %.0f 0 0 0 %.2f %.2f Z", x0, y0, r2, r2, x1, y1, x2, y2, r1, r1, x3, y3)
}
func arc(r, a0, a1 float64) string {
	x0, y0 := pt(r, a0)
	x1, y1 := pt(r, a1)
	return fmt.Sprintf("M%.2f %.2f A%.0f %.0f 0 0 1 %.2f %.2f", x0, y0, r, r, x1, y1)
}

var iconPath = regexp.MustCompile(`<path d="([^"]+)"`)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var c Content
	cf := os.Getenv("CONTENT")
	if cf == "" {
		cf = "content.json"
	}
	b, err := fsys.ReadFile("content/" + cf)
	must(err)
	must(json.Unmarshal(b, &c))
	for i := range c.Sectors {
		span := 360 / float64(len(c.Sectors))
		a0 := float64(i) * span
		c.Sectors[i].Index = i
		c.Sectors[i].Path = annulus(205, 285, a0+1.5, a0+span-1.5)
		c.Sectors[i].Label = arc(252, a0+14, a0+span-14)
	}
	v := View{C: c}
	// Variante de diseño (v2: a | b) y datos reales del build para el pie
	variant := strings.TrimSpace(os.Getenv("VARIANT"))
	if variant == "" {
		vb, err := fsys.ReadFile("content/variant.txt")
		must(err)
		variant = strings.TrimSpace(string(vb))
	}
	sha := os.Getenv("GITHUB_SHA")
	if sha == "" {
		if out, err := exec.Command("git", "rev-parse", "HEAD").Output(); err == nil {
			sha = strings.TrimSpace(string(out))
		} else {
			sha = "local"
		}
	}
	if len(sha) > 7 {
		sha = sha[:7]
	}
	v.B = BuildInfo{SHA: sha, Date: time.Now().UTC().Format("2006-01-02"), Variant: variant}
	// SKILLSORT=group: competencias juntas por tema (sin títulos), de mayor a menor nivel
	if os.Getenv("SKILLSORT") == "group" {
		v.B.Grouped = true
		order := map[string]int{}
		for i, g := range c.Groups {
			order[g.Key] = i
		}
		sort.SliceStable(v.C.Skills, func(i, j int) bool {
			a, b := v.C.Skills[i], v.C.Skills[j]
			if order[a.Group] != order[b.Group] {
				return order[a.Group] < order[b.Group]
			}
			return a.Level > b.Level
		})
	}
	for i := 0; i < 24; i++ {
		a := float64(i) * 15
		x1, y1 := pt(34, a)
		x2, y2 := pt(78, a)
		v.Spokes = append(v.Spokes, Spoke{x1, y1, x2, y2})
	}
	for i := 0; i < 6; i++ {
		x, y := pt(58, float64(i)*60+30)
		if i == 0 {
			v.Hex += fmt.Sprintf("%.2f,%.2f", x, y)
		} else {
			v.Hex += fmt.Sprintf(" %.2f,%.2f", x, y)
		}
	}
	for i := 0; i < 72; i++ {
		a := float64(i) * 5
		l := 6.0
		if i%6 == 0 {
			l = 12
		}
		x1, y1 := pt(292-l, a)
		x2, y2 := pt(292, a)
		v.Ticks = append(v.Ticks, Spoke{x1, y1, x2, y2})
	}
	for i := 0; i < len(c.Sectors); i++ {
		x1, y1 := pt(205, float64(i)*360/float64(len(c.Sectors)))
		x2, y2 := pt(285, float64(i)*360/float64(len(c.Sectors)))
		v.Bound = append(v.Bound, Spoke{x1, y1, x2, y2})
	}
	funcs := template.FuncMap{
		"l": func(x L) template.HTML {
			return template.HTML(`<span data-lang="es">` + template.HTMLEscapeString(x.ES) + `</span><span data-lang="en">` + template.HTMLEscapeString(x.EN) + `</span>`)
		},
		"ui": func(k string) template.HTML {
			x := c.UI[k]
			return template.HTML(`<span data-lang="es">` + template.HTMLEscapeString(x.ES) + `</span><span data-lang="en">` + template.HTMLEscapeString(x.EN) + `</span>`)
		},
		"en": func(x L) string { return x.EN },
		"icons": func(slugs []string) template.HTML {
			// Iconos de Simple Icons (CC0) en línea, coloreados con currentColor
			if len(slugs) == 0 {
				return template.HTML(`<svg class="ico ico-generic" viewBox="0 0 24 24" aria-hidden="true" focusable="false"><path d="M12 3l9 9-9 9-9-9z" fill="none" stroke="currentColor" stroke-width="2"/></svg>`)
			}
			var sb strings.Builder
			for _, s := range slugs {
				b, err := fsys.ReadFile("static/icons/" + s + ".svg")
				must(err)
				d := iconPath.FindSubmatch(b)
				if d == nil {
					panic("icono sin path: " + s)
				}
				sb.WriteString(`<svg class="ico" viewBox="0 0 24 24" aria-hidden="true" focusable="false"><path d="` + template.HTMLEscapeString(string(d[1])) + `"/></svg>`)
			}
			return template.HTML(sb.String())
		},
		"es":     func(x L) string { return x.ES },
		"banner": bannerSVG,
		"dots": func(level int) template.HTML {
			n := (level + 10) / 20
			var sb strings.Builder
			for i := 0; i < 5; i++ {
				if i < n {
					sb.WriteString(`<i class="on"></i>`)
				} else {
					sb.WriteString(`<i></i>`)
				}
			}
			return template.HTML(sb.String())
		},
	}
	t, err := template.New("index.html.tmpl").Funcs(funcs).ParseFS(fsys, "templates/index.html.tmpl")
	must(err)
	must(os.MkdirAll("dist", 0o755))
	f, err := os.Create("dist/index.html")
	must(err)
	must(t.Execute(f, v))
	must(f.Close())
	name := strings.TrimSpace(os.Getenv("THEME"))
	if name == "" {
		tb, err := fsys.ReadFile("content/theme.txt")
		must(err)
		name = strings.TrimSpace(string(tb))
	}
	tj, err := fsys.ReadFile("content/themes/" + name + ".json")
	must(err)
	var tv map[string]string
	must(json.Unmarshal(tj, &tv))
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
	must(os.WriteFile("dist/theme.css", []byte(sb.String()), 0o644))
	hd, err := fsys.ReadFile("static/headers.txt")
	must(err)
	// CSP sin 'unsafe-inline' para scripts: se calcula el hash de cada <script> en línea
	page, err := os.ReadFile("dist/index.html")
	must(err)
	var hashes []string
	for _, mm := range regexp.MustCompile(`(?s)<script>(.*?)</script>`).FindAllSubmatch(page, -1) {
		sum := sha256.Sum256(mm[1])
		hashes = append(hashes, "'sha256-"+base64.StdEncoding.EncodeToString(sum[:])+"'")
	}
	hd = []byte(strings.ReplaceAll(string(hd), "__SCRIPT_HASHES__", strings.Join(hashes, " ")))
	must(os.WriteFile("dist/_headers", hd, 0o644)) // cabeceras de seguridad de Cloudflare Pages
	// fuentes autoalojadas
	must(os.MkdirAll("dist/fonts", 0o755))
	must(fs.WalkDir(fsys, "static/fonts", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fsys.ReadFile(p)
		must(err)
		return os.WriteFile(filepath.Join("dist/fonts", d.Name()), b, 0o644)
	}))
	for _, n := range []string{"style.css", "app.js", "favicon.svg", "favicon-32.png", "apple-touch-icon.png", "og.png", "variant-" + variant + ".css"} {
		d, err := fsys.ReadFile("static/" + n)
		must(err)
		must(os.WriteFile(filepath.Join("dist", n), d, 0o644))
	}
	fmt.Println("dist/ generated")
}
