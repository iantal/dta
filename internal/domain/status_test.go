package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := DownloadSuccess
	r := s.String()

	assert.Equal(t, r, "DOWNLOAD_SUCCESS")
}
