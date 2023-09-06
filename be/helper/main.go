package helper

import(
    "strings"
    "time"
    "math/rand"
    "github.com/tealeg/xlsx"
)

func TrimDateFileName(s time.Time)(res string){
    timeString := s.Format("2006-01-02")
    trimmedDate := strings.Replace(timeString, "-", "", -1)
    return trimmedDate
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	// Daftar karakter yang mungkin
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Buat string acak
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func IsRowEmpty(row *xlsx.Row) bool {
	for _, cell := range row.Cells {
		if cell.String() != "" {
			return false
		}
	}
	return true
}
