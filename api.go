package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os/exec"
)

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func A03Injection(r *gin.Engine) {
	// Command Injection
	command := func(c *gin.Context) {
		host := `host ` + c.Param("host")
		out, err := exec.Command("/bin/bash", "-c", host).Output()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"msg": string(out)})
	}
	// SQL Injection
	//http://127.0.0.1:8080/injection/sql/3%20or%20password%20=%20'admin123'
	//http://127.0.0.1:8080/injection/sql/3;DROP%20TABLE%20users
	SQL := func(c *gin.Context) {
		db, err := sql.Open("sqlite3", "owasp.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		id := c.Param("id")
		query := fmt.Sprintf("SELECT * FROM users WHERE id = %s", id)
		var u struct {
			id       int
			user     string
			password string
		}
		err = db.QueryRow(query).Scan(&u.id, &u.user, &u.password)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"username": u.user})
	}

	injection := r.Group("/injection")
	{
		injection.GET("/command/:host", command)
		injection.GET("/sql/:id", SQL)
	}
}

func apiServer() {
	router := gin.Default()

	A03Injection(router)

	router.Run("127.0.0.1:8080")
}

func main() {
	apiServer()
}
