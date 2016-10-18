package gae

import (
	"reflect"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// Put stores an Entity into the the datastore and applies EntityPreparer and BeforePutter logic
// for all structs that implement those interfaces.
func Put(ctx context.Context, entity Entity) (key *datastore.Key, err error) {
	key = entity.Key(ctx)

	// Execute logic before Put
	err = beforePut(entity)
	if err != nil {
		return nil, err
	}

	// Execute Put
	key, err = datastore.Put(ctx, key, entity)
	if err != nil {
		return nil, err
	}

	// Execute logic to prepare entity
	prepare(key, entity)
	return key, nil
}

// Get retrieves an Entity from the datastore
func Get(ctx context.Context, entity Entity) error {
	key := entity.Key(ctx)
	if err := datastore.Get(ctx, key, entity); err != nil {
		return err
	}

	prepare(key, entity)
	return nil
}

// Delete removes an Entity from the datastore
func Delete(ctx context.Context, entity Entity) error {
	key := entity.Key(ctx)
	return datastore.Delete(ctx, key)
}

// RunQuery executes a query,  populates the provided slice and returns a Cursor
func RunQuery(ctx context.Context, q *datastore.Query, entities interface{}) (datastore.Cursor, error) {

	// This was inspired by the mongodb mbo library
	resultv := reflect.ValueOf(entities)
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	t := q.Run(ctx)
	for {
		elemp := reflect.New(elemt)
		key, err := t.Next(elemp.Interface())
		if err == datastore.Done {
			break
		}
		if err != nil {
			break
		}
		if e, ok := elemp.Elem().Addr().Interface().(Entity); ok {
			prepare(key, e)
		}
		slicev = reflect.Append(slicev, elemp.Elem())
	}

	resultv.Elem().Set(slicev)
	return t.Cursor()
}

// KeysOnly queries and only retrieves entity keys
func KeysOnly(ctx context.Context, q *datastore.Query) (keys []*datastore.Key, cursor datastore.Cursor, err error) {
	// Ensure keys only
	q = q.KeysOnly()

	t := q.Run(ctx)
	for {
		key, err := t.Next(nil)
		if err == datastore.Done {
			break
		}
		if err != nil {
			break
		}
		keys = append(keys, key)
		cursor, _ = t.Cursor()
	}

	return keys, cursor, err
}

func prepare(key *datastore.Key, entity Entity) {
	if prep, ok := entity.(EntityPreparer); ok {
		prep.Prepare(key)
	}
}

func beforePut(entity Entity) error {
	if prep, ok := entity.(BeforePutter); ok {
		return prep.BeforePut()
	}
	return nil
}
