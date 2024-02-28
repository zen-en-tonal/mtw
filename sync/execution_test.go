package sync

import (
	"errors"
	"testing"
)

func ok(_ interface{}) error {
	return nil
}

func fails(_ interface{}) error {
	return errors.New("")
}

func TestTryAll_OneFunction(t *testing.T) {
	if err := TryAll(nil, fails); err == nil {
		t.Error("shuold failed")
	}
}

func TestTryAll_Function_Ok(t *testing.T) {
	if err := TryAll(nil, ok, ok); err != nil {
		t.Error(err)
	}
}

func TestTryAll_Function_fail(t *testing.T) {
	if err := TryAll(nil, ok, fails); err == nil {
		t.Error("shuold failed")
	}
}

func TestTrySome_OneFunction(t *testing.T) {
	if err := TrySome(nil, fails); err == nil {
		t.Error("shuold failed")
	}
}

func TestTrySome_Functions_Ok(t *testing.T) {
	if err := TrySome(nil, ok, ok); err != nil {
		t.Error(err)
	}
}

func TestTrySome_Function_fail(t *testing.T) {
	if err := TrySome(nil, ok, fails); err == nil {
		t.Error("shuold failed")
	}
}
