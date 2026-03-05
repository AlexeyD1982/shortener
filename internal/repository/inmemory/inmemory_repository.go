package inmemory

import (
	"github.com/AlexeyD1982/shortener/internal/model"
	"github.com/AlexeyD1982/shortener/pkg/errors"
)

type InMemoryRepository struct {
	urls map[string]string
}

func NewStorage() *InMemoryRepository {
	return &InMemoryRepository{urls: make(map[string]string)}
}

func (r *InMemoryRepository) SaveURL(url *model.URL) error {
	if _, exist := r.urls[url.Short]; exist {
		return errors.ErrNotUnique
	}
	r.urls[url.Short] = url.Origin
	return nil
}

func (r *InMemoryRepository) GetURL(short string) (*model.URL, error) {
	origin, exist := r.urls[short]
	if !exist {
		return nil, errors.ErrNotFound
	}
	return &model.URL{Origin: origin, Short: short}, nil
}
