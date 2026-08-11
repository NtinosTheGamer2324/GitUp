package git

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gitup/src/helper"

	gitlib "github.com/go-git/go-git/v5"
)

type fileEntry struct {
	path     string
	status   gitlib.StatusCode
	selected bool
}

// SelectFilesInteractively shows the user every changed file in the
// repository and lets them toggle which ones should go into the commit.
//
// It returns:
//   - a non-empty file list if the user confirmed a selection
//   - (nil, nil) if the user cancelled, or there was nothing to select
//   - a non-nil error only on an actual failure (repo/status read errors)
func SelectFilesInteractively(folder string) ([]string, error) {
	repo, err := gitlib.PlainOpen(folder)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	if len(status) == 0 {
		helper.Log("There are no changes to select.")
		return nil, nil
	}

	entries := make([]*fileEntry, 0, len(status))
	for file, fileStatus := range status {
		entries = append(entries, &fileEntry{
			path:     file,
			status:   fileStatus.Worktree,
			selected: true,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].path < entries[j].path
	})

	reader := bufio.NewReader(os.Stdin)

	for {
		printFileSelection(entries)

		fmt.Print("\033[32mNumbers to toggle, 'a' all, 'n' none, 'c' commit, 'q' cancel:\033[0m ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch strings.ToLower(line) {
		case "a":
			setAll(entries, true)
			continue
		case "n":
			setAll(entries, false)
			continue
		case "c":
			selected := selectedPaths(entries)
			if len(selected) == 0 {
				helper.LogFail("No files selected. Select at least one file, or 'q' to cancel.")
				continue
			}
			return selected, nil
		case "q":
			helper.Log("Selection cancelled.")
			return nil, nil
		}

		if !toggleFromInput(entries, line) {
			helper.LogFail("Please enter file numbers (e.g. 1 3 4), or one of: a, n, c, q")
		}
	}
}

func toggleFromInput(entries []*fileEntry, line string) bool {
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == ',' || r == ' '
	})

	if len(parts) == 0 {
		return false
	}

	matchedAny := false

	for _, part := range parts {
		idx, err := strconv.Atoi(part)
		if err != nil || idx < 1 || idx > len(entries) {
			continue
		}

		entries[idx-1].selected = !entries[idx-1].selected
		matchedAny = true
	}

	return matchedAny
}

func setAll(entries []*fileEntry, selected bool) {
	for _, e := range entries {
		e.selected = selected
	}
}

func selectedPaths(entries []*fileEntry) []string {
	paths := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.selected {
			paths = append(paths, e.path)
		}
	}

	return paths
}

func printFileSelection(entries []*fileEntry) {
	fmt.Println()
	fmt.Println("\033[1mGitUp™ — Select files\033[0m")
	fmt.Println()

	for i, e := range entries {
		box := "☐"
		if e.selected {
			box = "☑"
		}

		fmt.Printf("  %2d. %s %c %s\n", i+1, box, e.status, e.path)
	}

	fmt.Println()
	fmt.Printf("Selected: %d file(s)\n", len(selectedPaths(entries)))
}
