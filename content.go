package main

import (
	"encoding/json"
	"fmt"
)

// L is a bilingual string: every visible text exists in Spanish and English.
type L struct{ ES, EN string }

// Content is the whole CV, loaded from content/content.json.
type Content struct {
	Name     string
	Email    string
	LinkedIn string
	GitHub   string
	Title    L
	Intro    L
	UI       map[string]L
	Sectors  []Sector
	Skills   []Skill
	Build    Build
}

// Sector is one slice of the experience disc (DevOps, Sysadmin…).
type Sector struct {
	ID      string
	Name    L
	Summary L
	Roles   []Role
	Certs   []L
	Sub     *Sub
	// Filled in by the generator, not by the JSON.
	Path  string
	Label string
	Index int
}

type Role struct {
	Company L
	Title   L
	Period  L
	Items   []L
}

// Sub is a nested block inside a sector (DevSecOps inside DevOps).
type Sub struct {
	Name  L
	Note  L
	Items []L
}

// Skill is one level bar. Col places it in one of the three skill columns.
type Skill struct {
	Name  L
	Level int
	Icons []string
	Col   int
}

// Build is the "how this site is built" section.
type Build struct {
	Text  []L
	Tools []Tool
}

type Tool struct {
	Name  string
	Icons []string
}

func loadContent(file string) (Content, error) {
	var c Content
	b, err := fsys.ReadFile("content/" + file)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, fmt.Errorf("%s: %w", file, err)
	}
	return c, nil
}
