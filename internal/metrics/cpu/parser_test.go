package cpu

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/strings"
	"github.com/sitnikovik/sysmon/internal/models"
)

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
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.CPUStats
		wantErr bool
	}{
		{
			name: "ok darwin",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					t.Helper()

					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdDarwin, strings.ToInterfaces(argsDarwin)...).
						Return(&cmd.Result{
							Bytes: []byte("CPU usage: 10.0% user, 20.0% sys, 70.0% idle"),
						}, nil)

					execer.EXPECT().
						OS().
						Return(os.Darwin)

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: models.CPUStats{
				User:   10.0,
				System: 20.0,
				Idle:   70.0,
			},
		},
		{
			name: "err darwin invalid cmd",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					t.Helper()

					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdDarwin, strings.ToInterfaces(argsDarwin)...).
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
					t.Helper()

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
			got, err := p.Parse(tt.args.ctx)

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
