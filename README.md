#gae crud

Simplifies basic crud operations when using the Go App Engine Datastore. 

# How to use

1.Define your structs
```go
    type Dog struct {
        ID int64
        Name string
        Age int64
        Modified time.Time
    }
```

2.Implement `Entity` interface. It's up to you how and what values to use for your `datastore.Key`
```go

func (d Dog) Key(ctx context.Context) *datastore.Key {
	if d.ID == 0 {
		return datastore.NewIncompleteKey(ctx, "Dog", nil)
	}
	return datastore.NewKey(ctx, "Dog", "", d.ID, nil)
}
```

3.Implement `Preparer` interface (optional). This is executed after your entity is retrieved. Put any
post processing logic like setting your struct to store the `int` or `string` id from the key
```go
func (d *Dog) Prepare(key *datastore.Key) {
	d.ID = key.IntID()
}
```

4.Implement `BeforePutter` interface (optional). This is executed before your entity is saved
```go
func (d *Dog) BeforePut() error {
    d.Modified = time.Now()
    return nil
}
```

5.Start using it!
```go
// Put
dog := Dog{Name: "Winston", Age: 2}
key, err := gae.Put(ctx, &dog)

// Get
dog := Dog{ID: 1234}
err := gae.Get(ctx, &dog)
```
