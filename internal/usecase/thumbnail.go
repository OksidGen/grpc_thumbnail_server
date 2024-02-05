package usecase

import (
	"fmt"
	"github.com/OksidGen/grpc_thumbnail/server/internal/repository"
	"io"
	"net/http"
	"strings"
)

type ThumbnailUsecase interface {
	GetThumbnail(videoURL string) ([]byte, error)
}
type thumbnailUsecase struct {
	repo repository.ThumbnailRepository
}

func NewThumbnailUsecase(repo repository.ThumbnailRepository) ThumbnailUsecase {
	return &thumbnailUsecase{repo: repo}
}

func (u *thumbnailUsecase) GetThumbnail(videoURL string) ([]byte, error) {
	cachedThumbnail, err := u.repo.GetThumbnail(videoURL)
	if err == nil {
		return cachedThumbnail, nil
	}

	thumbnailURL, err := getThumbnailURL(videoURL)
	if err != nil {
		return nil, err
	}

	thumbnailData, err := downloadThumbnail(thumbnailURL)
	if err != nil {
		return nil, err
	}

	err = u.repo.SaveThumbnail(videoURL, thumbnailData)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil
}
func getThumbnailURL(videoURL string) (string, error) {
	htmlContent, err := getHTMLContent(videoURL)
	if err != nil {
		return "", err
	}
	thumbnailURL, err := extractThumbnailURL(htmlContent)
	if err != nil {
		return "", err
	}
	return thumbnailURL, nil
}

func downloadThumbnail(thumbnailURL string) ([]byte, error) {
	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	thumbnailData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return thumbnailData, nil
}

func getHTMLContent(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	htmlBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(htmlBytes), nil
}

func extractThumbnailURL(htmlContent string) (string, error) {
	startTag := "<meta property=\"og:image\" content=\""
	endTag := "\">"

	startIndex := findSubstring(htmlContent, startTag)
	if startIndex == -1 {
		return "", fmt.Errorf("start tag not found")
	}

	startIndex += len(startTag)

	endIndex := strings.Index(htmlContent[startIndex:], endTag)
	if endIndex == -1 {
		return "", fmt.Errorf("end tag not found")
	}

	thumbnailURL := htmlContent[startIndex : startIndex+endIndex]

	return thumbnailURL, nil
}

func findSubstring(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
