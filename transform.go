package resvg

import (
	"context"

	"github.com/kanrichan/resvg-go/internal"
	"github.com/tetratelabs/wazero/api"
)

type transform func(context.Context, api.Module) (int32, error)

func TransformIdentity() transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformIdentity(ctx, mod)
	}
}

func TransformFromRow(sx float32, ky float32, kx float32, sy float32, tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRow(ctx, mod, sx, ky, kx, sy, tx, ty)
	}
}

func TransformFromTranslate(tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromTranslate(ctx, mod, tx, ty)
	}
}

func TransformFromScale(width float32, height float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromScale(ctx, mod, width, height)
	}
}

func TransformFromSkew(kx float32, ky float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromSkew(ctx, mod, kx, ky)
	}
}

func TransformFromRotate(angle float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRotate(ctx, mod, angle)
	}
}

func TransformFromRotateAt(angle float32, tx float32, ty float32) transform {
	return func(ctx context.Context, mod api.Module) (int32, error) {
		return internal.TinySkiaTransformFromRotateAt(ctx, mod, angle, tx, ty)
	}
}
