/**
 * Package util
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/18 14:26
 */

package util

import "math/rand"

func Rand(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.Intn(max)
}
