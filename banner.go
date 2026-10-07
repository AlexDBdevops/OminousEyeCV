package main

import "strings"

// Letras de 5x5 para el cartel de consola (echo) de la variante A.
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

// banner dibuja el texto con bloques, como lo imprimiría una herramienta tipo figlet.
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
