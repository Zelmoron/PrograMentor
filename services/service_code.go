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

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		errChan <- fmt.Errorf("failed to get absolute path: %w", err)
		return
	}

	sourceDir := filepath.Dir(absFilePath)

	targetPath := "/tmp/codes/" + filepath.Base(filePath)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",
		Cmd:   []string{"go", "run", targetPath},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourceDir,
				Target: "/tmp/codes",
			},
		},
	}, nil, nil, "")
	if err != nil {
		errChan <- err
		return
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		errChan <- err
		return
	}

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
	fmt.Println(buf.String())
	logs <- buf.String()

	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
		errChan <- err
		return
	}
}
