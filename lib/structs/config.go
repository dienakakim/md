package structs

// Config saves the current configuration of this server run.
type Config struct {
	DarkMode   bool
	FileName   string
	MathJax    bool
	StyleBytes []byte
}
