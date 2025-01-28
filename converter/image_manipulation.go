package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

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

type CompressionSettings struct {
	lossy      int
	colors     int
	dither     bool
	scale      int
	targetSize int64
}

const TARGET_SIZE = 200000

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

		if file.framecount > 59 {
			file = reduce_frames(file)
		}
		if file.size <= TARGET_SIZE {
			err := moveFile(filepath.Join("to-convert", file.filename), filepath.Join("converted", file.filename))
			if err != nil {
				panic(err)
			}
		} else {
			file = optimizeGIF(file)
			err := moveFile(filepath.Join("to-convert", file.filename), filepath.Join("converted", file.filename))
			if err != nil {
				panic(err)
			}
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

func compress_lossy(file File, settings CompressionSettings) File {
	//make lossy val and colors and dither setupable
	lossyarg := "--lossy=" + strconv.Itoa(settings.lossy)
	colorarg := "--colors=" + strconv.Itoa(settings.colors)
	ditherarg := ""
	if settings.dither {
		ditherarg = "--dither"
	}
	writeFile(file)
	executeCommand("gifsicle", "-b", file.filename, "--optimize=3", lossyarg, colorarg, ditherarg)
	return getFileSize(file)
}

func resize(file File, fac int) File {
	scalearg := "--scale=0." + strconv.Itoa(fac)
	executeCommand("gifsicle", "-b", file.filename, "--optimize=3", scalearg)
	return getFileSize(file)
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
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	file.gif_data = createGIF(file.frames)
	gif.EncodeAll(f, file.gif_data)

}

func getFileSize(file File) File {
	f, err := os.Open(filepath.Join("to-convert", file.filename))
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	file.size = fi.Size()
	fmt.Println(file.size/1024, "kb")
	return file
}

func optimizeGIF(file File) File {
	compressionSteps := []CompressionSettings{
		{lossy: 30, colors: 128, dither: true},
		{lossy: 50, colors: 64, dither: true},
		{lossy: 70, colors: 32, dither: true},
	}

	for _, settings := range compressionSteps {
		file = compress_lossy(file, settings)
		if file.size <= TARGET_SIZE {
			return file
		}
	}

	for file.size > TARGET_SIZE {
		file = resize(file, 8)
	}

	return file
}

func moveFile(sourcePath, destPath string) error {
	destPath = destPath[:len(destPath)-6] + ".gif"
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}
