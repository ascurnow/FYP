// StudentController
package controllers

import (
	//"fmt"
	//"fypTest/models"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	//"html/template"
	//"net/http"
	"labix.org/v2/mgo"
	//"os"
)

// Create a basic User struct
type User struct {
	Username string `form:"Username"`
	Email    string `form:"Email"`
	Password string `form:Password"`
}

// Function to return all Users in database (NOTE this should be edited to only students later)
func GetAllStudents(db *mgo.Database) []User {
	var studentList []User
	db.C("studentList").Find(nil).All(&studentList)
	return studentList
}

func StudentIndex(r render.Render, db *mgo.Database) {
	r.HTML(200, "Student/index", GetAllStudents(db))
}

//func StudentAdd(db *mgo.Database, res http.ResponseWriter, req *http.Request) {
//	db.C("Users").Insert()
//}

func StudentAddPage(r render.Render, db *mgo.Database) {
	r.HTML(200, "Student/add", nil)
}

func StudentRemovePage(r render.Render, db *mgo.Database) {
	r.HTML(200, "Student/remove", nil)
}
