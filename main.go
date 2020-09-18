// A Markdown server that uses the Goldmark engine and dark theme swappability.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	. "github.com/dienakakim/mds/lib/render"
	. "github.com/dienakakim/mds/lib/structs"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Help text
var helpText = `
Usage: ${prog} --file=FILE.md
       ${prog} --port 3000 --file=FILE.md

    --port      Port to serve from
    --dark      Display in dark theme
    --help      Show this help screen
`

// main is the driver code for the program.
func main() {
	// Flags
	help := flag.Bool("help", false, "show help")
	dark := flag.Bool("dark", false, "enable dark theme")
	mathMode := flag.Bool("math", true, "enable MathJax")
	port := flag.String("port", "8080", "server port")
	file := flag.String("file", "", "filename")
	flag.Parse()

	if *help {
		usage("")
	} else if *file == "" {
		usage("Please provide a file as an argument e.g. --file=README.md")
	}

	config := Config{DarkMode: *dark, FileName: *file, MathJax: *mathMode}

	// Create template
	var htmlTemplateBytes []byte
	if *dark {
		styleBytes := MustAsset("assets/dark.out.css")
		config.StyleBytes = styleBytes
	} else {
		styleBytes := MustAsset("assets/light.out.css")
		config.StyleBytes = styleBytes
	}
	htmlTemplateBytes = MustAsset("assets/index.gohtml")
	templ, err := template.New("md").Parse(string(htmlTemplateBytes))
	if err != nil {
		log.Fatal(err)
	}

	// goldmarkInitializer will initialize Goldmark with:
	// - GitHub Flavored Markdown
	// - MathJax
	// - Appropriate styling for given theme
	// - Allow custom HTML
	// - Auto heading ID generation
	goldmarkInitializer := func(style string) goldmark.Markdown {
		return goldmark.New(goldmark.WithExtensions(extension.GFM, mathjax.MathJax, highlighting.NewHighlighting(highlighting.WithStyle(style))), goldmark.WithRendererOptions(html.WithUnsafe()),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()))
	}
	gmLight := goldmarkInitializer("monokailight")
	gmDark := goldmarkInitializer("solarized-dark")

	// Create new ServeMux
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		if config.DarkMode {
			Render(w, r, gmDark, templ, config)
		} else {
			Render(w, r, gmLight, templ, config)
		}
	})
	sm.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		w.Header().Set("Content-Type", "image/x-icon")
		faviconIcoBytes := MustAsset("assets/favicon.ico")
		w.Write(faviconIcoBytes)
		return
	})

	// Initialize signal handler
	signals := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signals, os.Interrupt, os.Kill)
	paused := false
	go func(config *Config) {
		for {
			sig := <-signals
			switch sig {
			case os.Interrupt:
				if !paused {
					// Pause menu
					paused = true
					log.Println("Paused.")
					fmt.Println("Choose one of the below options:")
					fmt.Println("1. Toggle dark mode")
					fmt.Println("2. Change filename")
					fmt.Println("3. Exit")
					fmt.Printf("> ")
					var choiceStr string
					fmt.Scanln(&choiceStr)
					choiceStr = strings.TrimSpace(choiceStr)

					switch choiceStr {
					case "1":
						// Toggle dark mode
						config.DarkMode = !config.DarkMode
						status := "enabled"
						if config.DarkMode {
							config.StyleBytes, _ = assetsDarkOutCssBytes()
						} else {
							status = "disabled"
							config.StyleBytes, _ = assetsLightOutCssBytes()
						}
						log.Printf("Dark mode %s", status)
						paused = false
					case "2":
						// Change filename
						fmt.Println("Enter in new Markdown filename: ")
						fmt.Printf("> ")
						input := bufio.NewScanner(os.Stdin)
						if input.Scan() {
							config.FileName = strings.Trim(input.Text(), "\"")
							log.Printf("Filename changed to: \"%s\"", config.FileName)
						} else {
							log.Println("Filename unchanged")
						}
						paused = false
					case "3":
						// Exit, send signal
						signals <- os.Kill
					default:
						log.Printf("Invalid value: %s", choiceStr)
						paused = false
					}
				} else {
					// `paused` is true, so os.Interrupt received twice
					log.Fatal("Interrupt received twice. Force-closing.")
					signals <- os.Kill
				}
			case os.Kill:
				// Terminate now
				done <- true
				break
			}
		}
	}(&config)

	// Serve
	go func() {
		log.Printf("Starting server on port %s", *port)
		if err := http.ListenAndServe(":"+*port, sm); err != nil {
			log.Fatal(err)
		}
	}()

	// Block until receipt
	<-done
	// Done.
}

// usage displays the appropriate notice if the user did not specify the Markdown file to render.
func usage(note string) {
	if len(note) > 0 {
		fmt.Println("Error: " + note)
	}
	_, fileName := filepath.Split(os.Args[0])
	fmt.Println(strings.Replace(helpText, "${prog}", fileName, -1))
	os.Exit(0)
}
