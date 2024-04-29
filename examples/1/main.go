package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/TonyLeCode/go-auth/auth"
	"github.com/TonyLeCode/go-auth/sessions"
)

func protectedRoute(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	data := struct{ Session_ID string }{Session_ID: url.QueryEscape(cookie.Value)}
	tmpl := template.Must(template.ParseFiles("examples/1/protectedRoute.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, _ := r.Cookie("sessionID")
		if cookie != nil {
			http.Redirect(w, r, "/protectedRoute", http.StatusFound)
		}

		data := struct{ Name string }{Name: "hello"}
		tmpl := template.Must(template.ParseFiles("examples/1/login.html"))

		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")

		cookie, err := r.Cookie("sessionID")
		if err != nil || cookie.Value == "" {
			session := sessions.CreateSession()
			cookie := http.Cookie{
				Name:     "sessionID",
				Value:    url.QueryEscape(session.SessionID),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(time.Hour * 12),
			}
			http.SetCookie(w, &cookie)
		} else {
			// sessionID, _ := url.QueryUnescape(cookie.Value)
			http.Redirect(w, r, "/protectedRoute", http.StatusFound)
		}

		fmt.Println("email: ", email, "password: ", password)

	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := struct{ Success bool }{Success: false}
		tmpl := template.Must(template.ParseFiles("examples/1/register.html"))

		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email") // If you have an email field
		fmt.Println("username: ", username, "password: ", password, "email: ", email)

		var errorMessage string
		if username == "" {
			errorMessage = "Username cannot be empty"
		} else if len(username) < 5 {
			errorMessage = "Username must be at least 5 characters long"
		}

		if email != "" {
			// Use a regular expression to validate email format (example omitted)
		}

		if password == "" {
			errorMessage = "Password cannot be empty"
		}

		if errorMessage != "" {
			// Handle validation error:
			// - Display the error message to the user on the registration form.
			// - Consider using a struct to hold validation errors for better organization.
			fmt.Fprintf(w, "Error: %s", errorMessage)
			return
		}

		hashPassword := auth.HashPassword(password)
		fmt.Println("Hashed password: ", string(hashPassword))

		data := struct{ Success bool }{Success: true}
		tmpl := template.Must(template.ParseFiles("examples/1/register.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	port := "8090"
	address := "localhost:" + port
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/protectedRoute", protectedRoute)

	fmt.Println("Starting server on port: " + port)
	log.Fatal(http.ListenAndServe(address, nil))
}
