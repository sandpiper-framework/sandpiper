package sandpiper_test

import (
	"testing"

	"autocare.org/sandpiper/internal/model"
)

func TestPaginationTransform(t *testing.T) {
	p := &sandpiper.PaginationReq{
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
