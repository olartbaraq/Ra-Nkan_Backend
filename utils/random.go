package utils

import (
	"log"
	"math/rand"
	"strconv"
)

var alphabets = "abcdefghijklmnopqrstuvwxyz"
var numbers = "0123456789"

func RandomString(r int) string {
	wholeLetter := []rune{}
	k := len(alphabets)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		wholeLetter = append(wholeLetter, rune(alphabets[index]))
	}
	return string(wholeLetter)
}

func RandIntegers(r int) string {
	wholeFigure := []rune{}
	k := len(numbers)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		wholeFigure = append(wholeFigure, rune(numbers[index]))
	}
	return string(wholeFigure)
}

func randIntQty(r int) int32 {
	wholeQty := []rune{}
	k := len(numbers)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		wholeQty = append(wholeQty, rune(numbers[index]))
	}
	price := string(wholeQty)
	value, err := strconv.Atoi(price)
	if err != nil {
		log.Fatal("COuld not convert string to integer", err)
	}
	return int32(value)
}

func RandomEmail() string {
	return RandomString(7) + "@testing.com"
}

func RandomPhone() string {
	return RandIntegers(11)
}

func RandomName() string {
	return RandomString(5)
}

func RandomAddress() string {
	return RandomString(30)
}

func RandomText() string {
	return RandomString(100)
}

func RandomPrice() string {
	return RandIntegers(6) + "." + RandIntegers(2)
}

func RandomQty() int32 {
	return randIntQty(3)
}
