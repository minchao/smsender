package smsender

import (
	"fmt"
	"testing"

	"github.com/minchao/smsender/smsender/model"
)

func createRouter() Router {
	dummyBroker1 := NewDummyBroker("dummy1")
	dummyBroker2 := NewDummyBroker("dummy2")
	router := Router{}

	router.Add(model.NewRoute("default", `^\+.*`, dummyBroker1))
	router.Add(model.NewRoute("japan", `^\+81`, dummyBroker2))
	router.Add(model.NewRoute("taiwan", `^\+886`, dummyBroker2))
	router.Add(model.NewRoute("telco", `^\+886987`, dummyBroker2))
	router.Add(model.NewRoute("user", `^\+886987654321`, dummyBroker2))

	return router
}

func compareOrder(routes []*model.Route, expected []string) error {
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

func TestRouter_GetAll(t *testing.T) {
	router := createRouter()

	if err := compareOrder(router.GetAll(), []string{"user", "telco", "taiwan", "japan", "default"}); err != nil {
		t.Fatal(err)
	}
}

func TestRouter_Get(t *testing.T) {
	router := createRouter()

	route := router.Get("japan")
	if route == nil || route.Name != "japan" {
		t.Fatal("got wrong route")
	}
	route = router.Get("usa")
	if route != nil {
		t.Fatal("route should be nil")
	}
}

func TestRouter_Set(t *testing.T) {
	router := createRouter()
	broker := NewDummyBroker("dummy")

	route := model.NewRoute("taiwan", `^\+8869`, broker).SetFrom("sender")

	if err := router.Set(route.Name, route.Pattern, route.GetBroker(), route.From); err == nil {
		newRoute := router.Get("taiwan")
		if newRoute == nil {
			t.Fatal("route is not equal")
		}
		if newRoute.Name != route.Name {
			t.Fatal("route.Name is not equal")
		}
		if newRoute.Pattern != route.Pattern {
			t.Fatal("route.Pattern is not equal")
		}
		if newRoute.GetBroker() == nil || newRoute.GetBroker().Name() != route.GetBroker().Name() {
			t.Fatal("route.Broker is not equal")
		}
		if newRoute.From != route.From {
			t.Fatal("route.From is not equal")
		}
	}
	if err := router.Set("france", "", broker, ""); err == nil {
		t.Fatal("set route should be failed")
	}
}

func TestRouter_Remove(t *testing.T) {
	router := createRouter()

	router.Remove("telco")
	router.Remove("japan")
	if len(router.routes) != 3 {
		t.Fatal("remove route failed")
	}
	if err := compareOrder(router.routes, []string{"user", "taiwan", "default"}); err != nil {
		t.Fatal(err)
	}
}

func TestRouter_Reorder(t *testing.T) {
	var (
		dummyBroker = NewDummyBroker("dummy")
		router      = Router{}
	)

	for _, r := range []string{"D", "C", "B", "A"} {
		router.Add(model.NewRoute(r, "", dummyBroker))
	}

	if err := router.Reorder(-1, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.Reorder(4, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.Reorder(1, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.Reorder(0, 0, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.Reorder(1, 4, 0); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.Reorder(0, 1, -1); err == nil {
		t.Fatal("got incorrect error: nil")
	}
	if err := router.Reorder(0, 1, 5); err == nil {
		t.Fatal("got incorrect error: nil")
	}

	checkReorderRoutes(t, router, 1, 2, 3, []string{"A", "B", "C", "D"})
	checkReorderRoutes(t, router, 2, 2, 1, []string{"A", "C", "D", "B"})
	checkReorderRoutes(t, router, 0, 2, 4, []string{"C", "D", "A", "B"})
}

func checkReorderRoutes(t *testing.T, router Router, rangeStart, rangeLength, insertBefore int, expected []string) {
	if err := router.Reorder(rangeStart, rangeLength, insertBefore); err != nil {
		t.Fatalf("reorder routes error: %v", err)
	}
	if err := compareOrder(router.routes, expected); err != nil {
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
				t.Fatalf("test '%d' route.Name is not equal", i)
			}
			if test.broker != match.GetBroker().Name() {
				t.Fatalf("test '%d' route.Broker is not equal", i)
			}
		} else {
			if ok {
				t.Fatalf("test '%d' should not match", i)
			}
		}
	}
}
