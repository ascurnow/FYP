// Units.go
package models

import (
	//"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

/* Create a unit struct */
type Unit struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UUID     string        `bson:"UUID"`
	Unitname string        `form:"Unitname"`
	Students []string      `form:"Student"`
	Quizzes  []string      `form:"QuizName"`
}

/* Function that returns all units as a slice */
func GetAllUnits(db *mgo.Database) []Unit {
	var unitList []Unit
	db.C("units").Find(nil).All(&unitList)
	return unitList
}

/*
Function to remove a student from the unit
*/
func RemoveStudentFromUnit(unitname string, student string, db *mgo.Database) {
	/* Create Variables */
	var unit Unit
	var studentObj Student

	/* This section deals with removing the student from the unit database */
	db.C("units").Find(bson.M{"unitname": unitname}).One(&unit)

	/* Find the position of the student and remove it */
	for i := range unit.Students {
		if unit.Students[i] == student {
			unit.Students = unit.Students[:i+copy(unit.Students[i:], unit.Students[i+1:])]
			break
		}
	}

	/* Update Unit's Student List */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"students": unit.Students}},
		ReturnNew: true,
	}
	db.C("units").Find(bson.M{"unitname": unitname}).Apply(change, &unit)

	/* This section deals with removing the unit from the student */
	db.C("studentList").Find(bson.M{"username": student}).One(&studentObj)

	/* Find the position of the unit and remove it */
	for i := range studentObj.Units {
		if studentObj.Units[i] == unitname {
			studentObj.Units = studentObj.Units[:i+copy(studentObj.Units[i:], studentObj.Units[i+1:])]
			break
		}
	}

	/* Update the Student's Unit list */
	change2 := mgo.Change{
		Update:    bson.M{"$set": bson.M{"units": studentObj.Units}},
		ReturnNew: true,
	}
	db.C("studentList").Find(bson.M{"username": student}).Apply(change2, &studentObj)
}
