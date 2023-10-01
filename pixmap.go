package resvg

import "github.com/kanrichan/resvg-go/internal"

type Pixmap struct {
	wk  *Worker
	ptr int32
}

func (wk *Worker) NewPixmap(width uint32, height uint32) (*Pixmap, error) {
	if !wk.used.CompareAndSwap(false, true) {
		return nil, ErrWorkerIsBeingUsed
	}
	defer wk.used.Store(false)
	pm, err := internal.TinySkiaPixmapNew(wk.ctx, wk.mod, width, height)
	if err != nil {
		return nil, err
	}
	return &Pixmap{wk, pm}, nil
}

func (wk *Worker) NewPixmapDecodePNG(data []byte) (*Pixmap, error) {
	if !wk.used.CompareAndSwap(false, true) {
		return nil, ErrWorkerIsBeingUsed
	}
	defer wk.used.Store(false)
	pm, err := internal.TinySkiaPixmapDecodePNG(wk.ctx, wk.mod, data)
	if err != nil {
		return nil, err
	}
	return &Pixmap{wk, pm}, nil
}

func (pm *Pixmap) Close() error {
	if !pm.wk.used.CompareAndSwap(false, true) {
		return ErrWorkerIsBeingUsed
	}
	defer pm.wk.used.Store(false)
	if pm.ptr == 0 {
		return ErrPointerIsNil
	}
	err := internal.TinySkiaPixmapDelete(pm.wk.ctx, pm.wk.mod, pm.ptr)
	if err != nil {
		return err
	}
	pm.ptr = 0
	return nil
}

func (pm *Pixmap) EncodePNG() ([]byte, error) {
	if !pm.wk.used.CompareAndSwap(false, true) {
		return nil, ErrWorkerIsBeingUsed
	}
	defer pm.wk.used.Store(false)
	if pm.ptr == 0 {
		return nil, ErrPointerIsNil
	}
	return internal.TinySkiaPixmapEncodePng(pm.wk.ctx, pm.wk.mod, pm.ptr)
}
