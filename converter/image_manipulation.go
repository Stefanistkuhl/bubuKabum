package main

import (
	"fmt"
)

//give better name
type File struct {
	filename string
	size int64
	framecount 
}
func compress_emote(emote Emote) {
	var file File
	file.filename = emote.metadata.filename
	file.size = emote.metadata.size
	file.size = emote.metadata.frame_count

	if (file.framecount > 59){
		reduce_frames()
	}

	for {
		reduce_size = 
	}
}

func reduce_frames() {}

func compress_lossy() {}

func resize() {}
