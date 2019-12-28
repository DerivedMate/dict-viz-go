package main

import (
	"math"
	"strings"
)

func CosVecs(a, b []int) float64 {
	dotProd := 0
	sumA, sumB := 0, 0
	for i := range a {
		dotProd += a[i] * b[i]
		sumA += a[i] * a[i]
		sumB += b[i] * b[i]
	}

	den := math.Sqrt(float64(sumA)) * math.Sqrt(float64(sumB))

	if den == 0 {
		den = 1
	}

	return float64(dotProd) / den
}

func CosStrings(a, b string) float64 {
	a, b = strings.ToLower(a), strings.ToLower(b)

	wordsA, wordsB := strings.Split(a, ""), strings.Split(b, "")
	words := Unique(append(wordsA, wordsB...))
	vecA, vecB := make([]int, len(words)), make([]int, len(words))

	for i, w := range words {
		vecA[i] = Count(wordsA, w)
		vecB[i] = Count(wordsB, w)
	}

	return CosVecs(vecA, vecB)
}

func CosStringsAlphabet(a, b string) float64 {
	alphabet := "a &ecplfortiàmdvsh-,nk/gb1.y()523u4wjxqçíz'é’èê7!ôâîûöäáñ09?86ïùãńü"
	a, b = strings.ToLower(a), strings.ToLower(b)
	aLetters, bLetters := strings.Split(a, ""), strings.Split(b, "")

	aVec, bVec := make([]int, 0), make([]int, 0)

	for _, p := range alphabet {
		l := string(p)
		aSum, bSum := Count(aLetters, l), Count(bLetters, l)

		aVec = append(aVec, aSum)
		bVec = append(bVec, bSum)
	}

	return CosVecs(aVec, bVec)
}
