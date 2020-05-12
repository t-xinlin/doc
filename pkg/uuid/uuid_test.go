package uuid

import "testing"

func TestNewMD5(t *testing.T) {
	uuid := New()
	u := NewMD5(uuid, []byte("Hello"))
	t.Log(u.String())
}
