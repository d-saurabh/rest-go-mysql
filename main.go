package main

import "os"

func main() {
	a := App{}
	a.Init(os.Getenv("DB_USER_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	a.Start(":8000")

}