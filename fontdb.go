package resvg

import "github.com/kanrichan/resvg-go/internal"

// FontDB font database
type FontDB struct {
	wk  *Worker
	ptr int32
}

// NewFontDBDefault new a empty `FontDB` object in wasm.
// `FontDB` are not goroutine-safe, don't forget to close!
func (wk *Worker) NewFontDBDefault() (*FontDB, error) {
	if !wk.used.CompareAndSwap(false, true) {
		return nil, ErrWorkerIsBeingUsed
	}
	defer wk.used.Store(false)
	db, err := internal.FontdbDatabaseDefault(wk.ctx, wk.mod)
	if err != nil {
		return nil, err
	}
	return &FontDB{wk, db}, nil
}

// Close cloes the `FontDB` and recovers memory.
func (db *FontDB) Close() error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	err := internal.FontdbDatabaseDelete(db.wk.ctx, db.wk.mod, db.ptr)
	if err != nil {
		return err
	}
	db.ptr = 0
	return nil
}

// LoadFontFile loads font file into the `FontDB`.
func (db *FontDB) LoadFontFile(file string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseLoadFontFile(db.wk.ctx, db.wk.mod, db.ptr, file)
}

// LoadFontsDir loads font files from the selected directory into the `FontDB`.
func (db *FontDB) LoadFontsDir(dir string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseLoadFontsDir(db.wk.ctx, db.wk.mod, db.ptr, dir)
}

// LoadFontData loads font data into the `FontDB`.
func (db *FontDB) LoadFontData(data []byte) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseLoadFontData(db.wk.ctx, db.wk.mod, db.ptr, data)
}

// SetSerifFamily sets the family that will be used by `Family::Serif`.
func (db *FontDB) SetSerifFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetSerifFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

// SetSansSerifFamily sets the family that will be used by `Family::SansSerif`.
func (db *FontDB) SetSansSerifFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetSansSerifFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

// SetCursiveFamily sets the family that will be used by `Family::Cursive`.
func (db *FontDB) SetCursiveFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetCursiveFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

// SetFantasyFamily sets the family that will be used by `Family::Fantasy`.
func (db *FontDB) SetFantasyFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetFantasyFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

// SetMonospaceFamily sets the family that will be used by `Family::Monospace`.
func (db *FontDB) SetMonospaceFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetMonospaceFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

// Len returns the number of font faces in the `FontDB`
func (db *FontDB) Len() (int32, error) {
	if !db.wk.used.CompareAndSwap(false, true) {
		return 0, ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return 0, ErrPointerIsNil
	}
	return internal.FontdbDatabaseLen(db.wk.ctx, db.wk.mod, db.ptr)
}
