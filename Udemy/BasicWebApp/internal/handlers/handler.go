package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/config"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/models"
	"github.com/vanshaj/Microservice/Udemy/BasicWebApp/internal/render"
)

const basePath = "/home/vanshaj/Projects/Golang/Microservice/Udemy/BasicWebApp"

var repo *Handler

type Handler struct {
	App *config.AppConfig
}

func NewHandler(app *config.AppConfig) *Handler {
	return &Handler{
		App: app,
	}
}

func (r *Handler) Home(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/home.page.tmpl", nil)
}

func (r *Handler) About(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/about.page.tmpl", nil)
}
func (r *Handler) Data(w http.ResponseWriter, req *http.Request) {
	stringMap := map[string]string{
		"test": "hello",
	}
	td := &models.TemplateData{
		StringMap: stringMap,
	}
	render.RenderTemplate(w, basePath+"/templates/data.page.tmpl", td)
}
func (r *Handler) Reservation(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/make-reservation.page.tmpl", nil)
}
func (r *Handler) Generals(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/generals.page.tmpl", nil)
}
func (r *Handler) Majors(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/majors.page.tmpl", nil)
}
func (r *Handler) Availability(w http.ResponseWriter, req *http.Request) {
	token := nosurf.Token(req)
	fmt.Println(token)
	render.RenderTemplate(w, basePath+"/templates/search-availability.page.tmpl", &models.TemplateData{CSRFToken: token})
}
func (r *Handler) PostAvailability(w http.ResponseWriter, req *http.Request) {
	start := req.Form.Get("start")
	end := req.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Starts at %s , ends at %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Handler) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func (r *Handler) Contact(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, basePath+"/templates/contact.page.tmpl", nil)
}
