// +build amd64

#include "textflag.h"

// func AsmFunction() int
TEXT ·AsmFunction(SB), NOSPLIT, $0-8
    MOVQ $42, AX
    MOVQ AX, ret+0(FP)
    RET
