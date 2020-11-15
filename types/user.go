package types

/*
	Ogni utente sarà di tipo User, di default il tipo user fa riferimento 
	alla struct Globetrotter, contenente ogni campo fondamentale di ogni 
	utente.
	Nel momento in cui un utente si registra come Cicerone, l'Utente verrà 
	automaticamente "specializzato" in Cicerone e guadagnerà automaticamente
	gli attributi del Globetrotter + quelli specifici del Cicerone.

	Ogni Utente ha un metodo IsCicerone() di tipo booleano (in questo caso stiamo
	utilizzando il golang come se fosse un linguaggio ad oggetti, e le interfacce
	ci permettono di applicare il concetto di "polimorfismo").
	Nel caso un utente sia specializzato come Globetrotter essa restitirà false,
	viceversa true nel caso di Cicerone.
*/

type User interface{
	idUser			int
	username		string
	email 			string
	password 		string
}

type Person struct {
	User
    Nome 			string
    Cognome 		string
    Sesso 			string
    DataNascita		int64
    Email 			string
    Password 		string
}

type Couple struct {
	idCouple 		int
	father 			Person
	mother 			Person
}
