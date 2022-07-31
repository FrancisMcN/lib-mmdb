package main

import (
	"fmt"
	"github.com/FrancisMcN/lib-mmdb2/field"
	"github.com/FrancisMcN/lib-mmdb2/mmdb"
	"github.com/FrancisMcN/lib-mmdb2/trie"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	db := mmdb.NewMMDB()
	t := trie.NewTrie()
	db.PrefixTree = t
	//_, c, _ := net.ParseCIDR("1.1.1.0/24")
	//c.IP = c.IP.To16()
	//t.Insert(c, field.String("hello world"))
	//
	//_, c, _ = net.ParseCIDR("1.1.0.0/24")
	//c.IP = c.IP.To16()
	//t.Insert(c, field.String("hello world"))

	_, c, _ := net.ParseCIDR("8000::/1")
	c.IP = c.IP.To16()
	t.Insert(c, field.String("hello world"))

	_, c, _ = net.ParseCIDR("4000::/2")
	c.IP = c.IP.To16()
	t.Insert(c, field.String("hello world 123"))

	//t.Print()
	fmt.Println("----")
	t.Finalise()
	//t.Print()
	err := ioutil.WriteFile("test.mmdb", db.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(fmt.Sprintf("%x", t.Bytes()))
	//t.SetTotalId(big.NewInt(200))
	//t.Print()

	//fmt.Println("Hello world")
}
