package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func makeHandlerWithState(registrar IRegisterHandlers, pattern string, state Server, methods map[string]struct{}, handler func(Server, http.ResponseWriter, *http.Request)) {
	registrar.registerHandler(pattern, func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		_, match := methods[request.Method]
		if !match {
			http.Error(writer, "Invalid request method.", 405)
			return
		}
		if state.isShutdown() {
			http.Error(writer, "Server shutting down", 503)
			return
		}
		handler(state, writer, request)
		state.logResponse(pattern, int(1000000*time.Since(start).Seconds()))
	})
}

func shutdownHandler(state Server, writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "initiating shutdown")
	state.shutdown()
}

func statsHandler(state Server, writer http.ResponseWriter, request *http.Request) {
	stats := state.getStats("/hash")
	js, err := json.Marshal(stats)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func computeHashHandler(state Server, writer http.ResponseWriter, request *http.Request) {
	password := request.FormValue("password")
	if password == "" {
		http.Error(writer, "Invalid request, missing param 'password'", 400)
		return
	}

	index := state.doHash(password, 5)

	_, err := fmt.Fprint(writer, index)
	if err != nil {
		log.Print(err)
	}
}

func getHashHandler(state Server, writer http.ResponseWriter, request *http.Request) {
	index, err := strconv.Atoi(request.URL.Path[len("/hash/"):])
	if err != nil {
		http.Error(writer, "Invalid hash index.", 400)
		return
	}

	hashValue := state.getHash(index)
	if hashValue == "not found" {
		http.Error(writer, "Invalid hash index.", 400)
		return
	}
	fmt.Fprint(writer, hashValue)
}

func methods(args ...string) map[string]struct{} {
	elementMap := make(map[string]struct{})
	for _, method := range args {
		elementMap[method] = struct{}{}
	}
	return elementMap
}
