package main

import (
	"fmt"
	"math"
	"strings"
)

// The SVG uses a 600x600 viewBox centred on (300,300). Angles are in degrees,
// 0 at the top and growing clockwise, which is how the disc is read.

type Line struct{ X1, Y1, X2, Y2 float64 }

// Disc holds the precomputed geometry the template draws.
type Disc struct {
	Spokes []Line // iris fibres
	Hex    string // hexagon inside the iris (polygon points)
	Ticks  []Line // graduation marks on the outer rim
	Bounds []Line // separators between sectors
}

func pt(r, a float64) (float64, float64) {
	rad := a * math.Pi / 180
	return 300 + r*math.Sin(rad), 300 - r*math.Cos(rad)
}

// annulus returns the path of a ring segment between radii r1<r2 and angles a0<a1 (<180°).
func annulus(r1, r2, a0, a1 float64) string {
	x0, y0 := pt(r2, a0)
	x1, y1 := pt(r2, a1)
	x2, y2 := pt(r1, a1)
	x3, y3 := pt(r1, a0)
	return fmt.Sprintf("M%.2f %.2f A%.0f %.0f 0 0 1 %.2f %.2f L%.2f %.2f A%.0f %.0f 0 0 0 %.2f %.2f Z", x0, y0, r2, r2, x1, y1, x2, y2, r1, r1, x3, y3)
}

// arc returns a clockwise arc of radius r, used as a textPath for the sector labels.
func arc(r, a0, a1 float64) string {
	x0, y0 := pt(r, a0)
	x1, y1 := pt(r, a1)
	return fmt.Sprintf("M%.2f %.2f A%.0f %.0f 0 0 1 %.2f %.2f", x0, y0, r, r, x1, y1)
}

// layoutSectors splits the ring evenly between the sectors.
func layoutSectors(s []Sector) {
	span := 360 / float64(len(s))
	for i := range s {
		a0 := float64(i) * span
		s[i].Index = i
		s[i].Path = annulus(205, 285, a0+1.5, a0+span-1.5)
		s[i].Label = arc(252, a0+14, a0+span-14)
	}
}

func buildDisc(sectors int) Disc {
	var d Disc
	for i := 0; i < 24; i++ {
		a := float64(i) * 15
		x1, y1 := pt(34, a)
		x2, y2 := pt(78, a)
		d.Spokes = append(d.Spokes, Line{x1, y1, x2, y2})
	}
	hex := make([]string, 6)
	for i := range hex {
		x, y := pt(58, float64(i)*60+30)
		hex[i] = fmt.Sprintf("%.2f,%.2f", x, y)
	}
	d.Hex = strings.Join(hex, " ")
	for i := 0; i < 72; i++ {
		a := float64(i) * 5
		l := 6.0
		if i%6 == 0 {
			l = 12
		}
		x1, y1 := pt(292-l, a)
		x2, y2 := pt(292, a)
		d.Ticks = append(d.Ticks, Line{x1, y1, x2, y2})
	}
	for i := 0; i < sectors; i++ {
		a := float64(i) * 360 / float64(sectors)
		x1, y1 := pt(205, a)
		x2, y2 := pt(285, a)
		d.Bounds = append(d.Bounds, Line{x1, y1, x2, y2})
	}
	return d
}
