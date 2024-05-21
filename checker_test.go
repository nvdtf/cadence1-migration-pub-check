package main

import "testing"

func TestChecker(t *testing.T) {

	const testCode = `
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

	Check(testCode)
}
