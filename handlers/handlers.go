package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SalviCF/authorization-server/authorization"
)

// All documentation: https://casbin.io/docs/overview
// - API overview: https://casbin.io/docs/api-overview
// -- Global Management API: https://casbin.io/docs/management-api
// --- RBAC API: https://casbin.io/docs/rbac-api
// ---- RBAC with Tenants API: https://casbin.io/docs/rbac-with-domains-api, https://github.com/casbin/casbin/blob/master/rbac_api_with_domains.go

var e = authorization.InitEnforcer()

// Get Homepage. Ej: curl -X GET http://localhost:8080/
func GetHomepage(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "HOME\n")
	}
}

// POLICIES /////////////////////////////////////////////////////////////////////////////////////////////

// Read all policies of a tenant
// Ex: curl -X GET http://localhost:8080/tenants/tenant3/policies
func ReadPoliciesTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]
		policies := e.Enforcer.GetFilteredPolicy(1, tenant) // [role, tenant, resource, permission]
		json.NewEncoder(w).Encode(&policies)
	}
}

// Create a new policy in a tenant
// Ex: curl -X POST -d '{"role":"dev", "resource":"data2", "permission":"write"}' http://localhost:8080/tenants/tenant3/policies
func CreatePolicyTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		var policy map[string]string
		json.NewDecoder(r.Body).Decode(&policy)
		role, resource, permission := policy["role"], policy["resource"], policy["permission"]

		if !e.Enforcer.HasPolicy(role, tenant, resource, permission) { // create policy if doesn't exist
			added, err := e.Enforcer.AddPolicy(role, tenant, resource, permission)
			e.Enforcer.SavePolicy()
			if added {
				fmt.Fprint(w, "Created new authorization rule to the current tenant policy\n")
			} else {
				fmt.Fprintf(w, "error: %s", err)
			}
		} else {
			fmt.Fprint(w, "Authorization rule already exists in the current tenant policy\n")
		}
	}
}

// Remove a policy from a tenant
// Ex: curl -X DELETE -d '{"role":"dev", "resource":"data2", "permission":"write"}' http://localhost:8080/tenants/tenant3/policies
func DeletePolicyTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		var policy map[string]string
		json.NewDecoder(r.Body).Decode(&policy)
		role, resource, permission := policy["role"], policy["resource"], policy["permission"]

		if e.Enforcer.HasPolicy(role, tenant, resource, permission) { // delete policy if it exists
			deleted, err := e.Enforcer.RemovePolicy(role, tenant, resource, permission)
			e.Enforcer.SavePolicy()

			if deleted {
				fmt.Fprint(w, "Deleted authorization rule from the current tenant policy\n")
			} else {
				fmt.Fprintf(w, "error: %s", err)
			}
		} else {
			fmt.Fprint(w, "Authorization rule doesn't exist in the current tenant policy\n")
		}
	}
}

// Update an existing policy in a tenant
// Ex: curl -X PUT -d '{"old":{"role":"dev", "resource":"data2", "permission":"write"}, "new":{"role":"dev", "resource":"data2", "permission":"read"}}' http://localhost:8080/tenants/tenant3/policies
func UpdatePolicyTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		var policyPair map[string]map[string]string
		json.NewDecoder(r.Body).Decode(&policyPair)
		auxOldPolicy, auxNewPolicy := policyPair["old"], policyPair["new"]

		oldRole, oldResource, oldPermission := auxOldPolicy["role"], auxOldPolicy["resource"], auxOldPolicy["permission"]
		newRole, newResource, newPermission := auxNewPolicy["role"], auxNewPolicy["resource"], auxNewPolicy["permission"]
		oldPolicy := []string{oldRole, tenant, oldResource, oldPermission}
		newPolicy := []string{newRole, tenant, newResource, newPermission}

		if e.Enforcer.HasPolicy(oldPolicy) {
			if !e.Enforcer.HasPolicy(newPolicy) { // update policy if new policy doesn't exists
				updated, err := e.Enforcer.UpdatePolicy(oldPolicy, newPolicy)
				e.Enforcer.SavePolicy()
				if updated {
					fmt.Fprint(w, "Updated authorization rule in the current tenant policy\n")
				} else {
					fmt.Fprintf(w, "error: %s", err)
				}
			} else {
				fmt.Fprint(w, "Authorization rule already exists in the current tenant policy\n")
			}

		} else {
			fmt.Fprint(w, "Authorization rule doesn't exist in the current tenant policy\n")
		}
	}
}

// ROLES ////////////////////////////////////////////////////////////////////////////////////////////

// Read all roles of a tenant
// Ex: curl -X GET http://localhost:8080/tenants/tenant3/roles
// https://github.com/casbin/casbin/issues/655
// https://stackoverflow.com/questions/34018908/golang-why-dont-we-have-a-set-datastructure
func ReadRolesTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		policies := e.Enforcer.GetFilteredPolicy(1, tenant)                   // [role, tenant, resource, permission]
		roleDefs := e.Enforcer.GetFilteredNamedGroupingPolicy("g", 2, tenant) // [user role tenant]

		var allRoles []string

		mapPolicies := make(map[string]bool) // create a set using a map {"admin": true, "dev": true}
		for _, p := range policies {
			role := p[0]
			_, exist := mapPolicies[role]
			if !exist {
				mapPolicies[role] = true
				allRoles = append(allRoles, role)
			}
		}

		mapRoles := make(map[string]bool)
		for _, p := range roleDefs {
			role := p[1]
			_, exist := mapRoles[role]
			_, exist2 := mapPolicies[role]
			if !exist && !exist2 {
				mapRoles[role] = true
				allRoles = append(allRoles, role)
			}
		}

		fmt.Fprintf(w, "%v", allRoles)
	}
}

