package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
)

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func dbFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}

		rows, err := db.Query("SELECT tick FROM ticks")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		}
	}
}

// dbGetRouteFunc gets all data forwarded from a Blues Wireless route from the
// example database
func dbGetRouteFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Create the routedata table if it doesn't yet exist
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS routedata ( id serial primary key not null, routedata json not null );"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		rows, err := db.Query("SELECT routedata FROM routedata;")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading routedata: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var routedata string
			if err := rows.Scan(&routedata); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning routedata: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", routedata))
		}
	}
}

// dbPostRouteFunc allows a Blues Wireless route configured to this application
// to POST data to the example database
func dbPostRouteFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create the routedata table if it doesn't yet exist
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS routedata ( id serial primary key not null, routedata json not null );"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		//get the body as a string
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		newStr := buf.String()

		query := fmt.Sprintf("insert into routedata ( routedata ) values ( '%s' );", newStr)

		_, err := db.Exec(query)
		if err != nil {
			log.Printf("----- error inserting: %s", err)
			return
		}

		c.String(http.StatusOK, "accepted")
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.Run([]byte("**hi!**"))))
	})

	router.GET("/repeat", repeatHandler(repeat))

	router.GET("/db", dbFunc(db))

	//Blues Wireless route example methods
	router.GET("/route", dbGetRouteFunc(db))
	router.POST("/route", dbPostRouteFunc(db))

	router.Run(":" + port)
}
