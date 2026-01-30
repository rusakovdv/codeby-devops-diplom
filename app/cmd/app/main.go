package main

import "myapp/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
