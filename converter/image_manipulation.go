package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// give better name
type File struct {
	filename   string
	size       int64
	framecount int64
}

func compress_emote(emote Emote) {
	var file File
	file.filename = emote.emote_name + emote.metadata.filename
	fmt.Println(file.filename)
	file.size = emote.metadata.size
	file.size = emote.metadata.frame_count
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
	}
}

func reduce_frames(file File) {
	//get rid of this and use the image libary instead
	err := executeCommand("gifsicle", "--explode", file.filename)
	if err != nil {
		fmt.Println(err)
	}
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
