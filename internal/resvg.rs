use resvg::{usvg, tiny_skia};
use usvg::{fontdb, TreeTextToPath, TreeParsing};

#[no_mangle]
pub extern "C" fn fontdb_database_default() -> *mut fontdb::Database {
    let database = fontdb::Database::default();
    Box::into_raw(database.into())
}

#[no_mangle]
pub extern "C" fn fontdb_database_delete(database: *mut fontdb::Database) {
    let _ = unsafe { Box::from_raw(database) };
}

#[no_mangle]
pub extern "C" fn fontdb_database_load_font_data(database: &mut fontdb::Database, data_ptr: *mut u8, data_size: usize) {
    let data = unsafe { Vec::from_raw_parts(data_ptr, data_size, data_size) };
    database.load_font_data(data);
}

#[no_mangle]
pub extern "C" fn fontdb_database_load_font_file(database: &mut fontdb::Database, file_ptr: *mut u8, file_size: usize) {
    let file = unsafe { Vec::from_raw_parts(file_ptr, file_size, file_size) };
    let file = match String::from_utf8(file) {
        Ok(v) => v,
        Err(_) => return,
    };
    let path = std::path::Path::new(&file);
    let _ = database.load_font_file(path);
}

#[no_mangle]
pub extern "C" fn fontdb_database_load_fonts_dir(database: &mut fontdb::Database, dir_ptr: *mut u8, dir_size: usize) {
    let dir = unsafe { Vec::from_raw_parts(dir_ptr, dir_size, dir_size) };
    let dir: String = match String::from_utf8(dir) {
        Ok(v) => v,
        Err(_) => return,
    };
    let path = std::path::Path::new(&dir);
    database.load_fonts_dir(path);
}

#[no_mangle]
pub extern "C" fn fontdb_database_len(database: &mut fontdb::Database) -> usize {
    database.len()
}

#[no_mangle]
pub extern "C" fn fontdb_database_set_serif_family(database: &mut fontdb::Database, family_ptr: *mut u8, family_size: usize) {
    let family = unsafe { Vec::from_raw_parts(family_ptr, family_size, family_size) };
    let family: String = match String::from_utf8(family) {
        Ok(v) => v,
        Err(_) => return,
    };
    database.set_serif_family(family);
}

#[no_mangle]
pub extern "C" fn fontdb_database_set_sans_serif_family(database: &mut fontdb::Database, family_ptr: *mut u8, family_size: usize) {
    let family = unsafe { Vec::from_raw_parts(family_ptr, family_size, family_size) };
    let family: String = match String::from_utf8(family) {
        Ok(v) => v,
        Err(_) => return,
    };
    database.set_sans_serif_family(family);
}

#[no_mangle]
pub extern "C" fn fontdb_database_set_cursive_family(database: &mut fontdb::Database, family_ptr: *mut u8, family_size: usize) {
    let family = unsafe { Vec::from_raw_parts(family_ptr, family_size, family_size) };
    let family: String = match String::from_utf8(family) {
        Ok(v) => v,
        Err(_) => return,
    };
    database.set_cursive_family(family);
}

#[no_mangle]
pub extern "C" fn fontdb_database_set_fantasy_family(database: &mut fontdb::Database, family_ptr: *mut u8, family_size: usize) {
    let family = unsafe { Vec::from_raw_parts(family_ptr, family_size, family_size) };
    let family: String = match String::from_utf8(family) {
        Ok(v) => v,
        Err(_) => return,
    };
    database.set_fantasy_family(family);
}

#[no_mangle]
pub extern "C" fn fontdb_database_set_monospace_family(database: &mut fontdb::Database, family_ptr: *mut u8, family_size: usize) {
    let family = unsafe { Vec::from_raw_parts(family_ptr, family_size, family_size) };
    let family: String = match String::from_utf8(family) {
        Ok(v) => v,
        Err(_) => return,
    };
    database.set_monospace_family(family);
}

