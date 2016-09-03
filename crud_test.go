package gae

import (
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

type dog struct {
	ID    int64 `datastore:"-"`
	Name  string
	Color string
}

func (d dog) Key(ctx context.Context) *datastore.Key {
	if d.ID == 0 {
		return datastore.NewIncompleteKey(ctx, "Dog", nil)
	}
	return datastore.NewKey(ctx, "Dog", "", d.ID, nil)
}

func (d *dog) Prepare(key *datastore.Key) {
	fmt.Println("Dog Prepare", key.IntID())
	d.ID = key.IntID()
}

func TestPut(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	d := dog{Name: "Winston", Color: "Red/White"}
	key, err := Put(ctx, &d)
	if err != nil {
		t.Error(err)
	}

	if key.Incomplete() {
		t.Error("Key is incomplete")
	}

	savedDog := dog{}
	if err := datastore.Get(ctx, key, &savedDog); err != nil {
		t.Error(err)
	}

	if savedDog.Name != "Winston" {
		t.Errorf("Expected Name to be Wiston\n")
	}

	savedDog = dog{ID: key.IntID()}
	if err := Get(ctx, &savedDog); err != nil {
		t.Error("Could not retrieve saved dog", err)
	}
}

func TestGet(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	d := createDog("Finley")
	dogKey := datastore.NewIncompleteKey(ctx, "Dog", nil)
	key, err := datastore.Put(ctx, dogKey, d)
	if err != nil {
		log.Fatal(err)
	}
	d = &dog{ID: key.IntID()}

	if err := Get(ctx, d); err != nil {
		t.Fatal(err)
	}

	if d.Name != "Finley" {
		t.Fatal("Expected dog name to be Finley")
	}
}

func createDog(name string) *dog {
	return &dog{Name: name}
}
