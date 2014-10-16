// Data structure for students
package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"math"
)

// Create a basic student struct
type Student struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UUID         string        `bson:"UUID"`
	Username     string        `form:"Username"`
	Email        string        `form:"Email"`
	Password     []byte        `form:"Password"`
	Staff        string        `form:"1"`
	Units        []string      `form:"Unit"`
	Exp          float64       `form:100.2`
	Level        int           `form:10`
	Points       float64       `form:10.5`
	Description  string        `form:"Description"`
	Titleflag    bool
	Title        string   `form:"Title"`
	Achievements []string `form:"Achievemnt"`
}

// Function to return all Users in database
func GetAllStudents(db *mgo.Database) []Student {
	var studentList []Student
	db.C("studentList").Find(nil).All(&studentList)
	return studentList
}

/* Function to change a user's level */
func Level(username string, db *mgo.Database) Student {
	/* Create variables */
	var user Student
	var remainder float64

	/* Load information */
	db.C("studentList").Find(bson.M{"username": username}).One(&user)

	/* Calculate Level */
	remainder = math.Mod(user.Exp, 100)
	level := int((user.Exp - remainder) / 100)
	user.Level = level

	/* Update user */
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"level": user.Level}},
		ReturnNew: true,
	}
	db.C("studentList").Find(bson.M{"username": user.Username}).Apply(change, &user)

	return user
}
