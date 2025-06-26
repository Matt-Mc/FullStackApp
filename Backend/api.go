package main

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunAPI() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/api/bills", func(c *gin.Context) {
		var bills []Bill
		rows, err := appDB.Query("SELECT Session, Legisinfo_id, Introduced, Name, Number, Url FROM bills")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var bill Bill
			var nameEN string // Temporary variable for the English name
			err = rows.Scan(&bill.Session, &bill.Legisinfo_id, &bill.Introduced, &nameEN, &bill.Number, &bill.Url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			bill.Name.EN = nameEN // Assign the scanned string to the EN field
			bills = append(bills, bill)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bills": bills})
	})

	router.GET("/api/bills/:id", func(c *gin.Context) {
		id := c.Param("id")
		var bill Bill
		var nameEN string // Temporary variable for the English name
		row := appDB.QueryRow("SELECT Session, Legisinfo_id, Introduced, Name, Number, Url FROM bills WHERE Legisinfo_id = $1", id)
		err := row.Scan(&bill.Session, &bill.Legisinfo_id, &bill.Introduced, &nameEN, &bill.Number, &bill.Url)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		bill.Name.EN = nameEN // Assign the scanned string to the EN field
		c.JSON(http.StatusOK, gin.H{"bill": bill})
	})

	router.GET("/api/mps", func(c *gin.Context) {
		var mps []MP
		rows, err := appDB.Query("SELECT * FROM mps")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var mp MP
			err = rows.Scan(&mp.Name, &mp.CurrentParty.ShortName.EN, &mp.CurrentRiding.Name.EN, &mp.URL, &mp.Image, &mp.CurrentRiding.Province)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			mps = append(mps, mp)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"MPs": mps})
	})

	router.GET("/api/mps/:name", func(c *gin.Context) {
		name := c.Param("name")
		name, _ = url.PathUnescape(name)
		var mp MP
		row := appDB.QueryRow("SELECT * FROM mps WHERE name = $1", name)
		err := row.Scan(&mp.Name, &mp.CurrentParty.ShortName.EN, &mp.CurrentRiding.Name.EN, &mp.URL, &mp.Image, &mp.CurrentRiding.Province)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"MP": mp})
	})

	router.Run(":1500")
}
