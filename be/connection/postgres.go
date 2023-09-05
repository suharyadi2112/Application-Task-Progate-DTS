package connect

import(
	"fmt"
    "database/sql"
    _ "github.com/lib/pq"// postgres
    "github.com/joho/godotenv"
    "os"
    "log"
)

func Dbcon() (*sql.DB, error) {

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Fail Load .env")
    }

    // Get the value of an Environment Variable
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }


    return db, nil
}
