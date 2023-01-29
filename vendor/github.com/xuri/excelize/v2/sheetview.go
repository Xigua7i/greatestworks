// Copyright 2016 - 2023 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// Package excelize providing a set of functions that allow you to write to and
// read from XLAM / XLSM / XLSX / XLTM / XLTX files. Supports reading and
// writing spreadsheet documents generated by Microsoft Excel™ 2007 and later.
// Supports complex components by high compatibility, and provided streaming
// API for generating or reading data from a worksheet with huge amounts of
// data. This library needs Go version 1.16 or later.

package excelize

// getSheetView returns the SheetView object
func (f *File) getSheetView(sheet string, viewIndex int) (*xlsxSheetView, error) {
	ws, err := f.workSheetReader(sheet)
	if err != nil {
		return nil, err
	}
	if ws.SheetViews == nil {
		ws.SheetViews = &xlsxSheetViews{
			SheetView: []xlsxSheetView{{WorkbookViewID: 0}},
		}
	}
	if viewIndex < 0 {
		if viewIndex < -len(ws.SheetViews.SheetView) {
			return nil, newViewIdxError(viewIndex)
		}
		viewIndex = len(ws.SheetViews.SheetView) + viewIndex
	} else if viewIndex >= len(ws.SheetViews.SheetView) {
		return nil, newViewIdxError(viewIndex)
	}

	return &(ws.SheetViews.SheetView[viewIndex]), err
}

// setSheetView set sheet view by given options.
func (view *xlsxSheetView) setSheetView(opts *ViewOptions) {
	if opts.DefaultGridColor != nil {
		view.DefaultGridColor = opts.DefaultGridColor
	}
	if opts.RightToLeft != nil {
		view.RightToLeft = *opts.RightToLeft
	}
	if opts.ShowFormulas != nil {
		view.ShowFormulas = *opts.ShowFormulas
	}
	if opts.ShowGridLines != nil {
		view.ShowGridLines = opts.ShowGridLines
	}
	if opts.ShowRowColHeaders != nil {
		view.ShowRowColHeaders = opts.ShowRowColHeaders
	}
	if opts.ShowRuler != nil {
		view.ShowRuler = opts.ShowRuler
	}
	if opts.ShowZeros != nil {
		view.ShowZeros = opts.ShowZeros
	}
	if opts.TopLeftCell != nil {
		view.TopLeftCell = *opts.TopLeftCell
	}
	if opts.View != nil {
		if _, ok := map[string]interface{}{
			"normal":           nil,
			"pageLayout":       nil,
			"pageBreakPreview": nil,
		}[*opts.View]; ok {
			view.View = *opts.View
		}
	}
	if opts.ZoomScale != nil && *opts.ZoomScale >= 10 && *opts.ZoomScale <= 400 {
		view.ZoomScale = *opts.ZoomScale
	}
}

// SetSheetView sets sheet view options. The viewIndex may be negative and if
// so is counted backward (-1 is the last view).
func (f *File) SetSheetView(sheet string, viewIndex int, opts *ViewOptions) error {
	view, err := f.getSheetView(sheet, viewIndex)
	if err != nil {
		return err
	}
	if opts == nil {
		return err
	}
	view.setSheetView(opts)
	return nil
}

// GetSheetView gets the value of sheet view options. The viewIndex may be
// negative and if so is counted backward (-1 is the last view).
func (f *File) GetSheetView(sheet string, viewIndex int) (ViewOptions, error) {
	opts := ViewOptions{
		DefaultGridColor:  boolPtr(true),
		ShowFormulas:      boolPtr(true),
		ShowGridLines:     boolPtr(true),
		ShowRowColHeaders: boolPtr(true),
		ShowRuler:         boolPtr(true),
		ShowZeros:         boolPtr(true),
		View:              stringPtr("normal"),
		ZoomScale:         float64Ptr(100),
	}
	view, err := f.getSheetView(sheet, viewIndex)
	if err != nil {
		return opts, err
	}
	if view.DefaultGridColor != nil {
		opts.DefaultGridColor = view.DefaultGridColor
	}
	opts.RightToLeft = boolPtr(view.RightToLeft)
	opts.ShowFormulas = boolPtr(view.ShowFormulas)
	if view.ShowGridLines != nil {
		opts.ShowGridLines = view.ShowGridLines
	}
	if view.ShowRowColHeaders != nil {
		opts.ShowRowColHeaders = view.ShowRowColHeaders
	}
	if view.ShowRuler != nil {
		opts.ShowRuler = view.ShowRuler
	}
	if view.ShowZeros != nil {
		opts.ShowZeros = view.ShowZeros
	}
	opts.TopLeftCell = stringPtr(view.TopLeftCell)
	if view.View != "" {
		opts.View = stringPtr(view.View)
	}
	if view.ZoomScale >= 10 && view.ZoomScale <= 400 {
		opts.ZoomScale = float64Ptr(view.ZoomScale)
	}
	return opts, err
}
