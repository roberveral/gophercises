package urlshort

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Store is a interface which contains methods required for the interaction
// with a persistence which stores the mapping between short paths and urls.
// This allows to make quick implementations the urlshort.handler that are
// backed in different databases.
type Store interface {
	Get(path string) (string, bool)
	Put(path string, url string) error
}

// Main type of the urlshort handler service, which contains the store where
// the mappings are persisted. It's not exported because its only purpose is
// implement the http.Handler interface.
type handler struct {
	store Store
}

// New creates a new urlshort http.handler which uses the given store as persistence
// for the path mappings. The handler contains a REST API to create and obtain path
// mappings and the redirect logic which redirects the user to the appropiate URL given
// the short path.
//
// - POST /api/urls {"url": "..."} :- Creates a new short path for the given URL.
// - GET /api/urls/{shortPath} :- Obtains the URL associated with the given shortPath.
// - Any request to /{shortPath} :- Redirects the request to the URL associated with the
//                                  shortPath when exists.
func New(store Store) http.Handler {
	urlshortHandler := &handler{store}
	router := mux.NewRouter()

	router.HandleFunc("/api/urls/{path}", urlshortHandler.GetURL).Methods("GET")
	router.HandleFunc("/api/urls", urlshortHandler.CreateURL).Methods("POST")
	router.PathPrefix("/").Handler(urlshortHandler)

	return router
}

// http.Handler method implements the redirect logic.
func (h *handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path[1:]

	log.Infof("Redirect request for path '%s'", path)

	if dest, ok := h.store.Get(path); ok {
		http.Redirect(response, request, dest, http.StatusFound)
		return
	}

	http.NotFound(response, request)
}

// CreateURL is the function executed in each request to create a new short path for a URL.
func (h *handler) CreateURL(response http.ResponseWriter, request *http.Request) {
	var body createURLRequest

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(response, "Invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	path := getUniquePath(body.URL)

	log.Infof("Request to shorten URL '%s'. Path assigned: '%s'", body.URL, path)

	if err := h.store.Put(path, body.URL); err != nil {
		log.Errorf("Generated path '%s' cannot be stored due to: %+v", path, err)
		http.Error(response, "There was a collision with the generated path, please try again", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Location", "/"+path)
	response.WriteHeader(http.StatusCreated)
}

// GetURL is the function executed in each request to obtain the URL associated to a short path.
func (h *handler) GetURL(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	path := vars["path"]

	log.Infof("Request to obtain URL associated with path '%s'", path)

	if dest, ok := h.store.Get(path); ok {
		encoder := json.NewEncoder(response)
		encoder.Encode(urlMapping{path, dest})
		return
	}

	http.NotFound(response, request)
}

// Obtains the short path to associate to the given url
func getUniquePath(url string) string {
	// Include current timestamp in the hashed string to ensure that we can have multiple paths for the same url
	hashedURL := sha256.Sum256([]byte(url + time.Now().Format(time.RFC850)))
	return base64.URLEncoding.EncodeToString(hashedURL[:])[:9]
}

// Helper type to parse a CreateURL request.
type createURLRequest struct {
	URL string `json:"url,omitempty"`
}

// Helper type to return the GetURL response.
type urlMapping struct {
	Path string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}
