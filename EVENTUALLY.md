# Eventually

`assert.Eventually` waits for a user-supplied condition to become `true`. It is
most often used when code under test sets a value from another goroutine, so the
test needs to poll until the desired state is reached. This guide also covers
`assert.EventuallyWithT` and `assert.Never`.

## Variants at a glance

- `Eventually` polls until the condition returns `true` or the timeout fires.
- `EventuallyWithT` retries with a fresh `*assert.CollectT` each tick. It keeps
  only the last tick's errors.
- `Never` makes sure the condition stays `false` for the whole timeout.

## Signature and scheduling

```go
func Eventually(
    t TestingT,
    condition func() bool,
    waitFor time.Duration,
    tick time.Duration,
    msgAndArgs ...interface{}
) bool
```

- `condition` runs on its own goroutine. The first check happens immediately.
  Later checks run every `tick` while it keeps returning `false`.
- `waitFor` sets the maximum polling time. When the deadline passes, the
  assertion fails with "Condition never satisfied" and appends any
  `msgAndArgs` to the output.
- The return value is `true` when the condition succeeds before the timeout.
  It is `false` otherwise. The assertion also reports failures through `t`.

> [!Note]
> You must protect shared state between the test and the condition with
> mutexes, `atomic` types, or other synchronization tools.

## Exit and panic behavior

Since [PR #1809](https://github.com/stretchr/testify/pull/1809),
`assert.Eventually` handles each way the condition goroutine can finish:

- **Condition returns `true`:** Polling stops at once and the assertion passes.
- **Condition times out:** Polling keeps running until `waitFor` expires and the
  test fails with "Condition never satisfied".
- **Condition panics:** The panic is not recovered. The Go runtime prints the
  panic and stack trace, then stops the test run. This is the normal goroutine
  panic path.
- **Condition calls `runtime.Goexit`:** The assertion now fails immediately with
  "Condition exited unexpectedly". Earlier versions waited for `waitFor` and
  could hang after `t.FailNow()` or `require.*`.

### `EventuallyWithT` specifics

`assert.EventuallyWithT` runs the same polling loop but gives each tick a new
`*assert.CollectT` named `collect`:

- Returning from the closure without errors on `collect` marks the condition as
  satisfied and the assertion succeeds right away.
- Recording errors on `collect` (via `collect.Errorf`, `assert.*(collect, ...)`
  helpers, or `collect.FailNow()`) fails only that tick. Polling keeps going,
  and if the timeout hits, the last tick's errors replay on the parent `t`
  before "Condition never satisfied".
- Call `collect.FailNow()` to exit the tick quickly and move to the next poll.
- If the closure exits via `runtime.Goexit` without first recording errors on
  `collect`—for example, by calling `require.*` on the parent `t`—the assertion
  fails immediately with "Condition exited unexpectedly".
- Panics behave the same as in `assert.Eventually`. They are not recovered and
  stop the test process.

Use `collect` for assertions you want to retry on each tick. Call `require.*`
on the parent `t` when you want the test to stop immediately. The same rules
apply to `require.EventuallyWithT` and its helpers.

## `Never` specifics

`assert.Never` uses the same polling loop as `Eventually` but expects the
condition to stay `false`. If the condition ever returns `true`, the assertion
fails immediately with "Condition satisfied".

- If the condition panics, the panic is not recovered and the test process
  terminates.
- If the condition calls `runtime.Goexit`, the assertion fails immediately with
  "Condition exited unexpectedly".
- These behaviors match those of `assert.Eventually`.

> [!Note]
> `Never` only succeeds when it lasts the full timeout, so it cannot finish
> early. Prefer using `Eventually` to keep tests fast.

## Usage tips

- Pick a `tick` that balances quick feedback with the work the condition does.
  Very small ticks can turn into a busy loop.
- Run `go test -race` when `Eventually` works with other goroutines. Data races
  cause more flakiness than the assertion itself.
- Use `assert.EventuallyWithT` when you need richer diagnostics or multiple
  errors. Record the failures on the provided `CollectT` value.
- Call `require.*` on the parent `t` inside the condition when you need to stop
  the test immediately.
