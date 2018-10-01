package utils

import "fmt"

// "Would you like to edit the PKGBUILD? [y/N]"
func AskConfirmation(message string, defaultAnswer bool) (bool, error) {
	fmt.Println(message)
	var response string
	yes := []string{"y", "Y", "yes", "Yes", "YES"}
	no := []string{"n", "N", "no", "No", "NO"}
	_, err := fmt.Scanln(&response)
	if err != nil {
		return defaultAnswer, err
	}
	// FROM: https://stackoverflow.com/questions/23025694/is-there-no-xor-operator-for-booleans-in-golang
	// (X xor Y) -> X != Y
	if IsInSlice(response, yes) {
		return true, nil
	} else if IsInSlice(response, no) {
		return false, nil
	} else if response == "" {
		return defaultAnswer, nil
	}
	fmt.Println("Answer not understood, repeating.")
	return AskConfirmation(message, defaultAnswer)
}

func IsInSlice(word string, slice []string) bool {
	inSlice := false
	for _, sliceMemeber := range slice {
		if word == sliceMemeber {
			inSlice = true
		}
	}
	return inSlice
}
