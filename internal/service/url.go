package service

import (
	"fmt"

	"github.com/AlexeyD1982/shortener/internal/model"
)

type URLRepository interface {
	SaveURL(url *model.URL) error
	GetURL(short string) (*model.URL, error)
}

type URLService struct {
	urlRepository URLRepository
}

func NewURLService(urlRepository URLRepository) *URLService {
	return &URLService{urlRepository: urlRepository}
}

func (s *URLService) SaveURL(originURL string) string {
	var shortURL string
	for {
		shortURL = generateRandomString(8)
		if err := s.urlRepository.SaveURL(&model.URL{Origin: originURL, Short: shortURL}); err == nil {
			break
		}
	}
	return shortURL
}

func (s *URLService) ResolveURL(shortURL string) (string, error) {
	url, err := s.urlRepository.GetURL(shortURL)
	if err != nil {
		return "", fmt.Errorf("repository error: %w", err)
	}
	return url.Origin, nil
}
