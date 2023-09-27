package generateid

import "github.com/google/uuid"

func GenerateUUID() string {
	idTr := uuid.New().String()
	return idTr
}
