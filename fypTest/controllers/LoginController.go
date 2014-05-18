// LoginController
package controllers

import (
	//"fmt"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

func newLogin(r render.Render, db *mgo.Database) {
	r.HTML(200, "Login/new-login", nil)
}

func Login(r render.Render, db *mgo.Database) {
	r.HTML(200, "Login/login", nil)
}
