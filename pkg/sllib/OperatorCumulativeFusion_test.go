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

package sllib

import (
	"testing"
)

func TestCumulativeFusion(t *testing.T) {
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
		{"TestCumulativeFusion1",
			args{nil, nil},
			nil,
			true,
		},
		{"TestCumulativeFusion2",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			nil,
			true,
		},
		{"TestCumulativeFusion3",
			args{&Opinion{0, 1, 0, 0.5}, nil},
			nil,
			true,
		},

		//u1 = u2 = 0
		{"TestCumulativeFusion4",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}},
			&Opinion{0.5, 0.5, 0, 0.5},
			false,
		},

		//u1 = u2 = 1
		{"TestCumulativeFusion5",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{0, 0, 1, 0.5},
			false,
		},

		//general tests
		{"TestCumulativeFusion6",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestCumulativeFusion7",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestCumulativeFusion8",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestCumulativeFusion9",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0.6, 0.3, 0.1, 0},
			false,
		},
		{"TestCumulativeFusion10",
			args{&Opinion{0.6, 0.3, 0.1, 0}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			&Opinion{0.5129506008011, 0.4056074766355, 0.08144192256342, 0.08081395348837},
			false,
		},
		{"TestCumulativeFusion11",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 1.000}},
			&Opinion{0.3877797355899, 0.4558215600831, 0.156398704327, 0.7465267528829},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CumulativeFusion(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CumulativeFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("CumulativeFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
