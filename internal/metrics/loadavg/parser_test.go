package loadavg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
)

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
		want    LoadAverageStats
		wantErr bool
	}{
		{
			name: "ok darwin",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdUnix).
						Return(&cmd.Result{
							Bytes: []byte("23:04  up 42 days, 13:14, 1 user, load averages: 3.99 3.95 3.58"),
						}, nil).Once()

					execer.EXPECT().
						OS().
						Return(os.Darwin).Once()

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: LoadAverageStats{
				OneMinute:     3.99,
				FiveMinute:    3.95,
				FifteenMinute: 3.58,
			},
			wantErr: false,
		},
		{
			name: "ok linux",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdUnix).
						Return(&cmd.Result{
							Bytes: []byte("23:04  up 42 days, 13:14, 1 user, load averages: 3.99 3.95 3.58"),
						}, nil).Once()

					execer.EXPECT().
						OS().
						Return(os.Linux).Once()

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: LoadAverageStats{
				OneMinute:     3.99,
				FiveMinute:    3.95,
				FifteenMinute: 3.58,
			},
			wantErr: false,
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

			require.Equal(t, tt.wantErr, err != nil, "unexpected error: %v", err)
			require.Equal(t, tt.want, got)
		})
	}
}
