// Copyright 2016 mshk.top, lion@mshk.top
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commonstock

import (
	"goweb/pkg/commonmodels"
	"math"
)

type Boll struct {
	n int
	k float64
}

//NewBoll(20, 2)
func NewBoll(n int, k int32) *Boll {
	return &Boll{n: n, k: float64(k)}
}
func (this *Boll) sma(lines []*commonmodels.Kline) float64 {
	s := len(lines)
	var sum float64 = 0
	for i := 0; i < s; i++ {
		sum += float64(lines[i].Close)
	}
	return sum / float64(s)
}
func (this *Boll) dma(lines []*commonmodels.Kline, ma float64) float64 {
	s := len(lines)
	//log.Println(s)
	var sum float64 = 0
	for i := 0; i < s; i++ {
		sum += (lines[i].Close - ma) * (lines[i].Close - ma)
	}
	return math.Sqrt(sum / float64(this.n))
}

func (this *Boll) Boll(lines []*commonmodels.Kline) (mid []float64, up []float64, low []float64) {
	l := len(lines)

	mid = make([]float64, l)
	up = make([]float64, l)
	low = make([]float64, l)
	if l < this.n {
		return
	}
	for i := l - 1; i > this.n-1; i-- {
		ps := lines[(i - this.n + 1) : i+1]
		mid[i] = this.sma(ps)
		dm := this.k * this.dma(ps, mid[i])
		up[i] = mid[i] + dm
		low[i] = mid[i] - dm
	}

	return
}
