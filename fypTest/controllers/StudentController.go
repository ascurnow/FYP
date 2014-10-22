// StudentController
package controllers

import (
	"FYP/fypTest/models"
	//"fmt"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	//"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
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
	//unit := req.FormValue("unit")

	//println(name)
	//println(email)
	//println(staff)
	//println(unit)

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
// Renders a page asking to confirm the deletion
func StudentRemovePage(r render.Render, db *mgo.Database, req *http.Request) {
	var studentToRemove models.Student
	name := req.FormValue("selectRemoveStudent")
	//println(name)

	db.C("studentList").Find(bson.M{"username": name}).One(&studentToRemove)
	//println(studentToRemove.Username)
	r.HTML(200, "Student/remove", studentToRemove)
}

// Function to remove the student
func StudentRemove(r render.Render, db *mgo.Database, req *http.Request) {
	option := req.FormValue("confirmDelete")
	name := req.FormValue("studentName")
	email := req.FormValue("studentEmail")
	//println(option)
	//println(name)
	//println(email)
	if option == "TRUE" {
		db.C("studentList").Remove(bson.M{"username": name, "email": email})
		r.HTML(200, "Student/removalConfirmed", option)
	} else {
		r.HTML(200, "Student/notRemoved", option)
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

	/* Check if staff or student */
	if user.Staff == "1" {
		/* If points are high enough to buy a title display option to edit title */
		if user.Titleflag {
			r.HTML(200, "Profile/editProfileTitle", user)
		} else {
			r.HTML(200, "Profile/editProfileNoTitle", user)
		}
	} else {
		/* If points are high enough to buy a title display option to edit title */
		if user.Titleflag {
			r.HTML(200, "Profile/editProfileTitleStudent", user)
		} else {
			r.HTML(200, "Profile/editProfileNoTitleStudent", user)
		}
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
	if description == "" {
		user.Description = user.Description
	} else {
		user.Description = description
	}

	//fmt.Println(description)

	/* Check if the user is allowed a title */
	if user.Titleflag {
		title = req.FormValue("Title")
		//fmt.Println(title)
		/* Check if blank */
		if title == "" {
			user.Title = user.Title
		} else {
			user.Title = title
		}
	}
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"description": user.Description, "title": user.Title}},
		ReturnNew: true,
	}

	db.C("studentList").Find(bson.M{"UUID": user.UUID}).Apply(change, &user)

	/* Check if staff or student */
	if user.Staff == "1" {
		r.HTML(200, "Profile/editProfileFinal", nil)
	} else {
		r.HTML(200, "Profile/editProfileFinalStudent", nil)
	}

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
	staff := s.Get("Staff")
	//fmt.Println(staff)

	/* Load from database */
	db.C("studentList").Find(bson.M{"UUID": uuid}).One(&user)

	/* Check if staff or student */
	if staff == "1" {
		r.HTML(200, "Profile/profile", user)
	} else {
		r.HTML(200, "Profile/studentprofile", user)
	}
}

