package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (h *handler) Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("HOME")
	fmt.Println(h.cfg.Token.Access)
	fmt.Println(h.cfg.Token.Refresh)
}
