package book

import (
	"fmt"
	"github.com/behnambm/assignment/param"
	"github.com/behnambm/assignment/repo/sqlite"
	"log"
)

type Service struct {
	repo *sqlite.Repo
}

func New(repo *sqlite.Repo) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) BookList() ([]param.FullBookResponse, error) {
	booksFromDB, err := s.repo.BookList()
	if err != nil {
		log.Println("BOOK LIST SERVICE ERR", err)
		return []param.FullBookResponse{}, fmt.Errorf("couldn't get list of books")
	}

	var books []param.FullBookResponse
	for _, book := range booksFromDB {
		books = append(books, param.FullBookResponse{
			ID:          book.ID,
			Name:        book.Name,
			Price:       book.Price,
			IsPublished: book.IsPublished,
		})
	}
	return books, nil
}

func (s Service) PublishedBookList() ([]param.MinimalBookResponse, error) {
	booksFromDB, err := s.repo.PublishedBookList()
	if err != nil {
		log.Println("PUBLISHED BOOK LIST SERVICE ERR", err)
		return []param.MinimalBookResponse{}, fmt.Errorf("couldn't get list of books")
	}

	var books []param.MinimalBookResponse
	for _, book := range booksFromDB {
		books = append(books, param.MinimalBookResponse{
			ID:    book.ID,
			Name:  book.Name,
			Price: book.Price,
		})
	}
	return books, nil
}

func (s Service) CreateBook(param param.BookCreateRequest) (Book, error) {
	// check if the user is admin and the IsPublished is true then allow it otherwise make it not published
	bookFromDB, createErr := s.repo.CreateBook(param.Name, param.Price, param.IsPublished)
	if createErr != nil {
		log.Println("BOOK CREATE SERVICE ERR", createErr)
		return Book{}, fmt.Errorf("couldn't create book")
	}

	// if the user is admin then use the FullBookResponse
	return Book{
		ID:          bookFromDB.ID,
		Name:        bookFromDB.Name,
		Price:       bookFromDB.Price,
		IsPublished: bookFromDB.IsPublished,
	}, nil
}

func (s Service) UpdateBook(bookId int, param param.BookUpdateRequest) (Book, error) {
	bookFromDB, updateErr := s.repo.UpdateBook(bookId, param.Name, param.Price, param.IsPublished)
	if updateErr != nil {
		log.Println("BOOK UPDATE SERVICE ERR", updateErr)
		return Book{}, fmt.Errorf("couldn't update book")
	}

	// if the user is admin then use the FullBookResponse
	return Book{
		ID:    bookFromDB.ID,
		Name:  bookFromDB.Name,
		Price: bookFromDB.Price,
	}, nil
}

func (s Service) DeleteBook(bookId int) error {
	// check if the user is admin then allow to delete the book
	if deleteErr := s.repo.DeleteBook(bookId); deleteErr != nil {
		log.Println("BOOK DELETE SERVICE ERR", deleteErr)
		return fmt.Errorf("couldn't delete the book")
	}
	return nil
}
