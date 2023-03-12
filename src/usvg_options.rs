
#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_default")]
#[no_mangle]
pub fn usvg_options_default() -> *mut usvg::Options {
    let opt = usvg::Options::default();
    let ptr = Box::into_raw(opt.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_free")]
#[no_mangle]
pub fn usvg_options_free(opt: *mut usvg::Options) {
    let _ = unsafe { Box::from_raw(opt) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_set_dpi")]
#[no_mangle]
pub fn usvg_options_set_dpi(opt: &mut usvg::Options, dpi: f64) {
    opt.dpi = dpi;
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_set_font_family")]
#[no_mangle]
pub fn usvg_options_set_font_family(opt: &mut usvg::Options, font_family: (*mut u8, usize)) {
    let font_family = unsafe{ String::from_raw_parts(
        font_family.0, font_family.1, font_family.1) };
    opt.font_family = font_family;
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_set_font_size")]
#[no_mangle]
pub fn usvg_options_set_font_size(opt: &mut usvg::Options, font_size: f64) {
    opt.font_size = font_size;
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_set_keep_named_groups")]
#[no_mangle]
pub fn usvg_options_set_keep_named_groups(opt: &mut usvg::Options, keep_named_groups: bool) {
    opt.keep_named_groups = keep_named_groups;
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_options_set_default_size")]
#[no_mangle]
pub fn usvg_options_set_default_size(opt: &mut usvg::Options, size: *mut usvg::Size) {
    let size = unsafe{ Box::from_raw(size) };
    opt.default_size = *size;
}