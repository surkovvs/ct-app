package vector

import (
	"context"

	"github.com/surkovvs/ct-app/wgchan"
)

type basicFunc[T any] func(context.Context) T

type Vector[T any] struct {
	wait    <-chan struct{}
	release chan<- struct{}
	job     basicFunc[T]
	jobCtx  context.Context
	seq     *Vector[T]
	conc    []*Vector[T]

	report chan<- T

	// job   func(context.Context) error
	// log    appifaces.Logger
}

func (v *Vector[T]) Exec(ctx context.Context) {
	// высвобождаем блокировку, если можем
	defer func() {
		if v.release != nil {
			close(v.release)
		}
	}()

	if v.wait != nil {
		<-v.wait // блокируемся, если нужно
	}

	// log.Printf("executed: %p, %+v\n", v, *v) // TODO: remove

	// проверяем на закрытие контекста исполнения вектора
	select {
	case <-ctx.Done():
		return
	default:
	}

	// исполняем свою функцию
	jobDone := make(chan struct{})
	if v.job != nil {
		go func() {
			defer close(jobDone)
			v.report <- v.job(v.jobCtx)
		}()
	} else {
		close(jobDone)
	}

	// ожидаем завершения джобы, либо закрытия контекста исполнения
	select {
	case <-ctx.Done():
		return
	case <-jobDone:
	}

	// конкурентно запускаем векторы, и ждем исполнения
	if len(v.conc) != 0 {
		wgc := wgchan.NewWgChan()
		for _, cVector := range v.conc {
			wgc.Add(1)
			go func() {
				defer wgc.Done()
				cVector.Exec(ctx)
			}()
		}
		select {
		case <-ctx.Done():
			return
		case <-wgc.WaitChan():
		}
	}

	// // высвобождаем блокировку, если можем
	// defer func() {
	// 	if v.release != nil {
	// 		close(v.release)
	// 	}
	// }()

	// запускаем следующего
	if v.seq != nil {
		// go v.seq.Exec(ctx)
		v.seq.Exec(ctx)
	}
}

type Constructor[T any] struct {
	// log    appifaces.Logger
	report chan<- T
}

func NewConstructor[T any](report chan<- T) Constructor[T] {
	return Constructor[T]{
		report: report,
	}
}

func (c Constructor[T]) newVector() *Vector[T] {
	return &Vector[T]{
		report: c.report,
	}
}

func (c Constructor[T]) Sequentially(ctx context.Context, targets ...any) *Vector[T] {
	var tail, head *Vector[T]

	for i := len(targets) - 1; i >= 0; i-- {
		switch target := targets[i].(type) {
		case *Vector[T]:
			if target == nil {
				continue
			}
			if target.job != nil && target.jobCtx == nil {
				target.jobCtx = ctx
			}
			tail = target
		case func(context.Context) T:
			tail = c.newVector()
			tail.job = target
			tail.jobCtx = ctx
		}

		if tail.seq == nil {
			tail.seq = head
		} else {
			b := tail
			for b.seq != nil {
				b = b.seq
			}
			b.seq = head
		}

		head = tail
	}
	return tail
}

func (c Constructor[T]) Concurrently(ctx context.Context, heads ...any) *Vector[T] {
	res := c.newVector()
	for _, head := range heads {
		switch head := head.(type) {
		case *Vector[T]:
			if head.job != nil && head.jobCtx == nil {
				head.jobCtx = ctx
			}
			res.conc = append(res.conc, head)
		case func(context.Context) T:
			sub := c.newVector()
			sub.jobCtx = ctx
			sub.job = head
			res.conc = append(res.conc, sub)
		}
	}
	return res
}

func (c Constructor[T]) WithRelease(release chan<- struct{}, head any) *Vector[T] {
	var vector *Vector[T]
	switch head := head.(type) {
	case *Vector[T]:
		vector = head
	case func(context.Context) T:
		vector = c.newVector()
		vector.job = head
	default:
		panic("incorrect type passed to *vector.Consructor method")
	}
	vector.release = release

	return vector
}

func (c Constructor[T]) WithWait(wait <-chan struct{}, head any) *Vector[T] {
	var vector *Vector[T]
	switch head := head.(type) {
	case *Vector[T]:
		vector = head
	case func(context.Context) T:
		vector = c.newVector()
		vector.job = head
	default:
		panic("incorrect type passed to *vector.Consructor method")
	}
	vector.wait = wait

	return vector
}
