/**
 * Package goPing
 * @Author fengfeng.mei <Biophiliam@protonmail.com>
 * @Date 2024/12/17 13:44
 */

package goPing

import (
	"errors"
	probing "github.com/nvksie/pro-bing"
	"log"
	"searchlight/pkg/nic"
	"time"
)

func GPService(srcIp, tarIp string, count, size, interVal, timeOut int) error {
	p := probing.New(tarIp)

	p.Count = count
	p.Size = size
	p.Interval = time.Duration(interVal) * time.Second
	p.Timeout = time.Duration(timeOut) * time.Second
	// 当能获取到对应的网卡ip时，通过IP指定出口网卡
	sIp, err := nic.GetSameAvailableIp(srcIp, nic.GetIpType(tarIp))
	if err != nil {
		return err
	}
	if sIp != "" {
		p.Source = sIp
	}

	if interVal == 1 {
		log.Printf("place wait %d seconds. goPing is working ... ...\n", count)
	} else {
		log.Printf("place wait %d ~ %d seconds. goPing is working ... ...\n", count, count*interVal)
	}
	err = p.Run()
	if err != nil {
		return err
	}

	res := p.Statistics()

	if res == nil {
		return errors.New("goPing result is nil")
	}

	log.Printf("\n================light goPing result================\n")
	log.Printf("source ip: %s\n", p.Source)
	log.Printf("target ip: %s\n", res.IPAddr.String())
	log.Printf("sent packets count: %d\n", res.PacketsSent)
	log.Printf("received packets count: %d\n", res.PacketsRecv)
	log.Printf("received duplicates packets count: %d\n", res.PacketsRecvDuplicates)
	log.Printf("rtts: %f\n", func(rtts []time.Duration) []float64 {
		var rs []float64
		for _, v := range rtts {
			rs = append(rs, v.Seconds())
		}
		return rs
	}(res.Rtts))
	log.Printf("ttls: %d\n", res.TTLs)
	log.Printf("max rtts: %f\t", res.MaxRtt.Seconds())
	log.Printf("min rtts: %f\t", res.MinRtt.Seconds())
	log.Printf("avg rtts: %f\t", res.AvgRtt.Seconds())
	log.Printf("stdDev rtts: %f\n", res.StdDevRtt.Seconds())
	log.Printf("loss: %f\n", res.PacketLoss)
	return nil
}
