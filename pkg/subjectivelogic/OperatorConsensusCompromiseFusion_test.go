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

func TestConsensusCompromiseFusion(t *testing.T) {
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
		{"TestConsensusCompromiseFusion1",
			args{nil, nil},
			Opinion{},
			true,
		},
		{"TestConsensusCompromiseFusion2",
			args{nil, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},
		{"TestConsensusCompromiseFusion3",
			args{&Opinion{0, 1, 0, 0.5}, nil},
			Opinion{},
			true,
		},

		//null opinion input
		{"TestConsensusCompromiseFusion4",
			args{&Opinion{0, 0, 0, 0}, &Opinion{1, 0, 0, 0.5}},
			Opinion{},
			true,
		},

		//mismatched base rates
		{"TestConsensusCompromiseFusion5",
			args{&Opinion{0.5, 0.3, 0.2, 0.5}, &Opinion{0.5, 0.3, 0.2, 0.6}},
			Opinion{},
			true,
		},

		//u1 = u2 = 0, fully opposed opinions resolve to total uncertainty
		{"TestConsensusCompromiseFusion6",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 1, 0, 0.5}},
			Opinion{0, 0, 1, 0.5},
			false,
		},

		//u1 = u2 = 1
		{"TestConsensusCompromiseFusion7",
			args{&Opinion{0, 0, 1, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{0, 0, 1, 0.5},
			false,
		},

		//general tests
		{"TestConsensusCompromiseFusion8",
			args{&Opinion{1, 0, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{1, 0, 0, 0.5},
			false,
		},
		{"TestConsensusCompromiseFusion9",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0, 0, 1, 0.5}},
			Opinion{0, 1, 0, 0.5},
			false,
		},
		{"TestConsensusCompromiseFusion10",
			args{&Opinion{0, 1, 0, 0.5}, &Opinion{0.6, 0.3, 0.1, 0.5}},
			Opinion{0, 0.4, 0.6, 0.5},
			false,
		},
		{"TestConsensusCompromiseFusion11",
			args{&Opinion{0, 0, 1, 0.4}, &Opinion{0.6, 0.3, 0.1, 0.4}},
			Opinion{0.6, 0.3, 0.1, 0.4},
			false,
		},
		{"TestConsensusCompromiseFusion12",
			args{&Opinion{0.6, 0.3, 0.1, 0.4}, &Opinion{0.091, 0.604, 0.305, 0.4}},
			Opinion{0.3548491352338, 0.3516668086644, 0.2934840561018, 0.4},
			false,
		},
		{"TestConsensusCompromiseFusion13",
			args{&Opinion{0.091, 0.604, 0.305, 0.7}, &Opinion{0.53, 0.227, 0.243, 0.7}},
			Opinion{0.2991608404794, 0.3694237107969, 0.3314154487237, 0.7},
			false,
		},
		{"TestConsensusCompromiseFusion14",
			args{&Opinion{0.4, 0.0, 0.6, 0.5}, &Opinion{.7, .0, .3, .5}},
			Opinion{.82, .0, .18, .5},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConsensusCompromiseFusion(tt.args.opinion1, tt.args.opinion2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsensusCompromiseFusion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Compare(tt.want) {
				t.Errorf("ConsensusCompromiseFusion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConsensusCompromiseFusion(b *testing.B) {
	bmBinarySlFunc(ConsensusCompromiseFusion, b)
}