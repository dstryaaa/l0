package controllers

import (
	"html/template"
	"net/http"
)

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func ServeCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/styles.css")
}

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	var order_uid_from_form string
	if r.Method == http.MethodPost {
		r.ParseForm()

		order_uid_from_form = "/order/show/" + r.FormValue("order_uid")
	}

	http.Redirect(w, r, order_uid_from_form, http.StatusSeeOther)
}
