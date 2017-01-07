package smsender

import (
	"testing"
)

func TestReorderRoutes(t *testing.T) {
	var (
		dummyBroker = NewDummyBroker("dummy")
		router      = Router{}
	)

	for _, r := range []string{"D", "C", "B", "A"} {
		router.AddRoute(NewRoute(r, "", dummyBroker))
	}

	if err := router.reorder(-1, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.reorder(4, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.reorder(1, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.reorder(0, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.reorder(1, 4, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.reorder(0, 1, -1); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.reorder(0, 1, 5); err == nil {
		t.Fatal("got incorrect error: nil")
	}

	checkReorderRoutes(t, router, 1, 2, 3, []string{"A", "B", "C", "D"})
	checkReorderRoutes(t, router, 2, 2, 1, []string{"A", "C", "D", "B"})
	checkReorderRoutes(t, router, 0, 2, 4, []string{"C", "D", "A", "B"})
}

func checkReorderRoutes(t *testing.T, router Router, rangeStart, rangeLength, insertBefore int, expected []string) {
	err := router.reorder(rangeStart, rangeLength, insertBefore)
	if err != nil {
		t.Fatalf("reorder routes error: %v", err)
	}

	got := []string{}
	isNotMatch := false
	for i, route := range router.routes {
		got = append(got, route.Name)
		if route.Name != expected[i] {
			isNotMatch = true
		}
	}
	if isNotMatch {
		t.Fatalf("routes %v expected, but got %v", expected, got)
	}
}
