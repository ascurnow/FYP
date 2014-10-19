// Data structure for Quiz
package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Create a basic quiz struct
type Quiz struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Unitname     string        `form:"Unitname"`
	UUID         string        `bson:"UUID"`
	Quizname     string        `form:"Quizname"`
	Maxgrade     string        `form:"10"`
	Numquestions int           `form:10`
	Questions    []string      `form:"Questions"`
	Answers      []string      `form:"Answers"`
}

// Function to return all quizzes in database
func GetAllQuizzes(db *mgo.Database) []Quiz {
	var quizList []Quiz
	db.C("quiz").Find(nil).All(&quizList)
	return quizList
}
