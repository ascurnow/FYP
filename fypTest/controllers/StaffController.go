// StaffController.go
package controllers

import (
	"FYP/fypTest/models"
	"fmt"
	//"github.com/go-martini/martini"
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

/*
Function to render the Unit list page for Staff members.
This function renders the /unitList page
Gets a list of units from the database using models.GetAllUnits()
and passes them to the template renderer
*/
func StaffUnitList(r render.Render, db *mgo.Database) {
	//fmt.Println(models.GetAllUnits(db))
	r.HTML(200, "Staff/unitList", models.GetAllUnits(db))
}

/*
Function to render the staff homepage.
This function redners /staff_home page.
The page requires a list of the units that the staff member is a member of.
*/
func StaffHomePage(r render.Render, db *mgo.Database, s sessions.Session) {
	/* Create a variable to hold the staff's user data */
	var staffData models.Student

	/* Get the staff member's ID and Name */
	staffId := s.Get("userId") /* Get staff's username from the session */
	//fmt.Println(staffId)
	//staffUserName := s.Get("userName")
	//test := s.Get("userName")
	//fmt.Println(staffUserName)
	//fmt.Println(test)

	/* Get all data from the database */
	err := db.C("studentList").Find(bson.M{"UUID": staffId}).One(&staffData)
	/* This version of the search does it off their username and is inferior to the above. Once sessions for ids
	are fixed use the above version
	err := db.C("studentList").Find(bson.M{"username": staffUserName}).One(&staffData)
	*/

	/* Panic If there is an error in retriving the user's data */
	PanicIf(err)

	/* Render the HTML page */
	r.HTML(200, "Staff/homepage", staffData)
}

/*
 Function to render the staff unit page. Loads that users units and displays them.
This page functions very similarly to this section from the home page.
It loads the user's UUID from the session and then collects the user's information.
After this it then passes this to the template engine to be displayed.
*/
func StaffMyUnit(r render.Render, db *mgo.Database, s sessions.Session) {
	var user models.Student

	/* Retrive user's ID from sessions */
	userId := s.Get("userId")
	//username := s.Get("userName")
	//fmt.Println(username)

	/* Find the user and store his information in the user variable */
	db.C("studentList").Find(bson.M{"UUID": userId}).One(&user)
	//fmt.Println(user)
	//fmt.Println(user.Username)

	/* Check if staff or student */
	if user.Staff == "1" {
		r.HTML(200, "Staff/staffUnits", user)
	} else {
		r.HTML(200, "Student/studentUnits", user)
	}
}

/*
Function to render the staff version of the unit page.
Function renders /unit

This function gives the page all of the unit information including the unit name,
students currently enrolled in the unit, and the quizzes available for the unit.
NOTE: It might be worth adding another section on this page later for adding quizzes if this isn't done through another
page or function.
*/
func StaffUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variables as required for storing all unit information */
	var unit models.Unit
	var unitname string /* String for storing the unitname from the Form */
	/* Change this to use Unit ID later on when that is fixed */

	// Get the name of the unit to search for
	unitname = req.FormValue("selectUnit")

	// Request Unit information - based upon the selected unit
	err := db.C("units").Find(bson.M{"unitname": unitname}).One(&unit)
	PanicIf(err)

	//fmt.Println(err)
	//fmt.Println(unit)
	//fmt.Println(unit.getUnitname())
	//fmt.Println(unit.Unitname)

	r.HTML(200, "Staff/staffUnitPage", unit)
}

