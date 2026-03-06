package inmemory

import (
	"github.com/AlexeyD1982/shortener/internal/model"
	"github.com/AlexeyD1982/shortener/pkg/local_errors"
)

type Repository struct {
	urls map[string]string
}

func NewStorage() *Repository {
	return &Repository{urls: make(map[string]string)}
}

func NewStorageWithData(data map[string]string) *Repository {
	return &Repository{urls: data}
}

func (r *Repository) SaveURL(url *model.URL) error {
	if _, exist := r.urls[url.Short]; exist {
		return local_errors.ErrNotUnique
	}
	r.urls[url.Short] = url.Origin
	return nil
}

func (r *Repository) GetURL(short string) (*model.URL, error) {
	origin, exist := r.urls[short]
	if !exist {
		return nil, local_errors.ErrNotFound
	}
	return &model.URL{Origin: origin, Short: short}, nil
}
