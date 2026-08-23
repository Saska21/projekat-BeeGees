package beegees

import (
	"context" //služi za kontrolu životnog ciklusa goroutine-a čvora, Simulacija je završena, prekini rad.
	"fmt"     //ispis u terminalu
	"time"    //za tajmere

	"github.com/danilokacanski/da/week03_04_parallel/process" //ID cvora
)

// Struktura čvora, jedan čvor/repliku u distribuiranom sistemu
/*
Čvorovi ne dele jednu zajedničku memoriju za stanje algoritma. Svaki čvor ima svoju lokalnu kopiju informacija koje je do tog trenutka saznao.
*/
type BeeGeesNode struct {
	id       process.ProcessID   //ID čvora
	allIDs   []process.ProcessID //Spisak svih čvorova u sistemu
	behavior string              //ponasanje cvora

	Bparent *Block //blok na koji će se novi blok nadovezati
	QCanc   *QC    //QC koji potvrđuje Banc
	Banc    *Block //poslednji bezbedno potvrđen blok

	//lokalno stanje
	currentView   int                            //U kom view-u se čvor trenutno nalazi
	lastVoteHash  string                         // hash poslednjeg bloka za koji je ovaj čvor glasao (koristi se u NewView poruci), lastVote iz Algoritma 2, koristi se i kod materijalizacijr
	lastVotedView int                            // OVO JE POTREBNO DA BI ZNALI da li je liveness tajmer istekao zbog neformiranja QC ili zbog nepredlaganja bloka od strane lidera
	blocks        map[string]*Block              //Lokalna baza poznatih blokova, potrebno za pronalazenje parenta, grandparenta, Banc, commit rule, materijalizaciju..
	votes         map[string][]process.ProcessID //kluc=hash bloka, a vrednost=lista čvorova koji su glasali za taj blok; Ovde lider skuplja glasove, Kada broj glasova dostigne kvorum → formira se QC
	newViews      map[int][]*NewViewMessage      //Ovde lider čuva pristigle NewView poruke po view-u-NVset, kako bi znao sta su replike poslednje videle i za sta su poslednje glasale-QC materijalizacija

	committedHeight int         //najveća visina bloka koju je ovaj čvor već komitovao, sprečavaš da isti blok ili neki stariji blok ponovo ispisuješ kao commit
	viewTimer       *time.Timer //Liveness timeout za trenutni view-ako istekne pokreće se Slow View Change i prelazi se u naredni view. Nema napretka → menjamo view
	matTimer        *time.Timer // MAT tajmer iz Algoritma 3, Sačekaj još malo pre nego što predložiš novi blok, možda stignu zakašnjeli VoteRESP i moći ćemo da materijalizujemo QC.
}

// Konstruktor-ovde se pravi novi cvor. Konstruktor pravi objekat čvora, postavlja početni View na 1, inicijalizuje lokalne mape i pokreće viewTimer na 3 sekunde
func NewBeeGeesNode(id process.ProcessID, allIDs []process.ProcessID, behavior string) *BeeGeesNode {
	return &BeeGeesNode{
		id:          id,                                   // Postavlja jedinstveni ID ovog čvora (npr. "node-1")
		allIDs:      allIDs,                               // Pamti spisak svih čvorova u mreži (potrebno za kvorum i izbor lidera)
		behavior:    behavior,                             //definise se ponasanje cvora
		currentView: 1,                                    // Svaki čvor započinje rad u View-u 1
		blocks:      make(map[string]*Block),              // Inicijalizuje lokalnu bazu blokova (hash -> Block)
		votes:       make(map[string][]process.ProcessID), // Inicijalizuje mapu za sakupljanje glasova (hash -> lista glasača)
		newViews:    make(map[int][]*NewViewMessage),      // Inicijalizuje mapu za NewView poruke po View-u
		viewTimer:   time.NewTimer(3 * time.Second),       // Glavni liveness tajmer - ako isteknu 3s bez napretka, ide View Change
		matTimer:    time.NewTimer(100 * time.Hour),       // Inicijalno "ugašen" dok ga cvor eksplicitno ne resetuje kada postane lider
	}
}

