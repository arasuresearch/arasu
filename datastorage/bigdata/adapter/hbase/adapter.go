package hbase

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/arasuresearch/arasu/datastorage/bigdata/adapter/abstract"
	. "github.com/arasuresearch/arasu/datastorage/bigdata/adapter/hbase/thrift/Hbase"
	"log"
	"net"
)

type HbaseAdapter struct {
	abstract.AbstractAdapter
	Conn *HbaseClient
}

//creating new hbase adapter
//this adapter assumes to use hbase thrift2 binary protocal to communicate with hbase database

func New(abstractAdapter abstract.AbstractAdapter) (hbaseAdapter *HbaseAdapter) {
	hbaseAdapter = &HbaseAdapter{AbstractAdapter: abstractAdapter}
	tcpAddr, err := net.ResolveTCPAddr("tcp", hbaseAdapter.Address)

	if err != nil {
		log.Fatal(err)
	}

	transport, err := thrift.NewTSocket(tcpAddr.String())
	if err != nil {
		log.Fatal(err)
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	hbaseClient := NewHbaseClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Println("Error occured while connecting with bigdata datastore ", err, "please make sure bigdata is available")
		log.Fatal(err)
	}

	//defer transport.Close()
	hbaseAdapter.Conn = hbaseClient
	return
}

func NewAdapter(name string, conf string) *HbaseAdapter {
	return New(abstract.New(name, conf))
}
