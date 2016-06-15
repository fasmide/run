package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store struct {
	db *mgo.Session
}

func NewStore() *Store {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err.Error())
	}
	session.DB("run").C("measurements").EnsureIndexKey("+durationNs")
	return &Store{db: session}
}

func (s *Store) Add(doc *MeasurementEnded) (bson.ObjectId, error) {
	// Not sure how this method ensures unique ids but its good enough
	i := bson.NewObjectId()
	_, err := s.db.DB("run").C("measurements").Upsert(bson.M{"_id": i}, doc)
	return i, err
}

func (s *Store) GetHighscores() []MeasurementEnded {

	// Select all measurements in the last 7 days
	q := s.db.DB("run").C("measurements").Find(bson.M{
		"_id": bson.M{
			"gt": bson.NewObjectIdWithTime(time.Now().Add(time.Duration(time.Hour * -24 * 7))),
		},
	})

	q.Sort("duration")

	var results []MeasurementEnded
	err := q.All(&results)

	if err != nil {
		panic(err.Error())
	}
	return results
}
