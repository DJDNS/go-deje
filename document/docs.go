// DEJE Documents and elements (Events and Quorums).
//
// Document contents are reconstructed through *history*, and history
// is a combination of Events and Quorums. For more information about
// how events are applied to reconstruct the contents of a document,
// see the go-deje/state package.
//
//   The history model works as follows:
//
// Events act like commits, expressing a change/delta from a parent
// state. They can be simple primitives, or complex custom events,
// but even complex events end up being boiled down into primitives
// internally when you apply them.
//
// Quorums act as "anchor points" bridging between Events and
// Timestamps. They represent a consensus about which event is the
// correct "tip" event (latest officially accepted event). Only
// certain identities are allowed to contribute signatures to a
// quorum - the Acceptor list is part of the document content, and
// can be almost *guaranteed* to change over time.
//
// Timestamps are polled via some external registry, like the Bitcoin
// network. They tie Quorums (consensus about tip Event) to
// specific times, or rather, a chronological order. Earlier timestamps
// take precedence over later ones, so that a set of former Acceptors
// can't rewrite history.
//
// More info will be available as contextual documentation when the
// TimestampTracker code is written. The code will demonstrate how the
// combination of Events + Quorums + Timestamps is turned into a single
// chain of accepted valid history. The comments will explain how that
// algorithm *works*, and the reasoning behind it.
package document
