package runtime

import (
	"github.com/fudali113/good-job/storage"
	"os"
	"testing"
)

func TestCreateK8sRuntime(t *testing.T) {
	tpl, err := createTemplate()
	if err != nil {
		t.Error(err)
	}
	err = tpl.Execute(os.Stdout, map[string]interface{}{
		"config": storage.ExecConfig{
			Image: "oo",
			Cmd:   []string{"11", "22"},
			Args:  []string{"11", "22"},
			Env:   []string{"11", "22"},
		},
		"pipeline": "pipeline1",
		"job":      "job1",
		"name":     "pipeline1-job1",
	})
	if err != nil {
		t.Error(err)
	}
}
