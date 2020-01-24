package commands

import (
	"errors"
	"strings"
)

var (
	VidExts []string = []string{"mp4", "mov", "avi", "gif"}
	ImgExts []string = []string{"bmp", "dib", "jpeg", "jpg", "jpe", "jp2", "png", "pbm", "pgm", "ppm", "sr", "ras", "tiff", "tif"}
)

func CheckExt(fileExt string, exts []string) bool {
	for _, ext := range exts {
		if fileExt == ext {
			return true
		}
	}
	return false
}

func DetermineSubcommand(args []string) (Command, error) {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			if strings.HasPrefix(arg, "http") {
				cmd, _ := CommandMap["url"]
				return cmd, nil
			}
			fileExtension := strings.Split(arg, ".")[len(strings.Split(arg, "."))-1]
			if CheckExt(fileExtension, VidExts) {
				cmd, _ := CommandMap["video"]
				return cmd, nil
			} else if CheckExt(fileExtension, ImgExts) {
				cmd, _ := CommandMap["image"]
				return cmd, nil
			} else if len(strings.Split(arg, ".")) != 1 {
				return nil, errors.New("Invalid Format")
			} else {
				return nil, errors.New("Unknown Subcommand")
			}
		}
	}
	return nil, errors.New("Impossible")
}
