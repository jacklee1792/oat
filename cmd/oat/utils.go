package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// makeOutputPath creates an output path from the input path with the given suffix. For example,
// calling makeOutputPath with suffix "filtered" and inputPath "spec.json" will produce
// "spec-filtered.json" or "spec-filtered-n.json" for n as high as required to create a new file.
func makeOutputPath(in string, suffix string) string {
	ext := filepath.Ext(in)
	base := strings.TrimSuffix(in, ext)
	out := fmt.Sprintf("%s-%s%s", base, suffix, ext)
	for n := 0; ; n++ {
		if n >= 1 {
			out = fmt.Sprintf("%s-%s-%d%s", base, suffix, n, ext)
		}
		_, err := os.Stat(out)
		if errors.Is(err, os.ErrNotExist) {
			break
		}
		cobra.CheckErr(err)
	}
	return out
}
