package cpu

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/strings"
)

// // argsToInterfacesForOS maps command string arguments to interfaces for os
// func argsToInterfacesForOS(os string) []interface{} {
// 	res := make([]interface{}, len(argsByOS[os]))
// 	for i, arg := range argsByOS[os] {
// 		res[i] = interface{}(arg)
// 	}

// 	return res
// }

func TestNewParser(t *testing.T) {
	t.Parallel()

	t.Run("not nil on nil args", func(t *testing.T) {
		t.Parallel()
		assert.NotNil(t, NewParser(nil))
	})

	t.Run("with execer", func(t *testing.T) {
		t.Parallel()
		assert.NotNil(t, NewParser(cmd.NewExecer()))
	})
}

func Test_parser_Parse(t *testing.T) {
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
			name: "ok darwin",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdByOS[os.Darwin], strings.ToInterfaces(argsByOS[os.Darwin])...).
						Return(&cmd.Result{
							Bytes: []byte("CPU usage: 10.0% user, 20.0% sys, 70.0% idle"),
						}, nil)

					execer.EXPECT().
						OS().
						Return(os.Darwin)

					return execer
				},
			},
			want: CpuStats{
				User:   10.0,
				System: 20.0,
				Idle:   70.0,
			},
		},
		{
			name: "err darwin invalid cmd",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdByOS[os.Darwin], strings.ToInterfaces(argsByOS[os.Darwin])...).
						Return(nil, errors.New("invalid cmd"))

					execer.EXPECT().
						OS().
						Return(os.Darwin)

					return execer
				},
			},
			wantErr: true,
		},
		{
			name: "err unsupported os",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						OS().
						Return("testOS")

					return execer
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &parser{
				execer: tt.fields.execerMockFunc(t),
			}
			got, err := p.Parse()

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
