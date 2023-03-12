
#[cfg_attr(all(target_arch = "wasm32"), export_name = "__rust_bytes_new")]
#[no_mangle]
pub fn rust_bytes_new(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    std::mem::forget(buf);
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__rust_bytes_free")]
#[no_mangle]
pub fn rust_bytes_free(data: (*mut u8, usize)) {
    let _ = unsafe { Vec::from_raw_parts(data.0, 0, data.1) };
}
