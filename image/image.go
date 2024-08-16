package image

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func imageBase64(image_path string) (string, error) {
	image_bytes, err := os.ReadFile(image_path)
	if err != nil {
		return "", err
	}

	image_base64 := base64.StdEncoding.EncodeToString(image_bytes)

	return image_base64, nil
}

func imageCodeITerm2(imageBase64 string, width uint, height uint, doNotMoveCursor bool) (string, error) {
	image_code := "\x1b]1337;File=inline=1"

	if width > 0 {
		image_code += fmt.Sprintf(";width=%d", width)
	}

	if height > 0 {
		image_code += fmt.Sprintf(";height=%d", height)
	}

	if doNotMoveCursor {
		image_code += ";doNotMoveCursor=1"
	}

	image_code += ":"

	image_code += imageBase64

	image_code += "\a"

	return image_code, nil
}

func imageCodeKitty(imageBase64 string, width uint, height uint, doNotMoveCursor bool) (string, error) {
	image_code := "\x1b_Ga=T,f=100" // open

	if width > 0 {
		image_code += fmt.Sprintf(",c=%d", width)
	}

	if height > 0 {
		image_code += fmt.Sprintf(",r=%d", height)
	}

	if doNotMoveCursor {
		image_code += ",C=1"
	}

	image_code += ";"

	image_code += imageBase64

	image_code += "\x1b\\" // close

	return image_code, nil
}

func padImage(imageCode string, width uint, height uint, returnCursor bool) string {
	image_code_padded := ""

	// force the terminal to scroll
	// image_code_padded += strings.Repeat("\n", height)
	// image_code_padded += fmt.Sprintf("%v[%vA", EscapeChar, height)

	// add image code
	if returnCursor {
		// ! SOMEWHAT HACKY:
		// ! - ensure that enough lines are on screen to render image by inerting lines using `<ESC>[%dL`
		// ! - save cursor position using `<ESC>[s`
		// ! - render image
		// ! - restore cursor position using `<ESC>[u`
		// !   - position on screen and not position in content -> scrolling messes up
		// image_code_padded += fmt.Sprintf("%v[%dL", EscapeChar, height)
		image_code_padded += ansi.SaveCursorPosition
		image_code_padded += imageCode
		image_code_padded += ansi.RestoreCursorPosition
	} else {
		image_code_padded += imageCode
	}

	// add spaces to fill width and height
	empty_line := strings.Repeat(" ", int(width)) + "\n"
	empty_block := strings.Repeat(empty_line, int(height))
	image_code_padded += empty_block

	return image_code_padded
}

func scrollDown(lines int) string {
	return fmt.Sprintf("%v%v", strings.Repeat("\n", lines), ansi.CursorUp(lines))
}

func PathToImgBlock(imagePath string, width uint, height uint) (string, error) {
	image_base64, err := imageBase64(imagePath)
	if err != nil {
		return "", err
	}

	image_code, err := imageCodeITerm2(image_base64, width, height, true)
	if err != nil {
		return "", err
	}

	image_code_padded := padImage(image_code, width, height, false)

	return image_code_padded, nil
}
