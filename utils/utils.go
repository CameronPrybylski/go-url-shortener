package utils

import(
	"math/rand"
	"time"
	"strconv"
)

func ShortenURL(url string) string {
	shortURL := ""
	
	for i :=0 ; i < 6; i++{
		rand.Seed(time.Now().UnixNano()) // Seed to ensure different numbers each 
		randomNumber := rand.Intn(10)    // Generates a number between 0 and 9
		shortURL += strconv.Itoa(randomNumber)
	}

	return shortURL

}
