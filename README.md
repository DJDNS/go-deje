go-deje
=====

Golang library for DEJE Next, a protocol/technology for decentralized document hosting and concurrent editing.

Please note that this is not the same protocol as the original DEJE, which is vulnerable to a few security issues, and has never been completely implemented, essentially due to overengineering, and a conflict of concurrency models in Python.

go-deje requires you to be running a Bitcoin daemon, or at least point to one on another machine. Otherwise, it is self-contained, because of Golang's static linking (does require some dependencies to build).

## DEJE Next

 * Uses WAMP as a communications bus - actually more efficient, because of all the many-recipient messages involved in the consensus algorithm.
 * Uses Bitcoin blockchain as a distributed timestamping service - with a bit of work, could be swapped out for any comparable service.
 * Simpler protocol and serialization format. No subscriptions, no snapshots.
 * Two layers of protection: document acceptor consensus, and blockchain-based checkpoint ordering.
 * Can bootstrap from multiple download URLs to ensure majority agreement.

This protocol's flagship use will be a decentralized DNS database, but it has many other potential uses, including turn-based network gaming, and a successor technology to Google/Apache Wave called Orchard.

### Data model

A document is made of a [DAG][dag] of events, each of which is an action signed by its author. It also has timestamps, which use an external timestamping service, which bind events to a specific order.

The state structure, which is constructed from the application of a series of events upon an initial starting state, represents the contents of the document at the given point in time/history. This state is a JSON-compatible object, which includes metacontent such as the event handlers and permissions information.

#### Event

There will be certain built-in event types, like setting/copying/deleting values in the document. Basic stuff. These will have UPPERCASE names, like "SET". There will also, at some point, be a facility for custom event handlers written in Lua.

The primary reasoning for document-custom functions are permissions. Allowing everyone full write access is like letting other people log into your computer as root. Yuck. Custom event handlers allow you to make specific actions, like "edit my own comment", which allow people to interact with the document, without having access to the all-powerful building blocks of those actions. It also provides a mechanism for contextual validation- in a chess game, for example, whether a move with certain arguments is valid depends entirely on the state of the board.

The secondary reasoning is it allows for efficient expression of actions that are specific to the document. A good example - which is more network efficient, replacing a gigantic blob of text with an almost-identical one, or expressing that change in the form of a regex or diff?

Finally, it allows conceptually atomic (indivisible) changes to be atomic in the implementation. If you are expressing one *conceptual* change in the form of a bunch of low-level events, and someone builds off your halfway-broadcast event chain, and *their* chain becomes the official one... well, you just orphaned half of something that was intended to be transactional. That's one of the worst kinds of surprises, short of [sugar-free gummy bears][bears].

#### Timestamp

A single timestamp in the external timestamping service. Imposes a mostly-reasonable, somewhat-arbitrary order on when events happened, and this order allows us to pick a single official chain of events.

Timestamps are ordered first by their blockheight (timestamps in earlier blocks always happen before timestamps in later blocks), then by vote count, and then by a string sort of their hashes (as a tiebreaker for multiple timestamps in a single block).

This allows for odd orders of confirmation sometimes - a child event may be confirmed earlier than its parent, for example - but this is harmless, those events will only be applied once.

### What is the correct latest event?

We get all the timestamps for the document. Then we apply the following algorithm, iterating through timestamps in serialization order:

 * Get event data from peers, including ancestors.
 * Is the event a fast-forward of the current tip? If not, drop it.

Then just traverse the chain of the tip's history, ignore any events that are invalid, and there you go.

[dag]: https://en.wikipedia.org/wiki/Directed_acyclic_graph
[bears]: http://www.amazon.com/Haribo-Gummy-Candy-Sugarless-5-Pound/dp/B000EVQWKC/
