package util

import (
	"errors"
	"strings"
)

var (
	ErrFileNameFormat    = errors.New("file name must end with '_d{1,4}xd{1,4}.<extension>'")
	ErrWrongDimensions   = errors.New("wrong dimensions")
	ErrNotExist          = errors.New("image info not exists in the cache")
	ErrMethodNotAllowed  = errors.New("method %s not not allowed on uri %s")
	ErrPathVariableWrong = errors.New("path variable %s has wrong value")
	ErrWrongImageType    = errors.New("image must be of type either jpg or jpeg")
)

const (
	HTTP       = "http://"
	WIDTH      = "width"
	HEIGHT     = "height"
	URL        = "url"
	SLASH      = "/"
	JPG        = "jpg"
	JPEG       = "jpeg"
	UNDERSCORE = "_"
	DOT        = "."
	X          = "x"
)

func Substr(str string, start, end int) string {
	return strings.TrimSpace(str[start:end])
}
