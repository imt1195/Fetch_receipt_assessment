package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type items struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
type receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Total        string  `json:"total"`
	Items        []items `json:"items"`
}

func main() {
	receiptInfo, err := os.ReadFile("./examples/Test1.json")
	if err != nil {
		fmt.Println(err)
	}

	receiptData := receipt{}
	marshErr := json.Unmarshal([]byte(receiptInfo), &receiptData)

	if marshErr != nil {
		fmt.Println(marshErr)
	}

	calculatePoints(receiptData)
}

func calculatePoints(receiptData receipt) {
	var points = 0
	points += countChars(receiptData.Retailer)
	points += PointsFromTotal(receiptData.Total)
	points += purchaseDatePoints(receiptData.PurchaseDate)
	points += purchaseTimePoints(receiptData.PurchaseTime)
	points += PointsFromNumberOfItems(len(receiptData.Items))
	points += TrimmedLengthPoints(receiptData.Items)

	fmt.Println("Total Points:", points)
}

func countChars(retailer string) int {
	count := 0

	for _, char := range retailer {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	fmt.Println("countChars points: ", count)
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
	fmt.Println("Points From Total: ", points)
	return points
}

func PointsFromNumberOfItems(numItems int) int {
	fmt.Println("Number of items points: ", (numItems/2)*5)
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
	fmt.Println("Length points: ", points)
	return points
}

func purchaseDatePoints(date string) int {
	dayDigit, _ := strconv.Atoi(date[len(date)-1:])
	if dayDigit%2 == 1 {
		fmt.Println("Date points: ", 6)
		return 6
	}
	fmt.Println("Date points: ", 0)
	return 0
}

func purchaseTimePoints(time string) int {
	timeInt, _ := strconv.Atoi(strings.ReplaceAll(time, ":", ""))
	if (timeInt > 1400) && (timeInt < 1600) {
		fmt.Println("Time points: ", 10)
		return 10
	}
	fmt.Println("Time points: ", 0)
	return 0
}
