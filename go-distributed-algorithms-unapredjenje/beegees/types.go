package beegees //Svi fajlovi koji implementiraju algoritam biće u paketu beegees

import (
	"fmt"

	"github.com/danilokacanski/da/week03_04_parallel/process" //process.ProcessID predstavlja identitet čvora, koristimo tip koji je predvidjen za simulator
)

// Tipovi poruka (Section 4.2)
const (
	KindVoteReq  = "VOTE_REQ"  // Liderov redlog bloka(Proposal)
	KindVoteResp = "VOTE_RESP" // Validatoorov glas za blok(Vote)
	KindNewView  = "NEW_VIEW"  // Izveštaj za promenu lidera (koristi se u Slow View Change)
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
	Parent string            //Hash roditeljskog bloka
	QCanc  *QC               //blok nosi QC svog relevantnog prethodnika (omogucava proveru commit pravila,nastavak lanca posle view change-a)
	Txs    []string          //batch transakcija
	NVset  []*NewViewMessage //(možeš da prikažeš od kojih replika je lider dobio informacije i da debaguješ izbor parent bloka)
	Height int               //Visina bloka u lancu (potrebo za logove, poredjenje blokova, generisanje hash-a)
	Data   string            // DODATO: Ovo je nedostajalo za Byzantine simulaciju kako bi lider mogao različitim čvorovima poslati različite sadržaje
}

//Hash funk, Hash predstavlja identifikator bloka i koristi se kao pokazivač na prethodni blok u lancu.
/*Ako u istom view-u napraviš dva različita bloka iste visine (što se dešava kod Byzantine lidera), dobićeš isti hash, zato dodajemo i b.Data*/
func (b *Block) Hash() string {
	return fmt.Sprintf("H%d_V%d_%s", b.Height, b.View, b.Data)
}

// QC (Quorum Certificate), mozemo ispisati "QC for H2_V3 signed by [1 2 3]"
type QC struct {
	View      int                 //view u kom je predlozen blok koji QC potvrdjuje
	BlockHash string              // Blok koji QC potvrđuje (QC.B)
	Signers   []process.ProcessID //ko je glasao od cvorova
}

// BeeGeesMessage - šta putuje mrežom
type BeeGeesMessage struct {
	Kind        string            //koja je poruka u pitanju(VoteREQ, VoteRESP, NewView)
	View        int               // u kom je view-u poslata poruka
	Block       *Block            //Koristi se kod VOTE_REQ, odnosi se na predlozeni blok
	QC          *QC               //Koristi se kada lider prosleđuje QC ili kada se šalje uz predlog
	Sender      process.ProcessID //Ko je poslao poruku
	BlockHash   string            // hash bloka za koji validator glasa  (nedostajalo u handleVoteResp)
	NewViewData *NewViewMessage   // Ako je KindNewView
}

type NewViewMessage struct {
	ReplicaID    process.ProcessID //koja replika salje poruku
	View         int               //u kom view je poslata poruka
	LastVoteReq  *Block            //poslednji VoteREQ koji je replika prihvatila
	LastVoteResp string            // hash poslednjeg glasa koji je replika(validator) poslala
	LastQC       *QC               //replike tokom view change-a šalju i najnoviji QC koji znaju
}
