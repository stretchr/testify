# Eventually

`assert.Eventually` waits for a user supplied condition to become `true`. It is
most often used when the code under test sets a value from another goroutine and
the test needs to poll until the desired state is reached.

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

- `condition` is called on its own goroutine. The first evaluation happens
  immediately. Subsequent evaluations are triggered every `tick` as long as the
  condition keeps returning `false`.
- `waitFor` defines the maximum amount of time `Eventually` will spend polling. If
  the deadline expires, the assertion fails with "Condition never satisfied" and
  the optional `msgAndArgs` are appended to the failure output.
- The return value is `true` when the condition succeeds before the timeout and
  `false` otherwise. The assertion also reports the failure through `t`.
- All state that is shared between the test and the condition must be protected
  for concurrent access via mutexes, `atomic` types, and other synchronization
  mechanisms.

## Exit and panic behavior

Since [PR #1809](https://github.com/stretchr/testify/pull/1809) `assert.Eventually`
distinguishes the different ways the condition goroutine can terminate:

- **Condition returns `true`:** `Eventually` stops polling immediately and
  succeeds.
- **Condition times out:** `Eventually` keeps polling until `waitFor` elapses and
  then fails the test with "Condition never satisfied".
- **Condition panics:** The panic is *not* recovered. The Go runtime terminates
  the process, prints the panic message and stack trace to standard error, and
  the test run stops. This matches the normal behavior of panics in goroutines.
- **Condition calls `runtime.Goexit`:** `Eventually` now fails the test
  immediately with "Condition exited unexpectedly".
  Before [PR #1809](https://github.com/stretchr/testify/pull/1809)
  the assertion waited until `waitFor` expired, causing tests that called
  `t.FailNow()` (or `require.*` helpers that use it) to hang.
  The new behavior surfaces the failure as soon as it happens.

### `EventuallyWithT` specifics

`assert.EventuallyWithT` runs the same polling loop but supplies each tick with
a fresh `*assert.CollectT`:

- Returning from the closure without recording errors on `collect` marks the
  condition as satisfied and the assertion succeeds immediately.
- Recording errors on `collect` (via `collect.Errorf`, `assert.*(collect, ...)`
  helpers, or `collect.FailNow()`) marks just that tick as failed. The polling
  continues, and if the timeout expires, the errors captured during the final
  tick are replayed on the parent `t` before emitting "Condition never satisfied".
- If the closure exits via `runtime.Goexit` *without* first recording errors on
  `collect`—for example by calling `require.*` on the parent `t`—the assertion
  fails immediately with "Condition exited unexpectedly".
- Panics behave the same as in `assert.Eventually`: they are not recovered and
  crash the test process.

Use `collect` for tick-scoped assertions you want to keep retrying, and call
`require.*` on the parent `t` when you want the test to stop right away. The
same rules apply to `require.EventuallyWithT` and its helpers.

## `Never` specifics

`assert.Never` runs the same polling loop as `Eventually` but expects the
condition to always return `false`. If the condition ever returns `true`, the
assertion fails immediately with "Condition satisfied".

- If the condition panics, the panic is not recovered and the test process
  terminates.
- If the condition calls `runtime.Goexit`, the assertion fails immediately with
  "Condition exited unexpectedly".
- These behaviors match those of `assert.Eventually`.

> [!Note]
> Since `Never` needs to run until the timeout expires to be successful,
> it cannot succeed early like `Eventually`. Prefer `Eventually` when possible
> to keep tests fast.

## Usage tips

- Pick a `tick` that balances fast feedback with the overhead of running the
  condition. Extremely small ticks can create busy loops.
- Run your test suite with `go test -race`, esp when `Eventually` coordinates
  with other goroutines. Data races are a more common source of flakiness
  than the assertion logic itself.
- If the condition needs to report rich diagnostics or multiple errors, prefer
  `assert.EventuallyWithT` and record failures on the provided `CollectT` value.
- It is safe to call `require.*` helpers on the parent `t` from within the
  condition closure to stop the test immediately.
