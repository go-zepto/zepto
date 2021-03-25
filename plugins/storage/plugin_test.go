package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-zepto/zepto/plugins/storage"
	"github.com/go-zepto/zepto/plugins/storage/storagemock"
	"github.com/golang/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	sm := storagemock.NewMockStorage(ctrl)

	sm.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(&storage.FileInfo{
			AccessType:  storage.Public,
			Key:         "path/to/image.png",
			ContentType: "image/png",
			Size:        1024,
			Url:         "http://localhost/path/to/image.png",
		}, nil)

	p := storage.NewStoragePlugin(storage.Options{
		Storage: sm,
	})

	si := p.Instance().(storage.StorageInstance)
	si.Create(context.Background(), storage.CreateOptions{})
}

func TestDestroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	sm := storagemock.NewMockStorage(ctrl)

	sm.
		EXPECT().
		Destroy(gomock.Any(), gomock.Any()).
		Return(nil)

	p := storage.NewStoragePlugin(storage.Options{
		Storage: sm,
	})

	si := p.Instance().(storage.StorageInstance)
	si.Destroy(context.Background(), "path/to/image.png")
}

func TestGenerateSignedURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	sm := storagemock.NewMockStorage(ctrl)

	sm.
		EXPECT().
		GenerateSignedURL(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("https://secure.link/path/to/image.png", nil)

	p := storage.NewStoragePlugin(storage.Options{
		Storage: sm,
	})

	si := p.Instance().(storage.StorageInstance)
	si.GenerateSignedURL(context.Background(), "path/to/image.png", time.Minute*30)
}