/*
	Function to remove the selected student from the unit
*/
func StaffRemoveStudentFromUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create variables as required */
	//var unit models.Unit
	var studentToRemoveName string /* String for storing the students name */
	var unitToRemoveFrom string    /* String for storing the unit name */

	//fmt.Println(data)

	/* Get Name and Id of the student to remove and unit from which to remove from */
	studentToRemoveName, unitToRemoveFrom = req.FormValue("selectStudentToRemove"), req.FormValue("UnitName")
	//fmt.Println(studentToRemoveName)
	//fmt.Println(unitToRemoveFrom)

	/* Delete the student from the Unit */
	models.RemoveStudentFromUnit(unitToRemoveFrom, studentToRemoveName, db)

	//fmt.Println(err)
	//PanicIf(err)

	data := struct {
		Unitname string `form:"Unitname"`
	}{
		unitToRemoveFrom,
	}
	r.HTML(200, "Staff/staffRemoveStudentFromUnitPage", data)
}

/*
	Function to add students to a unit
	This function will display the /addStudentToUnit page. This page requires a list of the students in the database,
	the unit itself, and will then be required to determine the students who are not already in the unit and display these
	students. The user then selects the students he/she wishes to add to the unit in a form and they are added.
*/
func StaffAddStudentToUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create variables as required */
	var unit models.Unit
	var studentList []models.Student
	//var studentToAdd models.Student
	var studentsNotInUnit []string /* slice to store the UUID of the students not currently in the unit */
	//var studentsNotInUnitName []string /* slice to store the names of the students not currently in the unit */
	var unitName string   /* String for storing the unit's name */
	studentFound := false /* Flag for finding student */

	/* Create struct for pushing data to the template engine */
	type studentsAdd struct {
		Unitname string   `form:"Unitname"`
		Students []string `form:"Username"`
	}

	var studentsToAdd studentsAdd

	/* Get the unit's unitname */
	unitName = req.FormValue("UnitName")
	/* Find UUID */
	unitId := GetMD5Hash(unitName)
	/* Load the unit's information */
	err := db.C("units").Find(bson.M{"UUID": unitId}).One(&unit)

	PanicIf(err)

	/* Load all students */
	err = db.C("studentList").Find(nil).All(&studentList)

	/* Loop over student list and find the students not in the unit */
	for i := range studentList {
		for j := range unit.Students {
			if studentList[i].Username == unit.Students[j] { /* Check and see if the student is in the unit */
				studentFound = true /* If the student is found to be in the unit */
			}
		}
		if studentFound == false { /* If the student was not found */
			/* Add student to the list of students not in the unit */
			studentsNotInUnit = append(studentsNotInUnit, studentList[i].Username)
		}
		studentFound = false /* Reset flag */
	}
	/* Loop over UUID's in studentsNotInUnit and find their username */
	/*for k := range studentsNotInUnit {
		err = db.C("studentList").Find(bson.M{"username": studentsNotInUnit[k]}).One(&studentToAdd)
		PanicIf(err)
		studentsNotInUnitName = append(studentsNotInUnitName, studentToAdd.Username)
		fmt.Println(studentsNotInUnitName)
	}*/

	studentsToAdd.Unitname = unit.Unitname
	studentsToAdd.Students = studentsNotInUnit

	//fmt.Println(stufftoadd)
	//fmt.Println(studentsNotInUnit)
	//fmt.Println(err)
	PanicIf(err)

	r.HTML(200, "Staff/staffAddStudentToUnitPage", studentsToAdd)
}

