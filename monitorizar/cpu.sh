docker stats --format "{{.Name}}:{{.CPUPerc}}" | tee -a >(ts "%d,%m,%y,%H,%M,%S" > cpu.txt) 
