/**
 * Package nicTest
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/18 13:36
 */

package nicTest

import (
	"fmt"
	"math/rand/v2"
	"searchlight/pkg/nic"
	"testing"
)

func TestGetAllNic(t *testing.T) {
	ipv4s, ipv6s := nic.GetAllNiCs()
	fmt.Println("==========ipv4s==========")
	for k, v := range ipv4s {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println("==========ipv6s==========")
	for k, v := range ipv6s {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func TestGetIpType(t *testing.T) {
	ipType := nic.GetIpType("e8:61:1f:1b:07:4a")
	fmt.Println(ipType)
}

func TestRound(t *testing.T) {
	for i := 0; i < 100; i++ {
		//round := math.Round(100)
		round := rand.IntN(2)
		fmt.Printf("i:%d,round:%d\n", i, round)
	}
}
