package main

import "github.com/dronnix/search-accomodation/storage"

func main() {
	s := storage.NewIPLocationStorage(nil)
	s.MigrateUp(nil, "")
}
