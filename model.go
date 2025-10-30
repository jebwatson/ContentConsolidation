package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	detailView
	editView
)

type model struct {
	state          uint
	repo           *Repo
	content        []Content
	currentContent Content
	listIndex      int
	textarea       textarea.Model
	textinput      textinput.Model
}

func NewModel(repo *Repo) model {
	content, err := repo.GetContent()
	if err != nil {
		log.Fatalf("Unable to get content: %v", err)
	}

	return model{
		state:     listView,
		repo:      repo,
		content:   content,
		textarea:  textarea.New(),
		textinput: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currentContent = Content{}
				m.state = detailView
			// ... show input
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.content)-1 {
					m.listIndex++
				}
			case "enter":
				m.currentContent = m.content[m.listIndex]
				m.textinput.SetValue(m.currentContent.Title)
				m.textinput.Focus()
				m.textinput.CursorEnd()
				m.state = detailView
			}
		case detailView:
			switch key {
			case "enter":
				title := m.textinput.Value()
				if title != "" {
					m.currentContent.Title = title
					m.textarea.SetValue(m.currentContent.Location)
					m.textarea.Focus()
					m.textarea.CursorEnd()
					m.state = editView
				}
			case "esc":
				m.state = listView
			}
		case editView:
			switch key {
			case "ctrl+s":
				location := m.textarea.Value()
				m.currentContent.Location = location

				var err error
				if err = m.repo.SaveContentEntry(m.currentContent); err != nil {
					// TODO: handle error
					return m, tea.Quit
				}

				m.content, err = m.repo.GetContent()
				if err != nil {
					// TODO: handle error
					return m, tea.Quit
				}

				m.currentContent = Content{}
				m.state = listView
			case "esc":
				m.state = listView
			}
		}
	}
	return m, tea.Batch(cmds...)
}
