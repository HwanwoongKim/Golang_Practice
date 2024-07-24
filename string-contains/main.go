package main

import (
	"fmt"
	"sort"
)

// Simulated data structures for demonstration
var BBSGossipMember = [][]string{
	{"member1"},
	{"member2"},
	{"member3"},
}

var BBSSelfMember = struct {
	Endpoint string
}{
	Endpoint: "selfEndpoint",
}

// Function to get network members
func getNetworkMembers() ([]string, int) {
	var memberArr []string
	for _, member := range BBSGossipMember {
		memberArr = append(memberArr, member[0])
	}

	// Append self endpoint to network members endpoint array
	selfEndpoint := BBSSelfMember.Endpoint
	memberArr = append(memberArr, selfEndpoint)

	// Sort network member array
	sort.Strings(memberArr)
	memberLen := len(memberArr)

	return memberArr, memberLen
}

// Function to check if a string is in a slice of slices of strings
func contains(slice [][]string, value string) bool {
	for _, v := range slice {
		if len(v) > 0 && v[0] == value {
			return true
		}
	}
	return false
}

func main() {
	// Get previous network members
	prevMember, _ := getNetworkMembers()

	// Example data for wereInPrev
	wereInPrev := [][]string{
		{"peer0.org2.example.com:7056"},
	}

	I := 3 // Index to check
	fmt.Println(wereInPrev, prevMember[1])

	// Check if prevMember[I] is in wereInPrev
	if contains(wereInPrev, prevMember[I]) {
		fmt.Println("prevMember[I] is in wereInPrev")
	} else {
		fmt.Println("prevMember[I] is not in wereInPrev")
	}
}
