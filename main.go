package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
)

type model struct {
	entries []Website
	err     error
}

func updateEntries() tea.Msg {
	entries, err := ReadContentEntries()
	if err != nil {
		return errMsg{err}
	}

	return entriesMsg(entries)
}

type entriesMsg []Website

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return updateEntries
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case entriesMsg:
		m.entries = []Website(msg)
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if len(m.entries) <= 0 {
		return ""
	}

	items := make([]string, len(m.entries))
	for index, entry := range m.entries {
		items[index] = entry.Site
	}

	l := list.New(items)
	fmt.Println(l)
	return ""
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

/*
	p := tea.NewProgram(model{entries})

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

/*
func processCommand() {
	// Check if enough arguments were provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: cc <command> [arguments]\nCommands:\n  add <site> <description>\n  list\n  update <site> <description> <id>\n  delete <id>")
	}

	// Figure out what was asked for
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 4 {
			log.Fatal("Usage: cc add <site> <description>")
		}
		CreateContentEntry(os.Args[2], os.Args[3])
	case "list":
		ReadContentEntries()
	case "update":
		if len(os.Args) < 5 {
			log.Fatal("Usage: cc update <site> <description> <id>")
		}
		id, err := bson.ObjectIDFromHex(os.Args[4])
		if err != nil {
			log.Fatal("Could not convert id argument from string to int")
		}

		website := Website{ID: id, Site: os.Args[2], Description: os.Args[3]}
		UpdateContentEntry(website)
	case "delete":
		if len(os.Args) < 3 {
			log.Fatal("Usage: cc delete <id>")
		}
		id, err := bson.ObjectIDFromHex(os.Args[2])
		if err != nil {
			log.Fatal("Could not convert id argument from string to int")
		}

		DeleteContentEntry(id)
	default:
		log.Fatal("Unknown command: ", command)
	}
}
*/
