// strukture podataka kojima BeeGees predstavlja blokove, QC i poruke
package beegees //Svi fajlovi koji implementiraju algoritam biće u paketu beegees

import (
	"fmt" //ispisivanje informacija u terminal //koristiš za generisanje hash stringa

	"github.com/danilokacanski/da/week03_04_parallel/process" //process.ProcessID predstavlja identitet čvora, koristimo tip koji je predvidjen za simulator
)

// Tipovi poruka (Section 4.2)
const (
	KindVoteReq  = "VOTE_REQ"  // Liderov predlog bloka validatorima (Proposal)
	KindVoteResp = "VOTE_RESP" // Validatoorov glas za blok koji se salje narednom lideru(Vote)
	KindNewView  = "NEW_VIEW"  // Izveštaji validatora novom lideru (koristi se u Slow View Change)
)

// Tipovi ponašanja (za simulaciju grešaka)
const (
	BehaviorHonest    = "honest"    //normalan rad
	BehaviorSilent    = "silent"    //lider ne šalje poruke
	BehaviorByzantine = "byzantine" //lider šalje različite blokove
)

// Block struktura (Algorithm 1)
type Block struct {
	View   int               //view-u u kom je blok predložen
	Parent string            //Hash roditeljskog bloka, blok ne čuva direktno ceo roditeljski blok, već njegov identifikator
	QCanc  *QC               //blok nosi QC svog relevantnog prethodnika (omogucava proveru commit pravila,nastavak lanca posle view change-a)
	Txs    []string          //batch transakcija koje nosi blok
	NVset  []*NewViewMessage //Ovo je skup NewView poruka
	Height int               //Visina bloka u lancu
	Data   string            // DODATO: Ovo je nedostajalo za Byzantine simulaciju kako bi lider mogao različitim čvorovima poslati različite sadržaje
}

// Hash funk, Hash predstavlja identifikator bloka (H3_V3_TX_VIEW_3 ili H3_V3_FORK)
// Ako u istom view-u napraviš dva različita bloka iste visine (što se dešava kod Byzantine lidera), dobićeš isti hash, zato dodajemo i b.Data
func (b *Block) Hash() string {
	return fmt.Sprintf("H%d_V%d_%s", b.Height, b.View, b.Data)
}

// QC (Quorum Certificate)-dokaz da je dovoljan broj replika glasao za određeni blok
/*
QC:
View = 3
BlockHash = H3_V3_TX_VIEW_3
Signers = [node-1 node-2 node-4]
*/
type QC struct {
	View      int                 //view u kome je QC formiran
	BlockHash string              // Blok koji QC potvrđuje (QC.B), za koji blok je formiran QC
	Signers   []process.ProcessID //ko je glasao od cvorova
}

// BeeGeesMessage - šta putuje mrežom
type BeeGeesMessage struct {
	Kind        string            //koja je poruka u pitanju(VoteREQ, VoteRESP, NewView)
	View        int               // u kom je view-u poslata poruka
	Block       *Block            //Koristi se kod VOTE_REQ, odnosi se na predlozeni blok
	QC          *QC               //Koristi se kada lider prosleđuje QC ili kada se šalje uz predlog
	Sender      process.ProcessID //Ko je poslao poruku
	BlockHash   string            // hash bloka za koji validator glasa
	NewViewData *NewViewMessage   // Ako je KindNewView
}

// NewViewMessage (Algorithm 3)
// Novi lider iz ovoga pokušava da rekonstruiše šta se desilo u prethodnom view-u
type NewViewMessage struct {
	ReplicaID    process.ProcessID //koja replika salje poruku
	View         int               //u kom view je poslata poruka
	LastVoteReq  *Block            //poslednji VoteREQ koji je replika prihvatila. Ovo ti je kasnije bitno za: AHL → pronalaženje dva različita bloka u istom view-u
	LastVoteResp string            // hash poslednjeg bloka za koji je replika glasala
	LastQC       *QC               //Najnoviji QC koji replika poznaje.replike tokom view change-a šalju i najnoviji QC koji znaju
}
