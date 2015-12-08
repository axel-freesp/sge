package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

func aboutAbout(menu *GoAppMenu) {
	menu.aboutdialog.Run()
	menu.aboutdialog.Hide()
}

func aboutHelp(menu *GoAppMenu) {
	log.Println("aboutHelp")
}

func MenuAboutInit(menu *GoAppMenu) {
	menu.aboutAbout.Connect("activate", func() { aboutAbout(menu) })
	menu.aboutHelp.Connect("activate", func() { aboutHelp(menu) })

	menu.aboutdialog, _ = gtk.AboutDialogNew()
	menu.aboutdialog.SetCopyright("Copyright (C) 2015 by Axel von Blomberg. All rights reserved.")
	menu.aboutdialog.SetLicenseType(gtk.LICENSE_BSD)
	menu.aboutdialog.SetAuthors([]string{"Axel von Blomberg"})
	menu.aboutdialog.SetProgramName("sge (Signal Graph Editor)")
	iconPath := os.Getenv("SGE_ICON_PATH")
	if len(iconPath) > 0 {
		imageLogo, err := gdk.PixbufNewFromFile(fmt.Sprintf("%s/sge-logo-small.png", iconPath))
		if err == nil {
			menu.aboutdialog.SetLogo(imageLogo)
		}
	}
}
