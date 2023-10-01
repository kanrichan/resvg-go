package resvg

// ShapeRenderingMode a shape rendering method, `shape-rendering` attribute in the SVG.
type ShapeRenderingMode int32

// TextRenderingMode a text rendering method, `text-rendering` attribute in the SVG.
type TextRenderingMode int32

// ImageRenderingMode an image rendering method, `image-rendering` attribute in the SVG.
type ImageRenderingMode int32

const (
	// ShapeRenderingModeOptimizeSpeed OptimizeSpeed
	ShapeRenderingModeOptimizeSpeed ShapeRenderingMode = iota
	// ShapeRenderingModeCrispEdges CrispEdges
	ShapeRenderingModeCrispEdges
	// ShapeRenderingModeGeometricPrecision GeometricPrecision
	ShapeRenderingModeGeometricPrecision
)

const (
	// TextRenderingModeOptimizeSpeed OptimizeSpeed
	TextRenderingModeOptimizeSpeed TextRenderingMode = iota
	// TextRenderingModeOptimizeLegibility OptimizeLegibility
	TextRenderingModeOptimizeLegibility
	// TextRenderingModeGeometricPrecision GeometricPrecision
	TextRenderingModeGeometricPrecision
)

const (
	// ImageRenderingModeOptimizeQuality OptimizeQuality
	ImageRenderingModeOptimizeQuality ImageRenderingMode = iota
	// ImageRenderingModeOptimizeSpeed OptimizeSpeed
	ImageRenderingModeOptimizeSpeed
)

type Options struct {
	// ResourcesDir directory that will be used during relative paths resolving.
	// Expected to be the same as the directory that contains the SVG file,
	// but can be set to any.
	// Default: `None`
	ResourcesDir string

	// Dpi target DPI.
	// Impacts units conversion.
	// Default: 96.0
	Dpi float32

	// FontFamily a default font family.
	// Will be used when no `font-family` attribute is set in the SVG.
	// Default: Times New Roman
	FontFamily string

	// FontSize a default font size.
	// Will be used when no `font-size` attribute is set in the SVG.
	// Default: 12
	FontSize float32

	// Languages a list of languages.
	// Will be used to resolve a `systemLanguage` conditional attribute.
	// Format: en, en-US.
	// Default: `[en]`
	Languages []string

	// ShapeRenderingMode specifies the default shape rendering method.
	// Will be used when an SVG element's `shape-rendering` property is set to `auto`.
	// Default: GeometricPrecision
	ShapeRenderingMode ShapeRenderingMode

	// TextRenderingMode specifies the default text rendering method.
	// Will be used when an SVG element's `text-rendering` property is set to `auto`.
	// Default: OptimizeLegibility
	TextRenderingMode TextRenderingMode

	// ImageRenderingMode specifies the default image rendering method.
	// Will be used when an SVG element's `image-rendering` property is set to `auto`.
	// Default: OptimizeQuality
	ImageRenderingMode ImageRenderingMode

	// DefaultSizeWidth default viewport size to assume if there is no `viewBox` attribute and
	// the `width` attributes are relative.
	// Default: `100`
	DefaultSizeWidth float32

	// DefaultSizeHeight default viewport size to assume if there is no `viewBox` attribute and
	// the `height` attributes are relative.
	// Default: `100`
	DefaultSizeHeight float32
}
