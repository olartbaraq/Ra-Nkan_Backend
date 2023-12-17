package all_test

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/olartbaraq/spectrumshelf/api"
)

func TestAuthenticatedMiddleware(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := api.AuthenticatedMiddleware(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticatedMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
