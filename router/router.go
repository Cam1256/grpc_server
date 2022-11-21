package router

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/SalviCF/authorization-server/handlers"
)

type Route struct {
	Method  string
	Regex   *regexp.Regexp
	Handler func([]string) http.HandlerFunc
}

type Router struct {
	Routes []Route
}

// Routes definition using regular expressions
func InitRouter() *Router {
	return &Router{
		[]Route{
			// localhost:8080/
			{"GET",
				regexp.MustCompile(`^\/$`),
				handlers.GetHomepage},

			// http://localhost:8080/tenants/{tenantName}/policies
			{"GET",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/policies$`),
				handlers.ReadPoliciesTenant},

			// http://localhost:8080/tenants/{tenantName}/policies
			{"POST",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/policies$`),
				handlers.CreatePolicyTenant},

			// http://localhost:8080/tenants/{tenantName}/policies
			{"DELETE",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/policies$`),
				handlers.DeletePolicyTenant},

			// http://localhost:8080/tenants/{tenantName}/policies
			{"PUT",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/policies$`),
				handlers.UpdatePolicyTenant},

			// http://localhost:8080/tenants/{tenantName}/roles
			{"GET",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/roles$`),
				handlers.ReadRolesTenant},

			// http://localhost:8080/tenants/{tenantName}/roles
			{"POST",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/roles$`),
				handlers.CreateRoleTenant},

			// http://localhost:8080/tenants/{tenantName}/roles
			{"DELETE",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/roles$`),
				handlers.DeleteRoleTenant},

			// http://localhost:8080/tenants/{tenantName}/roles
			{"PUT",
				regexp.MustCompile(`^\/tenants\/([^/]+)\/roles$`),
				handlers.UpdateRoleTenant},

			// http://localhost:8080/users/{user}/tenants/{tenant}/resources/{resource}/permissions/{permission}
			{"GET",
				regexp.MustCompile(`^\/users\/([^/]+)\/tenants\/([^/]+)\/resources\/([^/]+)\/permissions\/([^/]+)$`),
				handlers.CheckPermissionUserResourceTenant},
		},
	}
}

func (router *Router) Serve(w http.ResponseWriter, r *http.Request) {
	var allowed []string
	for _, route := range router.Routes {
		matches := route.Regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.Method {
				allowed = append(allowed, route.Method)
				continue
			}
			w.WriteHeader(http.StatusOK)
			handler := route.Handler(matches[1:])
			handler(w, r)
			return
		}
	}
	if len(allowed) > 0 {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}
