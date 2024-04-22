//Copyright 2024 Institute of Distributed Systems, Ulm University
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package subjectivelogic

import (
	"testing"
)

var testOpinionsComult = []*Opinion{
	//valid Opinions
	{1, 0, 0, 0.5},
	{0, 1, 0, 0.5},
	{0, 0, 1, 0.5},
	{0.6, 0.3, 0.1, 0},
	{0.091, 0.604, 0.305, 0.4},
	{0.53, 0.227, 0.243, 1.000},
	{0.004, 0.5950000001, 0.4009999999, 0.00000000004334},
	{0, 0.999999999999999, 0.000000000000001, 0.5000000000000000000000000000000000000003},
	{0.000473, 0.555506, 0.444021, 0}}

func TestComultiplication(t *testing.T) {
	type args struct {
		opinion1 *Opinion
		opinion2 *Opinion
	}
	tests := []struct {
		name    string
		args    args
		want    Opinion
		wantErr bool
	}{
		//nil input
		{"TestComumtiplication1",
			args{nil, nil},
			Opinion{},
			true,
		},
		{"TestComumtiplication2",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},
		{"TestComumtiplication3",
			args{&Opinion{0, 1, 0, 0.5}, nil},
			Opinion{},
			true,
		},

		//general tests
		{"TestComumtiplication4",
			args{&Opinion{0.6, 0.3, 0.1, 0},
				&Opinion{0.000473, 0.555506, 0.444021, 0}},
			Opinion{},
			true,
		},
		{"TestComumtiplication5",
			args{&Opinion{0.53, 0.227, 0.243, 1.000},
				&Opinion{0.53, 0.227, 0.243, 1.000}},
			Opinion{},
			true,
		},
		{"TestComumtiplication6",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}},
			Opinion{1, 0, 0, 0.75},
			false,
		},
		{"TestComumtiplication7",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{1, 0, 0, 0.75},
			false,
		},
		{"TestComumtiplication8",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{0, 0.333333333333333, 0.666666666666666, 0.75},
			false,
		},
		{"TestComumtiplication9",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			Opinion{0.6, 0.4, 0, 0.5},
			false,
		},
		{"TestComumtiplication10",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			Opinion{0.6, 0, 0.4, 0.5},
			false,
		},
		{"TestComumtiplication11",
			args{&Opinion{0.6, 0.3, 0.1, 0}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			Opinion{0.6364, 0.2416, 0.122, 0.4},
			false,
		},
		{"TestComumtiplication12",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 1.000}},
			Opinion{0.57277, 0.178649, 0.248581, 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Comultiplication(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Comultiplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("Comultiplication() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkComultiplication(b *testing.B) {
	bmBinarySlFunc(Comultiplication, b)
}
