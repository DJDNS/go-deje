gdeje
=====

Golang library for DEJE Next, a protocol/technology for decentralized document hosting and concurrent editing.

Please note that this is not the same protocol as the original DEJE, which is vulnerable to a few security issues, and has never been completely implemented, essentially due to overengineering, and a conflict of concurrency models in Python.

gdeje requires you to be running a Bitcoin daemon, or at least point to one on another machine. Otherwise, it is self-contained, because of Golang's static linking (does require some dependencies to build).

## DEJE Next

 * Uses IRC as a communications bus - actually more efficient, because of all the many-recipient messages involved in the consensus algorithm.
 * Uses Bitcoin blockchain as a distributed timestamping service - with a bit of work, could be swapped out for any comparable service.
 * Simpler protocol and serialization format. No subscriptions, no snapshots.
 * Two layers of protection: document acceptor consensus, and blockchain-based checkpoint ordering.
 * Can bootstrap from multiple download URLs to ensure majority agreement.

### Data model

A document is made of a [DAG][dag] of events, each of which is an action signed by its author. It also has checkpoint objects, which periodically tie the dominant event chain into the external timestamping service. Individual events do not need acceptor consensus, but checkpoints (which are an accumulation of events) do.
