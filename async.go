package eventx

import (
	"context"
	"sync"

	"github.com/lcylpzls/errx"
)

// asyncDispatcher 管理异步队列与 worker 池。
type asyncDispatcher struct {
	queue   chan asyncTask
	stop    chan struct{}
	wg      sync.WaitGroup
	workers int
	onError func(error)
}

// asyncTask 是异步投递单元。
type asyncTask struct {
	ctx     context.Context
	topic   string
	payload any
}

// newAsyncDispatcher 创建异步分发器。
func newAsyncDispatcher(cfg config) *asyncDispatcher {
	return &asyncDispatcher{
		queue:   make(chan asyncTask, cfg.queueSize),
		stop:    make(chan struct{}),
		workers: cfg.workers,
		onError: cfg.onError,
	}
}

// start 启动 worker 池。
func (d *asyncDispatcher) start(b *Bus) {
	d.wg.Add(d.workers)
	for i := 0; i < d.workers; i++ {
		go d.run(b)
	}
}

// PublishAsync 异步发布事件：投递到队列后立即返回。
// 队列满时返回 EVENTX_QUEUE_FULL；handler 错误通过 WithErrorHandler 上报。
func (b *Bus) PublishAsync(ctx context.Context, topic string, payload any) error {
	if b == nil || b.async == nil {
		return errx.NewCode(CodeInvalidOption, "总线未初始化，请使用 New 构造")
	}
	ctx = ensureContext(ctx)
	if err := validatePublishTopic(topic); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return errx.WrapCode(err, CodeCancelled, "发布上下文已取消")
	}
	// 在读锁内完成关闭检查与入队：保证要么返回 BusClosed，
	// 要么任务在 Close 排空之前进入队列。
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return errx.NewCode(CodeBusClosed, "总线已关闭")
	}
	select {
	case b.async.queue <- asyncTask{ctx: ctx, topic: topic, payload: payload}:
		b.publishes.Add(1)
		return nil
	default:
		return errx.NewCode(CodeQueueFull, "异步队列已满")
	}
}

// stopAndWait 停止 worker 并等待队列中在途任务完成。
func (d *asyncDispatcher) stopAndWait() {
	close(d.stop)
	d.wg.Wait()
}

// run 是单个 worker 的主循环；收到停止信号后先排空剩余任务再退出。
func (d *asyncDispatcher) run(b *Bus) {
	defer d.wg.Done()
	for {
		select {
		case task := <-d.queue:
			d.process(b, task)
		case <-d.stop:
			for {
				select {
				case task := <-d.queue:
					d.process(b, task)
				default:
					return
				}
			}
		}
	}
}

// process 投递单个异步任务并上报错误。
func (d *asyncDispatcher) process(b *Bus, task asyncTask) {
	if err := deliver(b, task.ctx, task.topic, task.payload, b.snapshot(task.topic)); err != nil && d.onError != nil {
		d.onError(err)
	}
}
