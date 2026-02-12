package ai

import (
	"regexp"
	"strconv"
	"strings"
)

func parseAIResponse(content string) (comment string, score int, err error) {

	scoreRegex := regexp.MustCompile(`(?:分数|Score)[：:]\s*(\d+)`)
	matches := scoreRegex.FindStringSubmatch(content)

	if len(matches) > 1 {
		s, _ := strconv.Atoi(matches[1])
		if s >= 0 && s <= 100 {
			score = s
		}
	}

	lines := strings.Split(content, "\n")
	var comments []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" &&
			!strings.Contains(trimmed, "分数") &&
			!strings.Contains(trimmed, "Score") &&
			!strings.HasPrefix(trimmed, "评语：") {
			comments = append(comments, trimmed)
		} else if strings.HasPrefix(trimmed, "评语：") {
			comments = append(comments, strings.TrimPrefix(trimmed, "评语："))
		}
	}

	comment = strings.TrimSpace(strings.Join(comments, "\n"))
	if comment == "" {
		comment = "AI评价已生成"
	}

	return comment, score, nil
}
