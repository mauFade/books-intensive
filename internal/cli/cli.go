package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/mauFade/books-intensive/internal/service"
)

type BookCLI struct {
	service *service.BookService
}

func NewBookCLI(s *service.BookService) *BookCLI {
	return &BookCLI{
		service: s,
	}
}

func (c *BookCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go-books <command> [args]")

		return
	}

	command := os.Args[1]

	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("usage: go-books search <query>")
			return

		}

		books, err := c.searchBookByName(os.Args[2])

		if err != nil {
			fmt.Println("Error getting book by name: ", err.Error())
			return
		}

		fmt.Printf("%d books found: ", len(books))

		for _, book := range books {
			fmt.Printf("ID: %d, Title: %s, Author: %s, Genre: %s\n\n", book.ID, book.Title, book.Author, book.Genre)
		}
	}
}

func (c *BookCLI) searchBookByName(name string) ([]service.Book, error) {
	books, err := c.service.SearchBooksByName(name)

	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, errors.New("no books found")
	}

	return books, nil
}
