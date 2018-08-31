package main

import (
	"log"
	"net/http"
)

func main() {
	serverState := makeState()
	registrar := &handlerRegistrar{}
	makeHandlerWithState(registrar, "/hash", serverState, methods(http.MethodPost), computeHashHandler)
	makeHandlerWithState(registrar, "/hash/", serverState, methods(http.MethodGet), getHashHandler)
	makeHandlerWithState(registrar, "/stats", serverState, methods(http.MethodGet), statsHandler)
	makeHandlerWithState(registrar, "/shutdown", serverState, methods(http.MethodGet), shutdownHandler)
	log.Println("Listening on server 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type handlerRegistrar struct {
}

func (*handlerRegistrar) registerHandler(pattern string, handler func(writer http.ResponseWriter, request *http.Request)) {
	http.HandleFunc(pattern, handler)
}

type IRegisterHandlers interface {
	registerHandler(string, func(http.ResponseWriter, *http.Request))
}
