package douyu

import (
	"reflect"
	"testing"
)

func Test_number2bytes(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name   string
		args   args
		wantBs []byte
	}{
		{
			name: "Test_number2bytes01",
			args: args{
				num: 987654321,
			},
			wantBs: []byte("987654321"),
		},
		{
			name: "Test_number2bytes02",
			args: args{
				num: 9876543210,
			},
			wantBs: []byte("9876543210"),
		},
		{
			name: "Test_number2bytes02",
			args: args{
				num: 100,
			},
			wantBs: []byte("100"),
		},
	}
	for _, tt := range tests {
		if gotBs := number2bytes(tt.args.num); !reflect.DeepEqual(gotBs, tt.wantBs) {
			t.Errorf("%q. number2bytes() = %v, want %v", tt.name, gotBs, tt.wantBs)
		}
	}
}

func Test_bytes2number(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name       string
		args       args
		wantNumber int64
	}{
		// TODO: Add test cases.
		{
			name: "Bytes2Number01",
			args: args{
				data: []byte("123456"),
			},
			wantNumber: 123456,
		},
		{
			name: "Bytes2Number02",
			args: args{
				data: []byte("000"),
			},
			wantNumber: 0,
		},
		{
			name: "Bytes2Number02",
			args: args{
				data: []byte("12304560"),
			},
			wantNumber: 12304560,
		},
	}
	for _, tt := range tests {
		if gotNumber := bytes2number(tt.args.data); gotNumber != tt.wantNumber {
			t.Errorf("%q. bytes2number() = %v, want %v", tt.name, gotNumber, tt.wantNumber)
		}
	}
}
