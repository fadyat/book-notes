#include "textflag.h"

TEXT ·sqrt(SB),NOSPLIT,$0
	FMOVD	x+0(FP), F0
	FSQRTD	F0, F0
	FMOVD	F0, ret+8(FP)
	RET
