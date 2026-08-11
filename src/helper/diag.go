package helper

import "fmt"

type Response int

const (
	Y Response = iota
	N
)

func ConfirmationDiag(question string, consequences string, action string) Response {

	var input string

	fmt.Printf(question + "\n")
	fmt.Printf(consequences + "\n")
	fmt.Println()

retry:
	fmt.Printf("%s [y,N]: ", action)
	fmt.Scanln(&input)

	if input == "Y" || input == "y" {
		return Y
	} else if input == "N" || input == "n" {
		return N
	} else if input == "" {
		return N
	} else {
		fmt.Println("Please Provide one of the two options. N or Y")
		goto retry
	}
}
