package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexing(t *testing.T) {
	resultChannel := make(chan ResultMessage)
	go proccessMkv("testdata/_test.mkv", resultChannel)
	result := <-resultChannel
	assert.NotEqual(t, Error, result.Result)
	go proccessMkv("testdata/_test.mkv", resultChannel)
	result = <-resultChannel
	assert.Equal(t, AlreadyIndexed, result.Result)

}
