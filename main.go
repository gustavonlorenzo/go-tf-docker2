package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gustavonlorenzo/loggerator"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
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

	yamlFile, err := os.ReadFile("./infra/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Failed to Unmarshal file: %v", err)
	}

	return c
}

func getPets() []pet {
	var c conn
	c.getConf()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgres.Host, c.Postgres.Port, c.Postgres.Username, c.Postgres.Password, c.Postgres.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error establishing connection to pg backend.")
	}

	defer db.Close()

	var p pet
	var allPet []pet
	var id int

	sqlStatement := "SELECT * FROM pets"
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Invalid query statement")
	}

	for rows.Next() {
		err = rows.Scan(&id, &p.Name, &p.Age, &p.Species)
		if err != nil {
			log.Fatalf("Invalid data returned.")
		}
		allPet = append(allPet, pet{Name: p.Name, Age: p.Age, Species: p.Species})
	}

	return allPet

}

func addRecord(g *gin.Context) {
	var c conn
	c.getConf()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgres.Host, c.Postgres.Port, c.Postgres.Username, c.Postgres.Password, c.Postgres.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error")
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
		log.Fatalf("error")
	}

	g.IndentedJSON(http.StatusCreated, pet)
}

func delRecord(g *gin.Context) {

	var c conn
	c.getConf()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", c.Postgres.Host, c.Postgres.Port, c.Postgres.Username, c.Postgres.Password, c.Postgres.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error")
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
		log.Fatalf("error")
	}

	g.IndentedJSON(http.StatusOK, pet)
	// http.StatusAccepted
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("layout.html")

	router.GET("/pets", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"Title": "PETS",
			"Items": getPets(),
		})
	})
	router.POST("/pets", addRecord)
	router.POST("/pets/remove", delRecord)

	router.Run("localhost:8080")

}
