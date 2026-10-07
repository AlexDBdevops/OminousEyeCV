package main

import (
	"fmt"
	"html/template"
	"strings"
)

// 5x5 pixel glyphs for the console banner printed by `echo` in the hero.
var glyphs = map[rune][5]string{
	'B': {"####.", "#...#", "####.", "#...#", "####."},
	'I': {"#####", "..#..", "..#..", "..#..", "#####"},
	'E': {"#####", "#....", "####.", "#....", "#####"},
	'N': {"#...#", "##..#", "#.#.#", "#..##", "#...#"},
	'V': {"#...#", "#...#", "#...#", ".#.#.", "..#.."},
	'D': {"####.", "#...#", "#...#", "#...#", "####."},
	'O': {".###.", "#...#", "#...#", "#...#", ".###."},
	'S': {".####", "#....", ".###.", "....#", "####."},
	'W': {"#...#", "#...#", "#.#.#", "##.##", "#...#"},
	'L': {"#....", "#....", "#....", "#....", "#####"},
	'C': {".####", "#....", "#....", "#....", ".####"},
	'M': {"#...#", "##.##", "#.#.#", "#...#", "#...#"},
}

// banner draws text with block characters, figlet style. Unknown runes are skipped.
func banner(text string) string {
	var rows [5]strings.Builder
	for i, r := range strings.ToUpper(text) {
		g, ok := glyphs[r]
		if !ok {
			continue
		}
		for y := 0; y < 5; y++ {
			if i > 0 {
				rows[y].WriteString(" ")
			}
			rows[y].WriteString(strings.NewReplacer("#", "█", ".", " ").Replace(g[y]))
		}
	}
	out := make([]string, 5)
	for y := range rows {
		out[y] = strings.TrimRight(rows[y].String(), " ")
	}
	return strings.Join(out, "\n")
}

// bannerSVG draws the same banner as pixel rects: sharp at any size and independent of
// whether the font has the block glyph. Rows are .78 tall to leave scanline gaps.
func bannerSVG(text string) template.HTML {
	lines := strings.Split(banner(text), "\n")
	w := 0
	for _, l := range lines {
		if n := len([]rune(l)); n > w {
			w = n
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg viewBox="0 0 %d 5" style="--cols:%d" role="img" aria-label="%s" shape-rendering="crispEdges">`, w, w, template.HTMLEscapeString(text))
	for y, l := range lines {
		r := []rune(l)
		for x := 0; x < len(r); {
			if r[x] != '█' {
				x++
				continue
			}
			s := x
			for x < len(r) && r[x] == '█' {
				x++
			}
			fmt.Fprintf(&sb, `<rect x="%d" y="%d" width="%d" height=".78"/>`, s, y, x-s)
		}
	}
	sb.WriteString(`</svg>`)
	return template.HTML(sb.String())
}
