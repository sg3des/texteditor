package main

import (
	_ "bufio"
	"fmt"
	_ "github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"io/ioutil"
	"os"
	_ "path"
	"regexp"
	_ "sync"
	"time"
)

var (
	File     = Arguments(1)
	Text     = ""
	Buffer   = ""
	Modified = false

	Opentime time.Time

	Window         *gtk.Window
	Buffertextview *gtk.TextBuffer
	GtkTextview    *gtk.TextView
)

func main() {
	gtk.Init(nil)
	GtkWindow()

	go PrepareFile()
	go FileSync()

	gtk.Main()
}

func GtkWindow() {
	Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	Window.Connect("destroy", Quit)
	Window.SetDefaultSize(600, 250)

	vbox := gtk.NewVBox(false, 1)

	menubar := Menubar(vbox)
	textview := Textview(vbox)

	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(textview)

	Window.Add(vbox)
	Window.ShowAll()
}

func Menubar(vbox *gtk.VBox) *gtk.MenuBar {
	menubar := gtk.NewMenuBar()
	file := gtk.NewMenuItemWithMnemonic("_File")
	filemenu := gtk.NewMenu()

	accel_group := gtk.NewAccelGroup()
	Window.AddAccelGroup(accel_group)

	save := gtk.NewImageMenuItemFromStock(gtk.STOCK_SAVE, accel_group)
	save_as := gtk.NewImageMenuItemFromStock(gtk.STOCK_SAVE_AS, accel_group)
	open := gtk.NewImageMenuItemFromStock(gtk.STOCK_OPEN, accel_group)
	sep := gtk.NewSeparatorMenuItem()
	quit := gtk.NewImageMenuItemFromStock(gtk.STOCK_QUIT, accel_group)

	save.Connect("activate", Save)
	save_as.Connect("activate", Save_as)
	open.Connect("activate", Open)
	quit.Connect("activate", Quit)

	filemenu.Append(save)
	filemenu.Append(save_as)
	filemenu.Append(open)
	filemenu.Append(sep)
	filemenu.Append(quit)
	file.SetSubmenu(filemenu)
	menubar.Append(file)

	return menubar
}

func Textview(vbox *gtk.VBox) *gtk.VPaned {
	vpaned := gtk.NewVPaned()

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	GtkTextview = gtk.NewTextView()
	swin.Add(GtkTextview)

	Buffertextview = GtkTextview.GetBuffer()
	Buffertextview.Connect("changed", TextChanged)

	vpaned.Pack1(swin, false, false)
	return vpaned
}

func TextChanged() {
	fmt.Println("TextChanged")
	var start, end gtk.TextIter
	Buffertextview.GetStartIter(&start)
	Buffertextview.GetEndIter(&end)
	fmt.Println("Buffer:" + Buffer)
	fmt.Println("Text  :" + Text)
	Buffer = Buffertextview.GetText(&start, &end, false)
	fmt.Println("Buffer:" + Buffer)
	fmt.Println("Text  :" + Text)
	fmt.Println("---")
	if Buffer != Text {
		Modified = true
	} else {
		Modified = false
	}
	SetTitle()
}

func Save() {
	if len(File) > 0 {
		fmt.Println("Save")
		WriteFile()
	} else {
		Save_as()
	}
}

func Save_as() {
	fmt.Println("Save_as")

	dialog := gtk.NewFileChooserDialog("Save as...", Window, gtk.FILE_CHOOSER_ACTION_SAVE, gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_SAVE, gtk.RESPONSE_ACCEPT)
	if response := dialog.Run(); response == gtk.RESPONSE_ACCEPT {
		if File = dialog.GetFilename(); len(File) > 0 {
			WriteFile()
		}
	}

	dialog.Destroy()
}

func Open() {
	fmt.Println("Open")

	dialog := gtk.NewFileChooserDialog("Open", Window, gtk.FILE_CHOOSER_ACTION_OPEN, gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_ACCEPT)
	if response := dialog.Run(); response == gtk.RESPONSE_ACCEPT {
		if File = dialog.GetFilename(); len(File) > 0 {
			PrepareFile()
		}
	}

	dialog.Destroy()
}

func Quit() {
	fmt.Println("Quit")
	gtk.MainQuit()
}

func WriteFile() {
	err := ioutil.WriteFile(File, []byte(Buffer), 0755)
	Error(err)
	if stat, err := os.Stat(File); os.IsNotExist(err) {
		fmt.Println(File + " not exists")
	} else {
		Opentime = stat.ModTime()
		Text = Buffer
		Modified = false
		SetTitle()
	}
}

func PrepareFile() {
	if len(File) > 0 {
		if stat, err := os.Stat(File); os.IsNotExist(err) {
			fmt.Println(File + " not exists")
		} else {
			Text = ReadFile(File)
			// Buffer = Text
			Opentime = stat.ModTime()
			Modified = false
			SetTitle()
			SetText()
		}
	}
}

func ReadFile(File string) string {
	if text, err := ioutil.ReadFile(File); err == nil {
		return string(text)
	} else {
		Error(err)
	}
	return ""
}

func SetTitle() {
	if len(File) > 0 {
		Window.SetTitle(File)
	}
	if match, _ := regexp.MatchString("^\\*", Window.GetTitle()); !match && Modified {
		Window.SetTitle("*" + Window.GetTitle())
	}
}

func SetText() {
	fmt.Println("SetText")
	Buffertextview.SetText(Text)
	// GtkTextview.GetBuffer().SetText(Text)
}

func Arguments(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	} else {
		return ""
	}
}

func FileSync() {
	for {
		time.Sleep(2 * time.Second)
		if len(File) > 0 && Opentime.Before(Chtimes(File)) && !Modified {
			fmt.Println("Sync")
			PrepareFile()
		}
	}
	fmt.Println("")
}

func Chtimes(File string) time.Time {
	info, err := os.Stat(File)
	Error(err)
	// if info.ModTime().Equal(info.ModTime()) {
	// 	fmt.Println("after")
	// }
	return info.ModTime()
}

func Error(e error) {
	if e != nil {
		panic(e)
	}
}
