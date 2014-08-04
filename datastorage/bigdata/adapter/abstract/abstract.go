// --
// Copyright (c) 2014 Thaniyarasu Kannusamy <thaniyarasu@gmail.com>.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// ++
//

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
