package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	// "test_go1/audioSource"
	"test_go1/data"
)

var (
	getAudioSourceRe = regexp.MustCompile(`^\/audioSource\/.`)
	getPlaylist      = regexp.MustCompile(`^\/audio\/(\d+)$`)
	getAllPlaylists  = regexp.MustCompile(`^\/audio\/$`)
)

type audioStore struct {
	playlists []data.Playlist
	*sync.RWMutex
}

type sourceStore struct {
	*sync.RWMutex
}

type audioSourceHandler struct {
	store *sourceStore
}

type audioHandler struct {
	store *audioStore
}

func (h *audioSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeHTTP audioSourceHandler")
	switch {
	case r.Method == http.MethodGet && getAudioSourceRe.MatchString(r.URL.Path):
		h.GetSource(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}
func (h *audioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeHTTP audioHandler")
	switch {
	case r.Method == http.MethodGet && getAllPlaylists.MatchString(r.URL.Path):
		h.GetAll(w, r)
		return
	case r.Method == http.MethodGet && getPlaylist.MatchString(r.URL.Path):
		h.GetPlaylistById(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *audioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetAll")
	jsonBytes, err := json.Marshal(h.store.playlists)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func playlistById(givenId string, playlists []data.Playlist) *data.Playlist {
	for _, p := range playlists {
		if p.ID == givenId {
			return &p
		}
	}
	return nil
}
func (h *audioHandler) GetPlaylistById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetPlaylistById")
	substring := getPlaylist.FindStringSubmatch(r.URL.Path)
	playlist := playlistById(substring[1], h.store.playlists)
	jsonBytes, err := json.Marshal(playlist)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *audioSourceHandler) GetSource(w http.ResponseWriter, r *http.Request) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GetSource")
	var path = dir + "/" + r.URL.Path[1:]
	w.Header().Set("Content-type", "audio/mpeg")

	w.Header().Set("accept-ranges", "bytes")
	// w.Header().Set("Content-Range", "bytes 0-42119807/42119808")
	h.store.RLock()
	fmt.Println(path)
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	fi, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	size := fi.Size()
	w.Header().Set("Content-Length", strconv.FormatInt(int64(size), 10))
	h.store.RUnlock()
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func main() {
	fmt.Println("Started")
	sourceH := &audioSourceHandler{
		store: &sourceStore{
			RWMutex: &sync.RWMutex{},
		},
	}
	audioH := &audioHandler{
		store: &audioStore{
			playlists: audioSource.Playlists,
			RWMutex:   &sync.RWMutex{},
		}}
	mux := http.NewServeMux()
	mux.Handle("/audioSource/", sourceH)
	mux.Handle("/audio/", audioH)
	http.ListenAndServe("localhost:8080", mux)
}
