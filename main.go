package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	db "quizyy/database"
	props "quizyy/templates"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/static", "./static")

	t, err := template.New("").Funcs(template.FuncMap{"join": strings.Join, "lowercase": strings.ToLower, "omitEmptyStrings": omitEmptyStrings}).ParseGlob("templates/base/*.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[TEMPLATES] %s", t.DefinedTemplates())

	tHome, err := CopyAndParse(t, "templates/home.html")
	if err != nil {
		log.Fatal(err)
	}
	tQuizes, err := CopyAndParse(t, "templates/quizes.html")
	if err != nil {
		log.Fatal(err)
	}
	tCreate, err := CopyAndParse(t, "templates/create-quiz.html")
	if err != nil {
		log.Fatal(err)
	}
	tQuiz, err := CopyAndParse(t, "templates/quiz.html")
	if err != nil {
		log.Fatal(err)
	}

	sqlxx, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		tHome.ExecuteTemplate(c.Writer, "index.html", nil)
	})

	r.GET("/quizes", func(c *gin.Context) {
		quizes, err := sqlxx.GetQuizes()
		if err != nil {
			fmt.Printf("[ERROR] get quizes error %s", err)
		}

		data := props.Quizes{quizes}
		if c.GetHeader("HX-Request") != "" {
			err := tQuizes.ExecuteTemplate(c.Writer, "body", data)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			tQuizes.ExecuteTemplate(c.Writer, "index.html", data)
		}
	})

	r.GET("/answer", func(c *gin.Context) {
		id := c.Query("id")
		answer := c.Query("answer")
		fmt.Print("[INFO]", id, answer)
		if id == "" || answer == "" {
			fmt.Print("[ERROR] invalid params")
			c.JSON(400, gin.H{"error": "invalid params"})
			return
		}
		quiz, err := sqlxx.GetQuiz(id)
		if err != nil {
			fmt.Printf("[ERROR] get quiz error %s", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if quiz.Answer == answer {
			fmt.Print("[INFO] correct")
			c.Header("HX-Retarget", "main")
			if err != nil {
				fmt.Printf("[ERROR] parse id error %s", err)
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			tQuiz.ExecuteTemplate(c.Writer, "success.html", props.Success{Title: "Correct", ActionText: "Next", ActionUrl: "/quizes"})
		} else {
			fmt.Print("[INFO] wrong")
			err := tQuiz.ExecuteTemplate(c.Writer, "answer-button-wrong", answer)
			if err != nil {
				log.Print("[ERROR]", err)
				log.Fatal(err)
			}
		}
	})

	r.GET("/validate", func(c *gin.Context) {
		var quiz *db.QuizInput
		if c.ShouldBindQuery(&quiz) != nil {
			fmt.Print("[ERROR] invalid params")
			c.JSON(400, gin.H{"error": "invalid params"})
			return
		}

		err := db.Validate.Struct(quiz)

		for qKey := range c.Request.URL.Query() {
			c.Header("HX-Retarget", fmt.Sprintf("#%s-error", qKey))
			for _, err := range err.(validator.ValidationErrors) {
				fieldLower := strings.ToLower(err.Field())
				fmt.Printf("[INFO] validate, fieldLower: %s, tag: %s, qKey: %s", fieldLower, err.Tag(), qKey)
				if fieldLower != qKey {
					continue
				}
				tCreate.ExecuteTemplate(c.Writer, "text-error", db.QuizErrorMessages[fieldLower][err.Tag()])
			}
			tCreate.ExecuteTemplate(c.Writer, "nil", nil)
			break
		}
	})

	r.GET("/create", func(c *gin.Context) {
		data := props.CreateQuizForm{Inputs: []props.Input{{"Question", `name="question"`, ""}, {"Answer", `name="answer"`, ""}}, VariantsInput: props.VariantsInput{[]string{}, ""}}
		c.Header("Cache-Control", "public, max-age:60")
		if c.GetHeader("HX-Request") != "" {
			err := tCreate.ExecuteTemplate(c.Writer, "body", data)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = tCreate.ExecuteTemplate(c.Writer, "index.html", data)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	r.GET("/quizes/:id", func(c *gin.Context) {
		quiz, err := sqlxx.GetQuiz(c.Param("id"))
		if err != nil {
			fmt.Printf("[ERROR] get quiz error %s", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if c.GetHeader("HX-Request") != "" {
			err := tQuiz.ExecuteTemplate(c.Writer, "body", quiz)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			tQuiz.ExecuteTemplate(c.Writer, "index.html", quiz)
		}
	})

	r.POST("/quiz", func(c *gin.Context) {
		var quiz db.QuizInput
		err := c.Bind(&quiz)
		if err != nil {
			fmt.Printf("[ERROR] bind quiz error %s", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		errs := db.Validate.Struct(&quiz)
		errArr, ok := errs.(validator.ValidationErrors)
		if ok {
			fmt.Print("[ERRORS EXIST]")
			var quizErrors = map[string]string{}
			for _, err := range errArr {
				fieldLower := strings.ToLower(err.Field())
				errMess := db.QuizErrorMessages[fieldLower][err.Tag()]
				quizErrors[fieldLower] = errMess
			}
			c.Header("HX-Retarget", "main")
			tCreate.ExecuteTemplate(c.Writer, "create-quiz-form.html", props.CreateQuizForm{Inputs: []props.Input{{"Question", toHTMLAttr(map[string]string{"name": "question", "value": quiz.Question}), quizErrors["question"]}, {"Answer", toHTMLAttr(map[string]string{"name": "answer", "value": quiz.Answer}), quizErrors["answer"]}}, VariantsInput: props.VariantsInput{strings.Split(quiz.Variants, ","), quizErrors["variants"]}})
			return
		}
		formatedQuiz := db.BindQuizForm(c)
		if err := sqlxx.CreateQuiz(formatedQuiz); err != nil {
			fmt.Errorf("[ERROR] create quiz error %s", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err = tCreate.ExecuteTemplate(c.Writer, "success.html", props.Success{Title: "Well done!", ActionText: "Go to home", ActionUrl: "/"})
		if err != nil {
			log.Print("[ERROR]", err)
			log.Fatal(err)
		}
	})

	r.PATCH("/variants-input", func(c *gin.Context) {
		var combined []string
		if variants := c.PostForm("variants"); variants != "" {
			varinatsArr := strings.Split(variants, ",")
			combined = append(varinatsArr, c.PostForm("variant"))
		} else {
			combined = []string{c.PostForm("variant")}
		}
		tCreate.ExecuteTemplate(c.Writer, "variants-input-list", combined)
	})

	err = r.Run(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func CopyAndParse(temp *template.Template, files ...string) (*template.Template, error) {
	t, err := temp.Clone()
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		t, err = t.ParseFiles(file)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func omitEmptyStrings(arr []string) []string {
	var result []string
	for _, s := range arr {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func toHTMLAttr(attrs map[string]string) template.HTMLAttr {
	var result string
	for k, v := range attrs {
		result += k + `="` + v + `"`
	}
	return template.HTMLAttr(result)
}
