package particle

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"testing"
	"time"
)

var (
	pg = Generator{
		10, 5,
		2, 1,
		90, 3,
		0, 360,
		3, 0.5,
		50, 50,
		200,
		4, 3,
		4, 3,
		color.RGBA{100, 200, 200, 255},
		color.RGBA{0, 0, 0, 0},
		color.RGBA{0, 0, 0, 0},
		color.RGBA{0, 0, 0, 0},
		15, 9,
		Diamond,
	}
)

func BenchmarkParticles(b *testing.B) {
	curSeed := time.Now().UTC().UnixNano()
	rand.Seed(curSeed)

	pg.Rotation = pg.Rotation / 180 * math.Pi
	pg.RotationRand = pg.Rotation / 180 * math.Pi

	// Make a source
	ps := Source{
		Generator: pg,
		particles: make([]*Particle, 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		pg := ps.Generator

		newParticles := make([]*Particle, 0)

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

			newParticles = append(newParticles, &Particle{
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

	pg.Rotation = pg.Rotation / 180 * math.Pi
	pg.RotationRand = pg.Rotation / 180 * math.Pi

	// Make a source
	ps := Source{
		Generator: pg,
		particles: make([]*Particle, 0),
	}

	getRotateCh := make(chan *Particle)
	rotateChs := [channelCount]chan *Particle{}
	for i := 0; i < channelCount; i++ {
		rotateChs[i] = make(chan *Particle)
		go func(pg Generator, rotateCh chan *Particle, pCh chan *Particle) {
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
	newChs := [channelCount]chan bool{}
	for i := 0; i < channelCount; i++ {
		newChs[i] = make(chan bool)
		go func(pg Generator, newCh chan bool, pCh chan *Particle) {
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
		}(pg, newChs[i], getRotateCh)
	}

	doneCh := make(chan bool, 100)

	drawChs := [channelCount]chan *Particle{}
	for i := 0; i < channelCount; i++ {
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

	rotateAggregateCh := make(chan *Source)
	go func(rotateChs [channelCount]chan *Particle, aggCh chan *Source) {
		for {
			ps := <-aggCh
			for i, p := range ps.particles {
				rotateChs[i%channelCount] <- p
			}
		}
	}(rotateChs, rotateAggregateCh)

	newAggregateCh := make(chan int)
	go func(newChs [channelCount]chan bool, aggCh chan int) {
		for {
			newParticleCount := <-aggCh
			for i := 0; i < newParticleCount; i++ {
				newChs[i%channelCount] <- true
			}
		}
	}(newChs, newAggregateCh)

	drawAggregateCh := make(chan *Source)
	go func(drawChs [channelCount]chan *Particle, aggCh chan *Source) {
		for {
			ps := <-aggCh
			for i, p := range ps.particles {
				drawChs[i%channelCount] <- p
			}
		}
	}(drawChs, drawAggregateCh)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		rotateAggregateCh <- &ps

		pg := ps.Generator
		newParticleRand := roundFloat(floatFromSpread(pg.NewPerFrameRand))
		newParticleCount := int(pg.NewPerFrame) + newParticleRand

		newAggregateCh <- newParticleCount

		newParticles := make([]*Particle, 0)

		for i := 0; i < newParticleCount+len(ps.particles); i++ {
			next := <-getRotateCh
			if next != nil {
				newParticles = append(newParticles, next)
			}
		}
		ps.particles = newParticles

		drawAggregateCh <- &ps

		for range ps.particles {
			<-doneCh
		}
	}
}
