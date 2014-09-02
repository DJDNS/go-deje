// DEJE Documents and events.
//
// Document contents are reconstructed through *history*, and history
// is a graph of Events. For more information about how events are
// applied to reconstruct the contents of a document, see the
// go-deje/state package.
//
//   The history model works as follows:
//
// Events act like commits, expressing a change/delta from a parent
// state. They can be simple primitives, or complex custom events,
// but even complex events end up being boiled down into primitives
// internally when you apply them.
//
// Timestamps are polled via some external registry, like the Bitcoin
// network. They tie Events specific times, or rather, a chronological
// order. Earlier timestamps take precedence over later ones, so that
// a set of former Acceptors can't rewrite history.
//
// For more info, see the timestamps module.
package document
