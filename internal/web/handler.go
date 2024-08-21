package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mauFade/books-intensive/internal/service"
)

type BookHandler struct {
	service *service.BookService
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetBooks()

	if err != nil {
		http.Error(w, "Erro getting books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book service.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateBook(&book); err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookByID(id)

	if err != nil {
		http.Error(w, "failed to get book", http.StatusInternalServerError)
		return
	}

	if book == nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
