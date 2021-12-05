# Workflow for solution to Exercise 10.2

- Copy the `server` and `service` directory from the solution of Exercise 10.1
- Change the module name of `server` to `module github.com/practicalgo/book-exercise-solutions/chap10/exercise2`
- `waitForShutDown()` is where the logic is implemented (we use [signal.NotifyContext](https://pkg.go.dev/os/signal?utm_source=gopls#NotifyContext) 
   to setup signal handling)