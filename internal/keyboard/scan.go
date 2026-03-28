package keyboard

import "machine"

type Scanner struct {
	cols []machine.Pin
	rows []machine.Pin
	keys []bool // 事前確保（TinyGoでヒープ再確保を避ける）
}

func (s *Scanner) KeyCount() int { return len(s.keys) }

func (s *Scanner) Scan() []bool {
	for col, c := range s.cols {
		c.High()
		for row, r := range s.rows {
			s.keys[row*len(s.cols)+col] = r.Get()
		}
		c.Low()
	}
	return s.keys
}
