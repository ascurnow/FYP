// Data structure for Quiz
package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Create a basic quiz struct
type Quiz struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Unitname string        `form:"Unitname"`
	UUID     string        `bson:"UUID"`
	Quizname string        `form:"Quizname"`
	Maxgrade string        `form:"10"`
}

/* Quetsions and Answers associated with the quiz are linked by the Quiz UUID */

// Function to return all quizs in database (NOTE this should be edited to only quizzes in the selected unit later)
func GetAllQuizzes(db *mgo.Database) []Quiz {
	var quizList []Quiz
	db.C("quiz").Find(nil).All(&quizList)
	return quizList
}
