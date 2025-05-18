// a RBAC auth check by casbin
package casbin

import (
	"github.com/taluos/Malt/pkg/errors"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type RBACEnforcer struct {
	Enforcer *casbin.Enforcer
}

type EnforceMethod interface {
	VerifyAuth(role string, path string, method string) (bool, error)
	UpdateEnforcer() error
	AddRoleForUser(role string, user string) error
	DeleteRoleForUser(role string, user string) error
	AddPolicy(role string, path string, method string) error
	DeletePolicy(role string, path string, method string) error
	GetRolesForUser(user string) ([]string, error)
	GetUsersForRole(role string) ([]string, error)
}

func NewAdapter(db *gorm.DB, modelPath string, policyPath string) (*RBACEnforcer, error) {
	e, err := casbin.NewEnforcer(modelPath, policyPath)

	if db != nil {
		adapter, _ := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
		e, err = casbin.NewEnforcer(modelPath, adapter)
	}

	if err != nil {
		return nil, errors.New("casbin init error")
	}

	err = e.LoadPolicy()
	if err != nil {
		return nil, errors.New("casbin load policy error")
	}
	return &RBACEnforcer{Enforcer: e}, nil
}

func (e *RBACEnforcer) VerifyAuth(role string, path string, method string) (bool, error) {
	ok, err := e.Enforcer.Enforce(role, path, method)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (e *RBACEnforcer) UpdateEnforcer() error {
	err := e.Enforcer.LoadPolicy()
	if err != nil {
		return errors.Wrapf(err, "casbin load policy error")
	}
	return nil
}

func (e *RBACEnforcer) AddRoleForUser(role string, user string) error {
	_, err := e.Enforcer.AddRoleForUser(role, user)
	if err != nil {
		return errors.Wrapf(err, "casbin add role for user error")
	}
	return nil
}

func (e *RBACEnforcer) DeleteRoleForUser(role string, user string) error {
	_, err := e.Enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		return errors.Wrapf(err, "casbin delete role for user error")
	}
	return nil
}

func (e *RBACEnforcer) AddPolicy(role string, path string, method string) error {
	_, err := e.Enforcer.AddPolicy(role, path, method)
	if err != nil {
		return errors.Wrapf(err, "casbin add policy error")
	}
	return nil
}

func (e *RBACEnforcer) DeletePolicy(role string, path string, method string) error {
	_, err := e.Enforcer.RemovePolicy(role, path, method)
	if err != nil {
		return errors.Wrapf(err, "casbin delete policy error")
	}
	return nil
}

func (e *RBACEnforcer) GetRolesForUser(user string) ([]string, error) {
	roles, err := e.Enforcer.GetRolesForUser(user)
	if err != nil {
		return nil, errors.Wrapf(err, "casbin get roles for user error")
	}
	return roles, nil
}

func (e *RBACEnforcer) GetUsersForRole(role string) ([]string, error) {
	users, err := e.Enforcer.GetUsersForRole(role)
	if err != nil {
		return nil, errors.Wrapf(err, "casbin get users for role error")
	}
	return users, nil
}
