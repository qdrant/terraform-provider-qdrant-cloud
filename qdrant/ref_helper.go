package qdrant

import "time"

func newString(s string) *string {
	return &s
}

func derefString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func derefStringArray(ptr *[]string) []string {
	if ptr != nil {
		return *ptr
	}
	return nil
}

func newInt(i int) *int {
	return &i
}

func derefInt(ptr *int) int {
	if ptr != nil {
		return *ptr
	}
	return 0
}

func newTime(t time.Time) *time.Time {
	return &t
}

func derefTime(ptr *time.Time) time.Time {
	if ptr != nil {
		return *ptr
	}
	return time.Time{}
}
