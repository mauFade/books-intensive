package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Genre  string
}

type BookService struct {
	db *sql.DB
}

func NewBookService(d *sql.DB) *BookService {
	return &BookService{
		db: d,
	}
}

func (s *BookService) CreateBook(book *Book) error {
	query := "INSERT INTO books (title, author, genre) VALUES (?, ?, ?)"

	result, err := s.db.Exec(query, book.Title, book.Author, book.Genre)

	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return err
	}

	book.ID = int(lastID)

	return nil
}

func (s *BookService) GetBooks() ([]Book, error) {
	rows, err := s.db.Query("SELECT * FROM books")

	if err != nil {
		return nil, err
	}

	var books []Book

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *BookService) GetBookByID(id int) (*Book, error) {
	query := "SELECT * FROM books WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var book Book

	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &book, nil
}

func (s *BookService) GetBookByName(name string) ([]Book, error) {
	query := "SELECT * FROM books WHERE title LIKE ?"
	rows, err := s.db.Query(query, "%"+name+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book

		if err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *BookService) UpdateBook(book *Book) error {
	query := "UPDATE books SET title = ?, author = ?, genre = ? WHERE id = ?"
	_, err := s.db.Exec(query, book.Title, book.Author, book.Genre, book.ID)

	return err
}

func (s *BookService) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := s.db.Exec(query, id)

	return err
}

func (s *BookService) SearchBooksByName(name string) ([]Book, error) {
	query := "SELECT id, title, author, genre FROM books WHERE title LIKE ?"
	rows, err := s.db.Query(query, "%"+name+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book

		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *BookService) SimulateReading(bookId int, duration time.Duration, results chan<- string) {
	book, err := s.GetBookByID(bookId)

	if err != nil || book == nil {
		results <- fmt.Sprintf("Book with id %d not found", bookId)
	}

	time.Sleep(duration)

	results <- fmt.Sprintf("Book with title %s read", book.Title)
}

func (s *BookService) SimulateMultipleReadings(bookIDs []int, duration time.Duration) []string {
	results := make(chan string, len(bookIDs))

	for _, id := range bookIDs {
		go func(id int) {
			s.SimulateReading(id, duration, results)
		}(id)
	}

	var responses []string

	for range bookIDs {
		responses = append(responses, <-results)
	}

	close(results)

	return responses
}
