package structs_test

import (
	"testing"

	"github.com/ribice/gorsk/internal/platform/structs"
	"github.com/stretchr/testify/assert"
)

type carID int64

type mergeStruct struct {
	Name              string
	Surname           string
	SomeID            int64
	Cars              []string
	CarIDs            []carID
	CustomStruct      customStruct
	CustomStructSlice []customStruct
	IgnoredField      string
}

type customStruct struct {
	isLoved bool
	isHated bool
}

type mergeCmd struct {
	Name              *string
	Surname           *string
	SomeID            *int64
	Cars              []string
	CarIDs            []carID
	CustomStruct      *customStruct
	CustomStructSlice []customStruct
	IgnoredField      *string `structs:"-"`
}

func TestMerge(t *testing.T) {
	cases := map[string]struct {
		mstruct      mergeStruct
		cmd          func() mergeCmd
		mergedStruct mergeStruct
	}{
		"basic merge": {
			mstruct: mergeStruct{
				Name: "Name",
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				return mergeCmd{
					Surname: &surname,
				}
			},
			mergedStruct: mergeStruct{
				Name:    "Name",
				Surname: "Surname",
			},
		},
		"basic slice merge": {
			mstruct: mergeStruct{
				Name: "Dz",
			},
			cmd: func() mergeCmd {
				surname := "G"
				return mergeCmd{
					Surname: &surname,
					Cars:    []string{"peugeot", "citroen"},
				}
			},
			mergedStruct: mergeStruct{
				Name:    "Dz",
				Surname: "G",
				Cars:    []string{"peugeot", "citroen"},
			},
		},
		"custom slice merge": {
			mstruct: mergeStruct{
				Name: "Dz",
			},
			cmd: func() mergeCmd {
				surname := "G"
				return mergeCmd{
					Surname: &surname,
					Cars:    []string{"peugeot", "citroen"},
					CarIDs:  []carID{1, 2},
				}
			},
			mergedStruct: mergeStruct{
				Name:    "Dz",
				Surname: "G",
				Cars:    []string{"peugeot", "citroen"},
				CarIDs:  []carID{1, 2},
			},
		},
		"merge slice": {
			mstruct: mergeStruct{
				Name:   "Name",
				SomeID: 5,
				CarIDs: []carID{1, 2, 3},
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				someID := int64(3)
				return mergeCmd{
					Surname: &surname,
					SomeID:  &someID,
					CarIDs:  []carID{6, 7, 8},
				}
			},
			mergedStruct: mergeStruct{
				Name:    "Name",
				Surname: "Surname",
				SomeID:  3,
				CarIDs:  []carID{6, 7, 8},
			},
		},
		"test ignored fields": {
			mstruct: mergeStruct{
				Name:         "Name",
				SomeID:       5,
				IgnoredField: "ignored",
				CarIDs:       []carID{1, 2, 3},
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				someID := int64(3)
				ignored := "ignored-update"
				return mergeCmd{
					Surname:      &surname,
					SomeID:       &someID,
					IgnoredField: &ignored,
					CarIDs:       []carID{6, 7, 8},
				}
			},
			mergedStruct: mergeStruct{
				Name:         "Name",
				Surname:      "Surname",
				SomeID:       3,
				IgnoredField: "ignored",
				CarIDs:       []carID{6, 7, 8},
			},
		},
		"custom struct merge": {
			mstruct: mergeStruct{
				Name:         "Name",
				SomeID:       5,
				IgnoredField: "ignored",
				CarIDs:       []carID{1, 2, 3},
				CustomStruct: customStruct{
					isLoved: true,
					isHated: false,
				},
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				someID := int64(3)
				ignored := "ignored-update"
				return mergeCmd{
					Surname:      &surname,
					SomeID:       &someID,
					IgnoredField: &ignored,
					CarIDs:       []carID{6, 7, 8},
					CustomStruct: &customStruct{
						isLoved: false,
						isHated: false,
					},
				}
			},
			mergedStruct: mergeStruct{
				Name:         "Name",
				Surname:      "Surname",
				SomeID:       3,
				IgnoredField: "ignored",
				CarIDs:       []carID{6, 7, 8},
				CustomStruct: customStruct{
					isLoved: false,
					isHated: false,
				},
			},
		},
		"custom struct slice merge": {
			mstruct: mergeStruct{
				Name:         "Name",
				SomeID:       5,
				IgnoredField: "ignored",
				CarIDs:       []carID{1, 2, 3},
				CustomStructSlice: []customStruct{
					{
						isLoved: true,
						isHated: true,
					},
					{
						isLoved: false,
						isHated: false,
					},
				},
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				someID := int64(3)
				ignored := "ignored-update"
				return mergeCmd{
					Surname:      &surname,
					SomeID:       &someID,
					IgnoredField: &ignored,
					CarIDs:       []carID{6, 7, 8},
					CustomStructSlice: []customStruct{
						{
							isLoved: false,
							isHated: true,
						},
					},
				}
			},
			mergedStruct: mergeStruct{
				Name:         "Name",
				Surname:      "Surname",
				SomeID:       3,
				IgnoredField: "ignored",
				CarIDs:       []carID{6, 7, 8},
				CustomStructSlice: []customStruct{
					{
						isLoved: false,
						isHated: true,
					},
				},
			},
		},
		"nil struct merge": {
			mstruct: mergeStruct{
				Name:         "Name",
				SomeID:       5,
				IgnoredField: "ignored",
				CarIDs:       []carID{1, 2, 3},
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				someID := int64(3)
				ignored := "ignored-update"
				return mergeCmd{
					Surname:      &surname,
					SomeID:       &someID,
					IgnoredField: &ignored,
					CarIDs:       []carID{6, 7, 8},
					CustomStruct: nil,
				}
			},
			mergedStruct: mergeStruct{
				Name:         "Name",
				Surname:      "Surname",
				SomeID:       3,
				IgnoredField: "ignored",
				CarIDs:       []carID{6, 7, 8},
			},
		},
		"empty slice merge": {
			mstruct: mergeStruct{
				Name:         "Name",
				SomeID:       5,
				IgnoredField: "ignored",
				CarIDs:       []carID{1, 2, 3},
			},
			cmd: func() mergeCmd {
				surname := "Surname"
				someID := int64(3)
				ignored := "ignored-update"
				return mergeCmd{
					Surname:           &surname,
					SomeID:            &someID,
					IgnoredField:      &ignored,
					CarIDs:            []carID{6, 7, 8},
					CustomStructSlice: []customStruct{},
				}
			},
			mergedStruct: mergeStruct{
				Name:              "Name",
				Surname:           "Surname",
				SomeID:            3,
				IgnoredField:      "ignored",
				CarIDs:            []carID{6, 7, 8},
				CustomStructSlice: []customStruct{},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cmd := c.cmd()
			structs.Merge(&c.mstruct, &cmd)
			assert.Equal(t, c.mstruct, c.mergedStruct)
		})
	}
}
