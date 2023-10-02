package props

import (
	"html/template"
	db "quizyy/database"
)

type Success struct {
	Title string
	ActionText string
	ActionUrl string
}

type Input struct {
	Label string
	Attrs template.HTMLAttr
	ErrorText string
}
type VariantsInput struct {
	Variants []string
	ErrorText string
}

type CreateQuizForm struct {
	Inputs []Input
	VariantsInput VariantsInput
}

type Quizes struct {
	Quizes []db.Quiz
}
