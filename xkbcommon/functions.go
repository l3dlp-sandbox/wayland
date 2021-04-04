package xkbcommon

/*
#cgo pkg-config: xkbcommon
#cgo LDFLAGS: -ldl

#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"
import "runtime"

// ComposeTableNewFromLocale Creates a compose table for a given locale.
//
// The locale is used for searching the file-system for an appropriate
// Compose file.  The search order is described in Compose(5).  It is
// affected by the following environment variables:
//
// 1. `XCOMPOSEFILE` - see Compose(5).
// 2. `XDG_CONFIG_HOME` - before `$HOME/.XCompose` is checked,
// `$XDG_CONFIG_HOME/XCompose` is checked (with a fall back to
// `$HOME/.config/XCompose` if `XDG_CONFIG_HOME` is not defined).
// This is a libxkbcommon extension to the search procedure in
// Compose(5) (since libxkbcommon 1.0.0). Note that other
// implementations, such as libX11, might not find a Compose file in
// this path.
// 3. `HOME` - see Compose(5).
// 4. `XLOCALEDIR` - if set, used as the base directory for the system's
// X locale files, e.g. `/usr/share/X11/locale`, instead of the
// preconfigured directory.
//
// Parameter context
// The library context in which to create the compose table.
// Parameter locale
// The current locale.  See  compose-locale.
//
// The value is copied, so it is safe to pass the result of getenv(3)
// (or similar) without fear of it being invalidated by a subsequent
// setenv(3) (or similar).
// Parameter flags
// Optional flags for the compose table, or 0.
//
// It returns A compose table for the given locale, or nil if the
// compilation failed or a Compose file was not found.
func (context *Context) ComposeTableNewFromLocale(locale string, flags uint32) (ct *ComposeTable) {
	ct = new(ComposeTable)

	ct.ct = C.xkb_compose_table_new_from_locale(context.cx, C.CString(locale), flags)
	runtime.SetFinalizer(ct, composeTableUnref)
	return ct
}

// ComposeStateNew Creates a new compose state object.
//
// Parameter table
// The compose table the state will use.
// Parameter flags
// Optional flags for the compose state, or 0.
//
// It returns A new compose state, or nil on failure.
func ComposeStateNew(table *ComposeTable, flags uint32) (cs *ComposeState) {
	cs = new(ComposeState)

	cs.cs = C.xkb_compose_state_new(table.ct, flags)
	runtime.SetFinalizer(cs, composeStateUnref)
	return cs
}

// ComposeTableUnref Releases a reference on a compose table, and possibly free it.
//
// Parameter table The object.  If it is nil, this function does nothing.
func ComposeTableUnref(*ComposeTable) {
}

func composeTableUnref(table *ComposeTable) {

	C.xkb_compose_table_unref(table.ct)
	table.ct = nil
}

// ComposeStateUnref Releases a reference on a compose state object, and possibly free it.
//
// Parameter state The object.  If nil, do nothing.
func ComposeStateUnref(*ComposeState) {
}

func composeStateUnref(state *ComposeState) {

	C.xkb_compose_state_unref(state.cs)
	state.cs = nil
}

// GetStatus Gets the current status of the compose state machine.
//
// See ComposeStatus
func (state *ComposeState) GetStatus() ComposeStatus {
	return ComposeStatus(C.xkb_compose_state_get_status(state.cs))
}

// GetOneSym Gets the result keysym for a composed sequence.
//
// See  compose-overview for more details.  This function is only
// useful when the status is ComposeComposed.
//
// It returns The result keysym.  If the sequence is not complete, or does
// not specify a result keysym, returns KeyNoSymbol.
func (state *ComposeState) GetOneSym() uint32 {
	return uint32(uint(C.xkb_compose_state_get_one_sym(state.cs)))
}

// Feed Feeds one keysym to the Compose sequence state machine.
//
// This function can advance into a compose sequence, cancel a sequence,
// start a new sequence, or do nothing in particular.  The resulting
// status may be observed with ComposeStateGetStatus().
//
// Some keysyms, such as keysyms for modifier keys, are ignored - they
// have no effect on the status or otherwise.
//
// The following is a description of the possible status transitions, in
// the format CURRENT STATUS => NEXT STATUS, given a non-ignored input
// keysym `keysym`:
//
//
// 	NOTHING or CANCELLED or COMPOSED =>
// 	NOTHING   if keysym does not start a sequence.
// 	COMPOSING if keysym starts a sequence.
// 	COMPOSED  if keysym starts and terminates a single-keysym sequence.
// 	COMPOSING =>
// 	COMPOSING if keysym advances any of the currently possible
// 	       sequences but does not terminate any of them.
// 	COMPOSED  if keysym terminates one of the currently possible
// 	       sequences.
// 	CANCELLED if keysym does not advance any of the currently
// 	       possible sequences.
//
//
// The current Compose formats do not support multiple-keysyms.
// Therefore, if you are using a function such as StateKeyGetSyms()
// and it returns more than one keysym, consider feeding KeyNoSymbol
// instead.
//
// Parameter state
// The compose state object.
// Parameter keysym
// A keysym, usually obtained after a key-press event, with a
// function such as StateKeyGetOneSym().
//
// It returns Whether the keysym was ignored.  This is useful, for example,
// if you want to keep a record of the sequence matched thus far.
func (state *ComposeState) Feed(keysym uint32) ComposeFeedResult {
	return ComposeFeedResult(C.xkb_compose_state_feed(state.cs, C.uint(keysym)))
}

// KeymapUnref Releases a reference on a keymap, and possibly free it.
//
// Parameter keymap The keymap.  If it is nil, this function does nothing.
func KeymapUnref(*Keymap) {
}

func keymapUnref(keymap *Keymap) {

	C.xkb_keymap_unref(keymap.km)
	keymap.km = nil
}

// KeymapNewFromString Creates a keymap from a keymap string.
//
// This is just like xkb_keymap_new_from_file(), but instead of a file, gets
// the keymap as one enormous string.
//
// See xkb_keymap_new_from_file()
func (context *Context) KeymapNewFromString(str []byte, a uint32, b uint32) (km *Keymap) {

	km = new(Keymap)

	km.km = C.xkb_keymap_new_from_string(context.cx, C.CString(string(str)), a, b)
	runtime.SetFinalizer(km, keymapUnref)
	return km
}

// StateNew Creates a new keyboard state object.
//
// Parameter keymap The keymap which the state will use.
//
// It returns A new keyboard state object, or nil on failure.
func (keymap *Keymap) StateNew() (st *State) {

	st = new(State)

	st.st = C.xkb_state_new(keymap.km)
	runtime.SetFinalizer(st, stateUnref)
	return st
}

// StateUnref Releases a reference on a keyboard state object, and possibly free it.
//
// Parameter state The state.  If it is nil, this function does nothing.
func StateUnref(*State) {
}
func stateUnref(state *State) {

	C.xkb_state_unref(state.st)
	state.st = nil
}

// ModGetIndex Gets the index of a modifier by name.
//
// It returns The index.  If no modifier with this name exists, returns
// ModInvalid.
//
// see also xkb_mod_index_t
func (keymap *Keymap) ModGetIndex(mod string) uint {
	return uint(C.xkb_keymap_mod_get_index(keymap.km, C.CString(mod)))
}

// KeyGetSyms Gets the keysyms obtained from pressing a particular key in a given
// keyboard state.
//
// Get the keysyms for a key according to the current active layout,
// modifiers and shift level for the key, as determined by a keyboard
// state.
//
// Parameter[in]  state    The keyboard state object.
// Parameter[in]  key      The keycode of the key.
// Parameter[out] symOut  The first of an immutable array of keysyms corresponding the
// key in the given keyboard state.
//
// This function does not perform any keysym-transformations.
// (This might change).
//
// It returns The first keysym in the symOut array.  If no keysyms
// are produced by the key in the given keyboard state, returns KeyNoSymbol and sets
// ok to false.
func (state *State) KeyGetSyms(key uint32) (symOut uint32, ok bool) {
	var data *C.uint

	if 0 != uint32(C.xkb_state_key_get_syms(state.st, C.uint(key), &data)) {
		if data == nil {
			return KeyNoSymbol, false
		}
		return uint32(*data), true
	}
	return KeyNoSymbol, false
}

// KeyGetOneSym Gets the single keysym obtained from pressing a particular key in a
// given keyboard state.
//
// This function is similar to StateKeyGetSyms(), but intended
// for users which cannot or do not want to handle the case where
// multiple keysyms are returned (in which case this function is
// preferred).
//
// It returns The keysym.  If the key does not have exactly one keysym,
// returns KeyNoSymbol
//
// This function performs Capitalization keysym-transformations.
//
// see also StateKeyGetSyms()
func (state *State) KeyGetOneSym(code uint32) uint32 {
	return uint32(C.xkb_state_key_get_one_sym(state.st, C.uint(code)))
}

// KeyRepeats Determines whether a key should repeat or not.
//
// A keymap may specify different repeat behaviors for different keys.
// Most keys should generally exhibit repeat behavior; for example, holding
// the 'a' key down in a text editor should normally insert a single 'a'
// character every few milliseconds, until the key is released.  However,
// there are keys which should not or do not need to be repeated.  For
// example, repeating modifier keys such as Left/Right Shift or Caps Lock
// is not generally useful or desired.
//
// It returns 1 if the key should repeat, 0 otherwise.
func (keymap *Keymap) KeyRepeats(code uint32) bool {
	return int(C.xkb_keymap_key_repeats(keymap.km, C.uint(code))) != 0
}

// ContextNew Creates a new context.
//
// Parameter flags Optional flags for the context, or 0.
//
// It returns A new context, or nil on failure.
func ContextNew(flags uint32) (cx *Context) {

	cx = new(Context)

	cx.cx = C.xkb_context_new(flags)
	runtime.SetFinalizer(cx, contextUnref)
	return cx
}

// ContextUnref Releases a reference on a context, and possibly free it.
//
// Parameter context The context.  If it is nil, this function does nothing.
func ContextUnref(*Context) {
}

func contextUnref(context *Context) {

	C.xkb_context_unref(context.cx)
	context.cx = nil
}

// KeysymToUtf32 Gets the Unicode/UTF-32 representation of a keysym.
//
// It returns The Unicode/UTF-32 representation of keysym, which is also
// compatible with UCS-4.  If the keysym does not have a Unicode
// representation, returns 0.
//
// This function does not perform any keysym-transformations.
// Therefore, prefer to use StateKeyGetUtf32() if possible.
//
// See also StateKeyGetUtf32()
func KeysymToUtf32(keysym uint32) uint32 {
	return uint32(C.xkb_keysym_to_utf32(C.uint(keysym)))
}

// KeyGetUtf32 Gets the Unicode/UTF-32 codepoint obtained from pressing a particular
// key in a a given keyboard state.
//
// Returns The UTF-32 representation for the key, if it consists of only
// a single codepoint.  Otherwise, returns 0.
//
// This function performs Capitalization and Control, see
// keysym-transformations.
//
// since 0.4.1
func (state *State) KeyGetUtf32(keysym uint32) uint32 {
	return uint32(C.xkb_state_key_get_utf32(state.st, C.uint(keysym)))
}

// GetUtf8 Gets the result Unicode/UTF-8 string for a composed sequence.
//
// See compose-overview for more details.  This function is only
// useful when the status is ComposeComposed.
//
// parameter state
//     The compose state.
//
// Returns
//   The bytes required for the string, excluding the NUL byte.
//   If the sequence is not complete, or does not have a viable result
//   string, sets `buffer` to the empty string.
func (state *ComposeState) GetUtf8() (buffer []byte) {
	var l = int(C.xkb_compose_state_get_utf8(state.cs, (*C.char)(nil), C.ulong(0)))
	if l == 0 {
		return
	}
	l++
	var buf = make([]C.char, l)
	l = int(C.xkb_compose_state_get_utf8(state.cs, (*C.char)(&buf[0]), C.ulong(l)))
	for i := 0; i < l; i++ {
		buffer = append(buffer, byte(buf[i]))
	}
	return
}

// UpdateMask Updates a keyboard state from a set of explicit masks.
//
// This entry point is intended for window systems and the like, where a
// master process holds an State, then serializes it over a wire
// protocol, and clients then use the serialization to feed in to their own
// State.
//
// All parameters must always be passed, or the resulting state may be
// incoherent.
//
// The serialization is lossy and will not survive round trips; it must only
// be used to feed slave state objects, and must not be used to update the
// master state.
//
// If you do not fit the description above, you should use
// xkb_state_update_key() instead.  The two functions should not generally be
// used together.
//
// it returns A mask of state components that have changed as a result of
// the update.  If nothing in the state has changed, returns 0.
//
// see also xkb_state_component
// see also xkb_state_update_key
func (state *State) UpdateMask(depressedMods, latchedMods, lockedMods,
	depressedLayout, latchedLayout, lockedLayout uint32) uint32 {
	return uint32(C.xkb_state_update_mask(state.st, C.uint(depressedMods),
		C.uint(latchedMods), C.uint(lockedMods),
		C.uint(depressedLayout), C.uint(latchedLayout),
		C.uint(lockedLayout)))
}

// SerializeMods is The counterpart to xkb_state_update_mask for modifiers, to be used on
// the server side of serialization.
//
// parameter state      The keyboard state.
// parameter components A mask of the modifier state components to serialize.
// State components other than StateMods* are ignored.
// If StateModsEffective is included, all other state components are
// ignored.
//
// it returns A xkb_mod_mask_t representing the given components of the
// modifier state.
//
// This function should not be used in regular clients; please use the
// xkb_state_mod_*_is_active API instead.
func (state *State) SerializeMods(mods StateComponent) uint32 {
	return uint32(C.xkb_state_serialize_mods(state.st, uint32(mods)))
}
