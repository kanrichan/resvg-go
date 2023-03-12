
#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_fit_to_original")]
#[no_mangle]
pub fn usvg_fit_to_original() -> *mut usvg::FitTo {
    let ft = usvg::FitTo::Original;
    Box::into_raw(ft.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_fit_to_width")]
#[no_mangle]
pub fn usvg_fit_to_width(width: u32) -> *mut usvg::FitTo {
    let ft = usvg::FitTo::Width(width);
    Box::into_raw(ft.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_fit_to_height")]
#[no_mangle]
pub fn usvg_fit_to_height(height: u32) -> *mut usvg::FitTo {
    let ft = usvg::FitTo::Height(height);
    Box::into_raw(ft.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_fit_to_size")]
#[no_mangle]
pub fn usvg_fit_to_size(width: u32, height: u32) -> *mut usvg::FitTo {
    let ft = usvg::FitTo::Size(width, height);
    Box::into_raw(ft.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_fit_to_zoom")]
#[no_mangle]
pub fn usvg_fit_to_zoom(zoom: f32) -> *mut usvg::FitTo {
    let ft = usvg::FitTo::Zoom(zoom);
    Box::into_raw(ft.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_fit_to_free")]
#[no_mangle]
pub fn usvg_fit_to_free(ft: *mut usvg::FitTo) {
    let _ = unsafe { Box::from_raw(ft) };
}