package apidsl_test

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
	"github.com/goadesign/goa/dslengine"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MediaType", func() {
	var name string
	var dslFunc func()

	var mt *MediaTypeDefinition

	BeforeEach(func() {
		InitDesign()
		dslengine.Errors = nil
		name = ""
		dslFunc = nil
	})

	JustBeforeEach(func() {
		mt = MediaType(name, dslFunc)
		dslengine.Run()
		Ω(dslengine.Errors).ShouldNot(HaveOccurred())
	})

	Context("with no dsl and no identifier", func() {
		It("produces an error", func() {
			Ω(mt).ShouldNot(BeNil())
			Ω(mt.Validate()).Should(HaveOccurred())
		})
	})

	Context("with no dsl", func() {
		BeforeEach(func() {
			name = "application/foo"
		})

		It("produces an error", func() {
			Ω(mt).ShouldNot(BeNil())
			Ω(mt.Validate()).Should(HaveOccurred())
		})
	})

	Context("with attributes", func() {
		const attName = "att"

		BeforeEach(func() {
			name = "application/foo"
			dslFunc = func() {
				Attributes(func() {
					Attribute(attName)
				})
				View("default", func() { Attribute(attName) })
			}
		})

		It("sets the attributes", func() {
			Ω(mt).ShouldNot(BeNil())
			Ω(mt.Validate()).ShouldNot(HaveOccurred())
			Ω(mt.AttributeDefinition).ShouldNot(BeNil())
			Ω(mt.Type).Should(BeAssignableToTypeOf(Object{}))
			o := mt.Type.(Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(attName))
		})
	})

	Context("with a description", func() {
		const description = "desc"

		BeforeEach(func() {
			name = "application/foo"
			dslFunc = func() {
				Description(description)
				Attributes(func() {
					Attribute("attName")
				})
				View("default", func() { Attribute("attName") })
			}
		})

		It("sets the description", func() {
			Ω(mt).ShouldNot(BeNil())
			Ω(mt.Validate()).ShouldNot(HaveOccurred())
			Ω(mt.Description).Should(Equal(description))
		})
	})

	Context("with links", func() {
		const linkName = "link"
		var link1Name, link2Name string
		var link2View string
		var mt1, mt2 *MediaTypeDefinition

		BeforeEach(func() {
			name = "foo"
			link1Name = "l1"
			link2Name = "l2"
			link2View = "l2v"
			mt1 = NewMediaTypeDefinition("application/mt1", "application/mt1", func() {
				Attributes(func() {
					Attribute("foo")
				})
				View("default", func() {
					Attribute("foo")
				})
				View("link", func() {
					Attribute("foo")
				})
			})
			mt2 = NewMediaTypeDefinition("application/mt2", "application/mt2", func() {
				Attributes(func() {
					Attribute("foo")
				})
				View("l2v", func() {
					Attribute("foo")
				})
				View("default", func() {
					Attribute("foo")
				})
			})
			Design.MediaTypes = make(map[string]*MediaTypeDefinition)
			Design.MediaTypes["application/mt1"] = mt1
			Design.MediaTypes["application/mt2"] = mt2
			dslFunc = func() {
				Attributes(func() {
					Attributes(func() {
						Attribute(link1Name, mt1)
						Attribute(link2Name, mt2)
					})
					Links(func() {
						Link(link1Name)
						Link(link2Name, link2View)
					})
					View("default", func() {
						Attribute(link1Name)
						Attribute(link2Name)
					})
				})
			}
		})

		It("sets the links", func() {
			Ω(mt).ShouldNot(BeNil())
			Ω(dslengine.Errors).Should(BeEmpty())
			Ω(mt.Validate()).ShouldNot(HaveOccurred())
			Ω(mt.Links).ShouldNot(BeNil())
			Ω(mt.Links).Should(HaveLen(2))
			Ω(mt.Links).Should(HaveKey(link1Name))
			Ω(mt.Links[link1Name].Name).Should(Equal(link1Name))
			Ω(mt.Links[link1Name].View).Should(Equal("link"))
			Ω(mt.Links[link1Name].Parent).Should(Equal(mt))
			Ω(mt.Links[link2Name].Name).Should(Equal(link2Name))
			Ω(mt.Links[link2Name].View).Should(Equal(link2View))
			Ω(mt.Links[link2Name].Parent).Should(Equal(mt))
		})
	})

	Context("with views", func() {
		const viewName = "view"
		const viewAtt = "att"

		BeforeEach(func() {
			name = "application/foo"
			dslFunc = func() {
				Attributes(func() {
					Attribute(viewAtt)
				})
				View(viewName, func() {
					Attribute(viewAtt)
				})
				View("default", func() {
					Attribute(viewAtt)
				})
			}
		})

		It("sets the views", func() {
			Ω(mt).ShouldNot(BeNil())
			Ω(mt.Validate()).ShouldNot(HaveOccurred())
			Ω(mt.Views).ShouldNot(BeNil())
			Ω(mt.Views).Should(HaveLen(2))
			Ω(mt.Views).Should(HaveKey(viewName))
			v := mt.Views[viewName]
			Ω(v.Name).Should(Equal(viewName))
			Ω(v.Parent).Should(Equal(mt))
			Ω(v.AttributeDefinition).ShouldNot(BeNil())
			Ω(v.AttributeDefinition.Type).Should(BeAssignableToTypeOf(Object{}))
			o := v.AttributeDefinition.Type.(Object)
			Ω(o).Should(HaveLen(1))
			Ω(o).Should(HaveKey(viewAtt))
			Ω(o[viewAtt]).ShouldNot(BeNil())
			Ω(o[viewAtt].Type).Should(Equal(String))
		})
	})
})

