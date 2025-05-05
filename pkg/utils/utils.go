package utils

import (
	"log"
	"math"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	env := os.Getenv("SERVER_ENVIRONMENT")
	if env == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func IsFibonacci(n int) bool {
	// --- CHECK IF NUMBER IS NEGATIVE (FIBONACCI IS NOT DEFINED FOR NEGATIVE NUMBERS) ---
	if n < 0 {
		return false
	}

	// --- CHECK IF 5*n²+4 IS A PERFECT SQUARE ---
	check1 := 5*n*n + 4
	sqrt1 := int(math.Sqrt(float64(check1)))
	isPerfectSquare1 := sqrt1*sqrt1 == check1

	// --- CHECK IF 5*n²-4 IS A PERFECT SQUARE ---
	check2 := 5*n*n - 4
	sqrt2 := int(math.Sqrt(float64(check2)))
	isPerfectSquare2 := sqrt2*sqrt2 == check2

	// --- n IS FIBONACCI IF AT LEAST ONE OF THE CHECKS IS TRUE ---
	return isPerfectSquare1 || isPerfectSquare2
}
