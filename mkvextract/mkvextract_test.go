package mkvextract

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	info, err := ExtractMetadata("testdata/_test.mkv")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Mkv Tøst", info.Title)
	assert.Equal(t, 2, len(info.Chapters))
	assert.Equal(t, "Start", info.Chapters[0].Name)
	assert.Equal(t, time.Duration(0), info.Chapters[0].Start)
	assert.Equal(t, time.Duration(7*time.Second), info.Chapters[0].End)
	assert.Equal(t, "Jümp", info.Chapters[1].Name)
	assert.Equal(t, time.Duration(7*time.Second), info.Chapters[1].Start)
	assert.Equal(t, time.Duration(5*time.Second), info.Chapters[1].End)
}
