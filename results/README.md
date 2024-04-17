# Capability Analysis Results
We performed capability analysis on different Go packages, from different categories. The following are the most imported repos on GitHub, hosting a Go package.

### [Logrus](https://pkg.go.dev/github.com/sirupsen/logrus) 
    
**Category**: logging

**Imported by**: 178.805 packages

**Capabilities**:
        
    CAPABILITY_FILES: 123 references (119 direct, 4 transitive)
    CAPABILITY_NETWORK: 22 references (22 direct, 0 transitive)
    CAPABILITY_RUNTIME: 4 references (4 direct, 0 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 120 references (118 direct, 2 transitive)
    CAPABILITY_SYSTEM_CALLS: 127 references (0 direct, 127 transitive)
    CAPABILITY_UNANALYZED: 124 references (124 direct, 0 transitive)
    CAPABILITY_UNSAFE_POINTER: 1 references (0 direct, 1 transitive)
    CAPABILITY_REFLECT: 120 references (117 direct, 3 transitive)

### [cobra](https://pkg.go.dev/github.com/spf13/cobra) 
**Category**: productivity

**Imported by**: 125.960 packages

**Capabilities**:
        
    CAPABILITY_FILES: 123 references (76 direct, 47 transitive)
    CAPABILITY_NETWORK: 36 references (36 direct, 0 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 15 references (15 direct, 0 transitive)
    CAPABILITY_UNANALYZED: 78 references (10 direct, 68 transitive)
    CAPABILITY_UNSAFE_POINTER: 18 references (13 direct, 5 transitive)
    CAPABILITY_REFLECT: 22 references (13 direct, 9 transitive)

### [protobuf/reflect/protoreflect](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect) 

**Category**: productivity
    
**Imported by**: 81.933 packages

**Capabilities**:
        
    CAPABILITY_FILES: 149 references (22 direct, 127 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 2289 references (18 direct, 2271 transitive)
    CAPABILITY_MODIFY_SYSTEM_STATE: 3 references (3 direct, 0 transitive)
    CAPABILITY_UNANALYZED: 2259 references (236 direct, 2023 transitive)
    CAPABILITY_UNSAFE_POINTER: 4423 references (495 direct, 3928 transitive)
    CAPABILITY_REFLECT: 2371 references (345 direct, 2026 transitive)
    CAPABILITY_EXEC: 13 references (13 direct, 0 transitive)


### [glog](https://pkg.go.dev/github.com/golang/glog) 

**Category**: logging

**Imported by**: 73.536 packages

**Capabilities**:
    
    CAPABILITY_FILES: 78 references (75 direct, 3 transitive)
    CAPABILITY_RUNTIME: 12 references (12 direct, 0 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 73 references (8 direct, 65 transitive)
    CAPABILITY_MODIFY_SYSTEM_STATE: 1 references (1 direct, 0 transitive)
    CAPABILITY_SYSTEM_CALLS: 12 references (12 direct, 0 transitive)
    CAPABILITY_UNANALYZED: 77 references (76 direct, 1 transitive)
    CAPABILITY_UNSAFE_POINTER: 70 references (8 direct, 62 transitive)

### [k8s.io/client-go/rest](https://pkg.go.dev/k8s.io/client-go/rest) 

**Category**: cloud tools

**Imported by**: 57.095 packages

**Capabilities**:
    
    CAPABILITY_FILES: 3626 references (131 direct, 3495 transitive)
    CAPABILITY_NETWORK: 3594 references (120 direct, 3474 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 3586 references (59 direct, 3527 transitive)
    CAPABILITY_MODIFY_SYSTEM_STATE: 2 references (2 direct, 0 transitive)
    CAPABILITY_OPERATING_SYSTEM: 9 references (4 direct, 5 transitive)
    CAPABILITY_SYSTEM_CALLS: 3542 references (0 direct, 3542 transitive)
    CAPABILITY_ARBITRARY_EXECUTION: 586 references (0 direct, 586 transitive)
    CAPABILITY_UNANALYZED: 3593 references (74 direct, 3519 transitive)
    CAPABILITY_UNSAFE_POINTER: 4733 references (219 direct, 4514 transitive)
    CAPABILITY_REFLECT: 4558 references (60 direct, 4498 transitive)
    CAPABILITY_EXEC: 3535 references (5 direct, 3530 transitive)


### [Testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock) 
    
**Category**: testing

**Imported by**: 24.721 packages

**Capabilities**:
    
    CAPABILITY_FILES: 371 references (174 direct, 197 transitive)
    CAPABILITY_NETWORK: 50 references (26 direct, 24 transitive)
    CAPABILITY_UNANALYZED: 39 references (0 direct, 39 transitive)
    CAPABILITY_UNSAFE_POINTER: 110 references (34 direct, 76 transitive)
    CAPABILITY_REFLECT: 156 references (73 direct, 83 transitive)

### [errors](https://pkg.go.dev/github.com/juju/errors) 

**Category**: error handling

**Imported by**: 18.733 ackages

**Capabilities**:
    
    CAPABILITY_UNANALYZED: 36 references (36 direct, 0 transitive)

### [go-ethereum](https://pkg.go.dev/github.com/ethereum/go-ethereum) 

**Category**: productivity

**Imported by**: 11.189 packages

**Capabilities**:
    
    CAPABILITY_FILES: 3638 references (599 direct, 3039 transitive)
    CAPABILITY_NETWORK: 3516 references (700 direct, 2816 transitive)
    CAPABILITY_RUNTIME: 3370 references (360 direct, 3010 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 3439 references (92 direct, 3347 transitive)
    CAPABILITY_MODIFY_SYSTEM_STATE: 69 references (17 direct, 52 transitive)
    CAPABILITY_OPERATING_SYSTEM: 178 references (45 direct, 133 transitive)
    CAPABILITY_SYSTEM_CALLS: 3388 references (31 direct, 3357 transitive)
    CAPABILITY_ARBITRARY_EXECUTION: 3822 references (111 direct, 3711 transitive)
    CAPABILITY_CGO: 3505 references (39 direct, 3466 transitive)
    CAPABILITY_UNANALYZED: 3905 references (795 direct, 3110 transitive)
    CAPABILITY_UNSAFE_POINTER: 3851 references (151 direct, 3700 transitive)
    CAPABILITY_REFLECT: 3971 references (513 direct, 3458 transitive)
    CAPABILITY_EXEC: 3361 references (130 direct, 3231 transitive)

        
### [Gomega](https://pkg.go.dev/github.com/onsi/gomega) 

**Category**: testing

**Imported by**: 8.192 packages

**Capabilities**:
    
    CAPABILITY_FILES: 262 references (27 direct, 235 transitive)
    CAPABILITY_NETWORK: 241 references (78 direct, 163 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 69 references (17 direct, 52 transitive)
    CAPABILITY_OPERATING_SYSTEM: 10 references (10 direct, 0 transitive)
    CAPABILITY_SYSTEM_CALLS: 2 references (2 direct, 0 transitive)
    CAPABILITY_UNANALYZED: 93 references (26 direct, 67 transitive)
    CAPABILITY_UNSAFE_POINTER: 258 references (34 direct, 224 transitive)
    CAPABILITY_REFLECT: 274 references (106 direct, 168 transitive)
    CAPABILITY_EXEC: 15 references (15 direct, 0 transitive)


### [Ginko](https://pkg.go.dev/github.com/onsi/ginkgo) 
    
**Category**: testing

**Imported by**: 5.961 packages

**Capabilities**:
    
    CAPABILITY_FILES: 128 references (95 direct, 33 transitive)
    CAPABILITY_NETWORK: 122 references (80 direct, 42 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 119 references (74 direct, 45 transitive)
    CAPABILITY_MODIFY_SYSTEM_STATE: 29 references (11 direct, 18 transitive)
    CAPABILITY_OPERATING_SYSTEM: 32 references (21 direct, 11 transitive)
    CAPABILITY_SYSTEM_CALLS: 34 references (1 direct, 33 transitive)
    CAPABILITY_UNANALYZED: 198 references (87 direct, 111 transitive)
    CAPABILITY_UNSAFE_POINTER: 200 references (74 direct, 126 transitive)
    CAPABILITY_REFLECT: 96 references (61 direct, 35 transitive)
    CAPABILITY_EXEC: 63 references (45 direct, 18 transitive)


### [coredns/core/dnsserver](https://pkg.go.dev/github.com/coredns/coredns/core/dnsserver) 

**Category**: productivity

**Imported by**: 3.594 packages

**Capabilities**: 
    
    CAPABILITY_FILES: 410 references (57 direct, 353 transitive)
    CAPABILITY_NETWORK: 404 references (75 direct, 329 transitive)
    CAPABILITY_RUNTIME: 326 references (21 direct, 305 transitive)
    CAPABILITY_READ_SYSTEM_STATE: 336 references (27 direct, 309 transitive)
    CAPABILITY_MODIFY_SYSTEM_STATE: 73 references (3 direct, 70 transitive)
    CAPABILITY_OPERATING_SYSTEM: 309 references (0 direct, 309 transitive)
    CAPABILITY_SYSTEM_CALLS: 339 references (0 direct, 339 transitive)
    CAPABILITY_ARBITRARY_EXECUTION: 326 references (0 direct, 326 transitive)
    CAPABILITY_UNANALYZED: 435 references (54 direct, 381 transitive)
    CAPABILITY_UNSAFE_POINTER: 486 references (40 direct, 446 transitive)
    CAPABILITY_REFLECT: 336 references (6 direct, 330 transitive)
    CAPABILITY_EXEC: 309 references (0 direct, 309 transitive)


### [aws-sdk-go](https://pkg.go.dev/github.com/aws/aws-sdk-go/aws)


**Category**: cloud tools

**Imported by**: 35.193

**Capabilities**:

	TODO
