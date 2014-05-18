// fypMain.go
package main

import (
	"fypTest/controllers"
	"fypTest/models"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//"github.com/martini-contrib/sessionauth"
	//"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
)

// DB returns a martini.Handler
func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("MeLTS"))
		defer s.Close()
		c.Next()
	}
}

func main() {
	// Set up a classic martini handler
	m := martini.Classic()
	// Create the render and point it to our views folder (holds templates)
	m.Use(render.Renderer(render.Options{Directory: "views"}))
	// Tell martini to use our mongo DB
	m.Use(DB())

	// Set Up Session
	//store := sessions.NewCookieStore([]byte("secret123"))
	//m.Use(sessions.Sessions("my_session", store))
	//m.Use(sessionauth.Sessionuser(GenerateAnonymousUser))
	//sessionauth.RedirectUrl = "/new-login"
	//sessionauth.RedirectParam = "new-next"

	// Function calls

	m.Get("/students", controllers.StudentIndex)
	m.Get("/login", controllers.Login)
	//m.Post("/students/add", controllers.StudentAdd)
	m.Get("/students/add", controllers.StudentAddPage)
	m.Get("/students/remove", controllers.StudentRemovePage)

	//m.Use(martini.Static("assets"))

	// Change the port to listen on Port 8001 as specified by Jon
	//log.Fatal(http.ListenAndServe(":8001", m))

	m.Run()
}
