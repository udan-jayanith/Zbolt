package icons

import (
	"embed"
	"image"
	_ "image/png"
	"log"
	"sync"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed icons/*.png
var store embed.FS

type image_container struct {
	image *ebiten.Image
	t     time.Time
}

type cache_store struct {
	mutex  sync.Mutex
	images map[string]*image_container
}

func (cs *cache_store) Open(icon_name string) *ebiten.Image {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	icon_name = icon_name + ".png"
	image_container_ := cs.images[icon_name]
	if image_container_ == nil {
		file, err := store.Open("icons" + "/" + icon_name)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err.Error())
		}

		image_container_ = &image_container{
			image: ebiten.NewImageFromImage(img),
			t:     time.Now(),
		}
		cs.images[icon_name] = image_container_
	}

	image_container_.t = time.Now()
	return image_container_.image
}

func new() *cache_store {
	cs := cache_store{
		images: make(map[string]*image_container, 10),
	}

	/*
	go func() {
		tick := time.Tick(time.Second)
		for {
			current_time := <-tick
			cs.mutex.Lock()

			to_be_deleted := make([]string, 0, 2)
			for icon_name, ic := range cs.images {
				dur := current_time.Sub(ic.t)
				if dur.Seconds() > 1 {
					ic.image.Clear()
					to_be_deleted = append(to_be_deleted, icon_name)
				}
			}

			for _, icon_name := range to_be_deleted {
				fmt.Println("deleting", icon_name)
				delete(cs.images, icon_name)
			}
			fmt.Println(	len(cs.images))
			cs.mutex.Unlock()
		}
	}()
	 */

	return &cs
}

var (
	Store *cache_store = new()
)

type Icon struct {
	gui.DefaultWidget

	image_widget widget.Image
	IconName     string
	Point        *image.Point
}

func (icon *Icon) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	icon.image_widget.SetImage(Store.Open(icon.IconName))
	return nil
}

func (icon *Icon) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&icon.image_widget, widgetBounds.Bounds())
}

func (icon *Icon) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	if icon.Point == nil {
		u := widget.UnitSize(ctx)
		pt := image.Pt(u, u)
		icon.Point = &pt
	}
	return *icon.Point
}

func NewIcon(icon_name string, size int) *Icon {
	pt := image.Pt(size, size)
	return &Icon{
		IconName: icon_name,
		Point: &pt,
	}
}
