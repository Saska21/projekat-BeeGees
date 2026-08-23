// *pokretanje simulatora i podešavanje scenarija
package main

import (
	"fmt"  //ispisivanje informacija u terminal
	"time" //podešavanje trajanja simulacije i timeout-a

	"github.com/danilokacanski/da/week03_04_parallel/failures"      //ponašanje kvarova
	"github.com/danilokacanski/da/week03_04_parallel/link"          //mreža između čvorova, simulira komunikacionu mrežu između procesa
	"github.com/danilokacanski/da/week03_04_parallel/process"       //obezbedjuje identitet (ID) procesa/cvora, kao i strukturu poruke koju simulator koristi
	simrt "github.com/danilokacanski/da/week03_04_parallel/runtime" //pokreće sistem

	"github.com/danilokacanski/da/beegees"
)

func main() {
	// ==========================================================
	// IZABERI SCENARIO OVDE:
	// 1. "normal"    - Svi su pošteni
	// 2. "silent"    - Lider ćuti (Testira Timeout i Slow View Change)
	// 3. "byzantine" - Lider šalje različite blokove (Testira Safety/Equivocation)
	// ==========================================================
	selectedScenario := "byzantine"

	fmt.Println("================================================================")
	fmt.Printf("  BEEGEES SIMULATOR - Scenario: %s\n", selectedScenario)
	fmt.Println("================================================================")

	allIDs := []process.ProcessID{"node-1", "node-2", "node-3", "node-4"} //skup svih replika u sistemu; U simulaciji imamo četiri čvora, što odgovara konfiguraciji 3f+1 za toleranciju jednog vizantijskog čvora.

	//Ovde se gradi komunikacioni sloj simulatora. Simulator implementira mrežu kroz Fair-Loss, Stubborn i Perfect Link slojeve, tako da BeeGees ne mora direktno da se bavi detaljima prenosa poruka.
	fl := link.NewFairLossLink(0.0, 42) //Prvi parametar je verovatnoća gubitka poruke (0.0-nema gubitka poruke, 42 je seed za generator slučajnosti.). Ne želim da mreža nasumično gubi poruke; želim da testiram ponašanje BeeGees algoritma
	sl := link.NewStubbornLink(fl)      //Ako poruka nije uspešno isporučena, sloj je ponovo šalje
	pl := link.NewPerfectLink(sl)       //apstrakciju pouzdanog kanala/komunikacije
	fm := failures.NewNoFailure()       //simulator neće nasumično rušiti čvorove, napravila sam svoje ponašanje čvora za potrebe eksperimenta.

	//Ovo podešava izvršavanje cele simulacije. Pokreni mi ove čvorove, poveži ih preko ove mreže, koristi ovaj failure model i izvršavaj simulaciju pod ovim uslovima.“
	rt := simrt.NewRuntime(pl, fm,
		simrt.WithMaxDuration(15*time.Second), //Simulator se najkasnije nakon 15 sekundi završava
		simrt.WithIdleTimeout(5*time.Second),  // Ako se pet sekundi ništa ne dešava (niko ne prica), simulacija se završava.
		// Usporavamo retransmit da ne bi spamovao terminal dok mi spavamo (time.Sleep)
		//simrt.WithRetransmitInterval(200*time.Millisecond),
		simrt.WithRetransmitInterval(1*time.Second), //Simulator periodično pokušava ponovo da pošalje poruke. To je posebno bitno zbog Stubborn Link-a.
		simrt.WithVerbose(false),                    //Ne prikazuje sve detalje mrežnog saobraćaja, da terminal ne bude zatrpan-ne prikazuje mrezne poruke
	)

	//kreiranje cvorova
	for _, pid := range allIDs {
		behavior := beegees.BehaviorHonest

		// Postavljamo specifično ponašanje zavisno od scenarija
		if selectedScenario == "silent" && pid == "node-2" {
			behavior = beegees.BehaviorSilent
		} else if selectedScenario == "byzantine" && pid == "node-3" {
			behavior = beegees.BehaviorByzantine
		}
		fmt.Printf("Registrujem %s [%s]\n", pid, behavior)

		node := beegees.NewBeeGeesNode(pid, allIDs, behavior) //kreiranje beegees cvora gde svaki cvor dobija svoj id, listu svih cvorova i svoje ponasanje
		rt.Register(node)                                     //registracija cvora u simulatoru
	}

	fmt.Printf(">>> Pokrećem simulaciju [%s]...\n", selectedScenario)
	rt.Run() //pokreće se cela distribuirana simulacija
	fmt.Println("================================================================")
}

/*
Korak 1 — definišeš čvorove
Korak 2 — napraviš mrežu
Korak 3 — napraviš runtime
Korak 4 — svakom čvoru odrediš ponašanje
Korak 5 — napraviš BeeGees čvor
Korak 6 — ubaciš cvor u simulator
Korak 7 — pokreneš sve
*/
