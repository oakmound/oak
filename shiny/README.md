# Driver

This is a clone of github.com/oakmound/oak/v2/shiny, itself a fork of golang.org/exp/shiny.
The goal of this fork is to add additional window management functionality, and
focus the project down to just window management / common OS interfaces.

The goal of the clone here is to reduce iteration time of new features in oak. This
clone may be brought back to github.com/oakmound/oak/v2/shiny regularly, if we make
significant stable improvements.

## Long Term Plans

1. Standardize interfaces across OSes
2. Add new drivers for better performance, optional CGO
3. Add Fullscreen, screen movement, common screen options to all OSes
4. Add Window scaling types (bicubic, etc.) to all OSes
