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

func TestTrustDiscountingOppositeBelief(t *testing.T) {
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
		{"TestTrustDiscountingOppositeBelief",
			args{nil, nil},
			Opinion{},
			true},
		{"TestTrustDiscountingOppositeBelief",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{1, 0, 0, 0.5}, nil},
			Opinion{},
			true,
		},

		//general testing
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{1, 0, 0, 0.5}},
			Opinion{0, 0, 1, 0.5},
			false,
		},
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{1, 0, 0, 0.5}},
			Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{1, 0, 0, 0.5}},
			Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{0, 0.75, 0.25, 0.5}, &Opinion{1, 0, 0, 0.5}},
			Opinion{0, 0.75, 0.25, 0.5},
			false,
		},
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{0.3, 0.7, 0, 0.5}, &Opinion{1, 0, 0, 0.5}},
			Opinion{0.3, 0.7, 0, 0.5},
			false,
		},		
		{"TestTrustDiscountingOppositeBelief",
			args{&Opinion{0.2, 0.4, 0.4, 0.5}, &Opinion{1, 0, 0, 0.5}},
			Opinion{0.2, 0.4, 0.4, 0.5},
			false,
		},								
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TrustDiscountingOppositeBelief(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrustDiscountingOB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("TrustDiscountingOB() got = %v, want %v", got, tt.want)
			}
		})
	}
}



func BenchmarkTrustDiscountingOppositeBelief(b *testing.B) {
	bmBinarySlFunc(TrustDiscountingOppositeBelief, b)
}
