package gblib

import "context"

type UseCase[TIn any, TOut any] interface {
	Execute(ctx context.Context, input TIn) (TOut, error)
}

type handlerFunc func(context.Context, any) (any, error)

type interceptor func(context.Context, any, handlerFunc) (any, error)

type decoratedUseCase[TInput any, TOutput any] struct {
	interceptor interceptor
	uc          UseCase[TInput, TOutput]
}

func (dc decoratedUseCase[TIn, TOut]) Execute(ctx context.Context, input TIn) (TOut, error) {
	res, err := dc.interceptor(ctx, input, useCaseHandler(dc.uc))
	return res.(TOut), err
}

func useCaseHandler[TIn any, TOut any](uc UseCase[TIn, TOut]) handlerFunc {
	return func(ctx context.Context, i any) (any, error) {
		return uc.Execute(ctx, i.(TIn))
	}
}

func New[TIn any, TO any](useCase UseCase[TIn, TO], interceptors ...interceptor) UseCase[TIn, TO] {
	if len(interceptors) == 0 {
		return useCase
	}

	for i := len(interceptors) - 1; i >= 0; i-- {
		useCase = decoratedUseCase[TIn, TO]{
			uc:          useCase,
			interceptor: interceptors[i],
		}
	}

	return useCase
}
