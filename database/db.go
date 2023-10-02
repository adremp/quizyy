package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
)

type Quiz struct {
	*QuizInputFormated
	Id string
}
type QuizInput struct {
	Question string `json:"question" validate:"required,endswith=?" form:"question"`
	Answer   string `json:"answer" validate:"required" form:"answer"`
	Variants string `json:"variants" validate:"fieldcontains=Answer,comajoinedmin=3,comajoinedunic" form:"variants"`
}

type QuizInputFormated struct {
	Question string         `json:"question" form:"question"`
	Answer   string         `json:"answer" form:"answer"`
	Variants pg.StringArray `json:"variants" form:"variants"`
}

type Instance struct {
	*sqlx.DB
}

var (
	Validate          *validator.Validate
	QuizErrorMessages = map[string]map[string]string{
		"question": {
			"required": "Question is required",
			"endswith": `Question must end with "?"`,
		},
		"answer": {
			"required": "Answer is required",
		},
		"variants": {
			"comajoinedmin": "At least 3 variants are required",
			"fieldcontains": "Must contain answer",
			"comajoinedunic": "All variants must be unique",
		},
	}
)

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	Validate.RegisterValidation("comajoinedmin", func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Field().String(), ",")
		lenParam, err := strconv.Atoi(fl.Param())
		if err != nil {
			return false
		}
		return len(arr) >= lenParam
	})
	Validate.RegisterValidation("comajoinedunic", func(fl validator.FieldLevel) bool {
		arr := strings.Split(fl.Field().String(), ",")
		seen := make(map[string]bool, len(arr))

		for _, el := range arr {
			if seen[el] {
				return false
			}
			seen[el] = true
		}
		return true
	})
}

func New() (*Instance, error) {
	connection := fmt.Sprintf("user=%s database=%s password=%s host=%s port=%s sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"))
	sqlx, err := sqlx.Open("postgres", connection)
	if err != nil {
		log.Printf("sqlx connect error: %v", err)
		return nil, err
	}
	return &Instance{sqlx}, nil
}

func (db *Instance) CreateQuiz(quiz QuizInputFormated) error {
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

func BindQuizForm(c *gin.Context) QuizInputFormated {
	variantsArr := strings.Split(c.PostForm("variants"), ",")
	fmt.Printf("[variantsArr]: %v", variantsArr)
	return QuizInputFormated{Question: c.PostForm("question"), Answer: c.PostForm("answer"), Variants: pg.StringArray(variantsArr)}
}
