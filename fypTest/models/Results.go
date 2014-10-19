// Results
package models

import (
	//"fmt"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
)

/* Create a result struct */
type Result struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UUID        string        `bson:"UUID"`
	Quizuuid    string        `bson:"QUUID"`
	Quiz        string        `form:"Quizname"`
	Unit        string        `form:"Unitname"`
	Studentname string        `form:"Studentname"`
	Score       string        `form:"10"`
	Maxgrade    string        `form:"10"`
	Exp         float64       `form:25`
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func Experience(grade string, maxgrade string) float64 {
	/* Create Variables */
	var result float64
	var exp float64
	var max_exp float64

	max_exp = 25
	//fmt.Println(maxgrade)
	//fmt.Println(grade)

	/* Convert maxgrade from sting to int */
	maxgradeInt, err := strconv.Atoi(maxgrade)
	gradeInt, err2 := strconv.Atoi(grade)
	PanicIf(err)
	PanicIf(err2)

	result = float64(gradeInt) / float64(maxgradeInt)
	exp = max_exp * result

	return exp
}
