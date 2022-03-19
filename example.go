package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BIND QUERY STRING OR POST DATA
type Person struct {
	Name     string    `form:"name"`
	Addess   string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

//Bind HTML CHECKBOXES
type myForm struct {
	Colors []string `form:"colors[]"`
}

type StructA struct {
	FieldA string `form:"field_a"`
}

type StructB struct {
	NestedStruct StructA
	FiledB       string `form:"field_b"`
}

type StructC struct {
	NestedStructPointer *StructA
	FiledC              string `form:"field_c"`
}

type StructD struct {
	NestedAnonyStruct struct {
		FiledX string `form:"field_x"`
	}
	FiledD string `form:"field_x"`
}

type PersonUri struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func GetDataB(ctx *gin.Context) {
	var b StructB
	ctx.Bind(&b)
	ctx.JSON(http.StatusOK, gin.H{
		"a": b.NestedStruct,
		"b": b.FiledB,
	})
}

func GetDataC(ctx *gin.Context) {
	var c StructC
	ctx.Bind(&c)
	ctx.JSON(http.StatusOK, gin.H{
		"a": c.NestedStructPointer,
		"c": c.FiledC,
	})
}

func GetDataD(ctx *gin.Context) {
	var d StructD
	ctx.Bind(&d)
	ctx.JSON(http.StatusOK, gin.H{
		"x": d.NestedAnonyStruct,
		"d": d.FiledD,
	})
}

func indexHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "form.html", nil)
}

func formHnadler(ctx *gin.Context) {
	var fakeForm myForm
	ctx.Bind(&fakeForm)
	ctx.JSON(http.StatusOK, gin.H{
		"color": fakeForm.Colors,
	})
}

func getperson(ctx *gin.Context) {
	var person Person
	if ctx.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Addess)
		log.Println(person.Birthday)
	}
	ctx.String(http.StatusOK, "Success")
}

func main() {

	// gin.DisableConsoleColor()
	gin.ForceConsoleColor()
	r := gin.Default()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// r := gin.New()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ping",
		})
	})

	r.GET("/:name/:id", func(ctx *gin.Context) {
		var personUri PersonUri
		if err := ctx.ShouldBindUri(&personUri); err != nil {
			ctx.JSON(400, gin.H{"msg": err})
			return
		}
		ctx.JSON(200, gin.H{"name": personUri.Name, "uuid": personUri.ID})
	})

	// using AsciiJSON
	r.GET("/someJson", func(ctx *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>"}
		ctx.AsciiJSON(http.StatusOK, data)
	})

	//bind form-data request with custom struct
	r.GET("/getb", GetDataB)
	r.GET("/getc", GetDataC)
	r.GET("/getd", GetDataD)

	//BIND A CHECKBOX GROUP
	r.LoadHTMLGlob("views/*")
	r.GET("/", indexHandler)
	r.POST("/", formHnadler)

	//BIND QUERY STRING OR POST DATA
	r.GET("/getperson", getperson)

	// t, err := loadTemplate()

	// if err != nil {
	// 	panic(err)
	// }

	// r.SetHTMLTemplate(t)

	// r.GET("/", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "/html/index.tmpl", nil)
	// })

	// r.Run("localhost:8080")
	s.ListenAndServe()
}

// func loadTemplate() (*template.Template, error) {
// 	t := template.New("")
// 	for name, file := range Assets.Files {
// 		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
// 			continue
// 		}

// 		h, err := ioutil.ReadAll(file)
// 		if err != nil {
// 			return nil, err
// 		}

// 		t, err = t.New(name).Parse(string(h))

// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return t, nil
// }
