package mock

import "testing"

func TestParametrizedBodyComparator_Compare(t *testing.T) {
	type args struct {
		mock    []byte
		request []byte
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"Json parsing error. Request",
			args{
				mock:    []byte(`{"name": "alice"}`),
				request: []byte(`{"name": "alice"`),
			},
			false,
			true,
		},
		{
			"Json parsing error. Mock",
			args{
				mock:    []byte(`{"name": "alice"`),
				request: []byte(`{"name": "alice"}`),
			},
			false,
			true,
		},
		{
			"Single field without template with same names",
			args{
				mock:    []byte(`{"name": "alice"}`),
				request: []byte(`{"name": "alice"}`),
			},
			true,
			false,
		},
		{
			"Single field without template with different names",
			args{
				mock:    []byte(`{"name": "alice"}`),
				request: []byte(`{"name": "bob"}`),
			},
			false,
			false,
		},
		{
			"Single field with same type and order arrays",
			args{
				mock:    []byte(`{"name": [1,2,3]}`),
				request: []byte(`{"name": [1,2,3]}`),
			},
			true,
			false,
		},
		{
			"Single Field. Different len arrays",
			args{
				mock:    []byte(`{"name": [1,2,3]}`),
				request: []byte(`{"name": [1,2]}`),
			},
			false,
			false,
		},
		{
			"Root arrs. Different len arrays",
			args{
				mock:    []byte(`[1,2,3]`),
				request: []byte(`[1,2]`),
			},
			false,
			false,
		},
		{
			"Root arrs. Same arrays",
			args{
				mock:    []byte(`[1,2,3]`),
				request: []byte(`[1,2,3]`),
			},
			true,
			false,
		},
		{
			"Root arrs. Different arrays",
			args{
				mock:    []byte(`[1,2,3]`),
				request: []byte(`[1,2,4]`),
			},
			false,
			false,
		},
		{
			"Single field with same type and different order",
			args{
				mock:    []byte(`{"name": [1,2,3]}`),
				request: []byte(`{"name": [2,1,3]}`),
			},
			false,
			false,
		},
		{
			"Single field with different",
			args{
				mock:    []byte(`{"name": [false,true,false]}`),
				request: []byte(`{"name": [2,1,3]}`),
			},
			false,
			false,
		},
		{
			"Single level with template",
			args{
				mock:    []byte(`{"name": "{name}"}`),
				request: []byte(`{"name": "bob"}`),
			},
			true,
			false,
		},
		{
			"Two levels without template. Same",
			args{
				mock:    []byte(`{"name": {"name":"name", "arr": [1,2,3]}}`),
				request: []byte(`{"name": {"name":"name", "arr": [1,2,3]}}`),
			},
			true,
			false,
		},
		{
			"Two levels without template. Different in single value",
			args{
				mock:    []byte(`{"name": {"name": false, "arr": [1,2,3]}}`),
				request: []byte(`{"name": {"name": true, "arr": [1,2,3]}}`),
			},
			false,
			false,
		},
		{
			"Two levels without template. Different in arrays",
			args{
				mock:    []byte(`{"name": {"name": false, "arr": [1,2,3]}}`),
				request: []byte(`{"name": {"name": true, "arr": [1,3,3]}}`),
			},
			false,
			false,
		},
		{
			"Two levels with template. In value",
			args{
				mock:    []byte(`{"name": {"name": "{name}", "arr": [1,2,3]}}`),
				request: []byte(`{"name": {"name": true, "arr": [1,2,3]}}`),
			},
			true,
			false,
		},
		{
			"Two levels with template. In arr",
			args{
				mock:    []byte(`{"name": {"name": true, "arr": "{arr}"}}`),
				request: []byte(`{"name": {"name": true, "arr": [1,3,3]}}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. != op. Strings. False",
			args{
				mock:    []byte(`{"name": "{ name != alice }"}`),
				request: []byte(`{"name": "alice"}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. != op. Strings. True",
			args{
				mock:    []byte(`{"name": "{ name != alice }"}`),
				request: []byte(`{"name": "Bob"}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. != op. Int. True",
			args{
				mock:    []byte(`{"name": "{ name != 13 }"}`),
				request: []byte(`{"name": 12}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. != op. Int. False",
			args{
				mock:    []byte(`{"name": "{ name != 13 }"}`),
				request: []byte(`{"name": 123}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. > op. String. False",
			args{
				mock:    []byte(`{"name": "{ name > 13 }"}`),
				request: []byte(`{"name": "123"}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. >= op. String. False",
			args{
				mock:    []byte(`{"name": "{ name >= 13 }"}`),
				request: []byte(`{"name": "123"}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. <= op. String. False",
			args{
				mock:    []byte(`{"name": "{ name <= 13 }"}`),
				request: []byte(`{"name": "123"}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. = op. Undefined operator.",
			args{
				mock:    []byte(`{"name": "{ name = 13 }"}`),
				request: []byte(`{"name": 123}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. != op. Objects. Undefined type.",
			args{
				mock:    []byte(`{"name": "{ name != 13 }"}`),
				request: []byte(`{"name": { "name": "123" }}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. > op. Int. True",
			args{
				mock:    []byte(`{"name": "{ name > 13 }"}`),
				request: []byte(`{"name": 14}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. > op. Int. False",
			args{
				mock:    []byte(`{"name": "{ name > 13 }"}`),
				request: []byte(`{"name": 12}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. >= op. Int. True",
			args{
				mock:    []byte(`{"name": "{ name >= 13 }"}`),
				request: []byte(`{"name": 13}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. >= op. Int. False",
			args{
				mock:    []byte(`{"name": "{ name >= 13 }"}`),
				request: []byte(`{"name": 12}`),
			},
			false,
			false,
		},
		{
			"Expression interpreter. <= op. Int. True",
			args{
				mock:    []byte(`{"name": "{ name <= 13 }"}`),
				request: []byte(`{"name": 13}`),
			},
			true,
			false,
		},
		{
			"Expression interpreter. <= op. Int. False",
			args{
				mock:    []byte(`{"name": "{ name <= 13 }"}`),
				request: []byte(`{"name": 14}`),
			},
			false,
			false,
		},
		{
			"Objects with different fields count",
			args{
				mock:    []byte(`{"name": "123"}`),
				request: []byte(`{"name": 14, "field": 1}`),
			},
			false,
			false,
		},
		{
			"Objects with different nested objects. Request have different type",
			args{
				mock:    []byte(`{"name": {"name": ""}}`),
				request: []byte(`{"name": ""}`),
			},
			false,
			false,
		},
		{
			"Root Array With Nested Array. Same",
			args{
				mock:    []byte(`[[1,2,3],[1,2,3]]`),
				request: []byte(`[[1,2,3],[1,2,3]]`),
			},
			true,
			false,
		},
		{
			"Root Array With Nested Array. Different",
			args{
				mock:    []byte(`[[1,2,4],[1,2,3]]`),
				request: []byte(`[[1,2,3],[1,2,3]]`),
			},
			false,
			false,
		},
		{
			"Root Array With Objects. Same",
			args{
				mock:    []byte(`[{"name": ""}, {"name": ""}]`),
				request: []byte(`[{"name": ""}, {"name": ""}]`),
			},
			true,
			false,
		},
		{
			"Root Array With Objects. Different",
			args{
				mock:    []byte(`[{"name": ""}, {"name": ""}]`),
				request: []byte(`[{"name": ""}, {"name": "13"}]`),
			},
			false,
			false,
		},
		{
			"Root Array With Objects. DifferentTypes",
			args{
				mock:    []byte(`[{"name": ""}, {"name": ""}]`),
				request: []byte(`[1, 1]`),
			},
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmp := &ParametrizedBodyComparator{}
			got, err := cmp.Compare(tt.args.mock, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Compare() got = %v, want %v", got, tt.want)
			}
		})
	}
}
