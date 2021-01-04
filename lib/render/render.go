package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/dienakakim/mds/lib/structs"
	"github.com/yuin/goldmark"
)

// Default text to display when goldmark fails to render markdown
const errorText = "Failed to parse markdown"

// render uses the given Goldmark instance to render the HTML.
func Render(w http.ResponseWriter, r *http.Request, gm goldmark.Markdown, templ *template.Template, config Config) {
	content, err := ioutil.ReadFile(config.FileName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := fmt.Sprintf("404: \"%s\" cannot be opened", config.FileName)
		fmt.Fprintln(w, err)
		log.Println(err)
		return
	}
	if strings.HasSuffix(config.FileName, ".md") {
		// Markdown file
		var html bytes.Buffer
		if err := gm.Convert(content, &html); err != nil {
			log.Println(errorText)
		}
		_, fileName := filepath.Split(config.FileName)
		rendered := RenderedHTML{Body: template.HTML(html.String()), Style: template.CSS(string(config.StyleBytes)), FileName: fileName}
		templ.Execute(w, rendered)
	} else {
		// Arbitrary file
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write(content)
	}
}
