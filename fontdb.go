package resvg

import "github.com/kanrichan/resvg-go/internal"

type Fontdb struct {
	wk  *Worker
	ptr int32
}

func (wk *Worker) NewFontDBDefault() (*Fontdb, error) {
	if !wk.used.CompareAndSwap(false, true) {
		return nil, ErrWorkerIsBeingUsed
	}
	defer wk.used.Store(false)
	db, err := internal.FontdbDatabaseDefault(wk.ctx, wk.mod)
	if err != nil {
		return nil, err
	}
	return &Fontdb{wk, db}, nil
}

func (db *Fontdb) Close() error {
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

func (db *Fontdb) LoadFontFile(file string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseLoadFontFile(db.wk.ctx, db.wk.mod, db.ptr, file)
}

func (db *Fontdb) LoadFontsDir(dir string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseLoadFontsDir(db.wk.ctx, db.wk.mod, db.ptr, dir)
}

func (db *Fontdb) LoadFromData(data []byte) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseLoadFontData(db.wk.ctx, db.wk.mod, db.ptr, data)
}

func (db *Fontdb) SetSerifFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetSerifFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

func (db *Fontdb) SetSansSerifFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetSansSerifFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

func (db *Fontdb) SetCursiveFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetCursiveFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

func (db *Fontdb) SetFantasyFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetFantasyFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}

func (db *Fontdb) SetMonospaceFamily(family string) error {
	if !db.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer db.wk.used.Store(false)
	if db.ptr == 0 {
		return ErrPointerIsNil
	}
	return internal.FontdbDatabaseSetMonospaceFamily(db.wk.ctx, db.wk.mod, db.ptr, family)
}
