// A Markdown server that uses yuin's goldmark with GFM extensions.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

// Help text
var helpText = `
Usage: ${prog} --file=FILE.md
       ${prog} --port 3000 --file=FILE.md

    --port		Port to serve from
    --dark      Display in dark theme
	--help		Show this help screen
`

var errorText = "Failed to parse markdown"

type RenderedHTML struct {
	Body template.HTML
}

// main is the driver code for the program.
func main() {
	help := flag.Bool("help", false, "show help")
	dark := flag.Bool("dark", false, "enable dark theme")
	port := flag.String("port", "8080", "server port")
	file := flag.String("file", "", "filename")
	flag.Parse()

	if *help {
		usage("")
	} else if *file == "" {
		usage("Please provide a file as an argument e.g. --file=README.md")
	}

	// Create template
	var htmlTemplateBytes []byte
	if *dark {
		log.Fatal("Not implemented!")
	} else {
		htmlTemplateB, err := Asset("template.html")
		if err != nil {
			log.Fatal(err)
		}
		htmlTemplateBytes = htmlTemplateB
	}
	htmlTemplate := string(htmlTemplateBytes)
	templ, err := template.New("md").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// Create a goldmark instance with GFM extensions
	gm := goldmark.New(goldmark.WithExtensions(extension.GFM, highlighting.NewHighlighting(highlighting.WithStyle("monokailight"))))

	// Create new ServeMux
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, gm, *file, templ)
	})

	// Serve
	log.Printf("Starting server on port %s", *port)
	if err := http.ListenAndServe(":"+*port, sm); err != nil {
		log.Fatal(err)
	}
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

// render uses the given goldmark instance to render the HTML.
func render(w http.ResponseWriter, r *http.Request, gm goldmark.Markdown, filename string, templ *template.Template) {
	markdown, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var html bytes.Buffer
	if err := gm.Convert(markdown, &html); err != nil {
		log.Println(errorText)
	}
	rendered := RenderedHTML{Body: template.HTML(html.String())}
	templ.Execute(w, rendered)
}
