package main

import (
	"fmt"
	"sonpro1296/redis-idgaf/data_structures"
	"github.com/axiomhq/hyperloglog"
)

func main() {
	// skipList := data_structures.NewSkiplist(8, 0.5)
	// for i := 0; i < 10; i++ {
	// 	skipList.Add(fmt.Sprintf("abc%d", i), (float64(i)))
	// }

	// skipList.Display()
	// skipList.Delete("abc1")
	// skipList.Delete("abc4")
	// skipList.Display()

	hll := data_structures.NewHyperLogLog()
	for i := 0; i < 10_000_000; i++ {
		hll.Add([]byte(fmt.Sprintf("abc%d", i)))
	}
	fmt.Println(hll.Count())

	axiomHll := hyperloglog.NewNoSparse()
	for i := 0; i < 10_000_000; i++ {
		axiomHll.Insert([]byte(fmt.Sprintf("abc%d", i)))
	}
	fmt.Println(axiomHll.Estimate())



}
