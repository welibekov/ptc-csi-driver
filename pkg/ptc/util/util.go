package util

// Deref returns the value pointed to by the item pointer or a new zero value if nil.
func Deref[T any](item *T) T {
	if item == nil {
		return *new(T)
	}
	return *item
}

// Ptr returns a pointer to the item.
func Ptr[T any](item T) *T {
	return &item
}

func In[T comparable](what T, where []T) bool {
	for _, item := range where {
		if item == what {
			return true
		}
	}

	return false
}
