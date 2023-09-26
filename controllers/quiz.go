package quiz

type Quiz struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Variants []string `json:"variants"`
}

func Create()

func Get()