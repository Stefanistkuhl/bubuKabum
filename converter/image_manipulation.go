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
	gif_data   *gif.GIF
	size       int64
	framecount int64
	frames     []Frame
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
	file.framecount = emote.metadata.frame_count
	if strings.HasSuffix(file.filename, ".gif") {

		gifData, err := os.ReadFile(filepath.Join("to-convert", file.filename))
		if err != nil {
			panic(err)
		}
		frames, err := splitGif(gifData)
		file.frames = frames
		if err != nil {
			panic(err)
		}
		file.gif_data = createGIF(frames)

		reduceSizeFuncs := []func(file File) File{compress_lossy, resize}
		i := 0
		if file.framecount > 59 {
			file = reduce_frames(file)
		}
		if file.size <= 200000 {
			//do nothing
		}

		fmt.Println("size:", file.size/1024, "kb")
		for {
			//make it reduce colors -> 512 -> 256 -> 128 then resize
			file = reduceSizeFuncs[i](file)
			i = 1 - i
			if file.size <= 200000 {
				fmt.Println("size:", file.size/1024, "kb")
				break
			}
			//only temp ther until i actually compress shit
			break
		}
	}
}

func reduce_frames(file File) File {
	buf := new(bytes.Buffer)
	frames := file.frames
	originalLength := len(frames)
	targetFrameCount := 59
	newFrames := make([]Frame, targetFrameCount)

	for i := 0; i < targetFrameCount; i++ {
		position := int(float64(i) * float64(originalLength-1) / float64(targetFrameCount-1))
		newFrames[i] = frames[position]
	}

	frames = newFrames
	file.gif_data = createGIF(frames)
	file.frames = frames
	err := gif.EncodeAll(buf, file.gif_data)
	if err != nil {
		panic(err)
	}
	file.size = int64(buf.Len())
	return file

}

func compress_lossy(file File) File {
	//make lossy val and colors and dither setupable
	writeFile(file)
	executeCommand("gifsicle", "-b", file.filename, "--optimize=3", "--lossy=100", "--colors=25", "--dither")
	return file
}

func resize(file File) File {
	return file
}

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
func createGIF(frames []Frame) *gif.GIF {
	g := &gif.GIF{
		Image:    make([]*image.Paletted, len(frames)),
		Delay:    make([]int, len(frames)),
		Disposal: make([]byte, len(frames)),
	}

	for i, frame := range frames {
		g.Image[i] = frame.Image
		g.Delay[i] = frame.Delay
		g.Disposal[i] = frame.Disposal
	}

	return g
}

func writeFile(file File) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(cwd, "to-convert", file.filename)
	fmt.Println(path)
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	file.gif_data = createGIF(file.frames)
	gif.EncodeAll(f, file.gif_data)

}
