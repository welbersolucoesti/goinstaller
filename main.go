package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gocolly/colly"
)

func downloadFile(url string, path string) (filepath string, err error) {

	resp, err := http.Get(url)

	if err != nil {
		return filepath, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return filepath, fmt.Errorf(fmt.Sprintf("O site devolveu %d", resp.StatusCode))
	}

	splitedURL := strings.Split(url, "/")

	filepath = splitedURL[(len(splitedURL) - 1)]

	file, err := os.Create(filepath)

	if err != nil {
		return filepath, err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return filepath, err
	}

	return filepath, nil
}

func getURLFile() (linuxFileURL string) {

	pageLink := "https://go.dev/dl/"

	c := colly.NewCollector(colly.AllowedDomains("go.dev"))

	c.OnHTML(".downloadBox", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "linux") {
			linuxFileURL = fmt.Sprintf("https://go.dev%s", link)
		}
	})

	c.Visit(pageLink)

	return
}

func descompressAndInstall(filepath string) (err error) {

	var cmd *exec.Cmd

	cmd = exec.Command("rm", "-rf", "/usr/local/go")
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("tar", "-C", "/usr/local", "-xvf", filepath)
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("rm", "-rf", filepath)
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func main() {

	var err error

	filepath, err := downloadFile(getURLFile(), "./")

	if err != nil {
		log.Panic(err.Error())
	}

	err = descompressAndInstall(filepath)

	if err != nil {
		log.Panic(err.Error())
	}
}
