package api

import (
	"encoding/json"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/httprouter"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
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

func (h *Handler) SetRoutes(router *httprouter.Router) {
	setBucketRoutes(router, h)
	setKeyRoutes(router, h)
}

func setBucketRoutes(router *httprouter.Router, h *Handler) {
	router.GET("/v1/buckets", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketNames, err := h.service.BucketList()
		if err != nil {
			return nil, err
		}

		return ctx.OK(bucketNames)
	})

	router.POST("/v1/buckets", func(ctx *httprouter.Context) (*httprouter.Response, error) {
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

	router.GET("/v1/buckets/{bucket}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		bucket, err := h.service.GetBucket(bucketName)
		if err != nil {
			return nil, err
		}

		return ctx.OK(bucket)
	})

	router.DELETE("/v1/buckets/{bucket}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		err := h.service.DeleteBucket(bucketName)
		if err != nil {
			return nil, err
		}

		return ctx.NoContent()
	})
}

func setKeyRoutes(router *httprouter.Router, h *Handler) {
	router.GET("/v1/buckets/{bucket}/keys/{key}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		key := ctx.Request.PathValue("key")

		value, err := h.service.Value(bucketName, key)
		if err != nil {
			return nil, err
		}

		return ctx.OK(value)
	})

	router.PUT("/v1/buckets/{bucket}/keys/{key}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
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

	router.DELETE("/v1/buckets/{bucket}/keys/{key}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		key := ctx.Request.PathValue("key")

		err := h.service.DeleteValue(bucketName, key)
		if err != nil {
			return nil, err
		}

		return ctx.NoContent()
	})

	router.GET("/v1/buckets/{bucket}/keys", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		criteria := ctx.Request.URL.Query()

		keys, err := h.service.FindKeys(bucketName, criteria)
		if err != nil {
			return nil, err
		}

		return ctx.OK(keys)
	})
}
