package utils

import "fmt"

func PrintErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
