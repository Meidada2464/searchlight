/**
 * Package cmdPing
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/18 11:07
 */

package cmdPing

import (
	"errors"
	"fmt"
	"log"
	"math"
	"searchlight/pkg/nic"
	"searchlight/pkg/util"
	"strconv"
	"strings"
	"time"
)

func CPingServer(srcIp, tarIp string, count, size, interVal, timeOut int, pattern bool) error {
	// 数据校验
	ipType := nic.GetIpType(tarIp)
	if ipType == "" {
		return errors.New("target ip not available")
	}

	sIp, err := nic.GetSameAvailableIp(srcIp, ipType)
	if err != nil {
		return err
	}

	// assamble args
	args := []string{
		"-I", sIp,
		"-c", fmt.Sprintf("%d", count),
		"-s", fmt.Sprintf("%d", size),
		"-i", fmt.Sprintf("%d", interVal),
		"-t", fmt.Sprintf("%d", timeOut),
	}
	if ipType == "ipv6" {
		args = append(args, "-6")
	}
	args = append(args, tarIp)

	if interVal == 1 {
		log.Printf("place wait %d seconds. cmdPing is working ... ...\n", count)
	} else {
		log.Printf("place wait %d ~ %d seconds. cmdPing is working ... ...\n", count, count*interVal)
	}
	log.Println("ping", args)

	out, errOut, err := util.RunCommand("ping", args, 30*time.Second+5*time.Second)
	if err != nil {
		return err
	}

	// show details data
	if errOut != nil && pattern {
		fmt.Println(string(out))
		return nil
	}

	if errOut != nil {
		packetLoss, minRTT, avgRTT, maxRTT, mdevRTT, duration, err2 := parsePingOutput(string(out))
		if err2 != nil {
			fmt.Println(err2)
		}
		log.Printf("\n================light cmdPing result================\n")
		log.Printf("source ip: %s\n", sIp)
		log.Printf("target ip: %s\n", tarIp)
		log.Printf("sent packets count: %d\n", count)
		log.Printf("max rtts: %f\t", maxRTT)
		log.Printf("min rtts: %f\t", minRTT)
		log.Printf("avg rtts: %f\t", avgRTT)
		log.Printf("stdDev rtts: %f\n", mdevRTT)
		log.Printf("loss: %f\n", packetLoss)
		log.Printf("time: %s\n", duration)
		return nil
	}
	return nil
}

func parsePingOutput(output string) (float64, float64, float64, float64, float64, string, error) {
	var (
		packetLoss       float64
		minRTT           float64
		avgRTT           float64
		maxRTT           float64
		mdevRTT          float64
		duration         string
		invaliLossOutput = true
		invaliRTTOutput  = true
	)

	// 根据输出中的关键字进行分割和查找
	for _, line := range strings.Split(output, "\n") {
		// 解析丢包
		if strings.Contains(line, "packet loss") {
			// 如果找到包含 packet loss 的行，则进行解析
			// Darwin:
			// 5 packets transmitted, 5 packets received, 0.0% packet loss
			// 2 packets transmitted, 0 packets received, 100.0% packet loss
			// Linux:
			// 2 packets transmitted, 2 received, 0% packet loss, time 1001ms
			// 2 packets transmitted, 0 received, +2 errors, 100% packet loss, time 999ms
			plIdx := strings.LastIndex(line, "% packet loss")
			fields := strings.Fields(line[:plIdx])
			if len(fields) > 5 {
				lost, pErr := strconv.ParseFloat(fields[len(fields)-1], 64)
				if pErr == nil {
					packetLoss = lost
					invaliLossOutput = false
				}
			}
		}

		if strings.Contains(line, "min/avg/max/mdev") {
			fields := strings.Split(line, " ")
			if len(fields) > 3 {
				rttValues := strings.Split(fields[3], "/")
				if len(rttValues) == 4 {
					minRtt, errMin := strconv.ParseFloat(rttValues[0], 64)
					avgRtt, errAvg := strconv.ParseFloat(rttValues[1], 64)
					maxRtt, errMax := strconv.ParseFloat(rttValues[2], 64)
					mdevRtt, errMdev := strconv.ParseFloat(rttValues[3], 64)
					if errMin == nil && errAvg == nil && errMax == nil && errMdev == nil {
						minRTT = math.Round((minRtt/1000)*10000) / 10000
						avgRTT = math.Round((avgRtt/1000)*10000) / 10000
						maxRTT = math.Round((maxRtt/1000)*10000) / 10000
						mdevRTT = math.Round((mdevRtt/1000)*10000) / 10000
						invaliRTTOutput = false
						//log.Printf("min: %.3f ms, avg: %.3f ms, max: %.3f ms, mdev: %.3f ms\n", minRTT, avgRTT, maxRTT, mdevRTT)
					} else {
						log.Println("解析 RTT 值时出错")
					}
				}
			}
		}

		if strings.Contains(line, "time") {
			pos := strings.Index(line, "time")
			duration = line[pos+4:]
		}
	}

	if invaliRTTOutput || invaliLossOutput {
		return 100, 0, 0, 0, 0, duration, fmt.Errorf("invalid output: %s", output)
	}

	return packetLoss, minRTT, avgRTT, maxRTT, mdevRTT, duration, nil
}
