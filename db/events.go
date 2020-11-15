package db

import (
	"log"
	"database/sql" 
	"github.com/alexanderi96/leafnet/types"
)


func CreateEvent(e types.Evento) error {
	return gQuery("insert into Eventi(FkCiceroneEvento, DataInizioEvento, DataFineEvento, TitoloEvento, DescrizioneEvento, ItinerarioEvento, NumeroMinPart, NumeroMaxPart, CostoEvento, LuogoRitrovoEvento, PrenotabileFinoAl) values(?,?,?,?,?,?,?,?,?,?,?)",
					e.Creatore, e.DataInizio, e.DataFine, e.Titolo,
					e.Descrizione, e.Itinerario, e.MinPart,
					e.MaxPart, e.Costo, e.Indirizzo,
					e.DataScadenzaPren)
}

func GetEvents() (E []types.MiniEvento, e error) {
	var Evento types.MiniEvento
	var rows *sql.Rows

	basicSQL := "select IdEvento, TitoloEvento, DescrizioneEvento from Eventi"
	rows = database.query(basicSQL)

	defer rows.Close()
	for rows.Next() {
		Evento = types.MiniEvento {}
		err = rows.Scan(&Evento.IdEvento, &Evento.Titolo, &Evento.Descrizione)
		if err != nil {
			return E, e
		}
		E = append(E, Evento)
	}
	return E, nil

}

func GetEventById(id int) (Evento types.Evento, e error) {
	log.Println("Getting Event")
	var rows *sql.Rows

	basicSQL := "select IdEvento, FkCiceroneEvento, DataInizioEvento, DataFineEvento, TitoloEvento, DescrizioneEvento, ItinerarioEvento, NumeroMinPart, NumeroMaxPart, CostoEvento, LuogoRitrovoEvento, PrenotabileFinoAl from Eventi where IdEvento = ?"
	rows = database.query(basicSQL, id)

	defer rows.Close()
	if rows.Next() {

		e = rows.Scan(&Evento.IdEvento, &Evento.Creatore, &Evento.DataInizio, &Evento.DataFine, &Evento.Titolo, &Evento.Descrizione, &Evento.Itinerario, &Evento.MinPart, &Evento.MaxPart, &Evento.Costo, &Evento.Indirizzo, &Evento.DataScadenzaPren)
		if e != nil {
			Evento = types.Evento {}
			return Evento, e
		}
	}

	return Evento, nil
}

func DeleteEveryEvent() error {
	basicSQL := "delete from Eventi"
	return gQuery(basicSQL)
}

func DeleteEventById(id int) (e error) {
	basicSQL := "delete from Eventi where IdEvento = ?"
	e = gQuery(basicSQL, id)
	//TODO: notify every interested user that this event was deleted 
	return
}


//TODO: improve search query in order to filter the results
func SearchEvent(query string) (E []types.MiniEvento, e error) {
	var Evento types.MiniEvento
	var rows *sql.Rows

	basicSQL := "select IdEvento, TitoloEvento, DescrizioneEvento from Eventi where (TitoloEvento like '%" + query + "%' or DescrizioneEvento like '%" + query + "%' or ItinerarioEvento like '%" + query +"%' or LuogoRitrovoEvento like '%" + query + "%')"
	rows = database.query(basicSQL)

	defer rows.Close()
	for rows.Next() {
		Evento = types.MiniEvento {}
		err = rows.Scan(&Evento.IdEvento, &Evento.Titolo, &Evento.Descrizione)
		if err != nil {
			return E, e
		}
		E = append(E, Evento)
	}
	return E, nil

}
