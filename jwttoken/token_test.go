package jwttoken

import (
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	type args struct {
		username string
		id       int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "TestGenerate_case_1",
			args: args{
				username: "admin",
				id:       10002,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.username, tt.args.id, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
			//if got != tt.want {
			//	t.Errorf("Generate() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		username string
		id       int64
		token    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Claims
		wantErr bool
	}{
		{
			name: "TestParse_case_1",
			args: args{
				username: "admin",
				id:       10001,
				token:    "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := Generate(tt.args.username, tt.args.id, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			zaplogger.Sugar().Info("tokenString ", tokenString)
			tt.args.token = tokenString
			got, err := Parse(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			zaplogger.Sugar().Info("got ", got)
			if !reflect.DeepEqual(got.Username, tt.args.username) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_id(t *testing.T) {

	var id int64 = 12354
	name := "saklsaska"

	token, err := Generate(name, id, true)
	assert.NoError(t, err)

	claims, err := Parse(token)
	assert.NoError(t, err)

	assert.Equal(t, id, claims.UserId)
}
