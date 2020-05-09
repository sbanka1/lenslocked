package controllers

import (
	"lenslocked/views"
)

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

type Static struct {
	Home    *views.View
	Contact *views.View
}
