package main

import (
	"fmt"
	"gitup/src/auth"
	"gitup/src/git"
	"gitup/src/github"
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

	case "repo":
		if len(os.Args) < 3 {
			helper.LogFail("Please provide a repo option")
			return
		}

		switch os.Args[2] {
		case "create":
			folder := "."
			name := ""
			private := true

			for _, arg := range os.Args[3:] {
				switch arg {
				case "--public":
					private = false
				case "--private":
					private = true
				default:
					if name == "" {
						name = arg
					}
				}
			}

			if name == "" {
				projectName, err := git.ProjectName(folder)
				if err != nil {
					helper.LogError("Could not determine a repository name: %s", err)
					return
				}
				name = projectName
			}

			creds, err := auth.EnsureValidCredentials()
			if err != nil {
				helper.LogError("%s", err)
				return
			}
			if creds == nil {
				helper.LogFail("You're not logged in to GitHub.")
				helper.Log("Run 'gitup login' first.")
				return
			}

			visibility := "private"
			if !private {
				visibility = "public"
			}

			response := helper.ConfirmationDiag(
				"GitUp is about to create a new GitHub repository.",
				fmt.Sprintf("This will create a %s repository named '%s' under %s.", visibility, name, creds.Login),
				"Continue?",
			)

			if response == helper.N {
				helper.Log("Repository creation cancelled.")
				return
			}

			helper.Log("Creating repository '%s' on GitHub...", name)

			repository, err := github.CreateRepository(creds.AccessToken, name, private, "")
			if err != nil {
				helper.LogError("Failed to create repository: %s", err)
				return
			}

			helper.LogOk("Repository created: %s", repository.HTMLURL)

			if !git.IsDirGitProject(folder) {
				helper.Log("This folder isn't a GitUp/Git repository yet.")
				helper.Log("Run 'gitup init' here, then 'gitup repo create' again to link it automatically,")
				helper.Log("or add the remote manually:")
				helper.Log("  git remote add origin %s", repository.CloneURL)
				return
			}

			if git.HasRemote(folder, "origin") {
				existing, _ := git.GetRemoteURL(folder, "origin")

				response := helper.ConfirmationDiag(
					"This repository already has an 'origin' remote.",
					fmt.Sprintf("Current: %s\nNew:     %s", existing, repository.CloneURL),
					"Replace it?",
				)

				if response == helper.N {
					helper.Log("Kept the existing remote. The repository was still created on GitHub.")
					return
				}
			}

			if err := git.SetRemoteURL(folder, "origin", repository.CloneURL); err != nil {
				helper.LogError("Repository was created, but GitUp failed to configure the remote: %s", err)
				return
			}

			helper.LogOk("Remote 'origin' configured.")
			helper.Log("Run 'gitup publish' to push your code.")

		default:
			helper.LogFail("Unknown repo option: %s", os.Args[2])
		}

	case "publish":
		folder := "."
		remoteName := "origin"

		args := os.Args[2:]
		if len(args) >= 1 {
			remoteName = args[0]
		}
		if len(args) >= 2 {
			folder = args[1]
		}

		creds, err := auth.EnsureValidCredentials()
		if err != nil {
			helper.LogError("%s", err)
			return
		}
		if creds == nil {
			helper.LogFail("You're not logged in to GitHub.")
			helper.Log("Run 'gitup login' first.")
			return
		}

		if err := git.Publish(folder, remoteName, creds.Login, creds.AccessToken); err != nil {
			helper.LogError("Failed to publish: %s", err)
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
	fmt.Println("GitUp")
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
	fmt.Println("  repo create [name]          Create a GitHub repository and link it as origin")
	fmt.Println("                              (add --public to make it public; default is private)")
	fmt.Println("  publish [remote] [path]     Push the current branch (defaults: origin, .)")
	fmt.Println("  login                       Log in to GitHub")
	fmt.Println("  logout                      Log out of GitHub")
	fmt.Println("  whoami                      Show the logged-in GitHub identity")
}
