package marketplace

import "testing"

func TestPlanRotationKeepsOldDeploymentsInHandoff(t *testing.T) {
	cases := []struct {
		name        string
		deployments []Deployment
		want        string
	}{
		{"old credential is scheduled", []Deployment{{Name: "orders", APIKey: "old"}}, "migrate:orders"},
		{"new credential needs no migration", []Deployment{{Name: "catalog", APIKey: "new"}}, "retire-old-key"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := PlanRotation(tc.deployments, "old", "new")
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, step := range p.HandoffOrder {
				if step == tc.want {
					found = true
				}
			}
			if !found {
				t.Fatalf("handoff=%v", p.HandoffOrder)
			}
		})
	}
}