#[no_mangle]
pub extern "C" fn usvg_options_default() -> *mut usvg::Options {
    let options = usvg::Options::default();
    Box::into_raw(options.into())
}

#[no_mangle]
pub extern "C" fn usvg_options_delete(options: *mut usvg::Options) {
    let _ = unsafe { Box::from_raw(options) };
}

#[no_mangle]
pub extern "C" fn usvg_options_set_resources_dir(options: &mut usvg::Options, dir_ptr: *mut u8, dir_size: usize) {
    let dir = unsafe { Vec::from_raw_parts(dir_ptr, dir_size, dir_size) };
    let dir: String = match String::from_utf8(dir) {
        Ok(v) => v,
        Err(_) => return,
    };
    options.resources_dir = Some(std::path::PathBuf::from(dir));
}

#[no_mangle]
pub extern "C" fn usvg_options_set_dpi(options: &mut usvg::Options, dpi: f32) {
    options.dpi = dpi;
}

#[no_mangle]
pub extern "C" fn usvg_options_set_font_family(options: &mut usvg::Options, font_family_ptr: *mut u8, font_family_size: usize) {
    let font_family = unsafe { Vec::from_raw_parts(font_family_ptr, font_family_size, font_family_size) };
    let font_family: String = match String::from_utf8(font_family) {
        Ok(v) => v,
        Err(_) => return,
    };
    options.font_family = font_family;
}

#[no_mangle]
pub extern "C" fn usvg_options_set_font_size(options: &mut usvg::Options, font_size: f32) {
    options.font_size = font_size;
}

#[no_mangle]
pub extern "C" fn usvg_options_set_languages(options: &mut usvg::Options, languages_ptr: *mut u8, languages_size: usize) {
    let languages = unsafe { Vec::from_raw_parts(languages_ptr, languages_size, languages_size) };
    let languages: String = match String::from_utf8(languages) {
        Ok(v) => v,
        Err(_) => return,
    };
    let mut arr: Vec<String> = Vec::new();
    for token in languages.split_whitespace(){
        arr.push(token.to_owned());
    }
    options.languages = arr;
}

#[no_mangle]
pub extern "C" fn usvg_options_set_shape_rendering_mode(options: &mut usvg::Options, mode: i32) {
    options.shape_rendering = match mode {
        0 => usvg::ShapeRendering::OptimizeSpeed,
        1 => usvg::ShapeRendering::CrispEdges,
        2 => usvg::ShapeRendering::GeometricPrecision,
        _ => return,
    }
}

#[no_mangle]
pub extern "C" fn usvg_options_set_text_rendering_mode(options: &mut usvg::Options, mode: i32) {
    options.text_rendering = match mode as i32 {
        0 => usvg::TextRendering::OptimizeSpeed,
        1 => usvg::TextRendering::OptimizeLegibility,
        2 => usvg::TextRendering::GeometricPrecision,
        _ => return,
    }
}

#[no_mangle]
pub extern "C" fn usvg_options_set_image_rendering_mode(options: &mut usvg::Options, mode: i32) {
    options.image_rendering = match mode as i32 {
        0 => usvg::ImageRendering::OptimizeQuality,
        1 => usvg::ImageRendering::OptimizeSpeed,
        _ => return,
    }
}

#[no_mangle]
pub extern "C" fn usvg_options_set_default_size(options: &mut usvg::Options, width: f32, height: f32) {
    let size = match usvg::Size::from_wh(width, height) {
        Some(v) => v,
        None => return,
    };
    options.default_size = size;
}

