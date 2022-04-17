package oak

import (
	"fmt"
	"testing"
)

func TestInitFailures(t *testing.T) {
	t.Run("BadConfig", func(t *testing.T) {
		c1 := NewWindow()
		err := c1.Init("", func(c Config) (Config, error) {
			return c, fmt.Errorf("whoops")
		})
		if err == nil {
			t.Fatal("expected error to cascade down from init")
		}
	})
	t.Run("ParseDebugLevel", func(t *testing.T) {
		c1 := NewWindow()
		err := c1.Init("", func(c Config) (Config, error) {
			c.Debug.Level = "bogus"
			return c, nil
		})
		if err == nil {
			t.Fatal("expected error parsing debug level")
		}
	})
	t.Run("SetLanguageString", func(t *testing.T) {
		c1 := NewWindow()
		err := c1.Init("", func(c Config) (Config, error) {
			c.Language = "bogus"
			return c, nil
		})
		if err == nil {
			t.Fatal("expected error parsing language string")
		}
	})
}
