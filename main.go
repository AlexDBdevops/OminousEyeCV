package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed templates/*.tmpl static/* content/*.json content/themes/*.json content/theme.txt
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
}
type Content struct {
	Name    string
	Email   string
	Title   L
	Intro   L
	UI      map[string]L
	Sectors []Sector
	Skills  []Skill
}
type Spoke struct{ X1, Y1, X2, Y2 float64 }
type View struct {
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
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var c Content
	b, err := fsys.ReadFile("content/content.json")
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
		"es": func(x L) string { return x.ES },
	}
	t, err := template.New("index.html.tmpl").Funcs(funcs).ParseFS(fsys, "templates/index.html.tmpl")
	must(err)
	must(os.MkdirAll("dist", 0o755))
	f, err := os.Create("dist/index.html")
	must(err)
	defer f.Close()
	must(t.Execute(f, v))
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
	must(os.WriteFile("dist/_headers", hd, 0o644)) // cabeceras de seguridad de Cloudflare Pages
	for _, n := range []string{"style.css", "app.js"} {
		d, err := fsys.ReadFile("static/" + n)
		must(err)
		must(os.WriteFile(filepath.Join("dist", n), d, 0o644))
	}
	fmt.Println("dist/ generated")
}
