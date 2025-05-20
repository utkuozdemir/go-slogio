package slogio_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vgarvardt/slogex/observer"

	"github.com/utkuozdemir/go-slogio"
)

func TestWriter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc   string
		level  slog.Level // defaults to info
		writes []string
		want   []slog.Record
	}{
		{
			desc: "simple",
			writes: []string{
				"foo\n",
				"bar\n",
				"baz\n",
			},
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "foo"},
				{Level: slog.LevelInfo, Message: "bar"},
				{Level: slog.LevelInfo, Message: "baz"},
			},
		},
		{
			desc:  "level too low",
			level: slog.LevelDebug,
			writes: []string{
				"foo\n",
				"bar\n",
			},
			want: []slog.Record{},
		},
		{
			desc:  "multiple newlines in a message",
			level: slog.LevelWarn,
			writes: []string{
				"foo\nbar\n",
				"baz\n",
				"qux\nquux\n",
			},
			want: []slog.Record{
				{Level: slog.LevelWarn, Message: "foo"},
				{Level: slog.LevelWarn, Message: "bar"},
				{Level: slog.LevelWarn, Message: "baz"},
				{Level: slog.LevelWarn, Message: "qux"},
				{Level: slog.LevelWarn, Message: "quux"},
			},
		},
		{
			desc:  "message split across multiple writes",
			level: slog.LevelError,
			writes: []string{
				"foo",
				"bar\nbaz",
				"qux",
			},
			want: []slog.Record{
				{Level: slog.LevelError, Message: "foobar"},
				{Level: slog.LevelError, Message: "bazqux"},
			},
		},
		{
			desc: "blank lines in the middle",
			writes: []string{
				"foo\n\nbar\nbaz",
			},
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "foo"},
				{Level: slog.LevelInfo, Message: ""},
				{Level: slog.LevelInfo, Message: "bar"},
				{Level: slog.LevelInfo, Message: "baz"},
			},
		},
		{
			desc: "blank line at the end",
			writes: []string{
				"foo\nbar\nbaz\n",
			},
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "foo"},
				{Level: slog.LevelInfo, Message: "bar"},
				{Level: slog.LevelInfo, Message: "baz"},
			},
		},
		{
			desc: "multiple blank line at the end",
			writes: []string{
				"foo\nbar\nbaz\n\n",
			},
			want: []slog.Record{
				{Level: slog.LevelInfo, Message: "foo"},
				{Level: slog.LevelInfo, Message: "bar"},
				{Level: slog.LevelInfo, Message: "baz"},
				{Level: slog.LevelInfo, Message: ""},
			},
		},
	}

	for _, tt := range tests {
		tt := tt // for t.Parallel
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			handler, observed := observer.New(nil)

			w := slogio.Writer{
				Log:   slog.New(handler),
				Level: tt.level,
			}

			for _, s := range tt.writes {
				_, err := io.WriteString(&w, s)
				require.NoError(t, err, "Writer.Write failed.")
			}

			assert.NoError(t, w.Close(), "Writer.Close failed.")

			//// Turn []observer.LoggedRecord => []slog.Record
			got := make([]slog.Record, observed.Len())
			for i, ent := range observed.AllUntimed() {
				got[i] = ent.Record
			}
			assert.Equal(t, tt.want, got, "Logged entries do not match.")
		})
	}
}
