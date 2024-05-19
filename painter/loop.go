package painter

import (
	"architecture-lab-3/primitives"
	"image"
	"sync"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(ts primitives.TextureStateI)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	ts primitives.TextureStateI
	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start() {
	// TODO: стартувати цикл подій.

	l.ts = new(primitives.TextureState)

	l.stop = make(chan struct{})
	l.stopReq = false

	go func() {
		for !l.stopReq || !l.mq.empty() {
			op := l.mq.pull()

			if update := op.Do(l.ts); update {
				l.Receiver.Update(l.ts)
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(ts primitives.TextureStateI) {
		l.stopReq = true
	}))
	<-l.stop
}

// TODO: Реалізувати чергу подій.
type message struct {
	nextNode *message
	op       Operation
}

type messageQueue struct {
	mu         sync.Mutex
	pullSignal chan struct{}

	head, tail *message
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	thisNode := &message{op: op}

	if mq.head == nil {
		mq.head = thisNode
		mq.tail = thisNode
	} else {
		mq.tail.nextNode = thisNode
		mq.tail = thisNode
	}
	if mq.pullSignal != nil {
		close(mq.pullSignal)
		mq.pullSignal = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if mq.head == nil {
		mq.pullSignal = make(chan struct{})
		mq.mu.Unlock()

		<-mq.pullSignal

		mq.mu.Lock()
	}

	res := mq.head.op
	mq.head = mq.head.nextNode

	return res
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	return mq.head == nil
}
