package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// give better name
type File struct {
	filename   string
	size       int64
	framecount int64
}
type Frame struct {
	Image    *image.Paletted
	Delay    int
	Disposal byte
}

func compress_emote(emote Emote) {
	var file File
	fmt.Println("download done now in compress_emote() with", emote.emote_name)
	file.filename = emote.emote_name + emote.metadata.filename
	file.size = emote.metadata.size
	file.size = emote.metadata.frame_count
	if strings.HasSuffix(file.filename, ".gif") {
		reduceSizeFuncs := []func(file File){compress_lossy, resize}
		reduce_frames(file)
		i := 0
		if file.framecount > 59 {
			reduce_frames(file)
		}
		if file.size <= 200 {
			//do nothing
		}
		for {
			reduceSizeFuncs[i](file)
			i = 1 - i
			if file.size <= 200 {
				break
			}
			//only temp ther until i actually compress shit
			break
		}
	}
}

func reduce_frames(file File) {
	gifData, err := os.ReadFile(filepath.Join("to-convert", file.filename))
	if err != nil {
		panic(err)
	}
	frames, err := splitGif(gifData)
	if err != nil {
		panic(err)
	}

	if len(frames) <= 59 {
		return
	}

	originalLength := len(frames)
	targetFrameCount := 59
	newFrames := make([]Frame, targetFrameCount)

	for i := 0; i < targetFrameCount; i++ {
		position := int(float64(i) * float64(originalLength-1) / float64(targetFrameCount-1))
		newFrames[i] = frames[position]
	}
	frames = newFrames
}

func compress_lossy(file File) {}

func resize(file File) {}

func executeCommand(command string, args ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	cmd := exec.Command(command, args...)
	toconvertdir := filepath.Join(cwd, "to-convert")
	cmd.Dir = toconvertdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	return nil
}

func splitGif(gifData []byte) ([]Frame, error) {
	reader := bytes.NewReader(gifData)
	gif, err := gif.DecodeAll(reader)
	if err != nil {
		return nil, err
	}

	frames := make([]Frame, len(gif.Image))
	for i := range gif.Image {
		frames[i] = Frame{
			Image:    gif.Image[i],
			Delay:    gif.Delay[i],
			Disposal: gif.Disposal[i],
		}
	}
	return frames, nil
}

