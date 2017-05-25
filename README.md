# Saga - Columnar in-memory data store [![Build Status](https://travis-ci.org/CrimsonVoid/Saga.svg?branch=master)](https://travis-ci.org/CrimsonVoid/Saga)


### Design questions
- [ ] Should Table be immutable? We might be able get this rather cheaply?
- [ ] Should we have functions that are inefficient (such as NewFromMap)?


### Table creation
- [x] New(columns []string, values ...[]interface{})
- [x] NewFromMap(...map[string]interface{})
- [ ] Reserve ( `New(cols).Reserve(N).Insert(cols, vals...)` )

### Adding/Modifying row based data
- [ ] Insert Rows(columns []string, values ...[]interface{})
- [x] Insert Row(map[string]interface{})
- [x] Add Column (TODO - Should there really be a distinction between Add and Update?)
- [x] Update column
- [ ] Append tables

### Adding/Modifying column based data
- [ ] Rename columns (be sure to be smart if renamed colNames collide)
- [ ] Delete column (should not be an error if colName does not exist)
- [ ] Select columns (is it an error if colName does not exist?)

### Removing/Reordering data
- [ ] Filter
- [ ] Order by
- [ ] Map

### Combinators
- [ ] Joins
- [ ] Group By/aggregate