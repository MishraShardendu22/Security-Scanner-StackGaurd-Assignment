package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Name string
	Age  int
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
    <html>
      <head><title>User Info</title></head>
      <body>
        <h1>Hello {{ .Name }}</h1>
        <p>You are {{ .Age }} years old.</p>
      </body>
    </html>
    `
	t := template.Must(template.New("webpage").Parse(tmpl))

	user := User{Name: "Shardendu", Age: 25}
	t.Execute(w, user)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("This is a test commment")
	http.ListenAndServe(":8080", nil)
}
