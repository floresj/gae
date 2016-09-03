package gae

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Entity Represents an Entity that can be stored and queried. Each Entity is reponsible for defining
// how it generates its own key. The key is the unique identifer for a datastore object.
// There are a couple of ways of creating key. You can use a custom StringID (such as email address),
// IntID (any integer) or an IncompleteKey where the datastore automatically generates a key.

// Helpful links:
//
// https://cloud.google.com/appengine/docs/go/datastore/reference
type Entity interface {
	Key(context.Context) *datastore.Key
}

// EntityPreparer allows structs to provide additional logic after an Entity
// has been retrieved from the datastore
// An Entity that can contain additional processing after it has been retrieved
// from a Save, Get or Query. This can be to set struct properties that are for display purposes.
// For instance, all structs use this interface to properly populate a friendly Id representation
//
//
// Example:
//
//    func (p *Project) Prepare(key *datastore.Key) {
//        p.Id = key.Encode()
//        p.CompanyId = p.Company.Encode()
//    }
//
type EntityPreparer interface {
	Prepare(*datastore.Key)
}

// BeforePutter allows an Entity to execute logic before a Put is executed
type BeforePutter interface {
	BeforePut() error
}
