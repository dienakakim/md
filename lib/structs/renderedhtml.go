package structs

import "html/template"

// RenderedHTML is the template struct used for the templating engine.
type RenderedHTML struct {
	Body     template.HTML
	Style    template.CSS
	FileName string
}
