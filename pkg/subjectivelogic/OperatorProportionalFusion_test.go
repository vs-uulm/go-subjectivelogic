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

func TestProportionalFusion(t *testing.T) {
	type args struct {
		opinion1    *Opinion
		opinion2    *Opinion
		proportion1 float64
		proportion2 float64
	}
	tests := []struct {
		name    string
		args    args
		want    Opinion
		wantErr bool
	}{
		//nil input
		{"TestProportionalFusion1",
			args{nil, nil, 0.5, 0.5},
			Opinion{},
			true,
		},
		{"TestProportionalFusion2",
			args{nil, &Opinion{1, 0, 0, 0.5}, 0.5, 0.5},
			Opinion{},
			true,
		},
		{"TestProportionalFusion3",
			args{&Opinion{0, 1, 0, 0.5}, nil, 0.5, 0.5},
			Opinion{},
			true,
		},

		//null opinion input
		{"TestProportionalFusion4",
			args{&Opinion{0, 0, 0, 0}, &Opinion{1, 0, 0, 0.5}, 0.5, 0.5},
			Opinion{},
			true,
		},

		//proportions don't sum to 1
		{"TestProportionalFusion5",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.2, 0.5, 0.3, 0.6}, 0.5, 0.6},
			Opinion{},
			true,
		},

		//equal weighting reduces to the arithmetic mean
		{"TestProportionalFusion6",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.2, 0.5, 0.3, 0.6}, 0.5, 0.5},
			Opinion{0.4, 0.4, 0.2, 0.5},
			false,
		},

		//weight 1 on opinion1 leaves it unchanged
		{"TestProportionalFusion7",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.091, 0.604, 0.305, 0.7}, 1.0, 0.0},
			Opinion{0.6, 0.3, 0.1, 0.4},
			false,
		},

		//weight 1 on opinion2 leaves it unchanged
		{"TestProportionalFusion8",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.091, 0.604, 0.305, 0.7}, 0.0, 1.0},
			Opinion{0.091, 0.604, 0.305, 0.7},
			false,
		},

		//general test, unequal weighting
		{"TestProportionalFusion9",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.091, 0.604, 0.305, 0.7}, 0.3, 0.7},
			Opinion{0.2437, 0.5128, 0.2435, 0.61},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProportionalFusion(tt.args.opinion1, tt.args.opinion2, tt.args.proportion1, tt.args.proportion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProportionalFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("ProportionalFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkProportionalFusion(b *testing.B) {
	opinion1 := &Opinion{0.6, 0.3, 0.1, 0.4}
	opinion2 := &Opinion{0.091, 0.604, 0.305, 0.7}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ProportionalFusion(opinion1, opinion2, 0.3, 0.7)
	}
}