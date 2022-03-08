/*
Copyright (c) 2022 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package bump

import (
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
)

// FileBump defines how a version within a file will be matched through a regex
// and bumped using the provided version
type FileBump struct {
	Path    string
	Regex   string
	Version string
	Count   int
	SemVer  bool
}

// Task for bumping versions within files
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "bump"
}

// Skip running the task if no version has changed
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged || ctx.SkipBumps
}

// Run the task bumping the semantic version of any file identified within
// the uplift configuration file
func (t Task) Run(ctx *context.Context) error {
	if len(ctx.Config.Bumps) == 0 {
		log.Info("no files to bump")
		return nil
	}

	n := 0
	for _, bump := range ctx.Config.Bumps {
		ok, err := regexBump(ctx, bump.File, bump.Regex)
		if err != nil {
			return err
		}

		if ok {
			// Attempt to stage the changed file, if this isn't a dry run
			if !ctx.DryRun {
				if err := git.Stage(bump.File); err != nil {
					return err
				}
				log.WithField("file", bump.File).Info("successfully staged file")
			}
			n++
		}
	}

	if n > 0 {
		log.WithField("count", n).Debug("bumped files staged")
	}

	return nil
}
