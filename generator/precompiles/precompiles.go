// Copyright 2020 Marius van der Wijden
// This file is part of the fuzzy-vm library.
//
// The fuzzy-vm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The fuzzy-vm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the fuzzy-vm library. If not, see <http://www.gnu.org/licenses/>.

package precompiles

import (
	"math/big"

	"github.com/MariusVanDerWijden/FuzzyVM/filler"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/goevmlab/program"
)

var (
	precompiles = []precompile{
		new(ecdsaCaller),
		new(sha256Caller),
		new(ripemdCaller),
		new(identityCaller),
		new(bigModExpCaller),
		new(bn256Caller),
		new(bn256MulCaller),
		new(bn256PairingCaller),
		new(blake2fCaller),
	}
)

type precompile interface {
	call(p *program.Program, f *filler.Filler) error
}

type CallObj struct {
	Gas       *big.Int
	Address   common.Address
	Value     *big.Int
	InOffset  uint32
	InSize    uint32
	OutOffset uint32
	OutSize   uint32
}

func CallRandomizer(p *program.Program, f *filler.Filler, c CallObj) {
	switch f.Byte() % 3 {
	case 0:
		p.Call(c.Gas, c.Address, c.Value, c.InOffset, c.InSize, c.OutOffset, c.OutSize)
	case 1:
		p.CallCode(c.Gas, c.Address, c.Value, c.InOffset, c.InSize, c.OutOffset, c.OutSize)
	case 2:
		p.StaticCall(c.Gas, c.Address, c.InOffset, c.InSize, c.OutOffset, c.OutSize)
	}
}

func CallPrecompile(p *program.Program, f *filler.Filler) {
	// call a precompile
	var (
		idx  = int(f.Byte()) % len(precompiles)
		prec = precompiles[idx]
	)
	if err := prec.call(p, f); err != nil {
		panic(err)
	}
}