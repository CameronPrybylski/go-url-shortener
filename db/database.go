package db

import (
	"database/sql"
	"fmt"
	"go-url-shortener/model"
	"log"

	_ "github.com/lib/pq"
// "github.com/pelletier/go-toml/query"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "go_user"
	password = "Puma1234"
	dbname   = "url_shortener"
)

var Db *sql.DB

func InitDB() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Connected to PostgreSQL!")
}

func CheckURL(originalURL string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM urls
		WHERE original_url = $1
	`
	err := Db.QueryRow(query, originalURL).Scan(&count)
	if err != nil{
		log.Fatal("Error selecting URL:", err)
		return false, err
	}
	return count > 0, nil
}


func CreateURL(originalURL, shortenedURL string) ( model.URL, error){
	var newURL model.URL

	// Insert URL into the database
	query := `
		INSERT INTO urls (original_url, short_code)
		VALUES ($1, $2)
		RETURNING id, original_url, short_code, created_at
	`

	err := Db.QueryRow(query, originalURL, shortenedURL).Scan(&newURL.ID, &newURL.OriginalURL, &newURL.ShortenedCode, &newURL.CreatedAt)

	if err != nil{
		log.Fatal("Error inserting URL:", err)
		return model.URL{}, err
	}

	return newURL, nil
	
}

func GetURL(shortenedURL string) (model.URL, error){
	var newURL model.URL

	query := `
		SELECT id, original_url, short_code, created_at, visit_count
		FROM urls
		WHERE short_code = $1
	`
	err := Db.QueryRow(query, shortenedURL).Scan(&newURL.ID, &newURL.OriginalURL, &newURL.ShortenedCode, &newURL.CreatedAt, &newURL.VisitCount)

	if err != nil{
		log.Fatal("Error retrieving URL:", err)
		return model.URL{}, err
	}

	return newURL, nil
}

func GetFromORGURL(originalURL string) (model.URL, error){
	var newURL model.URL

	query := `
		SELECT id, original_url, short_code, created_at
		FROM urls
		WHERE original_url = $1
	`
	err := Db.QueryRow(query, originalURL).Scan(&newURL.ID, &newURL.OriginalURL, &newURL.ShortenedCode, &newURL.CreatedAt)

	if err != nil{
		log.Fatal("Error retrieving URL:", err)
		return model.URL{}, err
	}

	return newURL, nil
}


func IncrementVisitCount(originalURL string, visitCount int)(model.URL, error){
	visitCount++
	var url model.URL
	query := `
		UPDATE urls
		SET visit_count = $2, last_accessed = current_timestamp
		WHERE original_url = $1
		RETURNING id, original_url, short_code, created_at, visit_count
	`
	err := Db.QueryRow(query, originalURL, visitCount).Scan(&url.ID, &url.OriginalURL, &url.ShortenedCode, &url.CreatedAt, &url.VisitCount)
	if err != nil{
		log.Fatal("Error updating URL:", err)
		return model.URL{}, err
	}

	return url, nil

}

func TopXVisited(numOfURLs int, order string) ([]model.URL, error){
	var mostVisitedURLs []model.URL
	if order != "ASC" && order != "DESC" {
		return nil, fmt.Errorf("invalid order parameter: %s", order)
	}

	query := fmt.Sprintf(`
			SELECT id, original_url, short_code, created_at, visit_count 
			FROM urls 
			ORDER BY visit_count %s 
			LIMIT $1
		`, order) // Inject ASC or DESC safely
	rows, err := Db.Query(query, numOfURLs)
	if err != nil{
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		var url model.URL
		err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortenedCode, &url.CreatedAt, &url.VisitCount)
		if err != nil{
			return nil, err
		}
		mostVisitedURLs = append(mostVisitedURLs, url)
	}

	return mostVisitedURLs, nil
}

func LastAccesed(numOfURLs int, order string) ([]model.URL, error){
	var urls []model.URL
	if order != "ASC" && order != "DESC" {
		return nil, fmt.Errorf("invalid order parameter: %s", order)
	}
	query := fmt.Sprintf(`
			SELECT id, original_url, short_code, created_at, visit_count, last_accessed
			FROM urls 
			ORDER BY last_accessed %s 
			LIMIT $1
		`, order) // Inject ASC or DESC safely
	rows, err := Db.Query(query, numOfURLs)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var url model.URL
		err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortenedCode, &url.CreatedAt, &url.VisitCount, &url.LastAccessed)
		if err != nil{
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

