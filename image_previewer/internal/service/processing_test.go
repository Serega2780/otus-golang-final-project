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
		key, err := ips.SetKey(100, 50, "_gopher_original_1024x504.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/_gopher_original_100x50.jpg", key)
	})
	t.Run("gopher_1024x252.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_1024x252.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/gopher_100x50.jpg", key)
	})
	t.Run("gopher_2000x1000.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_2000x1000.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/gopher_100x50.jpg", key)
	})
	t.Run("gopher_200x700.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_200x700.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/gopher_100x50.jpg", key)
	})
	t.Run("gopher_256x126.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_256x126.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/gopher_100x50.jpg", key)
	})
	t.Run("gopher_333x666.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_333x666.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/gopher_100x50.jpg", key)
	})
	t.Run("gopher_500x500.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_500x500.jpg")
		require.Nil(t, err)
		require.Equal(t, "/tmp/gopher_100x50.jpg", key)
	})
	t.Run("gopher_50x50.jpg", func(t *testing.T) {
		key, err := ips.SetKey(100, 50, "gopher_50x50.jpg")
		require.Equal(t, util.ErrWrongDimensions, err)
		require.Equal(t, "", key)
	})
}
