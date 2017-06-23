package main

import (
	"github.com/raoptimus/gserv/service"
	"net"
	"net/http"
	"os"
	"strconv"
)

type (
	Server struct {
		mux *http.ServeMux
	}
)

func NewServer(addr string) *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	go s.listening(addr)
	return s
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	defer service.DontPanic() //todo 500 on recovery
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "close")

	//it's (haProxy) checks that the server is healthy
	if r.Method == "HEAD" {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	contentId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
	}
	if err != nil || contentId <= 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	var contentType ContentType

	switch r.FormValue("type") {
	case "article":
		{
			contentType = ArticleContentType
		}
	case "photoalbum":
		{
			contentType = PhotoAlbumContentType
		}
	case "video":
		{
			contentType = VideoContentType
		}
	default:
		http.Error(w, "", http.StatusNotFound)
		return
	}

	raw := NewStatRaw(contentId, contentType)
	raw.ViewCount++
	counter.Push(raw)

}

func (s *Server) listening(addr string) {
	os.Remove(addr)

	l, err := net.Listen("unix", addr)
	if err != nil {
		log.Println(err)
		return
	}

	s.mux.HandleFunc("/", s.serve)
	log.Fatal(http.Serve(l, s.mux))
}
