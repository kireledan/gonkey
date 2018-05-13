package utils

import (
	"testing"
)

func TestInventoryLoad(t *testing.T) {
	hosts := parseInventory("./test_assets/test_inventory")
	if hosts[0].ip != "10.25.168.188:22" {
		t.Errorf("Parsed IP was incorrect, got: %s, wanted %s", hosts[0].ip, "10.25.168.188:22")
	}
	if hosts[0].user != "ubuntu" {
		t.Errorf("Parsed user was incorrect, got: %s, wanted %s", hosts[0].user, "ubuntu")
	}
}
