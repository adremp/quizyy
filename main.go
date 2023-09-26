package main

import (
	db "counter/database"
	props "counter/templates"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	r := gin.Default()
	r.Static("/static", "./static")

	log.SetOutput(os.Stdout)
	gin.SetMode(gin.ReleaseMode)

	t, err := template.New("").Funcs(template.FuncMap{"join": strings.Join}).ParseGlob("templates/base/*.html")
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
			// idInt, err := strconv.Atoi(id)
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

	r.GET("/create", func(c *gin.Context) {
		data := props.CreateQuizForm{Inputs: []props.Input{{"Question", "question"}, {"Answer", "answer"}}, Variants: []string{}}
		c.Header("Cache-Control", "public, max-age:60")
		if c.GetHeader("HX-Request") != "" {
			err := tCreate.ExecuteTemplate(c.Writer, "body", data)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			tCreate.ExecuteTemplate(c.Writer, "index.html", data)
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
		quiz := db.Quiz{Question: c.PostForm("question"), Answer: c.PostForm("answer"), Variants: strings.Split(c.PostForm("variants"), ",")}
		if err := sqlxx.CreateQuiz(quiz); err != nil {
			fmt.Errorf("[ERROR] create quiz error %s", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		t.ExecuteTemplate(c.Writer, "success.html", props.Success{Title: "Well done!", ActionText: "Go to home", ActionUrl: "/"})
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

	// dataFunc := map[string]func(c *gin.Context) any{
	// 	"create-quiz-form": func(c *gin.Context) any {
	// 		return props.CreateQuizForm{
	// 			Inputs:   []props.Input{{"Question", "question"}, {"Answer", "answer"}},
	// 			Variants: []string{},
	// 		}
	// 	},
	// 	"quiz-list-main": func(c *gin.Context) interface{} {
	// 		data, err := sqlxx.GetQuizes()
	// 		if err != nil {
	// 			c.JSON(400, gin.H{"error": err.Error()})
	// 			return nil
	// 		}
	// 		return props.QuizListMain{Quizes: data}
	// 	},
	// 	"quiz-main": func(c *gin.Context) any {
	// 		data, err := sqlxx.GetQuiz(c.Query("id"))
	// 		if err != nil {
	// 			c.JSON(400, gin.H{"error": err.Error()})
	// 			return nil
	// 		}
	// 		return props.QuizMain{data}
	// 	},
	// 	"variants-item": func(c *gin.Context) any {
	// 		fmt.Print(c.Param("variant"), c.PostForm("variant"), )
	// 		return c.Query("variant")
	// 	},
	// }

	// r.GET("/t/:templ", func(c *gin.Context) {
	// 	c.HTML(200, c.Param("templ")+".html", dataFunc[c.Param("templ")](c))
	// })

	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "index.html", struct{ Text, Title string }{"Home", "Home"})
	// })

	// r.POST("/post-quiz", func(c *gin.Context) {
	// 	quiz := db.Quiz{Question: c.PostForm("question"), Answer: c.PostForm("answer")}
	// 	if err := sqlxx.CreateQuiz(quiz); err != nil {
	// 		c.JSON(400, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.HTML(200, "post-quiz-success.html", nil)
	// })

	r.Run()
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
