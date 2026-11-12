package vm

import (
	"fmt"
	"io"

	bc "gopy/internal/bytecode"
)

type builtin func([]any) (any, error)

func Run(code *bc.Code, output io.Writer) error {
	m := machine{globals: map[string]any{}, print: func(args []any) (any, error) {
		for i, arg := range args {
			switch arg {
			case nil:
				args[i] = "None"
			case true:
				args[i] = "True"
			case false:
				args[i] = "False"
			case (*bc.Code)(nil):
				args[i] = "<function>"
			default:
				if fn, ok := arg.(*bc.Code); ok {
					args[i] = "<function " + fn.Name + ">"
				}
			}
		}
		for i, arg := range args {
			if i > 0 {
				if _, err := fmt.Fprint(output, " "); err != nil {
					return nil, err
				}
			}
			if _, err := fmt.Fprint(output, arg); err != nil {
				return nil, err
			}
		}
		_, err := fmt.Fprintln(output)
		return nil, err
	}}
	_, err := m.frame(code, m.globals, 0)
	return err
}

type machine struct {
	globals map[string]any
	print   builtin
}

func (m machine) frame(code *bc.Code, locals map[string]any, depth int) (any, error) {
	if depth >= 1000 {
		return nil, fmt.Errorf("maximum recursion depth exceeded")
	}
	var stack []any
	push := func(v any) { stack = append(stack, v) }
	pop := func() any {
		i := len(stack) - 1
		v := stack[i]
		stack = stack[:i]
		return v
	}
	for pc := 0; pc < len(code.Instructions); {
		ins := code.Instructions[pc]
		pc++
		switch ins.Op {
		case bc.Push:
			push(ins.Arg)
		case bc.Load:
			name := ins.Arg.(string)
			value, ok := locals[name]
			if !ok {
				value, ok = m.globals[name]
			}
			if !ok {
				if name != "print" {
					return nil, fmt.Errorf("undefined name %q", name)
				}
				value = m.print
			}
			push(value)
		case bc.Store:
			locals[ins.Arg.(string)] = pop()
		case bc.Pop:
			pop()
		case bc.Unary:
			n, ok := pop().(int64)
			if !ok {
				return nil, fmt.Errorf("unary %s requires an integer", ins.Arg)
			}
			if ins.Arg == "-" {
				n = -n
			}
			push(n)
		case bc.Binary, bc.Compare:
			right := pop()
			left := pop()
			value, err := operation(ins.Arg.(string), left, right)
			if err != nil {
				return nil, err
			}
			push(value)
		case bc.Call:
			argc := ins.Arg.(int)
			start := len(stack) - argc
			args := stack[start:]
			stack = stack[:start]
			callee := pop()
			var value any
			var err error
			switch fn := callee.(type) {
			case builtin:
				value, err = fn(args)
			case *bc.Code:
				if argc != len(fn.Params) {
					return nil, fmt.Errorf("%s expects %d arguments, got %d", fn.Name, len(fn.Params), argc)
				}
				bindings := make(map[string]any, argc)
				for i, name := range fn.Params {
					bindings[name] = args[i]
				}
				value, err = m.frame(fn, bindings, depth+1)
			default:
				return nil, fmt.Errorf("value is not callable")
			}
			if err != nil {
				return nil, err
			}
			push(value)
		case bc.JumpFalse:
			switch pop() {
			case nil, false, int64(0):
				pc = ins.Arg.(int)
			}
		case bc.Jump:
			pc = ins.Arg.(int)
		case bc.Return:
			return pop(), nil
		}
	}
	return nil, nil
}

func operation(op string, left, right any) (any, error) {
	a, aOK := left.(int64)
	b, bOK := right.(int64)
	if !aOK || !bOK {
		return nil, fmt.Errorf("%s requires integer operands", op)
	}
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "==":
		return a == b, nil
	case "!=":
		return a != b, nil
	case "<":
		return a < b, nil
	case "<=":
		return a <= b, nil
	case ">":
		return a > b, nil
	case ">=":
		return a >= b, nil
	case "//", "%":
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		q, r := a/b, a%b
		if r != 0 && (r < 0) != (b < 0) {
			q--
			r += b
		}
		if op == "%" {
			return r, nil
		}
		return q, nil
	default:
		return nil, fmt.Errorf("unsupported operator %q", op)
	}
}
