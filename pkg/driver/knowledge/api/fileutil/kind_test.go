package fileutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/usesnipet/snipet/pkg/driver/knowledge"
)

func TestMapKind(t *testing.T) {
	tests := []struct {
		mediaType string
		want      knowledge.SourceItemKind
	}{
		{"", knowledge.SourceItemKindUnknown},
		{"text/plain", knowledge.SourceItemKindText},
		{"application/pdf", knowledge.SourceItemKindDocument},
		{"application/vnd.openxmlformats-officedocument.wordprocessingml.document", knowledge.SourceItemKindDocument},
		{"image/png", knowledge.SourceItemKindImage},
		{"image/svg+xml", knowledge.SourceItemKindImage},
		{"audio/mpeg", knowledge.SourceItemKindAudio},
		{"video/mp4", knowledge.SourceItemKindVideo},
		{"application/json", knowledge.SourceItemKindStructured},
		{"application/octet-stream", knowledge.SourceItemKindUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.mediaType, func(t *testing.T) {
			assert.Equal(t, tt.want, MapKind(tt.mediaType))
		})
	}
}

func TestMapKindEveryListedMediaTypeResolves(t *testing.T) {
	groups := []struct {
		mediaTypes []string
		want       knowledge.SourceItemKind
	}{
		{textMediaTypes, knowledge.SourceItemKindText},
		{documentMediaTypes, knowledge.SourceItemKindDocument},
		{imageMediaTypes, knowledge.SourceItemKindImage},
		{audioMediaTypes, knowledge.SourceItemKindAudio},
		{videoMediaTypes, knowledge.SourceItemKindVideo},
		{structuredMediaTypes, knowledge.SourceItemKindStructured},
	}

	for _, g := range groups {
		for _, mediaType := range g.mediaTypes {
			assert.Equalf(t, g.want, MapKind(mediaType), "mediaType %q", mediaType)
		}
	}
}

func TestMapKindHasNoCrossKindCollisions(t *testing.T) {
	groups := []struct {
		name       string
		mediaTypes []string
		kind       knowledge.SourceItemKind
	}{
		{"text", textMediaTypes, knowledge.SourceItemKindText},
		{"document", documentMediaTypes, knowledge.SourceItemKindDocument},
		{"image", imageMediaTypes, knowledge.SourceItemKindImage},
		{"audio", audioMediaTypes, knowledge.SourceItemKindAudio},
		{"video", videoMediaTypes, knowledge.SourceItemKindVideo},
		{"structured", structuredMediaTypes, knowledge.SourceItemKindStructured},
	}

	seen := make(map[string]string)
	for _, g := range groups {
		for _, mediaType := range g.mediaTypes {
			if owner, ok := seen[mediaType]; ok {
				assert.Failf(t, "cross-kind collision", "media type %q claimed by both %q and %q", mediaType, owner, g.name)
				continue
			}
			seen[mediaType] = g.name
		}
	}
}
