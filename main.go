package main

import (
	"bufio"
	"fmt"
	"github.com/FrancisMcN/lib-mmdb2/field"
	"github.com/FrancisMcN/lib-mmdb2/mmdb"
	"github.com/FrancisMcN/lib-mmdb2/trie"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {
	db := mmdb.NewMMDB()
	t := trie.NewTrie()
	db.PrefixTree = t
	_, c, _ := net.ParseCIDR("1.1.1.0/24")
	c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	t.Insert(c, field.String("hello world"))

	_, c, _ = net.ParseCIDR("1.1.0.0/24")
	c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	t.Insert(c, field.String("hello world"))

	_, c, _ = net.ParseCIDR("1.1.0.0/32")
	c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	t.Insert(c, field.String("hello world /32"))

	_, c, _ = net.ParseCIDR("1.1.0.0/32")
	c.IP = c.IP.To16()
	//c.IP[len(c.IP)-1-4] = 0
	//c.IP[len(c.IP)-1-5] = 0
	t.Insert(c, field.String("net:1.178.112.0/20, asn:AS12975"))

	//_, c, _ := net.ParseCIDR("8000::/1")
	//c.IP = c.IP.To16()
	//t.Insert(c, field.String("hello"))
	//
	//_, c, _ = net.ParseCIDR("4000::/2")
	//c.IP = c.IP.To16()
	//t.Insert(c, field.String("hello1"))

	//t.Print()
	fmt.Println("----")
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

		if res := db.Query(net.ParseIP(text[:len(text)-1])); res != nil {
			fmt.Println("--- found ip --- \n", res)
		} else {
			fmt.Println(fmt.Sprintf("ip '%s' not found", text[:len(text)-1]))
		}

	}
}