// vraca se id cvora
func (n *BeeGeesNode) ID() process.ProcessID { return n.id }

// glavna petlja cvora koja se izvršava kao posebna goroutine za svaki čvor. Time se ispunjava zahtev:“svaki čvor mora biti nezavisna izvršna jedinica.”
func (n *BeeGeesNode) Run(ctx context.Context, inbox <-chan process.Message, send func(process.Message)) {
	//Čvor koji je lider prvog view-a odmah predlaže prvi blok
	if n.isLeader(1) {
		n.propose(send, nil, nil) //qc=nil → genesis situacija; nvSet=nil → nema prethodnog view change-a
	}

	for { // Beskonačna petlja u kojoj čvor živi i čeka događaje
		select { //Čvor istovremeno čeka:gašenje simulatora, istek view timeout-a, istek MAT timeout-a, dolazak poruke. konkurentno izvršavanje
		case <-ctx.Done(): //gašenje simulatora
			return // Prekida se rad gorutine čvora

		case <-n.viewTimer.C: //istek view timeout-a
			// ALGORITAM 3, linija 6: View timer ističe
			n.startSlowViewChange(send)

		case <-n.matTimer.C: //istek MAT timeout-a
			// ALGORITAM 3, linija 31: MAT tajmer ističe
			n.onMatTimerExpired(send) // Prekida se rad gorutine čvora

		case msg := <-inbox: //dolazak poruke
			payload, ok := msg.Data.(BeeGeesMessage) //provera da li je stvarno BeeGeesMessage
			if !ok {
				continue //ako nije, ignoriše se poruka
			}
			n.handleMessage(payload, send) // Prosledi poruku
		}
	}
}

// rutiranje poruka
func (n *BeeGeesNode) handleMessage(m BeeGeesMessage, send func(process.Message)) {
	if m.View > n.currentView { //ako je poruka veceg view-a od trenutnog, pređi u njega. Omogućava brzu sinhronizaciju
		n.currentView = m.View
		n.viewTimer.Reset(3 * time.Second)
	}

	switch m.Kind { //proverava se tip BeeGees poruke, rutiranje poruke prema njenom tipu
	case KindVoteReq:
		n.handleVoteReq(m, send) // Lider predlaže blok (Proposal)
	case KindVoteResp:
		n.handleVoteResp(m, send) // Validator šalje glas (Vote)
	case KindNewView:
		n.handleNewView(m, send) // Izveštaj za Slow View Change
	}
}

// Obrada Predloga Bloka - ALGORITAM 2 i 4: Fast View Change i Commit Rule
func (n *BeeGeesNode) handleVoteReq(m BeeGeesMessage, send func(process.Message)) {
	// Alg 2, linija 14: Validacija predloga (B.v+1 = B'.v i B = B'.p)
	// Ovde smo malo pojednostavili uslov radi simulatora
	n.blocks[m.Block.Hash()] = m.Block //Čvor pamti primljeni blok
	n.viewTimer.Reset(3 * time.Second) // Resetuje tajmer jer je lider aktivan i poslao je blok

	// ALGORITAM 4: Commit Rule
	n.checkCommitRule(m.Block) // Čim vidi novi blok, proveri da li se neki stariji blok može komitovati

	if n.behavior == BehaviorSilent { //Čvor ne šalje glas: pao ili pokvaren
		return
	}

	// AZURIRANJE POSLEDNJEG GLASA (lastVoteHash)
	n.lastVoteHash = m.Block.Hash()
	n.lastVotedView = m.View // Pamti da se glasano u ovom View-u, bitno da bi znali posle uzrok liveness timeouta (ako smo glasali u trenutnom view onda je uzrok nekreiranje QC, a ako nismo onda je razlog to sto lider nije predlozio blok)
	fmt.Printf("[%s] Glasam za blok %s (H:%d)\n", n.id, n.lastVoteHash, m.Block.Height)

	// Alg 2, linija 18: Slanje glasa sledećem lideru
	nextLeader := n.getNextLeader(n.currentView + 1)                        //odredjivanje novog lidera za naredni view
	vote := process.NewMessage(n.id, nextLeader, "BEEGEES", BeeGeesMessage{ // Pravi poruku sa glasom (KindVoteResp)
		Kind: KindVoteResp, View: n.currentView, BlockHash: n.lastVoteHash, Sender: n.id,
	})
	send(vote) // Šalje glas direktno sledećem lideru
}

