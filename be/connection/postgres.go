package connect

import(
	"fmt"
    "database/sql"
    _ "github.com/lib/pq"// postgres
    "os"
)

func Dbcon() (*sql.DB, error) {

    // Set Environment Variables
    os.Setenv("SITE_TITLE", "Progate Application Task")
    os.Setenv("DB_HOST", "localhost")
    os.Setenv("DB_PORT", "5432")
    os.Setenv("DB_USERNAME", "postgres")
    os.Setenv("DB_PASSWORD", "12345678")
    os.Setenv("DB_NAME", "postgres")

    // Get the value of an Environment Variable
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    //https://zetcode.com/golang/string-format/ (Referensi string format %s, %d, dsb)

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }

    return db, nil
}
