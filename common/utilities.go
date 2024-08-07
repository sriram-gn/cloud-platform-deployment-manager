/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright(c) 2019-2024 Wind River Systems, Inc. */

package common

import (
	"github.com/gophercloud/gophercloud/starlingx/inventory/v1/addresspools"
	"github.com/gophercloud/gophercloud/starlingx/inventory/v1/networks"
	"net"
	"regexp"
)

// Determines if an address is an IPv4 address
func IsIPv4(address string) bool {
	x := net.ParseIP(address)
	if x != nil {
		if x.To4() != nil {
			return true
		}
	}
	return false
}

// Determines if an address is an IPv6 address
func IsIPv6(address string) bool {
	x := net.ParseIP(address)
	if x != nil {
		// The net package does not have a good way to determine if an address
		// is definitely an IPv6 address.  The best it can do at the moment is
		// tell if it cleanly converts to an IPv4 value so we are going to
		// assume that if it parsed as an address and wasn't an IPv4 address
		// then it must be an IPv6 address.
		if x.To4() == nil {
			return true
		}
	}
	return false
}

// ListDelta is a utility function which calculates the difference between two
// lists.  If elements in 'b' are not present in 'a' then they will appear in
// the 'added' list.  If elements in a are not present in b then they will
// appear in the 'removed' list.
func ListDelta(a, b []string) (added []string, removed []string, same []string) {
	added = make([]string, 0)
	removed = make([]string, 0)
	same = make([]string, 0)
	present := make(map[string]bool)

	for _, s := range a {
		found := false
		for _, x := range b {
			if s == x {
				present[x] = true
				found = true
				break
			}
		}

		if !found {
			removed = append(removed, s)
		}
	}

	for _, x := range b {
		if !present[x] {
			added = append(added, x)
		} else {
			same = append(same, x)
		}
	}

	return added, removed, same
}

// ListChanged is a utility function which determines if a list of names
// provided in a spec is equivalent to the list of names return by the system
// API.  Since the spec accepts nil as a list that wasn't specified we consider
// the nil case as an empty list when comparing against the system API.
func ListChanged(a, b []string) bool {
	if len(a) != len(b) {
		return true
	}

	added, removed, _ := ListDelta(a, b)

	return len(added) > 0 || len(removed) > 0
}

// ListIntersect is a utility function which determines if there is any
// commonality between two lists of strings.
func ListIntersect(a []string, b []string) ([]string, bool) {
	result := make([]string, 0)

	for _, x := range a {
		for _, y := range b {
			if x == y {
				result = append(result, y)
			}
		}
	}

	return result, len(result) > 0
}

// ComparePartitionPaths is a utility function that compares the disk portion
// of two partition paths.  It returns true if the disk portion of a and b
// match.
func ComparePartitionPaths(a, b string) bool {
	re := regexp.MustCompile("-part[0-9]*")
	return re.ReplaceAllString(a, "") == re.ReplaceAllString(b, "")
}

// ContainsString is a utility function that determines whether a string is
// included in the list of elements of a slice.
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// RemoveString is a utility function that removes a string from the list of
// elements of a slice.
func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// DedupeSlice is a utility function that removes a duplicated element from
// a slice.
func DedupeSlice[T comparable](sliceList []T) []T {
	dedupeMap := make(map[T]bool)
	list := []T{}

	for _, item := range sliceList {
		if _, value := dedupeMap[item]; !value {
			dedupeMap[item] = true
			list = append(list, item)
		}
	}

	return list
}

func NormalizeIPv6(address string) string {
	var ret_ip string
	ip := net.ParseIP(address)
	if ip == nil {
		// The address doesn't seem to be valid IPv6.
		// Returning the address as is because semantic checks are
		// beyond the scope of this function.
		// Use it with along with "IsIPv6" function instead.
		return address
	}
	ret_ip = ip.String()
	return ret_ip
}

// IsIPAddressSame compares two addresses and returns
// true when IP addresses are same and false when they
// are not the same.
func IsIPAddressSame(ip1 string, ip2 string) bool {
	if ip1 != "" && IsIPv6(ip1) {
		// Normalization of IP address is required in case of
		// IPv6 due to different notations that could be employed.
		// Example: fcff:1::0 is same as fcff:1:0:0:0::0 but string
		// comparison would fail if IPs are not normalized.
		ip1 = NormalizeIPv6(ip1)
	}

	if ip2 != "" && IsIPv6(ip2) {
		ip2 = NormalizeIPv6(ip2)
	}

	return ip1 == ip2
}

func GetSystemAddrPoolByUUID(addrpool_list []addresspools.AddressPool, addrpool_uuid string) *addresspools.AddressPool {
	for _, addrpool := range addrpool_list {
		if addrpool.ID == addrpool_uuid {
			return &addrpool
		}
	}
	return nil
}

func GetSystemAddrPoolByName(addrpool_list []addresspools.AddressPool, addrpool_name string) *addresspools.AddressPool {
	for _, addrpool := range addrpool_list {
		if addrpool.Name == addrpool_name {
			return &addrpool
		}
	}
	return nil
}

func GetSystemNetworkByUUID(network_list []networks.Network, network_uuid string) *networks.Network {
	for _, network := range network_list {
		if network.UUID == network_uuid {
			return &network
		}
	}
	return nil
}

func GetSystemNetworkByName(network_list []networks.Network, network_name string) *networks.Network {
	for _, network := range network_list {
		if network.Name == network_name {
			return &network
		}
	}
	return nil
}

func GetSystemNetworkByType(network_list []networks.Network, network_type string) *networks.Network {
	for _, network := range network_list {
		if network.Type == network_type {
			return &network
		}
	}
	return nil
}
