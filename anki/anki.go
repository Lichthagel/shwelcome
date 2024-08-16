package anki

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// check if line starts with # or is empty
		if len(scanner.Text()) == 0 || scanner.Text()[0] == '#' {
			continue
		}

		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

type AnkiCard struct {
	Word         string
	Reading      string
	Translations []string
}

func RemoveHTMLTags(s string) string {
	regexp := regexp.MustCompile("<[^>]*>")
	return regexp.ReplaceAllString(s, "")
}

func ParseLine(line string) (AnkiCard, error) {
	// split line by tab
	parts := strings.Split(line, "\t")

	// check if line has 3 parts
	if len(parts) < 4 {
		return AnkiCard{}, fmt.Errorf("line has %d parts, at least 4", len(parts))
	}

	translations := strings.Split(parts[3], "<br>")

	// trim & remove quotes
	for i, translation := range translations {
		translations[i] = strings.Trim(translation, "\"")
		translations[i] = strings.TrimSpace(translations[i])
	}

	return AnkiCard{
		Word:         parts[1],
		Reading:      parts[2],
		Translations: translations,
	}, nil
}

func RandomCard(path string) (AnkiCard, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return AnkiCard{}, err
	}

	// pick a random line
	idx := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(lines))
	line := lines[idx]

	return ParseLine(line)
}

func RenderTranslation(translation string) string {
	styleType := lipgloss.NewStyle().Foreground(lipgloss.Color(catppuccin.Mocha.Subtext0().Hex)).Italic(true)
	styleNumber := lipgloss.NewStyle().
		Background(lipgloss.Color(catppuccin.Mocha.Pink().Hex)).
		Foreground(lipgloss.Color(catppuccin.Mocha.Crust().Hex)).
		Width(3).
		Align(1).
		MarginRight(1)
	styleWord := lipgloss.NewStyle()

	if strings.Contains(translation, "78909C") {
		// word type
		translation = RemoveHTMLTags(translation)
		translation = strings.TrimSpace(translation)
		translation = styleType.Render(translation)
	} else {
		// translation
		split := strings.Split(translation, "</b>")

		if len(split) < 2 {
			translation = RemoveHTMLTags(translation)
			translation = strings.TrimSpace(translation)
			translation = styleWord.Render(translation)
		} else {
			// first element is the number of the translation
			number := split[0]
			number = RemoveHTMLTags(number)
			number = strings.TrimSpace(number)

			translation = strings.Join(split[1:], "")
			translation = RemoveHTMLTags(translation)
			translation = strings.TrimSpace(translation)

			translation = styleNumber.Render(number) + styleWord.Render(translation)
		}
	}

	return translation
}

func RenderCard(card AnkiCard) string {
	styleWord := lipgloss.NewStyle().Foreground(lipgloss.Color(catppuccin.Mocha.Pink().Hex))
	styleReading := lipgloss.NewStyle().Foreground(lipgloss.Color(catppuccin.Mocha.Subtext0().Hex))

	translations := ""
	for _, translation := range card.Translations {
		translations += RenderTranslation(translation) + "\n"
	}

	out := lipgloss.JoinVertical(0,
		styleWord.Render(card.Word)+" - "+styleReading.Render(card.Reading),
		translations,
	)

	// border := lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color(catppuccin.Mocha.Overlay2().Hex)).
	// 	BorderTop(true)

	return out
}
