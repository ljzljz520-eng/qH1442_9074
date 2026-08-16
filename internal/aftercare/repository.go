package aftercare

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrTicketNotFound = errors.New("ticket not found")

type Repository interface {
	Save(context.Context, Ticket) (Ticket, error)
	Find(context.Context, string) (Ticket, error)
}

type MemoryRepository struct {
	mu      sync.RWMutex
	nextID  int
	tickets map[string]Ticket
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{tickets: make(map[string]Ticket)}
}

func (r *MemoryRepository) Save(ctx context.Context, ticket Ticket) (Ticket, error) {
	if err := ctx.Err(); err != nil {
		return Ticket{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if ticket.ID == "" {
		r.nextID++
		ticket.ID = fmt.Sprintf("T-%06d", r.nextID)
	}
	r.tickets[ticket.ID] = ticket
	return ticket, nil
}

func (r *MemoryRepository) Find(ctx context.Context, id string) (Ticket, error) {
	if err := ctx.Err(); err != nil {
		return Ticket{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	ticket, ok := r.tickets[id]
	if !ok {
		return Ticket{}, ErrTicketNotFound
	}
	return ticket, nil
}
