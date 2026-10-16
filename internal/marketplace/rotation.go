package marketplace

import "fmt"

type Deployment struct{ Name, APIKey string }
type RotationPlan struct {
	SellerAssets []string
	BuyerUpdates []string
	HandoffOrder []string
}

func PlanRotation(deployments []Deployment, oldKey, newKey string) (RotationPlan, error) {
	if oldKey == "" || newKey == "" || oldKey == newKey {
		return RotationPlan{}, fmt.Errorf("distinct old and new keys are required")
	}
	p := RotationPlan{SellerAssets: []string{"catalog", "payout-profile"}, BuyerUpdates: []string{"checkout credentials", "order status webhook"}, HandoffOrder: []string{"provision-new-key", "deploy-new-key", "verify-usage", "retire-old-key"}}
	for _, d := range deployments {
		if d.APIKey == oldKey {
			p.HandoffOrder = append(p.HandoffOrder, "migrate:"+d.Name)
		}
	}
	return p, nil
}
