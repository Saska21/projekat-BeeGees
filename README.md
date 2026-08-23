# projekat-BeeGees

Implementacija i simulacija **BeeGees** konsenzus algoritma. 

Simulator modeluje mrežu od 4 nezavisna čvora (procesa), gde je svaki čvor zasebna Go goroutine koja komunicira isključivo razmenom poruka preko simulirane, nepouzdane mreže (FairLoss → Stubborn → Perfect Link). 
Podržana su tri scenarija izvršavanja: normalan rad, čvor koji ćuti (crash/silent) i vizantijski (zlonameran) lider koji šalje konfliktne predloge.

## Preduslovi

- **Go 1.21+** — testirano sa `go1.26.1`
- Nisu potrebne eksterne biblioteke — projekat koristi isključivo standardnu Go biblioteku i sopstvene interne pakete

## Pokretanje

### 1. Preuzmi projekat

Direktno sa GitHub-a: na stranici repozitorijuma klikni zeleno dugme **"Code" → "Download ZIP"**, pa raspakuj preuzeti fajl.
Na kraju na svom računaru treba da dobiješ folder **projekat-BeeGees-main** koji sadrži foldere **`go-distributed-algorithms-main` i go-distributed-algorithms-unapredjenje** — otvori jedan od njih u editoru (npr. VS Code).

### 2. Podesi scenario u `main.go`

U korenu projekta otvori `main.go` i izmeni promenljivu `selectedScenario`:

```go
selectedScenario := "byzantine"   // "normal" | "silent" | "byzantine"
```

Dostupne vrednosti:

| Vrednost | Opis |
|---|---|
| `"normal"` | Svi čvorovi su pošteni — demonstrira Fast View Change kroz normalan rad protokola |
| `"silent"` | Jedan čvor ne šalje predlog niti glas kada je lider — demonstrira Liveness Timeout i Slow View Change |
| `"byzantine"` | Jedan čvor kao lider šalje konfliktne predloge različitim čvorovima — demonstrira equivocation detekciju, AHL granu Commit Rule-a, blacklistovanje lidera i ravnomernu rotaciju preostalih čvorova |

Ako je izabran scenario `"silent"` ili `"byzantine"`, potrebno je odabrati i **koji konkretan čvor** dobija to ponašanje — u istom fajlu, u petlji koja registruje čvorove:

```go
if selectedScenario == "silent" && pid == "node-2" {
    behavior = beegees.BehaviorSilent
} else if selectedScenario == "byzantine" && pid == "node-3" {
    behavior = beegees.BehaviorByzantine
}
```

Izmeni `"node-2"` odnosno `"node-3"` na željeni čvor (`"node-1"`–`"node-4"`) ako želiš da testiraš drugačiju kombinaciju.

### 3. Pokreni simulaciju

Iz terminala, u korenu foldera `go-distributed-algorithms-main`(go-distributed-algorithms-unapredjenje):

```bash
go run main.go
```

Za novi scenario, ponovi korak 2 (sačuvaj izmene u `main.go`) i ponovo pokreni komandu iz koraka 3.

## Šta se ispisuje u konzoli

Simulacija se pokreće direktno u terminalu i ispisuje hronološki tok događaja u realnom vremenu (sa veštačkim pauzama radi lakšeg praćenja). Ključni simboli u ispisu:

| Oznaka | Značenje |
|---|---|
| `--- [VIEW N] LIDER ... predlaže blok ... ---` | Novi blok je predložen |
| `Glasam za blok ...` | Validator je glasao za predloženi blok |
| `FORMIRAN QC za blok ...` | Sakupljen je kvorum glasova (2f+1) i formiran je Quorum Certificate |
| `COMMIT (CHL) / (AHL)` | Blok je bezbedno komitovan (uzastopni QC-ovi, ili neuzastopni ali bez dokaza o prevari) |
| `Liveness Timeout!` | Istekao je view timer — pokreće se Slow View Change |
| `Ja sam lider View-a N. Pokrećem MAT tajmer...` | Novi lider je pokrenuo Materialization tajmer |
| `MAT istekao...` | MAT tajmer je istekao pre pune materijalizacije — lider šalje blok sa onim što ima |
| `Lider šalje KONFLIKTNE blokove!` | (samo u `byzantine` scenariju) Zlonameran lider šalje različite predloge |
| `AHL Provera (Rupa između View X i Y)...` | Otkrivena je "rupa" (neuzastopni QC-ovi) — pokreće se provera equivokacije |
| `ABORT: Otkriven dokaz o prevari!` | Equivokacija je pronađena — blok se NE komituje |
| `Lider ... je BLACKLISTOVAN` | Otkriveni vizantijski lider je trajno isključen iz rotacije budućih lidera |


## (Opciono) Podešavanje parametara simulacije

Nekoliko dodatnih parametara se, ako je potrebno, može podesiti u `main.go`, unutar poziva `simrt.NewRuntime(...)` — podrazumevane vrednosti su dovoljne za standardno pokretanje i nije ih neophodno menjati:

```go
rt := simrt.NewRuntime(pl, fm,
    simrt.WithMaxDuration(15*time.Second),      // maksimalno trajanje simulacije
    simrt.WithIdleTimeout(5*time.Second),       // gasi simulaciju ako niko ne komunicira ovoliko dugo
    simrt.WithRetransmitInterval(1*time.Second),// koliko često se ponavlja slanje neisporučenih poruka
    simrt.WithVerbose(false),                   // true = ispisuju se i mrežne poruke niskog nivoa
)
```

Procenat gubitka poruka na mreži podešava se u konstrukciji linka (prvi argument, `0.0` = bez gubitaka):

```go
fl := link.NewFairLossLink(0.0, 42)
```

## Autor
Aleksandra Radojičić, E2 103/2025


