package persistence

import "context"

type SnapshotRepository interface {
	Save(ctx context.Context, data map[string][]byte) error
	Load(ctx context.Context) (map[string][]byte, error)
}
