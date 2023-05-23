package database

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/catinapoke/telegram-file-bot/internal/common"
	_ "github.com/lib/pq"
)

type DatabaseOperator struct {
	db *sql.DB
}

func (op *DatabaseOperator) Start() error {
	// Enviroment varaiables
	addr, err := common.GetEnv("DATABASE_URL")
	if err != nil {
		return common.NewDatabaseError("can't get address of database!", err)
	}

	user, err := common.GetEnvFromFile("POSTGRES_USER_FILE")
	if err != nil {
		return common.NewDatabaseError("can't get useruser of database!", err)
	}

	password, err := common.GetEnvFromFile("POSTGRES_PASSWORD_FILE")
	if err != nil {
		return common.NewDatabaseError("can't get password of database!", err)
	}

	db, err := common.GetEnvFromFile("POSTGRES_DB_FILE")
	if err != nil {
		return common.NewDatabaseError("can't get database name!", err)
	}

	// Connect
	url := fmt.Sprintf("postgres://%s:%s@%s/%v?sslmode=disable", user, password, addr, db)
	op.db, err = sql.Open("postgres", url)

	if err != nil {
		return common.NewDatabaseError("Connect error", err)
	}

	if err = op.db.Ping(); err != nil {
		return common.NewDatabaseError("Ping error", err)
	}

	return nil
}

func (op *DatabaseOperator) Close() {
	op.db.Close()
}

func (op *DatabaseOperator) CreateUser(user User) error {
	if op.db == nil {
		return errNotInitialized()
	}

	err := op.insertUser(user)

	if err != nil {
		err = common.NewDatabaseError("Create user", err)
	}

	return err
}

func (op *DatabaseOperator) GetFiles(user int64) ([]FileName, error) {
	if op.db == nil {
		return nil, errNotInitialized()
	}

	var files []FileName
	var err error = nil // err := op.db.Model(&files).Select()

	if err != nil {
		err = common.NewDatabaseError("GetFiles", err)
	}

	return files, err
}

func (op *DatabaseOperator) DeleteFile(name string) error {
	if op.db == nil {
		return errNotInitialized()
	}
	/*
		file := FileName{
			Name: name,
		}
	*/

	var err error = nil //_, err := op.db.Model(file).Delete()

	if err != nil {
		err = common.NewDatabaseError("DeleteFile", err)
	}

	return err
}

func (op *DatabaseOperator) CreateFile(file FileName) error {
	if op.db == nil {
		return errNotInitialized()
	}

	var err error = op.insertFile(file)

	if err != nil {
		err = common.NewDatabaseError("CreateFile", err)
	}

	return err
}

func (op *DatabaseOperator) fillUsers(rows *sql.Rows) ([]User, error) {
	var current User
	set := make([]User, 8)

	for rows.Next() {
		err := rows.Scan(&current.Id, &current.Username, &current.StartUsage, &current.Key, &current.FirstName, &current.LastName, &current.LanguageCode)
		if err != nil {
			return nil, err
		}
		log.Println(current)
		set = append(set, current)
	}
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return set, nil
}

func (op *DatabaseOperator) insertUser(user User) error {
	stmt, err := op.db.Prepare(`INSERT INTO telegramfilebot.users (id, username, start_usage, key, firstname, lastname, language_code) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return fmt.Errorf("insertUser: Prepare: %w", err)
	}

	res, err := stmt.Exec(user.Id, user.Username, user.StartUsage, user.Key, user.FirstName, user.LastName, user.LanguageCode)
	if err != nil {
		return fmt.Errorf("insertUser: Exec: %w", err)
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("insertUser: affected users: %w", err)
	}
	log.Printf("affected = %d\n", rowCnt)
	return nil
}

func (op *DatabaseOperator) insertFile(file FileName) error {
	stmt, err := op.db.Prepare(`INSERT INTO telegramfilebot.filenames (name, id, share, owner) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("insertUser: Prepare: %w", err)
	}

	_, err = stmt.Exec(file.Name, file.Id, file.Share, file.Owner)
	if err != nil {
		return fmt.Errorf("insertUser: Exec: %w", err)
	}

	return nil
}

func errNotInitialized() error {
	return common.NewDatabaseError("Database isn't initialized!", nil)
}
