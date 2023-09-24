package render

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RenderTemplate(c *gin.Context, name string, data interface{}) {
	// Parse the template
	tmpl, err := template.ParseFiles("./templates/" + name + ".html")
	if err != nil {
		log.Println(name+" template Parse Files error:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return // Exit the function to avoid further execution
	}

	// Execute the template with the provided data
	if err := tmpl.Execute(c.Writer, data); err != nil {
		log.Println(name+" template Execute Template error:", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return // Exit the function to avoid further execution
	}

	log.Println(name + " page rendered successfully!")

}
