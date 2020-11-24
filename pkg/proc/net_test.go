package proc

import "testing"

func Test_extractSocketInode(t *testing.T) {
	exp := "51212"
	link := "socket:[51212]"
	inode := extractSocketInode(link)
	if inode != exp {
		t.Errorf("Expected %v but got %v", exp, inode)
	}
}
