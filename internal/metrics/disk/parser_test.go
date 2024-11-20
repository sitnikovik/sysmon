package disk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
	stringsUtils "github.com/sitnikovik/sysmon/internal/metrics/utils/strings"
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
		want    models.DiskStats
		wantErr bool
	}{
		{
			name: "darwin",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					t.Helper()

					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(unixCmdDiskLoad, stringsUtils.ToInterfaces(unixArgsDiskLoad)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"          disk0           disk1\n" +
									"KB/t tps  MB/s     KB/t tps  MB/s\n" +
									"  32  10   0.31      64  20   1.25\n",
							),
						}, nil).
						Once()

					execer.EXPECT().
						Exec(unixCmdDiskSpaceInodes, stringsUtils.ToInterfaces(unixArgsDiskSpaceInodes)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Filesystem      Inodes   IUsed   IFree IUse% Mounted on\n" +
									"/dev/sda1      3276800  1048576 2228224   32% /\n" +
									"tmpfs           128000     4000  124000    3% /run\n",
							),
						}, nil).
						Once()

					execer.EXPECT().
						Exec(unixCmdDiskSpace, stringsUtils.ToInterfaces(unixArgsDiskSpace)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Filesystem      Size  Used Avail Use% Mounted on\n" +
									"/dev/sda1        50G   20G   30G  40% /\n" +
									"tmpfs           500M  100M  400M  20% /run\n",
							),
						}, nil).Once()

					execer.EXPECT().
						OS().
						Return(os.Darwin)

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: models.DiskStats{
				Reads:             10,
				Writes:            20,
				ReadWriteKb:       10*32 + 20*64,
				TotalMb:           50 * 1024,
				UsedMb:            20 * 1024,
				UsedPercent:       40,
				UsedInodes:        1048576,
				UsedInodesPercent: 32,
			},
		},
		{
			name: "linux",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					t.Helper()

					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(unixCmdDiskLoad, stringsUtils.ToInterfaces(unixArgsDiskLoad)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Linux 4.15.0-112-generic (hostname) 	09/01/2021 	_x86_64_	(4 CPU)\n" +
									"\n" +
									"Device             tps    kB_read/s    kB_wrtn/s    kB_read    kB_wrtn\n" +
									"sda               1.00         20.00         350.00          50          10\n",
							),
						}, nil).
						Once()

					execer.EXPECT().
						Exec(unixCmdDiskSpaceInodes, stringsUtils.ToInterfaces(unixArgsDiskSpaceInodes)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Filesystem      Inodes   IUsed   IFree IUse% Mounted on\n" +
									"/dev/sda1      3276800  1048576 2228224   32% /\n" +
									"tmpfs           128000     4000  124000    3% /run\n",
							),
						}, nil).
						Once()

					execer.EXPECT().
						Exec(unixCmdDiskSpace, stringsUtils.ToInterfaces(unixArgsDiskSpace)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Filesystem      Size  Used Avail Use% Mounted on\n" +
									"/dev/sda1        50G   20G   30G  40% /\n" +
									"tmpfs           500M  100M  400M  20% /run\n",
							),
						}, nil).Once()

					execer.EXPECT().
						OS().
						Return(os.Linux)

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: models.DiskStats{
				Reads:             1,
				Writes:            350,
				ReadWriteKb:       20 + 350,
				TotalMb:           50 * 1024,
				UsedMb:            20 * 1024,
				UsedPercent:       40,
				UsedInodes:        1048576,
				UsedInodesPercent: 32,
			},
		},
		{
			name: "err windows unsupported",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					t.Helper()

					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						OS().
						Return(os.Windows)

					return execer
				},
			},
			args: args{
				ctx: context.Background(),
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

			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
