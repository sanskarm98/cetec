package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	NAME        string `json:"name"`
	PHONENUMBER string `json:"phone_number"`
	CITY        string `json:"city"`
	STATE       string `json:"state"`
	STREET1     string `json:"street1"`
	STREET2     string `json:"street2"`
	ZIPCODE     string `json:"zip_code"`
}

func main() {
	router := gin.Default()
	router.GET("/person/:person_id/info", getByPersonId)
	router.POST("/person/create", CreatePerson)

	router.Run("localhost:8080")
}
func connectdatabase() *sql.DB {
	db, err := sql.Open("mysql", "username:password@localhost/cetec")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func getByPersonId(c *gin.Context) {
	id := c.Param("person_id")
	db := connectdatabase()
	defer db.Close()

	sqlquery := "SELECT  pe.name,ph.number,add.city,adr.state,adr.street1,adr.street2,adr.zip_code FROM address_join AS aj INNER JOIN address AS adr ON adr.id = aj.address_id INNER JOIN person AS pe ON aj.person_id = pe.id INNER JOIN phone AS ph ON ph.person_id = aj.person_id where aj.person_id = ?;"
	res, err := db.Query(sqlquery, id)
	defer res.Close()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	for res.Next() {

		var person Person
		err := res.Scan(&person.NAME, &person.PHONENUMBER, &person.CITY, &person.STATE, &person.STREET1, &person.STREET2, &person.ZIPCODE)

		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"data": person})

	}
}

func CreatePerson(c *gin.Context) {
	var newperson Person
	var createdperson Person = Person{
		NAME:        "Sanskar",
		PHONENUMBER: "123-456-7890",
		CITY:        "Sacramento",
		STATE:       "CA",
		STREET1:     "112 Main St",
		STREET2:     "Apt 12",
		ZIPCODE:     "12345",
	}
	if err := c.BindJSON(&newperson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := connectdatabase()
	defer db.Close()

	sqlquery := "INSERT INTO person(name) values (?);"
	res, err := db.Exec(sqlquery, createdperson.NAME)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert in person  table failed"})
		return
	}
	personeid, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	sqlquery = "INSERT INTO phone( person_id,number) values (?,?);"
	_, err = db.Exec(sqlquery, personeid, createdperson.PHONENUMBER)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert in phone table failed"})
		return
	}

	sqlquery = "INSERT INTO address(city , state , street1 , street2 ,zip_code) values (?,?,?,?,?);"
	res, err = db.Exec(sqlquery, createdperson.CITY, createdperson.STATE, createdperson.STREET1, createdperson.STREET2, createdperson.ZIPCODE)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert in address table failed"})
		return
	}
	addressid, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	sqlquery = "INSERT INTO address_join( person_id,address_id) values (?,?);"
	_, err = db.Exec(sqlquery, personeid, addressid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert in address_join table failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": createdperson})

}
