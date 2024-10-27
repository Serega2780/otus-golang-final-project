package service

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"regexp"
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
	r     *regexp.Regexp
}

func NewImageProcessingService(l *logger.Logger, conf *config.CacheConf) *ImageProcessingService {
	return &ImageProcessingService{
		log:   l,
		dir:   conf.Dir,
		cache: lru.NewCache(conf.Capacity),
		r:     regexp.MustCompile(util.PATTERN),
	}
}

func (ips *ImageProcessingService) Resize(imageInfo *model.ImageInfo, resizedKey string) ([]byte, error) {
	originalFile := imageInfo.BasicFile
	original, err := ips.readFile(originalFile)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(bytes.NewReader(original))
	if err != nil {
		return nil, fmt.Errorf(util.ErrWrongImageFormat.Error(), err)
	}
	width, height := util.ParseKey(resizedKey)
	resizedImage := resize.Resize(width, height, img, resize.Lanczos2)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return nil, err
	}
	resp := buf.Bytes()
	name := util.Substr(originalFile, strings.LastIndex(originalFile, util.SLASH)+1,
		strings.LastIndex(originalFile, util.UNDERSCORE)+1)
	ext := util.Substr(originalFile, strings.LastIndex(originalFile, util.DOT),
		len(originalFile))
	resizedFile := name + resizedKey + ext
	err = ips.saveFile(imageInfo.BasicDir+resizedFile, resp)
	if err != nil {
		ips.log.Errorf("%s file save failed", resizedFile)
		return nil, err
	}
	imageInfo.Files[resizedKey] = imageInfo.BasicDir + resizedFile

	return resp, nil
}

func (ips *ImageProcessingService) AddRoot(src []byte, imageInfo *model.ImageInfo,
	infoKey string,
) (*model.ImageInfo, error) {
	ips.mkDir(imageInfo.BasicDir)
	err := ips.saveFile(imageInfo.BasicFile, src)
	if err != nil {
		return nil, err
	}
	_, oldest := ips.cache.Set(lru.Key(infoKey), imageInfo)
	if oldest != nil {
		info := oldest.(*model.ImageInfo)
		ips.removeDirRecursive(info.BasicDir)
	}

	return ips.Get(infoKey)
}

func (ips *ImageProcessingService) Get(key string) (*model.ImageInfo, error) {
	v, b := ips.cache.Get(lru.Key(key))
	if !b {
		return nil, util.ErrNotExist
	}
	imageInfo := v.(*model.ImageInfo)

	return imageInfo, nil
}

func (ips *ImageProcessingService) GetResized(imageInfo *model.ImageInfo, keyResized string) ([]byte, error) {
	fileResized := imageInfo.GetFile(keyResized)

	return ips.readFile(fileResized)
}

func (ips *ImageProcessingService) ProcessPath(path string) (string, string, *model.ImageInfo, error) {
	width, err := strconv.Atoi(util.Substr(path, 0, strings.Index(path, util.SLASH)))
	if err != nil {
		er := fmt.Errorf(util.ErrPathVariableWrong.Error(), util.WIDTH)
		ips.log.Errorf("%v", er)
		return "", "", nil, er
	}
	path, _ = strings.CutPrefix(path, util.Substr(path, 0, strings.Index(path, util.SLASH)+1))
	height, err := strconv.Atoi(util.Substr(path, 0, strings.Index(path, util.SLASH)))
	if err != nil {
		er := fmt.Errorf(util.ErrPathVariableWrong.Error(), util.HEIGHT)
		ips.log.Errorf("%v", er)
		return "", "", nil, er
	}
	path, _ = strings.CutPrefix(path, util.Substr(path, 0, strings.Index(path, util.SLASH)+1))
	if len(path) == 0 {
		er := fmt.Errorf(util.ErrPathVariableWrong.Error(), util.URL)
		ips.log.Errorf("%v", er)
		return "", "", nil, er
	}
	remoteHost, subDir, fileName, err := util.ParsePath(path)
	if err != nil {
		return "", "", nil, err
	}
	bd := ips.dir + remoteHost + util.SLASH + subDir + util.SLASH
	info := model.NewImageInfo(bd+fileName, bd)
	rk := strconv.Itoa(width) + util.UNDERSCORE + strconv.Itoa(height)

	return path, rk, info, nil
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

func (ips *ImageProcessingService) removeDirRecursive(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		ips.log.Errorf("%v", err)
	}
}

func (ips *ImageProcessingService) mkDir(dir string) {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		ips.log.Errorf("%v", err)
	}
}

func (ips *ImageProcessingService) closeFile(f *os.File, fileName string) {
	err := f.Close()
	if err != nil {
		ips.log.Errorf("error close %s %v", fileName, err)
	}
}
