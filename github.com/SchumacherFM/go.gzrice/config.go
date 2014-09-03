package gzrice

import "path"

// LocateMethod defines how a box is located.
type LocateMethod int

const (
	LocateFS       = LocateMethod(iota) // Locate on the filesystem.
	LocateEmbedded                      // Locate embedded boxes.
)

var compressFileExt map[string]bool

func init() {
	compressFileExt = map[string]bool{
		"css": true,
		"js":  true,
		"eot": true,
		"svg": true,
	}
}

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

func IsCompressingAllowed(filePath string) bool {
	ext := path.Ext(filePath) // returns .ext
	if len(ext) < 2 {
		return false
	}
	_, isSet := compressFileExt[ext[1:]]
	return isSet
}
