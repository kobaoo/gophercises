package main

import (
	"cyoa/internal/handler"
	"cyoa/internal/models"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 8080, "Port to run server")
	filename := flag.String("file", "gopher.json", "Path to json file with stories")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("failed to open file: %s", err))
	}
	var stories models.Story
	if err := json.NewDecoder(file).Decode(&stories); err != nil {
		exit(fmt.Sprintf("failed to decode json: %s", err))
	}
	tpl := template.Must(template.New("").Parse(tmpl))
	h := handler.NewHandler(stories, handler.WithPathFn(pathFn), handler.WithTemplate(tpl))

	fmt.Println("Starting the server on port:", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "story" || path == "/story" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var tmpl = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}

    <ul>
      {{range .Options}}
      <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
      {{end}}
    </ul>
  </body>
</html>`
