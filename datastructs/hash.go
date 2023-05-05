package datastructs

/*
 * go map data.
 * author: CC
 * email : 151503324@qq.com
 * date  : 2017.06.17
 */
/************************************************************************/
// constants, variables, structs, interfaces.
/************************************************************************/

// Hash TODO.
type Hash struct {
	table map[interface{}]interface{}
	count int32
}

/************************************************************************/
// export functions.
/************************************************************************/

// NewHash TODO.
func NewHash(count int32) *Hash {
	m := new(Hash)
	m.table = make(map[interface{}]interface{})
	m.count = count
	return m
}

// Add TODO.
func (owner *Hash) Add(key interface{}, value interface{}) {
	owner.table[key] = value
}

// Set TODO.
func (owner *Hash) Set(key interface{}, value interface{}) {
	owner.table[key] = value
}

// Get TODO.
func (owner *Hash) Get(key interface{}) interface{} {
	ret, ok := owner.table[key]
	if ok {
		return ret
	}
	return nil
}

// Del TODO.
func (owner *Hash) Del(key interface{}) {
	delete(owner.table, key)
}

// Len TODO.
func (owner *Hash) Len() int {
	return len(owner.table)
}

// ForRange TODO.
func (owner *Hash) ForRange(f func(key interface{}, value interface{})) {
	for k, v := range owner.table {
		f(k, v)
	}
}

// ForRangeWithBreak TODO.
func (owner *Hash) ForRangeWithBreak(f func(key interface{}, value interface{}) bool) {
	for k, v := range owner.table {
		if f(k, v) {
			break
		}
	}
}

/************************************************************************/
// moudule functions.
/************************************************************************/

/************************************************************************/
// unit tests.
/************************************************************************/
