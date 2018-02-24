package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mgutz/ansi"
)

type StdoutWriter struct {
	out           io.Writer
	noColor       bool
	noGroup       bool
	noStream      bool
	timestamp     bool
	ingestionTime bool
}

func NewStdoutWriter(noColor, noGroup, noStream, timestamp, ingestion bool) (*StdoutWriter, error) {
	w := &StdoutWriter{
		out:           os.Stdout,
		noColor:       noColor,
		noGroup:       noGroup,
		noStream:      noStream,
		timestamp:     timestamp,
		ingestionTime: ingestion,
	}

	return w, nil
}

// <TODO> buffering
// <TODO> condition optimization
func (w *StdoutWriter) Write(stream chan Event) error {
	buf := &bytes.Buffer{}

	for {
		event, ok := <-stream
		if !ok {
			return nil
		}
		buf.Reset()
		if !w.noGroup {
			if w.noColor {
				buf.WriteString(event.Group + " ")
			} else {
				buf.WriteString(ansi.Color(event.Group, "green") + " ")
			}
		}
		if !w.noStream {
			if w.noColor {
				buf.WriteString(event.Stream + " ")
			} else {
				buf.WriteString(ansi.Color(event.Stream, "cyan") + " ")
			}
		}
		if w.timestamp {
			if w.noColor {
				buf.WriteString(event.Timestamp.Format(time.RFC3339) + " ")
			} else {
				buf.WriteString(ansi.Color(event.Timestamp.Format(time.RFC3339), "yellow") + " ")
			}
		}
		if w.ingestionTime {
			if w.noColor {
				buf.WriteString(event.IngestionTime.Format(time.RFC3339) + " ")
			} else {
				buf.WriteString(ansi.Color(event.IngestionTime.Format(time.RFC3339), "blue") + " ")
			}
		}

		buf.WriteString(strings.TrimSpace(event.Message))
		buf.WriteString("\n")

		buf.WriteTo(w.out)
	}

	return nil
}
