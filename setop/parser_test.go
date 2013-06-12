package setop

import (
	"reflect"
	"testing"
)

func TestSetOpParser(t *testing.T) {
	var x2 float64 = 2
	var x3 float64 = 3
	op, err := NewSetOpParser("(U (I ccc aa (D ffff*2 gg)*3) (I:ConCat c23 b_ff) (X dbla e&44))").Parse()
	if err != nil {
		t.Error(err)
	}
	cmp := &SetOp{
		Type: Union,
		Sources: []SetOpSource{
			SetOpSource{
				SetOp: &SetOp{
					Type: Intersection,
					Sources: []SetOpSource{
						SetOpSource{Key: []byte("ccc")},
						SetOpSource{Key: []byte("aa")},
						SetOpSource{
							SetOp: &SetOp{
								Type: Difference,
								Sources: []SetOpSource{
									SetOpSource{Key: []byte("ffff"), Weight: &x2},
									SetOpSource{Key: []byte("gg")},
								},
							},
							Weight: &x3,
						},
					},
				},
			},
			SetOpSource{
				SetOp: &SetOp{
					Type:  Intersection,
					Merge: ConCat,
					Sources: []SetOpSource{
						SetOpSource{Key: []byte("c23")},
						SetOpSource{Key: []byte("b_ff")},
					},
				},
			},
			SetOpSource{
				SetOp: &SetOp{
					Type: Xor,
					Sources: []SetOpSource{
						SetOpSource{Key: []byte("dbla")},
						SetOpSource{Key: []byte("e&44")},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(op, cmp) {
		t.Errorf("%v and %v should be equal", op, cmp)
	}
}
