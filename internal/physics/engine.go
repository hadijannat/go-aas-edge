package physics

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hadijannat/go-aas-edge/internal/model"
)

// Run starts the simulation loop in a non-blocking goroutine.
// It simulates realistic sensor behavior with inertia and thermodynamics.
func Run(twin *model.TwinManager) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		// Internal Physics State
		currentRPM := 0.0
		currentTemp := 20.0

		for range ticker.C {
			// 1. Target RPM varies randomly (Load changes)
			targetRPM := 1000.0 + (rand.Float64() * 1500.0) // 1000 - 2500 RPM

			// 2. Simulate Inertia (Smooth transition to target)
			currentRPM = currentRPM*0.9 + targetRPM*0.1

			// 3. Thermodynamics (Temp follows RPM with lag)
			// Formula: Ambient (20) + Friction Heat (RPM factor)
			targetTemp := 20.0 + (currentRPM * 0.035)
			currentTemp = currentTemp*0.95 + targetTemp*0.05

			// 4. Update the Twin
			twin.UpdateState(
				fmt.Sprintf("%.2f", currentTemp),
				fmt.Sprintf("%d", int(currentRPM)),
			)
		}
	}()
}
