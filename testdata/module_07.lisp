(defcolumns (X :u8) (Y :u8))
;;
(module m1)
(defcolumns (X :u8) (Y :u8))
(defpermutation (A B) ((+ X) (+ Y)))
(defconstraint diag_ab () (- (shift A 1) B))
