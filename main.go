/**
 * Package main
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/11 23:42
 */

package main

import (
	"log"
	"os"
	"searchlight/cmd"
)

func main() {
	log.SetFlags(0)
	err := cmd.Execute(os.Args[1:])
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
