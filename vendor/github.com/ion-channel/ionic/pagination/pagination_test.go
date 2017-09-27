package pagination

import (
	"net/url"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestPagination(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Pagination", func() {
		g.Describe("Instantiation", func() {
			g.It("should return a new pagination object", func() {
				o := 20
				l := 10
				p := New(o, l)

				Expect(p.Limit).To(Equal(l))
				Expect(p.Offset).To(Equal(o))
			})

			g.It("should set limit to -1 if not positive", func() {
				o := 20
				l := -100
				p := New(o, l)

				Expect(p.Limit).NotTo(Equal(l))
				Expect(p.Limit).To(Equal(-1))
				Expect(p.Offset).To(Equal(o))
			})

			g.It("should prevent offset from being negative", func() {
				o := -1
				l := 10
				p := New(o, l)

				Expect(p.Limit).To(Equal(l))
				Expect(p.Offset).NotTo(Equal(o))
				Expect(p.Offset).To(Equal(0))
			})
		})

		g.Describe("URL Params", func() {
			g.It("should append the url pagination params", func() {
				params := &url.Values{}
				params.Set("foo", "bar")
				p := Pagination{Limit: 10, Offset: 20}

				p.AddParams(params)
				Expect(params.Get("foo")).To(Equal("bar"))
				Expect(params.Get("limit")).To(Equal("10"))
				Expect(params.Get("offset")).To(Equal("20"))
			})
		})

		g.Describe("Incrementing", func() {
			g.It("should increment the offset up", func() {
				l := 10
				o := 15
				p := Pagination{Limit: l, Offset: o}

				p.Up()
				Expect(p.Offset).To(Equal(o + l))
				Expect(p.Limit).To(Equal(l))
			})

			g.It("should ignore incrementing up if all items", func() {
				l := AllItems.Limit
				o := AllItems.Offset

				AllItems.Up()
				Expect(AllItems.Limit).To(Equal(l))
				Expect(AllItems.Offset).To(Equal(o))
			})

			g.It("should increment the offset down", func() {
				l := 10
				o := 15
				p := Pagination{Limit: l, Offset: o}

				p.Down()
				Expect(p.Offset).To(Equal(o - l))
				Expect(p.Limit).To(Equal(l))
			})

			g.It("should not increment past 0", func() {
				l := 10
				o := 15
				p := Pagination{Limit: l, Offset: o}

				p.Down()
				Expect(p.Offset).To(Equal(o - l))
				Expect(p.Limit).To(Equal(l))
				p.Down()
				Expect(p.Offset).To(Equal(0))
				Expect(p.Limit).To(Equal(l))
			})

			g.It("should ignore incrementing up if all items", func() {
				l := AllItems.Limit
				o := AllItems.Offset

				AllItems.Down()
				Expect(AllItems.Limit).To(Equal(l))
				Expect(AllItems.Offset).To(Equal(o))
			})
		})

		g.Describe("SQLification", func() {
			g.It("should convert to a sql string", func() {
				p := Pagination{Limit: 10, Offset: 15}

				sql := p.SQL()
				Expect(sql).To(ContainSubstring("OFFSET 15"))
				Expect(sql).To(ContainSubstring("LIMIT 10"))
			})

			g.It("should leave limit off the sql string if 0", func() {
				p := Pagination{Limit: 0, Offset: 15}

				sql := p.SQL()
				Expect(sql).To(ContainSubstring("OFFSET 15"))
				Expect(sql).NotTo(ContainSubstring("LIMIT"))
			})

			g.It("should be blank when all items", func() {
				sql := AllItems.SQL()
				Expect(sql).To(Equal(""))
			})
		})
	})
}
