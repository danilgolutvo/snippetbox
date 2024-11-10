package models

import (
	"github.com/stretchr/testify/assert"
	"snippetboxmod/internal/myassert"
	"testing"
)

func TestUserModelExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	tests := []struct {
		name   string
		UserID int
		want   bool
	}{
		{
			name:   "Valid ID",
			UserID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			UserID: 0,
			want:   false,
		},
		{
			name:   "non-existent ID",
			UserID: 2,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{DB: db}
			exists, err := m.Exists(tt.UserID)
			assert.Equal(t, exists, tt.want)
			myassert.NilError(t, err)
		})
	}
}
