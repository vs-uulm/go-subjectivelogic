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

var testOpinionsMult = []*Opinion{
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

var expectedOpinionsMult = []*Opinion{
	{0, 1, 0, 0.25},
	{0.333333333333333, 0, 0.666666666666666, 0.25},
	{0, 1, 0, 0.25},
	{0, 1, 0, 0},
	{0.225, 0.3, 0.425, 0.2},
	{0.1278, 0.7228, 0.1494, 0},
	{0.070343, 0.693892, 0.235765, 0.4}}

func TestMultiplication(t *testing.T) {
	type args struct {
		opinion1 *Opinion
		opinion2 *Opinion
	}
	tests := []struct {
		name    string
		args    args
		want    *Opinion
		wantErr bool
	}{
		//nil input
		{"TestMultiplication1",
			args{nil, nil},
			nil,
			true,
		},
		{"TestMultiplication2",
			args{nil, testOpinionsMult[1]},
			nil,
			true,
		},
		{"TestMultiplication3",
			args{testOpinionsMult[0], nil},
			nil,
			true,
		},

		//general tests
		{"TestMultiplication4",
			args{testOpinionsMult[0], testOpinionsMult[1]},
			&Opinion{0, 1, 0, 0.25},
			false,
		},
		{"TestMultiplication5",
			args{testOpinionsMult[0], testOpinionsMult[2]},
			&Opinion{0.333333333333333, 0, 0.666666666666666, 0.25},
			false,
		},
		{"TestMultiplication6",
			args{testOpinionsMult[1], testOpinionsMult[2]},
			&Opinion{0, 1, 0, 0.25},
			false,
		},
		{"TestMultiplication7",
			args{testOpinionsMult[1], testOpinionsMult[3]},
			&Opinion{0, 1, 0, 0},
			false,
		},
		{"TestMultiplication8",
			args{testOpinionsMult[2], testOpinionsMult[3]},
			&Opinion{0.3, 0.3, 0.4, 0},
			false,
		},
		{"TestMultiplication9",
			args{testOpinionsMult[3], testOpinionsMult[4]},
			&Opinion{0.1278, 0.7228, 0.1494, 0},
			false,
		},
		{"TestMultiplication10",
			args{testOpinionsMult[4], testOpinionsMult[5]},
			&Opinion{0.070343, 0.693892, 0.235765, 0.4},
			false,
		},
		{"TestMultiplication11",
			args{testOpinionsMult[5], testOpinionsMult[4]},
			&Opinion{0.070343, 0.693892, 0.235765, 0.4},
			false,
		},
		{"TestMultiplication11",
			args{testOpinionsMult[5], testOpinionsMult[5]},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Multiplication(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Multiplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("Multiplication() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMultiplication(b *testing.B) {
	bmBinarySlFunc(Multiplication, b)
}
