package main

import (
	"fmt"
	"github.com/nu11p01n73R/fuz/src"
	"os"
	"os/exec"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get the working directory")
		os.Exit(1)
	}

	logo := `
	   __| |  |_  )
	   _|  │  │  /
  	 _|   ____│___|
		`

	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		fmt.Println("Cannot open file because $EDITOR is not set")
		os.Exit(1)
	}
	cmd := exec.Command(editor)

	fuz.Fuz(pwd, logo, cmd)
}
