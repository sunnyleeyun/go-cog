package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Model struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UseCase     string `json:"use_case"`
	DockerTag   string `json:"docker_tag"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL")+"/CogGpt")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	router := gin.Default()

	router.GET("/models", func(context *gin.Context) {
		var (
			model  Model
			models []Model
		)
		rows, err := db.Query("select name, description, use_case, docker_tag from Model;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&model.Name, &model.Description, &model.UseCase, &model.DockerTag)
			models = append(models, model)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		context.JSON(http.StatusOK, gin.H{
			"result": models,
			"count":  len(models),
		})
	})

	router.Run("localhost:9090")
}
