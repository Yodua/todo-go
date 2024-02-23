package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor int
    lists []string
    items map[int][]string
    selected int
    mode string
}

func initialModel() model {
    files, err := os.ReadDir("./lists/")
    if err != nil {
        panic(err)
    }

    lists := []string{}
    items := make(map[int][]string)
    for i, f := range files {
        lists = append(lists, f.Name())

        file, err := os.Open(fmt.Sprintf("./lists/%s", f.Name()))
        if err != nil {
            panic(err)
        }
        defer file.Close()
        
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            items[i] = append(items[i], strings.Trim(scanner.Text(), "\n"))
        }
    }

    return model {
        lists: lists,
        items: items,
        mode: "list",
        cursor: 0,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit

        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        case "down", "j":
            if m.mode == "list" && m.cursor < len(m.lists)-1 ||
            m.mode == "item" && m.cursor < len(m.items[m.selected])-1 {
                m.cursor++
            }

        case "enter", " ":
            if m.mode == "list" {
                m.mode = "item"
                m.selected = m.cursor
                m.cursor = 0
            }

        case "b":
            if m.mode == "item" {
                m.mode = "list"
                m.cursor = 0
            }

        case "d":
            if m.mode == "item" {
                m.items[m.selected] = append(m.items[m.selected][:m.cursor], m.items[m.selected][m.cursor+1:]...)
                if m.cursor != 0 {
                    m.cursor--
                }
            }
        }
    }

    return m, nil
}

func (m model) View() string {
    s := ""
    switch m.mode {
    case "list":
        s += "Your lists:\n"
        for i, list := range m.lists {
            cursor := " "
            if i == m.cursor {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, list)
        }

    case "item":
        s+= fmt.Sprintf("%s:\n", m.lists[m.selected])
        for i, item := range m.items[m.selected] {
            cursor := " "
            if i == m.cursor {
                cursor = ">"
            }

            s += fmt.Sprintf("%s %s\n", cursor, item)
        }
    }

    s += "\nPress q to quit or b to go back to the list selection\n"

    return s
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
