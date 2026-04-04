package services

import (
	"context"
	"errors"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type bookService struct {
	repo output.BookRepository
}

// NewBookService returns an input.BookService backed by the given repository.
func NewBookService(repo output.BookRepository) input.BookService {
	return &bookService{repo: repo}
}

func (s *bookService) List(ctx context.Context, opts query.QueryOptions) ([]domain.Book, int64, error) {
	return s.repo.FindAll(ctx, opts)
}

func (s *bookService) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *bookService) Create(ctx context.Context, book *domain.Book) error {
	return s.repo.Save(ctx, book)
}

func (s *bookService) Update(ctx context.Context, id string, req *domain.UpdateBookReq) (*domain.Book, error) {
	book, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}

func (s *bookService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
