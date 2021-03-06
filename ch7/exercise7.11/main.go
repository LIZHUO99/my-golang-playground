package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	rawDB rawDatabase
	lock  sync.Mutex
}

type rawDatabase map[string]dollars

func (d *database) create(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.rawDB[item]
	if ok {
		fmt.Fprintf(w, "This item has already existed: %q\n", item)
		return
	}
	priceF, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Fprintf(w, "Price is invalid\n")
		return
	}
	d.rawDB[item] = dollars(priceF)
	fmt.Fprintf(w, "Successfully created item: %s\n", item)
}

func (d *database) read(w http.ResponseWriter, r *http.Request) {
	d.lock.Lock()
	defer d.lock.Unlock()
	for k, v := range d.rawDB {
		fmt.Fprintf(w, "%s %s\n", k, v)
	}
}

func (d *database) update(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.rawDB[item]
	if !ok {
		fmt.Fprintf(w, "This item is not exist.\n")
		return
	}
	priceF, err := strconv.ParseFloat(price, 32)
	if err != nil {
		fmt.Fprintf(w, "Price is invalid\n")
		return
	}
	d.rawDB[item] = dollars(priceF)
	fmt.Fprintf(w, "Price is updated to: %s", dollars(priceF))
	return
}

func (d *database) delete(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.rawDB[item]
	if !ok {
		fmt.Fprintf(w, "This item dose not exist.\n")
		return
	}
	delete(d.rawDB, item)
	fmt.Fprintf(w, "Item \"%s\" is deleted from database.\n", item)
	return
}

func main() {
	db := new(database)
	db.rawDB = rawDatabase{"shoes": 50, "socks": 5}
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/delete", db.delete)
	http.HandleFunc("/update", db.update)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
