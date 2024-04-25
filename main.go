package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gustavonlorenzo/loggerator"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type postgres struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int
	Dbname   string
}

type conn struct {
	Postgres postgres `yaml:"postgres"`
}

type pet struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Species string `json:"species"`
}

func (c *conn) getConf() *conn {

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		l.Error("Failed to read YAML file.")
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		l.Error("Failed to Unmarshal file.")
		panic(err)
	}

	return c
}

func getPets(g *gin.Context) {
	var c conn
	c.getConf()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgres.Host, c.Postgres.Port, c.Postgres.Username, c.Postgres.Password, c.Postgres.Dbname)

	db, err := sql.Open("postgres", psqlInfo)

	l.Info("Establishing connection to PG.")
	if err != nil {
		l.Info("Error establishing connection to pg backend.")
		panic(err)
	}

	defer db.Close()

	var p pet
	var allPet []pet
	var id int

	sqlStatement := "SELECT * FROM pets"
	rows, err := db.Query(sqlStatement)

	if err != nil {
		l.Error("Invalid query.")
	}

	for rows.Next() {
		err = rows.Scan(&id, &p.Name, &p.Age, &p.Species)
		if err != nil {
			fmt.Println("Invalid data returned.")
		}
		allPet = append(allPet, pet{Name: p.Name, Age: p.Age, Species: p.Species})
	}

	g.IndentedJSON(http.StatusOK, allPet)

}

func addRecord(g *gin.Context) {
	var c conn
	c.getConf()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgres.Host, c.Postgres.Port, c.Postgres.Username, c.Postgres.Password, c.Postgres.Dbname)

	db, err := sql.Open("postgres", psqlInfo)

	l.Info("Establishing connection to PG.")
	if err != nil {
		l.Error("Error establishing connection to pg backend.")
		panic(err)
	}

	defer db.Close()

	var pet pet
	if err := g.BindJSON(&pet); err != nil {
		return
	}

	sqlStatement := `
  INSERT INTO pets (name, age, species)
  VALUES ($1, $2, $3)`

	_, err = db.Exec(sqlStatement, pet.Name, pet.Age, pet.Species)
	if err != nil {
		fmt.Println("error with query")
	}

	l.Info("successfully created record.")
	g.IndentedJSON(http.StatusCreated, pet)
}

func delRecord(g *gin.Context) {

	var c conn
	c.getConf()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgres.Host, c.Postgres.Port, c.Postgres.Username, c.Postgres.Password, c.Postgres.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	l.Info("Establishing connection to PG.")
	if err != nil {
		l.Error("Error establishing connection to pg backend.")
		panic(err)
	}

	defer db.Close()

	var pet pet
	if err := g.BindJSON(&pet); err != nil {
		return
	}

	sqlStatement := `
  DELETE FROM pets WHERE name=($1)`

	_, err = db.Exec(sqlStatement, pet.Name)
	if err != nil {
		l.Error("Error deleting pet.")
	}

	g.IndentedJSON(http.StatusOK, pet)
}

var l *zap.Logger = loggerator.InitializeLogger()

func main() {
	l.Info("Initializing router.")
	router := gin.Default()

	router.GET("/pets", getPets)
	router.POST("/pets", addRecord)
	router.POST("/pets/remove", delRecord)

	router.Run("localhost:8080")

}
