package hclparser

import (
	"reflect"
	"sync"

	"github.com/containerd/errdefs"
	"github.com/zclconf/go-cty/cty"
)

type FromCtyValue interface {
	FromCtyValue(in cty.Value, path cty.Path) error
}

func impliedTypeExt(rt reflect.Type, _ cty.Path) (cty.Type, error) {
	if rt.AssignableTo(fromCtyValueType) {
		return fromCtyValueCapsuleType(rt), nil
	}
	return cty.NilType, errdefs.ErrNotImplemented
}

var (
	fromCtyValueType  = reflect.TypeFor[FromCtyValue]()
	fromCtyValueTypes sync.Map
)

func fromCtyValueCapsuleType(rt reflect.Type) cty.Type {
	if val, loaded := fromCtyValueTypes.Load(rt); loaded {
		return val.(cty.Type)
	}

	// First time used.
	ety := cty.CapsuleWithOps(rt.Name(), rt.Elem(), &cty.CapsuleOps{
		ConversionTo: func(_ cty.Type) func(cty.Value, cty.Path) (interface{}, error) {
			return func(in cty.Value, p cty.Path) (interface{}, error) {
				rv := reflect.New(rt.Elem()).Interface()
				if err := rv.(FromCtyValue).FromCtyValue(in, p); err != nil {
					return nil, err
				}
				return rv, nil
			}
		},
	})

	// Attempt to store the new type. Use whichever was loaded first
	// in the case of a race condition.
	val, _ := fromCtyValueTypes.LoadOrStore(rt, ety)
	return val.(cty.Type)
}
