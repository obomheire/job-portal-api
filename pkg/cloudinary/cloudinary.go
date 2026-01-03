package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Service struct {
	cld *cloudinary.Cloudinary
}

func NewService() (*Service, error) {
	url := os.Getenv("CLOUDINARY_URL")
	if url == "" {
		return nil, fmt.Errorf("CLOUDINARY_URL is not set")
	}
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}
	return &Service{cld: cld}, nil
}

func (s *Service) UploadImage(ctx context.Context, file multipart.File, filename string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "job-portal",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
