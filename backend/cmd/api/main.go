package main

import (
	"fmt"
	"os"

	"faculty/internal/app"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := app.Migrate(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	a, err := app.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer a.Close()

	if err := a.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
