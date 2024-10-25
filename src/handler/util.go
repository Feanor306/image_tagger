package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// MAX_SIZE is the max size of tag or media that can be fetched
const MAX_SIZE = 100

// render is a helper that simplifies most render calls by handlers
func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

// handleFileUpload will save file in data/files with new name
// naming convention {media.Id}.{file_extension}
// data/files dir will be created if missing
func handleFileUpload(c echo.Context, id string) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fns := strings.Split(file.Filename, ".")
	destPath := filepath.Join(".", "data", "files", fmt.Sprintf("%s.%s", id, fns[len(fns)-1]))

	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		return err
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
