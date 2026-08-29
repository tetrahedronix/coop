## I sistemi

Ogni sistema nell'architettura ECS dovrebbe operare solo quelle entità che possiedono i componenti richiesti, non su tutte indistintamente. Un'implementazione in cui si processano tutte le entità per ogni sistema è inefficiente e concettualmente sbagliato.

L'approccio corretto prevede che i sistemi definiscano una "firma" (o *signature*) dei componenti di cui hanno bisogno, e il mondo fornisca un modo per iterare solo le entità che possiedono quella firma. Anche nella variante Tharsis del paradigma, il sistema legge dal past e scrive nel future, ma deve operare solo sulle entità che hanno componenti pertinenti (per esempio, il `RenderinSystem` processa le entità che possiedono almeno il componente `Renderer`).

Quindi, il codice deve avere un meccanismo di query. Le opzioni possibili sono:

* Ogni sistema, nel suo metodo `Process`, controlla se l'entità ha i componenti necessari e, se non li ha, ritorna subito.
* È il loop che, prima di chiamare `Process`, filtra le entità in base a una funzione di match associata al sistema.
* Si utilizza un sistema di archetype/query.

Approccio scelto in Googol:

* Si definiscono i sistemi come interfacce che dichiarano una signature (insieme di tipi di componenti richiesti) e un metodo `Process(e *Entity)` che verrà chiamato solo per le entità che soddisfano la signature. Poi nel *loop*, prima di chiamare `Process`, si filtrano le entità in base alla signature di ogni sistema, in modo da iterare solo quelle pertinenti.
* Ogni sistema legge dal past (`GetPastComponent`) e scrive nel future (`GetWritableComponent`) solo per i componenti che sa gestire.


## La struttura World

In `game.go` si definisce la struttura `World` con un logger, un flag per indicare se è attivo, una lista di entità, un mutex per la concorrenza, opzionalmente un registro dei sistemi (ma questi si preferisce gestirli nel `Loop`):

```go
type World struct {
    // Proteggere a slice di entità con un mutex
	mu       sync.RWMutex
	entities []*ecs.Entity
	Logger   *log.Logger
	// Whether the world tick should execute
	enabled bool
}
```

## Game Loop

Per realizzare la struttura dati (con i suoi metodi) per la gestione del loop principale nel motore **Googol**, occorre considerare i seguenti aspetti:

* **Thread-safety**. Consentire l'aggiunta e rimozione dinamica dei sistemi durante l'esecuzione (usando `AddSystem`, `RemoveSystem`). Viene realizzato usando `sync.RwMutex` interno.
* **Gestione dei sistemi**. I sistemi devono essere iterati ogni tick e ricevere tutte le entità. Viene realzizato usando slice e implementando metodi o funzioni helper.
* **Aggiornamento del World**. Per rispettare il paradigma *double buffer* di Tarsis, si implementa una funzione `SwapBuffers()` in `googol/ecs/entity.go` e un metodo `SwapBuffers()` per il tipo `World` chiamato alla fine di ogni tick.
* **Frame rate**. Il loop principale *gira* a 60 FPS. Per garantire una frequenza di aggiornamento stabile si implementa (idealmente) con un ticker: `time.NewTicker(time.Second/60)`.
* **Separazione tra World e Loop**. Il loop conosce `World` (per ottenere le entità e fare lo swap), ma non i dettagli dei componenti; delega invece ai sistemi la logica di lettura e scrittura.

```go
type Loop struct {
	mu      sync.RWMutex
	systems []System
    world   *World    
	ticker  *time.Ticker
	stopCh  chan struct{}
	running bool
}
```

In questa struttura:

* `systems`: slice di sistemi che implementano `ecs.System`.
* `world`: puntatore al mondo di gioco, necessario per ottenere le entità e per scambiare i buffer.
* `ticker`: genera il segnale a 60 Hz.
* `stopCh`: canale per fermare il loop pulitamente.
* `running`: flag per sapere se il loop è in esecuzione.

### Costruttore

Il costruttore del loop riceve per argomento un mondo e restituisce una struttura pronta per attaccare entità ed eseguire sistemi.

```go
func NewLoop(w *World) *Loop {
	return &Loop{
		world:   w,
		systems: make([]ecs.System, 0),
		stopCh:  make(chan struct{}),
	}
}
```



