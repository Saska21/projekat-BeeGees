package main //*

import (
	"fmt"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/failures"
	"github.com/danilokacanski/da/week03_04_parallel/link"
	"github.com/danilokacanski/da/week03_04_parallel/process"
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime"

	"github.com/danilokacanski/da/beegees"
)

func main() {
	// ==========================================================
	// IZABERI SCENARIO OVDE:
	// 1. "normal"    - Svi su pošteni (Happy Path)
	// 2. "silent"    - Lider View-a 1 ćuti (Testira Timeout i Slow View Change)
	// 3. "byzantine" - Lider View-a 1 šalje različite blokove (Testira Safety/Equivocation)
	// ==========================================================
	selectedScenario := "byzantine"

	fmt.Println("================================================================")
	fmt.Printf("  BEEGEES SIMULATOR - Scenario: %s\n", selectedScenario)
	fmt.Println("================================================================")

	allIDs := []process.ProcessID{"node-1", "node-2", "node-3", "node-4"}
	fl := link.NewFairLossLink(0.0, 42) //0.0-nema kubitka poruka
	sl := link.NewStubbornLink(fl)
	pl := link.NewPerfectLink(sl)
	fm := failures.NewNoFailure()

	rt := simrt.NewRuntime(pl, fm,
		simrt.WithMaxDuration(15*time.Second),
		simrt.WithIdleTimeout(5*time.Second), // IdleTimeout je bitan: ako niko ne priča ovoliko dugo, simulator se gasi
		// Usporavamo retransmit da ne bi spamovao terminal dok mi spavamo (time.Sleep)
		//simrt.WithRetransmitInterval(200*time.Millisecond),
		simrt.WithRetransmitInterval(1*time.Second),
		simrt.WithVerbose(false), //da se vide i mrezne poruke
	)

	for _, pid := range allIDs {
		behavior := beegees.BehaviorHonest

		// Postavljamo specifično ponašanje zavisno od scenarija
		if selectedScenario == "silent" && pid == "node-2" {
			behavior = beegees.BehaviorSilent
		} else if selectedScenario == "byzantine" && pid == "node-3" {
			behavior = beegees.BehaviorByzantine
		}
		fmt.Printf("Registrujem %s [%s]\n", pid, behavior)

		node := beegees.NewBeeGeesNode(pid, allIDs, behavior)
		rt.Register(node)
	}

	fmt.Printf(">>> Pokrećem simulaciju [%s]...\n", selectedScenario)
	rt.Run()
	fmt.Println("================================================================")
}
