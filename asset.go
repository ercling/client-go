package sanity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sanity-io/client-go/api"
	"github.com/sanity-io/client-go/internal/requests"
)

type AssetType string

const (
	AssetTypeImage AssetType = "image"
	AssetTypeFile  AssetType = "file"
)

type UploadAssetOption func(req *requests.Request)

// WithContentType is an UploadAssetOption to set the content type of the upload.
func WithContentType(contentType string) UploadAssetOption {
	return func(req *requests.Request) {
		req.Header("Content-Type", contentType)
	}
}

func WithFileName(filename string) UploadAssetOption {
	return func(req *requests.Request) {
		req.Param("filename", filename)
	}
}

func (c *Client) Asset() *AssetBuilder {
	return &AssetBuilder{
		c: c,
	}
}

type AssetBuilder struct {
	c   *Client
	err error
	tag string
}

// Upload uploads the asset data. For the api reference see: https://www.sanity.io/docs/assets
func (ab *AssetBuilder) Upload(ctx context.Context, assetType AssetType, data []byte, opts ...UploadAssetOption) (*api.AssetUploadResponse, error) {
	if ab.err != nil {
		return nil, fmt.Errorf("asset builder: %w", ab.err)
	}

	assetEndpoint := "images"
	if assetType == AssetTypeFile {
		assetEndpoint = "files"
	}

	req := ab.c.newAPIRequest().
		Method(http.MethodPost).
		AppendPath("assets").
		AppendPath(assetEndpoint).
		AppendPath(ab.c.dataset).
		Body(data).
		Tag(ab.tag, ab.c.tag)

	// Apply all provided options (including setting the Content-Type)
	for _, opt := range opts {
		opt(req)
	}

	var resp api.AssetUploadResponse
	var respItem api.MutateResultItem

	if _, err := ab.c.do(ctx, req, &respItem); err != nil {
		return nil, fmt.Errorf("asset: %w", err)
	}

	err := respItem.Unmarshal(&resp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal: %w", err)
	}

	return &resp, nil
}
