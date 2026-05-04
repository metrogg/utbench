package evaluator

import (
	"os"
	"testing"
)

func TestSplitJavaSourceByClasses_InnerClasses(t *testing.T) {
	// Read the actual interface_mock_000.java which has inner classes
	data, err := os.ReadFile("../../datasets/java/java_code_files_self_contained/interface_mock/interface_mock_000.java")
	if err != nil {
		t.Skip("dataset not available:", err)
	}
	source := string(data)

	names := extractAllClassNamesFromSource(source)
	t.Logf("extractAllClassNames: %v", names)

	// Should only find the top-level class, not inner classes
	if len(names) != 1 {
		t.Errorf("expected 1 top-level class, got %d: %v", len(names), names)
	}
	if len(names) > 0 && names[0] != "AuthenticationService" {
		t.Errorf("expected AuthenticationService, got %s", names[0])
	}

	classes := splitJavaSourceByClasses(source)
	t.Logf("splitJavaSourceByClasses returned %d classes", len(classes))

	// Should only have 1 top-level class
	if len(classes) != 1 {
		t.Errorf("expected 1 class in split result, got %d", len(classes))
		for name := range classes {
			t.Logf("  found class: %s", name)
		}
	}

	if _, ok := classes["AuthenticationService"]; !ok {
		t.Error("expected AuthenticationService in split result")
	}

	// The source for AuthenticationService should contain the inner classes
	if src, ok := classes["AuthenticationService"]; ok {
		if len(src) < 500 {
			t.Errorf("AuthenticationService source too short (%d bytes), inner classes likely missing", len(src))
		}
		t.Logf("AuthenticationService source: %d bytes", len(src))
	}
}

func TestSplitJavaSourceByClasses_SingleClass(t *testing.T) {
	source := `import java.util.*;

class Solution {
    public int add(int a, int b) { return a + b; }
}`
	names := extractAllClassNamesFromSource(source)
	if len(names) != 1 || names[0] != "Solution" {
		t.Errorf("expected [Solution], got %v", names)
	}
}

func TestSplitJavaSourceByClasses_MultipleTopLevel(t *testing.T) {
	source := `import java.util.*;

class Animal {
    String name;
}

class Dog extends Animal {
    void bark() {}
}

class Cat extends Animal {
    void meow() {}
}`
	names := extractAllClassNamesFromSource(source)
	if len(names) != 3 {
		t.Errorf("expected 3 top-level classes, got %d: %v", len(names), names)
	}

	classes := splitJavaSourceByClasses(source)
	if len(classes) != 3 {
		t.Errorf("expected 3 classes in split, got %d", len(classes))
	}
	for _, name := range []string{"Animal", "Dog", "Cat"} {
		if _, ok := classes[name]; !ok {
			t.Errorf("missing class %s in split result", name)
		}
	}
}
