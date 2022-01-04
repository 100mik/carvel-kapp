// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package clusterapply

import (
	"fmt"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	ctlconf "github.com/k14s/kapp/pkg/kapp/config"
	ctldiff "github.com/k14s/kapp/pkg/kapp/diff"
)

type ChangeSetViewOpts struct {
	Summary bool
	Changes bool
	ctldiff.TextDiffViewOpts
}

type ChangeSetView struct {
	changeViews []ChangeView
	maskRules   []ctlconf.DiffMaskRule
	opts        ChangeSetViewOpts

	changesView *ChangesView
}

func NewChangeSetView(changeViews []ChangeView,
	maskRules []ctlconf.DiffMaskRule, opts ChangeSetViewOpts) *ChangeSetView {

	return &ChangeSetView{changeViews, maskRules, opts, nil}
}

func (v *ChangeSetView) Print(ui ui.UI) {
	if v.opts.Changes {
		for _, view := range v.changeViews {
			textDiffView := ctldiff.NewTextDiffView(view.ConfigurableTextDiff(), v.maskRules, v.opts.TextDiffViewOpts)
			ui.BeginLinef("@@ %s %s @@\n", applyOpCodeUI[view.ApplyOp()], view.Resource().Description())
			ui.PrintBlock([]byte(textDiffView.String()))
		}
	}

	v.changesView = &ChangesView{ChangeViews: v.changeViews, Sort: true, countsView: NewChangesCountsView()}

	if v.opts.Summary {
		v.changesView.Print(ui)
	}
}

const (
	// 0.5 MB in bytes
	diffChangesStringMaxSize = 524288
)

func (v *ChangeSetView) DiffChangesString() string {
	var diffChangesBuffer strings.Builder
	for _, view := range v.changeViews {
		textDiffView := ctldiff.NewTextDiffView(view.ConfigurableTextDiff(), v.maskRules, v.opts.TextDiffViewOpts)
		diffChangesBuffer.Write([]byte(fmt.Sprintf("@@ %s %s @@\n", applyOpCodeUI[view.ApplyOp()], view.Resource().Description())))
		diffChangesBuffer.WriteString(textDiffView.String())
	}

	diffChangesString := diffChangesBuffer.String()

	if len(diffChangesString)*4 > diffChangesStringMaxSize {
		diffChangesString = fmt.Sprintf("%s\n%s", diffChangesString[:130000], "Diff truncated as it exceeded max size (0.5 MB) ...")
	}

	return diffChangesString
}

func (v *ChangeSetView) Summary() string {
	return v.changesView.Summary() // assumes Print was used before
}
