
#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_default")]
#[no_mangle]
pub fn tiny_skia_transform_default() -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::default();
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_identity")]
#[no_mangle]
pub fn tiny_skia_transform_identity() -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::identity();
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_from_row")]
#[no_mangle]
pub fn tiny_skia_transform_from_row(sx: f32, ky: f32, kx: f32, sy: f32, tx: f32, ty: f32) -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::from_row(sx, ky, kx, sy, tx, ty);
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_from_translate")]
#[no_mangle]
pub fn tiny_skia_transform_from_translate(tx: f32, ty: f32) -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::from_translate(tx, ty);
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_from_scale")]
#[no_mangle]
pub fn tiny_skia_transform_from_scale(sx: f32, sy: f32) -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::from_scale(sx, sy);
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_from_skew")]
#[no_mangle]
pub fn tiny_skia_transform_from_skew(kx: f32, ky: f32) -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::from_skew(kx, ky);
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_from_rotate")]
#[no_mangle]
pub fn tiny_skia_transform_from_rotate(angle: f32) -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::from_rotate(angle);
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_from_rotate_at")]
#[no_mangle]
pub fn tiny_skia_transform_from_rotate_at(angle: f32, tx: f32, ty: f32) -> *mut tiny_skia::Transform {
    let tf = tiny_skia::Transform::from_rotate_at(angle, tx, ty);
    Box::into_raw(tf.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__tiny_skia_transform_free")]
#[no_mangle]
pub fn tiny_skia_transform_free(tf: *mut tiny_skia::Transform) {
    let _ = unsafe { Box::from_raw(tf) };
}