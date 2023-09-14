package database

import "log"

type DefaultDatabase interface {
	Connect()
	Insert(data interface{}, collection string) error
}

type FakeRepo struct {
}

func (f *FakeRepo) Connect() {
	log.Print("Fake repo conected")
}

func (f *FakeRepo) Insert(data interface{}, collection string) error {
	log.Printf("fake repo insert %s", data)
	return nil
}
