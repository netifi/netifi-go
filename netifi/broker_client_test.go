package netifi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShouldIncludeUri(t *testing.T) {
	_, e := New().Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "uri or discovery service"))
}

func TestShouldIncludeAccessKeyGreaterThanZero(t *testing.T) {
	_, e := New().Uri("tcp://localhost").Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "key must be greater"))
}

func TestShouldIncludeAccessToken(t *testing.T) {
	_, e := New().AccessKey(123).Uri("tcp://localhost").Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "must include access token"))
}

func TestShouldIncludeGroup(t *testing.T) {
	_, e := New().AccessKey(123).AccessToken([]byte("token")).Uri("tcp://localhost").Build()
	if e == nil {
		t.Fail()
		return
	}
	fmt.Println(e)
	assert.True(t, strings.Contains(e.Error(), "must include a group"))
}

func TestBuild(t *testing.T) {
	_, e := New().
		AccessKey(123).
		AccessToken([]byte("token")).
		Uri("tcp://locahost").
		Group("test").
		Build()
	if e != nil {
		fmt.Println(e)
		t.Fail()
		return
	}
}
