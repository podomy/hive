// Copyright (C) 2026 Podomy.
// SPDX-License-Identifier: AGPL-3.0-or-later

package journalview

import (
	"context"

	"github.com/podomy/hive/src/journal"
)

type View interface {
	Apply(ctx context.Context, event journal.Event) error
	Rebuild(ctx context.Context, jr journal.Reader) error
}
