package dataStruct

import (
	"github.com/go-xorm/core"
	"time"
)

func init() {
	core.NewPrefixMapper(core.SnakeMapper{}, "gctf_")
}

type User struct {
	Id              int64 `xorm:"autoincr pk 'id'"`
	Username        string `xorm:"unique"`
	Password          string
	Email           string `xorm:"unique"`
	RegisterTime   time.Time `xorm:"updated"`
	ProblemsID        string
	SolvedProblems string
	Score           int
}
type Problems struct {
	Id          int64 `xorm:"autoincr pk 'id'"`
	Name        string `xorm:"unique"`
	Description string //Problem Description
	Value       int //score
	Category    string
	Hidden      bool   //should be problem hide?
	Location    string //saved physical position
	Scale       int    `xorm:"default 0"`//score scale when each answer submit
}

type UserProblems struct {
	Id	int64 `xorm:"autoincr pk 'id'"`
	userId  int64 //foreignkey gctf_user.id
	Location     string //problem net location
	flag        string
	problemsId int64  //foreignkey gctf_problems.id

}


type Hints struct {
	Id int64	`xorm:"autoincr pk 'id'"`
	ProblemsId int64//foreignkey gctf_problems.id
	Hint string
	Cost int //cost score to get hint
}

type Tag struct {
	Id int64	`xorm:"autoincr pk 'id'"`
	problemsId int64//foreignkey gctf_problems.id
	Tag string
}

type Team struct {
	Id int64	`xorm:"autoincr pk 'id'"`
	Name string `xorm:"unique"`
	Member string
	Banned bool //if true this team can't login
}
