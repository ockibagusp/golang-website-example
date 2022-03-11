package template

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/labstack/echo/v4"
)

// Templates is a custom html/template renderer for Echo framework
type Templates struct {
	Templates map[string]*template.Template
}

// New Templates
func NewTemplates() *Templates {
	t := make(map[string]*template.Template)
	t["home.html"] = parseFilesBase("views/home.html")
	t["about.html"] = parseFilesBase("views/about.html")
	t["login.html"] = parseFileHTMLOnly("views/login.html")
	t["users/user-all.html"] = parseFilesBase("views/users/user-all.html")
	t["users/user-add.html"] = parseFilesBase("views/users/user-add.html", "views/users/user-form.html")
	t["users/user-read.html"] = parseFilesBase("views/users/user-read.html", "views/users/user-form.html")
	t["users/user-view.html"] = parseFilesBase("views/users/user-view.html", "views/users/user-form.html")
	/*
		t["users/user-view-password.html"] -> no
			{..,"status":500,"error":"html/template: \"users/user-view-password.html\" is undefined",..}
			why?
		t["user-view-password.html"] -> yes
	*/
	t["user-view-password.html"] = parseFileHTMLOnly("views/users/user-view-password.html")

	return &Templates{
		Templates: t,
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

			TODO:
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
// parseFilesBase("views/login.html") -> base.html. Yes, but nothing.
// parseFilesBase("views/users/user-view-password.html") -> base.html. Yes, but nothing.
// ---
// parseFileHTMLOnly("views/login.html") -> Yes, this is good.
// parseFileHTMLOnly("views/users/user-view-password.html") -> Yes, this is good.
func parseFilesBase(s string, t ...string) *template.Template {
	dir := rootedPathName()

	templateBase := template.New("base").Funcs(FuncMapMore())

	// t parseFilesBase, example "views/user-form.html"
	if len(t) == 1 {
		return template.Must(
			/* template.New("") and template.New("base"), same.

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
func parseFileHTMLOnly(name string) *template.Template {
	dir := rootedPathName()

	return template.Must(
		/* template.New("") and template.New("HTML_only"), same.

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

// Rooted Path Name
func rootedPathName() (dir string) {
	dir, err := os.Getwd()
	if err != nil {
		// TODO: Docker, Kubernetes ex.?
		log.Fatal(err)
	}

	// TODO: Linux and MacOS: ok. Windows: ...?
	regex := regexp.MustCompile("/tests$")
	match := regex.Match([]byte(dir))

	if match {
		_dir, err := os.Open(path.Join(dir, "../"))
		if err != nil {
			log.Fatal(err)
		}
		dir = _dir.Name()
	}
	return
}

func executable() {

}
