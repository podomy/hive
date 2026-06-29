// Copyright (C) 2026 Podomy.
// SPDX-License-Identifier: AGPL-3.0-or-later

package journalview

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	"github.com/podomy/hive/src/journal"
	"github.com/podomy/hive/src/kvstore"
)

func TestEventsByIDGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	kv := testKVStore(t)
	view := NewEventsByID(kv)

	event := journal.NewEvent(uuid.New(), "node started", json.RawMessage(`{}`))

	if err := view.Apply(ctx, event); err != nil {
		t.Fatalf("apply event: %v", err)
	}

	got, err := view.Get(ctx, event.ID)
	if err != nil {
		t.Fatalf("get event: %v", err)
	}
	if got == nil {
		t.Fatalf("expected event got nil")
	}
	if got.ID != event.ID {
		t.Fatalf("expected event ID %s, got %s", event.ID, got.ID)
	}
}

func testKVStore(t *testing.T) *kvstore.KVStore {
	t.Helper()

	kv, err := kvstore.OpenDBPath(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db path: %v", err)
	}

	t.Cleanup(func() {
		if err := kv.Close(); err != nil {
			t.Fatalf("test close db: %v", err)
		}
	})

	return kv
}
