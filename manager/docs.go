// Collections of Document elements (Events and Quorums).
//
// These are more than just map[string]Thingy. They allow for
// getting items according to groups, giving all the items access
// *to each other*, and providing convenient APIs for Contains()
// testing.
//
// That said, it may make sense in the future to simplify out
// Managers entirely, and use typed maps, if a solution can be found
// to the grouping functionality, without a lot of duplicate code.
package manager
