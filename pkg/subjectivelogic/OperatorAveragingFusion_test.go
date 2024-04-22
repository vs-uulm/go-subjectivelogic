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

func TestAveragingFusion(t *testing.T) {
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
		{"TestAveragingFusion1",
			args{nil, nil},
			Opinion{},
			true,
		},
		{"TestAveragingFusion",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},
		{"TestAveragingFusion3",
			args{&Opinion{0, 1, 0, 0.5}, nil},
			Opinion{},
			true,
		},

		//u1 = u2 = 0
		{"TestAveragingFusion4",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}},
			Opinion{0.5, 0.5, 0, 0.5},
			false,
		},

		//general tests
		{"TestAveragingFusion5",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestAveragingFusion6",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestAveragingFusion7",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			Opinion{0, 1, 0, 0.25},
			false,
		},
		{"TestAveragingFusion8",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			Opinion{0.545454545454545, 0.272727272727272, 0.181818181818181, 0.25},
			false,
		},
		{"TestAveragingFusion9",
			args{&Opinion{0.6, 0.3, 0.1, 0}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			Opinion{0.4743209876543, 0.3750617283951, 0.1506172839506, 0.2},
			false,
		},
		{"TestAveragingFusion10",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 1.000}},
			Opinion{0.3353339416058, 0.3941733576642, 0.2704927007299, 0.7},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AveragingFusion(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("AveragingFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("AveragingFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAveragingFusion(b *testing.B) {
	bmBinarySlFunc(AveragingFusion, b)
}
