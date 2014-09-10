package gzrice

import (
	"path"
)

// LocateMethod defines how a box is located.
type LocateMethod int

const (
	LocateFS = LocateMethod(iota) // Locate on the filesystem.
	LocateEmbedded                // Locate embedded boxes.
)

var (
	// 1 yes
	// -1 no
	// 0 not found so could be a html page without extension
	compressFileExt = map[string]int{
	"htm":  1,
	"html":  1,
	"css":  1,
	"js":   1,
	"eot":  1,
	"svg":  1,
	"png":  -1,
	"gif":  -1,
	"jpg":  -1,
	"ttf":  -1,
	"woff": -1,
	"ico":    -1,
}
)

// Config allows customizing the box lookup behavior.
type Config struct {
	// LocateOrder defines the priority order that boxes are searched for. By
	// default, the package global FindBox searches for embedded boxes first,
	// and then finally boxes on the filesystem.  That
	// search order may be customized by provided the ordered list here. Leaving
	// out a particular method will omit that from the search space. For
	// example, []LocateMethod{LocateEmbedded} will never search
	// the filesystem for boxes.
	LocateOrder []LocateMethod
}

// FindBox searches for boxes using the LocateOrder of the config.
func (c *Config) FindBox(boxName string) (*Box, error) {
	return findBox(boxName, c.LocateOrder)
}

func IsCompressingAllowed(filePath string) int {
	ext := path.Ext(filePath) // returns .ext
	if len(ext) < 2 {
		return 0
	}
	status, isSet := compressFileExt[ext[1:]]
	if true == isSet {
		return status
	}
	return 0
}
