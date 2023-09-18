package main

import (
	"github.com/fakeyanss/usg-go/cmd"
	"github.com/fakeyanss/usg-go/pkg/routines"
	"github.com/rs/zerolog/log"
)

func main() {
	routines.Recover()
	err := cmd.Main()
	if err != nil {
		log.Error().Err(err).Msg("Fail to start usg server")
	}
}
