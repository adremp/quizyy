package props

import db "quizyy/database"

type Success struct {
	Title string
	ActionText string
	ActionUrl string
}

type Index struct {
	Title string
}

type Input struct {
	Label string
	Name  string
}

type CreateQuizForm struct {
	Inputs []Input
	Variants []string
}

type Quizes struct {
	Quizes []db.Quiz
}
// type QuizListMain struct {
// 	Quizes []db.Quiz
// }

type QuizMain struct {
	db.Quiz
}
