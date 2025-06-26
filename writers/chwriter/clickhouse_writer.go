package chwriter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickhouseWriter struct {
	Conn        driver.Conn
	hosts       []string
	user        string
	passwd      string
	db          string
	table       string
	FieldMap    map[string]string // mapping json field -> ch column name
	ExtraFields []string          // additional fields to write into separate columns. Should be present in CH table
}

var defaultFieldMap = map[string]string{
	"time":           "time",
	"component_type": "component_type",
	"component":      "component_name",
	"hostname":       "component_address",
	"address":        "component_ip",
	"version":        "component_version",
	"level":          "severity",
	"message":        "message",
	"details":        "details"}

var _ io.Writer = &ClickhouseWriter{}

func NewClickhouseWriter(hosts []string, user, passwd, db, table string) (*ClickhouseWriter, error) {
	conn, err := connect(hosts, user, passwd, db)
	if err != nil {
		return nil, err
	}

	writer := ClickhouseWriter{
		Conn:   conn,
		hosts:  hosts,
		db:     db,
		table:  table,
		user:   user,
		passwd: passwd,
	}

	return &writer, nil
}

func connect(hosts []string, user, passwd, db string) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: hosts,
			Auth: clickhouse.Auth{
				Database: db,
				Username: user,
				Password: passwd,
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "go-logger-client", Version: "0.1"},
				},
			},
			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			// TLS: &tls.Config{
			// 	InsecureSkipVerify: true,
			// },
			MaxOpenConns: 1,
			DialTimeout:  time.Duration(10) * time.Second,
			Settings:     map[string]any{"input_format_skip_unknown_fields": 1},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}

func (w *ClickhouseWriter) Write(p []byte) (n int, err error) {
	if w.Conn == nil {
		w.Conn, err = connect(w.hosts, w.user, w.passwd, w.db)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}

	if w.FieldMap == nil {
		w.FieldMap = defaultFieldMap
	}

	cols := make([]string, 0, len(w.FieldMap)+len(w.ExtraFields))
	for _, col := range []string{"time", "component_type", "component", "hostname", "address", "version", "level", "message", "details"} {
		cols = append(cols, w.FieldMap[col])
	}
	cols = append(cols, w.ExtraFields...)
	colStr := strings.Join(cols, ",")
	query := fmt.Sprintf(`INSERT INTO %s.%s (%s)`, w.db, w.table, colStr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := make(map[string]interface{})
	err = json.Unmarshal(p, &row)
	if err != nil {
		fmt.Println(err)
	}

	batch, err := w.Conn.PrepareBatch(ctx, query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	timeStr, _ := row["time"].(string)
	ts, _ := time.Parse(time.RFC3339Nano, timeStr)

	extraValues := make([]any, 0, len(w.ExtraFields))
	details := make(map[string]interface{})
	for k, v := range row {
		switch k {
		case "time", "component", "hostname", "address", "version", "level", "message":
		default:
			if slices.Contains(w.ExtraFields, k) {
				extraValues = append(extraValues, row[k])
				continue
			}
			details[k] = v
		}
	}

	detailsStr, _ := json.Marshal(details)

	rowToAdd := []any{ts,
		"",
		row["component"],
		row["hostname"],
		row["address"],
		row["version"],
		row["level"],
		row["message"],
		detailsStr,
	}
	rowToAdd = append(rowToAdd, extraValues...)
	if err = batch.Append(rowToAdd...); err != nil {
		fmt.Println(err)
		return 0, err
	}

	if err = batch.Send(); err != nil {
		fmt.Println(err)
		return 0, err
	}

	return len(p), nil
}
