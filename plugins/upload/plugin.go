package upload

import (
	"errors"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/upload/storage"
	"github.com/go-zepto/zepto/web"
)

type UploadPlugin struct {
	instance    UploadInstance
	maxFileSize int64
}

type Options struct {
	// Storage is the provider where the files will be stored
	Storage storage.Storage
	// MaxFileSize in bytes. e.g, upload.Mb(10) or upload.Kb(500)
	MaxFileSize int64
}

func NewUploadPlugin(opts Options) *UploadPlugin {
	maxFileSize := opts.MaxFileSize
	if maxFileSize == 0 {
		maxFileSize = Mb(20)
	}
	return &UploadPlugin{
		instance: &defaultUploadInstance{
			opts.Storage,
		},
		maxFileSize: maxFileSize,
	}
}

func (u *UploadPlugin) Name() string {
	return "upload"
}

func (u *UploadPlugin) Instance() interface{} {
	return u.instance
}

func (u *UploadPlugin) PrependMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (u *UploadPlugin) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

// Kb to bytes
func Kb(size int64) int64 {
	return size << 10
}

// Mb to bytes
func Mb(size int64) int64 {
	return size << 20
}

func maxBodyMiddleware(maxFileSize int64) web.MiddlewareFunc {
	return func(next web.RouteHandler) web.RouteHandler {
		return func(ctx web.Context) error {
			r := ctx.Request()
			w := ctx.Response()
			r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
			if err := r.ParseForm(); err != nil {
				return errors.New(http.StatusText(http.StatusRequestEntityTooLarge))
			}
			return next(ctx)
		}
	}
}

func (u *UploadPlugin) OnCreated(z *zepto.Zepto) {
	uploadCtrl := func(ctx web.Context) error {
		r := ctx.Request()
		r.ParseMultipartForm(u.maxFileSize)
		file, handler, err := r.FormFile("file")
		if err != nil {
			return err
		}
		defer file.Close()
		mime, err := mimetype.DetectReader(file)
		if err != nil {
			return err
		}
		key := generateFileKey(handler, file, mime)
		res, err := u.instance.UploadFile(ctx, storage.UploadFileOptions{
			AccessType:  storage.Private,
			Key:         key,
			ContentType: mime.String(),
			Body:        file,
		})
		return ctx.RenderJson(res)
	}
	r := z.Router("/upload")
	r.Post("/", maxBodyMiddleware(u.maxFileSize)(uploadCtrl))
}

func (u *UploadPlugin) OnZeptoInit(z *zepto.Zepto) {}

func (u *UploadPlugin) OnZeptoStart(z *zepto.Zepto) {}

func (u *UploadPlugin) OnZeptoStop(z *zepto.Zepto) {}
