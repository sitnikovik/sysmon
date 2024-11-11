package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
)

// argsToInterfacesForOS maps command string arguments to interfaces for os
func argsToInterfacesForOS(os string) []interface{} {
	res := make([]interface{}, len(argsByOS[os]))
	for i, arg := range argsByOS[os] {
		res[i] = interface{}(arg)
	}

	return res
}

func Test_parser_parseForDarwin(t *testing.T) {
	t.Parallel()

	type fields struct {
		execerMockFunc func(t *testing.T) cmd.Execer
	}
	tests := []struct {
		name    string
		fields  fields
		want    CpuStats
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {

					execer := cmd.NewMockExecer(t)
					execer.EXPECT().
						Exec(cmdByOS[os.Darwin], argsToInterfacesForOS(os.Darwin)...).
						Return(&cmd.Result{
							Bytes: []byte("CPU usage: 10.0% user, 20.0% sys, 70.0% idle"),
						}, nil)

					return execer
				},
			},
			want: CpuStats{
				User:   10.0,
				System: 20.0,
				Idle:   70.0,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				execer: tt.fields.execerMockFunc(t),
			}
			got, err := p.parseForDarwin()

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
