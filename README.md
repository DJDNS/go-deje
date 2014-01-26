go-deje
=====

Golang library for DEJE Next, a protocol/technology for decentralized document hosting and concurrent editing.

Please note that this is not the same protocol as the original DEJE, which is vulnerable to a few security issues, and has never been completely implemented, essentially due to overengineering, and a conflict of concurrency models in Python.

go-deje requires you to be running a Bitcoin daemon, or at least point to one on another machine. Otherwise, it is self-contained, because of Golang's static linking (does require some dependencies to build).

## DEJE Next

 * Uses IRC as a communications bus - actually more efficient, because of all the many-recipient messages involved in the consensus algorithm.
 * Uses Bitcoin blockchain as a distributed timestamping service - with a bit of work, could be swapped out for any comparable service.
 * Simpler protocol and serialization format. No subscriptions, no snapshots.
 * Two layers of protection: document acceptor consensus, and blockchain-based checkpoint ordering.
 * Can bootstrap from multiple download URLs to ensure majority agreement.

### Data model

A document is made of a [DAG][dag] of events, each of which is an action signed by its author. It also has synchronization objects, which periodically tie the dominant event chain into the external timestamping service. Individual events do not need acceptor consensus, but syncs (which are an accumulation of events) do.

A document has in its content, an IRC channel location, and an odd number of download URLs. During bootstrapping, these URLs are downloaded, and compared to each other. If a majority are consistent (they do not have to be identical to be consistent - they just have to not contain conflicting information), we use the union of the consistent serializations as our starting point, and reject information that conflicts with that. If no consistent majority exists in the downloaded files, bootstrapping fails, protecting the user from disinformation.

The state structure, which is constructed from the application of a series of events upon an initial starting state, represents the contents of the document at the given point in time/history. This state is a JSON-compatible object, which includes metacontent such as the event handler, permissions information, and IRC channel location.

[dag]: https://en.wikipedia.org/wiki/Directed_acyclic_graph
