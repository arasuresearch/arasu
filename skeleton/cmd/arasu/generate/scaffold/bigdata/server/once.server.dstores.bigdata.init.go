package bigdata

import (
	"boot"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase"
)

var BigData *hbase.HbaseAdapter

func init() {
	name, conf := boot.App.DsNameConf("bigdata")
	BigData = hbase.NewAdapter(name, conf)
}
