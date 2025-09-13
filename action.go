package GAction

import (
	"sync"
)

type Action0 func()
type Action1[T any] func(T)
type Action2[A, B any] func(A, B)
type Action3[A, B, C any] func(A, B, C)
type Action4[A, B, C, D any] func(A, B, C, D)

func Combine0(list ...Action0) Action0 {
	return func() {
		for _, f := range list {
			if f != nil {
				f()
			}
		}
	}
}
func Combine1[T any](list ...Action1[T]) Action1[T] {
	return func(v T) {
		for _, f := range list {
			if f != nil {
				f(v)
			}
		}
	}
}
func Combine2[A, B any](list ...Action2[A, B]) Action2[A, B] {
	return func(a A, b B) {
		for _, f := range list {
			if f != nil {
				f(a, b)
			}
		}
	}
}
func Combine3[A, B, C any](list ...Action3[A, B, C]) Action3[A, B, C] {
	return func(a A, b B, c C) {
		for _, f := range list {
			if f != nil {
				f(a, b, c)
			}
		}
	}
}
func Combine4[A, B, C, D any](list ...Action4[A, B, C, D]) Action4[A, B, C, D] {
	return func(a A, b B, c C, d D) {
		for _, f := range list {
			if f != nil {
				f(a, b, c, d)
			}
		}
	}
}

type Invoker0 interface{ Invoke() }
type Invoker1[T any] interface{ Invoke(T) }
type Invoker2[A, B any] interface{ Invoke(A, B) }
type Invoker3[A, B, C any] interface{ Invoke(A, B, C) }
type Invoker4[A, B, C, D any] interface{ Invoke(A, B, C, D) }

type Delegate0 struct {
	mu sync.RWMutex
	f  Action0
}

func NewDelegate0() *Delegate0     { return &Delegate0{} }
func (d *Delegate0) Set(f Action0) { d.mu.Lock(); d.f = f; d.mu.Unlock() }
func (d *Delegate0) Invoke() {
	d.mu.RLock()
	f := d.f
	d.mu.RUnlock()
	if f != nil {
		f()
	}
}

type Delegate1[T any] struct {
	mu sync.RWMutex
	f  Action1[T]
}

func NewDelegate1[T any]() *Delegate1[T] { return &Delegate1[T]{} }
func (d *Delegate1[T]) Set(f Action1[T]) { d.mu.Lock(); d.f = f; d.mu.Unlock() }
func (d *Delegate1[T]) Invoke(v T) {
	d.mu.RLock()
	f := d.f
	d.mu.RUnlock()
	if f != nil {
		f(v)
	}
}

var _ Invoker1[int] = (*Delegate1[int])(nil)

type Delegate2[A, B any] struct {
	mu sync.RWMutex
	f  Action2[A, B]
}

func NewDelegate2[A, B any]() *Delegate2[A, B] { return &Delegate2[A, B]{} }
func (d *Delegate2[A, B]) Set(f Action2[A, B]) { d.mu.Lock(); d.f = f; d.mu.Unlock() }
func (d *Delegate2[A, B]) Invoke(a A, b B) {
	d.mu.RLock()
	f := d.f
	d.mu.RUnlock()
	if f != nil {
		f(a, b)
	}
}

type Delegate3[A, B, C any] struct {
	mu sync.RWMutex
	f  Action3[A, B, C]
}

func NewDelegate3[A, B, C any]() *Delegate3[A, B, C] { return &Delegate3[A, B, C]{} }
func (d *Delegate3[A, B, C]) Set(f Action3[A, B, C]) { d.mu.Lock(); d.f = f; d.mu.Unlock() }
func (d *Delegate3[A, B, C]) Invoke(a A, b B, c C) {
	d.mu.RLock()
	f := d.f
	d.mu.RUnlock()
	if f != nil {
		f(a, b, c)
	}
}

type Delegate4[A, B, C, D any] struct {
	mu sync.RWMutex
	f  Action4[A, B, C, D]
}

