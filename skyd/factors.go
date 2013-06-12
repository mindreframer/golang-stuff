package skyd

import (
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"strconv"
	"sync"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Factors object manages the factorization and defactorization of values.
type Factors struct {
	db    *levigo.DB
	ro    *levigo.ReadOptions
	wo    *levigo.WriteOptions
	path  string
	mutex sync.Mutex
}

//------------------------------------------------------------------------------
//
// Errors
//
//------------------------------------------------------------------------------

//--------------------------------------
// Factor Not Found
//--------------------------------------

func NewFactorNotFound(text string) error {
	return &FactorNotFound{text}
}

type FactorNotFound struct {
	s string
}

func (e *FactorNotFound) Error() string {
	return e.s
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewFactors returns a new Factors object.
func NewFactors(path string) *Factors {
	return &Factors{path: path}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// The path to the database on disk.
func (f *Factors) Path() string {
	return f.path
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// State
//--------------------------------------

// Opens the factors databse.
func (f *Factors) Open() error {
	if f.IsOpen() {
		return errors.New("skyd.Factors: Factors database is already open.")
	}

	// Open database.
	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(f.path, opts)
	if err != nil {
		f.Close()
		return fmt.Errorf("skyd.Factors: Unable to open database: %v", err)
	}
	f.db = db

	// Setup read and write options.
	f.ro = levigo.NewReadOptions()
	f.wo = levigo.NewWriteOptions()

	return nil
}

// Closes the factors database.
func (f *Factors) Close() {
	if f.db != nil {
		f.db.Close()
	}
	if f.ro != nil {
		f.ro.Close()
	}
	if f.wo != nil {
		f.wo.Close()
	}
}

// Returns whether the factors database is open.
func (f *Factors) IsOpen() bool {
	return f.db != nil
}

//--------------------------------------
// Keys
//--------------------------------------

// The key for a given namespace/id/value.
func (f *Factors) key(namespace string, id string, value string) string {
	return fmt.Sprintf("%s>%s:%s", namespace, id, value)
}

// The reverse key for a given namespace/id/value.
func (f *Factors) revkey(namespace string, id string, value uint64) string {
	return fmt.Sprintf("%s>%s:%d", namespace, id, value)
}

// The sequence key for a given namespace/id.
func (f *Factors) seqkey(namespace string, id string) string {
	return fmt.Sprintf("%s>%s!", namespace, id)
}

//--------------------------------------
// Factorization
//--------------------------------------

// Converts the defactorized value for a given id in a given namespace to its internal representation.
func (f *Factors) Factorize(namespace string, id string, value string, createIfMissing bool) (uint64, error) {
	// Blank is always zero.
	if value == "" {
		return 0, nil
	}

	// Otherwise find it in the LevelDB database.
	data, err := f.db.Get(f.ro, []byte(f.key(namespace, id, value)))
	if err != nil {
		return 0, err
	}
	// If key does exist then parse and return it.
	if data != nil {
		return strconv.ParseUint(string(data), 10, 64)
	}

	// Create a new factor if requested.
	if createIfMissing {
		return f.add(namespace, id, value)
	}

	err = NewFactorNotFound(fmt.Sprintf("skyd.Factors: Factor not found: %v", f.key(namespace, id, value)))
	return 0, err
}

// Adds a new factor to the database if it doesn't exist.
func (f *Factors) add(namespace string, id string, value string) (uint64, error) {
	// Lock while adding a new value.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Retry factorize within the context of the lock.
	sequence, err := f.Factorize(namespace, id, value, false)
	if err == nil {
		return sequence, nil
	} else if _, ok := err.(*FactorNotFound); !ok {
		return 0, err
	}

	// Retrieve next id in sequence.
	sequence, err = f.inc(namespace, id)
	if err != nil {
		return 0, err
	}

	// Save lookup and reverse lookup.
	err = f.db.Put(f.wo, []byte(f.key(namespace, id, value)), []byte(strconv.FormatUint(sequence, 10)))
	if err != nil {
		return 0, err
	}
	err = f.db.Put(f.wo, []byte(f.revkey(namespace, id, sequence)), []byte(value))
	if err != nil {
		return 0, err
	}

	return sequence, nil
}

// Converts the factorized value for a given id in a given namespace to its internal representation.
func (f *Factors) Defactorize(namespace string, id string, value uint64) (string, error) {
	// Blank is always zero.
	if value == 0 {
		return "", nil
	}

	// Find it in LevelDB.
	data, err := f.db.Get(f.ro, []byte(f.revkey(namespace, id, value)))
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", fmt.Errorf("skyd.Factors: Value does not exist: %v", f.revkey(namespace, id, value))
	}
	return string(data), nil
}

// Retrieves the next available sequence number within a namespace for an id.
func (f *Factors) inc(namespace string, id string) (uint64, error) {
	data, err := f.db.Get(f.ro, []byte(f.seqkey(namespace, id)))
	if err != nil {
		return 0, err
	}

	// Initialize key if it doesn't exist. Otherwise increment it.
	if data == nil {
		err := f.db.Put(f.wo, []byte(f.seqkey(namespace, id)), []byte("1"))
		if err != nil {
			return 0, err
		}
		return 1, nil
	}

	// Parse existing sequence.
	sequence, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("skyd.Factors: Unable to parse sequence: %v", data)
	}

	// Increment and save the new value.
	sequence += 1
	err = f.db.Put(f.wo, []byte(f.seqkey(namespace, id)), []byte(strconv.FormatUint(sequence, 10)))
	if err != nil {
		return 0, err
	}
	return sequence, nil
}
