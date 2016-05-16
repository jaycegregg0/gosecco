package compiler

import (
	"github.com/twtiger/gosecco/asm"
	"github.com/twtiger/gosecco/tree"
	. "gopkg.in/check.v1"
)

type JumpsSuite struct{}

var _ = Suite(&JumpsSuite{})

func (s *JumpsSuite) Test_maxSizeJumpSetsUnconditionalJumpPoint(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		DefaultPositiveAction: "allow", DefaultNegativeAction: "kill", DefaultPolicyAction: "kill",
		Rules: []*tree.Rule{
			&tree.Rule{
				Name: "write",
				Body: tree.BooleanLiteral{true},
			},
			&tree.Rule{
				Name: "vhangup",
				Body: tree.BooleanLiteral{true},
			},
			&tree.Rule{
				Name: "read",
				Body: tree.BooleanLiteral{true},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	4\n"+
		"jeq_k	01	00	C000003E\n"+
		"jmp	7\n"+
		"ld_abs	0\n"+
		"jeq_k	00	01	1\n"+
		"jmp	3\n"+
		"jeq_k	00	01	99\n"+
		"jmp	1\n"+
		"jeq_k	00	01	0\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsMulipleUnconditionalJumpPoint(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		DefaultPositiveAction: "allow", DefaultNegativeAction: "kill", DefaultPolicyAction: "kill",
		Rules: []*tree.Rule{
			&tree.Rule{
				Name: "write",
				Body: tree.BooleanLiteral{true},
			},
			&tree.Rule{
				Name: "read",
				Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	4\n"+
		"jeq_k	01	00	C000003E\n"+
		"jmp	C\n"+
		"ld_abs	0\n"+
		"jeq_k	00	01	1\n"+
		"jmp	8\n"+
		"jeq_k	01	00	0\n"+
		"jmp	5\n"+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	01	02\n"+
		"jmp	1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsWithTwoComparisons(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		DefaultPositiveAction: "allow", DefaultNegativeAction: "kill", DefaultPolicyAction: "kill",
		Rules: []*tree.Rule{
			&tree.Rule{
				Name: "write",
				Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
			&tree.Rule{
				Name: "read",
				Body: tree.Comparison{Op: tree.EQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	4\n"+
		"jeq_k	01	00	C000003E\n"+
		"jmp	13\n"+
		"ld_abs	0\n"+
		"jeq_k	01	00	1\n"+
		"jmp	7\n"+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jmp	B\n"+
		"jeq_x	00	01\n"+
		"jmp	8\n"+
		"jeq_k	01	00	0\n"+
		"jmp	5\n"+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	01	02\n"+
		"jmp	1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsWithNotEqual(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		DefaultPositiveAction: "allow", DefaultNegativeAction: "kill", DefaultPolicyAction: "kill",
		Rules: []*tree.Rule{
			&tree.Rule{
				Name: "write",
				Body: tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	4\n"+
		"jeq_k	01	00	C000003E\n"+
		"jmp	A\n"+
		"ld_abs	0\n"+
		"jeq_k	01	00	1\n"+
		"jmp	5\n"+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	02	01\n"+
		"jmp	1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}

func (s *JumpsSuite) Test_maxSizeJumpSetsWithNotEqualWithMoreThanOneRule(c *C) {
	ctx := createCompilerContext()
	ctx.maxJumpSize = 2

	p := tree.Policy{
		DefaultPositiveAction: "allow", DefaultNegativeAction: "kill", DefaultPolicyAction: "kill",
		Rules: []*tree.Rule{
			&tree.Rule{
				Name: "write",
				Body: tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
			&tree.Rule{
				Name: "read",
				Body: tree.Comparison{Op: tree.NEQL, Left: tree.NumericLiteral{42}, Right: tree.NumericLiteral{1}},
			},
		},
	}

	res, _ := ctx.compile(p)
	c.Assert(asm.Dump(res), Equals, ""+
		"ld_abs	4\n"+
		"jeq_k	01	00	C000003E\n"+
		"jmp	13\n"+
		"ld_abs	0\n"+
		"jeq_k	01	00	1\n"+
		"jmp	7\n"+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	00	01\n"+
		"jmp	A\n"+
		"jmp	8\n"+
		"jeq_k	01	00	0\n"+
		"jmp	5\n"+
		"ld_imm	1\n"+
		"st	0\n"+
		"ld_imm	2A\n"+
		"ldx_mem	0\n"+
		"jeq_x	02	01\n"+
		"jmp	1\n"+
		"ret_k	7FFF0000\n"+
		"ret_k	0\n")
}
