package influx

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/equipment"
	"github.com/robertzml/Glaucus/glog"
	"time"
)

const (
	packageName = "influx"
)

// InfluxDb 数据存储接口
type Repository struct {

	// 累积值存储队列
	cumulateChan chan *equipment.WaterHeaterCumulate
}

// 初始化Influxdb 相关channel
func InitFlux() *Repository {
	repo := new(Repository)
	repo.cumulateChan = make(chan *equipment.WaterHeaterCumulate, 10)

	return repo
}

// 用于继承equipment.Context
func (repo *Repository) Connect() {

}

/*
 保存热水器累积数据到channel
 */
func (repo *Repository) SaveCumulate(data *equipment.WaterHeaterCumulate) {
	repo.cumulateChan <- data
}

/*
存储数据到数据库
 */
func (repo *Repository) Process() {
	client := influxdb2.NewClient(base.DefaultConfig.InfluxAddress, base.DefaultConfig.InfluxToken)

	writeApi := client.WriteAPI(base.DefaultConfig.InfluxOrg, base.DefaultConfig.InfluxBucket)

	// Get errors channel
	errorsCh := writeApi.Errors()
	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			glog.Write(2, packageName, "process", fmt.Sprintf("write error: %s", err.Error()))
		}
	}()

	defer func() {
		writeApi.Flush()
		client.Close()
	}()

	for {
		select {
		case packet := <- repo.cumulateChan:
			p := influxdb2.NewPointWithMeasurement("cumulative").
				AddTag("serialNumber", packet.SerialNumber).
				AddTag("mainboardNumber", packet.MainboardNumber).
				AddField("cumulateHeatTime", packet.CumulateHeatTime).
				AddField("cumulateHotWater", packet.CumulateHotWater).
				AddField("cumulateWorkTime", packet.CumulateWorkTime).
				AddField("cumulateUsedPower", packet.CumulateUsedPower).
				AddField("cumulateSavePower", packet.CumulateSavePower).
				AddField("coldInTemp", packet.ColdInTemp).
				AddField("setTemp", packet.SetTemp).
				AddField("energySave", packet.EnergySave).
				SetTime(time.Now())

			writeApi.WritePoint(p)
		}
	}
}