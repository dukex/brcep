package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leogregianin/brcep/api"
	"github.com/leogregianin/brcep/api/cepaberto"
	"github.com/leogregianin/brcep/api/viacep"
	"github.com/leogregianin/brcep/config"
	"github.com/leogregianin/brcep/config/env"
	"github.com/leogregianin/brcep/config/flag"
	"github.com/leogregianin/brcep/handler"
)

func main() {
	fmt.Println(`   ___.                                  `)
	fmt.Println(`   \_ |_________   ____  ____ ______     `)
	fmt.Println(`    | __ \_  __ \_/ ___\/ __ \\____ \    `)
	fmt.Println(`    | \_\ \  | \/\  \__\  ___/|  |_> >   `)
	fmt.Println(`    |___  /__|    \___  >___  >   __/    `)
	fmt.Println(`        \/            \/    \/|__|       `)

	cfg, err := config.NewConfig([]config.Loader{
		flag.NewFlagLoader(),
		env.NewEnvLoader(),
	})
	if err != nil {
		panic(err)
	}

	gin.SetMode(cfg.GetGinOperationMode())

	var (
		cepApis = map[string]api.Api{
			viacep.ID: viacep.NewViaCepApi(cfg.ViaCepUrl, http.DefaultClient),
		}
	)

	if len(cfg.CepAbertoToken) > 0 {
		cepApis[cepaberto.ID] = cepaberto.NewCepAbertoApi(
			cfg.CepAbertoUrl,
			cfg.CepAbertoToken,
			http.DefaultClient)
	}

	var cepHandler = &handler.CepHandler{
		PreferredApi: cfg.PreferredAPI,
		CepApis:      cepApis,
	}

	router := gin.Default()
	router.Use(gin.ErrorLogger())
	router.GET("/:cep/json", cepHandler.Handle)

	server := &http.Server{
		Addr:           cfg.Address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
