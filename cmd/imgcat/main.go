package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trashhalo/imgcat"
	"github.com/trashhalo/imgcat/component"
)

const usage = `imgcat [pattern|url]

Examples:
    imgcat path/to/image.jpg
    imgcat *.jpg *.svg
    imgcat https://example.com/image.jpg
    imgcat https://dev.w3.org/SVG/tools/svgweb/samples/svg-files/couch.svg`

type imgModel struct {
	component.Model
	url string
}

func (m imgModel) View() string {
	return m.Model.View()
}

func (m imgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return m, m.Model.Redraw(uint(msg.Width), uint(msg.Height), m.url)
	default:
		fmt.Printf("hit default: %+v%+v", reflect.TypeOf(msg), msg)
		m.Model, cmd = m.Model.Update(msg)
		return m, cmd
	}

	mod, cmd := m.Model.Update(msg)
	m.Model = mod
	return m, cmd
}

func (m imgModel) Init() tea.Cmd {
	return m.Model.Init()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.String("url", "", "The URL of the image to load")
	flag.Parse()

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println(usage)
		os.Exit(0)
	}

	url := flag.Lookup("url").Value.String()

	// check if the --url flag is set
	if url != "" {
		// if so, use the URL to load the image
		fmt.Println("url flag is set", url)

		res, err := http.Get(url)
		if err != nil {
			fmt.Println("error getting image", err)
			os.Exit(1)
		}

		r := res.Body
		str, err := component.ReaderToimage(uint(3), uint(3), r)
		if err != nil {
			fmt.Println("error converting image", err)
			os.Exit(1)
		}

		fmt.Println(str)

		// image := component.New(1, 1, url)
		// p := imgModel{Model: image, url: url}
		// if err := tea.NewProgram(p).Start(); err != nil {
		// 	fmt.Printf("couldn't load image(s): %v", err)
		// 	os.Exit(1)
		// }

		os.Exit(0)
	}

	p := tea.NewProgram(imgcat.NewModel(os.Args[1:len(os.Args)]))
	p.EnterAltScreen()
	defer p.ExitAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
