package main

import (
	"encoding/json"
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

func displayTree(node Node, level int) string {
	var result string
	var formated map[string]interface{}
	// indent := strings.Repeat("  ", level)

	onEnterFrom := node.OnEnter["From"]
	onEnterTo := node.OnEnter["To"]
	onEnterType := node.OnEnter["Type"]
	onEnterValue := node.OnEnter["Value"]
	onEnterInput := node.OnEnter["Input"]

	onExitOutput := node.OnExit["Output"]

	//var style lipgloss.Style
	//switch onEnterType {
	//case "CALL":
	//	style = callStyle
	//case "STATICCALL":
	//	style = staticCallStyle
	//case "DELEGATECALL":
	//	style = delegateCallStyle
	//default:
	//	style = lipgloss.NewStyle()
	//}

	//fmt.Printf("%s%s\n", indent, style.Render(fmt.Sprintf("%s -> %s", onEnterFrom, onEnterTo)))

	formated = map[string]interface{}{
		"type":   onEnterType,
		"from":   onEnterFrom,
		"to":     onEnterTo,
		"input":  onEnterInput,
		"value":  onEnterValue,
		"output": onExitOutput,
	}

	// Append the formatted event to the file
	err := appendEventToFile("events.json", formated)
	if err != nil {
		return ""
	}
	if err != nil {
		return ""
	}

	for _, child := range node.Children {
		result += displayTree(child, level+1)
	}

	return result
}

func main() {
	file, err := os.ReadFile("./trace.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var trace Event
	err = json.Unmarshal(file, &trace)
	if err != nil {
		fmt.Println(err)
		return
	}

	var content string
	for _, node := range trace {
		content += displayTree(node, 0)
		content += "\n"
	}

	m := model{content: content}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	//abiFragments, err := decoder.ParseContractABI(decoder.ABI)
	//fun := decoder.GetFuncFragment(abiFragments)
	//
	//if err != nil {
	//	fmt.Println("Error parsing contract ABI:", err)
	//	return
	//}
	//
	//// Print the parsed ABI fragments
	//for _, fragment := range *fun {
	//	fmt.Printf("Name: %s\n", fragment.Name)
	//	fmt.Printf("Type: %s\n", fragment.Type)
	//	fmt.Printf("Inputs: %+v\n", fragment.Inputs)
	//	fmt.Printf("Outputs: %+v\n", fragment.Outputs)
	//	fmt.Printf("Constant: %t\n", fragment.Constant)
	//	fmt.Printf("Payable: %t\n", fragment.Payable)
	//	fmt.Printf("Anonymous: %t\n", fragment.Anonymous)
	//	fmt.Printf("StateMutability: %s\n", fragment.StateMutability)
	//	fmt.Println("------------------------")
	//}
}

func appendEventToFile(filename string, event map[string]interface{}) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := file.Write(eventJSON); err != nil {
		return err
	}
	if _, err := file.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
