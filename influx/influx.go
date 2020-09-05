package influx

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/robertzml/Glaucus/base"
	"time"
)

// influxdb channel
var influxCh  chan *Packet

// 初始化
func InitFlux() {
	influxCh = make(chan *Packet, 10)
}

// 写数据到channel 中
func Write(sn string, cumulateHeatTime int, cumulateHotWater int, cumulateWorkTime int, cumulateUsedPower int, cumulateSavePower int) {
	packet := Packet{SerialNumber: sn, CumulateHeatTime: cumulateHeatTime, CumulateHotWater: cumulateHotWater,
		CumulateWorkTime: cumulateWorkTime, CumulateUsedPower: cumulateUsedPower,CumulateSavePower: cumulateSavePower}
	influxCh <- &packet
}

// 从channel 中获取待写入数据并写入到数据库
func Read() {
	client := influxdb2.NewClient(base.DefaultConfig.InfluxAddress, base.DefaultConfig.InfluxToken)

	writeApi := client.WriteAPI(base.DefaultConfig.InfluxOrg, "Molan")

	defer func() {
		writeApi.Flush()
		client.Close()
	}()

	for {
		packet := <- influxCh

		p := influxdb2.NewPointWithMeasurement("cumulative").
			AddTag("sn", packet.SerialNumber).
			AddField("cumulateHeatTime", packet.CumulateHeatTime).
			AddField("cumulateHotWater", packet.CumulateHotWater).
			AddField("cumulateWorkTime", packet.CumulateWorkTime).
			AddField("cumulateUsedPower", packet.CumulateUsedPower).
			AddField("cumulateSavePower", packet.CumulateSavePower).
			SetTime(time.Now())

		writeApi.WritePoint(p)
	}
}
