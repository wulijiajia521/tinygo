package runtime

// This file implements Go interfaces.
//
// Interfaces are represented as a pair of {typecode, value}, where value can be
// anything (including non-pointers).

type _interface struct {
	typecode uint16
	value    *uint8
}

// Return true iff both interfaces are equal.
func interfaceEqual(x, y _interface) bool {
	if x.typecode != y.typecode {
		// Different dynamic type so always unequal.
		return false
	}
	if x.typecode == 0 {
		// Both interfaces are nil, so they are equal.
		return true
	}
	// TODO: depends on reflection.
	panic("unimplemented: interface equality")
}

// interfaceTypeAssert is called when a type assert without comma-ok still
// returns false.
func interfaceTypeAssert(ok bool) {
	if !ok {
		runtimePanic("type assert failed")
	}
}

// The following declarations are only used during IR construction. They are
// lowered to inline IR in the interface lowering pass.
// See compiler/interface-lowering.go for details.

type interfaceMethodInfo struct {
	signature *uint16 // external *i16 with a name identifying the Go function signature
	funcptr   *uint8  // bitcast from the actual function pointer
}

// Pseudo function call used while putting a concrete value in an interface,
// that must be lowered to a constant uint16.
func makeInterface(typecode *uint16, methodSet *interfaceMethodInfo) uint16

// Pseudo function call used during a type assert. A dumb implementation would
// do:
//
//     return *assertedType == actualType
//
// However, it is optimized during lowering, to emit false when this type assert
// can never happen.
func typeAssert(actualType uint16, assertedType *uint16) bool

// Pseudo function call that returns whether a given type implements all methods
// of the given interface.
func interfaceImplements(typecode uint16, interfaceMethodSet **uint16) bool

// Pseudo function that returns a function pointer to the method to call.
// See the interface lowering pass for how this is lowered to a real call.
func interfaceMethod(typecode uint16, interfaceMethodSet **uint16, signature *uint16) *uint8
