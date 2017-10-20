package fuz

import (
	"fmt"
	"github.com/nu11p01n73R/walker"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
)

// Fuz Modes
const (
	SEARCH = iota
	NORMAL = iota
)

// Event Type
const (
	NOP       = iota
	UPWARD    = iota
	DONWARD   = iota
	TOGGLE    = iota
	OPEN      = iota
	BACKSPACE = iota
	EXIT      = iota
)

const MAX_VIEWPORT_SIZE = 20

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
	}, []string{`\.git`, `.*\.sw.*`, `node_modules`, `vendor`})
}

func printHeader(logo string) {
	fmt.Println(logo)
	fmt.Println("-----------------------------------")
}

// Print the list of files with the cursor and search string.
func printList(files []string, cursorAt int, searchString string, mode int) {

	if len(files) == 0 {
		fmt.Println("  Oops!! No matches found.")
	}

	for i := 0; i < MAX_VIEWPORT_SIZE && i < len(files); i++ {
		cursor := "  "
		if cursorAt == i {
			cursor = " >"
		}
		fmt.Printf("%s %s\n", cursor, files[i])
	}

	currentMode := "S"
	if mode == NORMAL {
		currentMode = "N"
	}

	fmt.Println("-----------------------------------")
	fmt.Printf("[%s] %s", currentMode, searchString)
}

// Clears the screen
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func runCommand(cmd *exec.Cmd, file string) error {
	cmd.Args = append(cmd.Args, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

// Toggle the modes between search and normal
func toggleMode(curr int) int {
	if curr == SEARCH {
		return NORMAL
	} else {
		return SEARCH
	}
}

// Checks if all the characters in search
// is present in str in the same order.
// If the search is present, returns true.
// Else returns false.
func contains(str, search string) bool {
	var found, j int
	searchRune := []rune(search)
	strRune := []rune(str)

	for i := 0; i < len(searchRune); i++ {
		for ; j < len(strRune); j++ {
			if unicode.ToLower(searchRune[i]) ==
				unicode.ToLower(strRune[j]) {
				found++
				break
			}
		}
	}

	return found == len(searchRune)
}

// Filters the file array based on the input search string.
// Returns filtered array, where the search string is
// present in the file name.
func filterFiles(files []string, searchString string) []string {
	output := []string{}
	for i := 0; i < len(files); i++ {
		if contains(files[i], searchString) {
			output = append(output, files[i])
		}
	}
	return output
}

// Determins the event for key press by the user.
// The type of event depends upon the current mode.
// Some event like TOGGLE and OPEN doesn't depend on
// the mode.
// Params
//	key 	byte 	The key pressed by the user
//	mode	int	The current mode fuz is on.
// Return
//	int	The even identified based on key
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
		case 113:
			return EXIT
		}
	}

	if mode == SEARCH && key == 127 {
		return BACKSPACE
	}

	return NOP
}

// Get the size of the view port.
// This determins the total number of
// files to be shown by default.
// At max fuz shows the first 10 files
// in the list.
// Params
//	files []string Files to be displayed
// Return
//	int	Number of files in viewport
//	int	Position of the cursor initailly
func getViewPortSize(files []string) (int, int) {
	size := len(files)
	if size > MAX_VIEWPORT_SIZE {
		return MAX_VIEWPORT_SIZE, MAX_VIEWPORT_SIZE - 1
	}
	return size, size - 1
}

// Main wainting loop.
// Sets up the viewport with files.
// Waits for user input.
// Based on the key press, the current mode
// is changed. Further keypress actions depend
// upon the mode.
// Params
// 	files, list of files to be displayed on startup.
func viewPort(files []string, logo string, cmd *exec.Cmd) error {
	var searchString string
	var err error

	viewPortSize, cursorAt := getViewPortSize(files)
	fileCount := len(files)
	mode := SEARCH
	allFiles := files

	fileFlag := "-F"
	if runtime.GOOS == "darwin" {
		fileFlag = "-f"
	}

	exec.Command("stty", fileFlag, "/dev/tty", "cbreak", "min", "1").Run()
	err = exec.Command("/bin/stty", fileFlag, "/dev/tty", "-echo").Run()
	if err != nil {
		return err
	}

	char := make([]byte, 1)

keyWait:
	for {
		clearScreen()

		if len(logo) > 0 {
			printHeader(logo)
		}
		printList(files, cursorAt, searchString, mode)

		os.Stdin.Read(char)

		event := keyHandler(char[0], mode)

		switch event {
		case EXIT:
			break keyWait
		case UPWARD:
			if cursorAt == 0 {
				cursorAt = viewPortSize - 1
			} else {
				cursorAt--
			}
			break
		case DONWARD:
			cursorAt = (cursorAt + 1) % viewPortSize
			break
		case TOGGLE:
			mode = toggleMode(mode)
			break
		case OPEN:
			if cursorAt < fileCount {
				err = runCommand(cmd, files[cursorAt])
				break keyWait
			}
		case SEARCH:
			searchString += string(char)
			files = filterFiles(files, searchString)
			viewPortSize, cursorAt = getViewPortSize(files)
			fileCount = len(files)
		case BACKSPACE:
			if len(searchString) > 0 {
				searchString = searchString[:len(searchString)-1]
				files = filterFiles(allFiles, searchString)
				viewPortSize, cursorAt = getViewPortSize(files)
				fileCount = len(files)
			}
		}

	}
	return err
}

// Sets terminal echo on
func cleanUp() {
	fileFlag := "-F"
	if runtime.GOOS == "darwin" {
		fileFlag = "-f"
	}

	exec.Command("/bin/stty", fileFlag, "/dev/tty", "echo").Run()
}

// Entry point to fuz
// Walks the current directory and print out a list
// of files in viewport.
// Starts up search mode for the user to type in the
// search string.
// The files in viewport are filtered based on the
// search string.
// Args
// 	dir The directory which has to be listed.
//	logo The header logo to be printed.
//	command The command to be executed when selecting a
//	file at cursor.
func Fuz(dir string, logo string, command *exec.Cmd) error {
	defer cleanUp()

	files, err := intialWalk(dir)
	if err != nil {
		return err
	}

	err = viewPort(files, logo, command)
	if err != nil {
		return err
	}

	fmt.Println()
	return nil
}
