//go:build !linux || !amd64
// +build !linux !amd64

package tracers

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/virbicoin/go-virbicoin/common"
	"github.com/virbicoin/go-virbicoin/core/vm"
)

type Tracer struct{}

func New(code string, txCtx vm.TxContext) (*Tracer, error) {
	return &Tracer{}, nil
}

func (t *Tracer) CaptureStart(from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) error {
	return nil
}
func (t *Tracer) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, rStack *vm.ReturnStack, rdata []byte, contract *vm.Contract, depth int, err error) error {
	return nil
}
func (t *Tracer) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, memory *vm.Memory, stack *vm.Stack, rStack *vm.ReturnStack, contract *vm.Contract, depth int, err error) error {
	return nil
}
func (t *Tracer) CaptureEnd(output []byte, gasUsed uint64, t2 time.Duration, err error) error {
	return nil
}
func (t *Tracer) Stop(err error) {}
func (t *Tracer) GetResult() (json.RawMessage, error) {
	return nil, errors.New("JSトレーサはlinux/amd64のみサポートされています")
}
