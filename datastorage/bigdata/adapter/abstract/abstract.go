// --
// The BSD License (BSD)

// Copyright (c) 2015 Thaniyarasu Kannusamy <thaniyarasu@gmail.com> & Arasu Research Lab Pvt Ltd. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above copyright notice, this list of
//    conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//    * Neither Thaniyarasu Kannusamy <thaniyarasu@gmail.com>. nor ArasuResearch Inc may be used to endorse or promote products derived from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND AUTHOR
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// ++

package abstract

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

type AbstractAdapter struct {
	Name    string
	Conf    string
	Address string
	DbName  string
	Query   map[string][]string
	Config  map[string]interface{}

	SchemaTableName         string
	SchemaTableVersion      string
	SchemaTableColumnFamily string

	Urandom       *os.File
	UrandomStream []byte
}

func (a *AbstractAdapter) GenUUID() string {
	a.Urandom.Read(a.UrandomStream)
	return fmt.Sprintf("%x", a.UrandomStream)
}

func (a *AbstractAdapter) GetDbName() string {
	return a.DbName
}
func New(name string, conf string) AbstractAdapter {
	var abstractAdapter = AbstractAdapter{
		SchemaTableName:         "schema_migrations",
		SchemaTableVersion:      "versions",
		SchemaTableColumnFamily: "cf",
	}
	urandom, err := os.Open("/dev/urandom")
	abstractAdapter.Urandom = urandom
	if err != nil {
		log.Println(err)
	}

	abstractAdapter.UrandomStream = make([]byte, 8)

	address, dbName, query := parseConfiguration(conf)
	abstractAdapter.Name = name
	abstractAdapter.Address = address

	abstractAdapter.DbName = dbName
	abstractAdapter.Query = query
	abstractAdapter.Conf = conf
	return abstractAdapter
}

var ErrNoDb = errors.New("arasu: no database specified (like 'username:password@protocol(host:post)/database_name?arg=value&...')")

// TODO : parse the configuration in deeply like database/sql adapter.
func parseConfiguration(conf string) (address string, db_name string, res map[string][]string) {

	u, err := url.Parse(conf)
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(u.Opaque, "/") {
		opaques := strings.Split(u.Opaque, "/")
		if len(opaques) > 1 && len(opaques[1]) > 0 {
			address = strings.Split(conf, "/")[0] + "/" //"/?" + u.RawQuery //opaques[0]
			address = strings.Split(strings.Split(address, "(")[1], ")")[0]
			db_name = opaques[1]
		} else {
			log.Fatal(ErrNoDb)
		}

	} else {
		log.Fatal(ErrNoDb)
	}

	q := u.Query()
	res = map[string][]string{}
	for k, v := range q {
		res[k] = []string{}
		for _, e := range v {
			if strings.Contains(e, ",") {
				for _, e0 := range strings.Split(e, ",") {
					res[k] = append(res[k], e0)
				}
			} else {
				res[k] = append(res[k], e)
			}

		}
	}
	return
}