func (n *BeeGeesNode) handleVoteResp(m BeeGeesMessage, send func(process.Message)) {
	if !contains(n.votes[m.BlockHash], m.Sender) { //provera da se isti glas ne broji dva puta
		n.votes[m.BlockHash] = append(n.votes[m.BlockHash], m.Sender)
	}
	quorum := (len(n.allIDs) * 2 / 3) + 1

	if len(n.votes[m.BlockHash]) == quorum {
		if n.behavior != BehaviorSilent {
			n.currentView++
			qc := &QC{View: m.View, BlockHash: m.BlockHash, Signers: n.votes[m.BlockHash]} //Pravi novi Quorum Certificate (QC)
			fmt.Printf("[%s] FORMIRAN QC za blok %s \n", n.id, m.BlockHash)
			delete(n.votes, m.BlockHash) // Briše iskorišćene glasove iz memorije
			//n.currentView++
			n.propose(send, qc, nil) // Odmah predlaže novi blok sa tek kreiranim QC-om
		} else {
			// silent lider ne salje predlog novog bloka vec cuti
			fmt.Printf("  \n--- [VIEW %d] Lider [%s] je SILENT. Ima kvorum za QC, ali ga NEĆE OBJAVITI! Takodje, NE ŠALJE predlog!\n", n.currentView+1, n.id)
			// Ovde ne zovemo n.propose, pa on ostaje "zakucan"
		}
	}
}

// ALGORITAM 3: Slow View Change & Materialization
func (n *BeeGeesNode) handleNewView(m BeeGeesMessage, send func(process.Message)) {
	// AKO JE lider SILENT, NE RADI NIŠTA (kao da ne postoji)
	if n.behavior == BehaviorSilent {
		return
	}
	n.newViews[m.View] = append(n.newViews[m.View], m.NewViewData) //lider prikuplja NewView poruke u NVset
	quorum := (len(n.allIDs) * 2 / 3) + 1

	// Proveravamo uslove samo ako smo lider i ako imamo bar 2f+1 NewView poruka
	if n.isLeader(m.View) && len(n.newViews[m.View]) >= quorum {

		// Linija 19: Određujemo Bparent (najviši predloženi blok iz NewView poruka)
		n.Bparent = n.highVoteReq(n.newViews[m.View])

		// Linije 20-21: Inicijalni Banc i QCanc
		if n.Bparent != nil && n.Bparent.QCanc != nil {
			n.QCanc = n.Bparent.QCanc
			n.Banc = n.blocks[n.QCanc.BlockHash]
		} else {
			n.QCanc = nil
			n.Banc = nil
		}

		// --- MATERIJALIZACIJA (Linije 22-25) ---
		// Proveravamo da li u NVsetu postoji n-f glasova za neki blok koji je potomak Banc-a
		voteCounts := make(map[string]int)
		for _, nv := range n.newViews[m.View] {
			if nv.LastVoteResp != "" {
				voteCounts[nv.LastVoteResp]++ //Brojiš koliko replika je glasalo za koji blok.
			}
		}

		for blockHash, count := range voteCounts {
			// Linija 23: Ako imamo n-f glasova za neki blok
			if count >= quorum {
				potentialBanc := n.blocks[blockHash]
				// Proveravamo da li taj blok "ispunjava uslove" da postane novi Banc
				if potentialBanc != nil && (n.Banc == nil || n.isDescendantOf(potentialBanc, n.Banc)) && (n.Banc == nil || potentialBanc.Height > n.Banc.Height) {
					fmt.Printf("[%s] MATERIJALIZACIJA: Stvaram novi QC za blok %s\n", n.id, blockHash)

					signers := []process.ProcessID{} //glasaci za materijalizovani QC bloka
					for _, nv := range n.newViews[m.View] {
						if nv.LastVoteResp == blockHash {
							signers = append(signers, nv.ReplicaID)
						}
					}

					// Linija 24: CreateQC(NVset) - stvaramo novi QC
					n.QCanc = &QC{
						View:      m.View - 1,
						BlockHash: blockHash,
						Signers:   signers,
					}
					// Linija 25: Ažuriramo Banc
					n.Banc = potentialBanc
				}
			}
		}

		// Linija 26: Ako je Banc postao jednak Bparent, možemo odmah da šaljemo predlog novog bloka jer nema vise blokova za koje postoji potencijalna materijalizacija
		if n.Banc != nil && n.Bparent != nil && n.Banc.Hash() == n.Bparent.Hash() {
			fmt.Printf("[%s] Banc == Bparent! Prekidam MAT i šaljem blok.\n", n.id)
			n.matTimer.Stop()
			n.propose(send, n.QCanc, n.newViews[m.View]) //predlog novog bloka pri cemu se prosledjuje i QCanc i sakupljene NewView poruke za trenutni view

			//Čistimo sakupljene NewView poruke kako ne bi doslo do mesanja starih i novih glasova u sledecoj rundi
			delete(n.newViews, m.View) //DODATO RESETOVANJE SAKUPLJENIH GLASOVA
			//Da ne bi ostale stare reference za sledeći view
			n.Bparent = nil
			n.Banc = nil
			n.QCanc = nil

			return
		}
	}
}