func NewDelegate4[A, B, C, D any]() *Delegate4[A, B, C, D] { return &Delegate4[A, B, C, D]{} }
func (d *Delegate4[A, B, C, D]) Set(f Action4[A, B, C, D]) { d.mu.Lock(); d.f = f; d.mu.Unlock() }
func (d *Delegate4[A, B, C, D]) Invoke(a A, b B, c C, e D) {
	d.mu.RLock()
	f := d.f
	d.mu.RUnlock()
	if f != nil {
		f(a, b, c, e)
	}
}

type Subscription interface{ Unsubscribe() }

/* ------- Event0 ------- */

type Event0 interface {
	Subscribe(Action0) Subscription
	HasSubscribers() bool
}
type Emitter0 interface{ Emit() }

func NewEvent0() (Event0, Emitter0) {
	e := &event0{}
	return &event0Reader{e}, &event0Emitter{e}
}

type event0 struct {
	mu     sync.RWMutex
	nextID uint64
	list   []entry0
}
type entry0 struct {
	id uint64
	f  Action0
}

func (e *event0) subscribe(f Action0) Subscription {
	if f == nil {
		return noopSub{}
	}
	e.mu.Lock()
	id := e.nextID
	e.nextID++
	e.list = append(e.list, entry0{id, f})
	e.mu.Unlock()
	return &sub0{e: e, id: id}
}
func (e *event0) emit() {
	e.mu.RLock()
	snap := make([]Action0, 0, len(e.list))
	for _, it := range e.list {
		if it.f != nil {
			snap = append(snap, it.f)
		}
	}
	e.mu.RUnlock()
	for _, h := range snap {
		h()
	}
}
func (e *event0) has() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, it := range e.list {
		if it.f != nil {
			return true
		}
	}
	return false
}

type sub0 struct {
	once sync.Once
	e    *event0
	id   uint64
}

func (s *sub0) Unsubscribe() {
	s.once.Do(func() {
		s.e.mu.Lock()
		for i := range s.e.list {
			if s.e.list[i].id == s.id {
				s.e.list[i].f = nil
				break
			}
		}
		s.e.mu.Unlock()
	})
}

type noopSub struct{}

func (noopSub) Unsubscribe() {}

type event0Reader struct{ e *event0 }

func (r *event0Reader) Subscribe(f Action0) Subscription { return r.e.subscribe(f) }
func (r *event0Reader) HasSubscribers() bool             { return r.e.has() }

type event0Emitter struct{ e *event0 }

func (w *event0Emitter) Emit() { w.e.emit() }

/* ------- Event1 ------- */

type Event1[T any] interface {
	Subscribe(Action1[T]) Subscription
	HasSubscribers() bool
}
type Emitter1[T any] interface{ Emit(T) }

func NewEvent1[T any]() (Event1[T], Emitter1[T]) {
	e := &event1[T]{}
	return &event1Reader[T]{e}, &event1Emitter[T]{e}
}

type event1[T any] struct {
	mu     sync.RWMutex
	nextID uint64
	list   []entry1[T]
}
type entry1[T any] struct {
	id uint64
	f  Action1[T]
}

func (e *event1[T]) subscribe(f Action1[T]) Subscription {
	if f == nil {
		return noopSub{}
	}
	e.mu.Lock()
	id := e.nextID
	e.nextID++
	e.list = append(e.list, entry1[T]{id, f})
	e.mu.Unlock()
	return &sub1[T]{e: e, id: id}
}
func (e *event1[T]) emit(v T) {
	e.mu.RLock()
	snap := make([]Action1[T], 0, len(e.list))
	for _, it := range e.list {
		if it.f != nil {
			snap = append(snap, it.f)
		}
	}
	e.mu.RUnlock()
	for _, h := range snap {
		h(v)
	}
}
func (e *event1[T]) has() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, it := range e.list {
		if it.f != nil {
			return true
		}
	}
	return false
}

type sub1[T any] struct {
	once sync.Once
	e    *event1[T]
	id   uint64
}

func (s *sub1[T]) Unsubscribe() {
	s.once.Do(func() {
		s.e.mu.Lock()
		for i := range s.e.list {
			if s.e.list[i].id == s.id {
				s.e.list[i].f = nil
				break
			}
		}
		s.e.mu.Unlock()
	})
}

