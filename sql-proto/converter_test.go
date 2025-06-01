package sqlproto

import (
	"context"
	"testing"
)

func TestConverter(t *testing.T) {
	ctx := context.Background()

	srv := "postgres"
	dsn := "postgres://postgres:*Abcd123456@localhost:5432/example?sslmode=disable"
	outputPath := "./api/protos"
	_ = Convert(ctx, &srv, &dsn, &outputPath, nil, nil)
}
