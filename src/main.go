package main

import (
	"fmt"
	"gitup/src/auth"
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
		args := os.Args[2 : len(os.Args)-1]

		dynamicFileSelect := false

		if len(args) > 0 && args[0] == "--dynamic_file" {
			dynamicFileSelect = true
			args = args[1:]
		}

		if len(args) > 0 && args[0] == "--path" {
			if len(args) < 2 {
				helper.LogFail("Please provide a project path")
				return
			}

			folder = args[1]
			args = args[2:]
		}

		files := args

		if dynamicFileSelect {
			selected, err := git.SelectFilesInteractively(folder)
			if err != nil {
				helper.LogError("Failed to open file selection: %s", err)
				return
			}

			if selected == nil {
				return
			}

			files = selected
		}

		helper.Log("Creating commit: %s", message)

		err := git.Commit(folder, files, message)
		if err != nil {
			helper.LogError("Failed to create commit: %s", err)
		}

	case "login":
		if err := auth.Login(); err != nil {
			helper.LogError("Login failed: %s", err)
		}

	case "logout":
		if err := auth.Logout(); err != nil {
			helper.LogError("Logout failed: %s", err)
		}

	case "whoami":
		if err := auth.WhoAmI(); err != nil {
			helper.LogError("Failed to check identity: %s", err)
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
	fmt.Println("  init                        Create a GitUp repository")
	fmt.Println("  setup                       Configure a repository")
	fmt.Println("  setup gitignore             Create default gitignore")
	fmt.Println("  status                      Show repository status")
	fmt.Println("  commit <files> \"msg\"        Commit specific files")
	fmt.Println("  commit \"msg\"                Commit all changes")
	fmt.Println("  commit --dynamic_file \"msg\" Pick files interactively, then commit")
	fmt.Println("  login                       Log in to GitHub")
	fmt.Println("  logout                      Log out of GitHub")
	fmt.Println("  whoami                      Show the logged-in GitHub identity")
}
