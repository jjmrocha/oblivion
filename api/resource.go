package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/future"
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
		c, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		bucketNames, err := future.Async(func() ([]string, error) {
			return h.service.BucketList(c)
		}).Await()
		if err != nil {
			return nil, err
		}

		return ctx.OK(bucketNames)
	})

	router.POST("/v1/buckets", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		var request externalBucket

		err := json.NewDecoder(ctx.Request.Body).Decode(&request)
		if err != nil {
			return nil, apperror.BadRequestPaylod.WithCause(err)
		}

		if err := valid.BucketName(request.Name); err != nil {
			return nil, err
		}

		if err := valid.Schema(request.Schema); err != nil {
			return nil, err
		}

		c, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		bucket, err := h.service.AsyncCreateBucket(c, request.Name, request.Schema).Await()
		if err != nil {
			return nil, err
		}

		response := createExternalBucket(bucket)

		return ctx.Created(response)
	})

	router.GET("/v1/buckets/{bucket}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		c, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		bucket, err := h.service.GetBucket(c, bucketName)
		if err != nil {
			return nil, err
		}

		response := createExternalBucket(bucket)

		return ctx.OK(response)
	})

	router.DELETE("/v1/buckets/{bucket}", func(ctx *httprouter.Context) (*httprouter.Response, error) {
		bucketName := ctx.Request.PathValue("bucket")

		if err := valid.BucketName(bucketName); err != nil {
			return nil, err
		}

		c, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		err := h.service.DeleteBucket(c, bucketName)
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

		c, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		value, err := h.service.Value(c, bucketName, key)
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
			return nil, apperror.BadRequestPaylod.WithCause(err)
		}

		c, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		err = h.service.SetValue(c, bucketName, key, value)
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

		c, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		err := h.service.DeleteValue(c, bucketName, key)
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

		c, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()

		keys, err := h.service.FindKeys(c, bucketName, criteria)
		if err != nil {
			return nil, err
		}

		return ctx.OK(keys)
	})
}
