package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/logger"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/service"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/util"
	"golang.org/x/sync/singleflight"
)

var sfg = singleflight.Group{}

type ProxyHandler struct {
	ctx     context.Context
	log     *logger.Logger
	service service.ImageService
	client  *http.Client
}

func NewProxyHandler(ctx context.Context, logger *logger.Logger, service service.ImageService) *ProxyHandler {
	return &ProxyHandler{
		ctx:     ctx,
		log:     logger,
		service: service,
		client:  initClient(),
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		resp := fmt.Sprintf(util.ErrMethodNotAllowed.Error(), r.Method, r.URL.Path)
		http.Error(w, resp, http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	width, err := strconv.Atoi(util.Substr(r.URL.Path, 0, strings.Index(path, util.SLASH)))
	p.log.Info(width)
	if err != nil {
		resp := fmt.Sprintf(util.ErrPathVariableWrong.Error(), util.WIDTH)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	path, _ = strings.CutPrefix(path, util.Substr(path, 0, strings.Index(path, util.SLASH)+1))
	height, err := strconv.Atoi(util.Substr(path, 0, strings.Index(path, util.SLASH)))
	p.log.Info(height)
	if err != nil {
		resp := fmt.Sprintf(util.ErrPathVariableWrong.Error(), util.HEIGHT)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	path, _ = strings.CutPrefix(path, util.Substr(path, 0, strings.Index(path, util.SLASH)+1))
	if len(path) == 0 {
		resp := fmt.Sprintf(util.ErrPathVariableWrong.Error(), util.URL)
		http.Error(w, resp, http.StatusBadRequest)
		return
	}
	fileName := util.Substr(path, strings.LastIndex(path, util.SLASH)+1, len(path))

	ext := util.Substr(fileName, strings.LastIndex(fileName, util.DOT)+1, len(fileName))
	if ext != util.JPG && ext != util.JPEG {
		http.Error(w, util.ErrWrongImageType.Error(), http.StatusBadRequest)
		return
	}

	key, err := p.service.SetKey(width, height, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, headers, err := p.service.Get(key)
	if errors.Is(err, util.ErrNotExist) {
		var status int
		v, err, _ := sfg.Do(key, func() (interface{}, error) {
			status, b, err = p.proxyRequest(w, r, path)
			if err != nil {
				return nil, err
			}
			newBytes, err := p.service.Resize(b, uint(width), uint(height))
			if err != nil {
				return nil, err
			}
			w.Header().Del("Content-Length")
			w.Header().Set("Content-Length", strconv.Itoa(len(newBytes)))
			return p.service.Add(newBytes, w.Header(), key)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.sendResponse(w, v.([]byte), nil, status)
	} else {
		p.sendResponse(w, b, headers, http.StatusOK)
	}
}

func (p *ProxyHandler) sendResponse(w http.ResponseWriter, response []byte, headers http.Header, status int) {
	if headers != nil {
		copyHeaders(headers, w.Header())
	}
	if status != 0 {
		w.WriteHeader(status)
	}
	_, err := io.Copy(w, bytes.NewReader(response))
	if err != nil {
		p.log.Errorf("response write error: %v", err)
	}
}

func (p *ProxyHandler) proxyRequest(w http.ResponseWriter, r *http.Request, path string) (int, []byte, error) {
	targetURL := util.HTTP + path
	reqCtx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
	defer cancel()
	proxyReq, err := http.NewRequestWithContext(reqCtx, r.Method, targetURL, nil)
	if err != nil {
		p.log.Errorf("Error creating proxy request %v", err)
		return 0, nil, err
	}

	copyHeaders(r.Header, proxyReq.Header)

	resp, err := p.client.Do(proxyReq)
	if err != nil {
		p.log.Errorf("Error sending proxy request %v", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	b, err := io.ReadAll(resp.Body)

	return resp.StatusCode, b, err
}

func copyHeaders(src, dst http.Header) {
	for name, values := range src {
		for _, value := range values {
			dst.Add(name, value)
		}
	}
}

func initClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 10 * time.Second,
	}
	return &http.Client{
		Transport: tr,
	}
}
