package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store struct {
	db *mgo.Session
}

type MeasurementStats struct {
	NumLast24Hrs *DurationStat
	NumLast7Days *DurationStat
	NumLastMonth *DurationStat
	NumLastYear  *DurationStat
}
type DurationStat struct {
	Min   int64 `bson:"min"`
	Max   int64 `bson:"max"`
	Avg   int64 `bson:"avg"`
	Count int64 `bson:"count"`
}

func NewStore() *Store {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err.Error())
	}
	session.DB("run").C("measurements").EnsureIndexKey("+duration")
	return &Store{db: session}
}

func (s *Store) Add(doc *MeasurementEnded) (bson.ObjectId, error) {
	// Not sure how this method ensures unique ids but its good enough
	i := bson.NewObjectId()
	_, err := s.db.DB("run").C("measurements").Upsert(bson.M{"_id": i}, doc)
	return i, err
}

func (s *Store) GetStats() *MeasurementStats {

	return &MeasurementStats{
		NumLast24Hrs: s.getSingleStat(time.Duration(time.Hour * -24)),
		NumLast7Days: s.getSingleStat(time.Duration(time.Hour * -24 * 7)),
		NumLastMonth: s.getSingleStat(time.Duration(time.Hour * -24 * 7 * 30)),
		NumLastYear:  s.getSingleStat(time.Duration(time.Hour * -24 * 7 * 365)),
	}

}

func (s *Store) getSingleStat(t time.Duration) *DurationStat {
	c := s.db.DB("run").C("measurements")

	pipe := c.Pipe([]bson.M{
		bson.M{"$match": bson.M{"_id": bson.M{
			"$gt": bson.NewObjectIdWithTime(time.Now().Add(time.Duration(t))),
		}}},
		bson.M{"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
			"min":   bson.M{"$min": "$duration"},
			"max":   bson.M{"$max": "$duration"},
			"avg":   bson.M{"$avg": "$duration"},
		}},
	})

	var results DurationStat
	err := pipe.One(&results)
	fmt.Printf("results: %+v", results)
	if err != nil {
		return nil
	}

	return &results
}

func (s *Store) GetHighscores() []MeasurementEnded {

	// Select all measurements in the last 7 days
	q := s.db.DB("run").C("measurements").Find(bson.M{
		"_id": bson.M{
			"$gt": bson.NewObjectIdWithTime(time.Now().Add(time.Duration(time.Hour * -24 * 7))),
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
