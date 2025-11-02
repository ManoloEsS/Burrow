// /*
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>
// */
// package main
//
// import (
// 	"fmt"
// 	"os"
//
// 	"github.com/ManoloEsS/Burrow/cli"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	"github.com/charmbracelet/bubbles/viewport"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )
//
// func main() {
// 	// configPath := os.ExpandEnv()
// 	// cmd.Execute()
// 	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("There has been an error: %v", err)
// 		os.Exit(1)
// 	}
// }
//
// type Model struct {
// 	title     string
// 	textInput textinput.Model
// 	viewport  viewport.Model
// 	ready     bool
// 	err       error
// 	body      string
// }
//
// func initialModel() Model {
// 	ti := textinput.New()
// 	ti.Placeholder = "Url"
// 	ti.Focus()
// 	ti.Width = 40
//
// 	return Model{
// 		textInput: ti,
// 		err:       nil,
// 		body:      "",
// 	}
// }
//
// func (m Model) Init() tea.Cmd {
//
// 	return textinput.Blink
// }
//
// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var (
// 		cmd  tea.Cmd
// 		cmds []tea.Cmd
// 	)
//
// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		headerHeight := 5
// 		footerHeight := 2
// 		verticalMarginHeight := headerHeight + footerHeight
//
// 		if !m.ready {
// 			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
// 			m.viewport.YPosition = headerHeight
// 			m.ready = true
// 		} else {
// 			m.viewport.Width = msg.Width
// 			m.viewport.Height = msg.Height - verticalMarginHeight
// 		}
//
// 	case tea.KeyMsg:
// 		if m.textInput.Focused() {
// 			switch msg.Type {
// 			case tea.KeyCtrlC, tea.KeyEsc:
// 				return m, tea.Quit
// 			case tea.KeyEnter:
// 				url := m.textInput.Value()
// 				m.body = ""
// 				m.err = nil
// 				return m, cli.GetRequest(url)
// 			}
// 		} else {
// 			switch msg.String() {
// 			case "q", "ctrl+c", "esc":
// 				return m, tea.Quit
// 			case "i":
// 				m.textInput.Focus()
// 				return m, textinput.Blink
// 			}
// 		}
//
// 	case cli.RequestMsg:
// 		if msg.Err != nil {
// 			m.err = msg.Err
// 		}
//
// 		m.body = msg.Body
// 		m.viewport.SetContent(m.body)
// 		m.textInput.SetValue("")
// 		m.textInput.Blur()
// 		return m, nil
//
// 	}
//
// 	if !m.textInput.Focused() {
// 		m.viewport, cmd = m.viewport.Update(msg)
// 		cmds = append(cmds, cmd)
// 	} else {
// 		m.textInput, cmd = m.textInput.Update(msg)
// 		cmds = append(cmds, cmd)
// 	}
//
// 	return m, tea.Batch(cmds...)
// }
//
// func (m Model) View() string {
// 	if !m.ready {
// 		return "\n  Initializing..."
// 	}
//
// 	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
// 	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
//
// 	header := headerStyle.Render("Make Request") + "\n\n" + m.textInput.View()
//
// 	var help string
// 	if m.textInput.Focused() {
// 		help = helpStyle.Render("enter: submit • esc: quit")
// 	} else {
// 		help = helpStyle.Render("↑/↓: scroll • i: input mode • q/esc: quit")
// 	}
//
// 	var content string
// 	if m.err != nil {
// 		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
// 		content = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
// 	} else if m.body != "" {
// 		content = m.viewport.View()
// 	} else {
// 		content = helpStyle.Render("Enter a URL and press enter to make a request")
// 	}
//
// 	return fmt.Sprintf("%s\n\n%s\n\n%s", header, content, help)
// }
