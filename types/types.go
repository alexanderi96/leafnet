package types

type Evento struct {
	IdEvento 			int
	Creatore			int
	Citta 				int
	DataInizio 			int64
	DataFine			int64
	Titolo 				string
	Descrizione			string
	Itinerario 			string
	MinPart				int
	MaxPart				int
	Costo				int
	Indirizzo			string
	DataScadenzaPren 	int64
	Categoria			string
	Lingue				int
	Prenotazioni		[]int
}

type MiniEvento struct {
	IdEvento 	int
	Titolo		string
	Descrizione	string
}

type Regione struct {
	IdRegione int
	NomeRegione	string
	Province	[]*Provincia //contiene gli id delle province
}

type Provincia struct {
	IdProvincia int
	NomeProvincia string
	Citta []*Citta 
}

type Citta struct {
	IdCitta int
	NomeCitta string
	Cap 	[]int
}

type Lingua struct {
	IdLingua int
	NomeLingua string
}

type Prenotazioni struct {
	IdPrenotazione	int
	Utente 			int
	DataPrenotazione int64
	Accettazione	bool
}


//Context is the struct passed to templates
type Context struct {
	Events      []MiniEvento 	//using
	Event 		Evento
	Utente 		interface{} 		//using
	CSRFToken  	string
	Referer    	string
	IsCicerone	bool
	City 		[]Citta
	Category 	[]Category
}

type Category struct {
	IdCat 	int
	Nome 	string
}