// 3. Šta se dešava kad MAT ipak istekne (Linije 31-33)
// Ovo se dešava u Run() petlji u "case <-n.matTimer.C:"
// Lider šalje blok sa onim što je uspeo da skupi do tog trenutka
func (n *BeeGeesNode) onMatTimerExpired(send func(process.Message)) {
	fmt.Printf("[%s] MAT istekao. Šaljem blok sa trenutnim QCanc.\n", n.id)
	// Linije 32-33: CreateBlock i Send Vote-Req
	n.propose(send, n.QCanc, n.newViews[n.currentView])

	delete(n.newViews, n.currentView)
	n.Bparent = nil
	n.Banc = nil
	n.QCanc = nil
}

// ALGORITAM 4: Commit Rule
func (n *BeeGeesNode) checkCommitRule(b *Block) {
	// DODAJ OVO: Ako je čvor silent, on ne procesira čak ni commit
	if n.behavior == BehaviorSilent || n.behavior == BehaviorByzantine {
		return
	}

	/*
		proveravamo da li:
		-b ima QC,
		-roditelj ima QC,
		-postoji grandparent.
	*/
	if b.QCanc == nil {
		return
	}
	parent, okP := n.blocks[b.QCanc.BlockHash]
	if !okP || parent.QCanc == nil { //ako roditeljski blok ne postoji u lokalnoj bazi, ili postoji, ali nema QC, onda ne možemo dalje da proveravamo commit.
		return
	}

	grandParent, okG := n.blocks[parent.QCanc.BlockHash]
	if !okG {
		return
	}

	// Tačka 1: Uzastopni QC? (Algorithm 4, linija 6)
	if parent.View == grandParent.View+1 {
		if grandParent.Height > n.committedHeight {
			n.committedHeight = grandParent.Height
			fmt.Printf("[%s] COMMIT (CHL): Uzastopni QC potvrđuju %s (Visina %d) \n", n.id, grandParent.Hash(), grandParent.Height)
		}
	} else {
		// Tačka 2: Nisu uzastopni (AHL). Provera Equivocation (FindEquivProof)
		fmt.Printf("[%s] AHL Provera (Rupa između View %d i %d)...\n", n.id, grandParent.View, parent.View)
		if n.findEquivocationProof(grandParent, parent, parent.NVset) {
			fmt.Printf("[%s] ABORT: Otkriven dokaz o prevari! Blok %s (visina %d) se NE komituje.\n", n.id, grandParent.Hash(), grandParent.Height)
		} else {
			if grandParent.Height > n.committedHeight {
				n.committedHeight = grandParent.Height
				fmt.Printf("[%s] COMMIT (AHL): Nema dokaza o prevari, komitujem %s (Visina %d)\n", n.id, grandParent.Hash(), grandParent.Height)
			}
		}
	}
}

