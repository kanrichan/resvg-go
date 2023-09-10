use resvg::{usvg, tiny_skia};
use usvg::{fontdb, TreeTextToPath, TreeParsing};

pub struct Renderer {
    fontdb_database: fontdb::Database,
    options_resources_dir: Option<std::path::PathBuf>,
    options_dpi: Option<f32>,
    options_font_family: Option<String>,
    options_font_size: Option<f32>,
    options_languages: Option<Vec<String>>,
    options_default_size: Option<usvg::Size>,
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_new")]
#[no_mangle]
pub extern "C" fn renderer_new() -> *mut Renderer {
    let ctx = Renderer { 
        fontdb_database: fontdb::Database::new(),
        options_resources_dir: None,
        options_dpi: None,
        options_font_family: None,
        options_font_size: None,
        options_languages: None,
        options_default_size: None,
    };
    Box::into_raw(ctx.into())
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_delete")]
#[no_mangle]
pub extern "C" fn renderer_delete(ctx: *mut Renderer) {
    let _ = unsafe { Box::from_raw(ctx) };
}

impl Renderer {
    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_fontdb_load_font_data")]
    #[no_mangle]
    pub extern "C" fn fontdb_load_font_data(&mut self, data_ptr: *mut u8, data_size: usize) {
        let data = unsafe { Vec::from_raw_parts(data_ptr, data_size, data_size) };
        self.fontdb_database.load_font_data(data);
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_fontdb_load_font_file")]
    #[no_mangle]
    pub extern "C" fn fontdb_load_font_file(&mut self, file_ptr: *mut u8, file_size: usize) {
        let file = unsafe { String::from_raw_parts(file_ptr, file_size, file_size) };
        let path = std::path::Path::new(&file);
        self.fontdb_database.load_font_file(path).unwrap();
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_fontdb_load_fonts_dir")]
    #[no_mangle]
    pub extern "C" fn fontdb_load_fonts_dir(&mut self, dir_ptr: *mut u8, dir_size: usize) {
        let dir = unsafe { String::from_raw_parts(dir_ptr, dir_size, dir_size) };
        let path = std::path::Path::new(&dir);
        self.fontdb_database.load_fonts_dir(path);
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_options_resources_dir")]
    #[no_mangle]
    pub extern "C" fn options_resources_dir(&mut self, dir_ptr: *mut u8, dir_size: usize) {
        let dir = unsafe { String::from_raw_parts(dir_ptr, dir_size, dir_size) };
        self.options_resources_dir = Some(std::path::PathBuf::from(dir));
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_options_dpi")]
    #[no_mangle]
    pub extern "C" fn options_dpi(&mut self, dpi: f32) {
        self.options_dpi = Some(dpi);
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_options_font_family")]
    #[no_mangle]
    pub extern "C" fn options_font_family(&mut self, font_family_ptr: *mut u8, font_family_size: usize) {
        let font_family = unsafe { String::from_raw_parts(font_family_ptr, font_family_size, font_family_size) };
        self.options_font_family = Some(font_family);
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_options_font_size")]
    #[no_mangle]
    pub extern "C" fn options_font_size(&mut self, font_size: f32) {
        self.options_font_size = Some(font_size);
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_options_languages")]
    #[no_mangle]
    pub extern "C" fn options_languages(&mut self, languages_ptr: *mut u8, languages_size: usize) {
        let languages = unsafe { String::from_raw_parts(languages_ptr, languages_size, languages_size) };
        let mut arr: Vec<String> = Vec::new();
        for token in languages.split_whitespace(){
            arr.push(token.to_owned());
        }
        self.options_languages = Some(arr);
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_options_default_size")]
    #[no_mangle]
    pub extern "C" fn options_default_size(&mut self, width: f32, height: f32) {
        self.options_default_size = Some(usvg::Size::from_wh(width, height).unwrap());
    }

    #[cfg_attr(all(target_arch = "wasm32"), export_name = "__renderer_render")]
    #[no_mangle]
    pub extern "C" fn render(
        &mut self,
        svg_xml_ptr: *mut u8,
        svg_xml_size: usize,
        width: u32,
        height: u32,
    ) -> u64 {
        let svg_xml = unsafe { Vec::from_raw_parts(svg_xml_ptr, svg_xml_size, svg_xml_size) };
        let opts = usvg::Options {
            resources_dir: self.options_resources_dir.clone(),
            dpi: self.options_dpi.unwrap_or(96.0),
            font_family: self.options_font_family.clone().unwrap_or("Times New Roman".to_owned()),
            font_size: self.options_font_size.unwrap_or(12.0),
            languages: self.options_languages.clone().unwrap_or(vec!["en".to_string()]),
            shape_rendering: usvg::ShapeRendering::default(),
            text_rendering: usvg::TextRendering::default(),
            image_rendering: usvg::ImageRendering::default(),
            default_size: self.options_default_size.unwrap_or(usvg::Size::from_wh(100.0, 100.0).unwrap()),
            image_href_resolver: usvg::ImageHrefResolver::default(),
        };
        let path = opts.get_abs_path(std::path::Path::new("ferris_image"));
        println!("{:#?}", path);

        let mut tree = usvg::Tree::from_data(&svg_xml, &opts).unwrap();
        tree.convert_text(&self.fontdb_database);
        let rtree = resvg::Tree::from_usvg(&tree);

        let mut pixmap = tiny_skia::Pixmap::new(
            if width != 0 { width } else { tree.size.width() as u32},
            if height != 0 { height } else { tree.size.height() as u32},
        ).unwrap();

        let transform = tiny_skia::Transform::from_scale(
            pixmap.width() as f32 / rtree.size.width() as f32,
            pixmap.height() as f32 / rtree.size.height() as f32,
        );

        rtree.render(
            transform,
            &mut pixmap.as_mut(),
        );
        let mut ret = pixmap.encode_png().unwrap();
        let ptr = ret.as_mut_ptr();
        let size = ret.len();
        std::mem::forget(ret);

        ((ptr as u64) << 32) | (size as u64)
    }
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__wasm_bytes_malloc")]
#[no_mangle]
pub extern "C" fn wasm_bytes_malloc(size: usize) -> *mut u8 {
    let mut buf = Vec::with_capacity(size);
    let ptr = buf.as_mut_ptr();
    std::mem::forget(buf);
    ptr
}

#[cfg_attr(all(target_arch = "wasm32"), export_name = "__wasm_bytes_free")]
#[no_mangle]
pub extern "C" fn wasm_bytes_free(data_ptr: *mut u8, data_size: usize) {
    let _ = unsafe { Vec::from_raw_parts(data_ptr, 0, data_size) };
}
