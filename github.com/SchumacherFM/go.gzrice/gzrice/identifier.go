package main

import (
	"strconv"

	incremental "github.com/SchumacherFM/wanderlust/Godeps/_workspace/src/github.com/GeertJohan/go.incremental"
)

var identifierCount incremental.Uint64

func nextIdentifier() string {
	num := identifierCount.Next()
	return strconv.FormatUint(num, 36) // 0123456789abcdefghijklmnopqrstuvwxyz
}
