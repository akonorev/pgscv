package collector

import (
	"database/sql"
	"github.com/jackc/pgproto3/v2"
	"github.com/stretchr/testify/assert"
	"github.com/weaponry/pgscv/internal/model"
	"testing"
)

/* IMPORTANT: this test will fail if there are no functions stats in the databases or track_functions is disabled */

func TestPostgresFunctionsCollector_Update(t *testing.T) {
	var input = pipelineInput{
		required: []string{
			"postgres_function_calls_total",
			"postgres_function_total_time_seconds",
			"postgres_function_self_time_seconds",
		},
		collector: NewPostgresFunctionsCollector,
		service:   model.ServiceTypePostgresql,
	}

	pipeline(t, input)
}

func Test_parsePostgresFunctionsStat(t *testing.T) {
	var testCases = []struct {
		name string
		res  *model.PGResult
		want map[string]postgresFunctionStat
	}{
		{
			name: "normal output",
			res: &model.PGResult{
				Nrows: 3,
				Ncols: 6,
				Colnames: []pgproto3.FieldDescription{
					{Name: []byte("datname")}, {Name: []byte("schemaname")}, {Name: []byte("funcname")},
					{Name: []byte("calls")}, {Name: []byte("total_time")}, {Name: []byte("self_time")},
				},
				Rows: [][]sql.NullString{
					{
						{String: "testdb", Valid: true}, {String: "testschema1", Valid: true}, {String: "testfunction1", Valid: true},
						{String: "10", Valid: true}, {String: "1000", Valid: true}, {String: "900", Valid: true},
					},
					{
						{String: "testdb", Valid: true}, {String: "testschema2", Valid: true}, {String: "testfunction2", Valid: true},
						{String: "20", Valid: true}, {String: "2000", Valid: true}, {String: "700", Valid: true},
					},
					{
						{String: "testdb", Valid: true}, {String: "testschema3", Valid: true}, {String: "testfunction3", Valid: true},
						{String: "30", Valid: true}, {String: "3000", Valid: true}, {String: "600", Valid: true},
					},
				},
			},
			want: map[string]postgresFunctionStat{
				"testdb/testschema1/testfunction1": {
					datname: "testdb", schemaname: "testschema1", funcname: "testfunction1", calls: 10, totaltime: 1000, selftime: 900,
				},
				"testdb/testschema2/testfunction2": {
					datname: "testdb", schemaname: "testschema2", funcname: "testfunction2", calls: 20, totaltime: 2000, selftime: 700,
				},
				"testdb/testschema3/testfunction3": {
					datname: "testdb", schemaname: "testschema3", funcname: "testfunction3", calls: 30, totaltime: 3000, selftime: 600,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parsePostgresFunctionsStats(tc.res, []string{"usename", "schemaname", "funcname"})
			assert.EqualValues(t, tc.want, got)
		})
	}
}
