package main

import (
	"fmt"
	"go-url-shortener/db"
	"go-url-shortener/model"
	"go-url-shortener/utils"
	"log"
	"net/http"
	"strings"
)
func main(){
	db.InitDB()
	defer db.Db.Close()
	fmt.Println("Working")
	
	state := "contine"

	for state != "quit"{
		state = homeScreen()
	}
	
}

func homeScreen() string{
	fmt.Println("Welcome! Please select an option")
	fmt.Println("1. Create Short URL")
	fmt.Println("2. Anaylytics")
	fmt.Println("3. Quit")
	
	var choice int
	fmt.Scanln(&choice)
	state := "continue"
	switch(choice){
		case 1:
			state = createShortURL()
		case 2:
			analytics()
		case 3:
			state = "quit"
		default:
			fmt.Println("Please enter valid choice")
	}
	return state
}

func createShortURL() string {
	// Set up the routes
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
	return "continue"
}

func analytics(){
	fmt.Println("1. Most visited URLs")
	fmt.Println("2. Least visited URLs")
	fmt.Println("3. Most recently visited URLs")
	fmt.Println("4. Least recently visited URLs")
	var choice int
	fmt.Scanln(&choice)
	var numberURLs int
	fmt.Println("What number of urls would you like?")
	fmt.Scanln(&numberURLs)
	switch(choice){
		case 1:
			fmt.Println("Most visited")
			findXVisited(numberURLs, "DESC")
		case 2:
			fmt.Println("Least visited")
			findXVisited(numberURLs, "ASC")
		case 3:
			fmt.Println("Most recently visited")
			findRecentVisit(numberURLs, "DESC")
		case 4:
			fmt.Println("Least recently visited")
			findRecentVisit(numberURLs, "ASC")
		default:
			fmt.Println("Please enter valid choice")
	}
}

func findXVisited(numberURLS int, order string){
	urls, err := db.TopXVisited(numberURLS, order)
	if err != nil{
		return
	}
	printURLs(urls)
}

func findRecentVisit(numberURLs int, order string){
	urls, err := db.LastAccesed(numberURLs, order)
	if err != nil{
		return
	}
	printURLs(urls)
}

func printURLs(urls []model.URL){
	for i := range urls{
		fmt.Println("URL:", urls[i].OriginalURL, "Visits:", urls[i].VisitCount, "Last Accessed:", urls[i].LastAccessed)
	}
}

// shortenHandler creates a shortened URL from the original URL
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// Get the original URL from the form
	fmt.Println("Navigate to http://localhost:8080/ then return to terminal")
	var originalURL string
	fmt.Println("What is the original url?")
	fmt.Scanln(&originalURL)
	fmt.Println("The original url is", originalURL)

	//var shortenedURL string
	var er error
	var newURL model.URL
	result, err := db.CheckURL(originalURL)
	if(result && err == nil){
		newURL, er = db.GetFromORGURL(originalURL)
		if er != nil {
			http.Error(w, "Error Getting URL", http.StatusInternalServerError)
			return
		}
		fmt.Println("URL Already exists")
	}else{

		// Generate a shortened URL
		shortenedURL := utils.ShortenURL(originalURL)

		// Store the URL in the database
		newURL, err = db.CreateURL(originalURL, shortenedURL)
	}
	if err != nil {
		http.Error(w, "Error creating URL", http.StatusInternalServerError)
		return
	}

	// Return the shortened URL to the user
	fmt.Fprintf(w, "Shortened URL: http://localhost:8080/%s", newURL.ShortenedCode)
}

// redirectHandler redirects to the original URL when the shortened URL is accessed
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the shortened URL from the URL path
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	shortenedURL := strings.TrimPrefix(r.URL.Path, "/")
	// Look up the original URL in the database
	fmt.Println("Shortened url: ", shortenedURL)
	url, err1 := db.GetURL(shortenedURL)
	url, err2 := db.IncrementVisitCount(url.OriginalURL, url.VisitCount)
	if err1 != nil || err2 != nil{
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
