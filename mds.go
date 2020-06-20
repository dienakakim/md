// A Markdown server that uses the Goldmark engine and dark theme swappability.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
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

var errorText = "Failed to parse markdown"

// RenderedHTML is the template struct used for the templating engine.
type RenderedHTML struct {
	Body     template.HTML
	Style    template.CSS
	FileName string
}

// Config saves the current configuration of this server run.
type Config struct {
	DarkMode   bool
	FileName   string
	MathJax    bool
	StyleBytes []byte
}

// main is the driver code for the program.
func main() {
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

	// Create a Goldmark instance with custom settings
	gmLight := goldmark.New(goldmark.WithExtensions(extension.GFM, mathjax.MathJax, highlighting.NewHighlighting(highlighting.WithStyle("monokailight"))), goldmark.WithRendererOptions(html.WithUnsafe()))
	gmDark := goldmark.New(goldmark.WithExtensions(extension.GFM, mathjax.MathJax, highlighting.NewHighlighting(highlighting.WithStyle("solarized-dark"))), goldmark.WithRendererOptions(html.WithUnsafe()))

	// Create new ServeMux
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		if config.DarkMode {
			render(w, r, gmDark, templ, config)
		} else {
			render(w, r, gmLight, templ, config)
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
					var choice int
					fmt.Scanf("%d\n", &choice)
					switch choice {
					case 1:
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
					case 2:
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
					case 3:
						// Exit, send signal
						signals <- os.Kill
					default:
						log.Printf("Invalid value: %d", choice)
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

// render uses the given Goldmark instance to render the HTML.
func render(w http.ResponseWriter, r *http.Request, gm goldmark.Markdown, templ *template.Template, config Config) {
	markdown, err := ioutil.ReadFile(config.FileName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := fmt.Sprintf("File \"%s\" cannot be opened", config.FileName)
		fmt.Fprintln(w, err)
		log.Println(err)
		return
	}
	var html bytes.Buffer
	if err := gm.Convert(markdown, &html); err != nil {
		log.Println(errorText)
	}
	_, fileName := filepath.Split(config.FileName)
	rendered := RenderedHTML{Body: template.HTML(html.String()), Style: template.CSS(string(config.StyleBytes)), FileName: fileName}
	templ.Execute(w, rendered)
}
