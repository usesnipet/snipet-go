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
