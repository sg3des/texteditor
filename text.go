package main

import (
	_ "bufio"
	"fmt"
	_ "github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"io/ioutil"
	"os"
	_ "path"
)

func main() {
	initGtkWindow()
}

func initGtkWindow() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.Connect("destroy", gtk.MainQuit)
	window.SetDefaultSize(688, 250)

	vbox := gtk.NewVBox(false, 1)

	menubar := initMenubar(vbox) // gtk.NewMenuBar()
	textview := initTextview(vbox)

	vbox.PackStart(menubar, false, false, 0)
	vbox.Add(textview)

	window.Add(vbox)
	window.ShowAll()
	gtk.Main()
}

func initMenubar(vbox *gtk.VBox) *gtk.MenuBar {
	menubar := gtk.NewMenuBar()
	menubar_file := gtk.NewMenuItemWithMnemonic("_File")
	menubar_file_submenu := gtk.NewMenu()

	var menubar_file_submenu_Save *gtk.MenuItem
	menubar_file_submenu_Save = gtk.NewMenuItemWithMnemonic("_Save")
	menubar_file_submenu_Save.Connect("activate", save)

	var menubar_file_submenu_Exit *gtk.MenuItem
	menubar_file_submenu_Exit = gtk.NewMenuItemWithMnemonic("E_xit")
	menubar_file_submenu_Exit.Connect("activate", gtk.MainQuit)

	menubar_file_submenu.Append(menubar_file_submenu_Save)
	menubar_file_submenu.Append(menubar_file_submenu_Exit)
	menubar_file.SetSubmenu(menubar_file_submenu)
	menubar.Append(menubar_file)

	return menubar
}

func initTextview(vbox *gtk.VBox) *gtk.VPaned {
	vpaned := gtk.NewVPaned()
	return vpaned
}

func save() {
	fmt.Println("save")
}

func CreateMenu(w *gtk.Window, vbox *gtk.VBox) {
	action_group := gtk.NewActionGroup("my_group")
	ui_manager := CreateUIManager()
	accel_group := ui_manager.GetAccelGroup()
	w.AddAccelGroup(accel_group)
	AddFileMenuActions(action_group)
	ui_manager.InsertActionGroup(action_group, 0)
	menubar := ui_manager.GetWidget("/MenuBar")
	vbox.PackStart(menubar, false, false, 0)
}

func CreateTextField(vbox *gtk.VBox) {
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textview := gtk.NewTextView()
	swin.Add(textview)

	buffer := textview.GetBuffer()
	buffer.SetText(Text)

	buffer.Connect("changed", func() {
		var start, end gtk.TextIter
		buffer.GetStartIter(&start)
		buffer.GetEndIter(&end)
		Text = buffer.GetText(&start, &end, false)
	})

	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	vpaned.Pack1(swin, false, false)
	// vbox.PackStart(swin, false, false, 0)
}

func CreateUIManager() *gtk.UIManager {
	UI_INFO := `
<ui>
	<menubar name='MenuBar'>
		<menu action='FileMenu'>
			<menuitem action='FileSave' />
			<menuitem action='FileQuit' />
		</menu>
	</menubar>
</ui>
`
	ui_manager := gtk.NewUIManager()
	ui_manager.AddUIFromString(UI_INFO)
	return ui_manager
}

func AddFileMenuActions(action_group *gtk.ActionGroup) {
	action_group.AddAction(gtk.NewAction("FileMenu", "File", "", ""))

	action_filenewmenu := gtk.NewAction("FileNew", "", "", gtk.STOCK_NEW)
	action_group.AddAction(action_filenewmenu)

	action_filequit := gtk.NewAction("FileQuit", "", "", gtk.STOCK_QUIT)
	action_filequit.Connect("activate", OnMenuFileQuit)
	action_group.AddActionWithAccel(action_filequit, "")

	action_filesave := gtk.NewAction("FileSave", "", "", gtk.STOCK_SAVE)
	action_filesave.Connect("activate", OnMenuFileSave)
	action_group.AddActionWithAccel(action_filesave, "<ctrl>s")
}

func OnMenuFileNewGeneric() {
	fmt.Println("A File|New menu item was selected.")
}

func OnMenuFileQuit() {
	fmt.Println("quit app...")
	gtk.MainQuit()
}

