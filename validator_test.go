package validator

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "valid struct with valid nested struct",
			args: args{
				v: struct {
					InInt        int `validate:"in:20,25,30"`
					NestedStruct struct {
						MaxA int    `validate:"max4"`
						LenB string `validate:"len:3"`
					} `validate:""`
				}{
					InInt: 25,
					NestedStruct: struct {
						MaxA int    `validate:"max4"`
						LenB string `validate:"len:3"`
					}{MaxA: 2, LenB: "abc"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid struct with invalid nested struct",
			args: args{
				v: struct {
					InInt        int `validate:"in:20,25,30"`
					NestedStruct struct {
						MaxA int    `validate:"max:4"`
						LenB string `validate:"len:3"`
					} `validate:""`
				}{
					InInt: 25,
					NestedStruct: struct {
						MaxA int    `validate:"max:4"`
						LenB string `validate:"len:3"`
					}{MaxA: 6, LenB: "abcdef"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)

				errs := []string{
					"field: NestedStruct, value: {6 abcdef}, error: field: MaxA, value: 6, error: shouldn't be great than 4",
					"field: NestedStruct, value: {6 abcdef}, error: field: LenB, value: abcdef, error: should has fixed length 3",
				}
				assert.True(t, err.Error() == strings.Join(errs, "\n"))

				return true
			},
		},
		{
			name: "valid struct with invalid nested struct",
			args: args{
				v: struct {
					InInt        int `validate:"in:20,25,30"`
					NestedStruct struct {
						MaxA int    `validate:"max:4"`
						LenB string `validate:"len:3"`
					}
				}{
					InInt: 25,
					NestedStruct: struct {
						MaxA int    `validate:"max:4"`
						LenB string `validate:"len:3"`
					}{MaxA: 6, LenB: "abcdef"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid struct with slices",
			args: args{
				v: struct {
					Ints    []int    `validate:"in:20,25,30"`
					Strings []string `validate:"max:5"`
				}{
					Ints:    []int{20, 30, 20},
					Strings: []string{"abcd", "abc", "ab", "a", ""},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid struct with invalid string slice",
			args: args{
				v: struct {
					Ints    []int    `validate:"in:20,25,30"`
					Strings []string `validate:"max:2"`
				}{
					Ints:    []int{20, 30, 20},
					Strings: []string{"abcd", "abc", "ab", "a", ""},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)
				return true
			},
		},
		{
			name: "invalid struct with invalid int slice",
			args: args{
				v: struct {
					Ints    []int    `validate:"in:25,30"`
					Strings []string `validate:"max:5"`
				}{
					Ints:    []int{20, 30, 20},
					Strings: []string{"abcd", "abc", "ab", "a", ""},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)
				return true
			},
		},
		{
			name: "invalid struct with invalid slices",
			args: args{
				v: struct {
					Ints    []int    `validate:"in:25,30"`
					Strings []string `validate:"max:2"`
				}{
					Ints:    []int{20, 30, 20},
					Strings: []string{"abcd", "abc", "ab", "a", ""},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 4)

				errs := []string{
					"field: Ints, value: 20, error: should be one of {25, 30}",
					"field: Ints, value: 20, error: should be one of {25, 30}",
					"field: Strings, value: abcd, error: shouldn't be longer than 2",
					"field: Strings, value: abc, error: shouldn't be longer than 2",
				}
				assert.True(t, err.Error() == strings.Join(errs, "\n"))

				return true
			},
		},
		{
			name: "valid struct with 2 validation rules for 1 field",
			args: args{
				v: struct {
					IntA    int      `validate:"min:2;max:6"`
					Ints    []int    `validate:"min:2;max:6"`
					StringB string   `validate:"min:2;max:6"`
					Strings []string `validate:"min:2;max:6"`
				}{
					IntA:    3,
					Ints:    []int{2, 3, 4, 5, 6},
					StringB: "abc",
					Strings: []string{"aa", "bbb", "cccc", "ddddd", "eeeeee"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid struct with 2 validation rules for 1 field",
			args: args{
				v: struct {
					IntA    int      `validate:"min:2;max:4"`
					Ints    []int    `validate:"min:2;max:4"`
					StringB string   `validate:"min:2;max:4"`
					Strings []string `validate:"min:2;max:4"`
				}{
					IntA:    5,
					Ints:    []int{5, 6},
					StringB: "abcde",
					Strings: []string{"ddddd", "eeeeee"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)

				errs := []string{
					"field: IntA, value: 5, error: shouldn't be great than 4",
					"field: Ints, value: 5, error: shouldn't be great than 4",
					"field: Ints, value: 6, error: shouldn't be great than 4",
					"field: StringB, value: abcde, error: shouldn't be longer than 4",
					"field: Strings, value: ddddd, error: shouldn't be longer than 4",
					"field: Strings, value: eeeeee, error: shouldn't be longer than 4",
				}
				assert.True(t, err.Error() == strings.Join(errs, "\n"))

				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
