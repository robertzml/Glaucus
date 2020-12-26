package influx

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/robertzml/Glaucus/base"
	"github.com/robertzml/Glaucus/glog"
	"time"
)

const (
	packageName = "influx"
)

// InfluxDb 数据存储接口
type Repository struct {

	// 累积值存储队列
	cumulativeChan chan *influxPoint

	// 基础数据存储队列
	basicChan chan *influxPoint

	// 报警数据存储队列
	alarmChan chan *influxPoint

	// 关键数据存储队列
	keyChan chan *influxPoint
}

// Influx 数据点
type influxPoint struct {
	tags   map[string]string
	fields map[string]interface{}
}

// 初始化Influxdb 相关channel
func InitFlux() *Repository {
	repo := new(Repository)
	repo.cumulativeChan = make(chan *influxPoint, 10)
	repo.basicChan = make(chan *influxPoint, 5)
	repo.alarmChan = make(chan *influxPoint, 5)
	repo.keyChan = make(chan *influxPoint, 5)

	return repo
}

/**
保存热水器累积数据到channel
*/
func (repo *Repository) SaveCumulate(tags map[string]string, fields map[string]interface{}) {
	point := influxPoint{tags, fields}
	repo.cumulativeChan <- &point
}

/**
保存基础数据到channel
*/
func (repo *Repository) SaveBasic(tags map[string]string, fields map[string]interface{}) {
	point := influxPoint{tags, fields}
	repo.basicChan <- &point
}

/**
保存报警数据到channel
*/
func (repo *Repository) SaveAlarm(tags map[string]string, fields map[string]interface{}) {
	point := influxPoint{tags, fields}
	repo.alarmChan <- &point
}

/**
 保存关键数据到channel
 */
func (repo *Repository) SaveKey(tags map[string]string, fields map[string]interface{}) {
	point := influxPoint{tags, fields}
	repo.keyChan <- &point
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
		case packet := <-repo.cumulativeChan:
			p := influxdb2.NewPoint("cumulative", packet.tags, packet.fields, time.Now())
			writeApi.WritePoint(p)
		case packet := <-repo.basicChan:
			p := influxdb2.NewPoint("basic", packet.tags, packet.fields, time.Now())
			writeApi.WritePoint(p)
		case packet := <-repo.alarmChan:
			p := influxdb2.NewPoint("alarm", packet.tags, packet.fields, time.Now())
			writeApi.WritePoint(p)
		case packet := <-repo.keyChan:
			p := influxdb2.NewPoint("key_status", packet.tags, packet.fields, time.Now())
			writeApi.WritePoint(p)
		}
	}
}
