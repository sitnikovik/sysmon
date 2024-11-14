package disk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/sysmon/internal/metrics/utils/cmd"
	"github.com/sitnikovik/sysmon/internal/metrics/utils/os"
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
		want    models.DiskStats
		wantErr bool
	}{
		{
			name: "ok darwin",
			fields: fields{
				execerMockFunc: func(t *testing.T) cmd.Execer {
					execer := cmd.NewMockExecer(t)

					execer.EXPECT().
						Exec(darwinCmdDiskLoad, strings.ToInterfaces(darwinArgsDiskLoad)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"          disk0           disk1\n" +
									"KB/t tps  MB/s     KB/t tps  MB/s\n" +
									"  32  10   0.31      64  20   1.25\n",
							),
						}, nil).Once()

					execer.EXPECT().
						Exec(darwinCmdDiskSpaceInodes, strings.ToInterfaces(darwinArgsDiskSpaceInodes)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Filesystem      Inodes   IUsed   IFree IUse% Mounted on\n" +
									"/dev/sda1      3276800  1048576 2228224   32% /\n" +
									"tmpfs           128000     4000  124000    3% /run\n",
							),
						}, nil).Once()

					execer.EXPECT().
						Exec(darwinCmdDiskSpace, strings.ToInterfaces(darwinArgsDiskSpace)...).
						Return(&cmd.Result{
							Bytes: []byte(
								"Filesystem      Size  Used Avail Use% Mounted on\n" +
									"/dev/sda1        50G   20G   30G  40% /\n" +
									"tmpfs           500M  100M  400M  20% /run\n",
							),
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
			want: models.DiskStats{
				Reads:             10,
				Writes:            20,
				ReadWriteKB:       10*32 + 20*64,
				TotalMB:           50 * 1024,
				UsedMB:            20 * 1024,
				UsedPercent:       40,
				UsedInodes:        1048576,
				UsedInodesPercent: 32,
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

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
