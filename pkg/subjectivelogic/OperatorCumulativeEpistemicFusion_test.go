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

func TestCumulativeEpistemicFusion(t *testing.T) {
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
		{"TestCumulativeEpistemicFusion1",
			args{nil, nil},
			Opinion{},
			true,
		},
		{"TestCumulativeEpistemicFusion2",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},
		{"TestCumulativeEpistemicFusion3",
			args{&Opinion{0, 1, 0, 0.5}, nil},
			Opinion{},
			true,
		},

		//null opinion input
		{"TestCumulativeEpistemicFusion4",
			args{&Opinion{0, 0, 0, 0}, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},

		//fused base rate is exactly 0, which is undefined for epistemic fusion
		{"TestCumulativeEpistemicFusion5",
			args{&Opinion{0.5, 0.3, 0.2, 0}, &Opinion{0.5, 0.3, 0.2, 0}},
			Opinion{},
			true,
		},

		//fused base rate is exactly 1, which is undefined for epistemic fusion
		{"TestCumulativeEpistemicFusion6",
			args{&Opinion{0.5, 0.3, 0.2, 1}, &Opinion{0.5, 0.3, 0.2, 1}},
			Opinion{},
			true,
		},

		//u1 = u2 = 0, fully opposed opinions: cumulative fusion averages to
		//b=d=0.5, and the maximum-uncertainty projection collapses that to
		//total uncertainty at the same projected probability
		{"TestCumulativeEpistemicFusion7",
			args{&Opinion{1, 0, 0, 0.3}, &Opinion{0, 1, 0, 0.7}},
			Opinion{0, 0, 1, 0.5},
			false,
		},

		//general tests
		{"TestCumulativeEpistemicFusion8",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			Opinion{0.2425456163774, 0, 0.7574543836226, 0.4},
			false,
		},
		{"TestCumulativeEpistemicFusion9",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 0.6}},
			Opinion{0, 0.0913742781960, 0.9086257218040, 0.5155089176276},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CumulativeEpistemicFusion(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CumulativeEpistemicFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("CumulativeEpistemicFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCumulativeEpistemicFusion(b *testing.B) {
	bmBinarySlFunc(CumulativeEpistemicFusion, b)
}