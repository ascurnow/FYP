// StudentController
package controllers

import (
	"FYP/fypTest/models"
	//"fmt"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
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
	return studentList
}

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

// Function to update a students details (currently only email) based upon the
// information given in the form
func StudentUpdate(r render.Render, db *mgo.Database, req *http.Request) {
	var student models.Student
	name := req.FormValue("name")
	email := req.FormValue("email")
	//println(name)
	//println(email)
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"email": email}},
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
