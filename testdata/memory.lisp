;; An example of a simple memory implemented using permutation sorts.
;; Each row in the trace represents a read (RW=0) or write (RW=1) to
;; memory.  For a read, the VAL column identifies the value read at
;; ADDR (i.e. the read address).  For a write, the VAL column holds
;; the value being written at ADDR (i.e. the write address).
;;
;; The constriants should enforce that we cannot construct values "out
;; of thin air".  That is, every read for a given address matches the
;; last write (or 0 if there was no previous write).  More
;; specifically, the most recent PC value where that addres was
;; written.

;; Program Counter (always increases by one)
(column PC :u16)
;; Read/Write flag (0=READ, 1=WRITE)
(column RW :u1)
;; Address being Read/Written
(column ADDR :u32)
;; Value being Read/Written
(column VAL :u8)
;; Permutation
(permute (ADDR' PC' RW' VAL') (+ADDR +PC +RW +VAL))
;; PC[0]=0
(vanish:first heartbeat_1 PC)
;; PC[k]=0 || PC[k]=PC[k-1]+1
(vanish heartbeat_2 (* PC (- PC (+ 1 (shift PC -1)))))
;; PC[k]=0 ==> PC[k-1]=0
(vanish heartbeat_3 (if PC (shift PC -1)))
;; PC[k]=0 ==> (RW[k]=0 && ADDR[k]=0 && VAL[k]=0)
(vanish heartbeat_4 (if PC (+ RW ADDR VAL)))
;; ADDR'[k] != ADDR'[k-1] ==> (RW'[k]=1 || VAL'[k]=0)
(vanish first_read_1 (ifnot (- ADDR' (shift ADDR' -1)) (* (- 1 RW') VAL')))
;; (RW'[0]=1 || VAL'[0]=0)
(vanish:first first_read_2 (* (- 1 RW') VAL'))
;; ADDR'[k] == ADDR'[k-1] ==> (RW=1 || VAL'[k]=VAL'[k-1])
(vanish next_read (if (- ADDR' (shift ADDR' -1)) (* (- 1 RW') (- VAL' (shift VAL' -1)))))
