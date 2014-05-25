// LoginController
package controllers

import (
	"FYP/fypTest/models"
	"code.google.com/p/go.crypto/bcrypt"
	//"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func Login(r render.Render) {
	r.HTML(200, "Login/login", nil)
}

func PostLogin(rw http.ResponseWriter, req *http.Request, db *mgo.Database, s sessions.Session) {
	var user models.Student
	email, password := req.FormValue("email"), req.FormValue("password")
	err := db.C("studentList").Find(bson.M{"email": email}).One(&user)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
	}

	s.Set("userId", user.Username)
	http.Redirect(rw, req, "/students", http.StatusFound)
}

func Logout(rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	s.Delete("userId")
	http.Redirect(rw, req, "/", http.StatusFound)
}

func Signup(r render.Render) {
	r.HTML(200, "Login/register", nil)
}

func SignupPost(rw http.ResponseWriter, req *http.Request, db *mgo.Database) {
	name, email, password := req.FormValue("name"), req.FormValue("email"), req.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	PanicIf(err)
	//println(name)
	//println(email)
	//println(password)

	student := &models.Student{Username: name, Email: email, Password: hashedPassword}
	db.C("studentList").Insert(student)
	http.Redirect(rw, req, "/login", http.StatusFound)
}

func RequireLogin(rw http.ResponseWriter, req *http.Request, s sessions.Session, db *mgo.Database, c martini.Context) {
	var user models.Student
	err := db.C("studentList").Find(bson.M{"username": s.Get("userId")}).One(&user)

	if err != nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return
	}

	c.Map(user)
}
