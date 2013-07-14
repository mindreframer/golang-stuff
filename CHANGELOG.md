v0.3.3 - 22 Jun 2013
* Using TCP instead of UDP for communication between Listener and Replay
* Significally improved performance
* Fixed bugs causing locking and message droping (concurrency issues)
* Rewrote concurrency model to use more channels

v0.3 - 10 Jun 2013
* Use RAW_SOCKETS instead of tcpdump
* Own TCP stack
* All HTTP request types support
* Simplified request parsing