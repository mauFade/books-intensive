package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mauFade/books-intensive/internal/cli"
	"github.com/mauFade/books-intensive/internal/service"
	"github.com/mauFade/books-intensive/internal/web"
)

func main() {
	db, err := sql.Open("mysql", "books:books@tcp(host:3306)/books")

	if err != nil {
		panic((err))
	}

	defer db.Close()

	bookService := service.NewBookService(db)

	bookHandlers := web.NewBookHandler(bookService)

	if len(os.Args) > 1 && (os.Args[1] == "simulate" || os.Args[1] == "search") {
		bookCLI := cli.NewBookCLI(bookService)
		bookCLI.Run()
		return
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /books", bookHandlers.GetBooks)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("GET /books/{id}", bookHandlers.GetBookByID)
	router.HandleFunc("PUT /books/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", bookHandlers.DeleteBook)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
