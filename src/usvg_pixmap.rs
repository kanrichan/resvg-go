
#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_pixmap_new")]
#[no_mangle]
pub fn tiny_skia_pixmap_new(width: u32, height: u32) -> *mut tiny_skia::Pixmap {
    let pixmap = tiny_skia::Pixmap::new(width, height).unwrap();
    let ptr = Box::into_raw(pixmap.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_pixmap_free")]
#[no_mangle]
pub fn tiny_skia_pixmap_free(pixmap: *mut tiny_skia::Pixmap) {
    let _ = unsafe { Box::from_raw(pixmap) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_pixmap_encode_png")]
#[no_mangle]
pub fn tiny_skia_pixmap_encode_png(pixmap: &mut tiny_skia::Pixmap) -> Vec<u8> {
    pixmap.encode_png().unwrap()
}