package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunAPI() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/api/bills", func(c *gin.Context) {
		var bills []Bill
		rows, err := appDB.Query("SELECT * FROM bills")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for rows.Next() {
			var bill Bill
			err = rows.Scan(&bill.Name.EN, &bill.Session, &bill.Introduced, &bill.Legisinfo_id, &bill.Number, &bill.Url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rows.Close()
			bills = append(bills, bill)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bills": bills})
	})

	router.GET("/api/bills/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	router.GET("/api/mps", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	router.GET("/api/mps/:name", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	router.Run(":1500")
}
