package main

import (
	"fmt"
	"gitup/src/git"
	"gitup/src/helper"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "--help", "-h", "help":
		printHelp()

	case "init":
		folder := "."

		if len(os.Args) >= 3 {
			folder = os.Args[2]
		}

		helper.Log("Initializing git repository.")
		git.InitFolderWithGit(folder)

	case "setup":
		if len(os.Args) < 3 {
			helper.LogFail("Please provide a setup option")
			return
		}

		switch os.Args[2] {
		case "gitignore":
			folder := "."

			if len(os.Args) >= 4 {
				folder = os.Args[3]
			}

			if !git.IsDirGitProject(folder) {
				helper.LogFail("Directory is not a Git repository: %s", folder)
				return
			}

			git.WriteGitIgnore(folder)

		default:
			helper.LogFail("Unknown setup option: %s", os.Args[2])
		}
	case "commit":
		if len(os.Args) < 3 {
			helper.LogFail("Please provide a commit message")
			return
		}

		folder := "."
		message := os.Args[len(os.Args)-1]
		files := os.Args[2 : len(os.Args)-1]

		if len(files) > 0 && files[0] == "--path" {
			if len(files) < 2 {
				helper.LogFail("Please provide a project path")
				return
			}

			folder = files[1]
			files = files[2:]
		}

		helper.Log("Creating commit: %s", message)

		err := git.Commit(folder, files, message)
		if err != nil {
			helper.LogError("Failed to create commit: %s", err)
		}

	default:
		helper.LogFail("Unknown command: %s", os.Args[1])
	}
}

func printHelp() {
	fmt.Println("GitUp™")
	fmt.Println()
	fmt.Println("Usage: gitup <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init              Create a GitUp repository")
	fmt.Println("  setup             Configure a repository")
	fmt.Println("  setup gitingore   Create default gitignore")
	fmt.Println("  status            how repository status")
	fmt.Println("  commit            Commit selected files")
}
