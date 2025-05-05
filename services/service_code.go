package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/pkg/stdcopy"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func SaveUserCode(userID int, code string) (string, error) {
	dir := "/codes"

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

	hostCodePath := os.Getenv("HOST_CODE_PATH")
	if hostCodePath == "" {
		errChan <- fmt.Errorf("HOST_CODE_PATH is not set")
		return
	}

	fileName := filepath.Base(filePath)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",
		Cmd:   []string{"go", "run", "/tmp/codes/" + fileName},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostCodePath,
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

	// Fix 1: Use proper log reader handling
	logReader, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		errChan <- err
		return
	}
	defer logReader.Close()

	// Fix 2: Use the proper Docker log reader which handles headers
	output := make(chan string)
	go func() {
		// Use docker-specific stdout/stderr multiplexing reader
		stdcopy.StdCopy(os.Stdout, os.Stderr, logReader)

		// Collect logs separately for the channel
		logReader.Close()
		newLogReader, _ := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     false,
		})
		defer newLogReader.Close()

		// Use proper header handling
		var stdoutBuf, stderrBuf bytes.Buffer
		_, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, newLogReader)
		if err != nil {
			errChan <- err
			return
		}

		// Combine stdout and stderr
		var combinedOutput string
		if stdoutBuf.Len() > 0 {
			combinedOutput += stdoutBuf.String()
		}
		if stderrBuf.Len() > 0 {
			if len(combinedOutput) > 0 {
				combinedOutput += "\n"
			}
			combinedOutput += stderrBuf.String()
		}

		output <- combinedOutput
	}()

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			errChan <- err
			return
		}
	case <-statusCh:
		// Container finished
	}

	// Get the output
	result := <-output
	fmt.Println("Output:", result)
	log.Println("Output:", result)
	logs <- result

	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
		errChan <- err
		return
	}
}
