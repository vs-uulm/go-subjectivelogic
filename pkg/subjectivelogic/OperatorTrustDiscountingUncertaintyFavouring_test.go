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

func TestTrustDiscountingUncertaintyFavouring(t *testing.T) {
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
		{"TestTrustDiscountingUncertaintyFavouring1",
			args{nil, nil},
			Opinion{},
			true,
		},
		{"TestTrustDiscountingUncertaintyFavouring2",
			args{nil, &Opinion{0.6, 0.3, 0.1, 0.7}},
			Opinion{},
			true,
		},
		{"TestTrustDiscountingUncertaintyFavouring3",
			args{&Opinion{0.8, 0.1, 0.1, 0.5}, nil},
			Opinion{},
			true,
		},

		//null opinion input
		{"TestTrustDiscountingUncertaintyFavouring4",
			args{&Opinion{0, 0, 0, 0}, &Opinion{0.6, 0.3, 0.1, 0.7}},
			Opinion{},
			true,
		},

		//full trust in opinion1 leaves opinion2 unchanged (but for its base rate, which is always taken from opinion2)
		{"TestTrustDiscountingUncertaintyFavouring5",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0.7}},
			Opinion{0.6, 0.3, 0.1, 0.7},
			false,
		},

		//no trust in opinion1 (vacuous) yields total uncertainty about opinion2
		{"TestTrustDiscountingUncertaintyFavouring6",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0.7}},
			Opinion{0, 0, 1, 0.7},
			false,
		},

		//full distrust in opinion1 also yields total uncertainty about opinion2
		{"TestTrustDiscountingUncertaintyFavouring7",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0.7}},
			Opinion{0, 0, 1, 0.7},
			false,
		},

		//general test
		{"TestTrustDiscountingUncertaintyFavouring8",
			args{&Opinion{0.8, 0.1, 0.1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0.7}},
			Opinion{0.48, 0.24, 0.28, 0.7},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TrustDiscountingUncertaintyFavouring(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrustDiscountingUncertaintyFavouring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("TrustDiscountingUncertaintyFavouring() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTrustDiscountingUncertaintyFavouring(b *testing.B) {
	bmBinarySlFunc(TrustDiscountingUncertaintyFavouring, b)
}