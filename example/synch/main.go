package main

import (
	log "github.com/julio77it/logger"
)

func test() {
	log.Print("---------------------------------------")
	log.Trace("This a int %d", 0)
	log.Debug("This a int %d", 1)
	log.Info("This a int %d", 2)
	log.Warning("This a int %d", 3)
	log.Error("This a int %d", 4)
	if log.GetLevel() == log.FatalLvl {
		log.Fatal("This a int %d", 5)
	}
	log.Print("skip log.Fatal()")
}

func main() {
	log.SetLevel(log.TraceLvl)
	test()
	log.SetLevel(log.DebugLvl)
	test()
	log.SetLevel(log.InfoLvl)
	test()
	log.SetLevel(log.WarningLvl)
	test()
	log.SetLevel(log.ErrorLvl)
	test()
	log.SetLevel(log.FatalLvl)
	test()
}