#[no_mangle]
pub extern "C" fn tiny_skia_pixmap_new(width: u32, height: u32) -> *mut tiny_skia::Pixmap {
    let pixmap = match tiny_skia::Pixmap::new(width, height) {
        Some(v) => v,
        None => return std::ptr::null_mut(),
    };
    Box::into_raw(pixmap.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_pixmap_delete(pixmap: *mut tiny_skia::Pixmap) {
    let _ = unsafe { Box::from_raw(pixmap) };
}

#[no_mangle]
pub extern "C" fn tiny_skia_pixmap_encode_png(pixmap: &mut tiny_skia::Pixmap) -> u64 {
    let mut data = match pixmap.encode_png() {
        Ok(v) => v,
        Err(_) => return 0,
    };
    let ptr = data.as_mut_ptr();
    let size = data.len();
    std::mem::forget(data);
    ((ptr as u64) << 32) | (size as u64)
}

#[no_mangle]
pub extern "C" fn tiny_skia_pixmap_get_width(pixmap: &mut tiny_skia::Pixmap) -> u32 {
    pixmap.width()
}

#[no_mangle]
pub extern "C" fn tiny_skia_pixmap_get_height(pixmap: &mut tiny_skia::Pixmap) -> u32 {
    pixmap.height()
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_identity() -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::identity();
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_from_row(sx: f32, ky: f32, kx: f32, sy: f32, tx: f32, ty: f32) -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::from_row(sx, ky, kx, sy, tx, ty);
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_from_translate(tx: f32, ty: f32) -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::from_translate(tx, ty);
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_from_scale(width: f32, height: f32) -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::from_scale(width, height);
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_from_skew(kx: f32, ky: f32) -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::from_skew(kx, ky);
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_from_rotate(angle: f32) -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::from_rotate(angle);
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_from_rotate_at(angle: f32, tx: f32, ty: f32) -> *mut tiny_skia::Transform {
    let transform = tiny_skia::Transform::from_rotate_at(angle, tx, ty);
    Box::into_raw(transform.into())
}

#[no_mangle]
pub extern "C" fn tiny_skia_transform_delete(transform: *mut tiny_skia::Transform) {
    let _ = unsafe { Box::from_raw(transform) };
}

#[no_mangle]
pub extern "C" fn usvg_tree_from_data(data_ptr: *mut u8, data_size: usize, options: &mut usvg::Options) -> *mut usvg::Tree {
    let data = unsafe { Vec::from_raw_parts(data_ptr, data_size, data_size) };
    let tree = match usvg::Tree::from_data(&data, options) {
        Ok(v) => v,
        Err(_) => return std::ptr::null_mut(),
    };
    Box::into_raw(tree.into())
}

#[no_mangle]
pub extern "C" fn usvg_tree_delete(tree: *mut usvg::Tree) {
    let _ = unsafe { Box::from_raw(tree) };
}

#[no_mangle]
pub extern "C" fn usvg_tree_convert_text(tree: &mut usvg::Tree, database: &mut fontdb::Database) {
    tree.convert_text(database);
}

#[no_mangle]
pub extern "C" fn usvg_tree_get_size_width(tree: &mut usvg::Tree) -> f32 {
    tree.size.width()
}

#[no_mangle]
pub extern "C" fn usvg_tree_get_size_height(tree: &mut usvg::Tree) -> f32 {
    tree.size.height()
}

#[no_mangle]
pub extern "C" fn resvg_tree_from_usvg(tree: &usvg::Tree) -> *mut resvg::Tree {
    let rtree = resvg::Tree::from_usvg(tree);
    Box::into_raw(rtree.into())
}

#[no_mangle]
pub extern "C" fn resvg_tree_delete(rtree: *mut resvg::Tree) {
    let _ = unsafe { Box::from_raw(rtree) };
}

#[no_mangle]
pub extern "C" fn resvg_tree_render(rtree: &mut resvg::Tree, transform: &mut tiny_skia::Transform, pixmap: &mut tiny_skia::Pixmap) {
    rtree.render(
        *transform,
        &mut pixmap.as_mut(),
    );
}

#[no_mangle]
pub extern "C" fn memory_malloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    std::mem::forget(buf);
    ptr
}

#[no_mangle]
pub extern "C" fn memory_free(data_ptr: *mut u8, data_size: usize) {
    let _ = unsafe { Vec::from_raw_parts(data_ptr, 0, data_size) };
}