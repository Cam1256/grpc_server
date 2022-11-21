# Authorization Server

Casbin server with basic CRUD operations in Golang to implement Role-Based Access Control (RBAC) with Tenants using SQLite3 storage model.
### TO-DO

- [ ] Check for corrupted data in request. Include tables
- [ ] gRPC  
### Usage
Assuming our database is like this
```
p, admin, tenant1, data1, read
p, admin, tenant2, data2, read
p, dev, tenant3, data1, write

g, alice, admin, tenant1
g, alice, user, tenant2
g, bob, dev, tenant3
```

#### Read all policies of a tenant
```curl -X GET http://localhost:8080/tenants/tenant3/policies ```
#### Create a new policy in a tenant
```curl -X POST -d '{"role":"dev", "resource":"data2", "permission":"write"}' http://localhost:8080/tenants/tenant3/policies```
#### Remove a policy from a tenant
```curl -X DELETE -d '{"role":"dev", "resource":"data2", "permission":"write"}' http://localhost:8080/tenants/tenant3/policies```
#### Update an existing policy in a tenant
```curl -X PUT -d '{"old":{"role":"dev", "resource":"data2", "permission":"write"}, "new":{"role":"dev", "resource":"data2", "permission":"read"}}' http://localhost:8080/tenants/tenant3/policies```
#### Read all roles of a tenant
```curl -X GET http://localhost:8080/tenants/tenant3/roles```
#### Create a new role definition for a tenant
```curl -X POST -d '{"user":"alice", "role":"dev"}' http://localhost:8080/tenants/tenant3/roles```
#### Remove a role from a tenant (and authorization rules involving the role in that tenant)
```curl -X DELETE -d '{"role":"admin"}' http://localhost:8080/tenants/tenant3/roles```
#### Update an existing definition role in a tenant (and update authorization rules involving that role)
```curl -X PUT -d '{"old":"admin", "new":"dev"}' http://localhost:8080/tenants/tenant3/roles```
#### Check user permission on resource in a tenant
```curl -X GET http://localhost:8080/tenants/tenant3/tenants/tenant3/resources/data1/permissions/read```
