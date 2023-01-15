package cpu

// clear screen
func I00E0() {
  Display = [len(Display)]uint8{}
}

// jump
func I1NNN() {
}
