package database

import "time"

type User struct {
	Id           int64
	Username     string
	StartUsage   time.Time
	Key          string
	FirstName    string
	LastName     string
	LanguageCode string
}

type FileName struct {
	Name  string
	Id    int
	Share int8
	Owner int64
}

type FileShare struct {
	Name   string
	UserId int64
}
