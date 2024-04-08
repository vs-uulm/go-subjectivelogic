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

func TestWeightedFusion(t *testing.T) {
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
		{"TestWeightedFusion1",
			args{nil, nil},
			nil,
			true,
		},
		{"TestWeightedFusion2",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			nil,
			true,
		},
		{"TestWeightedFusion3",
			args{&Opinion{0, 1, 0, 0.5}, nil},
			nil,
			true,
		},

		//u1 = u2 = 0
		{"TestWeightedFusion4",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}},
			&Opinion{0.5, 0.5, 0, 0.5},
			false,
		},

		//u1 = u2 = 1
		{"TestWeightedFusion5",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{0, 0, 1, 0.5},
			false,
		},

		//general tests
		{"TestWeightedFusion6",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestWeightedFusion7",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestWeightedFusion8",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0, 1, 0, 0.2631578947368},
			false,
		},
		{"TestWeightedFusion9",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0.6, 0.3, 0.1, 0},
			false,
		},
		{"TestWeightedFusion10",
			args{&Opinion{0.6, 0.3, 0.1, 0}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			&Opinion{0.497164244186, 0.3614186046512, 0.1414171511628, 0.1742946708464},
			false,
		},
		{"TestWeightedFusion11",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 1.000}},
			&Opinion{0.3445420741927, 0.3862656902719, 0.2691922355354, 0.7128099173554},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WeightedFusion(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("WeightedFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.ComparePtr(tt.want) {
				t.Errorf("WeightedFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkWeightedFusion(b *testing.B) {
	bmBinarySlFunc(WeightedFusion, b)
}
