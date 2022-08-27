# gomatrix
Classic matrix screen implemented with goroutines.

I know it's not the best choice to use goroutines here. I had to use global mutex locking and you can still see some race conditions going on. 
When two goroutines take one column it is not garanteed than they will be called in a certain order. 
In fact, their executions are interleaved and you can see artifacts because of that.
