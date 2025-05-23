package utils

import (
	"log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)


func StringToPgTypeUUID(uuidStr string) (pgtype.UUID, error) {
	uuidVer, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Fatalf("Failed to parse UUID: %v", err)
		return pgtype.UUID{}, err
	}

	return UUIDToPgType(uuidVer), nil
}

// UUIDToPgType converts a uuid.UUID to a PostgreSQL-compatible pgtype.UUID
func UUIDToPgType(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes:  [16]byte(uuid),
		Valid: true,
	}
}
