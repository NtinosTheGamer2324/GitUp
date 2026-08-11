package helper

import "fmt"

type Response int

const (
	Y Response = iota
	N
)

func ConfirmationDiag(question string, consequences string, action string) Response {
	var input string

	fmt.Println()
	fmt.Printf("%s\n", Bold(Cyan("? "+question)))
	fmt.Printf("%s\n", Yellow(consequences))
	fmt.Println()

retry:
	fmt.Printf("%s %s: ", Green(Bold("→ "+action)), Dim("[y/N]"))
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
