package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func SaveUserCode(userID int, code string) (string, error) {
	dir := "./codes"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("could not create code directory: %w", err)
	}

	filePath := fmt.Sprintf("%s/%d.go", dir, userID)

	if err := os.WriteFile(filePath, []byte(code), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to save code to file: %w", err)
	}

	return filePath, nil
}

func StartUserCode(ctx context.Context, logs chan string, errChan chan error, filePath string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		errChan <- err
		return
	}
	defer os.Remove(filePath)

	// Получаем абсолютный путь к файлу
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		errChan <- fmt.Errorf("failed to get absolute path: %w", err)
		return
	}

	// Получаем абсолютный путь к директории с файлом
	sourceDir := filepath.Dir(absFilePath)

	// Абсолютный путь назначения внутри контейнера
	targetPath := "/tmp/codes/" + filepath.Base(filePath)

	// Создаем контейнер
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",                   // Используем golang:latest
		Cmd:   []string{"go", "run", targetPath}, // Путь внутри контейнера
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourceDir,    // Абсолютный путь на хосте
				Target: "/tmp/codes", // Абсолютный путь внутри контейнера
			},
		},
	}, nil, nil, "")
	if err != nil {
		errChan <- err
		return
	}

	// Запускаем контейнер
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		errChan <- err
		return
	}

	// Читаем логи контейнера
	logReader, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		errChan <- err
		return
	}
	defer logReader.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, logReader); err != nil {
		errChan <- err
		return
	}

	logs <- buf.String()

	// Удаляем контейнер после выполнения
	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
		errChan <- err
		return
	}
}
