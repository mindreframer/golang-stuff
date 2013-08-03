/*Mongodb utility methods*/
package hamster

import (
	"labix.org/v2/mgo"
)

//mongodb
type Db struct {
	Url        string
	MgoSession *mgo.Session
}

//get new db session
func (d *Db) GetSession() *mgo.Session {
	if d.MgoSession == nil {
		var err error
		d.MgoSession, err = mgo.Dial(d.Url)
		if err != nil {
			panic(err) // no, not really
		}
	}
	return d.MgoSession.Clone()

}
