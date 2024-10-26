package util

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrFileNameFormat    = errors.New("file name must must match the pattern '%s'")
	ErrFileNameTooShort  = errors.New("file name length must be greater 0")
	ErrWrongDimensions   = errors.New("wrong dimensions")
	ErrNotExist          = errors.New("image info not exists in the cache")
	ErrMethodNotAllowed  = errors.New("method %s not not allowed on uri %s")
	ErrPathVariableWrong = errors.New("path variable %s has wrong value")
	ErrWrongImageType    = errors.New("image must be of type either jpg or jpeg")
	ErrWrongImageFormat  = errors.New("the file received is not of valid format either jpg or jpeg: %v")
	ErrNon200Status      = errors.New("http response status differs from 2XX")
)

const (
	HTTP       = "http://"
	WIDTH      = "width"
	HEIGHT     = "height"
	URL        = "url"
	SLASH      = "/"
	JPG        = ".jpg"
	JPEG       = ".jpeg"
	UNDERSCORE = "_"
	DOT        = "."
	X          = "x"
	PATTERN    = ".+_\\d+x\\d+\\.(jpg|jpeg)$"
)

func Substr(str string, start, end int) string {
	return strings.TrimSpace(str[start:end])
}

func GetFileName(path string) string {
	return Substr(path, strings.LastIndex(path, SLASH)+1, len(path))
}

func ParseKey(resizedKey string) (width, height uint) {
	dims := strings.Split(resizedKey, UNDERSCORE)
	w, _ := strconv.Atoi(dims[0])
	h, _ := strconv.Atoi(dims[1])

	return uint(w), uint(h)
}