// --- POMOĆNE METODE (Logic) ---

// Određuje Bparent (najviši predloženi blok iz NewView poruka)
func (n *BeeGeesNode) highVoteReq(set []*NewViewMessage) *Block {
	// Implementacija Algorithm 3, linija 45
	var highest *Block

	for _, nv := range set { //Čim vidi da je LastVoteReq prazan (nil), preskačemo ga
		if nv.LastVoteReq == nil {
			continue
		}

		if highest == nil || nv.LastVoteReq.View > highest.View {
			highest = nv.LastVoteReq
		}
	}

	return highest
}

func (n *BeeGeesNode) findEquivocationProof(start, end *Block, nvSet []*NewViewMessage) bool {
	// Implementacija Algorithm 4, linija 26: FindEquivProof
	fmt.Printf(" NVset sadrži:\n")
	for _, nv := range nvSet {
		if nv.LastVoteReq != nil {
			fmt.Printf("     replika %s → view %d → %s\n",
				nv.ReplicaID,
				nv.LastVoteReq.View,
				nv.LastVoteReq.Hash())
		} else {
			fmt.Printf("     replika %s → nema LastVoteReq\n", nv.ReplicaID)
		}
	}
	// Tražimo dva različita bloka u istom View-u unutar NVset-a
	viewsSeen := make(map[int]string) //pravi se mapa view -> hash bloka (npr. view 3: "H1_V3_A")
	for _, nv := range nvSet {
		if nv.LastVoteReq != nil { //Neke replike možda nisu primile nijedan VoteREQ, pa proveravaš da nije nil
			v := nv.LastVoteReq.View //U kom view-u je predlog napravljen
			h := nv.LastVoteReq.Hash()
			if existingHash, ok := viewsSeen[v]; ok && existingHash != h { //Da li smo već videli neki blok za isti view
				fmt.Printf("Equivocation dokaz: view %d ima blokove %s i %s\n", v, existingHash, h)
				return true // Pronađena dva različita bloka za isti View!
			}
			viewsSeen[v] = h //Ako prvi put vidimo taj view, zapamtimo hash
		}
	}
	return false //nema dokaza o prevari
}

// provera da li je neki cvor lider
func (n *BeeGeesNode) isLeader(view int) bool { return n.allIDs[(view-1)%len(n.allIDs)] == n.id }

// Standardna Round-Robin rotacija lidera
func (n *BeeGeesNode) getNextLeader(view int) process.ProcessID {
	return n.allIDs[(view-1)%len(n.allIDs)]
}

// za određivanje visine u lancu bloka
func (n *BeeGeesNode) getNewHeight(qc *QC) int {
	if qc == nil {
		return 1 //na pocetku visina je 1
	}
	if b, ok := n.blocks[qc.BlockHash]; ok {
		return b.Height + 1 // Visina novog bloka je za 1 veća od bloka na koji pokazuje QC
	}
	return 1
}

// za liderovo predlaganje novog bloka
func (n *BeeGeesNode) propose(send func(process.Message), qc *QC, nvSet []*NewViewMessage) {
	//Kako bi videli svaki korak dodajemo pauzu
	time.Sleep(800 * time.Millisecond) // Pauza od 0.8 sekundi pre svakog novog predloga

	parentHash := ""
	if qc != nil {
		parentHash = qc.BlockHash
	}

	newBlock := &Block{View: n.currentView, Parent: parentHash, Height: n.getNewHeight(qc), QCanc: qc, NVset: nvSet, Data: fmt.Sprintf("TX_VIEW_%d", n.currentView)}
	if n.behavior == BehaviorByzantine {
		fmt.Printf("  \n--- [VIEW %d] Lider [%s] šalje KONFLIKTNE blokove!\n", n.currentView, n.id)
		n.sendByzantineProposals(newBlock, send)
		return
	}
	if n.behavior == BehaviorSilent {
		fmt.Printf("  \n--- [VIEW %d] Lider [%s] je SILENT i NE ŠALJE predlog!\n", n.currentView, n.id)
		return
	}
	fmt.Printf("\n--- [VIEW %d] LIDER %s predlaže blok %s (H:%d) ---\n", n.currentView, n.id, newBlock.Hash(), newBlock.Height)
	for _, rid := range n.allIDs {
		send(process.NewMessage(n.id, rid, "BEEGEES", BeeGeesMessage{Kind: KindVoteReq, View: n.currentView, Block: newBlock, Sender: n.id}))
	}
}

