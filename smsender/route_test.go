package smsender

import (
	"testing"
)

func TestReorderRoutes(t *testing.T) {
	var (
		dummyBroker = NewDummyBroker("dummy")
		routes      []*Route
	)

	for _, r := range []string{"A", "B", "C", "D"} {
		routes = append(routes, NewRoute(r, "", dummyBroker))
	}

	if _, err := reorderRoutes(routes, -1, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if _, err := reorderRoutes(routes, 4, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if _, err := reorderRoutes(routes, 1, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if _, err := reorderRoutes(routes, 0, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if _, err := reorderRoutes(routes, 1, 4, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if _, err := reorderRoutes(routes, 0, 1, -1); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if _, err := reorderRoutes(routes, 0, 1, 5); err == nil {
		t.Fatal("got incorrect error: nil")
	}

	checkReorderRoutes(t, routes, 1, 2, 3, []string{"A", "B", "C", "D"})
	checkReorderRoutes(t, routes, 2, 2, 1, []string{"A", "C", "D", "B"})
	checkReorderRoutes(t, routes, 0, 2, 4, []string{"C", "D", "A", "B"})
}

func checkReorderRoutes(t *testing.T, routes []*Route, rangeStart, rangeLength, insertBefore int, expected []string) {
	reordered, err := reorderRoutes(routes, rangeStart, rangeLength, insertBefore)
	if err != nil {
		t.Fatalf("reorder routes error: %v", err)
	}

	got := []string{}
	isNotMatch := false
	for i, route := range reordered {
		got = append(got, route.Name)
		if route.Name != expected[i] {
			isNotMatch = true
		}
	}
	if isNotMatch {
		t.Fatalf("routes %v expected, but got %v", expected, got)
	}
}