/*
Function to render /profile off a POST.
This function needs to load the user's profile page based off the selected info in the form
*/
func Profile(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create Variable */
	var user models.Student

	/* Get info from form */
	username := req.FormValue("username")
	staff := s.Get("Staff")
	//fmt.Println(staff)

	/* Load from database */
	db.C("studentList").Find(bson.M{"username": username}).One(&user)

	/* Check if staff or student */
	if staff == "1" {
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
	staff := s.Get("Staff")

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
			/* Check if staff or student */
			if staff == "1" {
				/* Display page */
				r.HTML(200, "Profile/purchaseTitle", nil)
			} else {
				r.HTML(200, "Profile/purchaseTitleStudent", nil)
			}
		} else {
			/* Check if staff or student */
			if staff == "1" {
				r.HTML(200, "Profile/cannotPurchaseTitle", nil)
			} else {
				r.HTML(200, "Profile/cannotPurchaseTitleStudent", nil)
			}
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

/*
Function to render the student version of the unit page.
Function renders /unitStudent

This function needs to show all quizzes in the unit and
*/
func StudentUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variables as required for storing all unit information */
	var unit models.Unit
	var unitname string /* String for storing the unitname from the Form */

	// Get the name of the unit to search for
	unitname = req.FormValue("selectUnit")

	// Request Unit information - based upon the selected unit
	err := db.C("units").Find(bson.M{"unitname": unitname}).One(&unit)
	//fmt.Println(err)
	PanicIf(err)

	//fmt.Println(err)
	//fmt.Println(unit)
	//fmt.Println(unit.getUnitname())
	//fmt.Println(unit.Unitname)

	r.HTML(200, "Student/studentUnit", unit)
}

/*
Function to render /studentQuizPage
This function needs to display a page of the results for the quiz and allow the user to take the quiz
*/
func StudentQuizPage(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create Variables */
	var results []models.Result
	//var maxGrade string
	//var resultsArray string
	//var nameArray string

	/* Create struct for pushing data to the template engine */
	type studentData struct {
		Unitname    string  `form:"Unitname"`
		Quizname    string  `form"Quizname"`
		Maxgrade    string  `form:"10"`
		Result      string  `form:"10"`
		Studentname string  `form:"Username"`
		Exp         float64 `form:25`
		UUID        string  `bson:"UUID"`
	}

	/* Get quiz name from the form */
	quizName := req.FormValue("selectQuiz")
	/* Get Unit name from the form */
	//unitName := req.FormValue("UnitName")

	/* Load the quiz results from results database */
	err := db.C("result").Find(bson.M{"quiz": quizName}).All(&results)
	PanicIf(err)

	s.Set("quizname", quizName)
	//fmt.Println(results)
	//maxGrade = results[0].Maxgrade

	var Data []studentData
	//resultsArray := make([]string, len(results))
	//nameArray := make([]string, len(results))

	/* Loop over all results, store value and names */
	for i := range results {
		//resultsArray = results[i].Score
		//nameArray = results[i].Studentname
		temp := studentData{Unitname: results[i].Unit, Quizname: quizName, Maxgrade: results[i].Maxgrade, Result: results[i].Score, Studentname: results[i].Studentname, Exp: results[i].Exp, UUID: results[i].UUID}
		Data = append(Data, temp)
	}

	/* If no results return just the quizname */
	if len(results) == 0 {
		noData := struct {
			Quizname string `form"Quizname"`
		}{
			quizName,
		}
		r.HTML(200, "Student/studentQuizPageNoResults", noData)
	} else { /* Send data struct */
		r.HTML(200, "Student/studentQuizPage", Data)
	}
}

/*
Function to render /studentlist
This function needs to return a list of all students and provide an option to look at their profiles.
*/
func StudentList(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	r.HTML(200, "Student/studentlist", GetAllStudents(db))
}

/*
Function to render /quiz
This function needs to return the questions for the quiz. It should keep looping on itself until all questions are done.
On the last question it should provide a page that has a final submit option for grading.
*/
func StudentQuiz(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variables */
	var quiz models.Quiz
	var result models.Result
	var student models.Student
	var qscore float64
	var mg float64
	var numq float64
	var oldScoreF float64

	username := s.Get("userName")

	/* Get form information */
	quizname := req.FormValue("quizname")
	questionnumber := req.FormValue("questionNum")

	/* Convert questionnumber to int */
	qnum, err1 := strconv.Atoi(questionnumber)
	PanicIf(err1)

	/* Load quiz from database */
	err := db.C("quiz").Find(bson.M{"quizname": quizname}).One(&quiz)
	PanicIf(err)
	//fmt.Println(quiz)

	/* Calculate score of each question */
	maxgrade := quiz.Maxgrade
	mgrade, err2 := strconv.Atoi(maxgrade)
	mg = float64(mgrade)
	//fmt.Println(err2)
	PanicIf(err2)
	numq = float64(quiz.Numquestions)
	qscore = mg / numq

	/* Create or update result information depending on qnum */
	if qnum == 0 {
		/* Coming from quiz select. Need to make a new result in the database */

		/* Create hash for result UUID - Name + Quiz */
		stringToHash := username.(string) + quiz.Quizname
		resultuuid := GetMD5Hash(stringToHash)

		/* Fill in information */
		result.UUID = resultuuid
		result.Quiz = quiz.Quizname
		result.Quizuuid = quiz.UUID
		result.Unit = quiz.Unitname
		result.Studentname = username.(string)
		result.Maxgrade = quiz.Maxgrade
		result.Exp = float64(0)
		result.Score = "0"

		/* Insert the document into the database */
		db.C("result").Insert(result)
	} else {
		/* Coming from in progress quiz. Need to update result information */

		/* Load result */
		err4 := db.C("result").Find(bson.M{"QUUID": quiz.UUID, "studentname": username}).One(&result)
		//fmt.Println(err4)
		PanicIf(err4)
		//fmt.Println(result)
		/* Get submitted answer */
		studentAnswer := req.FormValue("Answer")

		if studentAnswer == quiz.Answers[qnum-1] {
			/* Answer was correct - add score to the result */
			resultScore := result.Score
			//fmt.Println(reflect.TypeOf(resultScore))
			oldScoreInt, err3 := strconv.Atoi(resultScore)
			//fmt.Println(err3)
			PanicIf(err3)
			/* Add old score and new score */
			oldScoreF = float64(oldScoreInt)
			newScore := int(oldScoreF + qscore)
			result.Score = strconv.FormatInt(int64(newScore), 10)
			change := mgo.Change{
				Update:    bson.M{"$set": bson.M{"score": result.Score}},
				ReturnNew: true,
			}
			/* Push to database */
			db.C("result").Find(bson.M{"UUID": result.UUID}).Apply(change, &result)
		}
	}

	/* Check we have questions to answer */
	if qnum < quiz.Numquestions {
		/* Questions till to answer - get next questions and pass to template engine */
		currQuestion := quiz.Questions[qnum]
		//fmt.Println(currQuestion)
		/* Incrememt qnum */
		qnum++
		qnumS := strconv.FormatInt(int64(qnum), 10)

		/* Insert Question and new qnum into struct for pushing to page */
		data := struct {
			Quizname string `form:"Quizname"`
			Qnum     string `form:"Qnum"`
			Question string `form:"Question"`
		}{
			quiz.Quizname,
			qnumS,
			currQuestion,
		}

		r.HTML(200, "Student/quiz", data)
	} else {
		/* No more questions - show quiz finished screen */
		/* Load result information to display */
		db.C("result").Find(bson.M{"quizuuid": quiz.UUID, "studentname": username}).One(&result)
		exp := models.Experience(result.Score, quiz.Maxgrade)

		change := mgo.Change{
			Update:    bson.M{"$set": bson.M{"exp": exp}},
			ReturnNew: true,
		}
		db.C("result").Find(bson.M{"UUID": result.UUID}).Apply(change, &result)

		/* Update student exp and Level */
		/* Load student information */
		db.C("studentList").Find(bson.M{"username": username}).One(&student)

		/* Update user exp */
		student.Exp = student.Exp + result.Exp

		change2 := mgo.Change{
			Update:    bson.M{"$set": bson.M{"exp": student.Exp}},
			ReturnNew: true,
		}
		db.C("studentList").Find(bson.M{"username": student.Username}).Apply(change2, &student)

		/* Update level */
		student = models.Level(student.Username, db)

		r.HTML(200, "Student/quizDone", result)
	}
}

/*
Function to render /addUnitsToStudent
Needs to find all units that the student isn't enrolled in and display them.
*/
func StudentAddUnits(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create variables as required */
	var units []models.Unit
	var user models.Student
	unitFound := false          /* Flag for finding unit in user list */
	var unitsAvailable []string /* Slice to store the names of all units available */

	/* Create struct for pushing data to the template engine */
	type unitsToAdd struct {
		Username string   `form:"Username"`
		Units    []string `form:"Unitname"`
	}

	/* Get username from form */
	username := req.FormValue("username")

	/* Load info from database */
	err := db.C("studentList").Find(bson.M{"username": username}).One(&user)
	PanicIf(err)
	err = db.C("units").Find(nil).All(&units)
	PanicIf(err)

	/* Find units user is not currently in */
	for i := range units {
		for j := range user.Units {
			/* Check and see if the unit is in the student's list */
			if units[i].Unitname == user.Units[j] {
				/* Student is enrolled in the unit */
				unitFound = true
			}
		}
		/* If the student is not enrolled */
		if unitFound == false {
			/* Add unit to the list of units that can be enrolled in */
			unitsAvailable = append(unitsAvailable, units[i].Unitname)
		}
		unitFound = false /* Reset Flag */
	}

	/* Update unitsToAdd */
	var unitAdd unitsToAdd
	unitAdd.Units = unitsAvailable
	unitAdd.Username = username

	r.HTML(200, "Student/addUnits", unitAdd)
}

/*
Function to render /addUnitsToStudentFinal
This function takes inputs from the form and then enrolls the user into the selected units.
Function updates the user's information in the database and the Unit information
*/
func StudentAddUnitsFinal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create variables as required */
	var unit models.Unit
	var user models.Student
	var currUnit string

	/* Get Information from form */
	req.ParseForm()
	unitsToAdd := req.Form["selectUnits"]
	username := req.Form["username"]

	/* Load student information */
	err := db.C("studentList").Find(bson.M{"username": username[0]}).One(&user)
	PanicIf(err)

	/* Add all units in unitsToAdd to the user */
	for i := range unitsToAdd {
		user.Units = append(user.Units, unitsToAdd[i])
	}
	/* Create change variable to put in the information */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"units": user.Units}},
		ReturnNew: true,
	}
	/* Update the user with the new units */
	db.C("studentList").Find(bson.M{"username": username[0]}).Apply(change, &user)

	/* Update the Units student list */
	for j := range unitsToAdd {
		/* Load the Unit's student list */
		err = db.C("units").Find(bson.M{"unitname": unitsToAdd[j]}).One(&unit)
		PanicIf(err)

		/* Add the student to the Unit's list of student */
		unit.Students = append(unit.Students, username[0])

		/* Set change variable to update the database */
		change = mgo.Change{
			Update:    bson.M{"$set": bson.M{"students": unit.Students}},
			ReturnNew: true,
		}

		/* Update the unit's student list*/
		currUnit = unitsToAdd[j]
		db.C("units").Find(bson.M{"unitname": currUnit}).Apply(change, &unit)

	}
	/* Create temp struct to push data back to page */
	addUnitList := struct {
		Username string   `form:"Username"`
		Units    []string `form:"Unitname"`
	}{
		username[0],
		unitsToAdd,
	}

	r.HTML(200, "Student/addUnitsFinal", addUnitList)
}
