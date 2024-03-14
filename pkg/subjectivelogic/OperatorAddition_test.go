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

func TestAddition(t *testing.T) {
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

		//general tests
		{"TestWeightedFusion6",
			args{&Opinion{0.5, 0.3, 0.2, 0.5}, &Opinion{0, 0.7, 0.3, 0.2}},
			&Opinion{0.5, 0.2714285714286, 0.2285714285714, 0.7},
			false,
		},
		{"TestWeightedFusion7",
			args{&Opinion{0.3, 0.2, 0.5, 0.3}, &Opinion{0.1, 0.2, 0.7, 0.1}},
			&Opinion{0.4, 0.05, 0.55, 0.4},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Addition(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Addition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("Addition() got = %v, want %v", got, tt.want)
			}
		})
	}
}
