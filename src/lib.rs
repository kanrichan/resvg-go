use usvg_text_layout::{fontdb, TreeTextToPath};

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_bytes_allocate")]
#[no_mangle]
pub extern "C" fn wazero_allocate(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    std::mem::forget(buf);
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_bytes_free")]
#[no_mangle]
pub extern "C" fn wazero_free(ptr: *mut u8, size: usize){
    let _ = unsafe {Vec::from_raw_parts(ptr, 0, size)};
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_database_new")]
#[no_mangle]
pub fn database_new() -> *mut fontdb::Database {
    let db = fontdb::Database::new();
    let ptr = Box::into_raw(db.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_database_load_font_data")]
#[no_mangle]
pub fn database_load_font_data(db: *mut fontdb::Database, data_ptr: *mut u8, data_size: usize) {
    let mut box_db = unsafe { Box::from_raw(db) };
    let box_data = unsafe { Vec::from_raw_parts(data_ptr, data_size, data_size) };
    box_db.load_font_data(box_data.clone());
    box_db.set_serif_family(String::from("Source Han Serif CN Light"));
    box_db.set_sans_serif_family(String::from("Source Han Serif CN Light"));
    box_db.set_cursive_family(String::from("Source Han Serif CN Light"));
    box_db.set_fantasy_family(String::from("Source Han Serif CN Light"));
    box_db.set_monospace_family(String::from("Source Han Serif CN Light"));
    std::mem::forget(box_db);
    std::mem::forget(box_data);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_database_len")]
#[no_mangle]
pub fn database_len(db: *mut fontdb::Database) -> usize {
    let box_db = unsafe { Box::from_raw(db) };
    let size = box_db.len();
    std::mem::forget(box_db);
    size
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_database_free")]
#[no_mangle]
pub fn database_free(db: *mut fontdb::Database) {
    let _ = unsafe { Box::from_raw(db) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_size_new")]
#[no_mangle]
pub fn size_new(width: f64, height: f64) -> *mut usvg::Size {
    let size = usvg::Size::new(width, height);
    let ptr = Box::into_raw(size.unwrap().into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_size_width")]
#[no_mangle]
pub fn size_width(size: *mut usvg::Size) -> f64 {
    let box_size = unsafe { Box::from_raw(size) };
    let width = box_size.width();
    std::mem::forget(box_size);
    width
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_size_height")]
#[no_mangle]
pub fn size_height(size: *mut usvg::Size) -> f64 {
    let box_size = unsafe { Box::from_raw(size) };
    let width = box_size.height();
    std::mem::forget(box_size);
    width
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_size_to_screen_size")]
#[no_mangle]
pub fn size_to_screen_size(size: *mut usvg::Size) -> *mut usvg::ScreenSize {
    let box_size = unsafe { Box::from_raw(size) };
    let screen_size = box_size.to_screen_size();
    let ptr = Box::into_raw(screen_size.into());
    std::mem::forget(box_size);
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_size_free")]
#[no_mangle]
pub fn size_free(size: *mut usvg::Size) {
    let _ = unsafe { Box::from_raw(size) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_screen_size_new")]
#[no_mangle]
pub fn screen_size_new(width: u32, height: u32) -> *mut usvg::ScreenSize {
    let screen_size = usvg::ScreenSize::new(width, height);
    let ptr = Box::into_raw(screen_size.unwrap().into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_screen_size_width")]
#[no_mangle]
pub fn screen_size_width(screen_size: *mut usvg::ScreenSize) -> u32 {
    let box_screen_size = unsafe { Box::from_raw(screen_size) };
    let width = box_screen_size.width();
    std::mem::forget(box_screen_size);
    width
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_screen_size_height")]
#[no_mangle]
pub fn screen_size_height(screen_size: *mut usvg::ScreenSize) -> u32 {
    let box_screen_size = unsafe { Box::from_raw(screen_size) };
    let width = box_screen_size.height();
    std::mem::forget(box_screen_size);
    width
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_screen_size_free")]
#[no_mangle]
pub fn screen_size_free(screen_size: *mut usvg::ScreenSize) {
    let _ = unsafe { Box::from_raw(screen_size) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_option_default")]
#[no_mangle]
pub fn option_default() -> *mut usvg::Options {
    let opt = usvg::Options::default();
    let ptr = Box::into_raw(opt.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_option_set_font_family")]
#[no_mangle]
pub fn option_set_font_family(opt: *mut usvg::Options, font_family_ptr: *mut u8, font_family_size: usize) {
    let mut box_opt = unsafe { Box::from_raw(opt) };
    let box_font_family = unsafe { String::from_raw_parts(font_family_ptr, font_family_size, font_family_size) };
    box_opt.font_family = box_font_family.clone();
    box_opt.font_size = 12.0;
    std::mem::forget(box_opt);
    std::mem::forget(box_font_family);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_option_set_default_size")]
#[no_mangle]
pub fn option_set_default_size(opt: *mut usvg::Options, size: *mut usvg::Size) {
    let mut box_opt = unsafe { Box::from_raw(opt) };
    let box_size = unsafe { Box::from_raw(size) };
    box_opt.default_size = *box_size;
    std::mem::forget(box_opt);
    std::mem::forget(box_size);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_option_free")]
#[no_mangle]
pub fn option_free(opt: *mut usvg::Options) {
    let _ = unsafe { Box::from_raw(opt) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_tree_from_data")]
#[no_mangle]
pub fn tree_from_data(data_ptr: *mut u8, data_size: usize, opt: *mut usvg::Options) -> *mut usvg::Tree {
    let box_data = unsafe { Vec::from_raw_parts(data_ptr, data_size, data_size) };
    let box_opt = unsafe { Box::from_raw(opt) };
    let tree = usvg::Tree::from_data(&box_data, &box_opt).unwrap();
    let ptr = Box::into_raw(tree.into());
    std::mem::forget(box_data);
    std::mem::forget(box_opt);
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_tree_convert_text")]
#[no_mangle]
pub fn tree_convert_text(tree: *mut usvg::Tree, db: *mut fontdb::Database, keep_named_groups: bool) {
    let mut box_tree = unsafe { Box::from_raw(tree) };
    let box_db = unsafe { Box::from_raw(db) };
    box_tree.convert_text(&box_db, keep_named_groups);
    std::mem::forget(box_tree);
    std::mem::forget(box_db);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_tree_get_size")]
#[no_mangle]
pub fn tree_get_size(tree: *mut usvg::Tree) -> *mut usvg::Size {
    let box_tree = unsafe { Box::from_raw(tree) };
    let ptr = Box::into_raw(box_tree.size.clone().into());
    std::mem::forget(box_tree);
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_tree_free")]
#[no_mangle]
pub fn tree_free(tree: *mut usvg::Tree) {
    let _ = unsafe { Box::from_raw(tree) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_pixmap_new")]
#[no_mangle]
pub fn pixmap_new(width: u32, height: u32) -> *mut tiny_skia::Pixmap {
    let pixmap = tiny_skia::Pixmap::new(width, height).unwrap();
    let ptr = Box::into_raw(pixmap.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_pixmap_encode_png")]
#[no_mangle]
pub fn pixmap_encode_png(pixmap: *mut tiny_skia::Pixmap) -> (*mut u8, usize) {
    let box_pixmap = unsafe { Box::from_raw(pixmap) };
    let mut data = box_pixmap.encode_png().unwrap();
    let ptr = data.as_mut_ptr();
    let size = data.len();
    std::mem::forget(box_pixmap);
    std::mem::forget(data);
    (ptr, size)
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_pixmap_free")]
#[no_mangle]
pub fn pixmap_free(pixmap: *mut tiny_skia::Pixmap) {
    let _ = unsafe { Box::from_raw(pixmap) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__floattech_render")]
#[no_mangle]
pub fn render(tree: *mut usvg::Tree, pixmap: *mut tiny_skia::Pixmap) {
    let box_tree = unsafe { Box::from_raw(tree) };
    let mut box_pixmap = unsafe { Box::from_raw(pixmap) };
    resvg::render(&box_tree, usvg::FitTo::Zoom(1.0), tiny_skia::Transform::identity(), box_pixmap.as_mut().as_mut());
    std::mem::forget(box_tree);
    std::mem::forget(box_pixmap);
}