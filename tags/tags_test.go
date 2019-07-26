package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTag(t *testing.T) {
	tag := New("k1", "val")
	assert.Equal(t, "k1", tag.Key())
	assert.Equal(t, "val", tag.Value())
}

func TestOf(t *testing.T)  {
	tags, e := Of("k1", "v1", "k2", "v2")
	assert.Nil(t, e)
	assert.Equal(t, 2, len(tags))

	assert.Equal(t, "k1", tags[0].Key())
	assert.Equal(t, "v1", tags[0].Value())
	assert.Equal(t, "k2", tags[1].Key())
	assert.Equal(t, "v2", tags[1].Value())
}

func TestOfAcceptOnlyEvent(t *testing.T)  {
	_, e := Of("1", "2", "3")
	assert.NotNil(t, e)
}

func TestAnd(t *testing.T)  {
	tags, _ := Of("k1", "v1", "k2", "v2")

	t1 := New("k3", "v3")
	t2:= New("k4", "v4")

	tags = tags.And(t1, t2)

	assert.Equal(t, "k1", tags[0].Key())
	assert.Equal(t, "v1", tags[0].Value())

	assert.Equal(t, "k2", tags[1].Key())
	assert.Equal(t, "v2", tags[1].Value())

	assert.Equal(t, "k3", tags[2].Key())
	assert.Equal(t, "v3", tags[2].Value())

	assert.Equal(t, "k4", tags[3].Key())
	assert.Equal(t, "v4", tags[3].Value())
}

func TestConcat(t *testing.T)  {
	t1, _ := Of("k1", "v1", "k2", "v2")
	t2, _ := Of("k3", "v3", "k4", "v4")

	tags := t1.Contact(t2)

	assert.Equal(t, "k1", tags[0].Key())
	assert.Equal(t, "v1", tags[0].Value())

	assert.Equal(t, "k2", tags[1].Key())
	assert.Equal(t, "v2", tags[1].Value())

	assert.Equal(t, "k3", tags[2].Key())
	assert.Equal(t, "v3", tags[2].Value())

	assert.Equal(t, "k4", tags[3].Key())
	assert.Equal(t, "v4", tags[3].Value())
}