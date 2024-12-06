package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type items struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
type receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Total        string    `json:"total"`
	Items        []items   `json:"items"`
	Id           uuid.UUID `json:"id"`
	Points       int       `json:"points"`
}

func main() {
	//get json file of receipt
	receiptInfo, err := os.ReadFile("./examples/Test1.json")
	if err != nil {
		fmt.Println(err)
	}

	receiptData := receipt{}
	marshErr := json.Unmarshal([]byte(receiptInfo), &receiptData)

	if marshErr != nil {
		fmt.Println(marshErr)
	}

	fmt.Println(receiptData)
	var totalPoints = calculatePoints(receiptData)

	id := uuid.New()
	receiptData.Id = id
	receiptData.Points = totalPoints

	fmt.Println(receiptData)
	fmt.Println("Generated ID: ", receiptData.Id)
	fmt.Println("Points: ", receiptData.Points)

	router := gin.Default()
	router.POST("/receipts/process", POST_RECEIPT)
	// router.GET("/receipts/:id/points", GET_POINTS_BY_RECEIPT_ID)
	router.Run("localhost:9090")
}

// POST process receipts
func POST_RECEIPT(context *gin.Context /*data receipt*/) {
	var data receipt
	if err := context.BindJSON(&data); err != nil {
		return
	}
	context.IndentedJSON(http.StatusOK, data)
	// receipt, _ := json.Marshal(data)
	// req, err := http.NewRequest("POST", "https://localhost:9090/receipts/process", bytes.NewBuffer(receipt))
	// if err != nil {
	// 	return
	// }

	// req.Header.Set("Content-Type", "application/json")
}

// GET points
// func GET_POINTS_BY_RECEIPT_ID(context *gin.Context /*data receipt*/) {
// 	id := context.Param("id")
// 	//find receipt by the id
// 	// if id == data.Id.String(){
// 	// 	return data.Points
// 	// }
// 	//isolate the points section of the JSON or the struct
// 	//print to screen
// }

// Point calculation functions
func calculatePoints(receiptData receipt) int {
	var points = 0
	points += countChars(receiptData.Retailer)
	points += PointsFromTotal(receiptData.Total)
	points += purchaseDatePoints(receiptData.PurchaseDate)
	points += purchaseTimePoints(receiptData.PurchaseTime)
	points += PointsFromNumberOfItems(len(receiptData.Items))
	points += TrimmedLengthPoints(receiptData.Items)

	return points
}

func countChars(retailer string) int {
	count := 0

	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}

// convert float to int by multiplying by 100 and converting to int
func PointsFromTotal(total string) int {
	totalInt, _ := strconv.Atoi(strings.ReplaceAll(total, ".", ""))
	var points = 0
	if totalInt%100 == 0 {
		points += 50
	}
	if totalInt%25 == 0 {
		points += 25
	}
	return points
}

func PointsFromNumberOfItems(numItems int) int {
	return (numItems / 2) * 5
}

func TrimmedLengthPoints(items []items) int {
	var points = 0
	var trimmString = " "
	for _, item := range items {
		var trimmedString = strings.Trim(item.ShortDescription, trimmString)
		var trimStringLen = len(trimmedString)
		if trimStringLen%3 == 0 {
			var price, err = strconv.ParseFloat(item.Price, 64)
			if err != nil {
				fmt.Println(err)
			}
			var lengthPoints = math.Ceil(price * .2)
			points += int(lengthPoints)
		}
	}
	return points
}

func purchaseDatePoints(date string) int {
	dayDigit, _ := strconv.Atoi(date[len(date)-1:])
	if dayDigit%2 == 1 {
		return 6
	}
	return 0
}

func purchaseTimePoints(time string) int {
	timeInt, _ := strconv.Atoi(strings.ReplaceAll(time, ":", ""))
	if (timeInt > 1400) && (timeInt < 1600) {
		return 10
	}
	return 0
}
