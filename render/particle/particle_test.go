package particle

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkParticles(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	pg := ParticleGenerator{
		11, 3,
		0, 0,
		6, 3,
		0, 360,
		3, 0.5,
		50, 50,
		200,
		0, 0,
		0, 0,
		color.RGBA{100, 200, 200, 255},
		color.RGBA{0, 0, 0, 0},
		color.RGBA{100, 200, 200, 200},
		color.RGBA{0, 0, 0, 0},
		15, 9,
		Square,
	}

	pg.Rotation = pg.Rotation / 180 * math.Pi
	pg.RotationRand = pg.Rotation / 180 * math.Pi

	// Make a source
	ps := ParticleSource{
		Generator: pg,
		particles: make([]Particle, 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		pg := ps.Generator

		newParticles := make([]Particle, 0)

		for _, p := range ps.particles {

			// Ignore dead particles
			if p.life > 0 {

				// Move towards doom
				p.life--

				// Be dragged down by the weight of the soul
				p.velX += pg.GravityX
				p.velY += pg.GravityY

				// Apply rotational acceleration
				if pg.Rotation != 0 && pg.RotationRand != 0 {
					magnitude := math.Abs(p.velX) + math.Abs(p.velY)
					angle := math.Atan2(p.velX, p.velY)
					angle += pg.Rotation + floatFromSpread(pg.RotationRand)
					p.velX = math.Sin(angle)
					p.velY = math.Cos(angle)
					magnitude = magnitude / (math.Abs(p.velX) + math.Abs(p.velY))
					p.velX = p.velX * magnitude
					p.velY = p.velY * magnitude
				}

				p.x += p.velX
				p.y += p.velY

				newParticles = append(newParticles, p)
			}
		}

		// Regularly create particles (up until max particles)
		newParticleRand := roundFloat(floatFromSpread(pg.NewPerFrameRand))
		newParticleCount := int(pg.NewPerFrame) + newParticleRand
		for i := 0; i < newParticleCount; i++ {

			angle := (pg.Angle + floatFromSpread(pg.AngleRand)) * math.Pi / 180.0
			speed := pg.Speed + floatFromSpread(pg.SpeedRand)
			startLife := pg.LifeSpan + floatFromSpread(pg.LifeSpanRand)

			newParticles = append(newParticles, Particle{
				x:          pg.X + floatFromSpread(pg.SpreadX),
				y:          pg.Y + floatFromSpread(pg.SpreadY),
				velX:       speed * math.Cos(angle) * -1,
				velY:       speed * math.Sin(angle) * -1,
				startColor: randColor(pg.StartColor, pg.StartColorRand),
				endColor:   randColor(pg.EndColor, pg.EndColorRand),
				life:       startLife,
				totalLife:  startLife,
				size:       pg.Size + intFromSpread(pg.SizeRand),
			})
		}

		ps.particles = newParticles

		for _, p := range ps.particles {

			r, g, b, a := p.startColor.RGBA()
			r2, g2, b2, a2 := p.endColor.RGBA()
			progress := p.life / p.totalLife
			c := color.RGBA64{
				uint16OnScale(r, r2, progress),
				uint16OnScale(g, g2, progress),
				uint16OnScale(b, b2, progress),
				uint16OnScale(a, a2, progress),
			}

			img := image.NewRGBA64(image.Rect(0, 0, p.size, p.size))

			for i := 0; i < p.size; i++ {
				for j := 0; j < p.size; j++ {
					if ps.Generator.Shape(i, j, p.size) {
						img.SetRGBA64(i, j, c)
					}
				}
			}
		}
	}
}

func BenchmarkParticlesParallel(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	pg := ParticleGenerator{
		11, 3,
		0, 0,
		6, 3,
		0, 360,
		3, 0.5,
		50, 50,
		200,
		0, 0,
		0, 0,
		color.RGBA{100, 200, 200, 255},
		color.RGBA{0, 0, 0, 0},
		color.RGBA{100, 200, 200, 200},
		color.RGBA{0, 0, 0, 0},
		15, 9,
		Square,
	}

	pg.Rotation = pg.Rotation / 180 * math.Pi
	pg.RotationRand = pg.Rotation / 180 * math.Pi

	// Make a source
	ps := ParticleSource{
		Generator: pg,
		particles: make([]Particle, 0),
	}

	getRotateCh := make(chan *Particle, 100)
	getNewCh := make(chan *Particle, 100)
	rotateChs := [8]chan *Particle{}
	for i := 0; i < 8; i++ {
		rotateChs[i] = make(chan *Particle)
		go func(pg ParticleGenerator, rotateCh chan *Particle, pCh chan *Particle) {
			for {
				p := <-rotateCh
				// Ignore dead particles
				if p.life > 0 {

					// Move towards doom
					p.life--

					// Be dragged down by the weight of the soul
					p.velX += pg.GravityX
					p.velY += pg.GravityY

					// Apply rotational acceleration
					if pg.Rotation != 0 && pg.RotationRand != 0 {
						magnitude := math.Abs(p.velX) + math.Abs(p.velY)
						angle := math.Atan2(p.velX, p.velY)
						angle += pg.Rotation + floatFromSpread(pg.RotationRand)
						p.velX = math.Sin(angle)
						p.velY = math.Cos(angle)
						magnitude = magnitude / (math.Abs(p.velX) + math.Abs(p.velY))
						p.velX = p.velX * magnitude
						p.velY = p.velY * magnitude
					}

					p.x += p.velX
					p.y += p.velY

					pCh <- p
				} else {
					pCh <- nil
				}
			}
		}(pg, rotateChs[i], getRotateCh)
	}
	newChs := [8]chan bool{}
	for i := 0; i < 8; i++ {
		newChs[i] = make(chan bool)
		go func(pg ParticleGenerator, newCh chan bool, pCh chan *Particle) {
			for {
				<-newCh
				angle := (pg.Angle + floatFromSpread(pg.AngleRand)) * math.Pi / 180.0
				speed := pg.Speed + floatFromSpread(pg.SpeedRand)
				startLife := pg.LifeSpan + floatFromSpread(pg.LifeSpanRand)
				pCh <- &Particle{
					x:          pg.X + floatFromSpread(pg.SpreadX),
					y:          pg.Y + floatFromSpread(pg.SpreadY),
					velX:       speed * math.Cos(angle) * -1,
					velY:       speed * math.Sin(angle) * -1,
					startColor: randColor(pg.StartColor, pg.StartColorRand),
					endColor:   randColor(pg.EndColor, pg.EndColorRand),
					life:       startLife,
					totalLife:  startLife,
					size:       pg.Size + intFromSpread(pg.SizeRand),
				}
			}
		}(pg, newChs[i], getNewCh)
	}

	doneCh := make(chan bool, 100)

	drawChs := [8]chan *Particle{}
	for i := 0; i < 8; i++ {
		drawChs[i] = make(chan *Particle)
		go func(drawCh chan *Particle, doneCh chan bool) {
			for {
				p := <-drawCh
				r, g, b, a := p.startColor.RGBA()
				r2, g2, b2, a2 := p.endColor.RGBA()
				progress := p.life / p.totalLife
				c := color.RGBA64{
					uint16OnScale(r, r2, progress),
					uint16OnScale(g, g2, progress),
					uint16OnScale(b, b2, progress),
					uint16OnScale(a, a2, progress),
				}

				img := image.NewRGBA64(image.Rect(0, 0, p.size, p.size))

				for i := 0; i < p.size; i++ {
					for j := 0; j < p.size; j++ {
						if ps.Generator.Shape(i, j, p.size) {
							img.SetRGBA64(i, j, c)
						}
					}
				}
				doneCh <- true
			}
		}(drawChs[i], doneCh)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pg := ps.Generator

		for i, p := range ps.particles {
			rotateChs[i%8] <- &p
		}

		newParticleRand := roundFloat(floatFromSpread(pg.NewPerFrameRand))
		newParticleCount := int(pg.NewPerFrame) + newParticleRand

		for i := 0; i < newParticleCount; i++ {
			newChs[i%8] <- true
		}

		newParticles := make([]Particle, 0)

		for i := 0; i < newParticleCount; i++ {
			newParticles = append(newParticles, *<-getNewCh)
		}
		for i := 0; i < len(ps.particles); i++ {
			next := <-getRotateCh
			if next != nil {
				newParticles = append(newParticles, *next)
			}
		}
		ps.particles = newParticles

		for i, p := range ps.particles {
			drawChs[i%8] <- &p
		}

		for range ps.particles {
			<-doneCh
		}
	}
}
