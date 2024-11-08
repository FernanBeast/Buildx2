package buildflags

import (
	"encoding"
	"sync"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
	"github.com/zclconf/go-cty/cty/gocty"
)

type FromCtyValue interface {
	FromCtyValue(in cty.Value, path cty.Path) error
}

func (e *CacheOptionsEntry) FromCtyValue(in cty.Value, p cty.Path) error {
	conv, err := convert.Convert(in, cty.Map(cty.String))
	if err == nil {
		m := conv.AsValueMap()
		if err := getAndDelete(m, "type", &e.Type); err != nil {
			return err
		}
		e.Attrs = asMap(m)
		return e.validate(in)
	}
	return unmarshalTextFallback(in, e, err)
}

func (e *ExportEntry) FromCtyValue(in cty.Value, p cty.Path) error {
	conv, err := convert.Convert(in, cty.Map(cty.String))
	if err == nil {
		m := conv.AsValueMap()
		if err := getAndDelete(m, "type", &e.Type); err != nil {
			return err
		}
		if err := getAndDelete(m, "dest", &e.Destination); err != nil {
			return err
		}
		e.Attrs = asMap(m)
		return e.validate()
	}
	return unmarshalTextFallback(in, e, err)
}

var secretType = sync.OnceValue(func() cty.Type {
	return cty.ObjectWithOptionalAttrs(
		map[string]cty.Type{
			"id":  cty.String,
			"src": cty.String,
			"env": cty.String,
		},
		[]string{"id", "src", "env"},
	)
})

func (e *Secret) FromCtyValue(in cty.Value, p cty.Path) (err error) {
	conv, err := convert.Convert(in, secretType())
	if err == nil {
		if id := conv.GetAttr("id"); !id.IsNull() {
			e.ID = id.AsString()
		}
		if src := conv.GetAttr("src"); !src.IsNull() {
			e.FilePath = src.AsString()
		}
		if env := conv.GetAttr("env"); !env.IsNull() {
			e.Env = env.AsString()
		}
		return nil
	}
	return unmarshalTextFallback(in, e, err)
}

var sshType = sync.OnceValue(func() cty.Type {
	return cty.ObjectWithOptionalAttrs(
		map[string]cty.Type{
			"id":    cty.String,
			"paths": cty.List(cty.String),
		},
		[]string{"id", "paths"},
	)
})

func (e *SSH) FromCtyValue(in cty.Value, p cty.Path) (err error) {
	conv, err := convert.Convert(in, sshType())
	if err == nil {
		if id := conv.GetAttr("id"); !id.IsNull() {
			e.ID = id.AsString()
		}
		if paths := conv.GetAttr("paths"); !paths.IsNull() {
			if err := gocty.FromCtyValue(paths, &e.Paths); err != nil {
				return err
			}
		}
		return nil
	}
	return unmarshalTextFallback(in, e, err)
}

func getAndDelete(m map[string]cty.Value, attr string, gv interface{}) error {
	if v, ok := m[attr]; ok {
		delete(m, attr)
		return gocty.FromCtyValue(v, gv)
	}
	return nil
}

func asMap(m map[string]cty.Value) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v.AsString()
	}
	return out
}

func unmarshalTextFallback[V encoding.TextUnmarshaler](in cty.Value, v V, inErr error) (outErr error) {
	// Attempt to convert this type to a string.
	conv, err := convert.Convert(in, cty.String)
	if err != nil {
		// Cannot convert. Do not attempt to convert and return the original error.
		return inErr
	}

	// Conversion was successful. Use UnmarshalText on the string and return any
	// errors associated with that.
	return v.UnmarshalText([]byte(conv.AsString()))
}
