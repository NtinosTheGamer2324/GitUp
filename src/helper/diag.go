package helper

import "fmt"

type Response int

const (
	Y Response = iota
	N
)

func ConfirmationDiag(question string, consequences string, action string) Response {
	var input string

	fmt.Printf("\033[36m%s\033[0m\n", question)
	fmt.Printf("\033[33m%s\033[0m\n", consequences)
	fmt.Println()

retry:
	fmt.Printf("\033[32m%s\033[0m [y,N]: ", action)
	fmt.Scanln(&input)

	if input == "Y" || input == "y" {
		return Y
	} else if input == "N" || input == "n" || input == "" {
		return N
	} else {
		LogFail("Please provide one of the two options: Y or N")
		goto retry
	}
}
