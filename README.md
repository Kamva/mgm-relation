Mongo Go Models(mgm) relation is the package to implement relations for the 
[Mongo Go Models (mgm)](https://github.com/Kamva/mgm) package.

**Hooks**
- `BeforeSync () error` : calls before sync. if return error, ew cancel sync and return that error to the caller.
- `AfterSync () error` : calls after sync. if return error, we return that error to the caller.

**Important Note**
This package use Mongo Go Models native methods, so you can not expect to have behavior of `mgn` (like set `ID` on the model, or update `created_at`,`updated_at` fields...).
All sort of this actions must handle by your self (you can handle It by Default `Sync` hooks).

**TODO**
- [ ] We can set the foreign key field value by the reflection package, implement if you need it(find foreign key field on the related model by `bson` tag value).  
