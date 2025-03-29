# cue-schema

`cue-schema` validates backwards compatibility of [CUE language][cue] schemas. 

The intended usecase for this CLI is to transform other language APIs or schemas into a CUE description and then
use this CLI as a general tool for validating backwards compatibility of those converted schemas.

Example usage:

```bash
cue-schema breaking --old old.cue --new new.cue
```

[cue]: https://cuelang.org/