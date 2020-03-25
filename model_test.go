package gorsk_test

import (
	"context"
	"testing"

	"github.com/ribice/gorsk"
	"github.com/ribice/gorsk/pkg/utl/mock"
)

func TestBeforeInsert(t *testing.T) {
	base := &gorsk.Base{
		ID: 1,
	}
	base.BeforeInsert(context.TODO())
	if base.CreatedAt.IsZero() {
		t.Error("CreatedAt was not changed")
	}
	if base.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not changed")
	}
}

func TestBeforeUpdate(t *testing.T) {
	base := &gorsk.Base{
		ID:        1,
		CreatedAt: mock.TestTime(2000),
	}
	base.BeforeUpdate(context.TODO())
	if base.UpdatedAt == mock.TestTime(2001) {
		t.Error("UpdatedAt was not changed")
	}

}

func TestPaginationTransform(t *testing.T) {
	p := &gorsk.PaginationReq{
		Limit: 5000, Page: 5,
	}

	pag := p.Transform()

	if pag.Limit != 1000 {
		t.Error("Default limit not set")
	}

	if pag.Offset != 5000 {
		t.Error("Offset not set correctly")
	}

	p.Limit = 0
	newPag := p.Transform()

	if newPag.Limit != 100 {
		t.Error("Min limit not set")
	}

}
