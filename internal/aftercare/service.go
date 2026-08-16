package aftercare

import (
	"context"
	"errors"
	"unicode/utf8"
)

type Observer func(View)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	if repository == nil {
		panic("aftercare repository is required")
	}
	return &Service{repository: repository}
}

func (s *Service) Submit(ctx context.Context, description string, observer Observer) View {
	publish(observer, LoadingView())

	normalized := normalizeDescription(description)
	if normalized == "" {
		return ErrorView("description_required", "description is required")
	}

	normalized = limitDescription(normalized, DescriptionLimit)
	ticket, err := s.repository.Save(ctx, Ticket{
		Description:    normalized,
		CharacterCount: utf8.RuneCountInString(normalized),
	})
	if err != nil {
		return ErrorView("storage_error", "ticket could not be stored")
	}
	return ResultView(ticket)
}

func (s *Service) Find(ctx context.Context, id string) View {
	if id == "" {
		return ErrorView("ticket_id_required", "ticket id is required")
	}

	ticket, err := s.repository.Find(ctx, id)
	if errors.Is(err, ErrTicketNotFound) {
		return ErrorView("ticket_not_found", "ticket was not found")
	}
	if err != nil {
		return ErrorView("storage_error", "ticket could not be loaded")
	}
	return ResultView(ticket)
}

func publish(observer Observer, view View) {
	if observer != nil {
		observer(view)
	}
}
