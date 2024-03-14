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

func TestTrustDiscounting(t *testing.T) {
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
		{"TestTrustDiscounting",
			args{nil, nil},
			nil,
			true},
		{"TestTrustDiscounting",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			nil,
			true,
		},
		{"TestTrustDiscounting",
			args{&Opinion{1, 0, 0, 0.5}, nil},
			nil,
			true,
		},

		//general testing
		{"TestTrustDiscounting",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}},
			&Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			&Opinion{0, 0, 1, 0.5},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0, 0, 1, 0},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0.3, 0.15, 0.55, 0},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{0.6, 0.3, 0.1, 0}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			&Opinion{0.0546, 0.3624, 0.583, 0.4},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.6, 0.3, 0.1, 0}},
			&Opinion{0.1278, 0.0639, 0.8083, 0},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 1.000}},
			&Opinion{0.11289, 0.048351, 0.838759, 1},
			false,
		},
		{"TestTrustDiscounting",
			args{&Opinion{0.53, 0.227, 0.243, 1.000}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			&Opinion{0.070343, 0.466892, 0.462765, 0.4},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TrustDiscounting(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrustDiscounting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("TrustDiscounting() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiEdgeTrustDisc(t *testing.T) {
	type args struct {
		opinions []*Opinion
	}
	tests := []struct {
		name    string
		args    args
		want    *Opinion
		wantErr bool
	}{
		//nil input
		{"TestMultiEdgeTrustDisc1",
			args{nil},
			nil,
			true,
		},
		{"TestMultiEdgeTrustDisc2",
			args{[]*Opinion{}},
			nil,
			true,
		},
		{"TestMultiEdgeTrustDisc3",
			args{[]*Opinion{nil}},
			nil,
			true,
		},
		{"TestMultiEdgeTrustDisc4",
			args{[]*Opinion{nil, nil}},
			nil,
			true,
		},
		{"TestMultiEdgeTrustDisc5",
			args{[]*Opinion{nil, &Opinion{1, 0, 0, 0.5}}},
			nil,
			true,
		},
		{"TestMultiEdgeTrustDisc6",
			args{[]*Opinion{&Opinion{1, 0, 0, 0.5}, nil}},
			nil,
			true,
		},
		{"TestMultiEdgeTrustDisc7",
			args{[]*Opinion{&Opinion{1, 0, 0, 0.5}, nil, &Opinion{0, 1, 0, 0.5}}},
			nil,
			true,
		},

		//1 argument
		{"TestMultiEdgeTrustDisc8",
			args{[]*Opinion{&Opinion{1, 0, 0, 0.5}}},
			nil,
			true,
		},

		//general tests
		{"TestMultiEdgeTrustDisc9",
			args{[]*Opinion{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}}},
			&Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestMultiEdgeTrustDisc10",
			args{[]*Opinion{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}}},
			&Opinion{0, 0, 1, 0.5},
			false,
		},
		{"TestMultiEdgeTrustDisc11",
			args{[]*Opinion{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}}},
			&Opinion{0, 0, 1, 0},
			false,
		},
		{"TestMultiEdgeTrustDisc12",
			args{[]*Opinion{&Opinion{0, 0, 1, 0.5}, &Opinion{0.6, 0.3, 0.1, 0}}},
			&Opinion{0.3, 0.15, 0.55, 0},
			false,
		},
		{"TestMultiEdgeTrustDisc13",
			args{[]*Opinion{&Opinion{0.6, 0.3, 0.1, 0}, &Opinion{0.091, 0.604, 0.305, 0.4}}},
			&Opinion{0.0546, 0.3624, 0.583, 0.4},
			false,
		},
		{"TestMultiEdgeTrustDisc14",
			args{[]*Opinion{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.6, 0.3, 0.1, 0}}},
			&Opinion{0.1278, 0.0639, 0.8083, 0},
			false,
		},
		{"TestMultiEdgeTrustDisc15",
			args{[]*Opinion{&Opinion{0.091, 0.604, 0.305, 0.4}, &Opinion{0.53, 0.227, 0.243, 1.000}}},
			&Opinion{0.11289, 0.048351, 0.838759, 1},
			false,
		},
		{"TestMultiEdgeTrustDisc16",
			args{[]*Opinion{&Opinion{0.53, 0.227, 0.243, 1.000}, &Opinion{0.091, 0.604, 0.305, 0.4}}},
			&Opinion{0.070343, 0.466892, 0.462765, 0.4},
			false,
		},

		//3 inputs
		{"TestMultiEdgeTrustDisc17",
			args{[]*Opinion{&Opinion{0.53, 0.227, 0.243, 1.000},
				&Opinion{0.6, 0.3, 0.1, 0},
				&Opinion{0.091, 0.604, 0.305, 0.4}}},
			&Opinion{0.0422058, 0.2801352, 0.677659, 0.4},
			false,
		},

		//5 inputs
		{"TestMultiEdgeTrustDisc18",
			args{[]*Opinion{&Opinion{0, 0, 1, 0.5},
				&Opinion{0.53, 0.227, 0.243, 1.000},
				&Opinion{0.6, 0.3, 0.1, 0},
				&Opinion{0.5, 0.5, 0, 0.5},
				&Opinion{0.6, 0.3, 0.1, 0}}},
			&Opinion{0.06957, 0.034785, 0.895645, 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MultiEdgeTrustDisc(tt.args.opinions)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiEdgeTrustDisc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("MultiEdgeTrustDisc() got = %v, want %v", got, tt.want)
			}
		})
	}
}
