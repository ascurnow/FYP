// fypMain.go
package main

import (
	"FYP/fypTest/controllers"
	//"FYP/fypTest/models"
	"github.com/go-martini/martini"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
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
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("server_session", store))
	//m.Use(sessionauth.Sessionuser(GenerateAnonymousUser))
	//sessionauth.RedirectUrl = "/new-login"
	//sessionauth.RedirectParam = "new-next"

	// Function calls

	// Login Functions
	m.Get("/", controllers.RequireLogin)
	m.Get("/login", controllers.Login)
	m.Post("/login", controllers.PostLogin)
	m.Get("/logout", controllers.Logout)
	m.Get("/signup", controllers.Signup)
	m.Post("/signup", controllers.SignupPost)

	// Student Controller Functions
	m.Get("/students", controllers.StudentIndex)
	m.Post("/edit", controllers.StudentEdit)
	m.Post("/students/updateUser", controllers.StudentUpdate)
	m.Get("/students/add", controllers.StudentAddPage)
	m.Post("/students/remove", controllers.StudentRemovePage)
	m.Post("/students/remove/confirm", controllers.StudentRemove)
	m.Post("/students/add/new", controllers.StudentAdd)

	// Staff Controller Functions
	m.Get("/staff_home", controllers.StaffHomePage)
	m.Get("/unitList", controllers.StaffUnitList)
	m.Get("/my_units", controllers.StaffMyUnit)
	m.Post("/unit", controllers.StaffUnit)
	m.Post("/removeStudentFromUnit", controllers.StaffRemoveStudentFromUnit)
	m.Post("/addStudentToUnit", controllers.StaffAddStudentToUnit)
	m.Post("/addStudentsToUnitFinal", controllers.StaffAddStudentToUnitFinal)
	m.Post("/staffCreateUnit", controllers.StaffCreateUnit)
	m.Post("/staffQuizPage", controllers.StaffQuizPage)
	m.Post("/addResultForQuiz", controllers.StaffAddQuizResult)
	m.Post("/staffAddResultToQuizFinal", controllers.StaffAddQuizResultFinal)
	m.Post("/editResultInUnit", controllers.StaffEditResultInUnit)
	m.Post("/staffEditResultInQuizFinal", controllers.StaffEditResultInUnitFinal)
	m.Post("/addQuizToUnit", controllers.StaffAddQuizToUnit)
	m.Post("/staffAddQuizFinal", controllers.StaffAddQuizToUnitFinal)
	m.Post("/staffRemoveQuizPage", controllers.StaffRemoveQuiz)
	m.Post("/addPoints", controllers.StaffAddPoints)
	m.Post("/addPointsFinal", controllers.StaffAddPointsFinal)
	m.Post("/addExp", controllers.StaffAddExp)
	m.Post("/addExpFinal", controllers.StaffAddExpFinal)
	m.Post("/addAchievement", controllers.StaffAddAchievement)
	m.Post("/addAchievementFinal", controllers.StaffAddAchievementFinal)

	/* Profile Page */
	m.Get("/editProfilePage", controllers.UserProfileEdit)
	m.Post("/editProfileFinal", controllers.UserProfileEditFinal)
	m.Get("/profile", controllers.ProfilePersonal) /* On a get go to personal profile page */
	//m.Post("/profile", controllers.Profile) /* On a post load from form the profile to see */
	m.Get("/purchaseTitle", controllers.PurchaseTitle)

	// Homepage - Student
	m.Get("/home", controllers.StudentHomePage)

	m.Use(martini.Static("assets"))

	// Change the port to listen on Port 8001 as specified by Jon
	//log.Fatal(http.ListenAndServe(":8001", m))

	m.Run()
}
