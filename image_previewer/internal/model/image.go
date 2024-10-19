package model

import "net/http"

type ImageInfo struct {
	Headers http.Header
}

func NewImageInfo(headers http.Header) *ImageInfo {
	return &ImageInfo{
		Headers: headers,
	}
}
