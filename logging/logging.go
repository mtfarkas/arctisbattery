// Package logging sets up file-based logging so crashes and errors
// can be diagnosed even when the app runs as a tray application
// with no visible console.
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

const (
	logDirName    = "arctisbattery"
	maxLogFiles   = 5
	logTimeFormat = "2006-01-02_15-04-05"
)

var logFile *os.File

// Init configures the global logger to write to both stderr and a
// timestamped log file in the user's config/log directory.
// Returns a cleanup function that should be deferred in main().
func Init() (cleanup func(), err error) {
	logDir, err := logDirectory()
	if err != nil {
		return nil, fmt.Errorf("could not determine log directory: %w", err)
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create log directory: %w", err)
	}

	pruneOldLogs(logDir)

	filename := fmt.Sprintf("arctisbattery_%s.log", time.Now().Format(logTimeFormat))
	logPath := filepath.Join(logDir, filename)

	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stderr, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	log.Printf("Logging initialized: %s", logPath)

	cleanup = func() {
		if logFile != nil {
			logFile.Close()
		}
	}
	return cleanup, nil
}

// RecoverAndLog recovers from a panic, logs the stack trace, and
// optionally calls a callback. Use it as: defer logging.RecoverAndLog("name", nil)
func RecoverAndLog(goroutineName string, onPanic func(err interface{})) {
	if r := recover(); r != nil {
		RecoverAndLogValue(goroutineName, r)

		if onPanic != nil {
			onPanic(r)
		}
	}
}

// RecoverAndLogValue logs a panic value and stack trace. Use this when you've
// already called recover() yourself and need to log the result.
func RecoverAndLogValue(goroutineName string, r interface{}) {
	stack := string(debug.Stack())
	log.Printf("PANIC in %s: %v\n%s", goroutineName, r, stack)
}

func logDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// fall back to temp directory
		return filepath.Join(os.TempDir(), logDirName, "logs"), nil
	}
	return filepath.Join(configDir, logDirName, "logs"), nil
}

func pruneOldLogs(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var logFiles []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "arctisbattery_") && strings.HasSuffix(e.Name(), ".log") {
			logFiles = append(logFiles, e)
		}
	}

	if len(logFiles) <= maxLogFiles {
		return
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Name() < logFiles[j].Name()
	})

	for _, f := range logFiles[:len(logFiles)-maxLogFiles] {
		os.Remove(filepath.Join(dir, f.Name()))
	}
}
