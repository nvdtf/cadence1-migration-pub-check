package main

import (
	"fmt"
	"slices"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/cadence/runtime/sema"
)

const testCode1 = `
access(all) contract Test {

	access(all) resource interface TestPub {
		access(all) fun deposit()
	}

	access(all) resource R: TestPub {
		access(all) fun deposit() {
		}
		access(all) fun withdraw() {
		}
		access(self) fun cook() {
		}
	}

	access(all) resource R2 {
		access(all) fun withdraw() {
		}
	}
}
`

func main() {
	program, err := parser.ParseProgram(nil, []byte(testCode1), parser.Config{})
	if err != nil {
		panic(err)
	}

	interfaces := enumerateInterfaces(testCode1)
	// fmt.Println(interfaces)

	// find all resource declarations
	if program.SoleContractDeclaration() != nil {

		members := program.SoleContractDeclaration().Members

		for _, mem := range members.Composites() {
			if mem.CompositeKind == common.CompositeKindResource {

				// check if resource implements an interface
				if len(mem.Conformances) == 1 {

					implements := mem.Conformances[0].String()

					// fmt.Println("Found: " + mem.Identifier.String())
					// fmt.Println("Implements: " + implements)

					interfaceFunctions := interfaces[implements]

					// compare list of functions to find newly exposed functions
					for _, f := range mem.Members.Functions() {
						if !slices.Contains(interfaceFunctions, f.Identifier.String()) {

							// check public
							if f.Access == ast.AccessAll {
								fmt.Printf("Resource %s: %s is exposing %s\n", mem.Identifier.String(), implements, f.Identifier.String())
							}
						}
					}

				} else if len(mem.Conformances) > 1 {
					panic("Multiple conformances not supported")
				}

			}
		}
	}

	config := &sema.Config{}
	if config.AccessCheckMode == sema.AccessCheckModeDefault {
		config.AccessCheckMode = sema.AccessCheckModeNotSpecifiedUnrestricted
	}
	config.ExtendedElaborationEnabled = true

	checker, err := sema.NewChecker(
		program,
		common.StringLocation("test"),
		nil,
		config,
	)
	if err != nil {
		panic(err)
	}

	err = checker.Check()
	if err != nil {
		panic(err)
	}
}

// returns list of interfaces and their functions
func enumerateInterfaces(code string) map[string][]string {

	program, err := parser.ParseProgram(nil, []byte(code), parser.Config{})
	if err != nil {
		panic(err)
	}

	result := make(map[string][]string)

	if program.SoleContractDeclaration() != nil {

		members := program.SoleContractDeclaration().Members

		for _, mem := range members.Interfaces() {
			name := mem.Identifier.String()
			for _, f := range mem.Members.Functions() {
				result[name] = append(result[name], f.Identifier.String())
				// fmt.Println(f.Access)
			}
		}
	}

	return result
}
