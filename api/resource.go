package api

import (
	"encoding/json"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
	"github.com/jjmrocha/oblivion/router"
)

type Handler struct {
	service *bucket.BucketService
}

func NewHandler(bucketService *bucket.BucketService) *Handler {
	handler := Handler{
		service: bucketService,
	}

	return &handler
}

func (h *Handler) SetRoutes(mux *router.Multiplexer) {
	setBucketRoutes(mux, h)
	setKeyRoutes(mux, h)
}

func setBucketRoutes(mux *router.Multiplexer, h *Handler) {
	mux.GET("/v1/buckets", func(ctx *router.Context) (*router.Response, error) {
		bucketNames, err := h.service.BucketList()
		if err != nil {
			return nil, err
		}

		return ctx.OK(bucketNames)
	})

	mux.POST("/v1/buckets", func(ctx *router.Context) (*router.Response, error) {
		var request repo.Bucket

		err := json.NewDecoder(ctx.Request.Body).Decode(&request)
		if err != nil {
			return nil, apperror.New(apperror.BadRequestPaylod)
		}

		bucket, err := h.service.CreateBucket(request.Name, request.Schema)
		if err != nil {
			return nil, err
		}

		return ctx.Created(bucket)
	})

	mux.GET("/v1/buckets/{bucket}", func(ctx *router.Context) (*router.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		bucket, err := h.service.GetBucket(bucketName)
		if err != nil {
			return nil, err
		}

		return ctx.OK(bucket)
	})

	mux.DELETE("/v1/buckets/{bucket}", func(ctx *router.Context) (*router.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		err := h.service.DeleteBucket(bucketName)
		if err != nil {
			return nil, err
		}

		return ctx.NoContent()
	})
}

func setKeyRoutes(mux *router.Multiplexer, h *Handler) {
	mux.GET("/v1/buckets/{bucket}/keys/{key}", func(ctx *router.Context) (*router.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		key := ctx.Request.PathValue("key")

		value, err := h.service.Value(bucketName, key)
		if err != nil {
			return nil, err
		}

		return ctx.OK(value)
	})

	mux.PUT("/v1/buckets/{bucket}/keys/{key}", func(ctx *router.Context) (*router.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		key := ctx.Request.PathValue("key")

		var value model.Object

		err := json.NewDecoder(ctx.Request.Body).Decode(&value)
		if err != nil {
			return nil, apperror.New(apperror.BadRequestPaylod)
		}

		err = h.service.SetValue(bucketName, key, value)
		if err != nil {
			return nil, err
		}

		return ctx.NoContent()
	})

	mux.DELETE("/v1/buckets/{bucket}/keys/{key}", func(ctx *router.Context) (*router.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		key := ctx.Request.PathValue("key")

		err := h.service.DeleteValue(bucketName, key)
		if err != nil {
			return nil, err
		}

		return ctx.NoContent()
	})

	mux.GET("/v1/buckets/{bucket}/keys", func(ctx *router.Context) (*router.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		criteria := ctx.Request.URL.Query()

		keys, err := h.service.FindKeys(bucketName, criteria)
		if err != nil {
			return nil, err
		}

		return ctx.OK(keys)
	})
}
