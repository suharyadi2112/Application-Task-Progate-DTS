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

func CountNonEmptyRows(sheet *xlsx.Sheet) int {
    count := 0
    for i, row := range sheet.Rows {
        if i == 0 {
            // Skip the header row (if applicable)
            continue
        }
        // Check if the row is non-empty
        if !IsRowEmpty(row) {
            count++
        }
    }
    return count
}


func IsRowEmpty(row *xlsx.Row) bool {
	for _, cell := range row.Cells {
		if cell.String() != "" {
			return false
		}
	}
	return true
}

func IsValidExcelFile(fileName string) bool {
	// Daftar ekstensi file Excel yang diizinkan
	allowedExtensions := []string{".xlsx", ".xls"}

	// Mendapatkan ekstensi file dari nama file
	fileExtension := strings.ToLower(fileName[strings.LastIndex(fileName, "."):])

	// Memeriksa apakah ekstensi file ada dalam daftar yang diizinkan
	for _, ext := range allowedExtensions {
		if fileExtension == ext {
			return true // Ekstensi file valid
		}
	}

	return false // Ekstensi file tidak valid
}