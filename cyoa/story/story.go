package story

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// Story is a "Choose your own adventure" story which has a series of chapters
// and an entrypoint (the introductory chapter).
type Story struct {
	// Intro is the name of the introductory chapter.
	// It should be contained in the "Chapters" map.
	Intro string `json:"intro"`
	// Chapters is the collection of chapters of the story mapped by their name.
	Chapters map[string]Chapter `json:"chapters"`
}

// Chapter is a part of a Story, defined by the contents of this part of the
// Story and the options to move forward from this Chapter.
type Chapter struct {
	// Title is the title of the chapter
	Title string `json:"title"`
	// Paragraphs is the slice of paragraphs which forms the chapters' story.
	Paragraphs []string `json:"story"`
	// Options is the slice of possible options to move forward from this chapter.
	Options []Option `json:"options"`
}

// Option is a possible choice to continue the adventure
// from one chapter to another.
type Option struct {
	// Text is the option description.
	Text string `json:"text"`
	// Chapter is the name of the chapter where the option leads to.
	Chapter string `json:"arc"`
}

// FromJSON parses a Story from its JSON representation. It receives a
// reader which will be decoded as a Story. The contents of the reader must
// follow the following structure:
//
// 		{
//			"intro": "start",
//			"chapters": {
//				"start": {
//					"title": "My story",
//					"story": ["My content"],
//					"options": [{ "text": "Cyclic", "chapter": "start" }]
//				}
//			}
//    }
//
func FromJSON(reader io.Reader) (*Story, error) {
	decoder := json.NewDecoder(reader)

	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, errors.Wrap(err, "Invalid/malformed JSON Story file")
	}

	return &story, nil
}

// FindIntro obtains the introductory Chapter of the Story.
// If the introductory chapter is not found (The defined intro chapter is not
// present in the chapters list), false is returned in the second argument.
func (s *Story) FindIntro() (*Chapter, bool) {
	return s.FindChapter(s.Intro)
}

// FindChapter obtains a Chapter from the Story given its name. It returns
// false in the second argument if there isn't a chapter with the given name.
func (s *Story) FindChapter(name string) (*Chapter, bool) {
	if chapter, ok := s.Chapters[name]; ok {
		return &chapter, true
	}
	return nil, false
}