var _ = Describe("Duplicate media types", func() {
	var mt *MediaTypeDefinition
	var duplicate *MediaTypeDefinition
	const id = "application/foo"
	const attName = "bar"
	var dslFunc = func() {
		Attributes(func() {
			Attribute(attName)
		})
		View("default", func() { Attribute(attName) })
	}

	BeforeEach(func() {
		InitDesign()
		dslengine.Errors = nil
		mt = MediaType(id, dslFunc)
		Ω(dslengine.Errors).ShouldNot(HaveOccurred())
		duplicate = MediaType(id, dslFunc)
	})

	It("produces an error", func() {
		Ω(dslengine.Errors).Should(HaveOccurred())
	})

	Context("with a response definition using the duplicate", func() {
		BeforeEach(func() {
			Resource("foo", func() {
				Action("show", func() {
					Routing(GET(""))
					Response(OK, func() {
						Media(duplicate)
					})
				})
			})
		})

		It("does not panic", func() {
			Ω(func() { dslengine.Run() }).ShouldNot(Panic())
		})
	})
})

var _ = Describe("CollectionOf", func() {
	Context("used on a global variable", func() {
		var col *MediaTypeDefinition
		BeforeEach(func() {
			InitDesign()
			mt := MediaType("application/vnd.example", func() { Attribute("id") })
			dslengine.Errors = nil
			col = CollectionOf(mt)
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
		})

		It("produces a media type", func() {
			Ω(col).ShouldNot(BeNil())
			Ω(col.Identifier).ShouldNot(BeEmpty())
			Ω(col.TypeName).ShouldNot(BeEmpty())
			Ω(Design.MediaTypes).Should(HaveKey(col.Identifier))
		})
	})

	Context("defined with the media type identifier", func() {
		var col *MediaTypeDefinition
		BeforeEach(func() {
			InitDesign()
			MediaType("application/vnd.example+json", func() { Attribute("id") })
			col = MediaType("application/vnd.parent+json", func() { Attribute("mt", CollectionOf("application/vnd.example")) })
		})

		JustBeforeEach(func() {
			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
		})

		It("produces a media type", func() {
			Ω(col).ShouldNot(BeNil())
			Ω(col.Identifier).Should(Equal("application/vnd.parent+json"))
			Ω(col.TypeName).Should(Equal("Parent"))
			Ω(col.Type).ShouldNot(BeNil())
			Ω(col.Type.ToObject()).ShouldNot(BeNil())
			Ω(col.Type.ToObject()).Should(HaveKey("mt"))
			mt := col.Type.ToObject()["mt"]
			Ω(mt.Type).ShouldNot(BeNil())
			Ω(mt.Type).Should(BeAssignableToTypeOf(&MediaTypeDefinition{}))
			Ω(mt.Type.Name()).Should(Equal("array"))
			et := mt.Type.ToArray().ElemType
			Ω(et).ShouldNot(BeNil())
			Ω(et.Type).Should(BeAssignableToTypeOf(&MediaTypeDefinition{}))
			Ω(et.Type.(*MediaTypeDefinition).Identifier).Should(Equal("application/vnd.example+json"))
		})
	})
})

