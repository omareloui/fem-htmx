package main

import (
	"html/template"
	"io"
	"slices"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var currId int64 = 0

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("web/views/*.html")),
	}
}

type ErrorRes struct {
	Message string
}

type Contact struct {
	ID    int64
	Name  string
	Email string
}

type Contacts []Contact

func (c *Contacts) hasEmail(email string) bool {
	for _, contact := range *c {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newContact(name, email string) *Contact {
	currId++
	return &Contact{
		ID: currId, Name: name, Email: email,
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() *FormData {
	return &FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type State struct {
	Contacts Contacts
	Form     FormData
}

func newState() *State {
	return &State{
		Contacts: Contacts{
			*newContact("jd", "jd@email.com"),
			*newContact("clara", "cl@email.com"),
		},
	}
}

func main() {
	// TODO: update with native golang stuff (remove echo)
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	e.Static("/images", "web/images")
	e.Static("/css", "web/css")

	state := *newState()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", state)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if state.Contacts.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "this email is already in use"
			return c.Render(422, "form", formData)
		}

		contact := *newContact(name, email)
		state.Contacts = append(state.Contacts, contact)
		err := c.Render(200, "form", newFormData())
		if err != nil {
			return err
		}
		return c.Render(200, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		time.Sleep(time.Second * 3)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(422, "invalid id")
		}

		idx := slices.IndexFunc(state.Contacts, func(c Contact) bool {
			return c.ID == int64(id)
		})

		if idx == -1 {
			return c.String(404, "contact not found")
		}

		state.Contacts = append(state.Contacts[:idx], state.Contacts[idx+1:]...)

		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
