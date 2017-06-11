package main

import (
	"errors"
	"fmt"
	"github.com/nu11p01n73R/walker"
	"os"
	"os/exec"
	"strings"
)

const SEARCH = 0
const NORMAL = 1

// Event Type
const (
	NOP       = iota
	UPWARD    = iota
	DONWARD   = iota
	TOGGLE    = iota
	OPEN      = iota
	BACKSPACE = iota
)

// Walks the directory.
// Removes the prefix of current working dir,
// for more clear path strings
// Return
//	[]string List of formated file paths.
func intialWalk(dir string) ([]string, error) {
	prefix := dir + "/"
	return walker.Walk(dir, func(files []string) []string {
		for i, file := range files {
			files[i] = strings.TrimPrefix(file, prefix)
		}
		return files
	})
}

func printList(files []string, cursorAt int, searchString string) {
	fmt.Println()
	for i, file := range files {
		cursor := " "
		if cursorAt == i {
			cursor = ">"
		}
		fmt.Printf("%s %s\n", cursor, file)
	}
	fmt.Printf(">> %s", searchString)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func openEditor(file string) error {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		return errors.New("Cannot open the file because $EDITOR not set")
	}

	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func toggleMode(curr int) int {
	if curr == SEARCH {
		return NORMAL
	} else {
		return SEARCH
	}
}

func normalMode(char string, files []string, size int, cursorAt *int) (bool, error) {
	var err error
	var done bool

	switch char {
	case "j":
		// pressing j
		if *cursorAt == 0 {
			*cursorAt = size - 1
		} else {
			*cursorAt--
		}
		break
	case "k":
		*cursorAt = (*cursorAt + 1) % size
		break
	case "o":
		err = openEditor(files[*cursorAt])
		done = true
		break
	}
	return done, err
}

func contains(str, search string) bool {
	var j int
	for i := 0; i < len(search); i++ {
		for ; j < len(str); j++ {
			if str[j] == search[i] {
				break
			}
		}

		if j == len(str) {
			return false
		}
	}
	return true
}

func filterFiles(files []string, searchString string) []string {
	output := []string{}
	for i := 0; i < len(files); i++ { //&& contains(files[i], searchString); i++ {
		if contains(files[i], searchString) {
			output = append(output, files[i])
		}
	}
	return output
}

func keyHandler(key byte, mode int) int {
	switch key {
	case 27:
		return TOGGLE
	case 10:
		return OPEN
	}

	if mode == NORMAL {
		switch key {
		case 106:
			return DONWARD
		case 107:
			return UPWARD
		}
	}

	if mode == SEARCH && key == 127 {
		return BACKSPACE
	}

	return NOP
}

func viewPort(files []string) error {
	var err error
	var searchString string
	var done bool

	size := len(files)
	cursorAt := size - 1
	mode := SEARCH
	initialList := files

	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	char := make([]byte, 1)
userInput:
	for {
		//		clearScreen()
		//printList(files, cursorAt, searchString)

		os.Stdin.Read(char)

		event := keyHandler(char[0], mode)
		fmt.Println(event)
		continue

		fmt.Println("press", char)
		if char[0] == 27 {
			mode = toggleMode(mode)
			continue
		}

		if mode == NORMAL {
			done, err = normalMode(string(char), files, size, &cursorAt)
			if done {
				break userInput
			}
		} else {
			if char[0] == 127 {
				searchString = searchString[:len(searchString)-2]
				files = initialList
			} else {
				searchString += string(char)
			}
			files = filterFiles(files, searchString)
			size = len(files)
			cursorAt = size - 1
		}

	}
	return err
}

func cleanUp() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

// Handles any errors
func handleError(err error) {
	if err != nil {
		cleanUp()
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	pwd, err := os.Getwd()
	handleError(err)

	files, err := intialWalk(pwd)
	handleError(err)

	err = viewPort(files)
	handleError(err)

	cleanUp()
}