type event1Reader[T any] struct{ e *event1[T] }

func (r *event1Reader[T]) Subscribe(f Action1[T]) Subscription { return r.e.subscribe(f) }
func (r *event1Reader[T]) HasSubscribers() bool                { return r.e.has() }

type event1Emitter[T any] struct{ e *event1[T] }

func (w *event1Emitter[T]) Emit(v T) { w.e.emit(v) }

/* ------- Event2 ------- */

type Event2[A, B any] interface {
	Subscribe(Action2[A, B]) Subscription
	HasSubscribers() bool
}
type Emitter2[A, B any] interface{ Emit(A, B) }

func NewEvent2[A, B any]() (Event2[A, B], Emitter2[A, B]) {
	e := &event2[A, B]{}
	return &event2Reader[A, B]{e}, &event2Emitter[A, B]{e}
}

type event2[A, B any] struct {
	mu     sync.RWMutex
	nextID uint64
	list   []entry2[A, B]
}
type entry2[A, B any] struct {
	id uint64
	f  Action2[A, B]
}

func (e *event2[A, B]) subscribe(f Action2[A, B]) Subscription {
	if f == nil {
		return noopSub{}
	}
	e.mu.Lock()
	id := e.nextID
	e.nextID++
	e.list = append(e.list, entry2[A, B]{id, f})
	e.mu.Unlock()
	return &sub2[A, B]{e: e, id: id}
}
func (e *event2[A, B]) emit(a A, b B) {
	e.mu.RLock()
	snap := make([]Action2[A, B], 0, len(e.list))
	for _, it := range e.list {
		if it.f != nil {
			snap = append(snap, it.f)
		}
	}
	e.mu.RUnlock()
	for _, h := range snap {
		h(a, b)
	}
}
func (e *event2[A, B]) has() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, it := range e.list {
		if it.f != nil {
			return true
		}
	}
	return false
}

type sub2[A, B any] struct {
	once sync.Once
	e    *event2[A, B]
	id   uint64
}

func (s *sub2[A, B]) Unsubscribe() {
	s.once.Do(func() {
		s.e.mu.Lock()
		for i := range s.e.list {
			if s.e.list[i].id == s.id {
				s.e.list[i].f = nil
				break
			}
		}
		s.e.mu.Unlock()
	})
}

type event2Reader[A, B any] struct{ e *event2[A, B] }

func (r *event2Reader[A, B]) Subscribe(f Action2[A, B]) Subscription { return r.e.subscribe(f) }
func (r *event2Reader[A, B]) HasSubscribers() bool                   { return r.e.has() }

type event2Emitter[A, B any] struct{ e *event2[A, B] }

func (w *event2Emitter[A, B]) Emit(a A, b B) { w.e.emit(a, b) }

/* ------- Event3 ------- */

type Event3[A, B, C any] interface {
	Subscribe(Action3[A, B, C]) Subscription
	HasSubscribers() bool
}
type Emitter3[A, B, C any] interface{ Emit(A, B, C) }

func NewEvent3[A, B, C any]() (Event3[A, B, C], Emitter3[A, B, C]) {
	e := &event3[A, B, C]{}
	return &event3Reader[A, B, C]{e}, &event3Emitter[A, B, C]{e}
}

type event3[A, B, C any] struct {
	mu     sync.RWMutex
	nextID uint64
	list   []entry3[A, B, C]
}
type entry3[A, B, C any] struct {
	id uint64
	f  Action3[A, B, C]
}

func (e *event3[A, B, C]) subscribe(f Action3[A, B, C]) Subscription {
	if f == nil {
		return noopSub{}
	}
	e.mu.Lock()
	id := e.nextID
	e.nextID++
	e.list = append(e.list, entry3[A, B, C]{id, f})
	e.mu.Unlock()
	return &sub3[A, B, C]{e: e, id: id}
}
func (e *event3[A, B, C]) emit(a A, b B, c C) {
	e.mu.RLock()
	snap := make([]Action3[A, B, C], 0, len(e.list))
	for _, it := range e.list {
		if it.f != nil {
			snap = append(snap, it.f)
		}
	}
	e.mu.RUnlock()
	for _, h := range snap {
		h(a, b, c)
	}
}
func (e *event3[A, B, C]) has() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, it := range e.list {
		if it.f != nil {
			return true
		}
	}
	return false
}

