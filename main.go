package main

import (
	"bufio"
	"fmt"
	"github.com/FrancisMcN/lib-mmdb/field"
	"github.com/FrancisMcN/lib-mmdb/mmdb"
	"github.com/FrancisMcN/lib-mmdb/trie"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	db := mmdb.NewMMDB()
	t := trie.NewTrie()
	db.PrefixTree = t

	_, c, _ := net.ParseCIDR("1.0.0.0/24")
	c.IP = c.IP.To16()
	c.IP[len(c.IP)-1-4] = 0
	c.IP[len(c.IP)-1-5] = 0
	//m := make(map[field.Field]field.Field)
	m := field.NewMap()
	m.Put(field.String("country"), field.String("AU"))
	t.Insert(c, m)

	_, c, _ = net.ParseCIDR("1.0.1.0/24")
	c.IP = c.IP.To16()
	c.IP[len(c.IP)-1-4] = 0
	c.IP[len(c.IP)-1-5] = 0
	m = field.NewMap()
	m.Put(field.String("country"), field.String("CN"))
	t.Insert(c, m)

	_, c, _ = net.ParseCIDR("1.0.2.0/24")
	c.IP = c.IP.To16()
	c.IP[len(c.IP)-1-4] = 0
	c.IP[len(c.IP)-1-5] = 0
	m = field.NewMap()
	m.Put(field.String("country"), field.String("JP"))
	t.Insert(c, m)

	_, c, _ = net.ParseCIDR("1.0.3.0/24")
	c.IP = c.IP.To16()
	c.IP[len(c.IP)-1-4] = 0
	c.IP[len(c.IP)-1-5] = 0
	//m = make(map[field.Field]field.Field)
	m = field.NewMap()
	m.Put(field.String("country2"), field.String("US"))
	t.Insert(c, m)
	//t.Insert(c, field.String("hello world"))

	_, c, _ = net.ParseCIDR("1.0.4.0/24")
	c.IP = c.IP.To16()
	c.IP[len(c.IP)-1-4] = 0
	c.IP[len(c.IP)-1-5] = 0
	m = field.NewMap()
	m.Put(field.String("country2"), field.String("US"))
	t.Insert(c, m)

	//_, c, _ := net.ParseCIDR("1.1.1.0/24")
	//c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	////m := make(map[field.Field]field.Field)
	////m[field.String("test")] = field.String("hello world")
	////t.Insert(c, field.Map(m))
	//t.Insert(c, field.String("hello world"))
	//
	//_, c, _ = net.ParseCIDR("1.1.0.0/24")
	//c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	////m = make(map[field.Field]field.Field)
	////m[field.String("test")] = field.String("hello world")
	////t.Insert(c, field.Map(m))
	//t.Insert(c, field.String("hello world"))
	//
	//_, c, _ = net.ParseCIDR("2.1.0.0/24")
	//c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	////m = make(map[field.Field]field.Field)
	////m[field.String("test")] = field.String("hello world")
	////m[field.String("abc")] = field.String("xyz")
	////m[field.String("def")] = field.String("francis")
	////m[field.String("test")] = field.String("hello world")
	////m[field.String("123")] = field.String("456")
	////t.Insert(c, field.Map(m))
	//t.Insert(c, field.String("hello world"))
	//
	//_, c, _ = net.ParseCIDR("1.1.0.0/32")
	//c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	//m := make(map[field.Field]field.Field)
	//m[field.String("test")] = field.String("abc /32")
	//m[field.String("francis")] = field.String("mcnamee")
	//m[field.String("hello world")] = field.String("123")
	//t.Insert(c, field.Map(m))
	////t.Insert(c, field.String("hello /32"))
	//
	//_, c, _ = net.ParseCIDR("5.1.0.0/32")
	//c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	////m = make(map[field.Field]field.Field)
	////m[field.String("test")] = field.String("abc /32")
	////t.Insert(c, field.Map(m))
	//t.Insert(c, field.String("abc /32"))
	//
	//_, c, _ = net.ParseCIDR("10.0.0.0/24")
	//c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	//m = make(map[field.Field]field.Field)
	//m[field.String("test")] = field.String("abc /32")
	//m[field.String("francis")] = field.String("123")
	//t.Insert(c, field.Map(m))
	////t.Insert(c, field.String("francis"))

	//_, c, _ = net.ParseCIDR("1.1.0.0/32")
	//c.IP = c.IP.To16()
	////c.IP[len(c.IP)-1-4] = 0
	////c.IP[len(c.IP)-1-5] = 0
	//t.Insert(c, field.String("net:1.178.112.0/20, asn:AS12975"))

	//_, c, _ := net.ParseCIDR("8000::/1")
	//c.IP = c.IP.To16()
	//t.Insert(c, field.String("hello"))
	//
	//_, c, _ = net.ParseCIDR("4000::/2")
	//c.IP = c.IP.To16()
	//t.Insert(c, field.String("hello1"))

	//t.Print()
	//fmt.Println("----")
	t.Finalise()
	//t.Print()
	err := ioutil.WriteFile("test.mmdb", db.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	db.Load(db.Bytes())

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}

		ip := net.ParseIP(text[:len(text)-1])
		ip[len(ip)-1-4] = 0
		ip[len(ip)-1-5] = 0
		if res := db.Query(ip); res != nil {
			fmt.Println("--- found ip --- \n", res)
		} else {
			fmt.Println(fmt.Sprintf("ip '%s' not found", text[:len(text)-1]))
		}

	}
}
