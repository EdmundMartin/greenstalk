package protocol

import "testing"

func stringSlicesIdentical(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for idx, val := range s1 {
		if s2[idx] != val {
			return false
		}
	}
	return true
}

func TestClientConn_findTube(t *testing.T) {
	conn := &ClientConn{Watching: []string{
		"a", "b", "c", "d", "f", "g",
	}}
	testChar := "e"
	result := conn.findTube(testChar)
	if result != -1 {
		t.Errorf("char %s is not contained in watched tubes", testChar)
	}
	testChar = "b"
	result = conn.findTube(testChar)
	if result != 1 {
		t.Errorf("char %s expected at index 1", testChar)
	}
	testChar = "g"
	result = conn.findTube(testChar)
	if result != len(conn.Watching) - 1 {
		t.Errorf("char %s expected at index %d", testChar, len(conn.Watching)-1)
	}
}

func TestClientConn_insertWatchingAlreadyExists(t *testing.T) {
	conn := &ClientConn{Watching: []string{
		"a", "b", "c",  "d",
	}}
	testChar := "b"
	conn.insertWatching(testChar)
	if len(conn.Watching) != 4 {
		t.Errorf("excepcted watching to be length of 4, instead got length %d", len(conn.Watching))
	}
}

func TestClientConn_insertWatchingMaintainsOrder(t *testing.T) {
	conn := &ClientConn{Watching: []string{
		"a", "c", "e", "f",
	}}
	testChar := "b"
	conn.insertWatching(testChar)
	if !stringSlicesIdentical(conn.Watching, []string{"a", "b", "c", "e", "f"}) {
		t.Errorf("expected compared slices to be identical")
	}
}

