package exec

import (
	"bytes"
	"log"
	"math"
	"net/http"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/davidbyttow/govips/v2/vips"
)

func Bounce(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var imageURL string
	var err error

	if len(args) > 1 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments")
		return
	}

	imageURL, err = getImageFromMessage(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	} else if imageURL == "" {
		imageURL = args[0]
		_, err = url.ParseRequestURI(imageURL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "You need to provide a valid image or URL")
			return
		}
	}

	imageName := getImageNameFromMessage(m)

	go func() {
		resp, err := http.Get(imageURL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to download image")
			return
		}
		defer resp.Body.Close()

		// simplified version of https://github.com/esmBot/esmBot/blob/master/natives/bounce.cc

		// Loads directly from the reader
		image, err := vips.NewImageFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer image.Close()

		// Convert to sRGB colorspace
		if image.Interpretation() != vips.InterpretationSRGB {
			if err := image.ToColorSpace(vips.InterpretationSRGB); err != nil {
				log.Fatal(err)
			}
		}

		// Add alpha channel if not present
		if !image.HasAlpha() {
			if err := image.BandJoinConst([]float64{255}); err != nil {
				log.Fatal(err)
			}
		}

		width := image.Width()
		pageHeight := image.Height()
		nPages := 15 // not handling multiple pages

		// Normalize image if too big
		maxSize := math.Max(float64(width), float64(pageHeight))
		if maxSize > 800 {
			image.Resize(800/maxSize, vips.KernelAuto)
			width = image.Width()
			pageHeight = image.Height()
		}

		mult := 3.14 / float64(nPages)
		halfHeight := pageHeight / 2

		// Create frames
		frames := make([]*vips.ImageRef, nPages)
		for i := 0; i < nPages; i++ {
			imgFrame, _ := image.Copy()
			defer imgFrame.Close()

			// https://www.desmos.com/calculator/sj0gdtzapk
			height := int(float64(halfHeight) * (-math.Sin(float64(i)*mult) + 1))

			imgFrame.Embed(0, height, width, pageHeight+halfHeight, vips.ExtendBackground)

			frames[i] = imgFrame
		}

		final := frames[0]
		defer final.Close()

		final.ArrayJoin(frames[1:], 1)
		final.SetInt("page-height", pageHeight+halfHeight)

		// Set animation delay (50ms per frame)
		delays := make([]int, 30)
		for i := range delays {
			delays[i] = 50
		}
		if err := final.SetPageDelay(delays); err != nil {
			log.Fatal(err)
		}

		// Export to GIF
		ep := vips.NewGifExportParams()
		ep.Quality = 95

		gifBytes, _, err := final.ExportGIF(ep)
		if err != nil {
			log.Fatal(err)
		}

		s.ChannelFileSend(m.ChannelID, imageName+".gif", bytes.NewReader(gifBytes))
	}()
}
