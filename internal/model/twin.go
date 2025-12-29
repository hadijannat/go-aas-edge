package model

import (
	"encoding/json"
	"sync"

	"github.com/aas-core-works/aas-core3.0-golang/jsonization"
	"github.com/aas-core-works/aas-core3.0-golang/types"
)

// TwinManager encapsulates the AAS Environment and ensures thread safety.
type TwinManager struct {
	Environment *types.Environment
	mu          sync.RWMutex

	// Direct pointers to properties allow O(1) updates without tree traversal
	TempProp *types.Property
	RPMProp  *types.Property
}

// NewTwinManager programmatically instantiates the AAS V3.0 Object Model.
func NewTwinManager() *TwinManager {
	// 1. Define Properties (Sensors)
	// AAS values are always strings in the model (xsd:string, xsd:float, etc.)
	initialTemp := "20.0"
	tempProp := types.NewProperty(types.DataTypeDefXSDFloat)
	tempProp.SetIDShort(ptr("Temperature"))
	tempProp.SetValue(&initialTemp)

	initialRPM := "0"
	rpmProp := types.NewProperty(types.DataTypeDefXSDInt)
	rpmProp.SetIDShort(ptr("RPM"))
	rpmProp.SetValue(&initialRPM)

	// 2. Define Submodel (Logical Grouping)
	smID := "https://example.com/ids/sm/operational_data"
	sm := types.NewSubmodel(smID)
	sm.SetIDShort(ptr("OperationalData"))
	sm.SetSubmodelElements([]types.ISubmodelElement{tempProp, rpmProp})

	// 3. Define the Asset Administration Shell
	aasID := "https://example.com/ids/aas/smart_motor_001"
	globalAssetID := "https://example.com/ids/asset/motor_001"
	assetInfo := types.NewAssetInformation(types.AssetKindInstance)
	assetInfo.SetGlobalAssetID(&globalAssetID)

	shell := types.NewAssetAdministrationShell(aasID, assetInfo)
	shell.SetIDShort(ptr("SmartMotorTwin"))

	// Create reference to submodel
	submodelRef := types.NewReference(
		types.ReferenceTypesModelReference,
		[]types.IKey{types.NewKey(types.KeyTypesSubmodel, smID)},
	)
	shell.SetSubmodels([]types.IReference{submodelRef})

	// 4. Wrap in Environment (Root Serialization Object)
	env := types.NewEnvironment()
	env.SetAssetAdministrationShells([]types.IAssetAdministrationShell{shell})
	env.SetSubmodels([]types.ISubmodel{sm})

	return &TwinManager{
		Environment: env,
		TempProp:    tempProp,
		RPMProp:     rpmProp,
	}
}

// UpdateState safely updates the twin's internal state.
func (tm *TwinManager) UpdateState(temp string, rpm string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.TempProp.SetValue(&temp)
	tm.RPMProp.SetValue(&rpm)
}

// Snapshot provides a read-only view for the AI Agent.
func (tm *TwinManager) Snapshot() map[string]string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return map[string]string{
		"Temperature": *tm.TempProp.Value(),
		"RPM":         *tm.RPMProp.Value(),
	}
}

// SerializeToJSON returns the standard-compliant JSON.
func (tm *TwinManager) SerializeToJSON() ([]byte, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Convert to jsonable representation first
	jsonable, err := jsonization.ToJsonable(tm.Environment)
	if err != nil {
		return nil, err
	}

	// Then marshal to JSON bytes
	return json.Marshal(jsonable)
}

// Helper for AAS string pointers
func ptr(s string) *string { return &s }
