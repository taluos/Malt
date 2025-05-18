package metadata

import (
	"strings"
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string][]string

// New creates an MD from a given key-values map.
// Different the metadata.New from grpc.metadata, this function will not convert the key to lowercase.
// This function allow to input multi map[string] with the velue is []string
// In the return Metadata, a key may appear corresponding to multiple values.
func New(mds ...map[string][]string) Metadata {
	md := Metadata{}
	// Get metadata(type:map[string][]string) from mds
	for _, m := range mds {
		// Get degital key(string)-velues([]string) from m
		for k, vList := range m {
			// Get value from velues
			for _, v := range vList {
				md.Add(k, v)
			}
		}
	}
	return md
}

// Add adds the key, value pair to the header.
func (m Metadata) Add(key, value string) {
	if len(key) == 0 {
		return
	}

	m[strings.ToLower(key)] = append(m[strings.ToLower(key)], value)
}

// Get returns the value associated with the passed key.
func (m Metadata) Get(key string) string {
	v := m[strings.ToLower(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Set stores the key-value pair.
func (m Metadata) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	m[strings.ToLower(key)] = []string{value}
}

// Range iterate over element in metadata.
func (m Metadata) Range(f func(k string, v []string) bool) {
	for k, v := range m {
		if !f(k, v) {
			break
		}
	}
}

// Values returns a slice of values associated with the passed key.
func (m Metadata) Values(key string) []string {
	return m[strings.ToLower(key)]
}

// Clone returns a deep copy of Metadata
func (m Metadata) sollowClone() Metadata {
	md := make(Metadata, len(m))
	for k, v := range m {
		md[k] = v
	}
	return md
}

// DeepClone returns a deep copy of Metadata
func (md Metadata) DeepClone() Metadata {
	out := make(Metadata, len(md))
	for k, v := range md {
		out[k] = copyOf(v)
	}
	return out
}

func copyOf(v []string) []string {
	vals := make([]string, len(v))
	copy(vals, v)
	return vals
}
