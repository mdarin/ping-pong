# ping-pong
Supervised classical ping-pong using goroutines

You should see this output in the console when program done

<code>go run ping_pong.go</code>
```
[APP] started
[PING] started
[PING] initialize conversation signal to Pong
[PONG] started
[PONG] Start a converstaion
[PING] Pong message: Hello, Ping
[SUP] started
[PONG] Ping message: Hello, Pong!
[PONG] terminate conversation signal to Ping
[PING] Slowing down...
[PING] Signal Pong to stop too
[PING] trerminated
[SUP] workers done: 1
[PONG] Slowing down...
[PONG] trerminated
[SUP] workers done: 2
[SUP] all workers done!
[SUP] terminated
[APP] done
````
