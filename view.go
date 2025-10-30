package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	appNameStyle    = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	faintStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	enumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

func (m model) View() string {
	s := appNameStyle.Render("Content Consolidation") + "\n\n"

	if m.state == detailView {
		s += "Content title:\n\n"
		s += m.textinput.View() + "\n\n"
		s += faintStyle.Render("enter - save, esc - discard")
	}

	if m.state == editView {
		s += "Content location:\n\n"
		s += m.textarea.View() + "\n\n"
		s += faintStyle.Render("ctrl+s - save, esc - discard")
	}

	if m.state == listView {
		for i, c := range m.content {
			prefix := ""
			if i == m.listIndex {
				prefix = ">"
			}

			shortLocation := strings.ReplaceAll(c.Location, "\n", " ")
			if len(shortLocation) > 30 {
				shortLocation = shortLocation[:30]
			}

			s += enumeratorStyle.Render(prefix) + c.Title + " | " + faintStyle.Render(shortLocation) + "\n\n"
		}

		s += faintStyle.Render("n - new content, q - quit")
	}

	return s
}