type sub3[A, B, C any] struct {
	once sync.Once
	e    *event3[A, B, C]
	id   uint64
}

func (s *sub3[A, B, C]) Unsubscribe() {
	s.once.Do(func() {
		s.e.mu.Lock()
		for i := range s.e.list {
			if s.e.list[i].id == s.id {
				s.e.list[i].f = nil
				break
			}
		}
		s.e.mu.Unlock()
	})
}

type event3Reader[A, B, C any] struct{ e *event3[A, B, C] }

func (r *event3Reader[A, B, C]) Subscribe(f Action3[A, B, C]) Subscription { return r.e.subscribe(f) }
func (r *event3Reader[A, B, C]) HasSubscribers() bool                      { return r.e.has() }

type event3Emitter[A, B, C any] struct{ e *event3[A, B, C] }

func (w *event3Emitter[A, B, C]) Emit(a A, b B, c C) { w.e.emit(a, b, c) }

/* ------- Event4 ------- */

type Event4[A, B, C, D any] interface {
	Subscribe(Action4[A, B, C, D]) Subscription
	HasSubscribers() bool
}
type Emitter4[A, B, C, D any] interface{ Emit(A, B, C, D) }

func NewEvent4[A, B, C, D any]() (Event4[A, B, C, D], Emitter4[A, B, C, D]) {
	e := &event4[A, B, C, D]{}
	return &event4Reader[A, B, C, D]{e}, &event4Emitter[A, B, C, D]{e}
}

type event4[A, B, C, D any] struct {
	mu     sync.RWMutex
	nextID uint64
	list   []entry4[A, B, C, D]
}
type entry4[A, B, C, D any] struct {
	id uint64
	f  Action4[A, B, C, D]
}

func (e *event4[A, B, C, D]) subscribe(f Action4[A, B, C, D]) Subscription {
	if f == nil {
		return noopSub{}
	}
	e.mu.Lock()
	id := e.nextID
	e.nextID++
	e.list = append(e.list, entry4[A, B, C, D]{id, f})
	e.mu.Unlock()
	return &sub4[A, B, C, D]{e: e, id: id}
}
func (e *event4[A, B, C, D]) emit(a A, b B, c C, d D) {
	e.mu.RLock()
	snap := make([]Action4[A, B, C, D], 0, len(e.list))
	for _, it := range e.list {
		if it.f != nil {
			snap = append(snap, it.f)
		}
	}
	e.mu.RUnlock()
	for _, h := range snap {
		h(a, b, c, d)
	}
}
func (e *event4[A, B, C, D]) has() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, it := range e.list {
		if it.f != nil {
			return true
		}
	}
	return false
}

type sub4[A, B, C, D any] struct {
	once sync.Once
	e    *event4[A, B, C, D]
	id   uint64
}

func (s *sub4[A, B, C, D]) Unsubscribe() {
	s.once.Do(func() {
		s.e.mu.Lock()
		for i := range s.e.list {
			if s.e.list[i].id == s.id {
				s.e.list[i].f = nil
				break
			}
		}
		s.e.mu.Unlock()
	})
}

type event4Reader[A, B, C, D any] struct{ e *event4[A, B, C, D] }

func (r *event4Reader[A, B, C, D]) Subscribe(f Action4[A, B, C, D]) Subscription {
	return r.e.subscribe(f)
}
func (r *event4Reader[A, B, C, D]) HasSubscribers() bool { return r.e.has() }

type event4Emitter[A, B, C, D any] struct{ e *event4[A, B, C, D] }

func (w *event4Emitter[A, B, C, D]) Emit(a A, b B, c C, d D) { w.e.emit(a, b, c, d) }
