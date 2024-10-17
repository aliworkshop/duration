package duration

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    *Duration
		wantErr bool
	}{
		{
			name:    "invalid-duration-3",
			args:    args{d: "0SP0D"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "period-only",
			args: args{d: "4Y"},
			want: &Duration{
				Years: 4,
			},
			wantErr: false,
		},
		{
			name: "time-only-decimal",
			args: args{d: "2.5S"},
			want: &Duration{
				Seconds: 2.5,
			},
			wantErr: false,
		},
		{
			name: "full",
			args: args{d: "3Y6M4D12H30m5.5S"},
			want: &Duration{
				Years:   3,
				Months:  6,
				Days:    4,
				Hours:   12,
				Minutes: 30,
				Seconds: 5.5,
			},
			wantErr: false,
		},
		{
			name: "negative",
			args: args{d: "-5m"},
			want: &Duration{
				Minutes:  5,
				Negative: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.d)
			if (err != nil) != tt.wantErr {
				fmt.Println(tt)
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromTimeDuration(t *testing.T) {
	tests := []struct {
		give time.Duration
		want *Duration
	}{
		{
			give: 0,
			want: &Duration{},
		},
		{
			give: time.Minute * 94,
			want: &Duration{
				Hours:   1,
				Minutes: 34,
			},
		},
		{
			give: -time.Second * 10,
			want: &Duration{
				Seconds:  10,
				Negative: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.give.String(), func(t *testing.T) {
			got := FromTimeDuration(tt.give)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		give time.Duration
		want string
	}{
		{
			give: 0,
			want: "0S",
		},
		{
			give: time.Minute * 94,
			want: "1H34m",
		},
		{
			give: time.Hour * 72,
			want: "3D",
		},
		{
			give: time.Hour * 26,
			want: "1D2H",
		},
		{
			give: time.Second * 465461651,
			want: "14Y9M3D12H54m11S",
		},
		{
			give: -time.Hour * 99544,
			want: "-11Y4M1W4D",
		},
		{
			give: -time.Second * 10,
			want: "-10S",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Format(tt.give)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDuration_ToTimeDuration(t *testing.T) {
	type fields struct {
		Years    float64
		Months   float64
		Weeks    float64
		Days     float64
		Hours    float64
		Minutes  float64
		Seconds  float64
		Negative bool
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "seconds",
			fields: fields{
				Seconds: 33.3,
			},
			want: time.Second*33 + time.Millisecond*300,
		},
		{
			name: "hours, minutes, and seconds",
			fields: fields{
				Hours:   2,
				Minutes: 33,
				Seconds: 17,
			},
			want: time.Hour*2 + time.Minute*33 + time.Second*17,
		},
		{
			name: "days",
			fields: fields{
				Days: 2,
			},
			want: time.Hour * 24 * 2,
		},
		{
			name: "weeks",
			fields: fields{
				Weeks: 1,
			},
			want: time.Hour * 24 * 7,
		},
		{
			name: "fractional weeks",
			fields: fields{
				Weeks: 12.5,
			},
			want: time.Hour*24*7*12 + time.Hour*84,
		},
		{
			name: "negative",
			fields: fields{
				Hours:    2,
				Negative: true,
			},
			want: -time.Hour * 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := &Duration{
				Years:    tt.fields.Years,
				Months:   tt.fields.Months,
				Weeks:    tt.fields.Weeks,
				Days:     tt.fields.Days,
				Hours:    tt.fields.Hours,
				Minutes:  tt.fields.Minutes,
				Seconds:  tt.fields.Seconds,
				Negative: tt.fields.Negative,
			}
			if got := duration.ToTimeDuration(); got != tt.want {
				t.Errorf("ToTimeDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_String(t *testing.T) {
	duration, err := Parse("3Y6M4D12H30m5.5S")
	if err != nil {
		t.Fatal(err)
	}

	if duration.String() != "3Y6M4D12H30m5.5S" {
		t.Errorf("expected: %s, got: %s", "3Y6M4D12H30m5.5S", duration.String())
	}

	duration.Seconds = 33.3333

	if duration.String() != "3Y6M4D12H30m33.3333S" {
		t.Errorf("expected: %s, got: %s", "3Y6M4D12H30m33.3333S", duration.String())
	}

	smallDuration, err := Parse("0.0000000000001S")
	if err != nil {
		t.Fatal(err)
	}

	if smallDuration.String() != "0.0000000000001S" {
		t.Errorf("expected: %s, got: %s", "0.0000000000001S", smallDuration.String())
	}

	negativeDuration, err := Parse("-2H5m")
	if err != nil {
		t.Fatal(err)
	}

	if negativeDuration.String() != "-2H5m" {
		t.Errorf("expected: %s, got: %s", "-2H5m", negativeDuration.String())
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	td, err := Parse("3Y6M4D12H30m5.5S")
	if err != nil {
		t.Fatal(err)
	}

	jsonVal, err := json.Marshal(struct {
		Dur *Duration `json:"d"`
	}{Dur: td})
	if err != nil {
		t.Errorf("did not expect error: %s", err.Error())
	}
	if string(jsonVal) != `{"d":"3Y6M4D12H30m5.5S"}` {
		t.Errorf("expected: %s, got: %s", `{"d":"3Y6M4D12H30m5.5S"}`, string(jsonVal))
	}

	jsonVal, err = json.Marshal(struct {
		Dur Duration `json:"d"`
	}{Dur: *td})
	if err != nil {
		t.Errorf("did not expect error: %s", err.Error())
	}
	if string(jsonVal) != `{"d":"3Y6M4D12H30m5.5S"}` {
		t.Errorf("expected: %s, got: %s", `{"d":"3Y6M4D12H30m5.5S"}`, string(jsonVal))
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	jsonStr := `
		{
			"d": "3Y6M4D12H30m5.5S"
		}
	`
	expected, err := Parse("3Y6M4D12H30m5.5S")
	if err != nil {
		t.Fatal(err)
	}

	var durStructPtr struct {
		Dur *Duration `json:"d"`
	}
	err = json.Unmarshal([]byte(jsonStr), &durStructPtr)
	if err != nil {
		t.Errorf("did not expect error: %s", err.Error())
	}
	if !reflect.DeepEqual(durStructPtr.Dur, expected) {
		t.Errorf("JSON Unmarshal ptr got = %s, want %s", durStructPtr.Dur, expected)
	}

	var durStruct struct {
		Dur Duration `json:"d"`
	}
	err = json.Unmarshal([]byte(jsonStr), &durStruct)
	if err != nil {
		t.Errorf("did not expect error: %s", err.Error())
	}
	if !reflect.DeepEqual(durStruct.Dur, *expected) {
		t.Errorf("JSON Unmarshal ptr got = %s, want %s", &(durStruct.Dur), expected)
	}
}
