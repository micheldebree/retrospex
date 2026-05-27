## 1. Specification Deployment

- [x] 1.1 Create the directory `openspec/specs/koala-export/`
- [x] 1.2 Move the reverse-engineered `spec.md` from the change directory to `openspec/specs/koala-export/spec.md`

## 2. Verification of Existing Implementation

- [x] 2.1 Add a unit test in `internal/c64io/koala_test.go` (or update existing) to verify the 10003 byte file size.
- [x] 2.2 Add a unit test to verify the memory layout order (Load Address, Bitmap, Screen RAM, Color RAM, Background).
- [x] 2.3 Add a unit test to verify correct color packing in Screen RAM for bit patterns 01 and 10.
