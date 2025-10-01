package search

import (
	"html"
	"regexp"
	"strings"
)

// extractMovieDetails ported logic from your Python code
func ExtractMovieDetails(caption string) (movieName, year, quality string) {
	movieName = "Unknown"
	year = "Unknown"
	quality = "Unknown"

	if caption == "" {
		return
	}

	namePatterns := []string{
		`^([^.]*?)\s\d{4}`,
		`^([^.]*?)\s\(\d{4}\)`,
		`^([^.]*?)\s\d{3,4}p`,
		`^([^.]*?)(?=\s\d{4})`,
	}

	yearPatterns := []string{
		`\((\d{4})\)`,
		`(\d{4})`,
		`\d{2}(\d{2})\b`,
		`(\d{4})\b`,
		`(19\d{2}|20\d{2})`,
		`(\d{2})(\d{2})\b`,
		`\b(20[012]\d|19[5-9]\d)\b`,
		`Movie :-\s(.*?)\(\d{4}\)`,
	}

	qualityPatterns := []string{
		`\d{3,4}p`,
		`\b(4|8|10)K\b`,
		`\b(FHD|HD|SD)\b`,
	}

	for _, pat := range namePatterns {
		re := regexp.MustCompile(pat)
		if m := re.FindStringSubmatch(caption); len(m) > 1 {
			movieName = strings.ReplaceAll(m[1], ".", "")
			break
		}
	}

	for _, pat := range yearPatterns {
		re := regexp.MustCompile(pat)
		if m := re.FindStringSubmatch(caption); len(m) > 1 {
			year = m[1]
			break
		}
	}

	for _, pat := range qualityPatterns {
		re := regexp.MustCompile(pat)
		if m := re.FindStringSubmatch(caption); len(m) > 0 {
			quality = m[0]
			break
		}
	}

	// escape HTML for safe formatting
	movieName = html.EscapeString(strings.TrimSpace(movieName))
	year = html.EscapeString(strings.TrimSpace(year))
	quality = html.EscapeString(strings.TrimSpace(quality))

	return
}
