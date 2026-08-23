package beegees

import (
	"context" //služi za kontrolu životnog ciklusa goroutine-a čvora, Simulacija je završena, prekini rad.
	"fmt"
	"time"

	"github.com/danilokacanski/da/week03_04_parallel/process"
)

// Struktura čvora
type BeeGeesNode struct {
	id       process.ProcessID   //ID čvora
	allIDs   []process.ProcessID //Spisak svih čvorova u sistemu
	behavior string              //scenario

	Bparent *Block //blok na koji će se novi blok nadovezati
	QCanc   *QC    //QC koji potvrđuje Banc
	Banc    *Block //poslednji bezbedno potvrđen blok

	//lokalno stanje
	currentView   int                            //U kom view-u se čvor trenutno nalazi
	lastVoteHash  string                         // blok za koji je cvor poslednji put glasao (koristi se u NewView poruci), lastVote iz Algoritma 2
	lastVotedView int                            // OVO JE POTREBNO DA BI ZNALI da li je liveness tajmer istekao zbog ne formiranja QC ili zbog ne predlaganja bloka od strane lidera
	blocks        map[string]*Block              //Lokalna baza poznatih blokova
	votes         map[string][]process.ProcessID //Ovde lider skuplja glasove, Kada broj glasova dostigne kvorum → formira se QC
	newViews      map[int][]*NewViewMessage      //Ovde lider čuva pristigle NewView poruke po view-u-NVset

	committedHeight int
	viewTimer       *time.Timer //Liveness timeout-ako istekne pokreće se Slow View Change
	matTimer        *time.Timer // MAT tajmer iz Algoritma 3, Sačekaj još malo pre nego što predložiš novi blok, možda stignu zakašnjeli VoteRESP i moći ćemo da materijalizujemo QC.

	blacklisted       map[process.ProcessID]bool //UNAPREDJENJE, sluzi za belezenje zlonamernih cvorova, na osnovu kog se kasnije računa ko sme, a ko ne sme da bude lider
	lastCommittedHash string                     //UNAPREDJENJE, pamti hash poslednjeg uspešno komitovanog bloka

	epochs []rotationEpoch // NOVO: istorija segmenata rotacije, Svaki segment kaže: "od view-a X nadalje, lideri se biraju kružno kroz ovu listu aktivnih čvorova, počev od ove pozicije
}

// DODATO KAKO BI SE CVOROVI ISPRAVNO SMENJIVALI NAKON BLACKLISTOVANJA
// rotationEpoch opisuje jedan "segment" rotacije lidera koji važi
// od view-a fromView nadalje, dok se ne dogodi sledeći blacklist.
// Svaki put kada se otkrije vizantijski čvor, kreira se nova epoha
type rotationEpoch struct {
	fromView int                 // od kog view-a važi ovaj segment
	active   []process.ProcessID // aktivni čvorovi u ovom segmentu (fiksni redosled)
	offset   int                 // indeks u 'active' koji odgovara lideru za view=fromView, rotacija nastavlja glatko, bez skoka
}

// Konstruktor-ovde se pravi novi cvor
func NewBeeGeesNode(id process.ProcessID, allIDs []process.ProcessID, behavior string) *BeeGeesNode {
	n := &BeeGeesNode{
		id:          id,
		allIDs:      allIDs,
		behavior:    behavior,
		currentView: 1,
		blocks:      make(map[string]*Block),
		votes:       make(map[string][]process.ProcessID),
		newViews:    make(map[int][]*NewViewMessage),
		viewTimer:   time.NewTimer(3 * time.Second),
		matTimer:    time.NewTimer(100 * time.Hour), // Inicijalno "ugašen" dok ga lider eksplicitno ne resetuje

		blacklisted: make(map[process.ProcessID]bool), //UNAPREDJENJE
	}

	activeCopy := make([]process.ProcessID, len(allIDs)) //pravi se kopija allIDs da ne bi menjali direktno allIDs vec da bi kopiju menjali
	copy(activeCopy, allIDs)
	n.epochs = []rotationEpoch{{fromView: 1, active: activeCopy, offset: 0}} //u poecetnoj epohi svi cvorovi su aktivni i krecu iz view-a 1 i u view 1 je lider sa indeksom 0 odnosno node 1

	return n
}

// vraca se id cvora
func (n *BeeGeesNode) ID() process.ProcessID { return n.id }

