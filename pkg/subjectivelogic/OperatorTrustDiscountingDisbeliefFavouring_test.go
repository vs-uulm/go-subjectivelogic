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

func TestTrustDiscountingDisbeliefFavouring(t *testing.T) {
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
		{"TestTrustDiscountingDisbeliefFavouring1",
			args{nil, nil},
			Opinion{},
			true,
		},
		{"TestTrustDiscountingDisbeliefFavouring2",
			args{nil, &Opinion{0.2, 0.5, 0.3, 0.6}},
			Opinion{},
			true,
		},
		{"TestTrustDiscountingDisbeliefFavouring3",
			args{&Opinion{0.5, 0.3, 0.2, 0.4}, nil},
			Opinion{},
			true,
		},

		//null opinion input
		{"TestTrustDiscountingDisbeliefFavouring4",
			args{&Opinion{0, 0, 0, 0}, &Opinion{0.2, 0.5, 0.3, 0.6}},
			Opinion{},
			true,
		},

		//both belief mass and uncertainty mass are zero (both opinions are total
		//disbelief, so d1 == d2 and the tie-break rule takes a1)
		{"TestTrustDiscountingDisbeliefFavouring5",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0, 1, 0, 0.7}},
			Opinion{0, 1, 0, 0.5},
			false,
		},

		//general test, disbelief of opinion2 dominates, so its base rate is used
		{"TestTrustDiscountingDisbeliefFavouring6",
			args{&Opinion{0.5, 0.3, 0.2, 0.4}, &Opinion{0.2, 0.5, 0.3, 0.6}},
			Opinion{0.2916666666667, 0.5000000000000, 0.2083333333333, 0.6},
			false,
		},

		//general test, disbelief of opinion1 dominates, so its base rate is used
		{"TestTrustDiscountingDisbeliefFavouring7",
			args{&Opinion{0.2, 0.6, 0.2, 0.4}, &Opinion{0.5, 0.2, 0.3, 0.6}},
			Opinion{0.2333333333333, 0.6000000000000, 0.1666666666667, 0.4},
			false,
		},

		//disbelief is tied between opinion1 and opinion2, so the tie-break
		//rule takes opinion1's base rate
		{"TestTrustDiscountingDisbeliefFavouring8",
			args{&Opinion{0.5, 0.3, 0.2, 0.4}, &Opinion{0.3, 0.3, 0.4, 0.6}},
			Opinion{0.4, 0.3, 0.3, 0.4},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TrustDiscountingDisbeliefFavouring(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrustDiscountingDisbeliefFavouring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("TrustDiscountingDisbeliefFavouring() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTrustDiscountingDisbeliefFavouring(b *testing.B) {
	bmBinarySlFunc(TrustDiscountingDisbeliefFavouring, b)
}