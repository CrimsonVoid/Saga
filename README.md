# Saga - Columnar in-memory data store [![Build Status](https://travis-ci.org/CrimsonVoid/Saga.svg?branch=master)](https://travis-ci.org/CrimsonVoid/Saga)


### Design questions
- [ ] Should Table be immutable? We might be able get this rather cheaply?
- [x] ~~Should we have functions that are inefficient (such as NewFromMap)? No~~


### Table creation
- [x] New(columns []string, values ...[]interface{})
- [x] ~~NewFromMap(...map[string]interface{})~~
- [ ] Reserve ( `New(cols).Reserve(N).Insert(cols, vals...)` )

### Adding/Modifying row based data
- [x] Insert Rows(columns []string, values ...[]interface{})
- [x] ~~Insert Row(map[string]interface{})~~
- [ ] Append tables

### Adding/Modifying column based data
- [x] ~~Add Column (Removed)~~
- [x] Update column
- [ ] Rename columns (be sure to be smart if renamed colNames collide)
- [ ] Delete columns (should not be an error if colName does not exist)
- [ ] Select columns (is it an error if colName does not exist?)

### Removing/Reordering data
- [ ] Filter
- [ ] Order by
- [ ] Map

### Combinators
- [ ] Joins
- [ ] Group By/aggregate
