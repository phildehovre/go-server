package main

import "io"

type FileSystemPlayerStore struct {
	database io.Reader
}
