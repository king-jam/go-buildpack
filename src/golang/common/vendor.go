package common

import (
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Godep struct {
	ImportPath      string   `json:"ImportPath"`
	GoVersion       string   `json:"GoVersion"`
	Packages        []string `json:"Packages"`
	WorkspaceExists bool
}

func SelectVendorTool(c *libbuildpack.Compiler, godep *Godep) (string, error) {
	godepsJSONFile := filepath.Join(c.BuildDir, "Godeps", "Godeps.json")

	// godirFile := filepath.Join(gc.Compiler.BuildDir, ".godir")
	// isGodir, err := libbuildpack.FileExists(godirFile)
	// if err != nil {
	// 	return err
	// }
	// if isGodir {
	// 	gs.Compiler.Log.Error(godirError())
	// 	return errors.New(".godir deprecated")
	// }

	// isGB, err := gs.isGB()
	// if err != nil {
	// 	return err
	// }
	// if isGB {
	// 	gs.Compiler.Log.Error(gbError())
	// 	return errors.New("gb unsupported")
	// }

	isGodep, err := libbuildpack.FileExists(godepsJSONFile)
	if err != nil {
		return "", err
	}
	if isGodep {
		c.Log.BeginStep("Checking Godeps/Godeps.json file")

		err = libbuildpack.NewJSON().Load(filepath.Join(c.BuildDir, "Godeps", "Godeps.json"), godep)
		if err != nil {
			c.Log.Error("Bad Godeps/Godeps.json file")
			return "", err
		}

		godep.WorkspaceExists, err = libbuildpack.FileExists(filepath.Join(c.BuildDir, "Godeps", "_workspace", "src"))
		if err != nil {
			return "", err
		}

		return "godep", nil
	}

	glideFile := filepath.Join(c.BuildDir, "glide.yaml")
	isGlide, err := libbuildpack.FileExists(glideFile)
	if err != nil {
		return "", err
	}
	if isGlide {
		return "glide", nil
	}

	return "go_nativevendoring", nil
}
