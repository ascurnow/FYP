// Questions
package models

import (
	//"fmt"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

/* Create a question struct */
type Question struct {
	Id             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UUID           string        `bson:"UUID"`
	Questionnumber string        `form:"1"`
	Quiz           string        `form:"QuizName"`
	Unit           string        `form:"UnitName"`
	Question       string        `form:"Question"`
	Answers        []string      `form:"Question"`
	CorrectAnswer  []string      `form:"1"`
}
