// Package database maps ROM hashes to game names.
//
// Lookups stream the embedded database and return as soon as the hash is
// found, so the full list is never held in memory.
package database

import "errors"

var ErrNotFound = errors.New("not found")
