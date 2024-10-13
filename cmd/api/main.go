package main

func main() {
	cfg := config{
		addr: ":4000",
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	app.run(mux)

}
