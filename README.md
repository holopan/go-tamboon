# GO-TAMBOON ไปทำบุญ Challence

I use Hexagornal concept for this challence, the Hexagornal concept is easy way to change adaptors with no effect with business logic and tested.

### CONTENTS

* `internal/core/domain` - Structure model.
* `internal/core/ports` - Blue print of function for repositories and services.
* `internal/core/service` - Business logic.
* `internal/handler` - Hanler for call services.
* `internal/repository` - Database Gateway etc.
* `pkg/` - Helper method.
* `cmd/` - the main function.


### CONFIG
app:
  encrypt:
    caesar:
      shift: 128
  repository:
    omise:
      public: xxxxxxxxx
      secret: xxxxxxxxx
  currency: thb
  pool_size: 5