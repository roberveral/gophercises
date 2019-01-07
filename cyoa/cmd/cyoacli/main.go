package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/roberveral/gophercises/cyoa/story"
)

const chapterTemplate string = `
- {{.Title}}

{{range .Paragraphs}}
{{.}}

{{end}}

---------------------------------------

{{range $i, $option := .Options}}
  - [{{$i}}]: {{.Text}}
{{end}}
`

// Quick solution, code can be improved a lot
func main() {
	storyPath := flag.String("story", "gopher.json", "Path to the JSON definition of the Story")

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

	tpl := template.Must(template.New("").Parse(chapterTemplate))

	chapter, ok := myStory.FindIntro()
	if !ok {
		log.Fatal("Stories intro chapter is not defined")
		return
	}

	for {
		tpl.Execute(os.Stdout, chapter)
		if len(chapter.Options) == 0 {
			return
		}
		fmt.Print("Choose your option: ")
		var option int
		fmt.Scanf("%d\n", &option)
		chapter, _ = myStory.FindChapter(chapter.Options[option].Chapter)
	}
}
