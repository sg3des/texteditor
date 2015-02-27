package window

import (
	"fmt"
	"github.com/mattn/go-gtk/gtk"
)

var (
	Window         *gtk.Window
	Buffertextview *gtk.TextBuffer
	Textview       *gtk.TextView
	Findentry      *gtk.Entry
	Findcount      *gtk.Label

	TagProps1 = map[string]string{"background": "#FF0000"}
	TagProps2 = map[string]string{"background": "#FFAA00"}
	Tag1      *gtk.TextTag
	Tag2      *gtk.TextTag
)

func Init() {
	fmt.Println("window.Init")
	Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	texteditor.Save()
	// Connect[0] = Window
	// this.Save()
	// Window.Connect("destroy", Quit)
	Window.SetDefaultSize(600, 250)

	vbox := gtk.NewVBox(false, 1)

	menubar := gtkMenubar()
	textview := gtkTextview()
	find := gtkFind()

	vbox.PackStart(menubar, false, false, 0)
	vbox.PackEnd(find, false, false, 0)
	vbox.Add(textview)

	Window.Add(vbox)
	Window.ShowAll()

}

func gtkMenubar() *gtk.Widget {
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

	// save.Connect("activate", Save)
	// save_as.Connect("activate", Save_as)
	// open.Connect("activate", Open)
	// quit.Connect("activate", Quit)

	// find.Connect("activate", Find)
	// findnext.Connect("activate", Findnext)
	// replace.Connect("activate", Replace)
	// redo.Connect("activate", Redo)
	// undo.Connect("activate", Undo)

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

func gtkTextview() *gtk.VPaned {
	vpaned := gtk.NewVPaned()

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	Textview = gtk.NewTextView()
	Textview.SetWrapMode(gtk.WRAP_WORD)
	swin.Add(Textview)

	Buffertextview = Textview.GetBuffer()
	// Buffertextview.Connect("changed", TextChanged)

	Tag1 = Buffertextview.CreateTag("Tag1", TagProps1)
	Tag2 = Buffertextview.CreateTag("Tag2", TagProps2)

	vpaned.Pack1(swin, false, false)
	return vpaned
}

func gtkFind() *gtk.Frame {
	frame := gtk.NewFrame("Find")
	hbox := gtk.NewHBox(false, 5)

	Findentry = gtk.NewEntry()
	// Findentry.Connect("changed", Find)

	Findcount = gtk.NewLabel("")

	button := gtk.NewButtonWithLabel("Find")
	// button.Connect("activate", Find)

	hbox.Add(Findentry)
	hbox.PackEnd(button, false, false, 0)
	hbox.PackEnd(Findcount, false, false, 0)
	frame.Add(hbox)
	return frame
}
