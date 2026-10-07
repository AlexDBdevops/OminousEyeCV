package main

import (
	"bytes"
	"io/fs"
	"math"
	"strings"
	"testing"
)

func TestBannerShape(t *testing.T) {
	lines := strings.Split(banner("WELCOME"), "\n")
	if len(lines) != 5 {
		t.Fatalf("banner has %d rows, want 5", len(lines))
	}
	// 7 letters of 5 columns plus 6 one-column gaps.
	if w := len([]rune(lines[0])); w != 41 {
		t.Errorf("first row is %d wide, want 41", w)
	}
	if banner("?!") != "\n\n\n\n" {
		t.Errorf("unknown runes should be skipped")
	}
}

func TestBannerSVG(t *testing.T) {
	svg := string(bannerSVG("BIENVENIDOS"))
	if !strings.Contains(svg, `viewBox="0 0 65 5"`) {
		t.Errorf("unexpected viewBox in %q", svg[:80])
	}
	if !strings.Contains(svg, `aria-label="BIENVENIDOS"`) {
		t.Errorf("banner must keep an accessible label")
	}
	if n := strings.Count(svg, "<rect"); n < 30 {
		t.Errorf("only %d rects drawn", n)
	}
}

func TestPointsOnCircle(t *testing.T) {
	x, y := pt(100, 0)
	if x != 300 || y != 200 {
		t.Errorf("pt(100,0) = (%v,%v), want (300,200)", x, y)
	}
	x, y = pt(100, 90)
	if math.Abs(x-400) > 1e-9 || math.Abs(y-300) > 1e-9 {
		t.Errorf("pt(100,90) = (%v,%v), want (400,300)", x, y)
	}
}

func TestLayoutSectors(t *testing.T) {
	s := []Sector{{ID: "a"}, {ID: "b"}}
	layoutSectors(s)
	for i, sec := range s {
		if sec.Index != i || !strings.HasPrefix(sec.Path, "M") || !strings.HasSuffix(sec.Path, "Z") {
			t.Errorf("sector %d badly laid out: %+v", i, sec)
		}
	}
	if d := buildDisc(2); len(d.Spokes) != 24 || len(d.Ticks) != 72 || len(d.Bounds) != 2 {
		t.Errorf("unexpected disc geometry: %d spokes, %d ticks, %d bounds", len(d.Spokes), len(d.Ticks), len(d.Bounds))
	}
}

func TestScriptHashes(t *testing.T) {
	// sha256("a") in base64.
	got := scriptHashes([]byte(`<script>a</script><script src="x.js"></script>`))
	if got != "'sha256-ypeBEsobvcr6wjGzmiPcTaeG7/gUfE5yuYB3ha/uSLs='" {
		t.Errorf("scriptHashes = %s", got)
	}
}

func TestThemeCSSIsSorted(t *testing.T) {
	css, err := themeCSS([]byte(`{"b":"2","a":"1"}`))
	if err != nil || css != ":root{--a:1;--b:2;}\n" {
		t.Errorf("themeCSS = %q, %v", css, err)
	}
}

func TestContentIsComplete(t *testing.T) {
	c, err := loadContent("content.json")
	if err != nil {
		t.Fatal(err)
	}
	if c.Name == "" || c.Email == "" || len(c.Sectors) == 0 || len(c.Skills) == 0 {
		t.Fatal("content is missing basic fields")
	}
	check := func(where string, x L) {
		if strings.TrimSpace(x.ES) == "" || strings.TrimSpace(x.EN) == "" {
			t.Errorf("%s: missing translation %+v", where, x)
		}
	}
	for _, s := range c.Sectors {
		check(s.ID+" name", s.Name)
		for _, r := range s.Roles {
			for _, it := range r.Items {
				check(s.ID+" item", it)
			}
		}
	}
	for _, sk := range c.Skills {
		check("skill", sk.Name)
		if sk.Level < 0 || sk.Level > 100 || sk.Col < 0 || sk.Col > 2 {
			t.Errorf("skill %s out of range: level %d, col %d", sk.Name.EN, sk.Level, sk.Col)
		}
		for _, ic := range sk.Icons {
			if _, err := iconSVG(ic); err != nil {
				t.Errorf("skill %s: %v", sk.Name.EN, err)
			}
		}
	}
}

func TestAllIconsParse(t *testing.T) {
	entries, err := fs.ReadDir(fsys, "static/icons")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".svg") {
			continue
		}
		svg, err := iconSVG(strings.TrimSuffix(e.Name(), ".svg"))
		if err != nil {
			t.Error(err)
		} else if strings.Contains(svg, `fill="#`) {
			t.Errorf("%s keeps a hard-coded colour", e.Name())
		}
	}
}

func TestRenderIndex(t *testing.T) {
	c, err := loadContent("content.json")
	if err != nil {
		t.Fatal(err)
	}
	layoutSectors(c.Sectors)
	var b bytes.Buffer
	if err := renderIndex(&b, View{C: c, Disc: buildDisc(len(c.Sectors)), B: BuildInfo{SHA: "test", Date: "2026-01-01", Variant: "a"}}); err != nil {
		t.Fatal(err)
	}
	page := b.String()
	for _, want := range []string{"./AlexDBdevopsCV.sh", c.Email, "build test", `id="art"`} {
		if !strings.Contains(page, want) {
			t.Errorf("rendered page lacks %q", want)
		}
	}
	if scriptHashes(b.Bytes()) == "" {
		t.Error("page should have an inline script whose hash goes into the CSP")
	}
}
