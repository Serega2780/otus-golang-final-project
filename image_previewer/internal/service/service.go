package service

import "net/http"

type ImageService interface {
	SetKey(width, height int, file string) (string, error)
	Add(src []byte, headers http.Header, key string) ([]byte, error)
	Get(key string) ([]byte, http.Header, error)
	Resize(original []byte, width, height uint) ([]byte, error)
}
