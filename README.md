# Orvyl's Go Kt

### orvyl/id
Provides a Id generator that can produce numeric and alphanumeric seeded with the machine data and start datetime. It uses [sony/sonyflake] and [hashids] as the core libraries.

##### How to use
Please visit the said libraries to understand deeper how the Ids being generated.
```
package main

import "github.com/orvyl/kit/id"
import "log"

func main() {
    // If the first param is set to false, it will generate a numeric.
    idGen, err := id.NewGenerator(true, id.Settings{UseAWSData:false})

    if err != nil {
        log.Fatalf("Failed to generate Id generator %v\n", err)
    }

    id, err := idGen.Next()
    if err != nil {
        log.Panicf("Failed to get the next Id %v\n", err)
    }
    log.Printf("ID: %v", id) // Sample output: 230Ngm29xRPq
}

```
If `idGen`'s first param is set to false, it will produce a numeric number provided by the [sony/sonyflake] library. E.g `49394138646315010`. For the alphanumeric, it is also this library that provides the numeric id to be hashed:
```
sonyflake --> generate id --> 49394138646315010 --> hashids --> hash(49394138646315010) --> 230Ngm29xRPq
```

With this, we can assure that generated Ids are unique (both numeric and alphanum).

License
----

MIT

[sony/sonyflake]: <https://github.com/sony/sonyflake>
[hashids]: <http://hashids.org/>
