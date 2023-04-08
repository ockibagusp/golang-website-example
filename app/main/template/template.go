package template

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
)

// Templates is a custom html/template renderer for Echo framework
type Templates struct {
	Templates map[string]*template.Template
}

// append templates package controller anywhere
var AppendTemplates = make(map[string]*template.Template)

// New Templates
func NewTemplates() *Templates {
	return &Templates{
		Templates: AppendTemplates,
	}
}

// Render implement e.Renderer interface
func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return errors.New("Template not found -> " + name)
	}

	// Add global methods if data is a map
	if viewContext, isMap := data.(echo.Map); isMap {
		viewContext["reverse"] = c.Echo().Reverse

		/*
			@param /login and /users/view/:id/password
			- Is HTML Only
			is_html_only (bool): {true, false}

			Login name: 				"login.html -> ok
			UpdateUserByPassword name: 	"user-view-password.html" -> ok
										"users/user-view-password.html" -> no
			why?
		*/
		if viewContext["is_html_only"] == true {
			return tmpl.ExecuteTemplate(w, name, data)
		}
	}

	return tmpl.ExecuteTemplate(w, "base.html", data)
}

// Parse Files Base: not implement
//
// ParseFilesBase("views/login.html") -> base.html. Yes, but nothing.
// ParseFilesBase("views/users/user-view-password.html") -> base.html. Yes, but nothing.
// ---
// parseFileHTMLOnly("views/login.html") -> Yes, this is good.
// parseFileHTMLOnly("views/users/user-view-password.html") -> Yes, this is good.
func ParseFilesBase(s string, t ...string) *template.Template {
	dir := rootedPathName()

	templateBase := template.New("base").Funcs(FuncMapMore())

	// t parseFilesBase, example "views/user-form.html"
	if len(t) == 1 {
		return template.Must(
			/* template.New("") and template.New("base"), equal.

			t := make(map[string]*template.Template)
			t["home.html"] = parseFilesBase("views/home.html")
			...
			*/
			templateBase.ParseFiles(
				fmt.Sprintf("%s/%s", dir, s),
				fmt.Sprintf("%s/%s", dir, t[0]),
				fmt.Sprintf("%s/%s", dir, "views/base.html"),
			),
		)
	} else if len(t) >= 2 {
		panic("t [1] parseFilesBase, example \"views/users/user-form.html\"")
	}
	// "views/base.html"?
	return template.Must(
		templateBase.ParseFiles(
			fmt.Sprintf("%s/%s", dir, s),
			fmt.Sprintf("%s/%s", dir, "views/base.html"),
		),
	)
}

// Parse File HTML Only
func ParseFileHTMLOnly(name string) *template.Template {
	dir := rootedPathName()

	return template.Must(
		/* template.New("") and template.New("HTML_only"), equal.

		t := make(map[string]*template.Template)
		...
		t["login.html"] = parseFileHTMLOnly("views/login.html")
		...
		*/
		template.New("HTML_only").Funcs(FuncMapMore()).
			ParseFiles(
				fmt.Sprintf("%s/%s", dir, name),
			),
	)
}

const projectDirName = "golang-website-example"

// Rooted Path Name
func rootedPathName() string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	return string(rootPath)
}
