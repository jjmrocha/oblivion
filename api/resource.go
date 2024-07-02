package api

import (
	"encoding/json"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/httprouter"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/valid"
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
		var request bucketRepresentation

		err := json.NewDecoder(ctx.Request.Body).Decode(&request)
		if err != nil {
			return nil, apperror.BadRequestPaylod.NewErrorWithCause(err)
		}

		if err := valid.BucketName(request.Name); err != nil {
			return nil, err
		}

		if err := valid.Schema(request.Schema); err != nil {
			return nil, err
		}

		bucket, err := h.service.CreateBucket(request.Name, request.Schema)
		if err != nil {
			return nil, err
		}

		response := createBucketRepresentation(bucket)

		return ctx.Created(response)
	})

	router.GET("/v1/buckets/{bucket}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		bucket, err := h.service.GetBucket(bucketName)
		if err != nil {
			return nil, err
		}

		response := createBucketRepresentation(bucket)

		return ctx.OK(response)
	})

	router.DELETE("/v1/buckets/{bucket}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

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

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		if err := valid.Key(key); err != nil {
			return nil, err
		}

		value, err := h.service.Value(bucketName, key)
		if err != nil {
			return nil, err
		}

		return ctx.OK(value)
	})

	router.PUT("/v1/buckets/{bucket}/keys/{key}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		key := ctx.Request.PathValue("key")

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		if err := valid.Key(key); err != nil {
			return nil, err
		}

		var value model.Object

		err := json.NewDecoder(ctx.Request.Body).Decode(&value)
		if err != nil {
			return nil, apperror.BadRequestPaylod.NewErrorWithCause(err)
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

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		if err := valid.Key(key); err != nil {
			return nil, err
		}

		err := h.service.DeleteValue(bucketName, key)
		if err != nil {
			return nil, err
		}

		return ctx.NoContent()
	})

	router.GET("/v1/buckets/{bucket}/keys", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")
		criteria := ctx.Request.URL.Query()

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		keys, err := h.service.FindKeys(bucketName, criteria)
		if err != nil {
			return nil, err
		}

		return ctx.OK(keys)
	})
}