// Create a new role definition for a tenant
// Ex: curl -X POST -d '{"user":"alice", "role":"dev"}' http://localhost:8080/tenants/tenant3/roles
func CreateRoleTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		var roleDef map[string]string
		json.NewDecoder(r.Body).Decode(&roleDef)
		user, role := roleDef["user"], roleDef["role"]

		hasRole, _ := e.Enforcer.HasRoleForUser(user, role, tenant)
		if !hasRole { // create role if doesn't exist
			added, err := e.Enforcer.AddRoleForUserInDomain(user, role, tenant)
			e.Enforcer.SavePolicy()

			if added {
				fmt.Fprint(w, "Created new role definition in tenant\n")
			} else {
				fmt.Fprintf(w, "error: %s", err)
			}
		} else {
			fmt.Fprint(w, "Role already exists in the current tenant policy\n")
		}
	}
}

// Remove a role from a tenant. Also remove all authorization rules that involve that role in that tenant
// Ej: curl -X DELETE -d '{"role":"admin"}' http://localhost:8080/tenants/tenant3/roles
func DeleteRoleTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		var roleAux map[string]string
		json.NewDecoder(r.Body).Decode(&roleAux)
		role := roleAux["role"]

		// Remove role definitions in a tenant: g, _, role, _
		deletedRole := false
		users := e.Enforcer.GetAllUsersByDomain(tenant)
		for _, user := range users {
			hasRole, _ := e.Enforcer.HasRoleForUser(user, role, tenant)
			if hasRole {
				e.Enforcer.DeleteRoleForUserInDomain(user, role, tenant)
				deletedRole = true
			}
		}

		if deletedRole {
			fmt.Fprint(w, "Deleted role definition from the current tenant policy\n")
		} else {
			fmt.Fprint(w, "Role definition doesn't exist in the current tenant policy\n")
		}

		// Remove authorization rules that involved that role: p, role, _, _, _
		var deleted bool
		var err error

		policies := e.Enforcer.GetFilteredPolicy(0, role, tenant)
		if len(policies) > 0 { // delete role if it exists
			deleted, err = e.Enforcer.RemoveFilteredPolicy(0, role, tenant) // start form index 0 [role, tenant, resource, permission]
			if deleted {
				fmt.Fprint(w, "Deleted authorization rule/s involving the role from the current tenant policy\n")
			} else {
				fmt.Fprintf(w, "error: %s", err)
			}
		} else {
			fmt.Fprint(w, "Role not involved in any authorization rule in the current tenant policy\n")
		}

		e.Enforcer.SavePolicy()
	}
}

// Update an existing definition role in a tenant. Also update authorization rules involving that role
// Ej: curl -X PUT -d '{"old":"admin", "new":"dev"}' http://localhost:8080/tenants/tenant3/roles
func UpdateRoleTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := fields[0]

		var rolePair map[string]string
		json.NewDecoder(r.Body).Decode(&rolePair)
		oldRole := rolePair["old"]
		newRole := rolePair["new"]

		var updated bool
		var err error

		// Update authorization rules involving that role
		policies := e.Enforcer.GetFilteredPolicy(0, oldRole, tenant)

		if len(policies) > 0 {
			for _, policy := range policies {
				newPolicy := policy // policy is [role tenant resource permission]
				newPolicy[0] = newRole
				if !e.Enforcer.HasPolicy(newPolicy[0], newPolicy[1], newPolicy[2], newPolicy[3]) {
					updated, err = e.Enforcer.UpdatePolicy(policy, newPolicy) // bug: bool is contrary
				} else {
					fmt.Fprintf(w, "Authorization rule already exists in the current tenant policy\n")
				}
			}
		} else {
			fmt.Fprint(w, "Role not involved in any authorization rule in the current tenant policy\n")
		}

		if !updated { // need to do this beccause of a bug in casbin's returned bool
			fmt.Fprint(w, "Updated authorization rule/s in the current tenant policy\n")
		} else if err != nil {
			fmt.Fprintf(w, "error: %s", err)
		}

		// Update role definitions
		users := e.Enforcer.GetUsersForRoleInDomain(oldRole, tenant)

		if len(users) > 0 {
			for _, user := range users {
				oldRoleDef := []string{user, oldRole, tenant}
				newRoleDef := []string{user, newRole, tenant}
				hasRole, _ := e.Enforcer.HasRoleForUser(user, newRole, tenant)
				if !hasRole {
					updated, err = e.Enforcer.UpdateNamedGroupingPolicy("g", oldRoleDef, newRoleDef)
				} else {
					fmt.Fprint(w, "Role already exists in the current tenant policy\n")
				}
			}
		} else {
			fmt.Fprint(w, "Role definition doesn't exist in the current tenant policy\n")
		}

		if updated {
			fmt.Fprint(w, "Updated role definition in the current tenant policy\n")
		} else if err != nil {
			fmt.Fprintf(w, "error: %s", err)
		}

		e.Enforcer.SavePolicy()
	}
}

// Check user permission on resource in a tenant.
// curl -X GET http://localhost:8080/users/salva/tenants/tenant3/resources/data1/permissions/read
func CheckPermissionUserResourceTenant(fields []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, tenant, resource, permission := fields[0], fields[1], fields[2], fields[3]
		ok, err := e.Enforcer.Enforce(user, tenant, resource, permission)

		if err != nil {
			fmt.Fprintf(w, "error enforcer: %s", err)
			return
		}
		fmt.Fprintf(w, "Does %s have permission to %s %s in %s? %s\n",
			user, permission, resource, tenant, strconv.FormatBool(ok))

	}

}