// poziva se kada istekne liveness tajmer
func (n *BeeGeesNode) startSlowViewChange(send func(process.Message)) {
	if n.behavior == BehaviorSilent {
		// Ako je čvor silent, on ne radi nista
		return
	}

	relog := "Lider nije poslao predlog"
	if n.lastVotedView == n.currentView {
		relog = "QC nije formiran na vreme"
	}

	// Ispisujemo timeout
	fmt.Printf("[%s] Liveness Timeout! (%s). Pokrećem Slow View Change...\n", n.id, relog)

	// MALA PAUZA (100ms) da se svi timeout logovi u terminalu slegnu
	time.Sleep(100 * time.Millisecond)

	n.currentView++
	n.viewTimer.Reset(3 * time.Second)           //povećava se view
	nextLeader := n.getNextLeader(n.currentView) //bira se novi lider

	// AKO JE LIDER, PRVO NAJAVI MAT, PA ONDA ŠALJI PORUKE
	if n.isLeader(n.currentView) {
		fmt.Printf("[%s] Ja sam lider View-a %d. Pokrećem MAT tajmer...\n", n.id, n.currentView)
		n.matTimer.Reset(2 * time.Second)
	}

	// Sastavljanje NewView poruke sa najsvežijim lokalnim informacijama
	nv := BeeGeesMessage{Kind: KindNewView, View: n.currentView, Sender: n.id,
		NewViewData: &NewViewMessage{
			ReplicaID:    n.id,
			View:         n.currentView,
			LastVoteReq:  n.blocks[n.lastVoteHash], // Poslednji predloženi blok za koji smo glasali
			LastVoteResp: n.lastVoteHash,
			LastQC:       n.getHighestKnownQC(), // Najviši viđeni QC
		}} //šalje se NEW_VIEW
	// šaljemo NewView poruke lideru
	send(process.NewMessage(n.id, nextLeader, "BEEGEES", nv))
}

func (n *BeeGeesNode) sendByzantineProposals(b *Block, send func(process.Message)) {
	for i, rid := range n.allIDs {
		altBlock := *b
		if i >= len(n.allIDs)/2 { // Prvoj polovini čvorova šalje originalni podatak, drugoj polovini šalje "FORK"
			altBlock.Data = "FORK"
		}
		send(process.NewMessage(n.id, rid, "BEEGEES", BeeGeesMessage{Kind: KindVoteReq, View: n.currentView, Block: &altBlock, Sender: n.id}))
	}
}

// za proveru postojanja identifikatora u nizu čvorova
func contains(ids []process.ProcessID, id process.ProcessID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// Koristi se prilikom formiranja NewViewMessage kako bi se novi lider obavestio o najnovijem potvrdnom stanju u mreži
func (n *BeeGeesNode) getHighestKnownQC() *QC {
	var highest *QC

	for _, b := range n.blocks {
		if b.QCanc != nil {
			if highest == nil || b.QCanc.View > highest.View {
				highest = b.QCanc
			}
		}
	}

	return highest
}

// Iterativno se kreće unazad od potencijalnog potomka (child), prateći pokazivače na roditeljske blokove (cur.Parent), proveravajući da li se na toj stazi nalazi predak (ancestor)
func (n *BeeGeesNode) isDescendantOf(child *Block, ancestor *Block) bool {
	if child == nil || ancestor == nil {
		return false
	}

	cur := child
	for cur != nil {
		if cur.Hash() == ancestor.Hash() {
			return true // Uspešno smo rekreirali stazu od child-a do ancestor-a
		}

		if cur.Parent == "" {
			break // Stigli do korena lanca
		}

		cur = n.blocks[cur.Parent] // Idemo korak unazad preko roditeljskog heša
	}

	return false
}
