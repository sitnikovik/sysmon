package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/strings"
	"github.com/sitnikovik/sysmon/internal/models"
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
		want    models.MemoryStats
		wantErr bool
	}{
		{
			name: "ok darwin",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					t.Helper()

					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdDarwin).
						Return(&cmd.Result{
							Bytes: []byte(
								"Mach Virtual Memory Statistics: (page size of 16384 bytes)\n" +
									"Pages free:                          9680.\n" +
									"Pages active:                        776961.\n" +
									"Pages inactive:                      776084.\n" +
									"Pages speculative:                          228.\n" +
									"Pages throttled:                              0.\n" +
									"Pages wired down:                        154979.\n" +
									"Pages purgeable:                           2848.\n" +
									"\"Translation faults\":                3335501749.\n" +
									"Pages copy-on-write:                  435408207.\n" +
									"Pages zero filled:                   1419852705.\n" +
									"Pages reactivated:                     10693508.\n" +
									"Pages purged:                           7790169.\n" +
									"File-backed pages:                       469637.\n" +
									"Anonymous pages:                        1083636.\n" +
									"Pages stored in compressor:              626905.\n" +
									"Pages occupied by compressor:            322718.\n" +
									"Decompressions:                         2316431.\n" +
									"Compressions:                           4834820.\n" +
									"Pageins:                               13566120.\n" +
									"Pageouts:                                137914.\n" +
									"Swapins:                                      0.\n" +
									"Swapouts:                                     0.\n",
							),
						}, nil).Once()

					execer.EXPECT().OS().Return("darwin").Once()

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: models.MemoryStats{
				TotalMb:     26842,
				AvailableMb: 12277,
				FreeMb:      151,
				UsedMb:      14565,
				ActiveMb:    12140,
				InactiveMb:  12126,
				WiredMb:     2421,
			},
		},
		{
			name: "ok linux",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(cmdLinux, strings.ToInterfaces(cmdLinuxArgs)...).
						Return(&cmd.Result{
							// free -m
							Bytes: []byte(
								"              total        used        free      shared  buff/cache   available\n" +
									"Mem:          32159       10659        1234        1234       20265       20265\n" +
									"Swap:          2047          10        2037\n",
							),
						}, nil).
						Once()

					execer.EXPECT().
						OS().
						Return("linux")

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: models.MemoryStats{
				TotalMb:     32159,
				AvailableMb: 20265,
				FreeMb:      1234,
				UsedMb:      10659,
				BuffersMb:   20265,
				CachedMb:    20265,
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
			got, err := p.Parse(tt.args.ctx)

			require.NoError(t, err, "unexpected error = %v", err)
			require.Equal(t, tt.want, got)
		})
	}
}
