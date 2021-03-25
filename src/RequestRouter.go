package main

import (
	"fmt"
	"net/http"
)

type RequestRouter struct {
	FileServerPrefix string

	FileServer http.Handler
	MockerHandler func (http.ResponseWriter, *http.Request)
}

func (r *RequestRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	firstPathPart := ""

	fmt.Println(request.URL.Path)

	for _, it := range request.URL.Path {

		fmt.Println("loop", string(it))
		fmt.Println("loop res", firstPathPart)

		if it == '/' {

			if len(firstPathPart) == 0 {
				continue
			}

			break
		}
		firstPathPart += string(it)
	}

	fmt.Println(r.FileServerPrefix)

	if firstPathPart == r.FileServerPrefix {
		fmt.Println("Server", firstPathPart)
		r.FileServer.ServeHTTP(writer, request)
	} else {
		fmt.Println("Mocks", firstPathPart)
		r.MockerHandler(writer, request)
	}
}
