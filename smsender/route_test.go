package smsender

import (
	"fmt"
	"testing"
)

func createRouter() Router {
	dummyBroker1 := NewDummyBroker("dummy1")
	dummyBroker2 := NewDummyBroker("dummy2")
	router := Router{}

	router.AddRoute(NewRoute("default", `^\+.*`, dummyBroker1))
	router.AddRoute(NewRoute("japan", `^\+81`, dummyBroker2))
	router.AddRoute(NewRoute("taiwan", `^\+886`, dummyBroker2))
	router.AddRoute(NewRoute("telco", `^\+886987`, dummyBroker2))
	router.AddRoute(NewRoute("user", `^\+886987654321`, dummyBroker2))

	return router
}

func compare(routes []*Route, expected []string) error {
	got := []string{}
	isNotMatch := false
	for i, route := range routes {
		got = append(got, route.Name)
		if route.Name != expected[i] {
			isNotMatch = true
		}
	}
	if isNotMatch {
		return fmt.Errorf("routes expecting %v, but got %v", expected, got)
	}
	return nil
}

func TestRouter_GetRoute(t *testing.T) {
	router := createRouter()

	route := router.GetRoute("japan")
	if route == nil || route.Name != "japan" {
		t.Fatal("got wrong route")
	}
	route = router.GetRoute("usa")
	if route != nil {
		t.Fatal("route should be nil")
	}
}

func TestRouter_removeRoute(t *testing.T) {
	router := createRouter()

	router.removeRoute(1)
	router.removeRoute(2)
	if len(router.routes) != 3 {
		t.Fatal("removeRoute failed")
	}
	if err := compare(router.routes, []string{"user", "taiwan", "default"}); err != nil {
		t.Fatal(err)
	}
}

func TestRouter_GetRoutes(t *testing.T) {
	router := createRouter()

	if err := compare(router.GetRoutes(), []string{"user", "telco", "taiwan", "japan", "default"}); err != nil {
		t.Fatal(err)
	}
}

func TestRouter_Reorder(t *testing.T) {
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

	if err := compare(router.routes, expected); err != nil {
		t.Fatal(err)
	}
}

type routeTest struct {
	phone       string
	shouldMatch bool
	route       string
	broker      string
}

func TestRouter_Match(t *testing.T) {
	router := createRouter()

	tests := []routeTest{
		{
			phone:       "+886987654321",
			shouldMatch: true,
			route:       "user",
			broker:      "dummy2",
		},
		{
			phone:       "+886987654322",
			shouldMatch: true,
			route:       "telco",
			broker:      "dummy2",
		},
		{
			phone:       "+886900000001",
			shouldMatch: true,
			route:       "taiwan",
			broker:      "dummy2",
		},
		{
			phone:       "+819000000001",
			shouldMatch: true,
			route:       "japan",
			broker:      "dummy2",
		},
		{
			phone:       "+10000000001",
			shouldMatch: true,
			route:       "default",
			broker:      "dummy1",
		},
		{
			phone:       "woo",
			shouldMatch: false,
			route:       "",
			broker:      "",
		},
	}

	for i, test := range tests {
		match, ok := router.Match(test.phone)
		if test.shouldMatch {
			if !ok {
				t.Fatalf("test '%d' should match", i)
			}
			if test.route != match.Name {
				t.Fatalf("test '%d' route not match", i)
			}
			if test.broker != match.Broker.Name() {
				t.Fatalf("test '%d' broker not match", i)
			}
		} else {
			if ok {
				t.Fatalf("test '%d' should not match", i)
			}
		}
	}
}
