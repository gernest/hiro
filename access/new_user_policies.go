package access

import (
	"fmt"

	"github.com/gernest/hiro/resource"
	"github.com/ory/ladon"
	"github.com/satori/go.uuid"
)

// NewUserPolicies returns access control policies for a new user account. Users
// can only view/update/delete their own accounts.
func NewUserPolicies(userID uuid.UUID) []*ladon.DefaultPolicy {
	usr := userID.String()
	return []*ladon.DefaultPolicy{
		{
			ID: uuid.NewV4().String(),
			Description: fmt.Sprintf(`This allows %s to view ,update and delete his\her profile`,
				usr,
			),
			Subjects:  []string{usr},
			Effect:    ladon.AllowAccess,
			Resources: []string{resource.Profile},
			Actions:   []string{"view", "update", "delete"},
			Conditions: ladon.Conditions{
				"user": &ladon.EqualsSubjectCondition{},
			},
		},
		{
			ID: uuid.NewV4().String(),
			Description: fmt.Sprintf(`This allows %s to create,view ,update and delete his\her qrcodes`,
				usr,
			),
			Subjects:  []string{usr},
			Effect:    ladon.AllowAccess,
			Resources: []string{resource.QR},
			Actions: []string{"create", "view", "list", "update",
				"delete", "print"},
			Conditions: ladon.Conditions{
				"user": &ladon.EqualsSubjectCondition{},
			},
		},
		{
			ID: uuid.NewV4().String(),
			Description: fmt.Sprintf(`This allows %s to create,view ,update and delete his\her collections`,
				usr,
			),
			Subjects:  []string{usr},
			Effect:    ladon.AllowAccess,
			Resources: []string{resource.Collections},
			Actions:   []string{"create", "view", "list", "update", "delete"},
			Conditions: ladon.Conditions{
				"user": &ladon.EqualsSubjectCondition{},
			},
		},
	}
}
