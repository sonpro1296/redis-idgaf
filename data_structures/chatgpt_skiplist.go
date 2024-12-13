package data_structures

import (
	"math"
	"math/rand"
	"time"
)

// HyperLogLog structure with 128 buckets
type CHyperLogLog struct {
	buckets []uint8 // Array of buckets to store the max number of leading zeros
	m       int     // Number of buckets (128 in this case)
}

// NewHyperLogLog initializes a new HyperLogLog with 128 buckets
func NewCHyperLogLog() *CHyperLogLog {
	return &CHyperLogLog{
		buckets: make([]uint8, 128), // 128 buckets
		m:       128,                // Set the number of buckets
	}
}

// Hash function (simple hash to simulate hashing in HyperLogLog)
func chash(x string) uint32 {
	// For simplicity, we use a random hash here (you should use a real hash function like MurmurHash)
	rand.Seed(time.Now().UnixNano())
	return rand.Uint32()
}

// Add method to add an element to the HyperLogLog structure
func (h *CHyperLogLog) Add(element string) {
	// Hash the element
	hashValue := chash(element)

	// Find the bucket index (use the first 7 bits of the hash for 128 buckets)
	bucketIndex := hashValue >> (32 - 7) // 7 bits for 128 buckets (2^7 = 128)
	// Find the number of leading zeros in the remaining part of the hash value
	remainingBits := hashValue & ((1 << (32 - 7)) - 1)
	leadingZeros := countLeadingZeros(remainingBits)

	// Update the bucket with the max number of leading zeros
	if uint8(leadingZeros) > h.buckets[bucketIndex] {
		h.buckets[bucketIndex] = uint8(leadingZeros)
	}
}

// Count the number of leading zeros in a 32-bit number
func countLeadingZeros(x uint32) int {
	// The leading zeros in the hash can be found using a simple loop
	for i := 31; i >= 0; i-- {
		if (x>>i)&1 == 1 {
			return 31 - i
		}
	}
	return 32 // If no 1s, return 32 (max leading zeros)
}

// Estimate cardinality (distinct element count) using the HyperLogLog algorithm
func (h *CHyperLogLog) Estimate() float64 {
	// Compute the harmonic mean of the leading zero counts in the buckets
	Z := 0.0
	for _, count := range h.buckets {
		Z += 1.0 / math.Pow(2, float64(count))
	}

	// Apply the correction factor
	alphaM := 0.7213 / (1 + 1.079/float64(h.m))
	estimate := alphaM * float64(h.m*h.m) / Z
	return estimate
}

// func main() {
// 	// Create a new HyperLogLog structure with 128 buckets
// 	hll := NewHyperLogLog()

// 	// Add some elements to the HyperLogLog
// 	elements := []string{"apple", "banana", "cherry", "apple", "date", "banana"}
// 	for _, element := range elements {
// 		hll.Add(element)
// 	}

// 	// Estimate the cardinality (distinct elements)
// 	estimate := hll.Estimate()
// 	fmt.Printf("Estimated number of distinct elements: %f\n", estimate)
// }