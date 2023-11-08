package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/acifani/vita/lib/game"
)

var nullID = strings.Repeat("0", 32)

type UniversalDataRecord struct {
	ID       string
	TopID    string
	BottomID string
	LeftID   string
	RightID  string
	Cells    []byte
}

func NewUniversalDataRecord(id string) *UniversalDataRecord {
	return &UniversalDataRecord{
		ID:       id,
		TopID:    nullID,
		BottomID: nullID,
		LeftID:   nullID,
		RightID:  nullID,
		Cells:    make([]byte, height*width),
	}
}

func UniversalDataRecordFromStore(value []byte) *UniversalDataRecord {
	data := NewUniversalDataRecord(nullID)
	data.Write(value)

	return data
}

func UniverseFromDataRecord(data *UniversalDataRecord) *game.Universe {
	universe := game.NewUniverse(height, width)
	universe.Write(data.Cells)

	return universe
}

func (u *UniversalDataRecord) Read(p []byte) (n int, err error) {
	copy(p[:32], []byte(u.ID))
	copy(p[32:64], []byte(u.TopID))
	copy(p[64:96], []byte(u.BottomID))
	copy(p[96:128], []byte(u.LeftID))
	copy(p[128:160], []byte(u.RightID))
	copy(p[160:], u.Cells)

	return len(p), nil
}

func (u *UniversalDataRecord) Write(p []byte) (n int, err error) {
	u.ID = string(p[:32])
	u.TopID = string(p[32:64])
	u.BottomID = string(p[64:96])
	u.LeftID = string(p[96:128])
	u.RightID = string(p[128:160])
	copy(u.Cells, p[160:])

	return len(p), nil
}

func generateKey() string {
	var result [32]byte
	rand.Read(result[:])
	return hex.EncodeToString(result[:])
}
