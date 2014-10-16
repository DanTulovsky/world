package world

import "net/http"

func (w *World) HomeHandler(writer http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if err != nil {
			http.Error(writer, err.Error(), status)
		}
	}()

	w.Show(writer)
	w.ShowSettings(writer)
	w.ShowGrid(writer)

}