/*
	Function to add students to the unit - final phase
	Function renders /staffFinalAddStudentToUnitPage
	This function takes the student names and the unit name from the post function and then add these studnets to the unit.
	Updates the following : Unit.Students and Students.Units for each student.
*/
func StaffAddStudentToUnitFinal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create variables as required */
	var unit models.Unit
	var student models.Student

	/* Get student list and unitname */
	req.ParseForm()
	studentsToAdd := req.Form["selectStudentsToAdd"]
	unitName := req.Form["UnitName"]

	//fmt.Println(unitName)
	//fmt.Println(studentsToAdd)

	/* Get the UUID for the unit */
	unitId := GetMD5Hash(unitName[0])
	/* Load the unit information */
	err := db.C("units").Find(bson.M{"UUID": unitId}).One(&unit)
	PanicIf(err)

	//fmt.Println(unit.Students)

	/* Add all students in studentsToAdd to the unit */
	for i := range studentsToAdd {
		unit.Students = append(unit.Students, studentsToAdd[i])
		//fmt.Println(unit.Students)
	}
	/* Create change variable and put in the information */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"students": unit.Students}},
		ReturnNew: true,
	}
	/* Update the unit with the new students */
	db.C("units").Find(bson.M{"UUID": unitId}).Apply(change, &unit)

	/* Update the Student's Unit list */
	for j := range studentsToAdd {
		/* Load the students unit list */
		err = db.C("studentList").Find(bson.M{"username": studentsToAdd[j]}).One(&student)
		PanicIf(err)
		/* Add the unit to the students list of units */
		student.Units = append(student.Units, unitName[0])

		/* Set change variable to update the new complete unit list */
		change = mgo.Change{
			Update:    bson.M{"$set": bson.M{"units": student.Units}},
			ReturnNew: true,
		}
		/* Update the student's unit list */
		studentName := studentsToAdd[j]
		db.C("studentList").Find(bson.M{"username": studentName}).Apply(change, &student)

	}
	addStudentList := struct {
		Unitname string   `form:"Unitname"`
		Students []string `form:"Username"`
	}{
		unitName[0],
		studentsToAdd,
	}
	//fmt.Println(addStudentList)

	r.HTML(200, "Staff/staffAddStudentToUnitPageFinal", addStudentList)

}

/*
	Function to create a new unit
	This function renders /staffCreateUnit
	Function takes the unit name, checks to see if the unit is unique, and if it is makes a new unit
*/
func StaffCreateUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create Variables */
	var unit models.Unit
	var unitList []models.Unit
	copy := false

	/* Get unitname from the form */
	unitName := req.FormValue("unitCode")
	//fmt.Println(unitName)

	/* Create UUID for the unit */
	uuId := GetMD5Hash(unitName)
	/* Assign UUID and Unitname to the variable */
	unit.Unitname = unitName
	unit.UUID = uuId

	/* Check and see if the unit is already created */
	/* Load the current list of units */
	unitList = models.GetAllUnits(db)

	/* Loop over all units and check if the UUID of the unit to be created is present (therefore not unique) */
	for i := range unitList {
		if unitList[i].UUID == unit.UUID {
			copy = true
		}
	}
	//fmt.Println(unit)

	/* Create the unit in the database if the unit is not a copy */
	if copy == false {
		err := db.C("units").Insert(unit)
		PanicIf(err)
	}

	r.HTML(200, "Staff/staffCreateUnit", unit)
}

/*
Function to render /staffQuizPage
This function needs to display a page of the results for the quiz
*/
func StaffQuizPage(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
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
		r.HTML(200, "Staff/staffQuizPageNoResults", noData)
	} else { /* Send data struct */
		r.HTML(200, "Staff/staffQuizPage", Data)
	}

	//fmt.Println(Data)

}

/*
Function to add a result to a quiz.
This page needs the following info:
- Quizuuid
- Quizname
- Max grade possible for the quiz
- Unitname

Page has a form that collects the following:
- Student Name
- Score/Result
*/
func StaffAddQuizResult(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Get quizuuid and name from the form */
	quizname := req.FormValue("quizname")

	//fmt.Println(quizuuid)
	//fmt.Println(quizname)
	/* Look up quiz - find max_grade and Unitname. */
	maxGrade := "10"       // Temp hard code
	unitName := "FIT 1029" //Temp hard code

	/* Create data structure */
	quizData := struct {
		Quizname string `form:"Quizname"`
		Maxgrade string `form:"10"`
		Unitname string `form:"Unitname"`
	}{
		quizname,
		maxGrade,
		unitName,
	}
	//fmt.Println(quizData)

	/* Render page and push Quiz Name and Max grade to be displayed */
	r.HTML(200, "Staff/staffAddResultToQuiz", quizData)
}

