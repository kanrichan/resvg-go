use usvg_text_layout::TreeTextToPath;
use usvg_text_layout::fontdb;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_tree_from_data")]
#[no_mangle]
pub fn usvg_tree_from_data(data: (*mut u8, usize), opt: &mut usvg::Options) -> *mut usvg::Tree {
    let data = unsafe { Vec::from_raw_parts(data.0, data.1, data.1) };
    let tree = usvg::Tree::from_data(&data, opt).unwrap();
    let ptr = Box::into_raw(tree.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_tree_free")]
#[no_mangle]
pub fn usvg_tree_free(tree: *mut usvg::Tree) {
    let _ = unsafe { Box::from_raw(tree) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_tree_convert_text")]
#[no_mangle]
pub fn usvg_tree_convert_text(tree: &mut usvg::Tree, db: &mut fontdb::Database, keep_named_groups: bool) {
    tree.convert_text(db, keep_named_groups);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__usvg_tree_get_size_clone")]
#[no_mangle]
pub fn usvg_tree_get_size_clone(tree: &mut usvg::Tree) -> *mut usvg::Size {
    let ptr = Box::into_raw(tree.size.clone().into());
    ptr
}