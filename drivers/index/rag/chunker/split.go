package chunker

import (
	"strings"
	"unicode/utf8"
)

var defaultSeparators = []string{
	"\n\n",
	"\n",
	". ",
	"? ",
	"! ",
	"; ",
	", ",
	" ",
	"",
}

type textPart struct {
	Content string
	Start   int
	End     int
}

func splitTextParts(text string, cfg Config) []textPart {
	text = normalizeText(text, cfg.TrimWhitespace)
	if text == "" {
		return nil
	}

	contents := splitRecursive([]rune(text), cfg.ChunkSize, cfg.separators())
	if cfg.MinChunkSize > 0 {
		contents = mergeSmallChunks(contents, cfg.MinChunkSize, cfg.ChunkSize)
	}

	parts := make([]textPart, 0, len(contents))
	searchFrom := 0
	for _, content := range contents {
		if cfg.TrimWhitespace {
			content = normalizeText(content, true)
		}
		if content == "" {
			continue
		}

		start := findOffset(text, content, searchFrom)
		end := start + utf8.RuneCountInString(content)
		searchFrom = end

		parts = append(parts, textPart{
			Content: content,
			Start:   start,
			End:     end,
		})
	}

	if cfg.Overlap > 0 {
		for i := 1; i < len(parts); i++ {
			prefix := tailRunes(parts[i-1].Content, cfg.Overlap)
			parts[i].Content = prefix + parts[i].Content
		}
	}

	return parts
}

func normalizeText(text string, trim bool) string {
	if !trim {
		return text
	}
	return strings.TrimSpace(text)
}

func splitRecursive(text []rune, chunkSize int, separators []string) []string {
	if len(text) == 0 {
		return nil
	}
	if len(text) <= chunkSize {
		return []string{string(text)}
	}
	if len(separators) == 0 {
		return hardSplit(text, chunkSize)
	}

	separator, remaining := pickSeparator(text, separators)
	if separator == "" {
		return hardSplit(text, chunkSize)
	}

	parts := splitBySeparator(text, separator)
	var chunks []string
	var current []rune

	flush := func() {
		if len(current) == 0 {
			return
		}
		chunks = append(chunks, string(current))
		current = nil
	}

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		nextLen := len(current) + len(part)
		if nextLen <= chunkSize {
			current = append(current, part...)
			continue
		}

		if len(current) > 0 {
			flush()
		}

		if len(part) <= chunkSize {
			current = append(current[:0], part...)
			continue
		}

		chunks = append(chunks, splitRecursive(part, chunkSize, remaining)...)
	}

	flush()
	return chunks
}

func pickSeparator(text []rune, separators []string) (string, []string) {
	asString := string(text)
	for i, separator := range separators {
		if separator == "" {
			return "", separators[i+1:]
		}
		if strings.Contains(asString, separator) {
			return separator, separators[i+1:]
		}
	}
	return "", nil
}

func splitBySeparator(text []rune, separator string) [][]rune {
	sepRunes := []rune(separator)
	if len(sepRunes) == 0 {
		return [][]rune{text}
	}

	var parts [][]rune
	start := 0
	for i := 0; i <= len(text)-len(sepRunes); i++ {
		if !runeSliceEqual(text[i:i+len(sepRunes)], sepRunes) {
			continue
		}
		parts = append(parts, text[start:i])
		start = i + len(sepRunes)
		i += len(sepRunes) - 1
	}
	parts = append(parts, text[start:])
	return parts
}

func runeSliceEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func hardSplit(text []rune, chunkSize int) []string {
	var chunks []string
	for start := 0; start < len(text); start += chunkSize {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, string(text[start:end]))
	}
	return chunks
}

func mergeSmallChunks(chunks []string, minSize, maxSize int) []string {
	if len(chunks) <= 1 {
		return chunks
	}

	out := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		if len(out) == 0 {
			out = append(out, chunk)
			continue
		}

		last := out[len(out)-1]
		if utf8.RuneCountInString(chunk) >= minSize {
			out = append(out, chunk)
			continue
		}

		merged := last + chunk
		if utf8.RuneCountInString(merged) <= maxSize {
			out[len(out)-1] = merged
			continue
		}

		out = append(out, chunk)
	}

	return out
}

func tailRunes(text string, size int) string {
	runes := []rune(text)
	if len(runes) <= size {
		return text
	}
	return string(runes[len(runes)-size:])
}

func findOffset(source, chunk string, from int) int {
	sourceRunes := []rune(source)
	chunkRunes := []rune(chunk)
	if from >= len(sourceRunes) {
		return len(sourceRunes)
	}

	for i := from; i <= len(sourceRunes)-len(chunkRunes); i++ {
		if runeSliceEqual(sourceRunes[i:i+len(chunkRunes)], chunkRunes) {
			return i
		}
	}

	return from
}
