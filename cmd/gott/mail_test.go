package main

import (
	"testing"
)

func TestSendMail(t *testing.T) {
	err := sendMail("mrc1rz+8ca4wocgn0mz8@sharklasers.com", "Test Test 123", "Hello, this is just a test")
	if err != nil {
		t.Error(err)
	}
}
