package mkvextract

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoChapter(t *testing.T) {
	info, err := ExtractMetadata("testdata/_test_no_chapter.mkv")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Mkv Tøst", info.Title)
	assert.Equal(t, 0, len(info.Chapters))
}

func TestSingle(t *testing.T) {
	info, err := ExtractMetadata("testdata/_test_single.mkv")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Mkv Tøst", info.Title)
	assert.Equal(t, 1, len(info.Chapters))
	assert.Equal(t, "Start", info.Chapters[0].Name)
	assert.Equal(t, time.Duration(0), info.Chapters[0].Start)
	assert.Equal(t, time.Duration(5*time.Second), info.Chapters[0].End)
}

func TestNoTitle(t *testing.T) {
	_, err := ExtractMetadata("testdata/_test_no_title.mkv")
	if err == nil {
		t.Error("No Error, for erroneous file")
	}
}

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
