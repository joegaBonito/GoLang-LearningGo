package mydict

import "errors"

//Dictionary type
type Dictionary map[string]string

var errorNotFound = errors.New("Not found")
var errorWordExists = errors.New("That word already exists.")
var errorCantUpdate = errors.New("Cant update non-existing word")

//Search for a word.
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errorNotFound
}

//Add a new word with definition
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)

	switch err {
	case errorNotFound:
		d[word] = def
	case nil:
		return errorWordExists
	}
	return nil
}

//Update a word with new definition
func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = def
	case errorNotFound:
		return errorCantUpdate
	}
	return nil
}

//Delete a word
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
