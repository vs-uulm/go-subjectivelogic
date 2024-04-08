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

func TestComplement(t *testing.T) {
	type args struct {
		opinion *Opinion
	}
	tests := []struct {
		name    string
		args    args
		want    *Opinion
		wantErr bool
	}{
		//nil input
		{"TestComplement1",
			args{nil},
			nil,
			true,
		},

		//general tests
		{"TestComplement2",
			args{&Opinion{1, 0, 0, 0.5}},
			&Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestComplement3",
			args{&Opinion{0, 1, 0, 0.5}},
			&Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestComplement4",
			args{&Opinion{0, 0, 1, 0.5}},
			&Opinion{0, 0, 1, 0.5},
			false,
		},
		{"TestComplement5",
			args{&Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0.3, 0.6, 0.1, 1},
			false,
		},
		{"TestComplement6",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}},
			&Opinion{0.604, 0.091, 0.305, 0.6},
			false,
		},
		{"TestComplement7",
			args{&Opinion{0.53, 0.227, 0.243, 1.000}},
			&Opinion{0.227, 0.53, 0.243, 0.000},
			false,
		},
		{"TestComplement8",
			args{&Opinion{0.004, 0.5950000001, 0.4009999999, 0.00000000004334}},
			&Opinion{0.5950000001, 0.004, 0.4009999999, 0.99999999995666},
			false,
		},
		{"TestComplement9",
			args{&Opinion{0, 0.999999999999999, 0.000000000000001, 0.5000000000000000000000000000000000000003}},
			&Opinion{0.999999999999999, 0, 0.000000000000001, 0.4999999999999999999999999999999999999997},
			false,
		},
		{"TestComplement10",
			args{&Opinion{0.000473, 0.555506, 0.444021, 0}},
			&Opinion{0.555506, 0.000473, 0.444021, 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Complement(tt.args.opinion)
			if (err != nil) != tt.wantErr {
				t.Errorf("Complement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("Complement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkComplement(b *testing.B) {
	opinion1, err := NewOpinion(1, 0, 0, 0)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for range b.N {
		x, err := Complement(&opinion1)
		if err != nil {
			b.Error(err)
		}
		sink = x
	}
}
