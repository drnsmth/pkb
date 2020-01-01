package main

import (
	"bufio"
	"fmt"
	"os"
	"text/template"
	"time"
)

type Note struct {
	createdDate time.Time
	Content     string
}

func (note *Note) FormatCreatedDate() string {
	return note.createdDate.Format(time.RFC3339)
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(1)
	}

	home := os.Getenv("PKB_HOME")
	if home == "" {
		fmt.Println("Must specify a path to store notes using the environment variable PKB_HOME")
		os.Exit(1)
	}

	if _, err := os.Stat(home); os.IsNotExist(err) {
		fmt.Println(home, "does not exist")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "note":
		createNote(os.Args[2], home)
	case "link":
		createLink(os.Args[2], home)
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: pkb <command> <options>")
}

func createNote(content string, dir string) {
	note := Note{createdDate: time.Now(), Content: content}
	writeNote(dir, note)
}

func createLink(content string, dir string) {
	markdown := fmt.Sprintf("[Link](%s)", content)
	note := Note{createdDate: time.Now(), Content: markdown}
	writeNote(dir, note)
}

const markdownTemplate = `---
Created: {{.FormatCreatedDate}}
---

* {{.Content}}

`

func writeNote(dir string, note Note) {
	template, err := template.New("template").Parse(markdownTemplate)
	if err != nil {
		panic(err)
	}

	filename := note.createdDate.Format("2006-01-02T1504.md")
	file, err := os.Create(dir + "/" + filename)
	defer file.Close()

	out := bufio.NewWriter(file)

	err = template.Execute(out, &note)
	if err != nil {
		panic(err)
	}

	out.Flush()
}