// glavna petlja cvora koja se izvršava kao posebna goroutine za svaki čvor.svaki čvor mora biti nezavisna izvršna jedinica
func (n *BeeGeesNode) Run(ctx context.Context, inbox <-chan process.Message, send func(process.Message)) {
	//Čvor koji je lider prvog view-a odmah predlaže prvi blok
	if n.isLeader(1) {
		n.propose(send, nil, nil) //qc=nil → genesis situacija; nvSet=nil → nema prethodnog view change-a
	}

	for {
		select { //Čvor istovremeno čeka:gašenje simulatora, istek view timeout-a, istek MAT timeout-a, dolazak poruke. Ovo je baš ono što su tražili pod “konkurentno izvršavanje”.
		case <-ctx.Done(): //gašenje simulatora,
			return

		case <-n.viewTimer.C: //istek view timeout-a,
			// ALGORITAM 3, linija 6: View timer ističe
			n.startSlowViewChange(send)

		case <-n.matTimer.C: //istek MAT timeout-a,
			// ALGORITAM 3, linija 31: MAT tajmer ističe
			n.onMatTimerExpired(send)

		case msg := <-inbox: //dolazak poruke.
			payload, ok := msg.Data.(BeeGeesMessage) //proveravaš da li je stvarno BeeGeesMessage
			if !ok {
				continue //ako nije, ignoriše se poruka
			}
			n.handleMessage(payload, send)
		}
	}
}

func (n *BeeGeesNode) handleMessage(m BeeGeesMessage, send func(process.Message)) {
	if m.View > n.currentView { //ako vidiš viši view (ako je poruka veceg view-a), pređi u njega. važno za sinhronizaciju.
		n.currentView = m.View
		n.viewTimer.Reset(3 * time.Second)
	}

	switch m.Kind { //proverava se tip BeeGees poruke
	case KindVoteReq:
		n.handleVoteReq(m, send)
	case KindVoteResp:
		n.handleVoteResp(m, send)
	case KindNewView:
		n.handleNewView(m, send)
	}
}

// ALGORITAM 2 i 4: Fast View Change i Commit Rule
func (n *BeeGeesNode) handleVoteReq(m BeeGeesMessage, send func(process.Message)) {
	// Alg 2, linija 14: Validacija predloga (B.v+1 = B'.v i B = B'.p)
	// Ovde smo malo pojednostavili uslov radi simulatora
	n.blocks[m.Block.Hash()] = m.Block //Čvor pamti predlozeni blok
	n.viewTimer.Reset(3 * time.Second)

	// ALGORITAM 4: Commit Rule
	n.checkCommitRule(m.Block) //Commit proveravaš čim vidiš novi blok. To odgovara BeeGees logici gde novi QC može omogućiti commit ranijeg bloka.

	if n.behavior == BehaviorSilent { //Čvor ne šalje glas: crash ili neodgovaranje lidera ili vizantijski cvor
		return
	}

	// AZURIRANJE POSLEDNJEG GLASA (lastVoteHash)
	n.lastVoteHash = m.Block.Hash()
	n.lastVotedView = m.View // Pamti se da je cvor glasao u ovom View-u, bitno da bi znali posle uzrok liveness timeouta (ako smo glasali u trenutnom view onda je uzrok nekreiranje QC, a ako nismo onda je razlog to sto lider nije predlozio blok)
	fmt.Printf("  [%s] Glasam za blok %s (H:%d)\n", n.id, n.lastVoteHash, m.Block.Height)

	// Alg 2, linija 18: Slanje glasa sledećem lideru
	nextLeader := n.getNextLeader(n.currentView + 1)
	vote := process.NewMessage(n.id, nextLeader, "BEEGEES", BeeGeesMessage{
		Kind: KindVoteResp, View: n.currentView, BlockHash: n.lastVoteHash, Sender: n.id,
	})
	send(vote)
}

