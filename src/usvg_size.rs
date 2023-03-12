#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_size_new")]
#[no_mangle]
pub fn usvg_size_new(width: f64, height: f64) -> *mut usvg::Size {
    let size = usvg::Size::new(width, height).unwrap();
    let ptr = Box::into_raw(size.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_size_free")]
#[no_mangle]
pub fn usvg_size_free(size: *mut usvg::Size) {
    let _ = unsafe { Box::from_raw(size) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_size_width")]
#[no_mangle]
pub fn usvg_size_width(size: &usvg::Size) -> f64 {
    size.width()
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_size_height")]
#[no_mangle]
pub fn usvg_size_height(size: &usvg::Size) -> f64 {
    size.height()
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_size_to_screen_size")]
#[no_mangle]
pub fn usvg_size_to_screen_size(size: &usvg::Size) -> *mut usvg::ScreenSize {
    let screen_size = size.to_screen_size();
    let ptr = Box::into_raw(screen_size.into());
    ptr
}

