package upload

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid"
)

func generateFileKey(handler *multipart.FileHeader, file multipart.File) string {
	uuid := shortuuid.New()
	fl := handler.Filename
	ext := filepath.Ext(fl)
	fileName := strings.TrimSuffix(fl, ext)
	return slug.Make(fileName) + "-" + uuid + ext
}
