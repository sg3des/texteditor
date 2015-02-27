package main

import (
	_ "bufio"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"fmt"
	_ "github.com/mattn/go-gtk/gdk"
	_ "github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	xcharset "golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
	_ "path"
	"regexp"
	"strconv"
	"strings"
	_ "sync"
	"texteditor/window"
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

	CharsetDir = "/srv/go/.golib/src/code.google.com/p/go-charset/data/"
)

func main() {
	gtk.Init(nil)

	window.Init()

	go PrepareFile()
	go FileSync()
	go Histories()

	gtk.Main()
}

func TextChanged() {
	var start, end gtk.TextIter
	window.Buffertextview.GetStartIter(&start)
	window.Buffertextview.GetEndIter(&end)

	buffer := window.Buffertextview.GetText(&start, &end, false)
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

func Highlight(indexes [][]int) {
	fmt.Println("Highlight")
	fmt.Println(indexes)
	for index, element := range indexes {
		var start, end gtk.TextIter
		window.Buffertextview.GetIterAtOffset(&start, element[0])
		window.Buffertextview.GetIterAtOffset(&end, element[1])
		if index == 0 {
			window.Buffertextview.ApplyTag(window.Tag1, &start, &end)
		} else {
			window.Buffertextview.ApplyTag(window.Tag2, &start, &end)
		}
	}
}

func RemoveHighlight() {
	var start, end gtk.TextIter
	window.Buffertextview.GetStartIter(&start)
	window.Buffertextview.GetEndIter(&end)
	window.Buffertextview.RemoveAllTags(&start, &end)
}

func Histories() {
	for {
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

	dialog := gtk.NewFileChooserDialog("Save as...", window.Window, gtk.FILE_CHOOSER_ACTION_SAVE, gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_SAVE, gtk.RESPONSE_ACCEPT)
	if response := dialog.Run(); response == gtk.RESPONSE_ACCEPT {
		if File = dialog.GetFilename(); len(File) > 0 {
			WriteFile()
		}
	}

	dialog.Destroy()
}

func Open() {
	fmt.Println("Open")

	dialog := gtk.NewFileChooserDialog("Open", window.Window, gtk.FILE_CHOOSER_ACTION_OPEN, gtk.STOCK_CANCEL, gtk.RESPONSE_CANCEL, gtk.STOCK_OPEN, gtk.RESPONSE_ACCEPT)
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
	substr := window.Findentry.GetText()
	str_count := ""
	RemoveHighlight()
	if len(substr) > 0 {
		re := regexp.MustCompile("(?ims)" + regexp.QuoteMeta(substr))
		indexes := re.FindAllStringSubmatchIndex(string(Buffer), -1)
		indexes2 := re.FindAllIndex([]byte(Buffer), -1)
		fmt.Println(indexes2)
		fmt.Println([]byte(Buffer))
		Highlight(indexes)
		str_count = strconv.Itoa(len(indexes))
	}
	window.Findcount.SetLabel(str_count)
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

func SetTitle() {
	if len(File) > 0 {
		window.Window.SetTitle(File)
	}
	if match, _ := regexp.MatchString("^\\*", window.Window.GetTitle()); !match && Modified {
		window.Window.SetTitle("*" + window.Window.GetTitle())
	}
}

func SetText(text string) {
	fmt.Println("SetText")
	Buffer = text
	window.Buffertextview.SetText(text)
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
