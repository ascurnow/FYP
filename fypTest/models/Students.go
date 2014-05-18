// Data structure for students
package models

import{
	
}

// Create a basic student struct
type Student struct {
	Username string `form:"Username"`
	Email    string `form:"Email"`
	Password string `form:Password"`
}

// Function to return all Users in database (NOTE this should be edited to only students later)
func GetAllStudents(db *mgo.Database) []Student {
	var studentList []User
	db.C("studentList").Find(nil).All(&studentList)
	return studentList
}