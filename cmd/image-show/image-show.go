package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sclevine/agouti"
)

func main() {
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}

	page, err := driver.NewPage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	currentDirectory, _ := os.Getwd()
	imagePath := fmt.Sprintf("file:///%s/image.png", currentDirectory)

	if err := page.Navigate(imagePath); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	for {
		time.Sleep(time.Second)
		page.Refresh()
	}
}
