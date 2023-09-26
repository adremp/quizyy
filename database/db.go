package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
)

type Quiz struct {
	Id       string
	Question string
	Answer   string
	Variants pg.StringArray
}

type Instance struct {
	*sqlx.DB
}

func New() (*Instance, error) {
	connection := fmt.Sprintf("user=%s database=%s password=%s sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"))
	sqlx, err := sqlx.Open("postgres", connection)
	if err != nil {
		log.Printf("sqlx connect error: %v", err)
		return nil, err
	}
	return &Instance{sqlx}, nil
}

func (db *Instance) CreateQuiz(quiz Quiz) error {
	_, err := db.NamedExec("INSERT INTO quiz (question, answer, variants) VALUES (:question, :answer, :variants)", quiz)
	if err != nil {
		log.Printf("create quiz error: %v", err)
		return err
	}
	return nil
}

func (db *Instance) GetQuiz(quizId string) (Quiz, error) {
	var quiz []Quiz
	quizIdInt, err := strconv.Atoi(quizId)
	if err != nil {
		log.Printf("parse quizId error: %v", err)
		return Quiz{}, err
	}
	err = db.Select(&quiz, "SELECT * FROM quiz WHERE id = $1 LIMIT 1", quizIdInt)
	if err != nil {
		log.Printf("get quiz error: %v", err)
		return Quiz{}, err
	}
	return quiz[0], nil
}

func (db *Instance) GetQuizes() ([]Quiz, error) {
	var quizes []Quiz
	err := db.Select(&quizes, "SELECT * FROM quiz")
	if err != nil {
		log.Printf("get quizes error: %v", err)
		return nil, err
	}
	return quizes, nil
}
