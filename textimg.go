package textimg

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	//"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"strings"
	//"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type textimg struct {
	dstimg   *image.RGBA
	imgbg    *image.Uniform
	ttfont   *truetype.Font
	fontSize float64
	hasDraw  bool
}

// rgba 目标图片( image.NewRGBA(image.Rect(0, 0, 1000, 50)) )
// bg 目标图片背景色 ( image.White )
func New(rgba *image.RGBA, bg *image.Uniform) *textimg {
	if rgba == nil {
		rgba = image.NewRGBA(image.Rect(0, 0, 320, 240))
	}
	if bg == nil {
		bg = image.Transparent
	}

	return &textimg{dstimg: rgba, imgbg: bg}
}

// 字体从外部传入
func (t *textimg) SetFont(ttfont *truetype.Font) error {
	if ttfont == nil {
		return errors.New("font net found")
	}

	t.ttfont = ttfont
	return nil
}

// 字体从本地新建
func (t *textimg) SetFontFromPath(fontpath string) error {
	//fontPath := "/home/mingqing/Documents/codes/zujuan/src/siming/public/siming/question/font/simsun.ttc"
	fontBytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	t.ttfont = f
	return nil
}

func (t *textimg) SetFontSize(size float64) error {
	if size == 0 {
		size = 14
	}

	t.fontSize = size

	return nil
}

// 返回png格式
func (t *textimg) PNG() *bytes.Buffer {
	temp := bytes.NewBuffer(make([]byte, 0))
	png.Encode(temp, t.dstimg)
	return temp
}

// 生存png格式图片
func (t *textimg) TextToPNG(fg *image.Uniform, textline []string) *bytes.Buffer {
	t.DrawDstimg(fg, textline)
	return t.PNG()
}

// 写入数据到目标图片中
func (t *textimg) DrawDstimg(fg *image.Uniform, textline []string) *image.RGBA {
	pt := freetype.Pt(0, int(t.fontSize))
	c := t.createImg(fg)
	t.drawTextline(c, pt, textline)

	return t.dstimg
}

// 添加一个文字图片
func (t *textimg) AddImage(pt image.Point, rgba *image.RGBA) *image.RGBA {
	draw.Draw(t.dstimg, t.dstimg.Bounds(), rgba, rgba.Bounds().Min.Add(pt), draw.Over)
	return t.dstimg
}

// 添加一个图片来自html
// <img src="data:image/png;base64......">
func (t *textimg) AddImageFromHtmlSrcBase64(pt image.Point, imgbase64 string) {
	tt := strings.Replace(imgbase64, "data:image/png;base64,", "", -1)
	ddd, _ := base64.StdEncoding.DecodeString(tt)
	imgcache := bytes.NewBuffer(ddd)
	img, _, _ := image.Decode(imgcache)
	draw.Draw(t.dstimg, t.dstimg.Bounds(), img, t.dstimg.Bounds().Min.Add(pt), draw.Src)
	//draw.DrawMask(t.dstimg, t.dstimg.Bounds(), img, t.dstimg.Bounds().Min.Add(pt),
}
func (t *textimg) AddImageFromHtmlSrcBase64WH(pt image.Point, imgbase64 string, width, height uint) {
	tt := strings.Replace(imgbase64, "data:image/png;base64,", "", -1)
	ddd, _ := base64.StdEncoding.DecodeString(tt)
	imgcache := bytes.NewBuffer(ddd)
	img, _, _ := image.Decode(imgcache)
	m := resize.Resize(width, height, img, resize.Lanczos3)
	draw.Draw(t.dstimg, t.dstimg.Bounds(), m, t.dstimg.Bounds().Min.Add(pt), draw.Over)
}

func (t *textimg) AddTextline(pt image.Point, text string) {

}

// fg 目标图片中文字颜色( image.Black )
func (t *textimg) createImg(fg *image.Uniform) *freetype.Context {
	if fg == nil {
		fg = image.Black
	}

	draw.Draw(t.dstimg, t.dstimg.Bounds(), t.imgbg, image.ZP, draw.Over)
	c := freetype.NewContext()
	// fontSize and dpi are used to calculate scale 26.6 fixed point units in 1 em
	c.SetDPI(72)
	c.SetFont(t.ttfont)
	c.SetFontSize(t.fontSize)
	c.SetClip(t.dstimg.Bounds())
	c.SetDst(t.dstimg)
	c.SetSrc(fg)
	c.SetHinting(font.HintingFull)
	return c
	//c.DrawString(t.body, pt)
	//buf := bytes.NewBuffer(make([]byte, 0))
}

// 文字写入图片
func (t *textimg) drawTextline(ftc *freetype.Context, pt fixed.Point26_6, textlist []string) *freetype.Context {
	//pt := freetype.Pt(0, t.fontSize)
	for _, s := range textlist {
		ftc.DrawString(s, pt)
		pt.Y += ftc.PointToFixed(t.fontSize * 3)
	}

	return ftc
}
