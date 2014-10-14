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
	"strings"
	"time"
)

func Login(r render.Render) {
	r.HTML(200, "Login/login", nil)
}

func PostLogin(rw http.ResponseWriter, req *http.Request, db *mgo.Database, s sessions.Session) {
	var user models.Student
	email, password := req.FormValue("email"), req.FormValue("password")
	err := db.C("studentList").Find(bson.M{"email": email}).One(&user)
	//fmt.Println(email)
	//fmt.Println(user)

	// If email cannot be found or the hashed passwords do not match redirect to the login page
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
	}
	/* Set the users Username and Id in the session */

	s.Set("userName", user.Username)
	//fmt.Println(user.UUID)
	s.Set("userId", user.UUID)

	//test := s.Get("userName")
	//fmt.Println(test)
	//test2 := s.Get("userId")
	//fmt.Println(test2)

	/* Set up expire time for cookie and make the cookie */
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "cookie:username", Value: user.Username, Expires: expire}
	http.SetCookie(rw, &cookie)

	// redirect to the home page - based off if student or staff
	if strings.Contains(user.Staff, "1") {
		http.Redirect(rw, req, "/staff_home", http.StatusFound)
	} else {
		http.Redirect(rw, req, "/home", http.StatusFound)
	}

}

// Function for handling logging out of the server
func Logout(rw http.ResponseWriter, req *http.Request, s sessions.Session) {

	// Delete user session and cookie on logout
	s.Delete("userId")
	s.Delete("userName")
	expire := time.Now()
	cookie := http.Cookie{Name: "cookie:username", Value: "09b421fc90cd6242003c4e76946392a7079fd78b67b371fdf55", Expires: expire}
	http.SetCookie(rw, &cookie)

	http.Redirect(rw, req, "/", http.StatusFound)
}

func Signup(r render.Render) {
	r.HTML(200, "Login/register", nil)
}

func SignupPost(rw http.ResponseWriter, req *http.Request, db *mgo.Database) {

	/* Get name, email, and password from the form request */
	name, email, password := req.FormValue("name"), req.FormValue("email"), req.FormValue("password")

	/* Hash the password using bcrypt */
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	/* Create unique id */

	uuId := GetMD5Hash(name + password)

	PanicIf(err)
	//println(name)
	//println(email)
	//println(password)

	/* Create student as a Student struct and set UserName, Email, and Password to the inputs and hashed password */
	student := &models.Student{Username: name, Email: email, Password: hashedPassword, Staff: "0", UUID: uuId}

	/* Insert the new student into the database */
	db.C("studentList").Insert(student)
	/* Redirect the user to the login page */
	http.Redirect(rw, req, "/login", http.StatusFound)
}

func RequireLogin(rw http.ResponseWriter, req *http.Request, s sessions.Session, db *mgo.Database, c martini.Context) {
	var user models.Student
	err := db.C("studentList").Find(bson.M{"username": s.Get("userName")}).One(&user)

	if err != nil {
		http.Redirect(rw, req, "/login", http.StatusFound)
		return
	}

	// map the user to the context
	c.Map(user)
}
