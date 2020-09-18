package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	. "github.com/dienakakim/mds/lib/structs"
	"github.com/yuin/goldmark"
)

// Default text to display when goldmark fails to render markdown
const errorText = "Failed to parse markdown"

// render uses the given Goldmark instance to render the HTML.
func Render(w http.ResponseWriter, r *http.Request, gm goldmark.Markdown, templ *template.Template, config Config) {
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
