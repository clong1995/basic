package excel

import (
	"testing"
)

func Test_server_AllRows(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		args     args
		wantRows [][]string
		wantErr  bool
	}{
		{
			name: "本地文件测试",
			args: args{
				filename: "/Users/yuchenglong/Desktop/excel/组织架构_0_0_1.xlsx",
			},
		},
		{
			name: "网络文件测试",
			args: args{
				filename: "https://jifen-app.oss-cn-beijing.aliyuncs.com/temp/4BsRJoqH6To.xlsx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := server{}
			gotRows, err := s.AllRows(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*if !reflect.DeepEqual(gotRows, tt.wantRows) {
				t.Errorf("AllRows() gotRows = %v, want %v", gotRows, tt.wantRows)
			}*/
			t.Logf("AllRows() gotRows = %#v", gotRows)
		})
	}
}
