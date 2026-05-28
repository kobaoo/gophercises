package handler

import (
	"cyoa/internal/models"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultTemplateHandler))
}

var tpl *template.Template

var defaultTemplateHandler = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Choose Your Own Adventure</title>
    <style>
      body {
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        max-width: 800px;
        margin: 0 auto;
        padding: 20px;
        background: #faf9f6;
        color: #1a1a1a;
        line-height: 1.6;
      }

      h1 {
        font-size: 2.5em;
        border-bottom: 3px solid #4a5568;
        padding-bottom: 10px;
        margin-bottom: 30px;
        color: #2d3748;
      }

      p {
        font-size: 1.1em;
        margin-bottom: 20px;
      }

      ul {
        list-style: none;
        padding: 0;
        margin-top: 40px;
      }

      li {
        margin: 15px 0;
      }

      a {
        display: inline-block;
        background: #4a5568;
        color: white;
        text-decoration: none;
        padding: 12px 24px;
        border-radius: 6px;
        transition: background 0.2s;
      }

      a:hover {
        background: #2d3748;
      }

      .no-options {
        text-align: center;
        padding: 40px;
        background: #edf2f7;
        border-radius: 8px;
        margin-top: 40px;
      }

      .restart {
        margin-top: 20px;
        background: #48bb78;
      }

      .restart:hover {
        background: #38a169;
      }

      @media (max-width: 600px) {
        body {
          padding: 15px;
        }
        
        h1 {
          font-size: 1.8em;
        }
        
        a {
          display: block;
          text-align: center;
        }
      }
    </style>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}

    {{if .Options}}
      <ul>
        {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
      </ul>
    {{else}}
      <div class="no-options">
        <p>🏁 Конец приключения! 🏁</p>
        <a href="/intro" class="restart">Начать сначала</a>
      </div>
    {{end}}
  </body>
</html>
`

type HandlerOption func(*storyHandler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *storyHandler) {
		h.t = t
	}
}

func WithPathFn(fn func(*http.Request) string) HandlerOption {
	return func(h *storyHandler) {
		h.pathFn = fn
	}
}

func NewHandler(s models.Story, opts ...HandlerOption) http.Handler {
	h := storyHandler{s: s, t: tpl, pathFn: defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}

	return h
}

type storyHandler struct {
	s      models.Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

func (h storyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}
