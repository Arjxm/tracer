package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Event []Node

var FormatedEvent []interface{}

type Node struct {
	OnEnter  map[string]interface{} `json:"OnEnter"`
	Children []Node                 `json:"Children"`
	OnExit   map[string]interface{} `json:"OnExit"`
	OpCodes  []interface{}          `json:"OpCodes"`
}

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	callStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	staticCallStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	delegateCallStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := "Trace Viewer"
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func DisplayTree(node Node, level int) string {
	var result string
	indent := strings.Repeat("  ", level)

	onEnterFrom := node.OnEnter["From"]
	onEnterTo := node.OnEnter["To"]
	onEnterType := node.OnEnter["Type"].(string)
	onEnterValue := node.OnEnter["Value"]

	var onEnterValueStr string
	if onEnterValue == nil {
		onEnterValueStr = "%!s(<nil>)"
	} else {
		onEnterValueStr = fmt.Sprintf("%v", onEnterValue)
	}

	onExitOutput := node.OnExit["Output"]

	var style lipgloss.Style
	switch onEnterType {
	case "CALL":
		style = callStyle
	case "STATICCALL":
		style = staticCallStyle
	case "DELEGATECALL":
		style = delegateCallStyle
	default:
		style = lipgloss.NewStyle()
	}

	result += fmt.Sprintf("%s%s From: %s To: %s Value: %s -> Output: %s\n", indent, style.Render(onEnterType), onEnterFrom, onEnterTo, onEnterValueStr, onExitOutput)

	for _, child := range node.Children {
		result += DisplayTree(child, level+1)
	}

	return result
}

func Display(content string) {
	m := model{content: content}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
