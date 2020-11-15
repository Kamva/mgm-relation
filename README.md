Mongo Go Models(mgm) relation implements has-one,has-many relations for the 
[Mongo Go Models (mgm)](https://github.com/kamva/mgm) package.

**Hooks**
- `Syncing () error` : calls before sync. if return error, ew cancel sync and return that error to the caller.
- `Synced () error` : calls after sync. if return error, we return that error to the caller.

**Important Notes**: 
- This package use Mongo Go Models native methods, so you can not expect to have behavior of `mgn` (like set `ID` on the model, or update `created_at`,`updated_at` fields...).  
  You can write your sync hooks or use default `mgm-relation` implementation of sync hooks to handle it.

**TODO**
- [ ] We can also automatically set the foreign key field on each model before saving it. implement it if you need(find foreign key field on the related model by `bson` tag's value).