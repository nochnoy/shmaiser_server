package main

// ----------------------------------------------------------------------------
//	Types
// ----------------------------------------------------------------------------

// World - мир
type World struct {
	nextObjectID int
	Objects      map[int]*Object `json:"objects"`
}

// ----------------------------------------------------------------------------
//	Properties
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
//	Methods
// ----------------------------------------------------------------------------

// CreateWorld - создаём мир
func CreateWorld() *World {
	w := new(World)
	w.nextObjectID = 1
	w.Objects = make(map[int]*Object)
	return w
}

// CreateObject - создаём новый объект и добавляем в мир
func (w *World) CreateObject() *Object {
	b := new(Object)
	b.ID = w.nextObjectID
	w.Objects[b.ID] = b
	w.nextObjectID++
	return b
}

// RemoveObject - удаляем объект из мира
func (w *World) RemoveObject(b *Object) {
	delete(w.Objects, b.ID)
}
