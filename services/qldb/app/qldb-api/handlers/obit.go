package handlers

import (
	"context"
	"net/http"

	"github.com/obada-protocol/server-gateway/services/qldb/business/data/obit"
	"github.com/obada-protocol/server-gateway/services/qldb/foundation/web"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/trace"
)

type obitGroup struct {
	obit obit.Service
}

func(og obitGroup) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.obit.create")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var no obit.NewObit
	if err := web.Decode(r, &no); err != nil {
		return errors.Wrap(err, "decoding new obit")
	}

	err := og.obit.Create(ctx, v.TraceID, no)
	if err != nil {
		return errors.Wrapf(err, "creating new obit: %+v", no)
	}

	return web.Respond(ctx, w, "", http.StatusCreated)
}

func(og obitGroup) metadata(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.obit.metadata")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	m, err := og.obit.Metadata(ctx, v.TraceID, params["obitDID"])
	if err != nil {
		switch err {
		case obit.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case obit.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["obitDID"])
		}
	}

	return web.Respond(ctx, w, m, http.StatusOK)
}

func(og obitGroup) show(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.obit.show")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	o, err := og.obit.FindById(ctx, v.TraceID, params["obitDID"])
	if err != nil {
		switch err {
		case obit.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case obit.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["obitDID"])
		}
	}

	return web.Respond(ctx, w, o, http.StatusOK)
}

func(og obitGroup) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.obit.update")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	_, err := og.obit.FindById(ctx, v.TraceID, params["obitDID"])
	if err != nil {
		switch err {
		case obit.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case obit.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	var o obit.Obit
	if err := web.Decode(r, &o); err != nil {
		return errors.Wrap(err, "decoding obit")
	}

	err = og.obit.Update(ctx, v.TraceID, params["obitDID"], o)

	if err != nil {
		return errors.Wrapf(err, "creating new obit: %+v", o)
	}

	return web.Respond(ctx, w, "", http.StatusNoContent)
}

func(og obitGroup) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.obit.delete")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	params := web.Params(r)
	_, err := og.obit.FindById(ctx, v.TraceID, params["obitDID"])
	if err != nil {
		switch err {
		case obit.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case obit.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	err = og.obit.Delete(ctx, v.TraceID, params["obitDID"])

	if err != nil {
		return errors.Wrapf(err, "deleting obit: %+v", params["obitDID"])
	}

	return web.Respond(ctx, w, "", http.StatusNoContent)
}

func(og obitGroup) search(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.obit.search")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	obits, err := og.obit.FindBy(ctx, v.TraceID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, obits, http.StatusCreated)
}
