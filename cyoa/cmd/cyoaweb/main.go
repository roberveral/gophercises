package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/roberveral/gophercises/cyoa/story"
	"github.com/roberveral/gophercises/cyoa/web"
)

func main() {
	port := flag.Int("port", 8080, "Port to bind the server to")
	storyPath := flag.String("story", "gopher.json", "Path to the JSON definition of the Story")
	templatePath := flag.String("template", "chapter.html", "Path to the HTML template used to render each chapter of the story")

	flag.Parse()

	file, err := os.Open(*storyPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	myStory, err := story.FromJSON(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	tpl, err := template.ParseFiles(*templatePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Starting server in port %d", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), web.New(myStory, web.WithTemplate(tpl))))
}
