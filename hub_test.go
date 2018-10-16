package levelhub

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHub(t *testing.T) {
	Convey("hub", t, func() {
		key := []byte("a")
		Convey("NewLevelHub", func() {
			hub := NewLevelHub("tmp", nil)
			defer hub.Close()
			hub.Put(1, key, []byte("A"), nil)
			hub.Put(2, key, []byte("B"), nil)
			a, _ := hub.Get(1, key, nil)
			b, _ := hub.Get(2, key, nil)
			So(a[0], ShouldEqual, 'A')
			So(b[0], ShouldEqual, 'B')
		})
		Convey("Close", func() {
			hub := NewLevelHub("tmp", nil)
			defer hub.Close()
			for i := 1; i < 100; i++ {
				hub.Put(i, key, []byte("A"), nil)
			}
			for i := 1; i < 100; i++ {
				a, _ := hub.Get(1, key, nil)
				So(a[0], ShouldEqual, 'A')
			}
		})
	})
}
