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

This protocol's flagship use will be a decentralized DNS database, but it has many other potential uses, including turn-based network gaming, and a successor technology to Google/Apache Wave called Orchard.

### Data model

A document is made of a [DAG][dag] of events, each of which is an action signed by its author. It also has synchronization objects, which periodically tie the dominant event chain into the external timestamping service. Individual events do not need acceptor consensus, but syncs (which are an accumulation of events) do.

A document has in its content, an IRC channel location, and an odd number of download URLs. During bootstrapping, these URLs are downloaded, and compared to each other. If a majority are consistent (they do not have to be identical to be consistent - they just have to not contain conflicting information), we use the union of the consistent serializations as our starting point, and reject information that conflicts with that. If no consistent majority exists in the downloaded files, bootstrapping fails, protecting the user from disinformation.

The state structure, which is constructed from the application of a series of events upon an initial starting state, represents the contents of the document at the given point in time/history. This state is a JSON-compatible object, which includes metacontent such as the event handler, permissions information, and IRC channel location.

#### Event

 * ParentHash string
 * HandlerName string
 * Arguments map[string]interface{}

There will be certain built-in event types, like setting/copying/deleting values in the document. Basic stuff. These will have UPPERCASE names, like "SET". There will also, at some point, be a facility for custom event handlers written in Lua.

The primary reasoning for document-custom functions are permissions. Allowing everyone full write access is like letting other people log into your computer as root. Yuck. Custom event handlers allow you to make specific actions, like "edit my own comment", which allow people to interact with the document, without having access to the all-powerful building blocks of those actions. It also provides a mechanism for contextual validation- in a chess game, for example, whether a move with certain arguments is valid depends entirely on the state of the board.

The secondary reasoning is it allows for efficient expression of actions that are specific to the document. A good example - which is more network efficient, replacing a gigantic blob of text with an almost-identical one, or expressing that change in the form of a regex or diff?

Finally, it allows conceptually atomic (indivisible) changes to be atomic in the implementation. If you are expressing one *conceptual* change in the form of a bunch of low-level events, and someone builds off your halfway-broadcast event chain, and *their* chain becomes the official one... well, you just orphaned half of something that was intended to be transactional. That's one of the worst kinds of surprises, short of [sugar-free gummy bears][bears].

#### Sync

 * EventHash string
 * Signatures []Signature
 * Confirmation ProofOfExistence

There will be convenience methods involved to get the latest valid event, according to various degrees of pickiness.

#### Document

 * Channel IRCLocation
 * Downloads map[string][string]
 * Events map[string]Event
 * Syncs map[string]Sync

#### DocState

 * Version string // hash of last event applied, or "" if initial state
 * Content map[string]interface{}

### What is the correct latest event?

This question is really, what is the confirmed sync object with the most parents (i.e. deepest history chain)? The sync will point to (and therefore include) a specific event object, and imply that that chain of history is both valid and official.

To create a confirmed sync object:

 * Propose a sync object in the channel
 * Acceptors broadcast their signatures of the sync object
 * When there are enough sigs, the sync is "approved"
 * The approval hash is put into the timestamping service (Bitcoin blockchain)
 * The Proof of Existence is about the approval, rather than the sync object itself.

Forks in the chain of confirmed syncs will be rare, since it requires the document's paxos Acceptors to agree on creating the fork (or at least be confused enough to create a fork) in the first place. However, even then, we can resolve such things by simply treating the longest chain as correct.

There is still a conceivable attack where the longest chain is hidden, or constructed in secret, such that it can suddenly be revealed and replace the previously correct fork. This still requires 2 conditions:

 1. The bootstrap downloads are out of date or intentionally leave out the longest chain.
 2. All current paxos Acceptors in the document conspire, or at least a sufficient quorum of them.

The former can be solved by simply making sure that the Learners always update the bootstrap downloads every time a sync is confirmed. The latter is a false problem - if you cannot trust the Acceptors to not conspire in a secret channel, how can you trust them to perform their most basic validation functions? Such a document is doomed to failure anyways.

[dag]: https://en.wikipedia.org/wiki/Directed_acyclic_graph
[bears]: http://www.amazon.com/Haribo-Gummy-Candy-Sugarless-5-Pound/dp/B000EVQWKC/
