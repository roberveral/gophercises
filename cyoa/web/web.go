package web

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/roberveral/gophercises/cyoa/story"
)

// Default template used to render Chapters in the CYOA website
const defaultChapterTemplate string = `
<h1>{{.Title}}</h1>

{{range .Paragraphs}}
<p>{{.}}</p>
{{end}}

{{range .Options}}
<a style="display: block" href="/chapters/{{.Chapter}}">{{.Text}}</a>
{{end}}
`

// handler is an http.Handler implementation which renders and returns
// the proper chapter according to the path.
// 		/chapters/:name renders chapter 'name' of the story.
//		/ renders the intro chapter.
type handler struct {
	myStory         *story.Story
	chapterTemplate *template.Template
}

// HandlerOption is an alias for the functional options when creating a
// handler.
// (https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
type HandlerOption func(h *handler)

// WithTemplate is an option when creating a handler which makes it use
// the given Template instead of the default one.
func WithTemplate(tpl *template.Template) HandlerOption {
	return func(h *handler) {
		h.chapterTemplate = tpl
	}
}

// New creates a new http.Handler which renders and returns
// the proper chapter according to the path.
//
// 		/chapters/:name renders chapter 'name' of the story.
//		/ renders the intro chapter.
//
// The handler exposes the given story, and the options can be used to
// customize the created handler.
func New(myStory *story.Story, options ...HandlerOption) http.Handler {
	defaultTemplate := template.Must(template.New("").Parse(defaultChapterTemplate))
	h := &handler{myStory, defaultTemplate}

	for _, option := range options {
		option(h)
	}

	return h
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	pathPattern := regexp.MustCompile("^/chapters/(.*)$")

	var chapter *story.Chapter
	var ok bool

	if path == "" || path == "/" {
		chapter, ok = h.myStory.FindIntro()
	} else if matches := pathPattern.FindStringSubmatch(path); matches != nil {
		chapter, ok = h.myStory.FindChapter(matches[1])
	}

	if !ok {
		http.NotFound(rw, r)
		return
	}

	err := h.chapterTemplate.Execute(rw, chapter)
	if err != nil {
		http.Error(rw, "Something went wrong...", http.StatusInternalServerError)
	}
}
