// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"context"
	"strings"

	"github.com/kaleido-io/firefly/internal/config"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	rootLogger = logrus.NewEntry(logrus.StandardLogger())

	// L accesses the current logger from the context
	L = loggerFromContext
)

type (
	ctxLogKey struct{}
)

// WithLogger adds the specified logger to the context
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, ctxLogKey{}, logger)
}

// WithLogField adds the specified field to the logger in the context
func WithLogField(ctx context.Context, key, value string) context.Context {
	return WithLogger(ctx, loggerFromContext(ctx).WithField(key, value))
}

// LoggerFromContext returns the logger for the current context, or no logger if there is no context
func loggerFromContext(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(ctxLogKey{})
	if logger == nil {
		return rootLogger
	}
	return logger.(*logrus.Entry)
}

// SetupLogging initializes logging
func SetupLogging(ctx context.Context) {
	logrus.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   !config.GetBool(config.LogColor),
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		DisableSorting:  true,
		ForceFormatting: true,
		FullTimestamp:   true,
	})
	switch strings.ToLower(config.GetString(config.LogLevel)) {
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
		break
	}

	L(ctx).Debugf("Log level: %s", logrus.GetLevel())
}