func (n *BeeGeesNode) handleVoteResp(m BeeGeesMessage, send func(process.Message)) {
	if !contains(n.votes[m.BlockHash], m.Sender) { //provera da se isti glas ne broji dva puta
		n.votes[m.BlockHash] = append(n.votes[m.BlockHash], m.Sender)
	}
	quorum := (len(n.allIDs) * 2 / 3) + 1
	if len(n.votes[m.BlockHash]) == quorum {
		if n.behavior != BehaviorSilent {
			//n.currentView++
			qc := &QC{View: m.View, BlockHash: m.BlockHash, Signers: n.votes[m.BlockHash]}
			fmt.Printf("[%s] FORMIRAN QC za blok %s \n", n.id, m.BlockHash) //u ispisu se navodi view u kom je blok predlozen(m.View), a ne view u kom se nalazimo
			delete(n.votes, m.BlockHash)
			n.currentView++
			n.propose(send, qc, nil)
		} else {
			//silent lider ne salje predlog novog bloka
			fmt.Printf("\n--- [VIEW %d] Lider [%s] je SILENT. Ima kvorum za QC, ali ga NEĆE OBJAVITI! Takodje, NE ŠALJE predlog!\n", n.currentView+1, n.id)
			// Ovde ne zovemo n.propose, pa on ostaje "zakucan"
		}
	}
}

// ALGORITAM 3: Slow View Change & Materialization
func (n *BeeGeesNode) handleNewView(m BeeGeesMessage, send func(process.Message)) {
	// AKO JE ČVOR SILENT, NE RADI NIŠTA (kao da ne postoji)
	if n.behavior == BehaviorSilent {
		return
	}
	n.newViews[m.View] = append(n.newViews[m.View], m.NewViewData) //prikupljamo NewView poruke u NVset
	quorum := (len(n.allIDs) * 2 / 3) + 1

	// Proveravamo uslove samo ako smo lider i ako imamo bar 2f+1 poruka
	if n.isLeader(m.View) && len(n.newViews[m.View]) >= quorum {

		// Linija 19: Određujemo Bparent (najviši predloženi blok)
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

					signers := []process.ProcessID{}
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

		// Linija 26: Ako je Banc postao jednak Bparent, možemo odmah da šaljemo
		if n.Banc != nil && n.Bparent != nil && n.Banc.Hash() == n.Bparent.Hash() {
			fmt.Printf("[%s] Banc == Bparent! Prekidam MAT i šaljem blok.\n", n.id)
			n.matTimer.Stop()
			n.propose(send, n.QCanc, n.newViews[m.View])

			//Čistimo NewView poruke za ovaj View
			delete(n.newViews, m.View) //DODATO RESETOVANJE SAKUPLJENIH GLASOVA
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

// ALGORITAM 4: Commit Rule (Sve stavke koje si navela)
func (n *BeeGeesNode) checkCommitRule(b *Block) {
	// Ako je čvor silent, on ne procesira čak ni commit
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
			n.lastCommittedHash = grandParent.Hash() // UNAPREDJENJE
			fmt.Printf(" [%s] COMMIT (CHL): Uzastopni QC potvrđuju %s (Visina %d) \n", n.id, grandParent.Hash(), grandParent.Height)
		}
	} else {
		// Tačka 2: Nisu uzastopni (AHL). Provera Equivocation (FindEquivProof)
		fmt.Printf("[%s] AHL Provera (Rupa između View %d i %d)...\n", n.id, grandParent.View, parent.View)
		if n.findEquivocationProof(grandParent, parent, parent.NVset) {
			// View u kome je nastala equivocacija
			equivView := parent.View - 1 //ekvivokacije se dogodila u view pre parenta
			n.blacklistLeader(equivView)
			fmt.Printf("[%s] ABORT: Otkriven dokaz o prevari! Blok %s (visina %d) se NE komituje.\n", n.id, grandParent.Hash(), grandParent.Height)
		} else {
			if grandParent.Height > n.committedHeight {
				n.committedHeight = grandParent.Height
				n.lastCommittedHash = grandParent.Hash() // UNAPREDJENJE
				fmt.Printf("[%s] COMMIT (AHL): Nema dokaza o prevari, komitujem %s (Visina %d)\n", n.id, grandParent.Hash(), grandParent.Height)
			}
		}
	}
}

// --- POMOĆNE METODE (Logic) ---

func (n *BeeGeesNode) highVoteReq(set []*NewViewMessage) *Block {
	// Implementacija Algorithm 3, linija 45
	var highest *Block

	for _, nv := range set { //Čim vidi da je LastVoteReq prazan (nil), ona ga preskače pomoću continue
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
	fmt.Printf("NVset sadrži:\n")
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
			if existingHash, ok := viewsSeen[v]; ok && existingHash != h {
				fmt.Printf("Equivocation dokaz: view %d ima blokove %s i %s\n", v, existingHash, h)
				return true // Pronađena dva različita bloka za isti View!
			}
			viewsSeen[v] = h //Ako prvi put vidimo taj view, zapamtimo hash.
		}
	}
	return false

}

