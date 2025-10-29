package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/guruakashsm/logvault/filter"
	"github.com/guruakashsm/logvault/metrics"
	"github.com/guruakashsm/logvault/model"
	"github.com/guruakashsm/logvault/parser"
	"github.com/guruakashsm/logvault/segmenter"
	"github.com/guruakashsm/logvault/store"
	"github.com/guruakashsm/logvault/utils"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	logstore := &model.LogStore{
		Metrics: &model.LogStoreMetrics{
			ByLevel:     make(map[model.LogLevel]int),
			ByComponent: make(map[string]int),
			ByHost:      make(map[string]int),
			ByRequestID: make(map[string]int),
		},
	}

	logParser := parser.NewRegexParser(`(?P<timestamp>\d{4}[-/]\d{2}[-/]\d{2} \d{2}:\d{2}:\d{2}\.\d+) \| (?P<level>\w+) \| (?P<component>[\w-]+) \| host=(?P<host>[\w-]+) \| request_id=(?P<request_id>[\w-]+) \| msg="(?P<message>.*?)"`)
	metricsUpdater := metrics.NewMetricsUpdater()
	storeHandler := store.NewStoreHandler(logstore, metricsUpdater)

	folderPath := "../logs"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		slog.Error("Failed to read folder", "path", folderPath, "err", err)
		return
	}

	slog.Info("Starting log processing...", "files_found", len(files))

	for i, file := range files {
		if file.IsDir() {
			continue
		}

		segment := &model.Segment{
			ID:       logstore.Metrics.TotalSegments + 1,
			FileName: file.Name(),
			FilePath: filepath.Join(folderPath, file.Name()),
			Metrics: &model.SegmentMetrics{
				ByLevel:     make(map[model.LogLevel]int),
				ByComponent: make(map[string]int),
				ByHost:      make(map[string]int),
				ByRequestID: make(map[string]int),
			},
			Index: &model.SegmentIndex{
				ByLevel:     make(map[model.LogLevel][]int),
				ByComponent: make(map[string][]int),
				ByHost:      make(map[string][]int),
				ByRequestID: make(map[string][]int),
			},
		}

		slog.Info("Processing file", "index", i+1, "file", file.Name())

		segmentHandler := segmenter.NewSegmentHandler(logParser, metricsUpdater, segment)

		f, err := os.Open(segment.FilePath)
		if err != nil {
			slog.Warn("Failed to open file", "file", segment.FilePath, "err", err)
			continue
		}

		scanner := bufio.NewScanner(f)
		lineCount := 0
		for scanner.Scan() {
			line := scanner.Text()
			if err := segmentHandler.AddLine(line); err != nil {
				slog.Warn("Failed to add line", "line_number", lineCount, "file", file.Name(), "err", err)
			}
			lineCount++
		}
		f.Close()

		if err := scanner.Err(); err != nil {
			slog.Warn("Error reading file", "file", segment.FilePath, "err", err)
		}

		if err := storeHandler.AddSegment(segment); err != nil {
			slog.Warn("Failed to add segment", "file", segment.FilePath, "err", err)
		}

		slog.Info("Completed file", "file", file.Name(), "lines_parsed", lineCount)
	}

	fmt.Print("\033[H\033[2J")

	slog.Info("Log processing completed")
	slog.Info("Metrics",
		"Total_Logs", logstore.Metrics.TotalLogs,
		"Total_Segments", logstore.Metrics.TotalSegments,
		"Total_Level", len(logstore.Metrics.ByLevel),
		"Total_Component", len(logstore.Metrics.ByComponent),
		"Total_Host", len(logstore.Metrics.ByHost),
		"Total_Request_ID", len(logstore.Metrics.ByRequestID),
	)

	startCLI(logstore)
}

