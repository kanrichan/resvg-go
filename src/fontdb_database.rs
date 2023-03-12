use usvg_text_layout::fontdb;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_new")]
#[no_mangle]
pub fn fontdb_database_new() -> *mut fontdb::Database {
    let db = fontdb::Database::new();
    let ptr = Box::into_raw(db.into());
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_free")]
#[no_mangle]
pub fn database_free(db: *mut fontdb::Database) {
    let _ = unsafe { Box::from_raw(db) };
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_load_font_data")]
#[no_mangle]
pub fn fontdb_database_load_font_data(db: &mut fontdb::Database, data: (*mut u8, usize)) {
    let data = unsafe{ Vec::from_raw_parts(
        data.0, data.1, data.1) };
    db.load_font_data(data);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_len")]
#[no_mangle]
pub fn fontdb_database_len(db: &fontdb::Database) -> usize {
    db.len()
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_set_serif_family")]
#[no_mangle]
pub fn fontdb_database_set_serif_family(db: &mut fontdb::Database, serif_family: (*mut u8, usize)) {
    let serif_family = unsafe{ String::from_raw_parts(
        serif_family.0, serif_family.1, serif_family.1) };
    db.set_serif_family(serif_family);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_set_sans_serif_family")]
#[no_mangle]
pub fn fontdb_database_set_sans_serif_family(db: &mut fontdb::Database, sans_serif_family: (*mut u8, usize)) {
    let sans_serif_family = unsafe{ String::from_raw_parts(
        sans_serif_family.0, sans_serif_family.1, sans_serif_family.1) };
    db.set_sans_serif_family(sans_serif_family);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_set_cursive_family")]
#[no_mangle]
pub fn fontdb_database_set_cursive_family(db: &mut fontdb::Database, cursive_family: (*mut u8, usize)) {
    let cursive_family = unsafe{ String::from_raw_parts(
        cursive_family.0, cursive_family.1, cursive_family.1) };
    db.set_cursive_family(cursive_family);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_set_fantasy_family")]
#[no_mangle]
pub fn fontdb_database_set_fantasy_family(db: &mut fontdb::Database, fantasy_family: (*mut u8, usize)) {
    let fantasy_family = unsafe{ String::from_raw_parts(
        fantasy_family.0, fantasy_family.1, fantasy_family.1) };
    db.set_fantasy_family(fantasy_family);
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__fontdb_database_set_monospace_family")]
#[no_mangle]
pub fn fontdb_database_set_monospace_family(db: &mut fontdb::Database, monospace_family: (*mut u8, usize)) {
    let monospace_family = unsafe{ String::from_raw_parts(
        monospace_family.0, monospace_family.1, monospace_family.1) };
    db.set_monospace_family(monospace_family);
}