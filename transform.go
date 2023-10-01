package resvg

import (
	"context"

	"github.com/kanrichan/resvg-go/internal"
	"github.com/tetratelabs/wazero/api"
)

type transform func(context.Context, api.Module) (int32, error)

// TransformIdentity creates an identity transform.
func TransformIdentity() transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformIdentity(ctx, mod)
	}
}

// TransformFromRow creates a new `Transform`.
// We are using column-major-column-vector matrix notation, therefore it's ky-kx, not kx-ky.
func TransformFromRow(sx float32, ky float32, kx float32, sy float32, tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRow(ctx, mod, sx, ky, kx, sy, tx, ty)
	}
}

// TransformFromTranslate creates a new translating `Transform`.
func TransformFromTranslate(tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromTranslate(ctx, mod, tx, ty)
	}
}

// TransformFromScale creates a new scaling `Transform`.
func TransformFromScale(width float32, height float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromScale(ctx, mod, width, height)
	}
}

// TransformFromSkew creates a new skewing `Transform`.
func TransformFromSkew(kx float32, ky float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromSkew(ctx, mod, kx, ky)
	}
}

// TransformFromRotate creates a new rotating `Transform`.
// `angle` in degrees.
func TransformFromRotate(angle float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRotate(ctx, mod, angle)
	}
}

// TransformFromRotateAt creates a new rotating `Transform` at the specified position.
// `angle` in degrees.
func TransformFromRotateAt(angle float32, tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRotateAt(ctx, mod, angle, tx, ty)
	}
}
