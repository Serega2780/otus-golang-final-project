package service

import (
	"testing"

	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/config"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/logger"
	"github.com/Serega2780/otus-golang-final-project/image_previewer/internal/util"
	"github.com/stretchr/testify/require"
)

func TestSetKey(t *testing.T) {
	ips := NewImageProcessingService(logger.New(&config.LoggerConf{
		Level:        "info",
		Format:       "text",
		File:         "image_previewer_proxy.log",
		LogToFile:    false,
		LogToConsole: true,
	}), &config.CacheConf{Dir: "/tmp/", Capacity: 10})
	t.Run("_gopher_original_1024x504.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/_gopher_original_1024x504.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/_gopher_original_1024x504.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/_gopher_original_1024x504/_gopher_original_1024x504.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/_gopher_original_1024x504/", newImageInfo.BasicDir)
	})
	t.Run("gopher_1024x252.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/gopher_1024x252.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/gopher_1024x252.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/gopher_1024x252/gopher_1024x252.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/gopher_1024x252/", newImageInfo.BasicDir)
	})
	t.Run("gopher_2000x1000.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/gopher_2000x1000.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/gopher_2000x1000.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/gopher_2000x1000/gopher_2000x1000.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/gopher_2000x1000/", newImageInfo.BasicDir)
	})
	t.Run("gopher_200x700.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/gopher_200x700.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/gopher_200x700.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/gopher_200x700/gopher_200x700.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/gopher_200x700/", newImageInfo.BasicDir)
	})
	t.Run("gopher_256x126.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/gopher_256x126.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/gopher_256x126.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/gopher_256x126/gopher_256x126.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/gopher_256x126/", newImageInfo.BasicDir)
	})
	t.Run("gopher_333x666.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/gopher_333x666.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/gopher_333x666.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/gopher_333x666/gopher_333x666.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/gopher_333x666/", newImageInfo.BasicDir)
	})
	t.Run("gopher_500x500.jpg", func(t *testing.T) {
		infoKey, resizedKey, newImageInfo, err := ips.ProcessPath("100/50/test.ru/gopher_500x500.jpg")
		require.Nil(t, err)
		require.Equal(t, "test.ru/gopher_500x500.jpg", infoKey)
		require.Equal(t, "100_50", resizedKey)
		require.Equal(t, "/tmp/test.ru/gopher_500x500/gopher_500x500.jpg", newImageInfo.BasicFile)
		require.Equal(t, "/tmp/test.ru/gopher_500x500/", newImageInfo.BasicDir)
	})
	t.Run("gopher_50x50.jpg", func(t *testing.T) {
		_, _, _, err := ips.ProcessPath("100/50/test.ru/gopher_50x50.jpg")
		require.Equal(t, util.ErrWrongDimensions, err)
	})
}