/*
Function for rendering /staffAddResultToQuizFinal
This function needs to create a new result in the data base using the form inputs
*/
func StaffAddQuizResultFinal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Get information from form */
	quizname := req.FormValue("Quizname")
	unitname := req.FormValue("Unitname")
	grade := req.FormValue("StudentResult")
	studentName := req.FormValue("StudentName")
	maxGrade := req.FormValue("MaxGrade")

	//fmt.Println(studentName)

	/* Get exp change */
	exp := models.Experience(grade, maxGrade)
	/* Get UUID */
	uuidString := studentName + quizname
	uuid := GetMD5Hash(uuidString)

	/* Create result variable from models.Result */
	result := &models.Result{Quiz: quizname, Unit: unitname, Studentname: studentName, Score: grade, UUID: uuid, Maxgrade: maxGrade, Exp: exp}

	/* Push new result to database */
	db.C("result").Insert(result)

	/* Update student's XP */
	var student models.Student
	/* Look up student's current XP */
	db.C("studentList").Find(bson.M{"username": studentName}).One(&student)
	fmt.Println(student)

	exp = student.Exp + exp

	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"exp": exp}},
		ReturnNew: true,
	}
	db.C("studentList").Find(bson.M{"username": studentName}).Apply(change, &student)

	/* Run EXP checker function to udpate achievements etc */
	/* TODO Check LEVEL */
	/* TODO CHECK ACHIEVEMENTS */

	/* Go to staff home page */
	http.Redirect(rw, req, "/staff_home", http.StatusFound)
}

/*
Function for /editResultInUnit
Function takes the student's name and quizname from the form, loads the results document details
*/
func StaffEditResultInUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variables */
	var result models.Result

	/* Get info from form */
	studentname := req.FormValue("studentToEdit")
	quizname := s.Get("quizname")
	//fmt.Println(studentname)
	//fmt.Println(quizname)

	/* Look up result */
	err := db.C("result").Find(bson.M{"quiz": quizname, "studentname": studentname}).One(&result)
	PanicIf(err)

	//fmt.Println(result)

	r.HTML(200, "Staff/staffEditResultInQuiz", result)
}

/*
Function for /staffEditResultInQuizFinal
Function takes the inputs from the form and updates the database wiht the details.
*/
func StaffEditResultInUnitFinal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request, s sessions.Session) {
	/* Create variables */
	var result models.Result

	/* Get Information from form */
	studentname := req.FormValue("Studentname")
	quizname := req.FormValue("Quiz")
	grade := req.FormValue("StudentResult")
	maxGrade := req.FormValue("MaxGrade")

	fmt.Println(grade)

	/* Calculate new exp */
	exp := models.Experience(grade, maxGrade)
	fmt.Println(exp)

	var oldResult models.Result

	db.C("result").Find(bson.M{"quiz": quizname, "studentname": studentname}).One(&oldResult)

	/* Update the result document */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"exp": exp, "score": grade}},
		ReturnNew: true,
	}
	db.C("result").Find(bson.M{"quiz": quizname, "studentname": studentname}).Apply(change, &result)

	/* Calculate exp difference */
	oldExp := oldResult.Exp
	difference := oldExp - exp

	/* Load Total experience from database */
	var student models.Student
	/* Look up student's current XP */
	db.C("studentList").Find(bson.M{"username": studentname}).One(&student)

	/* New total experience */
	updateExp := student.Exp - difference

	change2 := mgo.Change{
		Update:    bson.M{"$set": bson.M{"exp": updateExp}},
		ReturnNew: true,
	}
	/* Update student exp */
	db.C("studentList").Find(bson.M{"username": studentname}).Apply(change2, &student)

	/* Clean up sessions */
	s.Delete("quizname")

	/* Go to staff home page */
	r.HTML(200, "Staff/staffEditResultInQuizFinal", nil)
}

