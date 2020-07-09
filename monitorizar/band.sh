docker stats --format "{{.Name}}:{{.NetIO}}" | tee -a >(ts "%d,%m,%y,%H,%M,%S" > band.txt)
