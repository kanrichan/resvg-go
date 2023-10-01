package resvg

type ShapeRenderingMode int32
type TextRenderingMode int32
type ImageRenderingMode int32

const (
	ShapeRenderingModeOptimizeSpeed ShapeRenderingMode = iota
	ShapeRenderingModeCrispEdges
	ShapeRenderingModeGeometricPrecision
)

const (
	TextRenderingModeOptimizeSpeed TextRenderingMode = iota
	TextRenderingModeOptimizeLegibility
	TextRenderingModeGeometricPrecision
)

const (
	ImageRenderingModeOptimizeQuality ImageRenderingMode = iota
	ImageRenderingModeOptimizeSpeed
)

type Options struct {
	ResourcesDir       string
	Dpi                float32
	FontFamily         string
	FontSize           float32
	Languages          []string
	ShapeRenderingMode ShapeRenderingMode
	TextRenderingMode  TextRenderingMode
	ImageRenderingMode ImageRenderingMode
	DefaultSizeWidth   float32
	DefaultSizeHeight  float32
}
