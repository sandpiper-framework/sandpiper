// Copyright Auto Care Association. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.md file.

package sandpiper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

// GrainType is an enum defining grain content
type GrainType int16

const ( // Never change these values, only add to list
	GrainUnknown        GrainType = 0
	GrainAcesFile       GrainType = 1
	GrainAcesItem       GrainType = 2
	GrainAssetFile      GrainType = 3
	GrainPartsproFile   GrainType = 4
	GrainPartsproItem   GrainType = 5
	GrainPiesFile       GrainType = 6
	GrainPiesItem       GrainType = 7
	GrainPiesMarketcopy GrainType = 8
	GrainPiesPricesheet GrainType = 9
)

// GrainTypeStrings are the json values for GrainType enum
var GrainTypeStrings = []string{
	"not_set",
	"aces_file",
	"aces_item",
	"asset_file",
	"partspro_file",
	"partspro_item",
	"pies_file",
	"pies_item",
	"pies_marketcopy",
	"pies_pricesheet",
}

// EncodingMethod is an enum describing how the payload is encoded
type EncodingMethod int16

const ( // Never change these values, only add to list
  EncUnknown EncodingMethod = 0
	EncRaw     EncodingMethod = 1
	EncB64     EncodingMethod = 2
	EncGzipB64 EncodingMethod = 3
)

var EncodingMethodStrings = []string {
	"not_set",
	"raw",
	"b64",
	"gzipb64",
}

// Grain represents the sandpiper syncable-object
type Grain struct {
	ID        uuid.UUID      `json:"id" pg:",pk"`
	SliceID   *uuid.UUID     `json:"slice_id,omitempty"` // must be pointer for omitempty to work here!
	Type      GrainType      `json:"grain_type" pg:"grain_type"`
	Key       string         `json:"grain_key" pg:"grain_key"`
	Encoding  EncodingMethod `json:"encoding"`
	Payload   PayloadData    `json:"payload,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	Slice     *Slice         `json:"slice"` // has-one relation
}

// compile-time check variables for model hooks (which take no memory)
var _ orm.BeforeInsertHook = (*Grain)(nil)

// BeforeInsert hooks into insert operations, setting createdAt and updatedAt to current time
func (b *Grain) BeforeInsert(ctx context.Context) (context.Context, error) {
	b.CreatedAt = time.Now()
	return ctx, nil
}

// MarshalJSON method is a custom decoding for GrainType enum
func (g GrainType) MarshalJSON() ([]byte, error) {
	i := int(g)
	if i < 0 || i >= len(GrainTypeStrings) {
		return nil, errors.New(fmt.Sprintf("unknown GrainType (%d)", i))
	}
	return json.Marshal(GrainTypeStrings[i])
}

// UnmarshalJSON method is a custom encoding for GrainType enum
func (g *GrainType) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for i, v := range GrainTypeStrings {
		if v == s {
			*g = GrainType(i)
			return nil
		}
	}
	*g = GrainType(0)
	return errors.New(fmt.Sprintf("unknown GrainType (%s)", s))
}

// MarshalJSON method is a custom decoding for EncodingMethod enum
func (g EncodingMethod) MarshalJSON() ([]byte, error) {
	i := int(g)
	if i < 0 || i >= len(EncodingMethodStrings) {
		return nil, errors.New(fmt.Sprintf("unknown EncodingMethod (%d)", i))
	}
	return json.Marshal(EncodingMethodStrings[i])
}

// UnmarshalJSON method is a custom encoding for EncodingMethod enum
func (g *EncodingMethod) UnmarshalJSON(b []byte) error {
	var s string

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	for i, v := range EncodingMethodStrings {
		if v == s {
			*g = EncodingMethod(i)
			return nil
		}
	}
	*g = EncodingMethod(0)
	return errors.New(fmt.Sprintf("unknown EncodingMethod (%s)", s))
}