package service

import (
	"bytes"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/config"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/logger"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/lru"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/model"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/util"
	"github.com/nfnt/resize"
)

type ImageProcessingService struct {
	log   *logger.Logger
	dir   string
	cache lru.Cache
}

func NewImageProcessingService(l *logger.Logger, conf *config.CacheConf) *ImageProcessingService {
	return &ImageProcessingService{
		log:   l,
		dir:   conf.Dir,
		cache: lru.NewCache(conf.Capacity),
	}
}

func (ips *ImageProcessingService) Resize(original []byte, width, height uint) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(original))
	if err != nil {
		return nil, err
	}
	resizedImage := resize.Resize(width, height, img, resize.Lanczos2)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (ips *ImageProcessingService) Add(src []byte, headers http.Header, key string) ([]byte, error) {
	image := model.NewImageInfo(headers)
	_ = ips.cache.Set(lru.Key(key), image)
	err := ips.saveFile(key, src)
	if err != nil {
		return nil, err
	}

	return src, nil
}

func (ips *ImageProcessingService) Get(key string) ([]byte, http.Header, error) {
	v, b := ips.cache.Get(lru.Key(key))
	if !b {
		return nil, nil, util.ErrNotExist
	}
	image := v.(*model.ImageInfo)
	body, err := ips.readFile(key)
	if err != nil {
		return nil, nil, err
	}

	return body, image.Headers, nil
}

func (ips *ImageProcessingService) SetKey(width, height int, file string) (string, error) {
	underscore := strings.LastIndex(file, util.UNDERSCORE)
	if underscore == -1 {
		return "", util.ErrFileNameFormat
	}
	dot := strings.LastIndex(file, util.DOT)
	if dot == -1 {
		return "", util.ErrFileNameFormat
	}
	name := util.Substr(file, 0, underscore+1)
	ext := util.Substr(file, dot, len(file))
	dimensions := util.Substr(file, underscore+1, dot)
	dims := strings.Split(dimensions, util.X)
	w, err := strconv.Atoi(dims[0])
	if err != nil {
		return "", err
	}
	h, err := strconv.Atoi(dims[1])
	if err != nil {
		return "", err
	}
	if width > w || height > h {
		return "", util.ErrWrongDimensions
	}
	return ips.dir + name + strconv.Itoa(width) + util.X + strconv.Itoa(height) + ext, nil
}

func (ips *ImageProcessingService) saveFile(fileName string, data []byte) error {
	f, err := openFile(fileName)
	if err != nil && os.IsNotExist(err) {
		f, err = os.Create(fileName)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		if err != nil {
			return err
		}
		defer ips.closeFile(f, fileName)
		return nil
	}
	if err != nil {
		return err
	}
	defer ips.closeFile(f, fileName)

	return nil
}

func (ips *ImageProcessingService) readFile(fileName string) ([]byte, error) {
	f, err := openFile(fileName)
	if err != nil {
		return nil, err
	}
	defer ips.closeFile(f, fileName)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func openFile(fileName string) (*os.File, error) {
	return os.Open(fileName)
}

func (ips *ImageProcessingService) closeFile(f *os.File, fileName string) {
	err := f.Close()
	if err != nil {
		ips.log.Errorf("error close %s %v", fileName, err)
	}
}