/*
This function renders /addQuizToUnit
Function for adding a quiz to a unit
*/
func StaffAddQuizToUnit(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create Quiz Variable */
	var quiz models.Quiz
	var checkQuiz []models.Quiz
	found := false

	/* Get Unitname from form */
	unitname := req.FormValue("UnitName")
	quizname := req.FormValue("Quizname")

	//fmt.Println(unitname)
	//fmt.Println(quizname)

	/* Get hash for quiz UUID - quizname + unitname */
	stringToHash := quizname + unitname
	quizuuid := GetMD5Hash(stringToHash)
	//fmt.Println(quizuuid)

	quiz.Unitname = unitname
	quiz.Quizname = quizname
	quiz.UUID = quizuuid

	fmt.Println(quiz)

	/* Check if this a duplicate UUID */
	db.C("quizzes").Find(bson.M{"unitname": unitname}).All(&checkQuiz)

	for i := range checkQuiz {
		if checkQuiz[i].UUID == quizuuid {
			found = true
		} else {
			found = false
		}
	}
	if found {
		/* Quiz already made - go to a page saying this */
		r.HTML(200, "Staff/QuizAlreadyMade", nil)
	} else {
		/* Quiz not found - go to page to make quiz */
		r.HTML(200, "Staff/staffAddQuizToUnit", quiz)
	}
}

/*
Function renders /stafAddQuizFinal
Function takes the inputs from the form and adds the quiz, questions, and answers to the database.
*/

func StaffAddQuizToUnitFinal(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create quiz and unit vars */
	var quiz models.Quiz
	var unit models.Unit

	/* Get information from form */
	quizname := req.FormValue("Quizname")
	quizuuid := req.FormValue("UUID")
	unitname := req.FormValue("Unitname")
	maxgrade := req.FormValue("MaxGrade")

	/* Update quiz variable */
	quiz.Quizname = quizname
	quiz.Unitname = unitname
	quiz.Maxgrade = maxgrade
	quiz.UUID = quizuuid

	/* Store quiz in database */
	db.C("quiz").Insert(quiz)

	/* Load Unit information */
	db.C("units").Find(bson.M{"unitname": quiz.Unitname}).One(&unit)

	unit.Quizzes = append(unit.Quizzes, quizname)

	/* Update the result document */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"quizzes": unit.Quizzes}},
		ReturnNew: true,
	}
	/* Update Unit information with new quiz */
	db.C("units").Find(bson.M{"unitname": unitname}).Apply(change, &unit)

	/* Process and store questions in database */

	/* Process and store answers in database */

	r.HTML(200, "Staff/staffAddQuizToUnitFinal", nil)
}

/*
Function dealing with removing a quiz from the database and unit
*/
func StaffRemoveQuiz(r render.Render, db *mgo.Database, rw http.ResponseWriter, req *http.Request) {
	/* Create variables */
	var unit models.Unit

	/* Get info from the form */
	quizname := req.FormValue("selectQuizToRemove")
	unitname := req.FormValue("UnitName")

	/* This section deals with removing the quiz from the unit database */
	db.C("units").Find(bson.M{"unitname": unitname}).One(&unit)

	/* Find the position of the quiz and remove it */
	for i := range unit.Quizzes {
		if unit.Quizzes[i] == quizname {
			unit.Quizzes = unit.Quizzes[:i+copy(unit.Quizzes[i:], unit.Quizzes[i+1:])]
			break
		}
	}
	/* Update Unit's Quiz List */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"quizzes": unit.Quizzes}},
		ReturnNew: true,
	}
	db.C("units").Find(bson.M{"unitname": unitname}).Apply(change, &unit)

	/* Get hash */
	stringToHash := quizname + unitname
	quizuuid := GetMD5Hash(stringToHash)

	/* Delete quiz from quiz database */
	db.C("quiz").Remove(bson.M{"UUID": quizuuid})

	r.HTML(200, "Staff/staffRemoveQuiz", nil)
}
