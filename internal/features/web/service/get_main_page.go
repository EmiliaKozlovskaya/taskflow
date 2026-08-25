package web_service

import (
	"fmt"
	"os"
	"path"
)

func (s *WebService) GetMainPage() ([]byte, error) {
	//1. htmlFilePath
	htmlFilePath := path.Join(
		os.Getenv("PROJECT_ROOT"),
		"/public/index.html",
	)
	//2. html, err := repository.Get(htmlFilePath)
	html, err := s.webRepository.GetFile(htmlFilePath)
	if err != nil {
		return nil, fmt.Errorf(
			"get file from repository: %w",
			err,
		)
	}
	//3. return html
	return html, nil
}