var _ = Describe("Example", func() {
	Context("defined examples in a media type", func() {
		var mt *MediaTypeDefinition

		BeforeEach(func() {
			InitDesign()

			mt = MediaType("application/vnd.example+json", func() {
				Attribute("test1", String, "test1 desc", func() {
					Example("test1")
				})

				Attribute("test2", String, "test2 desc", func() {
					Example(None)
				})

				Attribute("test3", Integer, "test3 desc", func() {
					Minimum(1)
				})
			})
		})

		It("produces a media type with examples", func() {
			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
			Ω(mt).ShouldNot(BeNil())
			attr := mt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = mt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = mt.Type.ToObject()["test3"]
			Ω(attr.Example).Should(BeNumerically(">", 0))
		})

		It("produces media type examples from the linked media type", func() {
			pmt := MediaType("application/vnd.example.parent+json", func() {
				Attribute("test1", String, "test1 desc", func() {
					Example("test1")
				})

				Attribute("test2", String, "test2 desc", func() {
					Example(None)
				})

				Attribute("test3", Integer, "test3 desc", func() {
					Minimum(1)
				})

				Attribute("test4", mt, "test4 desc")
			})

			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
			Ω(mt).ShouldNot(BeNil())
			attr := mt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = mt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = mt.Type.ToObject()["test3"]
			Ω(attr.Example).Should(BeNumerically(">", 0))

			Ω(pmt).ShouldNot(BeNil())
			attr = pmt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = pmt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = pmt.Type.ToObject()["test3"]
			Ω(attr.Example).Should(BeNumerically(">", 0))
			attr = pmt.Type.ToObject()["test4"]
			Ω(attr.Example).ShouldNot(BeNil())
			attrChild, pass := attr.Example.(map[string]interface{})
			Ω(pass).Should(BeTrue())
			Ω(attrChild["test1"]).Should(Equal("test1"))
			Ω(attrChild["test2"]).Should(BeNil())
			Ω(attrChild["test3"]).Should(BeNumerically(">", 0))
		})

		It("produces media type examples from the linked media type collection with custom examples", func() {
			pmt := MediaType("application/vnd.example.parent+json", func() {
				Attribute("test1", String, "test1 desc", func() {
					Example("test1")
				})

				Attribute("test2", String, "test2 desc", func() {
					Example(None)
				})

				Attribute("test3", String, "test3 desc", func() {
					Pattern("^1$")
				})

				Attribute("test4", CollectionOf(mt), "test4 desc")
			})

			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
			Ω(mt).ShouldNot(BeNil())
			attr := mt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = mt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = mt.Type.ToObject()["test3"]
			Ω(attr.Example).Should(BeNumerically(">", 0))

			Ω(pmt).ShouldNot(BeNil())
			attr = pmt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = pmt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = pmt.Type.ToObject()["test3"]
			Ω(attr.Example).Should(Equal("1"))
			attr = pmt.Type.ToObject()["test4"]
			Ω(attr.Example).ShouldNot(BeNil())
			attrChildren, pass := attr.Example.([]interface{})
			Ω(pass).Should(BeTrue())
			Ω(attrChildren).Should(HaveLen(1))
			attrChild, pass := attrChildren[0].(map[string]interface{})
			Ω(pass).Should(BeTrue())
			Ω(attrChild["test1"]).Should(Equal("test1"))
			Ω(attrChild["test2"]).Should(BeNil())
			Ω(attrChild["test3"]).Should(BeNumerically(">", 0))
		})

		It("produces media type examples from the linked media type without custom examples", func() {
			cmt := MediaType("application/vnd.example.child+json", func() {
				Attribute("test1", String, "test1 desc")
			})

			pmt := MediaType("application/vnd.example.parent+json", func() {
				Attribute("test1", String, "test1 desc", func() {
					Example("test1")
				})

				Attribute("test2", String, "test2 desc", func() {
					Example(None)
				})

				Attribute("test3", cmt, "test3 desc")
			})

			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
			Ω(cmt).ShouldNot(BeNil())
			attr := cmt.Type.ToObject()["test1"]
			cexample := attr.Example
			Ω(cexample).ShouldNot(BeEmpty())

			Ω(pmt).ShouldNot(BeNil())
			attr = pmt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = pmt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = pmt.Type.ToObject()["test3"]
			Ω(attr.Example).ShouldNot(BeNil())
			attrChild, pass := attr.Example.(map[string]interface{})
			Ω(pass).Should(BeTrue())
			Ω(attrChild["test1"]).ShouldNot(Equal(cexample))
		})

		It("produces media type examples from the linked media type collection without custom examples", func() {
			cmt := MediaType("application/vnd.example.child+json", func() {
				Attribute("test1", String, "test1 desc")
			})

			pmt := MediaType("application/vnd.example.parent+json", func() {
				Attribute("test1", String, "test1 desc", func() {
					Example("test1")
				})

				Attribute("test2", String, "test2 desc", func() {
					Example(None)
				})

				Attribute("test3", CollectionOf(cmt), "test3 desc")
			})

			dslengine.Run()
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
			Ω(cmt).ShouldNot(BeNil())
			attr := cmt.Type.ToObject()["test1"]
			cexample := attr.Example
			Ω(cexample).ShouldNot(BeEmpty())

			Ω(pmt).ShouldNot(BeNil())
			attr = pmt.Type.ToObject()["test1"]
			Ω(attr.Example).Should(Equal("test1"))
			attr = pmt.Type.ToObject()["test2"]
			Ω(attr.Example).Should(BeNil())
			attr = pmt.Type.ToObject()["test3"]
			Ω(attr.Example).ShouldNot(BeNil())
			attrChildren, pass := attr.Example.([]interface{})
			Ω(pass).Should(BeTrue())
			Ω(len(attrChildren)).Should(BeNumerically(">", 0))
			attrChild, pass := attrChildren[0].(map[string]interface{})
			Ω(pass).Should(BeTrue())
			Ω(attrChild["test1"]).ShouldNot(Equal(cexample))
		})
	})
})
