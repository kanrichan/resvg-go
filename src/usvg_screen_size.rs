
#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_screen_size_new")]
#[no_mangle]
pub fn usvg_screen_size_new(width: u32, height: u32) -> *mut usvg::ScreenSize {
    let screen_size = usvg::ScreenSize::new(width, height);
    let ptr = Box::into_raw(screen_size.unwrap().into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_screen_size_free")]
#[no_mangle]
pub fn usvg_screen_size_free(screen_size: *mut usvg::ScreenSize) {
    let _ = unsafe { Box::from_raw(screen_size) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_screen_size_width")]
#[no_mangle]
pub fn usvg_screen_size_width(screen_size: &mut usvg::ScreenSize) -> u32 {
    screen_size.width()
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_screen_size_height")]
#[no_mangle]
pub fn usvg_screen_size_height(screen_size: &mut usvg::ScreenSize) -> u32 {
    screen_size.height()
}