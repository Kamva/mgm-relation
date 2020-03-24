Mongo Go Models(mgm) relation implements has-one,has-many relations for the 
[Mongo Go Models (mgm)](https://github.com/Kamva/mgm) package.

**Hooks**
- `BeforeSync () error` : calls before sync. if return error, ew cancel sync and return that error to the caller.
- `AfterSync () error` : calls after sync. if return error, we return that error to the caller.

**Important Notes**: 
- This package use Mongo Go Models native methods, so you can not expect to have behavior of `mgn` (like set `ID` on the model, or update `created_at`,`updated_at` fields...).  
  You can write your sync hooks or use default `mgm-relation` implementation of sync hooks to handle it.

**TODO**
- [ ] We can set the foreign key field, implement if you need it(find foreign key field on the related model by `bson` tag's value).
- [ ] Add `Insert`,`Update` methods to the HasMany model to add/update single model in HasMany relation.
