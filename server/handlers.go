package server

import (
	"net/http"
)

func Home_Handler(w http.ResponseWriter, r *http.Request){

	render_Template(w,"home.html", nil)
}