func OnMenuFileSave() {
	if len(File) == 0 {
		OnMenuFileSaveAs()
	} else {
		fmt.Println("Save: " + File)
		err := ioutil.WriteFile(File, []byte(Text), 0777)
		Check(err)
	}
}

func OnMenuFileSaveAs() {
	fmt.Println("SaveAs")
	// fmt.Println(Window)
	// saveas := gtk.NewFileChooserDialog("Save as...", Window, gtk.FILE_CHOOSER_ACTION_SAVE, gtk.STOCK_SAVE_AS, gtk.RESPONSE_ACCEPT)
	// saveas.Response(func() {
	// 	fmt.Println(saveas.GetFilename())
	// 	saveas.Destroy()
	// })
	// saveas.Run()
}

//--------------------------------------------------------
// Main
//--------------------------------------------------------

type Gtk_struct struct {
	Window *gtk.Window
}

var (
	File = Arguments(1)
	Text = CheckFile(File)
)

func CheckFile(file string) string {
	if err := FileExists(file); err == nil {
		// fmt.Println("CheckFile: " + file)
		return OpenFile(file)
	} else {
		// fmt.Println("CheckFile else: " + file)
		return ""
	}
	// return ""
}

func Arguments(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	} else {
		return ""
	}
}

func FileExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
}

func OpenFile(file string) string {
	fmt.Println("Open: " + file)
	text, err := ioutil.ReadFile(file)
	Check(err)
	return string(text)
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// func Save(file File) {
// 	fmt.Println(file.Text)
// 	fmt.Println("save")
// 	err := ioutil.WriteFile(file.Filepath, file.Text, 0777)
// 	Check(err)
// }

// func Frame(file *File) {
// 	gtk.Init(nil)
// 	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
// 	vbox := gtk.NewVBox(false, 1)

// 	menubar := gtk.NewMenuBar()
// 	vbox.PackStart(menubar, false, false, 0)

// 	vpaned := gtk.NewVPaned()
// 	vbox.Add(vpaned)
// 	//--------------------------------------------------------
// 	// Text
// 	//--------------------------------------------------------
// 	swin := gtk.NewScrolledWindow(nil, nil)
// 	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
// 	textview := gtk.NewTextView()
// 	swin.Add(textview)

// 	buffer := textview.GetBuffer()
// 	buffer.SetText(string(file.Text))

// 	vpaned.Pack1(swin, false, false)

// 	//--------------------------------------------------------
// 	// MenuBar
// 	//--------------------------------------------------------
// 	// var menuitem *gtk.MenuItem

// 	// menu_file := gtk.NewMenuItemWithMnemonic("_File")
// 	// menubar.Append(menu_file)
// 	// submenu_file := gtk.NewMenu()
// 	// menu_file.SetSubmenu(submenu_file)

// 	// menuitem_save := gtk.NewMenuItemWithMnemonic("_Save")
// 	// menuitem_save.Connect("activate", func() {
// 	// 	var start, end gtk.TextIter
// 	// 	buffer.GetStartIter(&start)
// 	// 	buffer.GetEndIter(&end)
// 	// 	text := buffer.GetText(&start, &end, false)
// 	// 	file := File{
// 	// 		file.Filename,
// 	// 		file.Filepath,
// 	// 		[]byte(text),
// 	// 	}
// 	// 	save(file)
// 	// })
// 	// menu_file.AddActionWithAccel(menuitem_save, "<control>S")
// 	// submenu_file.Append(menuitem_save)

// 	// menuitem_exit := gtk.NewMenuItemWithMnemonic("E_xit")
// 	// menuitem_exit.Connect("activate", func() {
// 	// 	gtk.MainQuit()
// 	// })
// 	// submenu_file.Append(menuitem_exit)

// 	//--------------------------------------------------------
// 	// EndFrame
// 	//--------------------------------------------------------

// 	window.SetPosition(gtk.WIN_POS_CENTER)
// 	window.SetTitle(file.Filepath)
// 	window.SetIconName("text-plain")
// 	window.Connect("destroy", gtk.MainQuit)

// 	window.Add(vbox)

// 	window.SetSizeRequest(250, 100)
// 	window.ShowAll()
// 	gtk.Main()
// }
