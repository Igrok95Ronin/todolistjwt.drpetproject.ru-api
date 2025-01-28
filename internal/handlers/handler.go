package handlers

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Router(router *httprouter.Router)
}
