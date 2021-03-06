// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package builder

import (
	bldr "github.com/arduino/arduino-cli/arduino/builder"
	"github.com/arduino/arduino-cli/arduino/sketch"
	"github.com/arduino/arduino-cli/legacy/builder/constants"
	"github.com/arduino/arduino-cli/legacy/builder/i18n"
	"github.com/arduino/arduino-cli/legacy/builder/types"
)

type ContainerAddPrototypes struct{}

func (s *ContainerAddPrototypes) Run(ctx *types.Context) error {
	// Generate the full pathname for the preproc output file
	if err := ctx.PreprocPath.MkdirAll(); err != nil {
		return i18n.WrapError(err)
	}
	targetFilePath := ctx.PreprocPath.Join(constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E)

	// Run preprocessor
	sourceFile := ctx.SketchBuildPath.Join(ctx.Sketch.MainFile.Name.Base() + ".cpp")
	if err := GCCPreprocRunner(ctx, sourceFile, targetFilePath, ctx.IncludeFolders); err != nil {
		return i18n.WrapError(err)
	}

	commands := []types.Command{
		&ReadFileAndStoreInContext{FileToRead: targetFilePath, Target: &ctx.SourceGccMinusE},
		&FilterSketchSource{Source: &ctx.SourceGccMinusE},
		&CTagsTargetFileSaver{Source: &ctx.SourceGccMinusE, TargetFileName: constants.FILE_CTAGS_TARGET_FOR_GCC_MINUS_E},
		&CTagsRunner{},
		&PrototypesAdder{},
	}

	for _, command := range commands {
		PrintRingNameIfDebug(ctx, command)
		err := command.Run(ctx)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	if err := bldr.SketchSaveItemCpp(&sketch.Item{ctx.Sketch.MainFile.Name.String(), []byte(ctx.Source)}, ctx.SketchBuildPath.String()); err != nil {
		return i18n.WrapError(err)
	}

	return nil
}
