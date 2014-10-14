// StudentController
package controllers

import (
	"FYP/fypTest/models"
	"fmt"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	//"os"
)

// Function to return all Users in database
func GetAllStudents(db *mgo.Database) []models.Student {
	var studentList []models.Student
	db.C("studentList").Find(nil).All(&studentList)
	//fmt.Println(studentList)
	return studentList
}

// Function to render the index page
func StudentIndex(r render.Render, db *mgo.Database) {
	r.HTML(200, "Student/index", GetAllStudents(db))
}

// Get the student name from drop list in /students, find this students data
// Display it for editing
func StudentEdit(r render.Render, db *mgo.Database, req *http.Request) {
	name := req.FormValue("selectEditStudent")
	//println(name)
	var student models.Student
	db.C("studentList").Find(bson.M{"username": name}).One(&student)
	r.HTML(200, "Student/edit", student)
}

// Function to update a students details (currently only email, staff status, and soon units) based upon the
// information given in the form
func StudentUpdate(r render.Render, db *mgo.Database, req *http.Request) {
	var student models.Student
	name := req.FormValue("name")
	email := req.FormValue("email")
	staff := req.FormValue("staff")
	unit := req.FormValue("unit")

	//println(name)
	//println(email)
	println(staff)
	println(unit)

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"email": email, "staff": staff}},
		ReturnNew: true,
	}
	db.C("studentList").Find(bson.M{"username": name}).Apply(change, &student)
	//println(student.Username)
	//println(student.Email)
	r.HTML(200, "Student/updateUser", student)
}

// Function to add a student to the database
func StudentAdd(r render.Render, db *mgo.Database, req *http.Request) {
	name := req.FormValue("newStudentName")
	email := req.FormValue("newStudentEmail")
	//println(name)
	//println(email)
	student := models.Student{Username: name, Email: email}
	db.C("studentList").Insert(student)
	r.HTML(200, "Student/add_new", student)
}

// Function for /student/add
func StudentAddPage(r render.Render, db *mgo.Database) {
	r.HTML(200, "Student/add", nil)
}

// Function to remove the student: Finds the student based upon the name
// (should be unique - CHANGE if name is no longer unique)
// Renders a page asking for confirm the deletion
func StudentRemovePage(r render.Render, db *mgo.Database, req *http.Request) {
	var studentToRemove models.Student
	name := req.FormValue("selectRemoveStudent")
	//println(name)
	//println(bson.M{"Username": "blah2"})
	db.C("studentList").Find(bson.M{"username": name}).One(&studentToRemove)
	//println(studentToRemove.Username)
	r.HTML(200, "Student/remove", studentToRemove)
}

// Function to remove the student
func StudentRemove(r render.Render, db *mgo.Database, req *http.Request) {
	option := req.FormValue("confirmDelete")
	name := req.FormValue("studentName")   // DODGY AS FK
	email := req.FormValue("studentEmail") // DODGY AS FK
	//println(option)
	//println(name)
	//println(email)
	if option == "TRUE" {
		db.C("studentList").Remove(bson.M{"username": name, "email": email})
		r.HTML(200, "Student/removal_confirmed", option)
	} else {
		r.HTML(200, "Student/removal_confirmed", option)
	}
}

/*
Function to render /editProfilePage
This function uses the session to determine the current user and allows the editing of his profile page.
Renders different editing pages depending if the user has bought a title yet
*/
func UserProfileEdit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variable */
	var user models.Student

	/* Get Username and UUID from session */
	//username := s.Get("username")
	uuid := s.Get("userId")

	/* Load from database */
	db.C("studentList").Find(bson.M{"UUID": uuid}).One(&user)

	/* If points are high enough to buy a title display option to edit title */
	if user.Titleflag {
		r.HTML(200, "Profile/editProfileTitle", user)
	} else {
		r.HTML(200, "Profile/editProfileNoTitle", user)
	}

}

/*
Function to render /editProfileFinal
This function takes the inputs and updates the database
*/
func UserProfileEditFinal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variable */
	var user models.Student
	var title string

	/* Get Username and UUID from session */
	uuid := s.Get("userId")

	/* Load from database */
	db.C("studentList").Find(bson.M{"UUID": uuid}).One(&user)

	/* Load description from form */
	description := req.FormValue("Description")
	user.Description = description
	fmt.Println(description)
	/* Check if the user is allowed a title */
	if user.Titleflag {
		title = req.FormValue("Title")
		fmt.Println(title)
		user.Title = title
	}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"description": user.Description, "title": user.Title}},
		ReturnNew: true,
	}

	db.C("studentList").Find(bson.M{"UUID": user.UUID}).Apply(change, &user)

	r.HTML(200, "Profile/editProfileFinal", nil)
}

/*
Function to render /profile for a get request.
This function should load the user's profile page
*/
func ProfilePersonal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variable */
	var user models.Student

	/* Get UUID */
	uuid := s.Get("userId")

	/* Load from database */
	db.C("studentList").Find(bson.M{"UUID": uuid}).One(&user)

	/* Check if staff or student */
	if user.Staff == "1" {
		r.HTML(200, "Profile/profile", user)
	} else {
		r.HTML(200, "Profile/studentprofile", user)
	}
}

/*
Function to render /purchase title
Function should load the user's informatin, check if they have enough points for a title and then either
edit the Titleflag to true or leave as false. Depening on the case display the corerct webpage
*/
func PurchaseTitle(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variable */
	var user models.Student

	/* Get UUID */
	uuid := s.Get("userId")

	/* Load from database */
	db.C("studentList").Find(bson.M{"UUID": uuid}).One(&user)

	/* Check if they already have a title */
	if user.Titleflag {
		r.HTML(200, "Profile/cannotPurchaseTitle", nil)
	} else {

		/* Check if over the points threshold */
		if user.Points >= 200 {
			/* Set flag to true */
			user.Titleflag = true
			/* Reduce points */
			user.Points = user.Points - 200
			/* Update database */
			change := mgo.Change{
				Update:    bson.M{"$set": bson.M{"titleflag": user.Titleflag, "points": user.Points}},
				ReturnNew: true,
			}
			db.C("studentList").Find(bson.M{"UUID": user.UUID}).Apply(change, &user)
			/* Display page */
			r.HTML(200, "Profile/purchaseTitle", nil)
		} else {
			r.HTML(200, "Profile/cannotPurchaseTitle", nil)
		}
	}
}

/*
Function to render the homepage "/home"
Function should look up the user and then display the relevant information
*/
func StudentHomePage(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create user Variable */
	var user models.Student

	/* Get UUID from sessions */
	uuid := s.Get("userId")

	/* Load from database */
	db.C("studentList").Find(bson.M{"UUID": uuid}).One(&user)

	r.HTML(200, "homepage", user)
}
