package upload

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid"
)

func generateFileKey(handler *multipart.FileHeader, file multipart.File, mime *mimetype.MIME) string {
	ext := mime.Extension()
	uuid := shortuuid.New()
	fl := handler.Filename
	fileName := strings.TrimSuffix(fl, filepath.Ext(fl))
	return slug.Make(fileName) + "-" + uuid + ext
}
