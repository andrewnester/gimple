package gimple

import "testing"

func TestGetSetService(t *testing.T) {
	gimple := New()

	_, err := gimple.GetService("test")
	if err == nil {
		t.Error("Service expected not to exist but it does")
	}

	expected := "My Service"
	gimple.SetService("test", func(g *Gimple) interface{} {
		return expected
	})

	actual, err := gimple.GetService("test")
	if err != nil {
		t.Error("Service expected to exist but it doesn't")
	}

	if actual != expected {
		t.Error("Expected and actual returned services mismatched")
	}

	// more complex example
	gimple.SetService("math-service", func(g *Gimple) interface{} {
		return math{1, 2}
	})

	actual, err = gimple.GetService("math-service")
	if err != nil {
		t.Error("Service expected to exist but it doesn't")
	}

	service, ok := actual.(math)
	if !ok {
		t.Error("Wrong service type")
	}

	if service.sum() != 3 {
		t.Errorf("Expected and actual returned services mismatched - %d is not equal to 3", service.sum())
	}
}

func TestServiceExists(t *testing.T) {
	gimple := New()

	if gimple.ServiceExists("test") {
		t.Error("Service expected not to exist but it does")
	}

	gimple.SetService("test", func(g *Gimple) interface{} {
		return "My Service"
	})

	if !gimple.ServiceExists("test") {
		t.Error("Service expected to exist but it doesn't")
	}
}

func TestUnsetService(t *testing.T) {
	gimple := New()

	gimple.SetService("test", func(g *Gimple) interface{} {
		return "My Service"
	})

	gimple.UnsetService("test")

	if gimple.ServiceExists("test") {
		t.Error("Service expected not to exist but it does")
	}
}

func TestShare(t *testing.T) {
	gimple := New()

	expected := math{1, 2}
	factory := gimple.Share(func(g *Gimple) interface{} {
		return expected
	})

	actual1 := factory(nil)
	actual2 := factory(nil)

	if actual1 != actual2 || actual1 != expected {
		t.Error("Share returned different object")
	}
}

func TestProtect(t *testing.T) {
	gimple := New()

	callback := func(g *Gimple) interface{} {
		return "My Service"
	}

	factory := gimple.Protect(callback)
	actual1, ok1 := factory(nil).(Callable)
	actual2, ok2 := factory(nil).(Callable)

	if !ok1 || !ok2 {
		t.Error("Protect returned different object")
	}

	if actual1(nil) != actual2(nil) || actual1(nil) != callback(nil) {
		t.Error("Protect returned different object")
	}
}

func TestRaw(t *testing.T) {
	gimple := New()

	_, err := gimple.Raw("math-service")
	if err == nil {
		t.Error("Service expected not to exist but it does")
	}

	factory := func(g *Gimple) interface{} {
		return "My Service"
	}
	gimple.SetService("math-service", factory)
	gimple.SetService("math-service2", func(g *Gimple) interface{} {
		return "My Service 2"
	})

	actual, err := gimple.Raw("math-service")
	if err != nil {
		t.Error("Service must exist but it doesn't")
	}

	if factory(nil) != actual(nil) {
		t.Error("Raw returned wrong object")
	}
}

func TestExtend(t *testing.T) {

	gimple := New()

	extender := func(origin interface{}, context *Gimple) interface{} {
		m := origin.(math)
		return math{m.sum(), 10}
	}

	_, err := gimple.Extend("math-service", extender)
	if err == nil {
		t.Error("Service must not exist but it does")
	}

	factory := func(g *Gimple) interface{} {
		return math{1, 2}
	}
	gimple.SetService("math-service", factory)

	actual, err := gimple.Extend("math-service", extender)
	if err != nil {
		t.Error("Service must exist but it doesn't")
	}

	m := actual(nil).(math)
	if m.sum() != 13 {
		t.Error("Wrong service returned after extending")
	}
}

func TestKeys(t *testing.T) {
	gimple := New()

	keys := gimple.Keys()
	if len(keys) != 0 {
		t.Error("Wrong keys returned")
	}

	gimple.SetService("test 1", func(g *Gimple) interface{} {
		return "test 1"
	})
	gimple.SetService("test 2", func(g *Gimple) interface{} {
		return "test 2"
	})
	gimple.SetService("test 1", func(g *Gimple) interface{} {
		return "test 111"
	})

	keys = gimple.Keys()
	if len(keys) != 2 {
		t.Error("Wrong keys returned")
	}

	if !inSlice("test 1", keys) {
		t.Error("Wrong keys returned")
	}

	if !inSlice("test 2", keys) {
		t.Error("Wrong keys returned")
	}

	gimple.UnsetService("test 2")

	keys = gimple.Keys()
	if len(keys) != 1 {
		t.Error("Wrong keys returned")
	}

	if !inSlice("test 1", keys) {
		t.Error("Wrong keys returned")
	}
}

func inSlice(s string, slice []string) bool {
	for _, value := range slice {
		if s == value {
			return true
		}
	}
	return false
}

type math struct {
	a int
	b int
}

func (m *math) sum() int {
	return m.a + m.b
}
