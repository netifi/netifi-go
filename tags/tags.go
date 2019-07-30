package tags

import "github.com/pkg/errors"

type Tag [2]string

func Empty() Tags {
	return make(Tags, 0)
}

func New(k string, v string) *Tag {
	return &Tag{k, v}
}

func (t *Tag) Key() string {
	return t[0]
}

func (t *Tag) Value() string {
	return t[1]
}

type Tags []*Tag

func Of(keyValues ...string) (Tags, error) {
	l := len(keyValues)
	if l == 0 {
		return Tags{}, nil
	}
	if l%2 == 1 {
		return nil, errors.New("size must be even, it is a set of key=value pairs")
	}
	t := make(Tags, l/2)
	for i := 0; i < l; i += 2 {
		t[i/2] = New(keyValues[i], keyValues[i+1])
	}
	return t, nil
}

func (t Tags) And(tags ...*Tag) Tags {
	if len(tags) == 0 {
		return t
	}

	return append(t, tags...)
}

func (t Tags) Contact(tags ...Tags) Tags {
	f := t
	for _, r := range tags {
		f = append(f, r...)
	}

	return f
}
