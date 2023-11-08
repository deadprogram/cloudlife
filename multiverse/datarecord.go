package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/acifani/vita/lib/game"
)

var nullID = strings.Repeat("0", 32)

type DataRecord struct {
	ID       string
	TopID    string
	BottomID string
	LeftID   string
	RightID  string
	Cells    []byte
}

func NewDataRecord(id string) *DataRecord {
	return &DataRecord{
		ID:       id,
		TopID:    nullID,
		BottomID: nullID,
		LeftID:   nullID,
		RightID:  nullID,
		Cells:    make([]byte, height*width),
	}
}

func DataRecordFromStore(value []byte) *DataRecord {
	data := NewDataRecord(nullID)
	data.Write(value)

	return data
}

func UniverseFromDataRecord(data *DataRecord) *game.Universe {
	universe := game.NewUniverse(height, width)
	universe.Write(data.Cells)

	return universe
}

func (dr *DataRecord) Read(p []byte) (n int, err error) {
	copy(p[:32], []byte(dr.ID))
	copy(p[32:64], []byte(dr.TopID))
	copy(p[64:96], []byte(dr.BottomID))
	copy(p[96:128], []byte(dr.LeftID))
	copy(p[128:160], []byte(dr.RightID))
	copy(p[160:], dr.Cells)

	return len(p), nil
}

func (dr *DataRecord) Write(p []byte) (n int, err error) {
	dr.ID = string(p[:32])
	dr.TopID = string(p[32:64])
	dr.BottomID = string(p[64:96])
	dr.LeftID = string(p[96:128])
	dr.RightID = string(p[128:160])
	copy(dr.Cells, p[160:])

	return len(p), nil
}

func generateKey() string {
	var result [32]byte
	rand.Read(result[:])
	return hex.EncodeToString(result[:])
}
