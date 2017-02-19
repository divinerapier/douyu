package douyu

import (
	"reflect"
	"testing"
)

func TestOpenDanmu(t *testing.T) {
	type args struct {
		rid int64
	}
	tests := []struct {
		name    string
		args    args
		wantDy  *Douyu
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotDy, err := OpenDanmu(tt.args.rid)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. OpenDanmu() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotDy, tt.wantDy) {
			t.Errorf("%q. OpenDanmu() = %v, want %v", tt.name, gotDy, tt.wantDy)
		}
	}
}

func TestDouyu_login(t *testing.T) {
	tests := []struct {
		name string
		dy   *Douyu
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt.dy.login()
	}
}

func Test_parseLoginResponse(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want *LoginResponse
	}{
		// TODO: Add test cases.
		{
			name: "ParseLoginResponse01",
			args: args{
				data: []byte("type@=loginres/userid@=123450/roomgroup@=123450/pg@=123450/sessionid@=123450/username@=hello/nickname@=world/is_signined@=1/signin_count@=123450/live_stat@=1/npv@=1/best_dlev@=123450/cur_lev@=123450/"),
			},

			want: &LoginResponse{
				Type:            []byte("loginres"),
				Username:        []byte("hello"),
				Nickname:        []byte("world"),
				UserID:          123450,
				RoomGroup:       123450,
				Pg:              123450,
				SessionID:       123450,
				SignedCount:     123450,
				BestDlev:        123450,
				CurLev:          123450,
				IsSigned:        true,
				LiveStat:        true,
				NeedPhoneVerify: true,
			},
		},
	}
	for _, tt := range tests {
		if got := parseLoginResponse(tt.args.data); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. parseLoginResponse() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
