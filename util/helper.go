package util

import (
	"angsuran-service/internal/controller/request"
	"errors"
	"log"
	"math"
	"strings"
	"time"
)

func RoundToTwoDecimal(value float64) float64 {
	return Round(value, 2)
}

func Round(value float64, decimalPlaces int) float64 {
	precision := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*precision) / precision
}

func Pow(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

func ValidateRequestBody(requestBody request.AngsuranRequest) error {
	var missingFields []string

	if requestBody.Plafond == 0 {
		missingFields = append(missingFields, "plafond")
	}
	if requestBody.LamaPinjaman == 0 {
		missingFields = append(missingFields, "lama_pinjaman")
	}
	if requestBody.Bunga == 0 {
		missingFields = append(missingFields, "bunga")
	}
	if requestBody.TanggalMulai == "" {
		missingFields = append(missingFields, "tanggal_mulai")
	}

	if len(missingFields) > 0 {
		errorMessage := "Bad request: Field(s) " + strings.Join(missingFields, ", ") + " is/are required"
		return errors.New(errorMessage)
	}

	return nil
}

func ParseTanggal(tanggalStr string) time.Time {
	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		log.Printf("Error parsing tanggal: %v\n", err)
		return time.Time{}
	}
	return tanggal
}
