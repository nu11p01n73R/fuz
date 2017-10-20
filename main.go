package main

import (
	"errors"
	"fmt"
	"github.com/nu11p01n73R/fuz/src"
	"os"
	"os/exec"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	pwd, err := os.Getwd()
	handleError(err)

	logo := `
	   __| |  |_  )
	   _|  │  │  /
  	 _|   ____│___|
		`

	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		handleError(
			errors.New("Cannot open file because $EDITOR is not set"))

	}
	cmd := exec.Command(editor)

	err = fuz.Fuz(pwd, logo, cmd)
	handleError(err)
}
