mod fontdb_database;
mod rust_bytes;
mod tiny_skia_transform;
mod usvg_fit_to;
mod usvg_options;
mod usvg_pixmap;
mod usvg_screen_size;
mod usvg_size;
mod usvg_tree;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__resvg_render")]
#[no_mangle]
pub fn resvg_render(tree: &mut usvg::Tree, ft: *mut usvg::FitTo, tf: *mut tiny_skia::Transform, pixmap: &mut tiny_skia::Pixmap) {
    let ft = unsafe { Box::from_raw(ft) };
    let tf = unsafe { Box::from_raw(tf) };
    resvg::render(tree, *ft, *tf, pixmap.as_mut());
}