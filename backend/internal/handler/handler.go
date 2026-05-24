package handler

import (
	"database/sql"
	"net/http"
)

type Renderer interface {
	RenderTemplate(w http.ResponseWriter, name string, data interface{})
	RenderLogin(w http.ResponseWriter, name string, data interface{})
}

type Handler struct {
	db             *sql.DB
	renderer       Renderer
	FacilitiesJSON string
}

func New(db *sql.DB, r Renderer, facilitiesJSON string) *Handler {
	return &Handler{db: db, renderer: r, FacilitiesJSON: facilitiesJSON}
}
