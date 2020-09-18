package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Currently supported archs and os's
var (
	ARCHS = []string{"386", "amd64", "arm", "arm64"}
	OS    = []string{"linux", "windows", "darwin"}
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Program name and version number required, e.g. \"mds\" and \"v1.1.0\"")
		os.Exit(1)
	}

	for _, a := range ARCHS {
		for _, o := range OS {
			var suffix string
			if o == "windows" {
				suffix = ".exe"
			} else {
				suffix = ""
			}

			fileName := fmt.Sprintf("%s-%s-%s-%s%s", os.Args[1], os.Args[2], a, o, suffix)
			fmt.Printf("Building %s\n", fileName)

			build := exec.Command(fmt.Sprintf("go"), "build", "-o", fileName)
			build.Env = append(os.Environ(), "GOARCH="+a, "GOOS="+o)
			build.Run()
		}
	}
}
