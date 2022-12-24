package shorten

import (
	"github.com/gudimz/urlShortener/internal/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	getShorten    = "/:id"
	createShorten = "/create"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(getShorten, h.GetShortenById)
	router.POST(createShorten, h.CreateShorten)
	router.DELETE(getShorten, h.DeleteShortenById)
}

func (h *handler) GetShortenById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//TODO: Implementation
	w.Write([]byte("this is get shorten by id"))
}

func (h *handler) CreateShorten(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//TODO: Implementation
	w.Write([]byte("this is add shorten"))
}
func (h *handler) DeleteShortenById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//TODO: Implementation
	w.Write([]byte("this is delete shorten by id"))
}
