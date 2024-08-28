package main

import (
	"io"
	"os"
)

/*
This code defines a struct named 'tape' that wraps an io.ReadWriteSeeker.
The custom Write method ensures that before any write operation occurs,
the file's position is reset to the beginning (using Seek).
This behavior simulates overwriting the content from the start of the file
each time Write is called, similar to how a tape would operate.
*/
type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
