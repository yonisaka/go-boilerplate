package locales

import (
	"context"
	"testing"
)

var data = map[string]string{
	"id.api.document.status.signed":    "Ditandatangani",
	"id.api.document.status.signed_by": "Ditandatangani oleh {{.privyID}}",
	"en.api.document.status.signed":    "Signed",
	"en.api.document.status.signed_by": "Signed by {{.privyID}}",
}

func Test_translator_Translate(t1 *testing.T) {
	type fields struct {
		data map[string]string
	}

	ctx := WithAcceptLanguage(context.Background(), "enx")
	ctxID := WithAcceptLanguage(context.Background(), "id")

	type args struct {
		ctx    context.Context
		key    string
		params map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "message without param",
			fields: fields{
				data: data,
			},
			args: args{
				ctx:    ctx,
				key:    "api.document.status.signed",
				params: nil,
			},
			want: "Signed",
		},
		{
			name: "message with param",
			fields: fields{
				data: data,
			},
			args: args{
				ctx: ctx,
				key: "api.document.status.signed_by",
				params: map[string]interface{}{
					"privyID": "DA787",
				},
			},
			want: "Signed by DA787",
		},
		{
			name: "message without param",
			fields: fields{
				data: data,
			},
			args: args{
				ctx:    ctxID,
				key:    "api.document.status.signed",
				params: nil,
			},
			want: "Ditandatangani",
		},
		{
			name: "message with param",
			fields: fields{
				data: data,
			},
			args: args{
				ctx: ctxID,
				key: "api.document.status.signed_by",
				params: map[string]interface{}{
					"privyID": "DA787",
				},
			},
			want: "Ditandatangani oleh DA787",
		},
		{
			name: "message with privy id but param is nil",
			fields: fields{
				data: data,
			},
			args: args{
				ctx:    ctxID,
				key:    "api.document.status.signed_by",
				params: nil,
			},
			want: "Ditandatangani oleh <no value>",
		},
		{
			name: "message with param and invalid key",
			fields: fields{
				data: data,
			},
			args: args{
				ctx: ctxID,
				key: "api.document.status.signed_by_x",
				params: map[string]interface{}{
					"privyID": "DA787",
				},
			},
			want: "invalid key : api.document.status.signed_by_x",
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := translator{
				data:        tt.fields.data,
				defaultLang: "en",
				langs:       []string{"en", "id"},
			}

			if got := t.Translate(tt.args.ctx, tt.args.key, tt.args.params); got != tt.want {
				t1.Errorf("Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}
