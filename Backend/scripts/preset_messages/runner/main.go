package main

import (
	"log"

	scripts "github.com/Petroviiic/ScreenScore/scripts/preset_messages"
)

func main() {
	if err := scripts.Seed(); err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
