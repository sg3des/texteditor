package main

import (
	_ "bufio"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	_ "github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	xcharset "golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
	_ "path"
	"regexp"
	"strings"
	_ "sync"
	"time"
)

var (
	File        = Arguments(1)
	Text        = ""
	Buffer      = ""
	History     = make(map[int]string)
	History_key = 0
	Modified    = false

	Opentime time.Time

	Window         *gtk.Window
	Buffertextview *gtk.TextBuffer
	GtkTextview    *gtk.TextView

	CharsetDir = "/srv/go/.golib/src/code.google.com/p/go-charset/data/"
)

func main() {
	gtk.Init(nil)
	GtkWindow()

	go PrepareFile()
	go FileSync()
	go Histories()

	gtk.Main()
}

func GtkWindow() {
	Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	Window.Connect("destroy", Quit)
	Window.SetDefaultSize(600, 250)

	vbox := gtk.NewVBox(false, 1)

	menubar := Menubar()
	notice := Notice()
	textview := Textview()

	vbox.PackStart(menubar, false, false, 0)
	vbox.PackStart(notice, false, false, 0)
	vbox.Add(textview)

	Window.Add(vbox)
	Window.ShowAll()
}

func Menubar() *gtk.Widget {
	ui_xml := `
<ui>
	<menubar name='MenuBar'>
		<menu action='File'>
			<menuitem action='Save' />
			<menuitem action='Save_as' />
			<menuitem action='Open' />
			<separator />
			<menuitem action='Quit' />
		</menu>
		<menu action='Edit'>
			<menuitem action='Find' />
			<menuitem action='Find_next' />
			<menuitem action='Replace' />
			<separator />
			<menuitem action='Redo' />
			<menuitem action='Undo' />
		</menu>
	</menubar>
</ui>
	`
	ui := gtk.NewUIManager()
	ui.AddUIFromString(ui_xml)
	action_group := gtk.NewActionGroup("MenuBar")
	accel_group := ui.GetAccelGroup()
	Window.AddAccelGroup(accel_group)

	action_group.AddAction(gtk.NewAction("File", "File", "", ""))
	action_group.AddAction(gtk.NewAction("Edit", "Edit", "", ""))

	save := gtk.NewAction("Save", "Save", "", gtk.STOCK_SAVE)
	save_as := gtk.NewAction("Save_as", "Save as", "", gtk.STOCK_SAVE_AS)
	open := gtk.NewAction("Open", "Open", "", gtk.STOCK_OPEN)
	quit := gtk.NewAction("Quit", "Quit", "", gtk.STOCK_QUIT)

	find := gtk.NewAction("Find", "Find", "", gtk.STOCK_FIND)
	findnext := gtk.NewAction("Find_next", "Find next", "", "")
	replace := gtk.NewAction("Replace", "Replace", "", gtk.STOCK_FIND_AND_REPLACE)
	redo := gtk.NewAction("Redo", "Redo", "", gtk.STOCK_REDO)
	undo := gtk.NewAction("Undo", "Undo", "", gtk.STOCK_UNDO)

	save.Connect("activate", Save)
	save_as.Connect("activate", Save_as)
	open.Connect("activate", Open)
	quit.Connect("activate", Quit)

	find.Connect("activate", Find)
	findnext.Connect("activate", Findnext)
	replace.Connect("activate", Replace)
	redo.Connect("activate", Redo)
	undo.Connect("activate", Undo)

	action_group.AddActionWithAccel(save, "<control>S")
	action_group.AddAction(save_as)
	action_group.AddActionWithAccel(open, "<control>O")
	action_group.AddActionWithAccel(quit, "<control>Q")

	action_group.AddActionWithAccel(find, "<control>F")
	action_group.AddActionWithAccel(findnext, "<control>G")
	action_group.AddActionWithAccel(replace, "<control>R")
	action_group.AddActionWithAccel(redo, "<control>Y")
	action_group.AddActionWithAccel(undo, "<control>Z")

	ui.InsertActionGroup(action_group, 0)
	menubar := ui.GetWidget("/MenuBar")

	return menubar
}

func Notice() *gtk.HBox {
	toolbar := gtk.NewHBox(false, 5)
	// button1 := gtk.NewButtonWithLabel("asd1")
	// button2 := gtk.NewButtonWithLabel("asd2")
	// button3 := gtk.NewButtonWithLabel("asd3")
	// button4 := gtk.NewButtonWithLabel("asd4")
	// toolbar.PackStart(button1, false, false, 0)
	// toolbar.PackStart(button2, false, false, 0)
	// toolbar.PackStart(button3, false, false, 0)
	// toolbar.PackStart(button4, false, false, 0)

	return toolbar
}

func Textview() *gtk.VPaned {
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
	// fmt.Println("TextChanged")
	var start, end gtk.TextIter
	Buffertextview.GetStartIter(&start)
	Buffertextview.GetEndIter(&end)

	buffer := Buffertextview.GetText(&start, &end, false)
	fmt.Println("Buffer: " + Buffer)
	fmt.Println("buffer: " + buffer)
	if buffer != Buffer && len(buffer) > 0 {
		fmt.Println("Update Buffer")
		Buffer = buffer
		History_key = len(History)
	}

	if Buffer != Text {
		Modified = true
	} else {
		Modified = false
	}
	SetTitle()
}

func Histories() {
	for {
		fmt.Println(History_key)
		if len(Buffer) > 0 && Buffer != History[len(History)-1] && History_key == len(History) {
			fmt.Println("add to History")
			History[len(History)] = Buffer
		}
		time.Sleep(2 * time.Second)
	}
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

func Find() {
	fmt.Println("Find")
}

func Findnext() {
	fmt.Println("Findnext")
}

func Replace() {
	fmt.Println("Replace")
}

func Redo() {
	fmt.Println("Redo")
	if History_key == -1 {
		History_key = len(History) - 1
	}
	if len(History)-1 > History_key {
		History_key++
		SetText(History[History_key])
	}
	fmt.Println(History_key)
}

func Undo() {
	fmt.Println("Undo")
	if History_key == -1 {
		History_key = len(History) - 1
	}
	if History_key > 0 && len(History) > 0 {
		History_key--
		SetText(History[History_key])
	}
	fmt.Println(History_key)
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
			SetText(Text)
		}
	}
}

func ReadFile(File string) string {
	if text, err := ioutil.ReadFile(File); err == nil {
		text := CharsetConverter(text)
		return text
	} else {
		Error(err)
	}
	return ""
}

func CharsetConverter(text []byte) string {
	_, charsetName, _ := xcharset.DetermineEncoding(text, "[]byte")
	if charsetName != "utf-8" {
		r, err := charset.NewReader("windows-1251", strings.NewReader(string(text)))
		Error(err)
		text, err = ioutil.ReadAll(r)
		Error(err)
	}
	return string(text)
}

// func Notice() {

// }

func SetTitle() {
	if len(File) > 0 {
		Window.SetTitle(File)
	}
	if match, _ := regexp.MatchString("^\\*", Window.GetTitle()); !match && Modified {
		Window.SetTitle("*" + Window.GetTitle())
	}
}

func SetText(text string) {
	fmt.Println("SetText")
	Buffer = text
	Buffertextview.SetText(text)
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
