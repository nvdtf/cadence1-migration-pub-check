package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
)

type CheckerResult struct {
	resourceName string
	functionName string
}

func Check(code string) []CheckerResult {
	res := make([]CheckerResult, 0)

	program, err := parser.ParseProgram(nil, []byte(code), parser.Config{})
	if err != nil {
		fmt.Println(err)
		return res
	}

	interfaces := EnumerateInterfaces(code)
	// fmt.Println(interfaces)

	// find all resource declarations
	if program.SoleContractDeclaration() != nil {

		members := program.SoleContractDeclaration().Members

		for _, mem := range members.Composites() {
			if mem.CompositeKind == common.CompositeKindResource {

				// check if resource implements an interface
				var unionInterfaceFunctions []string
				for _, intName := range mem.Conformances {
					intNameStr := intName.String()
					functions, exists := interfaces[intNameStr]
					if !exists {
						functions, err = GetExternalInterfaceFunctions(code, intNameStr)
						if err != nil {
							fmt.Println(err)
							continue
						}
					}
					unionInterfaceFunctions = append(unionInterfaceFunctions, functions...)
				}

				// compare list of functions to find newly exposed functions
				if len(unionInterfaceFunctions) > 0 {
					for _, f := range mem.Members.Functions() {
						if !slices.Contains(unionInterfaceFunctions, f.Identifier.String()) {
							// check ignore list
							if ShouldIgnoreFunctionName(f.Identifier.String()) {
								continue
							}

							// check public
							if f.Access == ast.AccessAll {
								// fmt.Printf("Resource %s is exposing %s\n", mem.Identifier.String(), f.Identifier.String())
								res = append(res, CheckerResult{
									resourceName: mem.Identifier.String(),
									functionName: f.Identifier.String(),
								})
							}
						}
					}
				}
			}
		}
	}

	return res

	// config := &sema.Config{}
	// if config.AccessCheckMode == sema.AccessCheckModeDefault {
	// 	config.AccessCheckMode = sema.AccessCheckModeNotSpecifiedUnrestricted
	// }
	// config.ExtendedElaborationEnabled = true

	// checker, err := sema.NewChecker(
	// 	program,
	// 	common.StringLocation("test"),
	// 	nil,
	// 	config,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// err = checker.Check()
	// if err != nil {
	// 	panic(err)
	// }
}

// returns list of interfaces and their functions
func EnumerateInterfaces(code string) map[string][]string {

	program, err := parser.ParseProgram(nil, []byte(code), parser.Config{})
	if err != nil {
		panic(err)
	}

	result := make(map[string][]string)

	var members *ast.Members
	if program.SoleContractDeclaration() != nil {
		members = program.SoleContractDeclaration().Members
	} else if program.SoleContractInterfaceDeclaration() != nil {
		members = program.SoleContractInterfaceDeclaration().Members
	}

	for _, mem := range members.Interfaces() {
		name := mem.Identifier.String()
		result[name] = []string{}
		for _, f := range mem.Members.Functions() {
			result[name] = append(result[name], f.Identifier.String())
			// fmt.Println(f.Access)
		}
	}

	return result
}

func GetExternalInterfaceFunctions(code string, name string) ([]string, error) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return []string{}, fmt.Errorf("invalid interface name: %s", name)
	}
	contractName := parts[0]
	interfaceName := parts[1]

	// find import address
	r, _ := regexp.Compile(`import (?P<Contract>\w*) from 0x(?P<Address>[0-9a-f]*)`)
	matches := r.FindAllStringSubmatch(code, -1)
	address := ""

	for i := range matches {
		if matches[i][1] == contractName {
			address = matches[i][2]
		}
	}
	if len(address) == 0 {
		return []string{}, fmt.Errorf("interface %s not staged", name)
	}

	// get staged contract
	depCode, found := GetStagedContract(address, contractName)
	if !found {
		return []string{}, fmt.Errorf("staged dependency not found: %s.%s", address, contractName)
	}

	interfaces := EnumerateInterfaces(depCode)

	res, ok := interfaces[interfaceName]
	if !ok {
		return []string{}, fmt.Errorf("interface %s not implemented in %s", interfaceName, contractName)
	}

	return res, nil
}

func ShouldIgnoreFunctionName(name string) bool {
	// ignore list
	ignoreList := []string{
		"borrowViewResolver",
		"resolveView",
		"borrowNFT",
		"createEmptyCollection",
		"assetPath",
		"name",
		"description",
		"mediumimage",
		"thumbnail",
		"video",
		"fetchPrice",
	}
	if slices.Contains(ignoreList, name) {
		return true
	}

	// ignore read prefixes
	prefixList := []string{
		"get",
		"is",
		"has",
	}
	for _, p := range prefixList {
		if HasPrefix(name, p) {
			return true
		}
	}

	return false
}

func HasPrefix(s string, prefix string) bool {
	// prefix + capital letter
	r, _ := regexp.Compile(fmt.Sprintf("^%s[A-Z]", prefix))
	matches := r.FindAllStringSubmatch(s, -1)
	return len(matches) > 0
}
