// Copyright (C) 2026 Podomy.
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/podomy/hive/src/logs"
)

func main() {
	logger, syncLogs, err := logs.Init()
	if err != nil {
		// logger has not been initialized here; this is the only case where log is acceptable.
		log.Fatal(err)
	}
	defer func() {
		if err := syncLogs(); err != nil {
			logger.Warn("log sync failed", zap.Error(err))
		}
	}()

	logger.Info("node runtime started")
}
