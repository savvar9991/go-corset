(defpurefun ((vanishes! :@loob) x) x)

(defcolumns (A :@loob) B)
(defconstraint c1 ()
  (if A
      (vanishes! 0)
      (vanishes! B)))
