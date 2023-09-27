package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

var Receipts []Receipt

// item represents data of an item
type Item struct {
	ShortDescription string `json:"shortDescription" validate:"regexp=^[\\w\\s\\-]+$"`
	Price            string `json:"price" validate:"regexp=^\\d+\\.\\d{2}$"`
}

// receipt represents data about a receipt
type Receipt struct {
	id           string `json:"_id"`
	Retailer     string `json:"retailer" binding:"required" validate:"regexp=^\S+$"`
	PurchaseDate string `json:"purchaseDate" binding:"required" validate:"regexp=^20[0-9]{2}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$"`
	PurchaseTime string `json:"purchaseTime" binding:"required" validate:"regexp=^([0-1][0-9]|2[0-3]):[0-5][0-9]$"`
	Items        []Item `json:"items" binding:"required"`
	Total        string `json:"total" binding:"required" validate:"regexp=^\\d+\\.\\d{2}$"`
}

// postReceipt add a receipt from JSON received in the request body
func postReceipt(c *gin.Context) {
	var receipt Receipt
	id := uuid.New().String()
	if err := c.BindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "The receipt is invalid",
		})
		return
	}

	receipt.id = id

	if err := validator.Validate(receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "The receipt is invalid",
		})
		return
	}

	Receipts = append(Receipts, receipt)
	c.JSON(http.StatusOK, gin.H{
		"id": string(receipt.id),
	})
}

func getPoints(c *gin.Context) {
	id := c.Param("id")

	points := 0
	mod := 0.25

	for _, a := range Receipts {
		if a.id == id {

			points += len(strip(a.Retailer))

			r, _ := regexp.Compile("^\\d+\\.0{2}$")
			if r.MatchString(a.Total) {
				points += 50
			}

			f, err := strconv.ParseFloat(a.Total, 64)
			if err != nil {
				c.JSON(http.StatusConflict, gin.H{
					"description": "Inconsistency",
				})
				return
			}

			if math.Mod(f, mod) == 0 {
				points += 25
			}

			points += (len(a.Items) / 2) * 5

			odd, _ := regexp.Compile("[13579]$")
			if odd.MatchString(a.PurchaseDate) {
				points += 6
			}

			time, _ := regexp.Compile("14:(0[1-9]|[1-5][0-9])|15:[0-5][0-9]")
			if time.MatchString(a.PurchaseTime) {
				points += 10
			}

			for _, b := range a.Items {

				des := len(strings.TrimSpace(b.ShortDescription))

				p, err := strconv.ParseFloat(b.Price, 64)
				if err != nil {
					c.JSON(http.StatusConflict, gin.H{
						"description": "Inconsistency",
					})
					return
				}

				if math.Mod(float64(des), 3) == 0 {
					points += int(math.Ceil(0.2 * p))
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"points": points,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"description": "No receipt found for that id",
	})

}

// strio function that removes non
// alfanumeric chars from a string
func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func main() {

	router := gin.Default()

	router.POST("/receipts/process", postReceipt)

	router.GET("/receipts/:id/points", getPoints)

	router.Run("localhost:8080")

}
