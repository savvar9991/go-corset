(column X)
(module m1)
(column ST)
(column X :u16)
(permute (Y) (+X))
;; Ensure sorted column increments by 1
(vanish increment (* ST (- (shift Y 1) (+ 1 Y))))
