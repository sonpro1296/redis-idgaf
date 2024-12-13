package data_structures

import (
	"hash"
	"math"

	"github.com/spaolacci/murmur3"
)

const BucketFactor = 14

type HyperLogLog struct {
	buckets     []int
	hasher      hash.Hash64
	emptyBucket map[int]struct{}
}

func NewHyperLogLog() *HyperLogLog {
	numBuckets := 1 << BucketFactor
	// fmt.Println(numBuckets)
	emptyBucket := make(map[int]struct{}, numBuckets)
	for i := 0; i < numBuckets; i++ {
		emptyBucket[i] = struct{}{}
	}
	return &HyperLogLog{
		buckets:     make([]int, numBuckets),
		hasher:      murmur3.New64(),
		emptyBucket: emptyBucket,
	}
}

func (h *HyperLogLog) countLeadingZeros(hash uint64) int {
	// count := 0
	indicator := 1<<(64-BucketFactor-1) - 1

	// for i := 64 - BucketFactor - 1; i >= 0; i-- {
	// 	// fmt.Printf("%032b\n", indicator)
	// 	if int64(hash)&indicator != 0 {
	// 		break
	// 	}
	// 	count++
	// 	indicator >>= 1
	// }
	rawResult := hash & uint64(indicator)

	return 64 - BucketFactor - 1 - int(math.Floor(math.Log2(float64(rawResult))))
}

func (h *HyperLogLog) correct(estimate float64) int64 {
	if float64(estimate) <= 2.5*float64(len(h.buckets)) {
		return int64(float64(estimate) * math.Log(float64(len(h.buckets))/float64(len(h.emptyBucket))))
	}
	if estimate > 10_000_000 {
		return int64(-math.Pow(2, -32) * math.Log(1-float64(estimate)*math.Pow(2, -32)))
	}
	return int64(estimate)

}

func (h *HyperLogLog) Add(key []byte) {
	h.hasher.Write(key)
	hashResult := h.hasher.Sum64()
	defer h.hasher.Reset()

	bucketIndex := hashResult >> (64 - BucketFactor)

	delete(h.emptyBucket, int(bucketIndex))
	leadingZeros := h.countLeadingZeros(hashResult)

	h.buckets[bucketIndex] = max(h.buckets[bucketIndex], leadingZeros)
}

func (h *HyperLogLog) Count() int64 {
	// calculate harmonic mean

	alpha := 0.7213 / (1 + 1.079/float64(len(h.buckets)))
	var harmonicMean float64
	// fmt.Println(len(h.emptyBucket))
	for _, count := range h.buckets {
		// fmt.Println(count)
		harmonicMean += 1 / math.Pow(2.0, float64(count))
	}
	// harmonicMean
	// fmt.Println(harmonicMean)

	// calculate estimate
	estimate := alpha * float64(len(h.buckets)) * (float64(len(h.buckets) - len(h.emptyBucket))) / harmonicMean
	// fmt.Println(estimate)
	return int64(estimate)

}