func startCLI(logstore *model.LogStore) {
	reader := bufio.NewReader(os.Stdin)
	slog.Info("Interactive Log Filter CLI started. Type 'help' for commands.")

	var lf filter.LogFilter
	filterhandler := filter.NewFilterHandler(logstore)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)

		switch {
		case cmd == "exit":
			slog.Info("👋 Exiting CLI...")
			return

		case cmd == "help":
			showHelp()

		case cmd == "clear":
			lf = filter.LogFilter{}
			slog.Info("✅ Cleared all filters")

		case strings.HasPrefix(cmd, "start="):
			t, _ := utils.ParseTime(strings.TrimPrefix(cmd, "start="), time.DateTime, time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano, time.DateOnly, time.TimeOnly)
			lf.StartTime = t
			slog.Info("🕒 Start time set", "time", t)

		case strings.HasPrefix(cmd, "end="):
			t, _ := utils.ParseTime(strings.TrimPrefix(cmd, "end="), time.DateTime, time.Layout, time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano, time.DateOnly, time.TimeOnly)
			lf.EndTime = t
			slog.Info("🕒 End time set", "time", t)

		case strings.HasPrefix(cmd, "level="):
			parts := strings.Split(strings.TrimPrefix(cmd, "level="), ",")
			var levels []model.LogLevel
			for _, p := range parts {
				levels = append(levels, model.LogLevel(strings.TrimSpace(p)))
			}
			lf.Levels = levels
			slog.Info("📊 Levels set", "levels", levels)

		case strings.HasPrefix(cmd, "component="):
			lf.Components = strings.Split(strings.TrimPrefix(cmd, "component="), ",")
			slog.Info("⚙️ Components set", "components", lf.Components)

		case strings.HasPrefix(cmd, "host="):
			lf.Hosts = strings.Split(strings.TrimPrefix(cmd, "host="), ",")
			slog.Info("💻 Hosts set", "hosts", lf.Hosts)

		case strings.HasPrefix(cmd, "req="):
			lf.RequestIDs = strings.Split(strings.TrimPrefix(cmd, "req="), ",")
			slog.Info("🔍 RequestIDs set", "request_ids", lf.RequestIDs)

		case strings.HasPrefix(cmd, "msg="):
			v := strings.Trim(strings.TrimPrefix(cmd, "msg="), `"`)
			lf.MessageContains = &v
			slog.Info("🗒️ Message filter set", "contains", v)

		case strings.HasPrefix(cmd, "log="):
			v := strings.Trim(strings.TrimPrefix(cmd, "log="), `"`)
			lf.LogContains = &v
			slog.Info("📄 Log text filter set", "contains", v)

		case cmd == "show":
			results := filterhandler.Filter(lf)
			slog.Info("✅ Found matching logs", "count", len(results))
			printSample(results, 500)

		default:
			slog.Warn("⚠️ Unknown command", "cmd", cmd)
		}
	}
}

func showHelp() {
	fmt.Printf(`
Available commands:
  start=<time>         - Set start time (e.g. start=2025-10-26T14:00:00Z)
  end=<time>           - Set end time
  level=INFO,ERROR     - Filter by log levels
  component=api,db     - Filter by components
  host=web01,worker02  - Filter by hosts
  req=req-123,req-456  - Filter by request IDs
  msg="partial text"   - Message contains
  log="full log text"  - Raw log text contains
  show                 - Execute current filters
  clear                - Clear filters
  help                 - Show this help
  exit                 - Exit CLI
`)
}

func printSample(entries []model.LogEntry, n int) {
	if len(entries) == 0 {
		slog.Warn("No logs found")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	limit := min(len(entries), n)

	fmt.Println("📄 ---- Sample Logs ---- 📄")

	for i := range limit {
		e := entries[i]

		var level model.LogLevel
		switch e.Level {
		case "INFO":
			level = "\033[34mINFO \033[0m" // Blue
		case "WARN":
			level = "\033[33mWARN \033[0m" // Yellow
		case "ERROR":
			level = "\033[31mERROR\033[0m" // Red
		default:
			level = "\033[37m" + e.Level + "\033[0m" // White
		}

		fmt.Printf(" [ %s ] [ %-5s ] %-10s | %-10s | %s | %s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			level,
			e.Component,
			e.Host,
			e.RequestID,
			e.Message,
		)

	}

	fmt.Println("X-X-X ---------------------- X-X-X")
}
