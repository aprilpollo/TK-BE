package services

import (
	"context"

	"aprilpollo/internal/core/domain"
	"aprilpollo/internal/core/ports/input"
	"aprilpollo/internal/core/ports/output"
	"aprilpollo/internal/pkg/query"
)

type calendarService struct {
	repo output.CalendarRepository
}

func NewCalendarService(repo output.CalendarRepository) input.CalendarService {
	return &calendarService{repo: repo}
}

func (s *calendarService) List(ctx context.Context, opts query.QueryOptions, projectID int64) ([]domain.Task, int64, error) {
	return s.repo.Find(ctx, opts, projectID)
}

func (s *calendarService) ListStatus(ctx context.Context, opts query.QueryOptions, projectID int64) ([]domain.TaskStatus, error) {
	return s.repo.FindStatus(ctx, opts, projectID)
}

func (s *calendarService) ListPriority(ctx context.Context) ([]domain.TaskPriority, error) {
	return s.repo.FindPriority(ctx)
}
