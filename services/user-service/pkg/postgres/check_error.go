package postgres

import "github.com/lib/pq"

func IsDuplicateKeyError(err error) bool {
    if pgErr, ok := err.(*pq.Error); ok {
        return pgErr.Code == "23505" // unique_violation
    }
    
    return false
}