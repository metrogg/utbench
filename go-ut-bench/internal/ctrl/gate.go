// Package ctrl 提供任务运行时控制原语：暂停/恢复。
// 通过 context 传递 Gate 实例，让运行时的 worker 循环在处理
// 每个单位任务之前调用 ctrl.Wait(ctx)，当管理端请求暂停时自动阻塞。
package ctrl

import (
	"context"
	"sync"
)

// Gate 定义了可由外部切换状态的闸门。
type Gate interface {
	// Wait 在 Gate 处于暂停状态时阻塞，直到 Resume 被调用或 ctx 结束。
	// Gate 未暂停时立即返回 nil。
	Wait(ctx context.Context) error
	Pause()
	Resume()
	Paused() bool
}

// ChanGate 是基于 channel 的默认实现。
// 使用示例：
//
//	g := ctrl.NewChanGate()
//	ctx := ctrl.WithGate(parentCtx, g)
//	g.Pause()   // 所有之后调用 ctrl.Wait(ctx) 的 goroutine 会阻塞
//	g.Resume()  // 放行所有等待者
type ChanGate struct {
	mu     sync.Mutex
	paused bool
	// wait 由 Pause 创建、Resume 关闭。未暂停时为 nil。
	wait chan struct{}
}

// NewChanGate 创建一个初始"未暂停"的闸门。
func NewChanGate() *ChanGate { return &ChanGate{} }

// Pause 切换到暂停状态；重复调用幂等。
func (g *ChanGate) Pause() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.paused {
		return
	}
	g.paused = true
	g.wait = make(chan struct{})
}

// Resume 放行所有当前等待者；重复调用幂等。
func (g *ChanGate) Resume() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.paused {
		return
	}
	close(g.wait)
	g.wait = nil
	g.paused = false
}

// Paused 返回当前是否暂停。
func (g *ChanGate) Paused() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.paused
}

// Wait 在暂停时阻塞；ctx 结束时返回 ctx.Err()。
func (g *ChanGate) Wait(ctx context.Context) error {
	for {
		g.mu.Lock()
		if !g.paused {
			g.mu.Unlock()
			return nil
		}
		ch := g.wait
		g.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			// 被 Resume 唤醒，回到循环顶部再次确认
		}
	}
}

// ─── context 绑定 ────────────────────────────────────────────────────────────

type ctxKey struct{}

// WithGate 将 Gate 关联到 ctx；传入 nil 等价于原 ctx。
func WithGate(ctx context.Context, g Gate) context.Context {
	if g == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, g)
}

// Wait 是便利函数：若 ctx 上绑定了 Gate 则等待放行，否则立即返回。
func Wait(ctx context.Context) error {
	g, _ := ctx.Value(ctxKey{}).(Gate)
	if g == nil {
		return nil
	}
	return g.Wait(ctx)
}
