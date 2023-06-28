package locales

import "testing"

func TestParseIOSLang(t *testing.T) {
	type args struct {
		in  string
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "en",
			args: args{
				in:  "en",
				def: "en",
			},
			want: "en",
		},
		{
			name: "id",
			args: args{
				in:  "id",
				def: "en",
			},
			want: "id",
		},
		{
			name: "en-ID;q=1.0, id-ID;q=0.9",
			args: args{
				in:  "en-ID;q=1.0, id-ID;q=0.9",
				def: "en",
			},
			want: "en",
		},
		{
			name: "id-ID;q=1.0, en-ID;q=0.9",
			args: args{
				in:  "id-ID;q=1.0, en-ID;q=0.9",
				def: "en",
			},
			want: "id",
		},
		{
			name: "en-ID;q=0.9, id-ID;q=1.0",
			args: args{
				in:  "en-ID;q=0.9, id-ID;q=1.0",
				def: "en",
			},
			want: "id",
		},
		{
			name: "invalid accept lang",
			args: args{
				in:  "any string",
				def: "en",
			},
			want: "en",
		},
		{
			name: "invalid accept lang",
			args: args{
				in:  "idx",
				def: "en",
			},
			want: "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseIOSLang(tt.args.in, tt.args.def); got != tt.want {
				t.Errorf("ParseIOSLang() = %v, want %v", got, tt.want)
			}
		})
	}
}