// UNAPREDJENJE
func (n *BeeGeesNode) isLeader(view int) bool {
	return n.getRotationLeader(view) == n.id
}

func (n *BeeGeesNode) getNextLeader(view int) process.ProcessID {
	return n.getRotationLeader(view)
}
func (n *BeeGeesNode) getNewHeight(qc *QC) int {
	// Ako nema QC, kreni od poslednjeg komitovanog bloka
	if qc == nil {
		return n.committedHeight + 1
	}
	if b, ok := n.blocks[qc.BlockHash]; ok {
		return b.Height + 1
	}
	return n.committedHeight + 1
}
func (n *BeeGeesNode) propose(send func(process.Message), qc *QC, nvSet []*NewViewMessage) {
	//Kako bi videli svaki korak dodajemo pauzu
	time.Sleep(800 * time.Millisecond) // Pauza od 0.8 sekundi pre svakog novog predloga

	parentHash := ""
	if qc != nil {
		parentHash = qc.BlockHash
	}

	newBlock := &Block{View: n.currentView, Parent: parentHash, Height: n.getNewHeight(qc), QCanc: qc, NVset: nvSet, Data: fmt.Sprintf("TX_VIEW_%d", n.currentView)}
	if n.behavior == BehaviorByzantine {
		fmt.Printf("\n--- [VIEW %d] Lider [%s] šalje KONFLIKTNE blokove!\n", n.currentView, n.id)
		n.sendByzantineProposals(newBlock, send)
		return
	}
	if n.behavior == BehaviorSilent {
		fmt.Printf("\n--- [VIEW %d] Lider [%s] je SILENT i NE ŠALJE predlog!\n", n.currentView, n.id)
		return
	}
	fmt.Printf("\n--- [VIEW %d] LIDER %s predlaže blok %s (H:%d) ---\n", n.currentView, n.id, newBlock.Hash(), newBlock.Height)
	for _, rid := range n.allIDs {
		send(process.NewMessage(n.id, rid, "BEEGEES", BeeGeesMessage{Kind: KindVoteReq, View: n.currentView, Block: newBlock, Sender: n.id}))
	}
}
func (n *BeeGeesNode) startSlowViewChange(send func(process.Message)) {
	if n.behavior == BehaviorSilent {
		// Ako je čvor silent, on ne šalje ni NewView poruke
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

	//nv := BeeGeesMessage{Kind: KindNewView, View: n.currentView, Sender: n.id, NewViewData: &NewViewMessage{ReplicaID: n.id, View: n.currentView, LastVoteResp: n.lastVoteHash}} //šalje se NEW_VIEW
	nv := BeeGeesMessage{Kind: KindNewView, View: n.currentView, Sender: n.id,
		NewViewData: &NewViewMessage{
			ReplicaID:    n.id,
			View:         n.currentView,
			LastVoteReq:  n.blocks[n.lastVoteHash], // ako postoji
			LastVoteResp: n.lastVoteHash,
			LastQC:       n.getHighestKnownQC(),
		}} //šalje se NEW_VIEW
	// validatori salju lideru new view poruke
	send(process.NewMessage(n.id, nextLeader, "BEEGEES", nv))
}

func (n *BeeGeesNode) sendByzantineProposals(b *Block, send func(process.Message)) {
	for i, rid := range n.allIDs {
		altBlock := *b
		if i >= len(n.allIDs)/2 {
			altBlock.Data = "FORK"
		}
		send(process.NewMessage(n.id, rid, "BEEGEES", BeeGeesMessage{Kind: KindVoteReq, View: n.currentView, Block: &altBlock, Sender: n.id}))
	}
}

func contains(ids []process.ProcessID, id process.ProcessID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// imzena u startSlowViewChange
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

// proverava da li se child blok nastavlja na pretka
func (n *BeeGeesNode) isDescendantOf(child *Block, ancestor *Block) bool {
	if child == nil || ancestor == nil {
		return false
	}

	cur := child
	for cur != nil {
		if cur.Hash() == ancestor.Hash() {
			return true
		}

		if cur.Parent == "" {
			break
		}

		cur = n.blocks[cur.Parent]
	}

	return false
}

//UNAPREDJENJE!!!
/*
Ovo:
-izračuna ko je bio lider tog view-a,
-doda ga u blacklist,
-ispiše poruku (samo prvi put).
*/
func (n *BeeGeesNode) blacklistLeader(view int) {
	leader := n.getLeader(view) // ko je vizantijski cvor
	if leader == "" {
		return
	}
	if n.blacklisted[leader] {
		return // već obrađeno, ne diramo epoch listu ponovo
	}
	n.blacklisted[leader] = true //ako nije obradjeno sada ga dodajemo u blacklistu

	// Lider po TRENUTNO važećem segmentu (potrenutnim pravilima), za view koji je upravo u toku —
	// odatle nastavljamo rotaciju, bez ikakvog "skoka" ili ponavljanja.
	lastLeader := n.getRotationLeader(n.currentView) //ko je trenutni lider

	newActive := n.activeIDs()                             // sada već isključuje i tek blacklistovanog
	nextLeader := n.nextActiveAfter(lastLeader, newActive) //ko je naredni lider nove epohe

	offset := 0
	for i, x := range newActive {
		if x == nextLeader {
			offset = i //i=pozicija nextLeader-a u newActive
			break
		}
	}

	n.epochs = append(n.epochs, rotationEpoch{
		fromView: n.currentView + 1,
		active:   newActive,
		offset:   offset,
	})

	fmt.Printf("[%s] Lider %s iz View-a %d je BLACKLISTOVAN zbog equivocacije!\n",
		n.id, leader, view)
}

// UNAPREDJENJE-funkcija koja vraća aktivne čvorove
func (n *BeeGeesNode) activeIDs() []process.ProcessID {
	active := []process.ProcessID{}
	for _, id := range n.allIDs {
		if !n.blacklisted[id] {
			active = append(active, id)
		}
	}
	return active
}

// odredjuje ko je zaista predlozio konfliktne blokove od 4 cvora na pocetku
func (n *BeeGeesNode) getLeader(view int) process.ProcessID {
	return n.allIDs[(view-1)%len(n.allIDs)]
}

// currentEpochFor pronalazi poslednji segment rotacije čiji fromView <= view, pronalazi kojoj epohi pripada neki view
func (n *BeeGeesNode) currentEpochFor(view int) rotationEpoch {
	best := n.epochs[0]
	for _, e := range n.epochs {
		if e.fromView <= view {
			best = e
		} else {
			break
		}
	}
	return best
}

// Ko je stvarni lider view-a X, uzimajući u obzir sve dosadašnje blacklist-eve, koristi se kada treba da znamo ko je lider
// ko je lider za view-a broj X", i to računa preko formule i podataka koje je nextActiveAfter već jednom postavio (offset).
func (n *BeeGeesNode) getRotationLeader(view int) process.ProcessID {
	e := n.currentEpochFor(view)
	if len(e.active) == 0 {
		return ""
	}
	idx := (view - e.fromView + e.offset) % len(e.active)
	return e.active[idx]
}

// U originalnom kružnom poretku allIDs, ko dolazi odmah posle čvora id, preskačući blacklistovane
// Koristi se samo jednom da "zasadi" početnu tačku novog segmenta.
// o dolazi odmah posle node-1, ako preskačem sve koji su na blacklist-i, ne gleda uopste koji je view u pitanju vec samo prati redosled cvorova
// poziva se samo po jednom kad se blacklistovanje, a u svim ostalim slucajevima se koristi getRotationLeader
func (n *BeeGeesNode) nextActiveAfter(id process.ProcessID, active []process.ProcessID) process.ProcessID {
	total := len(n.allIDs)
	startIdx := 0
	for i, x := range n.allIDs {
		if x == id {
			startIdx = i //nalazimo poziciju u allIDs na kojoj se nalazi neki cvor (poslednji lider)
			break
		}
	}
	for i := 1; i <= total; i++ {
		cand := n.allIDs[(startIdx+i)%total] //kandidat je prvi sledeci aktivan cvor krenuvsi od cvora sa id=startIdx
		if contains(active, cand) {          //ako je blacklistovan nece nici u active
			return cand
		}
	}
	return ""
}
