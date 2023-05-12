package main

import (
	"time"

	"github.com/go-pg/pg"
)

type DatabaseOperator struct {
	db *pg.DB
}

type Users struct {
	Id           int64
	Username     string
	StartUsage   time.Time
	Key          string
	FirstName    string
	LastName     string
	LanguageCode string
}

type Filenames struct {
	Name  string
	Id    int
	Share int8
	Owner int64
}

type FileShare struct {
	Name   string
	UserId int64
}

func (op *DatabaseOperator) Start() error {
	addr, err := GetEnv("DATABASE_URL")
	if err != nil {
		return NewDatabaseError("can't get address of database!", err)
	}

	user, err := GetEnvFromFile("POSTGRES_USER_FILE")
	if err != nil {
		return NewDatabaseError("can't get useruser of database!", err)
	}

	password, err := GetEnvFromFile("POSTGRES_PASSWORD_FILE")
	if err != nil {
		return NewDatabaseError("can't get password of database!", err)
	}

	db, err := GetEnvFromFile("POSTGRES_DB_FILE")
	if err != nil {
		return NewDatabaseError("can't get database name!", err)
	}

	op.db = pg.Connect(&pg.Options{
		Addr:     addr,
		User:     user,
		Password: password,
		Database: db,
	})
}

func (op *DatabaseOperator) Close() {
	op.db.Close()
}

func (op *DatabaseOperator) CreateUser(user *Users) error {
	if op.db == nil {
		return errNotInitialized()
	}

	_, err := op.db.Model(user).Insert()

	if err != nil {
		err = NewDatabaseError("Create user", err)
	}

	return err
}

func (op *DatabaseOperator) GetFiles(user int64) ([]Filenames, error) {
	if op.db == nil {
		return nil, errNotInitialized()
	}

	var files []Filenames
	err := op.db.Model(&files).Select()

	if err != nil {
		err = NewDatabaseError("GetFiles", err)
	}

	return files, nil
}

func (op *DatabaseOperator) DeleteFile(name string) error {
	if op.db == nil {
		return errNotInitialized()
	}

	file := Filenames{
		Name: name,
	}

	_, err := op.db.Model(file).Delete()

	if err != nil {
		err = NewDatabaseError("DeleteFile", err)
	}

	return err
}

func (op *DatabaseOperator) CreateFile(file Filenames) error {
	if op.db == nil {
		return errNotInitialized()
	}

	_, err := op.db.Model(file).Insert()

	if err != nil {
		err = NewDatabaseError("CreateFile", err)
	}

	return err
}

func errNotInitialized() error {
	return NewDatabaseError("Database isn't initialized!", nil)
}
